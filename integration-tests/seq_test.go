package integrationtests

import (
	"testing"
	"time"

	"gitlab.com/tibwere/comunigo/proto"
	gp "google.golang.org/protobuf/proto"
)

// Test d'integrazione per l'invio singolo (sequencer)
func TestSingleSendSequencer(t *testing.T) {
	sendSequencer(t, false)
}

// Test d'integrazione per l'invio multiplo (sequencer)
func TestMultipleSendSequencer(t *testing.T) {
	sendSequencer(t, true)
}

// Test sul funzionamento del multicast basato sull'uso del sequencer
func sendSequencer(t *testing.T, parallel bool) {
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

	ml, err := retrieveMessagesSEQ(users)
	if err != nil {
		t.Fatalf("Unable to retrieve messages (%v)", err)
	}

	ref := ml[users[0].GetName()]
	for _, u := range users[1:] {
		actual := ml[u.GetName()]
		for i := range actual {
			if !gp.Equal(ref[i], actual[i]) {
				t.Fatalf("%v-th message for %v: %v | %v-th message for %v: %v",
					i, u.GetName(), actual[i], i, users[0].GetName(), ref[i])
			}
		}
	}
}

// Funzione che permette di fare il retrieve di messaggi ricevuti dal sequencer
func retrieveMessagesSEQ(users []*User) (map[string][]*proto.SequencerMessage, error) {
	var ml = make(map[string][]*proto.SequencerMessage)
	var err error

	for _, u := range users {
		ml[u.GetName()], err = u.GetMessagesSEQ()
		if err != nil {
			return ml, err
		}
	}

	return ml, err
}
