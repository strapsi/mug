/*

Package cmd : hash command

Copyright © 2021 m.vondergruen@protonmail.com

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
	"crypto/sha512"
	"fmt"

	"github.com/spf13/cobra"

	// "golang.design/x/clipboard"
	"mug/mp"
)

// hashCmd represents the hash command
var hashCmd = &cobra.Command{
	Use:     "hash",
	Aliases: mp.Aliases["hash"],
	Short:   "hash the input string",
	Long:    `returns a sha512 hash of the given string`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			mp.ExitWithError("no string input provided")
		}
		hashCommand(args)
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
}

func hashCommand(args []string) {
	sha512Hash := sha512.New()
	sha512Hash.Write([]byte(args[0]))
	hashString := fmt.Sprintf("%x", sha512Hash.Sum(nil))
	// clipboard.Write(clipboard.FmtText, []byte(hashString))
	// clipboard.Read(clipboard.FmtText)
	fmt.Println(hashString)
}
