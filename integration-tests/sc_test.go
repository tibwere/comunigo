package integrationtests

import (
	"testing"
	"time"

	"gitlab.com/tibwere/comunigo/proto"
	gp "google.golang.org/protobuf/proto"
)

// Test d'integrazione per l'invio singolo (scalar)
func TestSingleSendScalar(t *testing.T) {
	sendScalar(t, false)
}

// Test d'integrazione per l'invio multiplo (scalar)
func TestMultipleSendScalar(t *testing.T) {
	sendScalar(t, true)
}

// Test sul funzionamento del multicast basato sull'uso di clock logico scalare
func sendScalar(t *testing.T, parallel bool) {
	users, err := Registration()
	if err != nil {
		t.Fatalf("Unable to sign test peers (%v)", err)
	}

	err = SendMessages(users, parallel, START_DELAY_INTERVAL, END_DELAY_INTERVAL)
	if err != nil {
		t.Fatalf("Unable to send messages (%v)", err)
	}

	t.Log("Waiting for complete delivery of messages ...")
	time.Sleep(END_DELAY_INTERVAL * time.Second)

	ml, err := retrieveMessagesSC(users)
	if err != nil {
		t.Fatalf("Unable to retrieve messages (%v)", err)
	}

	length := getMaxCommonIndex(users, ml)
	ref := ml[users[0].GetName()]
	for _, u := range users[1:] {
		actual := ml[u.GetName()]
		for i := 0; i < length; i++ {
			if !gp.Equal(ref[i], actual[i]) {
				t.Fatalf("%v-th message for %v: %v | %v-th message for %v: %v",
					i, u.GetName(), actual[i], i, users[0].GetName(), ref[i])
			}
		}
	}

}

// Funzione che permette di effettuare il retrieve della lunghezza
// della sottolista di messaggi consegnati comune a tutti i peer
func getMaxCommonIndex(users []*User, ml map[string][]*proto.ScalarClockMessage) int {
	curr := len(ml[users[0].GetName()])
	for _, u := range users[1:] {
		if len(ml[u.GetName()]) < curr {
			curr = len(ml[u.GetName()])
		}
	}

	return curr
}

// Funzione che permette di fare il retrieve di messaggi basati su clock logico scalare
func retrieveMessagesSC(users []*User) (map[string][]*proto.ScalarClockMessage, error) {
	var ml = make(map[string][]*proto.ScalarClockMessage)
	var err error

	for _, u := range users {
		ml[u.GetName()], err = u.GetMessagesSC()
		if err != nil {
			return ml, err
		}
	}

	return ml, err
}
