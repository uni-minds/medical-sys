package rpc

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func (this *RpcFunc) RunCmd(req RpcRequest, res *RpcResponse) (err error) {
	fmt.Println("Rpc func1:", req.Command)
	fmt.Println("Rpc param:", req.Param)

	params := make([]string, 0)
	for _, p := range req.Param {
		if p != "" {
			params = append(params, p)
		}
	}

	switch req.Command {
	case "test_e":
		res.Error = "test_ok"
		res.Result = "test_result"

	case "test_ok":
		res.Error = ""
		res.Result = "test_ok"

	case "group":
		res.Result, err = parseGroup(params)

	case "user":
		res.Result, err = parseUser(params)

	case "media":
		res.Result, err = parseMedia(params)

	case "label":
		res.Result, err = parseLabel(params)

	case "import":
		res.Result, err = parseImport(params)

	case "pacs":
		res.Result, err = parsePacs(params)

	case "genjson":
		res.Result, err = parseGenJson(params)

	case "progress":
		res.Result = fmt.Sprintf("1:标注中\n2:标注完成\n3:审阅中\n4:审阅完成，拒绝\n5:标注修改中\n6:标注完成修改，提交审阅\n7:审阅接受，最终状态")

	default:
		res.Result = fmt.Sprintf("unsupported command: %s\nsupport user | group | media | label | import | progress", req.Command)

	}

	if err != nil {
		res.Error = err.Error()
	}
	return nil
}

func RpcServer() {
	rpc.Register(new(RpcFunc))
	lis, err := net.Listen("tcp", "127.0.0.1:8096")
	if err != nil {
		log("F", "rpc:", err.Error())
		return
	}

	log("i", "rpc_server start")

	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}

		go func(conn net.Conn) {
			log("i", "rpc_connect")
			jsonrpc.ServeConn(conn)
		}(conn)
	}
}
