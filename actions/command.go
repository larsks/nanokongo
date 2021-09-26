package actions

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type (
	CommandAction struct {
		KeySpecs []string
	}
)

func (action CommandAction) Act(value int) error {
	log.Warn().Msgf("execute command action (%d)", value)
	return nil
}

func (action CommandAction) String() string {
	return "command"
}

func NewCommandAction(args yaml.Node) (Action, error) {
	action := CommandAction{}

	if err := args.Decode(&action.KeySpecs); err != nil {
		return nil, err
	}

	log.Debug().Msgf("sendkeys config: %+v", action)

	return action, nil
}
