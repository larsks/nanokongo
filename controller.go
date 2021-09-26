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
		config   *Config
		Driver   *driver.Driver
		Port     midi.In
		Channel  uint8
		Controls map[uint8]*Control
	}

	ControlTypeEnum int

	Control struct {
		Number    uint8
		Type      ControlTypeEnum
		LastValue uint8
		Scale     *ScaleSpec
		OnPress   []actions.Action
		OnRelease []actions.Action
		OnChange  []actions.Action
	}

	ScaleSpec struct {
		MinOutput int
		MaxOutput int
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
		config:   cfg,
		Controls: make(map[uint8]*Control),
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

func NewControl(controlnumber uint8, controltype ControlTypeEnum, scalerange []int) *Control {
	control := Control{
		Number: controlnumber,
		Type:   controltype,
	}

	if len(scalerange) == 2 {
		control.Scale = &ScaleSpec{
			MinOutput: scalerange[0],
			MaxOutput: scalerange[1],
		}
	}

	return &control
}

func (control *Control) ScaleValue(value uint8) int {
	var newval int

	if control.Scale == nil {
		newval = int(value)
	} else {
		newval = int(
			(float32(value)/float32(127))*
				float32(control.Scale.MaxOutput-control.Scale.MinOutput)) + control.Scale.MinOutput
	}

	return newval
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
	controller.Channel = controller.config.Channel

	for number, controlspec := range controller.config.Controls {
		var err error

		log.Debug().Msgf("found entry for control %d", number)
		control := NewControl(number,
			ControlTypeFromName(controlspec.Type),
			controlspec.ScaleRange)
		controller.Controls[number] = control
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
		matched, err := filepath.Match(controller.config.Device, in.String())
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

func (controller *Controller) HandleControlChange(pos *reader.Position, channelNum, controlNum, value uint8) {
	log := log.With().
		Int("channel", int(channelNum)).
		Int("control", int(controlNum)).
		Int("value", int(value)).
		Logger()

	if channelNum != controller.Channel {
		log.Info().Msgf("not listening on channel")
		return
	}

	control, exists := controller.Controls[controlNum]
	if !exists {
		log.Info().Msgf("no matching control configure")
		return
	}

	log = log.With().Int("lastvalue", int(control.LastValue)).Logger()

	switch control.Type {
	case ControlTypeButton:
		log.Debug().Msgf("handling button")
		if value != 0 && value != 127 {
			log.Warn().Msgf("value out of range")
			return
		}

		if value == 0 && control.LastValue == 127 {
			for _, action := range control.OnRelease {
				action.Act(int(value), int(control.LastValue))
			}
		} else if value == 127 && control.LastValue == 0 {
			for _, action := range control.OnPress {
				action.Act(int(value), int(control.LastValue))
			}
		}

		if value != control.LastValue {
			for _, action := range control.OnChange {
				action.Act(int(value), int(control.LastValue))
			}
		}
	case ControlTypeKnob:
		svalue := control.ScaleValue(value)
		slastvalue := control.ScaleValue(control.LastValue)
		log.Debug().Msgf("handling knob")
		if value != control.LastValue {
			for _, action := range control.OnChange {
				action.Act(svalue, slastvalue)
			}
		}
	default:
		panic(fmt.Errorf("unknown control type: %d", control.Type))
	}

	control.LastValue = value
}
