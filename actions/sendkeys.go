package actions

//go:generate ./genkeycodes.sh

import (
	"fmt"
	"strings"

	"github.com/bendahl/uinput"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var Kbd uinput.Keyboard

type (
	SendKeysAction struct {
		KeySpecs []KeySpec
	}

	KeySpec struct {
		Keys []string
		Mods []string
	}
)

func (action SendKeysAction) Act(value, lastvalue int) error {
	log.Warn().Msgf("execute sendkeys action (%d)", value)
	for _, keyspec := range action.KeySpecs {
		for _, modname := range keyspec.Mods {
			log.Debug().Msgf("keydown modifier %s", modname)
			mod, exists := keycodes[strings.ToLower(modname)]
			if !exists {
				log.Warn().Msgf("no such key: %s", modname)
				return fmt.Errorf("no such key: %s", modname)
			}

			Kbd.KeyDown(mod)
		}

		for _, keyname := range keyspec.Keys {
			log.Debug().Msgf("press key %s", keyname)
			key, exists := keycodes[strings.ToLower(keyname)]
			if !exists {
				log.Warn().Msgf("no such key: %s", keyname)
				return fmt.Errorf("no such key: %s", keyname)
			}

			Kbd.KeyPress(key)
		}

		for _, modname := range keyspec.Mods {
			log.Debug().Msgf("keyup modifier %s", modname)
			mod, exists := keycodes[strings.ToLower(modname)]
			if !exists {
				log.Warn().Msgf("no such key: %s", modname)
				return fmt.Errorf("no such key: %s", modname)
			}

			Kbd.KeyUp(mod)
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

	if err := args.Decode(&action.KeySpecs); err != nil {
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
