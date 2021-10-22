package integrationtests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"gitlab.com/tibwere/comunigo/proto"
)

type User struct {
	Name string
	Port uint16
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

func (u *User) GetMessagesSC() ([]*proto.ScalarClockMessage, error) {
	res := []*proto.ScalarClockMessage{}

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
