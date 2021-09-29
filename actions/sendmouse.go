package actions

import (
	"fmt"
	"strings"

	"github.com/bendahl/uinput"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var Mouse uinput.Mouse

type (
	// The sendMouse action lets you send mouse inputs in reaction to MIDI
	// control change messages.
	//
	// Supported attributes are:
	//
	// - click: name of a button to click (one of "left", "right")
	// - press: name of a button to press
	// - release: name of a button to release
	// - x: relative x movement
	// - y: relative y movement
	// - wheelx: relative horizontal wheel movement
	// - wheely: relative vertical wheel movement
	//
	// Example config:
	//
	//	controls:
	//	  58:
	//	    type: button
	//	    onRelease:
	//	      - sendMouse:
	//	          click: left
	SendMouseAction struct {
		Click   string
		Press   string
		Release string
		X       int
		Y       int
		WheelX  int
		WheelY  int
	}
)

func (action SendMouseAction) Act(value, lastvalue int) error {
	var dir int
	var err error

	log.Warn().Msgf("execute sendmouse action (%+v)", action)
	if value > lastvalue {
		dir = 1
	} else {
		dir = -1
	}

	if action.X != 0 {
		x := action.X

		log.Debug().Int("x", x).Int("direction", dir).Msg("move x")

		if dir < 0 {
			err = Mouse.MoveLeft(int32(x))
		} else {
			err = Mouse.MoveRight(int32(x))
		}
		if err != nil {
			return err
		}
	}

	if action.Y != 0 {
		y := action.Y

		log.Debug().Int("y", y).Int("direction", dir).Msg("move y")

		if dir < 0 {
			err = Mouse.MoveUp(int32(y))
		} else {
			err = Mouse.MoveDown(int32(y))
		}

		if err != nil {
			return err
		}
	}

	if action.WheelX != 0 {
		x := action.WheelX

		log.Debug().Int("x", x).Int("direction", dir).Msg("wheel x")
		if err = Mouse.Wheel(true, int32(dir*x)); err != nil {
			return err
		}
	}

	if action.WheelY != 0 {
		y := action.WheelY

		log.Debug().Int("y", y).Int("direction", dir).Msg("wheel y")
		if err = Mouse.Wheel(false, int32(dir*y)); err != nil {
			return err
		}
	}

	if action.Press != "" {
		buttons := strings.Split(action.Press, "+")
		for _, button := range buttons {
			log.Debug().Str("button", button).Msg("press")

			switch button {
			case "left":
				err = Mouse.LeftPress()
			case "right":
				err = Mouse.RightPress()
			default:
				panic(fmt.Errorf("unknown button: %s", button))
			}

			if err != nil {
				return err
			}
		}
	}

	if action.Release != "" {
		buttons := strings.Split(action.Release, "+")
		for _, button := range buttons {
			log.Debug().Str("button", button).Msg("release")

			switch button {
			case "left":
				err = Mouse.LeftRelease()
			case "right":
				err = Mouse.RightRelease()
			default:
				panic(fmt.Errorf("unknown button: %s", button))
			}

			if err != nil {
				return err
			}
		}
	}

	if action.Click != "" {
		buttons := strings.Split(action.Click, "+")
		for _, button := range buttons {
			log.Debug().Str("button", button).Msg("click")

			switch button {
			case "left":
				err = Mouse.LeftClick()
			case "right":
				err = Mouse.RightClick()
			default:
				panic(fmt.Errorf("unknown button: %s", button))
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewSendMouseAction(args yaml.Node) (Action, error) {
	action := SendMouseAction{}

	if Mouse == nil {
		if err := initVirtualMouse(); err != nil {
			return nil, err
		}
	}

	if err := args.Decode(&action); err != nil {
		return nil, err
	}

	log.Debug().Msgf("sendmouse config: %+v", action)

	return action, nil
}

func initVirtualMouse() error {
	log.Info().Msgf("initializing virtual mouse")
	mouse, err := uinput.CreateMouse("/dev/uinput", []byte("nanokongo-mouse"))
	if err != nil {
		return err
	}

	Mouse = mouse
	return nil
}
