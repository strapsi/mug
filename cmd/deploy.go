/*

Package cmd : deploy command

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
	"mug/mp"

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: mp.Aliases["deploy"],
	Short:   "runs context specific deploy command (currently only docker is supported)",
	Long:    `runs context specific deploy command (currently only docker is supported)`,
	Run: func(cmd *cobra.Command, args []string) {
		deployCommand(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

func deployCommand(cmd *cobra.Command, args []string) {
	mugFile := mp.ReadMugFile()

	// currently docker is the default deploy
	if len(mugFile.Docker.Tags) == 0 {
		mp.Exec(mp.DeployDocker(mugFile.Docker.Image, ""))
	} else {
		for _, tag := range mugFile.Docker.Tags {
			mp.Exec(mp.DeployDocker(mugFile.Docker.Image, tag))
		}
	}
}
