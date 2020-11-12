/*
Copyright © 2020 NAME HERE ayush20100@hotmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ayushsherpa111/goPass/kit"
	rpcCon "github.com/ayushsherpa111/goPass/rpc_util"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var (
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get your passwords [Decrypted].",
		Long: `Generates a table containing your passwords that you have provided or want to search.
The Passwords will be in plain text so use at your own risk.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fuzzyMatch.Key, _ = os.UserHomeDir()
			fuzzyMatch.Tbl = "passwords"
			var comms = make([]kit.Password, 0)
			var e error = errors.New("")
			var pass []byte
			for {
				if allFlag {
					fuzzyMatch.KeyVal = make(map[string]interface{})
				} else {
					keys := []string{"pid", "Username", "Email", "Site"}
					fuzzyMatch.KeyVal = kit.ZipMap(keys, searchVals, basePattern)
					// fuzzyMatch.Pattern = []interface{}{basePattern, basePattern, basePattern, basePattern, basePattern}
				}
				e = rpcCon.RPC_CLIENT.GetPass(fuzzyMatch, &comms)
				if e == nil {
					break
				}
				// probably an error about master password missing
				pass, e = kit.GetPass("Enter Master Password  ")
				e = rpcCon.RPC_CLIENT.ReEnterPassword(kit.ReEnter{Path: fuzzyMatch.Key, Password: pass})
				if e != nil {
					fmt.Println(e.Error())
					break
				}
			}
			Printtable(comms)
			return nil
		},
	}

	fuzzyMatch  kit.MatchPass
	allFlag     bool
	basePattern string
	pid         int
	username    string
	email       string
	site        string
	searchVals  = []interface{}{&pid, &username, &email, &site}
)

func init() {
	getCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Use this flag to get all of your passwords.")
	getCmd.Flags().StringVarP(&basePattern, "match", "b", "", "A base pattern that matches username, email, and website if not specified individually.")
	getCmd.Flags().IntVarP(&pid, "id", "i", 0, "Id of a password in the database.")
	getCmd.Flags().StringVarP(&username, "uname", "u", "", "Pattern to match username [basePattern will be used if not provided]")
	getCmd.Flags().StringVarP(&email, "email", "e", "", "Pattern to match Email [basePattern will be used if not provided]")
	getCmd.Flags().StringVarP(&site, "site", "s", "", "Pattern to match Website [basePattern will be used if not provided]")
	rootCmd.AddCommand(getCmd)
}

func Printtable(data []kit.Password) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	// t.SetAutoIndex(true)
	t.AppendHeader(table.Row{"Pid", "Email", "Username", "Password", "Site"})
	for _, v := range data {
		t.Style().Color.IndexColumn = append(t.Style().Color.IndexColumn, text.FgWhite, text.FgWhite, text.FgRed, text.FgWhite)
		t.AppendRow(table.Row{v.Pid, v.Email, v.Username, string(v.Password), v.Site})
	}
	t.AppendFooter(table.Row{0, "", "", "Total:", len(data)})
	t.SetStyle(table.StyleDouble)
	t.Render()
}
