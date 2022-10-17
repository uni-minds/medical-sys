package module_rpc

import (
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/logger"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Server struct {
	Instance   *rpc.Server
	ListenPort int
	chanStop   chan int
	Logger     *logger.Logger
}

func (this *RpcFunc) RunCmd(req RpcRequest, res *RpcResponse) (err error) {
	log.Debug(fmt.Sprintf("RPC_CALL: %s %v", req.Command, req.Param))

	params := make([]string, 0)
	for _, p := range req.Param {
		if p != "" {
			params = append(params, p)
		}
	}

	switch req.Command {
	case "GetVersion":
		res.Error = ""
		res.Result = global.GetVersionString()

	case "test_e":
		res.Error = "test_ok"
		res.Result = "test_result"

	case "test_ok":
		res.Error = ""
		res.Result = "test_ok"

	case "group":
		res.Result, err = this.ParseGroup(params)

	case "user":
		res.Result, err = this.ParseUser(params)

	case "media":
		res.Result, err = this.ParseMedia(params)

	case "stream":
		res.Result, err = this.ParseStream(params)

	case "pacs":
		res.Result, err = this.ParsePacs(params)

	case "genjson":
		res.Result, err = this.ParseGenJson(params)

	case "progress":
		res.Result = fmt.Sprintf("1:标注中\n2:标注完成\n3:审阅中\n4:审阅完成，拒绝\n5:标注修改中\n6:标注完成修改，提交审阅\n7:审阅接受，最终状态")

	case "cmds", "help":
		res.Result = fmt.Sprintf("commands: user | group | media | stream | label | import | pacs | progress | genjson")

	default:
		res.Result = fmt.Sprintf("unsupported: %s, use help to get support commands", req.Command)

	}

	if err != nil {
		res.Error = err.Error()
	}
	return nil
}

func CreateServer(port int) *Server {
	var srv Server

	instance := rpc.NewServer()
	instance.Register(new(RpcFunc))

	srv.Instance = instance
	srv.ListenPort = port
	srv.chanStop = make(chan int, 0)
	srv.Logger = logger.NewLogger("RPC")
	if log == nil {
		log = logger.NewLogger("RPCF")
	}
	return &srv
}

func (srv *Server) Start() {
	go srv.Core()
}

func (srv *Server) Stop() {
	srv.Logger.Log("i", "rpc_server sending stop signal")
	go func() {
		srv.chanStop <- 1
		srv.Logger.Log("i", "rpc_server stop signal sent")
	}()
}

func (srv *Server) Core() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.ListenPort))
	if err != nil {
		srv.Logger.Error(err.Error())
		return
	}

	for flagRun := true; flagRun; {
		select {
		case <-srv.chanStop:
			srv.Logger.Warn("rpc_server recv stop signal")
			flagRun = false
			continue

		default:
			conn, err := lis.Accept()
			if err != nil {
				continue
			}

			go func(conn net.Conn) {
				srv.Logger.Printf("rpc_connect")
				srv.Instance.ServeCodec(jsonrpc.NewServerCodec(conn))
			}(conn)
		}
	}

	srv.Logger.Warn("rpc_server stopped")
}
