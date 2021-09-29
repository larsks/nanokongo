// The actions module implements actions that nanokongo can trigger
// in response to MIDI control change messages.
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

// RegisterAction associates an Action constructor with the specified name.
func RegisterAction(name string, constructor ActionConstructor) {
	log.Debug().Str("action", name).Msgf("register action")
	ActionMap[name] = constructor
}

func RegisterActions() {
	RegisterAction("sendKeys", NewSendKeysAction)
	RegisterAction("command", NewCommandAction)
	RegisterAction("sendMouse", NewSendMouseAction)
	RegisterAction("sendMidi", NewSendMidiAction)
}

// LookUpAction looks for a named action in the list of registered actions
// and returns the appropriate constructor.
func LookupAction(want string) ActionConstructor {
	for have, constructor := range ActionMap {
		log.Trace().Str("have", have).Str("want", want).Send()
		if have == want {
			return constructor
		}
	}
	return nil
}
