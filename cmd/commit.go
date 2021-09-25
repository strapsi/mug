/*

Package cmd : commit

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
	"fmt"
	"mug/mp"

	"github.com/spf13/cobra"
)

var addAll bool
var overrideType string

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:     "commit",
	Aliases: mp.Aliases["commit"],
	Short:   "beng commit helper",
	Long:    `add correct commit tags for beng repositories`,
	Run: func(cmd *cobra.Command, args []string) {
		commitCommand(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().BoolVarP(&addAll, "add", "a", false, "add files prior to commit")
	commitCmd.Flags().StringVarP(&overrideType, "type", "t", "", "override type of commit")
}

func commitCommand(cmd *cobra.Command, args []string) {
	fmt.Println("running git commit")
	add, commit := mp.FrdCommit(args, mp.CurrentGitBranch(), addAll, overrideType)
	if len(add) > 0 {
		mp.Exec(add)
	}
	mp.Exec(commit)
}
