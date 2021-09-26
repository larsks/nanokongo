#!/bin/sh

set -e

if ! [ -f keycodes.go.txt ]; then
	curl -o keycodes.go.txt -sfL https://raw.githubusercontent.com/bendahl/uinput/master/keycodes.go
fi

(
cat <<EOT
package actions

import "github.com/bendahl/uinput"

var keycodes map[string]int = map[string]int{
EOT

awk '
	/^\t*Key.* =/ {
		keyname = tolower(substr($1,4))
		printf "\t\"%s\": uinput.%s,\n", keyname, $1
	}
' keycodes.go.txt

cat <<EOT
}
EOT
) > keycodes.go
