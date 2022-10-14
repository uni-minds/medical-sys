package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/tools"
)

func GetBlockchainNodelist(ctx *gin.Context) {
	resp, _, _ := tools.HttpGet(fmt.Sprintf("http://%s/api/v1/node/list", edaAddress))

	if resp.Code != 200 {
		ctx.JSON(http.StatusOK, FailReturn(400, resp.Message))
	} else {
		var list []global.NodeInfo
		bs, _ := json.Marshal(resp.Data)
		err := json.Unmarshal(bs, &list)
		if err != nil {
			log("e", "get nodelist:", err.Error())
			ctx.JSON(http.StatusOK, FailReturn(400, err.Error()))
		} else {
			var callback []global.NodeStatus
			for _, info := range list {
				alive := time.Now().Sub(info.LastTalk) < 35*time.Second
				if !alive {
					log("w", "node timeout:", info.Name, time.Now().Sub(info.LastTalk).Seconds(), info.LastTalk)
				}
				s := global.NodeStatus{
					Name:   info.Name,
					Alive:  alive,
					IP:     info.IP,
					Height: info.BlockHeight,
				}
				if info.IP == "localhost" {
					callback = append([]global.NodeStatus{s}, callback...)
				} else {
					callback = append(callback, s)
				}
			}
			ctx.JSON(http.StatusOK, SuccessReturn(callback))
		}
	}
}

func GetBlockchainTPS(ctx *gin.Context) {
	ip := ctx.Query("addr")
	if ip == "" {
		ip = edaAddress
	}
	url := fmt.Sprintf("http://%s/api/v1/node/tps", ip)
	resp, _, _ := tools.HttpGet(url)
	if resp.Code != 200 {
		ctx.JSON(http.StatusOK, FailReturn(400, resp.Message))
	} else {
		ctx.JSON(http.StatusOK, SuccessReturn(resp.Data))
	}
}

func GetBlockHeight(ctx *gin.Context) {
	h := ctx.Query("height")
	resp, _, _ := tools.HttpGet(fmt.Sprintf("http://%s/api/v1/block/record/height/%s", edaAddress, h))
	if resp.Code != 200 {
		ctx.JSON(http.StatusOK, FailReturn(400, resp.Message))
	} else {
		ctx.JSON(http.StatusOK, SuccessReturn(resp.Data))
	}
}

func WebSocket(ctx *gin.Context) {
	// change the reqest to websocket model
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(ctx.Writer, ctx.Request, nil)
	if error != nil {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	// websocket connect
	client := &Client{ID: "abab", Socket: conn, Send: make(chan []byte)}

	Manager.Register <- client

	go client.Read()
	go client.Write()
}
