package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/larsks/nanokongo/decouple"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	err := decouple.Load()
	decouple.SetPrefix("NANOKONGO_")

	loglevel, _ := decouple.GetIntInRange("LOGLEVEL", 1, -1, 5)
	cfgfilename, _ := decouple.GetString("CONFIG", "config.yml")
	debug, _ := decouple.GetBool("DEBUG", false)

	zerolog.SetGlobalLevel(zerolog.Level(loglevel))
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err != nil {
		log.Debug().Err(err).Send()
	}

	defer func() {
		// Set NANOKONGO_DEBUG=true if you want to see
		// backtraces for panics.
		if err := recover(); !debug && err != nil {
			log.Fatal().Err(err.(error)).Send()
		} else {
			panic(err)
		}
	}()

	cfg, err := ReadConfigFromFile(cfgfilename)
	must(err)

	ctrl, err := NewController(cfg)
	must(err)
	must(ctrl.Open())
	defer ctrl.Close()

	log.Info().Msgf("listening for events")
	must(ctrl.Listen())

	for {
		time.Sleep(1 * time.Second)
	}
}
