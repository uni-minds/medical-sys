package main

type RpcFunc struct {
}

type RpcRequest struct {
	Command string
	Param   []string
}

type RpcResponse struct {
	Result string
	Error  string
}
