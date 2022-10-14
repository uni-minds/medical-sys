package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type DeviceStatus struct {
	DeviceName   string
	DeviceStatus string
}

func MediBoxGetHandler(ctx *gin.Context) {
	devid := ctx.Param("devid")
	ops := ctx.Param("ops")

	token := ctx.Query("token")
	ordtime := ctx.Query("time")

	log.Println("OPS=", ops, "; DEV=", devid, "; TOKEN=", token, "; T=", ordtime)

	switch ops {
	case "record-start":
		log.Println("Start record.")
		ctx.JSON(http.StatusOK, SuccessReturn("OK"))

	case "record-stop":
		log.Println("Stop record.")
		ctx.JSON(http.StatusOK, SuccessReturn("UUID_DEMO_VID"))

	case "capture":
		log.Println("Cap.")
		ctx.JSON(http.StatusOK, SuccessReturn("UUID_DEMO_PIC"))

	case "status":
		log.Println("Status")
		s := DeviceStatus{
			DeviceName:   "B205",
			DeviceStatus: "busy",
		}
		ctx.JSON(http.StatusOK, SuccessReturn(s))

	default:
		ctx.JSON(http.StatusOK, FailReturn(ops))
	}
}
