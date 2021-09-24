package main

import (
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	driver "gitlab.com/gomidi/rtmididrv"
)

type (
	Controller struct {
		Config *Config
		Driver *driver.Driver
		Port   midi.In
	}
)

func NewController(cfg *Config) (*Controller, error) {
	ctrl := Controller{
		Config: cfg,
	}

	drv, err := driver.New()
	if err != nil {
		return nil, err
	}
	ctrl.Driver = drv

	return &ctrl, nil
}

func (ctrl *Controller) Open() error {
	ins, err := ctrl.Driver.Ins()
	if err != nil {
		return err
	}

	var selected midi.In
	for _, in := range ins {
		log.Debug().Str("portname", in.String()).Msg("looking for device")
		matched, err := filepath.Match(ctrl.Config.Device, in.String())
		if err != nil {
			return err
		}

		if matched {
			selected = in
			break
		}
	}

	if selected == nil {
		return fmt.Errorf("unable to find device")
	}

	log.Debug().Str("portname", selected.String()).Msg("found device")

	if err = selected.Open(); err != nil {
		return err
	}

	ctrl.Port = selected

	return nil

}

func (ctrl *Controller) Close() {
	ctrl.Port.Close()
	ctrl.Driver.Close()
}

func (ctrl *Controller) Listen() error {
	rd := reader.New(
		reader.NoLogger(),
		reader.ControlChange(ctrl.HandleControlChange),
	)

	return rd.ListenTo(ctrl.Port)
}

func (ctrl *Controller) HandleControlChange(pos *reader.Position, channel, control, value uint8) {
	log := log.With().
		Int("channel", int(channel)).
		Int("control", int(control)).
		Int("value", int(value)).
		Logger()

	log.Debug().Msg("scanning config")

	for _, c := range ctrl.Config.Controls {
		log = log.With().Str("type", c.Type).Logger()
		if c.Control == control {
			log.Debug().Msg("found match")
		}
	}
}
