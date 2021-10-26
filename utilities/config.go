package utilities

import (
	"errors"
	"fmt"
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
	EnvChatPort      = "CHAT_PORT"
	EnvRegHostname   = "REG_HOSTNAME"
	EnvSeqHostname   = "SEQ_HOSTNAME"
	EnvSize          = "SIZE"
	EnvEnableVerbose = "VERBOSE"
	EnvRedisHostname = "REDIS_HOSTNAME"
	EnvTypeOfService = "TOS"
)

type RegistrationServerConfig struct {
	RegPort       uint16
	ChatGroupSize uint16
	SeqHostname   string
	TypeOfService string
}

type SequencerServerConfig struct {
	ChatPort      uint16
	RegPort       uint16
	ChatGroupSize uint16
	TypeOfService string
}

type PeerConfig struct {
	WebServerPort uint16
	RegPort       uint16
	ChatPort      uint16
	ChatGroupSize uint16
	RegHostname   string
	SeqHostname   string
	RedisHostname string
	Verbose       bool
	TypeOfService string
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
		return 0, fmt.Errorf("%v [%v]", ErrEnvNotFound, envVar)
	}
}

func SetupSequencer() (*SequencerServerConfig, error) {
	c := &SequencerServerConfig{}

	val, err := parseUint16FromEnv(EnvChatPort)
	if err != nil {
		return c, err
	}
	c.ChatPort = val

	val, err = parseUint16FromEnv(EnvRegPort)
	if err != nil {
		return c, err
	}
	c.RegPort = val

	val, err = parseUint16FromEnv(EnvSize)
	if err != nil {
		return c, err
	}
	c.ChatGroupSize = val

	tos, isPresent := os.LookupEnv(EnvTypeOfService)
	if !isPresent {
		return c, fmt.Errorf("%v [TOS]", ErrEnvNotFound)
	} else {
		c.TypeOfService = strings.ToLower(tos)
	}

	return c, nil
}

func SetupRegistrationServer() (*RegistrationServerConfig, error) {
	c := &RegistrationServerConfig{}

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

	rhost, isPresent := os.LookupEnv(EnvSeqHostname)
	if !isPresent {
		return c, ErrEnvNotFound
	}
	c.SeqHostname = rhost

	tos, isPresent := os.LookupEnv(EnvTypeOfService)
	if !isPresent {
		return c, fmt.Errorf("%v [TOS]", ErrEnvNotFound)
	} else {
		c.TypeOfService = strings.ToLower(tos)
	}

	return c, nil
}

func SetupPeer() (*PeerConfig, error) {
	c := &PeerConfig{}

	val, err := parseUint16FromEnv(EnvRegPort)
	if err != nil {
		return c, err
	}
	c.RegPort = val

	val, err = parseUint16FromEnv(EnvChatPort)
	if err != nil {
		return c, err
	}
	c.ChatPort = val

	val, err = parseUint16FromEnv(EnvSize)
	if err != nil {
		return c, err
	}
	c.ChatGroupSize = val

	enable, isPresent := os.LookupEnv(EnvEnableVerbose)
	if !isPresent {
		c.Verbose = false
	} else {
		if strings.ToLower(enable) == "true" {
			c.Verbose = true
		} else {
			c.Verbose = false
		}
	}

	val, err = parseUint16FromEnv(EnvWsPort)
	if err != nil {
		return c, err
	}
	c.WebServerPort = val

	rhost, isPresent := os.LookupEnv(EnvRegHostname)
	if !isPresent {
		return c, ErrEnvNotFound
	}
	c.RegHostname = rhost

	rhost, isPresent = os.LookupEnv(EnvSeqHostname)
	if !isPresent {
		return c, ErrEnvNotFound
	}
	c.SeqHostname = rhost

	rhost, isPresent = os.LookupEnv(EnvRedisHostname)
	if !isPresent {
		return c, ErrEnvNotFound
	}
	c.RedisHostname = rhost

	tos, isPresent := os.LookupEnv(EnvTypeOfService)
	if !isPresent {
		return c, fmt.Errorf("%v [TOS]", ErrEnvNotFound)
	} else {
		c.TypeOfService = strings.ToLower(tos)
	}

	return c, nil
}
