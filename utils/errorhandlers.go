package utils

import (
	"errors"
	"log"
	"os"
)

// ServerlessErr records errors from the Serverless binary
type ServerlessErr struct {
	Msg string
}

func (e *ServerlessErr) Error() string {
	return e.Msg
}

// Check is a util function to panic when there is an error.
func Check(e error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("runtime panic : %v", err)
		}
	}()

	if e != nil {
		log.Printf("%v", e)
		erro := errors.New("Error happened during execution, please type 'wskdeploy -h' for help messages.")
		log.Printf("%v", erro)
		os.Exit(1)

	}

}
