package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"uni-minds.com/liuxy/medical-sys/tools"
)

func UiManageGetHandler(ctx *gin.Context) {
	action := ctx.Param("class")
	switch action {
	case "algo":
		ctx.HTML(http.StatusOK, "manage-algo.html", gin.H{
			"title":   "算法管理 | Medi-sys",
			"page_id": "manage-algo",
		})
		break

	case "user":
		ctx.HTML(http.StatusOK, "manage-user.html", gin.H{
			"title":   "用户管理| Medi-sys",
			"page_id": "manage-user",
		})
		break

	case "group":
		ctx.HTML(http.StatusOK, "manage-group.html", gin.H{
			"page_id": "manage-group",
			"title":   "群组管理 ｜ Medi-sys",
		})
		break

	case "media":
		ctx.HTML(http.StatusOK, "manage-media.html", gin.H{
			"page_id": "manage-media",
			"title":   "Medi-sys | 媒体管理",
		})
		break

	case "blockchain":
		ctx.HTML(http.StatusOK, "manage-blockchain.html", gin.H{
			"page_id": "manage-blockchain",
			"title":   "Medi-sys | 节点管理",
		})

		switch ctx.Query("action") {
		case "tps-stress":
			go func() {
				log("w", "start tps stress")
				nodes := []string{"52"}
				for i := 0; i < 10; i++ {
					log("i", "100tx/node * i=", i)
					for _, node := range nodes {
						go func(node string, i int) {
							url := fmt.Sprintf("http://192.168.1.%s/api/v1/demo/tps-stress?count=1000&context=%s", node, time.Now().Format("2006-01-02_15:04:05"))
							resp, _, err := tools.HttpGet(url)
							if err != nil {
								log("e", err.Error())
							} else if resp.Code != 200 {
								log("e", resp.Message)
							}
						}(node, i)
					}
					time.Sleep(2 * time.Second)

				}
				log("w", "tps stress finish. 100tx * 30time * 4node = 15000tx")
			}()
		}
		break
	case "browser":
		ctx.HTML(http.StatusOK, "manage-browser.html", gin.H{
			"page_id": "manage-browser",
			"title":   "Medi-sys | 区块链浏览器",
		})
		break

	case "upload":
		break
	}
	return
}
