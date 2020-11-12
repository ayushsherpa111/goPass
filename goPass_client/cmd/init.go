package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"

	"github.com/ayushsherpa111/goPass/kit"
	rpcCon "github.com/ayushsherpa111/goPass/rpc_util"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

const (
	ROUNDS   = 1400
	KEY_LEN  = 32
	SALT_LEN = 16
	CONFIG   = "config.json"
	VAULT    = ".vault.db"
)

var (
	home_dir, _ = os.UserHomeDir()
	PATH        = path.Join(home_dir, ".goPass")
	InitCmd     = &cobra.Command{
		Use:     "init",
		Short:   "Initialize your Repo",
		PreRunE: PassPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Send the user data to the RPC
			// add the users home path in the object
			currentUser.HomePath, _ = os.UserHomeDir()
			if u, err := rpcCon.RPC_CLIENT.Init_User(currentUser); err != nil {
				log.Fatalln(err.Error())
			} else {
				fmt.Print("\nPassword Manager Initialized\n")
				PrintHello(u)
			}
			return nil
		},
	}
	currentUser  = kit.User{}
	key          string
	configFile   *os.File
	email_REG, _ = regexp.Compile("[a-zA-Z0-9!_]+\\@[a-zA-Z0-9]+(\\.[a-z]{2,}){1,}")
)

func PassPreRun(cmd *cobra.Command, args []string) error {
	// check for username
	if len(currentUser.Username) < 3 {
		return errors.New("Invalid Username. Must be atleast 3 characters")
	}

	// check for email
	if !email_REG.MatchString(currentUser.Email) {
		return errors.New("Invalid Email. Make sure the email is valid")
	}

	var err error
	// get plain text Password
	currentUser.Key, err = kit.GetPass("Enter your Password [Minimum 8 Characters Long]: ")
	if err != nil {
		return err
	}
	// load data into the config file
	// currentUser.GenerateHash()

	return nil
}

func init() {
	InitCmd.Flags().StringVarP(&currentUser.Username, "username", "u", "", "The username that you use for the accounts")
	InitCmd.Flags().StringVarP(&currentUser.Email, "email", "E", "", "Your Email address a particular account")
}

func PrintHello(u kit.User) {
	welcomeFig := figure.NewColorFigure("Initialized Vault", "speed", "red", false)
	welcomeFig.Print()
	fmt.Printf(`
Welcome %s, Your Vault has been succesfully set up. 
You are now ready to use the CLI based password manager
`, u.Username)
}
