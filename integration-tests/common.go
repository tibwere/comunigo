package integrationtests

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"
)

const (
	N_MESSAGES_FOR_PEER = 5
)

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

func SendMessages(users []*User, parallel bool) error {
	if parallel {
		eg, _ := errgroup.WithContext(context.Background())

		for i := 0; i < N_MESSAGES_FOR_PEER; i++ {
			for _, u := range users {
				currUser := u
				eg.Go(func() error {
					return currUser.SendMessage(fmt.Sprintf("Message from %v", currUser.Name))
				})
			}
		}

		return eg.Wait()
	} else {
		for i := 0; i < N_MESSAGES_FOR_PEER; i++ {
			for _, u := range users {
				if err := u.SendMessage(fmt.Sprintf("Message from %v", u.Name)); err != nil {
					return err
				}

				//time.Sleep(5000 * time.Millisecond)
			}
		}
		return nil
	}
}

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

func generateUsers(ports []uint16) []*User {
	var generated []*User
	for i := range ports {
		generated = append(generated, &User{
			Name: fmt.Sprintf("test-%v", i),
			Port: ports[i],
		})
	}
	return generated
}
