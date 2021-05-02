/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"os"
	"github.com/spf13/cobra"
	"ninja/mp"
)

var PreferNpm bool

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "runs default development environment for working directory",
	Long: `Checks which project we are in and runs the according developent script
e.g. gradle bootRun or ng serve`,
	Run: func(cmd *cobra.Command, args []string) {
		runCommand(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&PreferNpm, "npm", "n", false, "prefer npm over ng to run project")
	// runCmd.MarkFlagRequired("dbl")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runCommand(cmd *cobra.Command, args []string) {
	if (mp.IsProjectType("angular") && !PreferNpm) { 
		fmt.Println("running ng serve")		
		mp.Exec(append([]string{"ng", "serve"}, args...))
		os.Exit(0)
	}
	if (mp.IsProjectType("npm")) { 
		fmt.Println("running npm start")		
		mp.Exec([]string{"npm", "start"})
		os.Exit(0)
	}
	if (mp.IsProjectType("go")) {
		mp.Exec([]string{"go", "run", "main.go"})
		os.Exit(0)
	}
}








