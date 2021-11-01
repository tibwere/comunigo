// Package contenente alcuni test d'integrazione
// volti a verificare la corretta implementazione degli
// algoritmi di multicast
package integrationtests

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"
)

// Alcuni parametri fixed su cui si basano i test
const (
	N_MESSAGES_FOR_PEER  = 5
	START_DELAY_INTERVAL = 0
	END_DELAY_INTERVAL   = 5
)

// Funzione wrapper per la registrazione del gruppo di multicast
func Registration() ([]*User, error) {
	users := []*User{}

	ports, err := getPorts()
	if err != nil {
		return users, err
	}

	users = generateUsers(ports)

	eg, _ := errgroup.WithContext(context.Background())
	for _, u := range users {
		currUser := u
		eg.Go(func() error {
			return currUser.Sign()
		})
	}

	err = eg.Wait()
	return users, err
}

// Funzione wrapper per l'invio dei messaggi da parte dei peer afferenti
// al gruppo in parallelo o meno
func SendMessages(users []*User, parallel bool, start int, end int) error {
	if parallel {
		eg, _ := errgroup.WithContext(context.Background())

		for i := 0; i < N_MESSAGES_FOR_PEER; i++ {
			for _, u := range users {
				currUser := u
				eg.Go(func() error {
					return currUser.SendMessage(fmt.Sprintf("Message from %v", currUser.GetName()), start, end)
				})
			}
		}

		return eg.Wait()
	} else {
		for i := 0; i < N_MESSAGES_FOR_PEER; i++ {
			for _, u := range users {
				if err := u.SendMessage(fmt.Sprintf("Message from %v", u.GetName()), start, end); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// Funzione che permette di effettuare il retrieve delle porte da
// contattare a partire dall'ambiente
func getPorts() ([]uint16, error) {
	ports := []uint16{}

	portsStr, ok := os.LookupEnv("COMUNIGO_TEST_PORTS")
	if !ok {
		return ports, fmt.Errorf("no ports found")
	} else {
		for _, pStr := range strings.Split(portsStr, ",") {
			p, err := strconv.ParseUint(pStr, 10, 16)
			if err != nil {
				return ports, fmt.Errorf("unable to parse %v", pStr)
			}

			ports = append(ports, uint16(p))
		}
	}

	return ports, nil
}

// Funzione che permette di generare un pool di utenti
// da utilizzare nei casi di test
func generateUsers(ports []uint16) []*User {
	var generated []*User
	for i := range ports {
		generated = append(generated, &User{
			name: fmt.Sprintf("test-%v", i),
			port: ports[i],
		})
	}
	return generated
}
