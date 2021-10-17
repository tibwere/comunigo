package utilities

import (
	"fmt"
	"log"
	"os"
)

func InitLogger(name string) error {
	logFile, err := os.OpenFile(
		fmt.Sprintf("/logs/%v.log", name),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		return err
	}

	log.SetOutput(logFile)
	return nil
}
