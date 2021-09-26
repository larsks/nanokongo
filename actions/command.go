package actions

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type (
	CommandAction struct {
		Args []string
	}
)

func substituteValue(oldargs []string, value int) []string {
	var args []string

	valString := fmt.Sprintf("%d", value)

	for _, arg := range oldargs {
		args = append(args, strings.Replace(arg, "{value}", valString, -1))
	}

	return args
}

func (action CommandAction) Act(value, lastvalue int) error {
	args := substituteValue(action.Args, value)
	err := exec.Command(args[0], args[1:]...).Run()
	log.Warn().Msgf("execute command action: %s", args)
	if err != nil {
		return err
	}

	return nil
}

func NewCommandAction(args yaml.Node) (Action, error) {
	action := CommandAction{}

	if err := args.Decode(&action.Args); err != nil {
		return nil, err
	}

	log.Debug().Msgf("command config: %+v", action)

	return action, nil
}
