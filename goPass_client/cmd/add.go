package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ayushsherpa111/goPass/kit"
	rpcCon "github.com/ayushsherpa111/goPass/rpc_util"
	"github.com/spf13/cobra"
)

var (
	AddCmd = &cobra.Command{
		Use:    "add",
		Short:  "Add a new password to your vault",
		Long:   "Add a new password for an application along with other details",
		PreRun: AddPre,
		Run:    AddPass,
		PostRun: func(cmd *cobra.Command, args []string) {
			var e error
			key, _ := os.UserHomeDir()
		loop:
			for {
				switch {
				case len(csvPath) > 0:
					e = SendCsv(csvPath, key)
					if e != nil {
						fmt.Println(e.Error())
					} else {
						break
					}
					break loop
				case len(csvPath) == 0:
					e = SendPass(newPassword, key)
					if e != nil {
						fmt.Println(e.Error())
					} else {
						fmt.Println("Password added")
						break loop
					}
				}
				if pass, e := kit.GetPass("Enter Master Password  "); e != nil {
					fmt.Println(e.Error())
					break
				} else {
					if e = rpcCon.RPC_CLIENT.ReEnterPassword(kit.ReEnter{Path: key, Password: pass}); e != nil {
						fmt.Println(e.Error())
						break
					}
				}
			}
		},
	}
	userDetails kit.User
	newPassword kit.Password
	csvPath     string
)

func SendCsv(csvPath string, key string) error {
	csvPath, _ = filepath.Abs(csvPath)
	csvPayload := kit.CSVData{CsvPath: csvPath, Key: key}
	return rpcCon.RPC_CLIENT.AddCSV(csvPayload)
}

func SendPass(pass kit.Password, key string) error {
	pass.Path = key
	return rpcCon.RPC_CLIENT.AddPassword(pass)
}

func init() {
	AddCmd.Flags().StringVarP(&csvPath, "csv", "f", "", "Import passwords from another password manager [CSV]")
}

func AddPre(cmd *cobra.Command, args []string) {
	// Get User Details first
	var err error
	if userDetails, err = rpcCon.RPC_CLIENT.Get_User_Info(kit.CONF_DIR); err != nil {
		log.Fatalln(err.Error())
	}
}

func AddPass(cmd *cobra.Command, args []string) {
	// prompt user for their passwords
	if len(csvPath) != 0 {
		return
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter the Username [%s] ", userDetails.Username)
	if text, _ := reader.ReadString(10); text != "\n" {
		newPassword.Username = strings.Trim(text, "\n")
	} else {
		newPassword.Username = userDetails.Username
	}
	fmt.Printf("Enter the Email address [%s] ", userDetails.Email)
	if text, _ := reader.ReadString(10); text != "\n" {
		newPassword.Email = strings.Trim(text, "\n")
	} else {
		newPassword.Email = userDetails.Email
	}
	fmt.Print("Enter the Website/Application this Password corresponds to  ")
	if text, _ := reader.ReadString(10); text != "\n" {
		newPassword.Site = strings.Trim(text, "\n")
	} else {
		log.Fatalln("Website/Application Required")
	}
	if pass, err := kit.GetPass("Enter Password For This Website  "); err != nil {
		log.Fatalln("Password Required", err.Error())
	} else {
		newPassword.Password = pass
	}
}
