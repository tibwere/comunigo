package integrationtests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"gitlab.com/tibwere/comunigo/proto"
)

// In ottica OO, oggetto che rappresenta
// Un utente che utilizza l'applicazione comuniGO
type User struct {
	name string
	port uint16
}

// "Metodo della classe User" che permette di effettuare
// il retrieve del nome dell'utente corrente
func (u *User) GetName() string {
	return u.name
}

// "Metodo della classe User" che permette di inviare
// un messaggio sfruttando la route "/send"
func (u *User) SendMessage(body string, start int, end int) error {
	params := url.Values{}
	params.Set("message", body)
	params.Set("delay", fmt.Sprintf("%v:%v", start, end))

	response, err := http.PostForm(
		fmt.Sprintf("http://localhost:%v/send", u.port),
		params,
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// "Metodo della classe User" che permette di effettuare il retrive
// dei messaggi sfruttando la route "/list (sequencer)"
func (u *User) GetMessagesSEQ() ([]*proto.SequencerMessage, error) {
	res := []*proto.SequencerMessage{}

	response, err := http.Get(fmt.Sprintf("http://localhost:%v/list", u.port))
	if err != nil {
		return res, err
	}
	defer response.Body.Close()

	d := json.NewDecoder(response.Body)
	err = d.Decode(&res)
	return res, err
}

// "Metodo della classe User" che permette di effettuare il retrive
// dei messaggi sfruttando la route "/list (scalar)"
func (u *User) GetMessagesSC() ([]*proto.ScalarClockMessage, error) {
	res := []*proto.ScalarClockMessage{}

	response, err := http.Get(fmt.Sprintf("http://localhost:%v/list", u.port))
	if err != nil {
		return res, err
	}
	defer response.Body.Close()

	d := json.NewDecoder(response.Body)
	err = d.Decode(&res)
	return res, err
}

// "Metodo della classe User" che permette di effettuare il retrive
// dei messaggi sfruttando la route "/list (vectorial)"
func (u *User) GetMessagesVC() ([]*proto.VectorialClockMessage, error) {
	res := []*proto.VectorialClockMessage{}

	response, err := http.Get(fmt.Sprintf("http://localhost:%v/list", u.port))
	if err != nil {
		return res, err
	}
	defer response.Body.Close()

	d := json.NewDecoder(response.Body)
	err = d.Decode(&res)
	return res, err
}

// "Metodo della classe User" che permette di registrare
// un utente sfruttando la route "/sign"
func (u *User) Sign() error {
	params := url.Values{}
	params.Set("username", u.name)

	response, err := http.PostForm(
		fmt.Sprintf("http://localhost:%v/sign", u.port),
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
		return fmt.Errorf("User '%v' error message: %v", u.name, res.Message)
	}
}
