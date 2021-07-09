/*

Package cmd : build command

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
	"time"
)

var preferNpmBuild bool
var ngProfile string
var ignoreTestOnBuild bool
var goBuildTarget string
var useNativeGradleForBuild bool
var startTime time.Time

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b"},
	Short:   "build project",
	Long:    `detects the type of project we are in and builds it`,
	Run: func(cmd *cobra.Command, args []string) {
		// activate verbose mode for build (just for printing build time)
		err := rootCmd.Flags().Set("verbose", "true")
		mp.CheckErrorExit(err)
		buildCommand(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().BoolVarP(&preferNpmBuild, "npm", "n", false, "prefer \"npm run build\" over \"ng build\" to build project")
	buildCmd.Flags().StringVarP(&ngProfile, "profile", "p", "", "specify angular profile to build with e.g. prod")
	buildCmd.Flags().BoolVarP(&ignoreTestOnBuild, "ignore-tests", "i", false, "ignore tests on build")
	buildCmd.Flags().StringVarP(&goBuildTarget, "target", "t", "", "go build target [linux | windows]")
	buildCmd.Flags().BoolVarP(&useNativeGradleForBuild, "gradle", "g", false, "use native gradle binary instead of gradlew")
}

func buildCommand(cmd *cobra.Command, args []string) {
	if mp.IsProjectType("angular") && !preferNpmBuild {
		profile := "prod"
		if ngProfile != "" {
			profile = ngProfile
		}
		profile = "--" + profile

		fmt.Println("running ng build " + profile)
		mp.Exec(append([]string{"ng", "build", profile}, args...))
	} else if mp.IsProjectType("npm") {
		fmt.Println("running npm run build")
		mp.Exec([]string{"npm", "run", "build"})
	} else if mp.IsProjectType("gradle") {
		fmt.Println("running gradle clean build")
		cleanBuild := append(mp.Gradle(!useNativeGradleForBuild), "clean", "build")
		if ignoreTestOnBuild {
			cleanBuild = append(cleanBuild, "-x", "test")
		}
		mp.Exec(cleanBuild)
	} else if mp.IsProjectType("go") {
		mugFile := mp.ReadMugFile()
		fmt.Println("running go build")
		var env []string
		if goBuildTarget != "" {
			if goBuildTarget == "linux" {
				env = append(env, "GOOS=linux", "GOARCH=amd64")
			}
			if goBuildTarget == "windows" {
				env = append(env, "GOOS=windows", "GOARCH=amd64")
			}
		}
		env = append(env, "BUILD_TIME="+time.Now().String())
		command := append([]string{"go", "build"}, mugFile.Build.Args...)
		command = append(command, args...)
		mp.ExecEnv(command, env)
	}
}
