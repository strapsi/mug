/*

Package cmd : init command

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
	"os"
)

var createConfigDir string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Hidden: true,
	Use:   "init",
	Short: "initialize mug",
	Long:  `creates the mug config file. default location is $HOME but can be specified`,
	Run: func(cmd *cobra.Command, args []string) {
		initCommand(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&createConfigDir, "config-dir", "", "", "location of config file")
}

func initCommand(cmd *cobra.Command, args []string) {
	fmt.Println("initializing mug")
	var dir string
	if createConfigDir == "" {
		if mp.IsWindows() {
			dir = os.Getenv("UserProfile")
		} else {
			dir = os.Getenv("HOME")
		}
	} else {
		dir = createConfigDir
	}
	mp.CreateConfigFile(dir)
}
