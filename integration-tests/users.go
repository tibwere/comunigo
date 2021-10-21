package integrationtests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"gitlab.com/tibwere/comunigo/proto"
	"golang.org/x/sync/errgroup"
)

type User struct {
	Name string
	Port uint16
}

func GenerateUsers(number int, ports []uint16) []*User {
	var generated []*User
	for i := 0; i < number; i++ {
		generated = append(generated, &User{
			Name: fmt.Sprintf("test-%v", i),
			Port: ports[i],
		})
	}
	return generated
}

func (u *User) SendMessage(body string) error {
	params := url.Values{}
	params.Set("message", body)

	response, err := http.PostForm(
		fmt.Sprintf("http://localhost:%v/send", u.Port),
		params,
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

func (u *User) GetMessagesSEQ() ([]*proto.SequencerMessage, error) {
	res := []*proto.SequencerMessage{}

	response, err := http.Get(fmt.Sprintf("http://localhost:%v/list", u.Port))
	if err != nil {
		return res, err
	}
	defer response.Body.Close()

	d := json.NewDecoder(response.Body)
	err = d.Decode(&res)
	return res, err
}

func (u *User) Sign() error {

	params := url.Values{}
	params.Set("username", u.Name)

	response, err := http.PostForm(
		fmt.Sprintf("http://localhost:%v/sign", u.Port),
		params,
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	res := struct {
		Status  string
		Message string
	}{}

	d := json.NewDecoder(response.Body)
	if err = d.Decode(&res); err != nil {
		return err
	}

	if res.Status == "SUCCESS" {
		return nil
	} else {
		return fmt.Errorf("User '%v' error message: %v", u.Name, res.Message)
	}
}

func SignUsersParallel(users []*User) error {

	eg, _ := errgroup.WithContext(context.Background())
	for _, u := range users {
		currUser := u
		eg.Go(func() error {
			return currUser.Sign()
		})
	}

	return eg.Wait()
}

func SendMessagesParallel(users []*User) error {
	eg, _ := errgroup.WithContext(context.Background())
	for _, u := range users {
		currUser := u
		eg.Go(func() error {
			return currUser.SendMessage(fmt.Sprintf("Message from %v", currUser.Name))
		})
	}

	return eg.Wait()
}
