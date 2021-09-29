# nanokongo: generate input events from midi control changes

## Configuration overview

The nanokongo configuration file is a [YAML][] file. By default
nanokongo will look for `~/.config/nanokongo/config.yml`, but you can
specify a path to a different file in the `NANOKONGO_CONFIG`
environment variable.

[yaml]: https://en.wikipedia.org/wiki/YAML

The configuration file has three top-level keys:

- `device`: name of the MIDI device from which to read messages. This
  may be any substring that uniquely identifies the target device; for
  example, on my system, I have:

  ```
  $ midicat  ins
  MIDI inputs
  [0] Midi Through:Midi Through Port-0 14:0
  [1] nanoKONTROL2:nanoKONTROL2 nanoKONTROL2 _ CTR 36:0
  ```

  So I set `device` like this:

  ```
  device: nanoKONTROL2
  ```

- `channel`: the MIDI channel on which to listen. Defaults to `0` if
  not provided.

- `controls`: A mapping of control numbers to actions. For example, to
  send `F1` whenever a button with MIDI control 48 is pushed:

  ```
  controls:
    48:
      type: button
      onRelease:
        - sendKeys:
            - keys: [f1]
  ```

## Describing actions

The values in the `controls` map describe the actions that will be
triggered by MIDI events.  A control mapping may have the following
fields:

- `type`: This may be either `button` or `knob`.

  A `button` control has `onRelease`, `onPress`, and `onChange` triggers.

  A `knob` control has an `onChange` trigger.

- `onPress`: This describes actions that occur when a button is
  pressed.

- `onRelease`: This describes actions that occur when a button is
  released.

- `onChange`: This describes actions that occur when the value of a
  MIDI control changes. The `onChange` trigger for `buttons` will only
  run if the control value is `127` or `0`; for `knobs` it runs on
  every change.

## Available actions

### sendKeys

Send keystrokes with optional modifier keys.

### sendMouse

Send mouse events.

### command

Run commands.

### sendMidi

Send MIDI control change messages.

## Example configuration file
