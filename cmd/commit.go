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
	"github.com/spf13/cobra"
	"mug/mp"
	"strings"
)

var addAll bool
var overrideType string

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "beng commit helper",
	Long:  `add correct commit tags for beng repositories`,
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
	output := mp.ExecEnvResult([]string{"git", "branch", "--show-current"}, []string{})
	branch := fmt.Sprintf("%s", output)
	message := composeMessage(branch, overrideType, args)
	var command []string
	if addAll {
		command = []string{"git", "add", "."}
		mp.Exec(command)
	}
	command = []string{"git", "commit", "-m", message}
	mp.Exec(command)
}

func composeMessage(branch string, overrideType string, args []string) string {
	if len(args) < 1 {
		mp.ExitWithError("wrong number of arguments. expecting message")
	}
	commitType, id := parseBranch(branch, overrideType)
	command := fmt.Sprintf("[%s][FRD-%s] %s\n", commitType, id, args[0])
	return command
}

func parseCommitType(commitType string) string {
	switch strings.ToLower(commitType) {
	case "feature", "f":
		return "FEATURE"
	case "refactor", "r":
		return "REFACTOR"
	case "intern", "i":
		return "INTERN"
	case "style", "s":
		return "STYLE"
	case "bugfix", "b":
		return "BUGFIX"
	default:
		mp.ExitWithError("unknown commit type " + commitType)
	}
	return ""
}

func parseBranch(branch string, overrideType string) (string, string) {
	parts := strings.Split(branch, "-")
	if len(parts) < 3 {
		mp.ExitWithError("branch has wrong format <type>-<id>-cool-feature")
	}
	var returnType string
	if overrideType != "" {
		returnType = overrideType
	} else {
		returnType = parts[0]
	}
	return parseCommitType(returnType), parts[1]
}
