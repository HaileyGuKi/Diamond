package utils

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func HandleErr(message string, err error) {
	if err != nil {
		logrus.Errorf(fmt.Sprintf("%s: %v", message, err))
		os.Exit(1)
	}
}
