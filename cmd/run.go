/*

Package cmd : run command

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

var preferNpmStart bool
var springProfile string
var useNativeGradleForRun bool

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: mp.Aliases["run"],
	Short:   "runs default development environment for working directory",
	Long: `Checks which project we are in and runs the according run command
e.g. gradle bootRun or ng serve`,
	Run: func(cmd *cobra.Command, args []string) {
		runCommand(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&preferNpmStart, "npm", "n", false, "prefer npm over ng to run project")
	runCmd.Flags().StringVarP(&springProfile, "profile", "p", "", "spring profile")
	runCmd.Flags().BoolVarP(&useNativeGradleForRun, "gradle", "g", false, "use native gradle")
}

func runCommand(cmd *cobra.Command, args []string) {
	if mp.IsProjectType("angular") && !preferNpmStart {
		fmt.Println("running ng serve")
		mp.Exec(mp.RunAngular(args))
	}
	if mp.IsProjectType("npm") {
		fmt.Println("running npm start")
		mp.Exec(mp.RunNpm(args))
	}
	if mp.IsProjectType("gradle") {
		fmt.Println("running gradlew bootRun")
		mp.Exec(mp.RunGradle(args, useNativeGradleForRun, springProfile, mp.IsWindows()))
	}
	if mp.IsProjectType("go") {
		fmt.Println("running go run")
		mp.Exec(mp.RunGo(args))
	}
}
