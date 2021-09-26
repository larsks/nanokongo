// Inspired by https://github.com/henriquebastos/python-decouple
package decouple

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var nameprefix string

func SetPrefix(prefix string) {
	nameprefix = prefix
}

func LookupEnv(name string) (string, bool) {
	name = fmt.Sprintf("%s%s", nameprefix, name)
	return os.LookupEnv(name)
}

func GetString(name, defval string) (string, bool) {
	val, exists := LookupEnv(name)
	if !exists {
		return defval, false
	}

	return val, true
}

func GetInt(name string, defval int) (int, bool) {
	val, exists := LookupEnv(name)
	if !exists {
		return defval, false
	}

	ret, err := strconv.ParseInt(val, 0, 0)
	if err != nil {
		return defval, false
	}

	return int(ret), true
}

func GetIntInRange(name string, defval, minval, maxval int) (int, bool) {
	ret, exists := GetInt(name, defval)

	switch {
	case ret < minval:
		ret = minval
	case ret > maxval:
		ret = maxval
	}

	return int(ret), exists
}

func GetBool(name string, defval bool) (bool, bool) {
	val, exists := LookupEnv(name)
	if !exists {
		return defval, false
	}

	ret, err := strconv.ParseBool(val)
	if err != nil {
		return defval, false
	}

	return ret, true
}

func Load(filenames ...string) {
	godotenv.Load(filenames...)
}
