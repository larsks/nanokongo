package actions

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type (
	SendKeysAction struct {
		KeySpecs []string
	}
)

func (action SendKeysAction) Act() error {
	log.Warn().Msgf("execute sendkeys action")
	return nil
}

func (action SendKeysAction) String() string {
	return "sendKeys"
}

func NewSendKeysAction(args yaml.Node) (Action, error) {
	action := SendKeysAction{}

	if err := args.Decode(&action.KeySpecs); err != nil {
		return nil, err
	}

	log.Debug().Msgf("sendkeys config: %+v", action)

	return action, nil
}
