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

var preferNpmBuild bool
var ngProfile string
var ignoreTestOnBuild bool
var goBuildTarget string

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build the current project",
	Long: `builds the current project we are in. can also build docker containers
e.g. gradlew clean build or ng build --prod`,
	Run: func(cmd *cobra.Command, args []string) {
		buildCommand(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().BoolVarP(&preferNpmBuild, "npm", "n", false, "prefer npm run build over ng build to run project")
	buildCmd.Flags().StringVarP(&ngProfile, "profile", "p", "", "angular profile e.g. prod")
	buildCmd.Flags().BoolVarP(&ignoreTestOnBuild, "ignore-tests", "i", false, "ignore tests on build")
	buildCmd.Flags().StringVarP(&goBuildTarget, "target", "t", "", "go build target linux | windows")
}

func buildCommand(cmd *cobra.Command, args []string) {
	if (mp.IsProjectType("angular") && !preferNpmBuild) { 
		profile := "prod"
		if ngProfile != "" {
			profile = ngProfile
		}
		profile = "--" + profile
	
		fmt.Println("running ng build " + profile)
		mp.Exec(append([]string{"ng", "build", profile}, args...))
		os.Exit(0)
	}
	if (mp.IsProjectType("npm")) { 
		fmt.Println("running npm run build")		
		mp.Exec([]string{"npm", "run", "build"})
		os.Exit(0)
	}
	if (mp.IsProjectType("gradle")) {
		fmt.Println("running gradlew clean build")
		var bootRun []string
		if mp.IsWindows() {
			bootRun = []string{"cmd.exe", "/C", "gradlew.bat"}			
		} else {		
			bootRun = []string{"sh", "gradlew"}
		}
		bootRun = append(bootRun, "clean", "build")
		if ignoreTestOnBuild {
			bootRun = append(bootRun, "-x", "test")	
		}
		mp.Exec(bootRun)
		os.Exit(0)
	}
	if (mp.IsProjectType("go")) {
		var env []string
		if goBuildTarget != "" {
			if goBuildTarget == "linux" {
				env = []string{"GOOS=linux", "GOARCH=amd64"}
			}
			if goBuildTarget == "windows" {
				env = []string{"GOOS=windows", "GOARCH=amd64"}
			}
		}
			
		mp.ExecEnv([]string{"go", "build"}, env)
		os.Exit(0)
	}
}
