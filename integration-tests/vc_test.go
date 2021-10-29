package integrationtests

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"gitlab.com/tibwere/comunigo/proto"
	"golang.org/x/sync/errgroup"
)

// Test d'integrazione per l'invio singolo (vectorial)
func TestSingleSendVectorial(t *testing.T) {
	sendVectorial(t, false)
}

// Test d'integrazione per l'invio multiplo (vectorial)
func TestMultipleSendVectorial(t *testing.T) {
	sendVectorial(t, true)
}

// Test sul funzionamento del multicast basato sull'uso del clock logico vettoriale
func sendVectorial(t *testing.T, parallel bool) {
	users, err := Registration()
	if err != nil {
		t.Fatalf("Unable to sign test peers (%v)", err)
	}

	eg, _ := errgroup.WithContext(context.Background())

	t.Log("Sending standard messages")
	for _, u := range users {
		if parallel {
			currentUser := u
			eg.Go(func() error {
				return sendStandardMessage(currentUser)
			})
		} else {
			if err = sendStandardMessage(u); err != nil {
				t.Fatalf(err.Error())
			}
		}
	}

	t.Log("Waiting for complete delivery of messages")
	if parallel {
		time.Sleep(500 * time.Millisecond)
	} else {
		time.Sleep(100 * time.Millisecond)
	}

	t.Log("Send summary messages")
	for _, u := range users {
		if parallel {
			currentUser := u
			eg.Go(func() error {
				return sendSummaryMessage(currentUser)
			})

		} else {
			if err := sendSummaryMessage(u); err != nil {
				t.Fatalf(err.Error())
			}
		}
	}

	t.Log("Waiting for complete delivery of summary messages")
	time.Sleep(3 * time.Second)

	t.Log("Retrieve updated message list")
	for _, u := range users {
		if parallel {
			currentUser := u
			eg.Go(func() error {
				return verifyIfCorrect(currentUser)
			})
		} else {
			if err := verifyIfCorrect(u); err != nil {
				t.Fatalf(err.Error())
			}
		}
	}
}

// Funzione che permette di inviare, per conto di un determinato utente, un messaggio "standard"
// ovvero contenente il nome del mittente
func sendStandardMessage(u *User) error {
	if err := u.SendMessage(u.GetName(), START_DELAY_INTERVAL, END_DELAY_INTERVAL); err != nil {
		return fmt.Errorf("Unable to send messages (%v)", err)
	}
	return nil
}

// Funzione che permette di inviare, per conto di un determinato utente, un messaggio "riassuntivo"
// ovvero contenente l'elenco dei messaggi visti separati da :
func sendSummaryMessage(u *User) error {
	ml, err := u.GetMessagesVC()
	if err != nil {
		return fmt.Errorf("Unable to retrieve messages (%v)", err)
	}

	summaryBody := summaryOfRetrievedMessages(ml)
	if err = u.SendMessage(summaryBody, START_DELAY_INTERVAL, END_DELAY_INTERVAL); err != nil {
		return fmt.Errorf("Unable to send messages (%v)", err)
	}

	return nil
}

// Funzione wrapper per la verifica della correttezza
func verifyIfCorrect(u *User) error {
	ml, err := u.GetMessagesVC()
	if err != nil {
		return fmt.Errorf("Unable to retrieve messages (%v)", err)
	}

	return checkCausalConsistency(ml)
}

// Funzione che costruisce il corpo del messaggio riassuntivo
func summaryOfRetrievedMessages(ml []*proto.VectorialClockMessage) string {
	newBody := ""

	for _, mess := range ml {
		// sono interessato unicamente ai messaggi prima dei riassunti
		// perché su quelli voglio testare l'effettivo rispetto della causalità
		if len(strings.Split(mess.GetBody(), ":")) == 1 {
			newBody += mess.GetBody() + ":"
		}
	}

	if len(newBody) > 0 {
		if newBody[len(newBody)-1] == ':' {
			return newBody[:len(newBody)-1]
		}
	}

	return newBody
}

// Funzione che verifica se la lista di messaggi ricevuta come parametro
// rispetta l'ordinamento causale
func checkCausalConsistency(ml []*proto.VectorialClockMessage) error {
	before := []string{}

	for _, mess := range ml {
		parts := strings.Split(mess.GetBody(), ":")
		if len(parts) > 1 {
			// Per la consistenza causale questi messaggi possono essere
			// ricevuti in un ordine arbitrario dai vari peer, l'importante
			// è che precedano il messaggio riassuntivo
			sort.Strings(before)
			sort.Strings(parts)

			if len(before) != len(parts) {
				return fmt.Errorf("Message seen before are: %v but in the summary there is %v", before, parts)
			}

			for i := range before {
				if before[i] != parts[i] {
					return fmt.Errorf("Message seen before are: %v but in the summary there is %v", before, parts)
				}
			}
		} else {
			before = append(before, parts[0])
		}
	}

	return nil
}
