package server

import (
	"log"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func StartProfiler() {
	r := gin.Default()
	pprof.Register(r, "")

	log.Println("starting profiling server at 0.0.0.0:6060")
	r.Run(":6060")
}
