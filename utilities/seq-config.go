package utilities

import (
	"fmt"
	"os"
	"strings"
)

// In ottica OO, oggetto mantenente i
// metadati per la configurazione del
// sequencer
type SequencerConfig struct {
	chatPort      uint16
	regPort       uint16
	chatGroupSize uint16
	tos           TypeOfService
}

// "Costruttore" dell'oggetto SequencerServerConfig
func InitSequencerConfig() (*SequencerConfig, error) {
	c := &SequencerConfig{}

	val, err := parseUint16FromEnv(EnvChatPort)
	if err != nil {
		return c, err
	}
	c.chatPort = val

	val, err = parseUint16FromEnv(EnvRegPort)
	if err != nil {
		return c, err
	}
	c.regPort = val

	val, err = parseUint16FromEnv(EnvSize)
	if err != nil {
		return c, err
	}
	c.chatGroupSize = val

	tos, isPresent := os.LookupEnv(EnvTypeOfService)
	if !isPresent {
		return c, fmt.Errorf("%v [TOS]", ErrEnvNotFound)
	} else {
		c.tos = setTOS(strings.ToLower(tos))
	}

	return c, nil
}

// "Metodo della classe SequencerConfig" per il retrieve della
// porta da usare per comunicare con i peer
func (c *SequencerConfig) GetToPeersPort() uint16 {
	return c.chatPort
}

// "Metodo della classe SequencerConfig" per il retrieve della
// porta da usare per comunicare con il nodo di registrazione
func (c *SequencerConfig) GetToRegistryPort() uint16 {
	return c.regPort
}

// "Metodo della classe SequencerConfig" per il retrieve della
// dimensione del gruppo di multicast
func (c *SequencerConfig) GetMulticastGroupSize() uint16 {
	return c.chatGroupSize
}

// "Metodo della classe SequencerConfig" per il retrieve del
// tipo di servizio scelto
func (c *SequencerConfig) GetTOS() TypeOfService {
	return c.tos
}
