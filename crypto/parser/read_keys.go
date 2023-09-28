package parser

import (
	"os"

	"github.com/sirupsen/logrus"
)

func ReadKeys(private_path string, public_path string) ([]byte, []byte, error) {
	var err error
	pubKey, err := os.ReadFile(public_path)
	if err != nil {
		logrus.Errorf("Error while reading signing keys : %v", err.Error())
		return nil, nil, err
	}
	privKey, err := os.ReadFile(private_path)
	if err != nil {
		logrus.Errorf("Error while reading signing keys : %v", err.Error())
		return nil, nil, err
	}

	return privKey, pubKey, nil
}
