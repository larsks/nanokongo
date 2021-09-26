package main

import (
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	driver "gitlab.com/gomidi/rtmididrv"
	yaml "gopkg.in/yaml.v3"

	"github.com/larsks/nanokongo/actions"
)

const (
	ControlTypeButton = iota
	ControlTypeKnob
)

type (
	Controller struct {
		Config   *Config
		Driver   *driver.Driver
		Port     midi.In
		Controls []Control
	}

	ControlTypeEnum int

	Control struct {
		Number    uint8
		Type      ControlTypeEnum
		OnPress   []actions.Action
		OnRelease []actions.Action
		OnChange  []actions.Action
	}
)

func (t ControlTypeEnum) String() string {
	var name string

	switch t {
	case ControlTypeButton:
		name = "Button"
	case ControlTypeKnob:
		name = "Knob"
	default:
		panic(fmt.Errorf("unknown control type: %d", t))
	}

	return name
}

func NewController(cfg *Config) (*Controller, error) {
	controller := Controller{
		Config: cfg,
	}

	drv, err := driver.New()
	if err != nil {
		return nil, err
	}
	controller.Driver = drv
	if err := controller.ProcessConfig(); err != nil {
		return nil, err
	}

	return &controller, nil
}

func NewControl(controlnumber uint8, controltype ControlTypeEnum) Control {
	return Control{
		Number: controlnumber,
		Type:   controltype,
	}
}

func ControlTypeFromName(t string) ControlTypeEnum {
	var res ControlTypeEnum

	switch t {
	case "button":
		res = ControlTypeButton
	case "knob":
		res = ControlTypeKnob
	default:
		panic(fmt.Errorf("unknown control type: %s", t))
	}

	return res
}

func buildActionList(spec []map[string]yaml.Node) ([]actions.Action, error) {
	var actionlist []actions.Action

	for _, actionspec := range spec {
		if len(actionspec) != 1 {
			return nil, fmt.Errorf("invalid action spec: %+v", actionspec)
		}
		log.Debug().Msgf("action: %+v", actionspec)
		for name, config := range actionspec {
			constructor := actions.LookupAction(name)
			if constructor == nil {
				log.Warn().Msgf("%s: unimplemented", name)
				continue
			}
			action, err := constructor(config)
			if err != nil {
				return nil, err
			}
			actionlist = append(actionlist, action)
		}
	}

	return actionlist, nil
}

func (controller *Controller) ProcessConfig() error {
	for _, controlspec := range controller.Config.Controls {
		var err error

		log.Debug().Msgf("found entry for control %d", controlspec.Control)
		control := NewControl(controlspec.Control,
			ControlTypeFromName(controlspec.Type))
		controller.Controls = append(controller.Controls, control)
		log.Debug().Msgf("control: %+v", control)

		if control.Type == ControlTypeButton {
			control.OnRelease, err = buildActionList(controlspec.OnRelease)
			if err != nil {
				return err
			}

			control.OnPress, err = buildActionList(controlspec.OnPress)
			if err != nil {
				return err
			}
		} else if control.Type == ControlTypeKnob {
			control.OnChange, err = buildActionList(controlspec.OnChange)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (controller *Controller) Open() error {
	ins, err := controller.Driver.Ins()
	if err != nil {
		return err
	}

	var selected midi.In
	for _, in := range ins {
		log.Debug().Str("portname", in.String()).Msg("looking for device")
		matched, err := filepath.Match(controller.Config.Device, in.String())
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

	controller.Port = selected

	return nil

}

func (controller *Controller) Close() {
	controller.Port.Close()
	controller.Driver.Close()
}

func (controller *Controller) Listen() error {
	rd := reader.New(
		reader.NoLogger(),
		reader.ControlChange(controller.HandleControlChange),
	)

	return rd.ListenTo(controller.Port)
}

func (controller *Controller) HandleControlChange(pos *reader.Position, channel, control, value uint8) {
	log := log.With().
		Int("channel", int(channel)).
		Int("control", int(control)).
		Int("value", int(value)).
		Logger()

	log.Debug().Msg("scanning config")

	for _, c := range controller.Config.Controls {
		log = log.With().Str("type", c.Type).Logger()
		if c.Control == control {
			log.Debug().Msg("found match")
		}
	}
}
