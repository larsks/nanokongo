// The sendkeys action allows you to send keystrokes in resposne to
// MIDI control change events.
//
// See https://github.com/bendahl/uinput/blob/master/keycodes.go for a
// list of valid keycodes. Remove the "Key" part of the name, and
// convert the rest to lower case, so "KeyA" -> "a", and
// "KeyLeftshift" -> "leftshift".
//
// Example configuration:
//
//	controls:
//	  41:
//	    type: button
//	    onRelease:
//	      - sendKeys:
//	          keys: [leftshift+o, o, d, d, b, i, t]
//                delay: 100
//
package actions

// Generate keycodes.go, which maps key names to keycodes.
//go:generate ./genkeycodes.sh

import (
	"fmt"
	"strings"
	"time"

	"github.com/bendahl/uinput"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var Kbd uinput.Keyboard

type (
	SendKeysAction struct {
		Keys  []string
		Delay int
	}
)

func (action SendKeysAction) Act(value, lastvalue int) error {
	log.Warn().Msgf("execute sendkeys action (%+v)", action)
	for _, keyspec := range action.Keys {
		keys := strings.Split(keyspec, "+")

		// Press keys
		for i := range keys {
			key := keys[i]
			log.Debug().Msgf("keydown %s", key)
			keycode, exists := keycodes[strings.ToLower(key)]
			if !exists {
				return fmt.Errorf("no such key: %s", key)
			}

			if err := Kbd.KeyDown(keycode); err != nil {
				return err
			}
		}

		// Release keys
		for i := range keys {
			key := keys[len(keys)-1-i]
			log.Debug().Msgf("keyup %s", key)
			keycode, exists := keycodes[strings.ToLower(key)]
			if !exists {
				return fmt.Errorf("no such key: %s", key)
			}

			if err := Kbd.KeyUp(keycode); err != nil {
				return err
			}
		}

		if action.Delay > 0 {
			time.Sleep(time.Duration(action.Delay) * time.Millisecond)
		}
	}

	return nil
}

func NewSendKeysAction(args yaml.Node) (Action, error) {
	action := SendKeysAction{}

	if Kbd == nil {
		if err := initVirtualKeyboard(); err != nil {
			return nil, err
		}
	}

	if err := args.Decode(&action); err != nil {
		return nil, err
	}

	log.Debug().Msgf("sendkeys config: %+v", action)

	return action, nil
}

func initVirtualKeyboard() error {
	log.Info().Msgf("initializing virtual keyboard")
	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("nanokongo-kbd"))
	if err != nil {
		return err
	}

	Kbd = keyboard
	return nil
}
