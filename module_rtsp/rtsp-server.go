package module_rtsp

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/logger"
	"gitee.com/uni-minds/utils/tools"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	Logger *logger.Logger

	TCPListener    *net.TCPListener
	TCPPort        int
	Stoped         bool
	pushers        map[string]*Pusher // Path <-> Pusher
	pushersLock    sync.RWMutex
	addPusherCh    chan *Pusher
	removePusherCh chan *Pusher
}

type FfCmdInfo struct {
	Cmd    *exec.Cmd
	Folder string
}

var Instance *Server

func CreateServer(port int) *Server {
	Instance = &Server{
		Stoped:         true,
		TCPPort:        port,
		pushers:        make(map[string]*Pusher),
		addPusherCh:    make(chan *Pusher),
		removePusherCh: make(chan *Pusher),
	}
	Instance.Logger = logger.NewLogger("RTSP")
	return Instance
}

func (server *Server) Start() (err error) {
	var (
		logger   = server.Logger
		addr     *net.TCPAddr
		listener *net.TCPListener
	)
	if addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", server.TCPPort)); err != nil {
		return
	}
	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		return
	}

	localRecord := global.GetRtspSettings().SaveStreamToLocal
	ffmpeg := global.GetRtspSettings().FFmpegPath
	uploadRoot := global.GetRtspSettings().UploadRoot
	tsDurationSecond := global.GetRtspSettings().TsDurationSecond
	SaveStreamToLocal := false
	if localRecord && (len(ffmpeg) > 0) && len(uploadRoot) > 0 {
		err = tools.EnsureDir(uploadRoot)
		if err != nil {
			logger.Printf("Create uploadRoot[%s] err:%s.", uploadRoot, err.Error())
		} else {
			SaveStreamToLocal = true
		}
	}
	go func() { // save to local.
		pusher2ffmpegMap := make(map[*Pusher]FfCmdInfo)
		if SaveStreamToLocal {
			logger.Printf("rtsp local save enabled")
		}
		var pusher *Pusher
		addChnOk := true
		removeChnOk := true
		for addChnOk || removeChnOk {
			select {
			case pusher, addChnOk = <-server.addPusherCh:
				if SaveStreamToLocal {
					if addChnOk {
						dir := path.Join(uploadRoot, pusher.Path(), time.Now().Format("20060102150405"))
						if err = tools.EnsureDir(dir); err != nil {
							logger.Error(fmt.Sprintf("EnsureDir:[%s] err:%v.", dir, err))
							continue
						}
						m3u8path := path.Join(dir, fmt.Sprintf("video.m3u8"))
						port := pusher.Server().TCPPort
						rtsp := fmt.Sprintf("rtsp://localhost:%d%s", port, pusher.Path())
						params := []string{"-fflags", "genpts", "-rtsp_transport", "tcp", "-i", rtsp, "-hls_time", strconv.Itoa(tsDurationSecond), "-hls_list_size", "0", m3u8path}
						paramStr := global.GetRtspSettings().FFmpegEncoder
						if paramStr != "" && paramStr != "default" {
							paramsOfThisPath := strings.Split(paramStr, " ")
							params = append(params[:6], append(paramsOfThisPath, params[6:]...)...)
						}
						// ffmpeg -i ~/Downloads/720p.mp4 -s 640x360 -g 15 -c:a aac -hls_time 5 -hls_list_size 0 record.m3u8
						cmd := exec.Command(ffmpeg, params...)
						f, err := os.OpenFile(path.Join(dir, "log.txt"), os.O_RDWR|os.O_CREATE, 0755)
						if err == nil {
							cmd.Stdout = f
							cmd.Stderr = f
						}
						err = cmd.Start()
						if err != nil {
							logger.Error(fmt.Sprintf("Start ffmpeg err: %v", err.Error()))
						}
						pusher2ffmpegMap[pusher] = FfCmdInfo{
							Cmd:    cmd,
							Folder: dir,
						}
						logger.Printf("add ffmpeg [%v] to pull stream from pusher[%v]", cmd, pusher.Path())
					} else {
						logger.Println("addPusherChan closed")
					}
				}

			case pusher, removeChnOk = <-server.removePusherCh:
				if SaveStreamToLocal {
					if removeChnOk {
						info := pusher2ffmpegMap[pusher]
						cmd := info.Cmd
						proc := cmd.Process
						if proc != nil {
							logger.Printf("prepare to SIGTERM to proc[%d]", proc.Pid)
							proc.Signal(syscall.SIGTERM)
							proc.Wait()

							URI := strings.TrimLeft(pusher.Path(), "/")

							err := RtspSaveToDatabase(URI, info.Folder)
							if err != nil {
								logger.Error(fmt.Sprintf("SaveToDB: %s", err.Error()))
							}

							logger.Printf("proc[%d] terminate.", proc.Pid)
						}
						delete(pusher2ffmpegMap, pusher)
						logger.Printf("delete ffmpeg from pull stream from pusher[%s]", pusher.Path())
					} else {
						for _, info := range pusher2ffmpegMap {
							cmd := info.Cmd
							proc := cmd.Process
							if proc != nil {
								logger.Printf("prepare to SIGTERM to proc[%d]", proc.Pid)
								proc.Signal(syscall.SIGTERM)
							}
						}
						pusher2ffmpegMap = make(map[*Pusher]FfCmdInfo)
						logger.Println("removePusherChan closed")
					}
				}
			}
		}
	}()

	server.Stoped = false
	server.TCPListener = listener
	logger.Printf("server start on port: %d", server.TCPPort)
	networkBuffer := global.GetRtspSettings().NetworkBuffer
	for !server.Stoped {
		var conn net.Conn
		if conn, err = server.TCPListener.Accept(); err != nil {
			logger.Error(err.Error())
			continue
		}
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			if err = tcpConn.SetReadBuffer(networkBuffer); err != nil {
				logger.Error(fmt.Sprintf("rtsp server conn set read buffer error, %v", err))
			}
			if err = tcpConn.SetWriteBuffer(networkBuffer); err != nil {
				logger.Error(fmt.Sprintf("rtsp server conn set write buffer error, %v", err))
			}
		}

		session := NewSession(server, conn)
		go session.Start()
	}
	return
}

