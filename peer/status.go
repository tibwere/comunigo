// Package top-level per la gestione della logica
// del componente peer dell'applicazione comuiniGO
package peer

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/go-redis/redis/v8"
	"gitlab.com/tibwere/comunigo/proto"
)

// In ottica OO, oggetto che racchiude i metadati
// basici che identificano lo stato di un peer all'interno
// del sistema
type Status struct {
	currentUsername string
	otherMembers    []*proto.PeerInfo
	datastore       *redis.Client
	frontBackCh     chan string
	exposedIP       string
}

// "Costruttore" dell'ogetto Status
func Init(redisAddr string) (*Status, error) {
	ip, err := retrieveIP()
	if err != nil {
		return nil, err
	} else {
		s := &Status{
			currentUsername: "",
			otherMembers:    []*proto.PeerInfo{},
			frontBackCh:     make(chan string),
			exposedIP:       ip,
		}

		s.initDatastore(redisAddr)
		return s, nil
	}
}

// "Metodo della classe Status" che permette di conoscere l'IP del
// nodo su cui è in esecuzione il peer
func (s *Status) GetExposedIP() string {
	return s.exposedIP
}

// "Metodo della classe Status" che permette di conoscere
// la lista dei peer connessi al gruppo di multicast
func (s *Status) GetOtherMembers() []*proto.PeerInfo {
	return s.otherMembers
}

// "Metodo della classe Status" che permette di effettuare il retrieve
// di uno specifico peer all'interno della lista di cui si specifica l'indice
func (s *Status) GetSpecificMember(index int) *proto.PeerInfo {
	return s.otherMembers[index]
}

// "Metodo della classe Status" che permette di settare l'username
// nel caso in cui esso non sia stato già settato in precedenza
// altrimenti viene generato un errore
func (s *Status) SetUsername(username string) error {
	if s.currentUsername != "" {
		return fmt.Errorf("unable to set twice username")
	} else {
		s.currentUsername = username
		return nil
	}
}

// "Metodo della classe Status" che permette di conoscere
// l'username corrente
func (s *Status) GetCurrentUsername() string {
	return s.currentUsername
}

// "Metodo della classe Status" che permette di verificare
// se sul peer è avvenuta già la registrazione o meno
func (s *Status) NotYetSigned() bool {
	return s.currentUsername == ""
}

// "Metodo della classe Status" che permette di prelevare messaggi
// dal canale adibito alla comunicazione frontend/backend
func (s *Status) GetFromFrontendBackendChannel() <-chan string {
	return s.frontBackCh
}

// "Metodo della classe Status" che permette di inserire messaggi
// nel canale adibito alla comunicazione frontend/backend
func (s *Status) PushIntoFrontendBackendChannel(message string) {
	s.frontBackCh <- message
}

// "Metodo della classe Status" che permette di aggiungere un nuovo
// membro alla lista dei partecipanti al gruppo di multicast
func (s *Status) InsertNewMember(newMember *proto.PeerInfo) {
	s.otherMembers = append(s.otherMembers, newMember)
}

// Funzione che permette di effettuare il retrieve dell'IP
// del nodo su cui deve andare in startup il peer corrente
func retrieveIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if !strings.Contains(addr.String(), "127.0.0.1") {
			return strings.Split(addr.String(), "/")[0], nil
		}
	}

	return "", errors.New("no public IP addresses found")
}
