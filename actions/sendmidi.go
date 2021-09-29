package actions

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/rtmididrv"
	"gopkg.in/yaml.v3"
)

type (
	SendMidiAction struct {
		DeviceName    string    `yaml:"device"`
		ChannelNumber uint8     `yaml:"channel"`
		Port          *midi.Out `yaml:"-"`
		Messages      []ControlChangeSpec
	}

	ControlChangeSpec struct {
		Control uint8
		Value   *uint8
	}
)

func (action SendMidiAction) Act(value, lastvalue int) error {
	wr := writer.New(*action.Port)
	wr.SetChannel(action.ChannelNumber)

	for _, cc := range action.Messages {
		if cc.Value != nil {
			value = int(*cc.Value)
		}

		log.Debug().
			Int("channel", int(action.ChannelNumber)).
			Int("control", int(cc.Control)).
			Int("value", value).Msgf("sending midi control change")
		if err := writer.ControlChange(wr, cc.Control, uint8(value)); err != nil {
			return err
		}
	}

	return nil
}

func NewSendMidiAction(args yaml.Node) (Action, error) {
	action := SendMidiAction{}

	if err := args.Decode(&action); err != nil {
		return nil, err
	}

	drv, err := driver.New()
	if err != nil {
		return nil, err
	}

	port, err := midi.OpenOut(drv, -1, action.DeviceName)
	if err != nil {
		return nil, err
	}
	action.Port = &port
	log.Debug().Str("port", port.String()).Msgf("opened output port")

	return action, nil
}
