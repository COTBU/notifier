package main

import (
	"flag"

	"github.com/gin-gonic/gin"

	"SOTBI/notifier/config"
)

func main() {
	path := ""
	flag.StringVar(&path, "config", "config/notifier.yaml", "config path")
	flag.Parse()

	cfg, err := config.Get(path)
	if err != nil {
		panic(err)
	}

	//srv, err := server.New(cfg)
	//if err != nil {
	//	panic(err)
	//}

	engine := gin.Default()
	//
	//notifyGroup := engine.Group("notifier")
	//notifyGroup.POST("", srv.NewEvent)

	if err := engine.Run("localhost:" + cfg.Service.Port); err != nil {
		panic(err)
	}
}
