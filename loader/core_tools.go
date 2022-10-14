package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"net/rpc/jsonrpc"
	"os"
	"strings"
)

func main() {
	conn, err := jsonrpc.Dial("tcp", "127.0.0.1:8096")
	if err != nil {
		fmt.Println("dail:", err.Error())
		os.Exit(1)
	}

	var res RpcResponse
	var input string

	for {
		fmt.Printf("rpc> ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = strings.TrimSpace(scanner.Text())
		p := strings.Split(input, " ")
		if len(p) == 0 {
			continue

		} else if p[0] == "exit" {
			os.Exit(0)

		} else {
			req := RpcRequest{
				Command: p[0],
				Param:   p[1:],
			}
			red := color.New(color.FgRed).PrintfFunc()
			yellow := color.New(color.FgYellow).PrintfFunc()

			if err = conn.Call("RpcFunc.RunCmd", req, &res); err != nil {
				red("E_RPC_CALL:\n")
				red("%s", err.Error())

			} else if res.Error != "" {
				red("E_RPC_RUN:\n")
				yellow("%s", res.Error)
			} else {
				red("RPC_RUN:\n")
				fmt.Println(res.Result)
			}
		}
	}
}
