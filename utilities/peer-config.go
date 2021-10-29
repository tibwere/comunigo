package utilities

import (
	"fmt"
	"os"
	"strings"
)

// In ottica OO, oggetto mantenente i
// metadati per la configurazione del
// peer
type PeerConfig struct {
	wsPort        uint16
	regPort       uint16
	chatPort      uint16
	chatGroupSize uint16
	regAddress    string
	seqAddress    string
	redisAddress  string
	verbose       bool
	tos           TypeOfService
}

// "Costruttore" dell'ogggetto PeerConfig
func InitPeerConfig() (*PeerConfig, error) {
	c := &PeerConfig{}

	val, err := parseUint16FromEnv(EnvRegPort)
	if err != nil {
		return c, err
	}
	c.regPort = val

	val, err = parseUint16FromEnv(EnvChatPort)
	if err != nil {
		return c, err
	}
	c.chatPort = val

	val, err = parseUint16FromEnv(EnvSize)
	if err != nil {
		return c, err
	}
	c.chatGroupSize = val

	enable, isPresent := os.LookupEnv(EnvEnableVerbose)
	if !isPresent {
		c.verbose = false
	} else {
		if strings.ToLower(enable) == "true" {
			c.verbose = true
		} else {
			c.verbose = false
		}
	}

	val, err = parseUint16FromEnv(EnvWsPort)
	if err != nil {
		return c, err
	}
	c.wsPort = val

	rhost, isPresent := os.LookupEnv(EnvRegHostname)
	if !isPresent {
		return c, ErrEnvNotFound
	}
	c.regAddress = rhost

	rhost, isPresent = os.LookupEnv(EnvSeqHostname)
	if !isPresent {
		return c, ErrEnvNotFound
	}
	c.seqAddress = rhost

	rhost, isPresent = os.LookupEnv(EnvRedisHostname)
	if !isPresent {
		return c, ErrEnvNotFound
	}
	c.redisAddress = rhost

	tos, isPresent := os.LookupEnv(EnvTypeOfService)
	if !isPresent {
		return c, fmt.Errorf("%v [TOS]", ErrEnvNotFound)
	} else {
		c.tos = setTOS(strings.ToLower(tos))
	}

	return c, nil
}

// "Metodo della classe PeerConfig" per il retrieve della porta
// esposta del web server
func (c *PeerConfig) GetWebServerPort() uint16 {
	return c.wsPort
}

// "Metodo della classe PeerConfig" per il retrieve della porta
// esposta dal server di registrazione
func (c *PeerConfig) GetRegistrationPort() uint16 {
	return c.regPort
}

// "Metodo della classe PeerConfig" per il retrieve della porta
// da usare nella comunicazione effettiva
func (c *PeerConfig) GetChatPort() uint16 {
	return c.chatPort
}

// "Metodo della classe PeerConfig" per il retrieve della dimensione
// del gruppo di multicast
func (c *PeerConfig) GetMulticastGroupSize() uint16 {
	return c.chatGroupSize
}

// "Metodo della classe PeerConfig" per il retrieve dell'indirizzo
// del nodo di registrazione
func (c *PeerConfig) GetRegistrationAddress() string {
	return c.regAddress
}

// "Metodo della classe PeerConfig" per il retrieve dell'indirizzo
// del sequencer
func (c *PeerConfig) GetSequencerAddress() string {
	return c.seqAddress
}

// "Metodo della classe PeerConfig" per il retrieve dell'indirizzo
// del datastore
func (c *PeerConfig) GetRedisAddress() string {
	return c.redisAddress
}

// "Metodo della classe PeerConfig" per verificare
// se richiesto o meno il verbose
func (c *PeerConfig) NeedVerbose() bool {
	return c.verbose
}

// "Metodo della classe PeerConfig" per il retrieve del
// tipo di servizio scelto
func (c *PeerConfig) GetTOS() TypeOfService {
	return c.tos
}
