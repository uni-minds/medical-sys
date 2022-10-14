package main

import (
	"bufio"
	"flag"
	"fmt"
	"gitee.com/uni-minds/medical-sys/module_rpc"
	"github.com/fatih/color"
	"net/rpc/jsonrpc"
	"os"
	"strings"
)

var rpcServer = "127.0.0.1:8096"

func main() {
	flag.StringVar(&rpcServer, "s", rpcServer, "rpc server")
	flag.Parse()

	var resp module_rpc.RpcResponse

	cRed := color.New(color.FgRed)
	cYellow := color.New(color.FgYellow)
	cGreen := color.New(color.FgGreen)

	flagKeepRun := true

	for flagKeepRun {
		conn, err := jsonrpc.Dial("tcp", rpcServer)
		if err != nil {
			panic(err.Error())
		}

		// 握手
		if err = conn.Call("RpcFunc.RunCmd", module_rpc.RpcRequest{Command: "GetVersion"}, &resp); err != nil {
			cRed.Printf("E_RPC_CALL: %s\n", err.Error())
			os.Exit(-1)
		} else if resp.Error != "" {
			cYellow.Printf("E_ShakeHands: %s\n", resp.Error)
			os.Exit(-1)
		} else {
			cGreen.Println("RPC_CONN_ESTABLISH server version:", resp.Result)
		}

		for {
			fmt.Printf("module_rpc> ")

			//if _, err = fmt.Scanln(&input); err != nil {
			//	cRed.Println(err.Error())
			//	continue
			//}

			scanner := bufio.NewReader(os.Stdin)
			input, _ := scanner.ReadString('\n')
			input = strings.TrimSpace(input)

			p := strings.Split(input, " ")
			if len(p) == 0 {
				continue

			} else if p[0] == "exit" {
				flagKeepRun = false
				break

			} else {
				req := module_rpc.RpcRequest{
					Command: p[0],
					Param:   p[1:],
				}

				if err = conn.Call("RpcFunc.RunCmd", req, &resp); err != nil {
					cRed.Println("E_RPC_CALL:", err.Error())
					break

				} else if resp.Error != "" {
					cYellow.Println(resp.Error)
					if resp.Result != "" {
						fmt.Println(resp.Result)
					}

				} else {
					cGreen.Println(resp.Result)
				}
			}
		}

		if err = conn.Close(); err != nil {
			cRed.Println("E_RPC_CONN_DOWN:", err.Error())
		} else {
			fmt.Println("RPC_CONN_STOPED")
		}
	}

	cGreen.Println("RPC_CLIENT_EXIT")
}
