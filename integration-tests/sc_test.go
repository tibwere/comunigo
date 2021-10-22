package integrationtests

import (
	"testing"
	"time"

	"gitlab.com/tibwere/comunigo/proto"
	gp "google.golang.org/protobuf/proto"
)

func TestSingleSendScalar(t *testing.T) {
	sendScalar(t, false)
}

func TestMultipleSendScalar(t *testing.T) {
	sendScalar(t, true)
}

func sendScalar(t *testing.T, parallel bool) {
	users, err := Registration()
	if err != nil {
		t.Fatalf("Unable to sign test peers (%v)", err)
	}

	time.Sleep(3 * time.Second)
	t.Log("Inizio l'invio")

	err = SendMessages(users, parallel)
	if err != nil {
		t.Fatalf("Unable to send messages (%v)", err)
	}

	t.Log("Waiting for complete delivery of messages ...")
	time.Sleep(3 * time.Second)

	ml, err := retrieveMessagesSC(users)
	if err != nil {
		t.Fatalf("Unable to retrieve messages (%v)", err)
	}

	length := getMaxCommonIndex(users, ml)
	ref := ml[users[0].Name]
	for _, u := range users[1:] {
		actual := ml[u.Name]
		for i := 0; i < length; i++ {
			if !gp.Equal(ref[i], actual[i]) {
				t.Fatalf("%v-th message for %v: %v | %v-th message for %v: %v",
					i, u.Name, actual[i], i, users[0].Name, ref[i])
			}
		}
	}

}

func getMaxCommonIndex(users []*User, ml map[string][]*proto.ScalarClockMessage) int {
	curr := len(ml[users[0].Name])
	for _, u := range users[1:] {
		if len(ml[u.Name]) < curr {
			curr = len(ml[u.Name])
		}
	}

	return curr
}

func retrieveMessagesSC(users []*User) (map[string][]*proto.ScalarClockMessage, error) {
	var ml = make(map[string][]*proto.ScalarClockMessage)
	var err error

	for _, u := range users {
		ml[u.Name], err = u.GetMessagesSC()
		if err != nil {
			return ml, err
		}
	}

	return ml, err
}
