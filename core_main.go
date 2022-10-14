package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"gitee.com/uni-minds/medical-sys/controller"
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/logger"
	"gitee.com/uni-minds/medical-sys/manager"
	"gitee.com/uni-minds/medical-sys/module"
	"gitee.com/uni-minds/medical-sys/module_router"
	"gitee.com/uni-minds/medical-sys/module_rpc"
	"gitee.com/uni-minds/medical-sys/module_rtsp"
	"gitee.com/uni-minds/utils/tools"
	figure "github.com/common-nighthawk/go-figure"
	"github.com/kardianos/service"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

type program struct {
	rtspPort   int
	rtspServer *module_rtsp.Server

	httpPort   int
	httpServer *http.Server

	rpcPort   int
	rpcServer *module_rpc.Server

	logger *logger.Logger
}

var (
	_BUILD_VER_  = "3.0-dev"
	_BUILD_TIME_ = "2022-09-01 00:01:02"
	_BUILD_REV_  = "90b0388"
)

const DefCfgFile = "/data/medisys/config.yaml"

func (p *program) StartRTSP() (err error) {
	p.rtspServer = module_rtsp.CreateServer(p.rtspPort)
	if p.rtspServer == nil {
		err = fmt.Errorf("RTSP Server Not Found")
		return
	}
	sport := ""
	if p.rtspPort != 554 {
		sport = fmt.Sprintf(":%d", p.rtspPort)
	}
	link := fmt.Sprintf("rtsp://%s%s", tools.LocalIP(), sport)
	p.logger.Println("rtsp server start -->", link)

	go func() {
		if err = p.rtspServer.Start(); err != nil {
			p.logger.Println("start rtsp server error", err)
		}
		p.logger.Println("rtsp server end")
	}()
	return
}

func (p *program) StopRTSP() (err error) {
	if p.rtspServer == nil {
		err = fmt.Errorf("RTSP Server Not Found")
		return
	}
	p.rtspServer.Stop()
	return
}

func (p *program) StartRouter() (err error) {
	p.httpServer = module_router.CreateServer(p.httpPort)
	if p.httpServer == nil {
		return errors.New("HTTP Server Not Found")
	}

	sport := ""
	if p.httpPort != 80 {
		sport = fmt.Sprintf(":%d", p.httpPort)
	}

	link := fmt.Sprintf("http://%s%s", tools.LocalIP(), sport)
	p.logger.Println("http server start -->", link)
	go func() {
		if err = p.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("start http server error", err)
		}
		p.logger.Println("http server end")
	}()
	return
}

func (p *program) StopRouter() (err error) {
	if p.httpServer == nil {
		err = fmt.Errorf("HTTP Server Not Found")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = p.httpServer.Shutdown(ctx); err != nil {
		return
	}
	return
}

func (p *program) StartRPC() (err error) {
	p.rpcServer = module_rpc.CreateServer(p.rpcPort)
	if p.rpcServer == nil {
		return errors.New("RPC Server Not Found")
	}

	sport := fmt.Sprintf(":%d", p.rpcPort)

	link := fmt.Sprintf("rpc://%s%s", tools.LocalIP(), sport)
	p.logger.Println("rpc  server start -->", link)
	p.rpcServer.Start()
	return nil
}

func (p *program) StopRPC() (err error) {
	if p.rpcServer == nil {
		err = fmt.Errorf("RPC Server Not Found")
		return
	}
	p.rpcServer.Stop()
	return
}

func (p *program) Start(s service.Service) (err error) {
	logger := p.logger
	if tools.IsPortInUse(p.rtspPort) {
		return fmt.Errorf("RTSP port[%d] In Use", p.rtspPort)
	}
	if tools.IsPortInUse(p.httpPort) {
		return fmt.Errorf("HTTP port[%d] In Use", p.httpPort)
	}
	if tools.IsPortInUse(p.rpcPort) {
		return fmt.Errorf("RPC port[%d] In Use", p.rtspPort)
	}

	if p.httpPort > 0 {
		err = p.StartRouter()
		if err != nil {
			logger.Error(fmt.Sprintf("Error start http server: %s", err.Error()))
			return err
		}
	} else {
		logger.Printf("Ignore start http server")
	}

	if p.rpcPort > 0 {
		err = p.StartRPC()
		if err != nil {
			logger.Error(fmt.Sprintf("Error start rpc server: %s", err.Error()))
			return err
		}
	} else {
		logger.Printf("Ignore start rpc server")
	}

	if p.rtspPort > 0 {
		err = p.StartRTSP()
		if err != nil {
			logger.Error(fmt.Sprintf("Error start rtsp server: %s", err.Error()))
			return err
		}
	} else {
		logger.Printf("Ignore start rtsp server")
	}

	return nil
}

func (p *program) Stop(s service.Service) (err error) {
	p.logger.Warn("Stop")
	return nil
}

func main() {
	var argVerbose, argDebug, argPrintScreen bool
	var argCfg, logFile string
	var appConfig global.AppSettings

	flag.BoolVar(&argDebug, "d", false, "Debug mode")
	flag.BoolVar(&argVerbose, "v", false, "Verbose")
	flag.BoolVar(&argPrintScreen, "p", false, "Print log to screen")
	flag.StringVar(&argCfg, "c", "", "Configure file")
	flag.Parse()

	if argCfg != "" {
		appConfig = global.Init(argCfg)
	} else if _, err := os.Stat(filepath.Dir(DefCfgFile)); err == nil {
		appConfig = global.Init(DefCfgFile)
	} else {
		appConfig = global.Init("./config.yaml")
	}

	global.FlagSetDebug(argDebug)
	global.FlagSetVerbose(argVerbose)

	global.SetBuildTime(_BUILD_TIME_)
	global.SetGitCommit(_BUILD_REV_)
	global.SetVersion(_BUILD_VER_)

	if !argPrintScreen {
		logFile = path.Join(appConfig.Paths.Log, time.Now().Format("medisys_20060102T150405.log"))
	}

	if err := logger.Init(logFile, argVerbose); err != nil {
		panic(err)
	} else if err = manager.Init(); err != nil {
		panic(err)
	} else if err = database.Init(); err != nil {
		panic(err)
	} else if err = module.Init(); err != nil {
		panic(err)
	} else if err = controller.Init(); err != nil {
		panic(err)
	}

	p := &program{
		rtspPort: appConfig.Ports.RTSP,
		httpPort: appConfig.Ports.HTTP,
		rpcPort:  appConfig.Ports.RPC,
		logger:   logger.NewLogger("MAIN"),
	}

	s, err := service.New(p, &service.Config{
		Name:        "Medisys",
		DisplayName: "Medisys",
		Description: "Medical System Service",
	})

	if err != nil {
		p.logger.Error(err.Error())
	}

	figure.NewFigure("Medisys", "", false).Print()
	fmt.Printf("Beihang University Medical System %s\n", global.GetVersionString())

	if err = s.Run(); err != nil {
		panic(err)
	}
}
