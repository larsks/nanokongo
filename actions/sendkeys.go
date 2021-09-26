package actions

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
		KeySpecs []KeySpec
	}

	KeySpec struct {
		Keys  []string
		Mods  []string
		Delay int
	}
)

func (action SendKeysAction) Act(value, lastvalue int) error {
	for _, keyspec := range action.KeySpecs {
		log.Warn().Msgf("execute sendkeys action (%+v)", keyspec)
		for _, modname := range keyspec.Mods {
			log.Debug().Msgf("keydown modifier %s", modname)
			mod, exists := keycodes[strings.ToLower(modname)]
			if !exists {
				log.Warn().Msgf("no such key: %s", modname)
				return fmt.Errorf("no such key: %s", modname)
			}

			if err := Kbd.KeyDown(mod); err != nil {
				return err
			}
		}

		for _, keyname := range keyspec.Keys {
			log.Debug().Msgf("press key %s", keyname)
			key, exists := keycodes[strings.ToLower(keyname)]
			if !exists {
				log.Warn().Msgf("no such key: %s", keyname)
				return fmt.Errorf("no such key: %s", keyname)
			}

			if err := Kbd.KeyPress(key); err != nil {
				return err
			}
			if keyspec.Delay > 0 {
				time.Sleep(time.Duration(keyspec.Delay) * time.Millisecond)
			}
		}

		for _, modname := range keyspec.Mods {
			log.Debug().Msgf("keyup modifier %s", modname)
			mod, exists := keycodes[strings.ToLower(modname)]
			if !exists {
				log.Warn().Msgf("no such key: %s", modname)
				return fmt.Errorf("no such key: %s", modname)
			}

			if err := Kbd.KeyUp(mod); err != nil {
				return err
			}
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
