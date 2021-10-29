package utilities

import (
	"fmt"
	"os"
	"strings"
)

// In ottica OO, oggetto mantenente i
// metadati per la configurazione del
// nodo di registrazione
type RegistrationServiceConfig struct {
	regPort       uint16
	chatGroupSize uint16
	seqAddress    string
	tos           TypeOfService
}

// "Costruttoure" dell'oggetto RegistrationServerConfig
func InitRegistrationServiceConfig() (*RegistrationServiceConfig, error) {
	c := &RegistrationServiceConfig{}

	val, err := parseUint16FromEnv(EnvRegPort)
	if err != nil {
		return c, err
	}
	c.regPort = val

	val, err = parseUint16FromEnv(EnvSize)
	if err != nil {
		return c, err
	}
	c.chatGroupSize = val

	rhost, isPresent := os.LookupEnv(EnvSeqHostname)
	if !isPresent {
		return c, ErrEnvNotFound
	}
	c.seqAddress = rhost

	tos, isPresent := os.LookupEnv(EnvTypeOfService)
	if !isPresent {
		return c, fmt.Errorf("%v [TOS]", ErrEnvNotFound)
	} else {
		c.tos = setTOS(strings.ToLower(tos))
	}

	return c, nil
}

// "Metodo della classe RegistrationServiceConfig" per il retrieve della porta
// esposta in fase di registrazione
func (c *RegistrationServiceConfig) GetExposedPort() uint16 {
	return c.regPort
}

// "Metodo della classe RegistrationServiceConfig" per il retrieve della
// dimensione del gruppo di multicast
func (c *RegistrationServiceConfig) GetMulticastGroupSize() uint16 {
	return c.chatGroupSize
}

// "Metodo della classe RegistrationServiceConfig" per il retrieve dell'indirizzo
// del sequencer
func (c *RegistrationServiceConfig) GetSequencerAddress() string {
	return c.seqAddress
}

// "Metodo della classe RegistrationServiceConfig" per il retrieve del
// tipo di servizio scelto
func (c *RegistrationServiceConfig) GetTOS() TypeOfService {
	return c.tos
}
