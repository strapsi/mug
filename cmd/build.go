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
	"mug/mp"

	"github.com/spf13/cobra"
)

var preferNpmBuild bool
var ngProfile string
var ignoreTestOnBuild bool
var goBuildTarget string
var useNativeGradleForBuild bool
var dockerBuild bool

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
	buildCmd.Flags().BoolVarP(&dockerBuild, "docker", "d", false, "additionally build docker container. you must provide infos in .mugfile")
}

func buildCommand(cmd *cobra.Command, args []string) {
	mugFile := mp.ReadMugFile()
	if mp.IsProjectType("angular") && !preferNpmBuild {
		fmt.Println("running ng build " + ngProfile)
		mp.Exec(mp.BuildAngular(args, ngProfile))
	} else if mp.IsProjectType("npm") {
		fmt.Println("running npm run build")
		mp.Exec(mp.BuildNpm(args))
	} else if mp.IsProjectType("gradle") {
		fmt.Println("running gradle clean build")
		mp.Exec(mp.BuildGradle(args, useNativeGradleForBuild, ignoreTestOnBuild, mp.IsWindows()))
	} else if mp.IsProjectType("go") {
		fmt.Println("running go build")
		mp.ExecEnv(mp.BuildGo(args, mugFile.Build.Args, goBuildTarget))
	}
	if dockerBuild {
		fmt.Print("running docker build")
		mp.Exec(mp.BuildDocker(mugFile.Docker.Image, mugFile.Docker.Tags, mp.WorkingDirectory(), args))
	}
}
