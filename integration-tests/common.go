package integrationtests

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TestConfiguration struct {
	NumberOfUsers int
	Ports         []uint16
}

func GetConfiguration() (*TestConfiguration, error) {

	tc := &TestConfiguration{
		NumberOfUsers: 0,
		Ports:         []uint16{},
	}

	portsStr, ok := os.LookupEnv("COMUNIGO_TEST_PORTS")
	if !ok {
		return tc, fmt.Errorf("no ports found")
	} else {
		for _, pStr := range strings.Split(portsStr, ",") {
			p, err := strconv.ParseUint(pStr, 10, 16)
			if err != nil {
				return tc, fmt.Errorf("unable to parse %v", pStr)
			}

			tc.Ports = append(tc.Ports, uint16(p))
		}
	}

	tc.NumberOfUsers = len(tc.Ports)

	return tc, nil
}

func RegistrationHandler() ([]*User, error) {
	users := []*User{}

	tc, err := GetConfiguration()
	if err != nil {
		return users, err
	}

	users = GenerateUsers(tc.NumberOfUsers, tc.Ports)

	err = SignUsersParallel(users)
	if err != nil {
		return users, err
	}

	return users, nil
}

func CompareMessageListsSEQ(users []*User) (string, bool) {
	mlRef, err := users[0].GetMessagesSEQ()
	if err != nil {
		return fmt.Sprintf("[USER %v] Unable to retrieve messages (%v)", users[0].Name, err), false
	}
	for _, u := range users[1:] {
		mlActual, err := u.GetMessagesSEQ()
		if err != nil {
			return fmt.Sprintf("[USER %v] Unable to retrieve messages (%v)", u.Name, err), false
		}

		for i := range mlActual {
			refTimestamp := mlRef[i].Timestamp
			actualTimestamp := mlActual[i].Timestamp

			refFrom := mlRef[i].From
			actualFrom := mlActual[i].From

			refBody := mlRef[i].Body
			actualBody := mlActual[i].Body

			if refTimestamp != actualTimestamp || refFrom != actualFrom || refBody != actualBody {
				return fmt.Sprintf("User %v's %v-th message is: %v - User %v's %v-th message instead is: %v MISMATCH",
					users[0], i, mlRef[i], u, i, mlActual[i]), false
			}
		}
	}

	return "PASS", true
}
