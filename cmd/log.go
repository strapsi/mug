/*
Package cmd : log

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
	"github.com/strapsi/go-docker"
	"mug/mp"
	"regexp"
)

var gitLogFormat = "%C(yellow)%h%C(reset) %C(auto)%d%C(reset) %s %C(blue)%cr%C(reset) by %C(green)%cn%C(reset)"

var dockerLogName string
var logLimit string
var dockerLogFollow bool
var withGitFileNames bool

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:     "log",
	Aliases: []string{"l"},
	Short:   "output log",
	Long:    `outputs log`,
	Run: func(cmd *cobra.Command, args []string) {
		logCommand(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().StringVarP(&dockerLogName, "docker", "d", "", "name of container to log. can be fuzzy")
	logCmd.Flags().StringVarP(&logLimit, "limit", "l", "10", "only print last x logs")
	logCmd.Flags().BoolVarP(&dockerLogFollow, "follow", "f", false, "follow log output")
	logCmd.Flags().BoolVarP(&withGitFileNames, "file-names", "n", false, "show changed files in git log")
}

func logCommand(cmd *cobra.Command, args []string) {
	if dockerLogName != "" {
		dockerLog()
	} else {
		gitLog() // default to git log
	}
}

func gitLog() {
	cmd := []string{"git", "log", "--format=" + gitLogFormat + "", "--graph", "-" + logLimit}
	if withGitFileNames {
		cmd = append(cmd, "--name-only")
	}
	mp.Exec(cmd)
}

func dockerLog() {
	searchString := ""
	for _, c := range dockerLogName {
		searchString += `\S*` + string(c)
	}
	searchString += `\S*`

	containers, err := docker.Ps(&docker.PsOptions{})
	mp.CheckErrorExit(err)
	for _, container := range containers {
		for _, name := range container.Names {
			match, err := regexp.Match(searchString, []byte(name))
			mp.CheckErrorExit(err)
			if match {
				command := []string{"docker", "logs", container.ID, "--tail", logLimit}
				if dockerLogFollow {
					command = append(command, "-f")
				}
				mp.Exec(command)
			}
		}
	}
	fmt.Println("no containers are running")
}
