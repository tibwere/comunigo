// Package contenente alcune funzionalità di utility
// comuni ai vari componenti dell'architettura come:
// retrieve della configurazione dall'ambiente,
// configurazione dell'attività di logging su file
// inizilizzazione del contesto cancellabile all'arrivo
// di un segnale tra SIGINT e SIGTERM
package utilities

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// Prefisso del messaggio d'errore da restituire nel caso in cui
// la variabile d'ambiente cercata non è attualmente presente
var (
	ErrEnvNotFound = errors.New("environment variable unset")
)

// Ridefinizione di un tipo di dato intero da
// poter utilizzare a mo' di enumerato per
// discriminare la scelta del tipo di servizio da offrire
type TypeOfService uint8

const (
	TOS_CS_SEQUENCER  TypeOfService = 0
	TOS_P2P_SCALAR    TypeOfService = 1
	TOS_P2P_VECTORIAL TypeOfService = 2
	TOS_INVALID       TypeOfService = 3
)

// Variabili d'ambiente da utilizzare nella configurazione
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

// Funzione di utility che permette di estrarre un intero unsigned a 16 bit
// a partire dalla stringa della variabile d'ambiente corrispondente
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

// Funzione che permette di settare il tipo di servizio
// a partire dalla variabile d'ambiente corrispondente
func setTOS(env string) TypeOfService {
	switch env {
	case "sequencer":
		return TOS_CS_SEQUENCER
	case "scalar":
		return TOS_P2P_SCALAR
	case "vectorial":
		return TOS_P2P_VECTORIAL
	default:
		return TOS_INVALID
	}
}

func (t TypeOfService) ToString() string {
	switch t {
	case TOS_CS_SEQUENCER:
		return "SEQUENCER"
	case TOS_P2P_SCALAR:
		return "SCALAR"
	case TOS_P2P_VECTORIAL:
		return "VECTORIAL"
	default:
		return "INVALID"
	}
}
