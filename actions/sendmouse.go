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
	SendMouseAction struct {
		MouseSpec MouseSpec
	}

	MouseSpec struct {
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

	log.Warn().Msgf("execute sendmouse action (%+v)", action.MouseSpec)
	if value > lastvalue {
		dir = 1
	} else {
		dir = -1
	}

	if action.MouseSpec.X != 0 {
		x := action.MouseSpec.X

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

	if action.MouseSpec.Y != 0 {
		y := action.MouseSpec.Y

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

	if action.MouseSpec.WheelX != 0 {
		x := action.MouseSpec.WheelX

		log.Debug().Int("x", x).Int("direction", dir).Msg("wheel x")
		if err = Mouse.Wheel(true, int32(dir*x)); err != nil {
			return err
		}
	}

	if action.MouseSpec.WheelY != 0 {
		y := action.MouseSpec.WheelY

		log.Debug().Int("y", y).Int("direction", dir).Msg("wheel y")
		if err = Mouse.Wheel(false, int32(dir*y)); err != nil {
			return err
		}
	}

	if action.MouseSpec.Press != "" {
		buttons := strings.Split(action.MouseSpec.Press, "+")
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

	if action.MouseSpec.Release != "" {
		buttons := strings.Split(action.MouseSpec.Release, "+")
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

	if action.MouseSpec.Click != "" {
		buttons := strings.Split(action.MouseSpec.Click, "+")
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

	if err := args.Decode(&action.MouseSpec); err != nil {
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
