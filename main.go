package main

import (
	"flag"

	"github.com/COTBU/notifier/config"
	"github.com/COTBU/notifier/service"
)

func main() {
	path := ""
	flag.StringVar(&path, "config", "config/config.yaml", "config path")
	flag.Parse()

	cfg, err := config.Get(path)
	if err != nil {
		panic(err)
	}

	srv := service.New(cfg)

	//engine := gin.Default()
	////
	////notifyGroup := engine.Group("notifier")
	////notifyGroup.POST("", srv.NewEvent)
	//
	//if err := engine.Run("localhost:" + cfg.Service.Port); err != nil {
	//	panic(err)
	//}

	if err := srv.RunConsumer(); err != nil {
		panic(err)
	}

	if err := srv.CloseClient(); err != nil {
		panic(err)
	}
}
