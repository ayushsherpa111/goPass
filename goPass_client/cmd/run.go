package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

type Test struct {
	Key  string
	Pass []byte
}

var (
	rootCmd = &cobra.Command{
		Use:   "goPass",
		Short: "A Password manager for the CLI",
		Long:  "goPass is a password manager that runs entirely on your CLI",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("run")
			// t := Test{}
			// t.Key, _ = os.UserHomeDir()
			// t.Pass = []byte("mypassword")
			// conn, _ := rpc.DialHTTP("tcp", ":5555")
			// var reply kit.User
			// conn.Call("API.TestRoute", t, &reply)
			// fmt.Println("the reply ", reply)
		},
	}
)

func init() {
	rootCmd.AddCommand(InitCmd)
	rootCmd.AddCommand(AddCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
