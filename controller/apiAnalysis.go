package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type DeepBuild struct {
	StudiesUID string
	SeriesUID  string
	L          DeepBuildPoint
	R          DeepBuildPoint
	T          string
}

type DeepBuildPoint struct {
	Index int `json:"i"`
	X     int `json:"x"`
	Y     int `json:"y"`
}

func AnalysisCtPost(ctx *gin.Context) {
	class := ctx.Param("class")
	mode := ctx.Param("mode")
	switch class {
	case "cta":
		switch mode {
		case "deepbuild":
			var data DeepBuild
			err := ctx.BindJSON(&data)
			fmt.Println("Deepbuild:", data, err)
		}

	}
}
