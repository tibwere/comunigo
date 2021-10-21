package integrationtests

import (
	"testing"
	"time"
)

func TestMultipleSendSequencer(t *testing.T) {

	users, err := RegistrationHandler()
	if err != nil {
		t.Fatalf("Unable to sign test peers (%v)", err)
	}

	if err = SendMessagesParallel(users); err != nil {
		t.Fatalf("Unable to send messages (%v)", err)
	}

	t.Log("Waiting for complete delivery of messages ...")
	time.Sleep(5 * time.Second)

	fail, passed := CompareMessageListsSEQ(users)
	if !passed {
		t.Fatalf(fail)
	}
}
