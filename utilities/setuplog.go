package utilities

import (
	"fmt"
	"log"
	"os"
)

// Funzione che permette di inizializzare
// l'attività di logging su file anziché su STDIN
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
