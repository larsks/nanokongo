package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	/*
		defer func() {
			if err := recover(); err != nil {
				log.Fatal().Err(err.(error)).Send()
			}
		}()
	*/

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := ReadConfigFromFile("config.yml")
	must(err)
	fmt.Printf("%+v\n", cfg)

	ctrl, err := NewController(cfg)
	must(err)
	must(ctrl.Open())
	defer ctrl.Close()

	must(ctrl.Listen())

	for {
		time.Sleep(1 * time.Second)
	}
}
