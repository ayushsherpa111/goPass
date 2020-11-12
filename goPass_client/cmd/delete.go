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
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var (
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a ",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == 0 {
				return errors.New("Id required")
			}
			criter := kit.MatchPass{}
			criter.Key, _ = os.UserHomeDir()
			criter.Tbl = "passwords"
			criter.KeyVal = map[string]interface{}{"pid": id}
			if count, e := rpcCon.RPC_CLIENT.DeleteItem(criter); e != nil {
				return e
			} else if count > 0 {
				fmt.Println(count, "Password(s) deleted")
			}
			return nil
		},
	}
	id int
)

func init() {
	deleteCmd.Flags().IntVarP(&id, "id", "i", 0, "The id of the password in the DataBase")
	rootCmd.AddCommand(deleteCmd)
}
