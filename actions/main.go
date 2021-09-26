package actions

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type (
	Action interface {
		Act(int, int) error
	}

	ActionConstructor func(args yaml.Node) (Action, error)
)

var ActionMap map[string]ActionConstructor = make(map[string]ActionConstructor)

func RegisterAction(name string, constructor ActionConstructor) {
	ActionMap[name] = constructor
}

func LookupAction(want string) ActionConstructor {
	for have, constructor := range ActionMap {
		log.Debug().Str("have", have).Str("want", want).Send()
		if have == want {
			return constructor
		}
	}
	return nil
}

func init() {
	RegisterAction("sendKeys", NewSendKeysAction)
	RegisterAction("command", NewCommandAction)
	RegisterAction("sendMouse", NewSendMouseAction)
}
