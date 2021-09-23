package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

var (
	ErrEnvNotFound = errors.New("environment variable unset")
)

var (
	EnvWsPort        = "WS_PORT"
	EnvRegPort       = "REG_PORT"
	EnvRegHostname   = "REG_HOSTNAME"
	EnvSize          = "SIZE"
	EnvEnableVerbose = "VERBOSE"
)

type Configuration struct {
	WebServerPort uint16
	RegPort       uint16
	ChatGroupSize uint16
	RegHostname   string
	EnableVerbose bool
}

func parseUint16FromEnv(envVar string) (uint16, error) {
	portStr, isPresent := os.LookupEnv(envVar)
	if isPresent {
		portUint16, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			return 0, err
		}
		// la funzione ParseUint pur prendendo come parametro la dimensione della variabile
		// restituisce comunque un uint64 per cui necessita di cast
		return uint16(portUint16), nil
	} else {
		return 0, ErrEnvNotFound
	}
}

func SetupRegistrationServerConfiguration() (*Configuration, error) {
	c := &Configuration{}

	val, err := parseUint16FromEnv(EnvRegPort)
	if err != nil {
		return c, err
	}
	c.RegPort = val

	val, err = parseUint16FromEnv(EnvSize)
	if err != nil {
		return c, err
	}
	c.ChatGroupSize = val

	enable, isPresent := os.LookupEnv(EnvEnableVerbose)
	if !isPresent {
		c.EnableVerbose = false
	} else {
		if strings.ToLower(enable) == "true" {
			c.EnableVerbose = true
		} else {
			c.EnableVerbose = false
		}
	}

	return c, nil
}

func SetupPeerConfiguration() (*Configuration, error) {

	c, err := SetupRegistrationServerConfiguration()
	if err != nil {
		return c, nil
	}

	val, err := parseUint16FromEnv(EnvWsPort)
	if err != nil {
		return c, err
	}
	c.WebServerPort = val

	rhost, isPresent := os.LookupEnv(EnvRegHostname)
	if !isPresent {
		return c, ErrEnvNotFound
	}
	c.RegHostname = rhost

	return c, nil
}
