package integrationtests

import (
	"fmt"
	"testing"
	"time"
)

func TestSingleSendSequencer(t *testing.T) {

	users, err := RegistrationHandler()
	if err != nil {
		t.Fatalf("Unable to sign test peers (%v)", err)
	}

	for _, u := range users {
		if err = u.SendMessage(fmt.Sprintf("Message from %v", u.Name)); err != nil {
			t.Fatalf("Unable to send messages (%v)", err)
		}
	}

	t.Log("Waiting for complete delivery of messages ...")
	time.Sleep(5 * time.Second)

	fail, passed := CompareMessageListsSEQ(users)
	if !passed {
		t.Fatalf(fail)
	}
}
