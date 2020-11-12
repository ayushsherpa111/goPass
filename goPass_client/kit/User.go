package kit

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	CONF_DIR, _ = os.UserHomeDir()
)

type User struct {
	Username string
	Email    string
	HomePath string
	Key      []byte
}

func GetPass(msg string) ([]byte, error) {
	fmt.Print(msg)
	if bytePass, err := terminal.ReadPassword(0); err != nil {
		return []byte(""), err
	} else {
		if len(bytePass) <= 8 {
			return []byte(""), errors.New("Password must contain atleast 8 characters")
		}
		// put cursor in the next line
		fmt.Println()
		return bytePass, nil
	}
}
