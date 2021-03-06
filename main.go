package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/larsks/go-decouple"
	"github.com/larsks/nanokongo/actions"
	"github.com/larsks/nanokongo/version"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func default_config_file() string {
	path := "nanokongo.yml"

	if val := os.Getenv("XDG_CONFIG_HOME"); val != "" {
		path = filepath.Join(val, "nanokongo", "config.yml")
	} else if val, err := os.UserHomeDir(); err == nil && val != "" {
		path = filepath.Join(val, ".config", "nanokongo", "config.yml")
	}

	log.Debug().Str("path", path).Msgf("default config path")
	return path
}

func main() {
	err := decouple.Load()
	decouple.SetPrefix("NANOKONGO_")

	loglevel, _ := decouple.GetIntInRange("LOGLEVEL", 1, -1, 5)
	zerolog.SetGlobalLevel(zerolog.Level(loglevel))
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfgfilename, _ := decouple.GetString("CONFIG", default_config_file())
	debug, _ := decouple.GetBool("DEBUG", false)

	if err != nil {
		log.Debug().Err(err).Send()
	}

	log.Info().
		Str("version", version.BuildVersion).
		Str("ref", version.BuildRef).
		Str("date", version.BuildDate).Msgf("starting")

	defer func() {
		// Set NANOKONGO_DEBUG=true if you want to see
		// backtraces for panics.
		if err := recover(); !debug && err != nil {
			log.Fatal().Err(err.(error)).Send()
		} else {
			panic(err)
		}
	}()

	actions.RegisterActions()

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