func (server *Server) Stop() {
	server.Logger.Printf("rtsp server stop on %d", server.TCPPort)
	server.Stoped = true
	if server.TCPListener != nil {
		server.TCPListener.Close()
		server.TCPListener = nil
	}
	server.pushersLock.Lock()
	server.pushers = make(map[string]*Pusher)
	server.pushersLock.Unlock()

	close(server.addPusherCh)
	close(server.removePusherCh)
}

func (server *Server) AddPusher(pusher *Pusher) bool {
	added := false
	server.pushersLock.Lock()
	_, ok := server.pushers[pusher.Path()]
	if !ok {
		server.pushers[pusher.Path()] = pusher
		server.Logger.Printf("%v start, now pusher size[%d]", pusher, len(server.pushers))
		added = true
	} else {
		added = false
	}
	server.pushersLock.Unlock()
	if added {
		go pusher.Start()
		server.addPusherCh <- pusher
	}
	return added
}

func (server *Server) TryAttachToPusher(session *Session) (int, *Pusher) {
	server.pushersLock.Lock()
	attached := 0
	var pusher *Pusher = nil
	if _pusher, ok := server.pushers[session.Path]; ok {
		if _pusher.RebindSession(session) {
			session.Logger.Printf("Attached to a pusher")
			attached = 1
			pusher = _pusher
		} else {
			attached = -1
		}
	}
	server.pushersLock.Unlock()
	return attached, pusher
}

func (server *Server) RemovePusher(pusher *Pusher) {
	removed := false
	server.pushersLock.Lock()
	if _pusher, ok := server.pushers[pusher.Path()]; ok && pusher.ID() == _pusher.ID() {
		delete(server.pushers, pusher.Path())
		server.Logger.Printf("%v end, now pusher size[%d]", pusher, len(server.pushers))
		removed = true
	}
	server.pushersLock.Unlock()
	if removed {
		server.removePusherCh <- pusher
	}
}

func (server *Server) GetPusher(path string) (pusher *Pusher) {
	server.pushersLock.RLock()
	pusher = server.pushers[path]
	server.pushersLock.RUnlock()
	return
}

func (server *Server) GetPushers() (pushers map[string]*Pusher) {
	pushers = make(map[string]*Pusher)
	server.pushersLock.RLock()
	for k, v := range server.pushers {
		pushers[k] = v
	}
	server.pushersLock.RUnlock()
	return
}

func (server *Server) GetPusherSize() (size int) {
	server.pushersLock.RLock()
	size = len(server.pushers)
	server.pushersLock.RUnlock()
	return
}
