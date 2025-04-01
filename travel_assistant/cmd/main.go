package main

import (
	"context"
	"flag"

	"github.com/Wenrh2004/travel_assistant/cmd/wire"
	"github.com/Wenrh2004/travel_assistant/pkg/util/config"
	"github.com/Wenrh2004/travel_assistant/pkg/util/log"
)

func main() {
	var envConf = flag.String("conf", "internal/config/config.yml", "config path, eg: -conf ./config/config.yml")
	flag.Parse()

	conf := config.NewConfig(*envConf)
	logger := log.NewLog(conf)
	app, cleanup, err := wire.NewWire(conf, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
}
