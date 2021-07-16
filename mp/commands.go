/*

Package mp : command helper functions

Copyright © 2021 m.vondergruen@protonmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this mugFile except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package mp

import (
	"fmt"
	"strings"
)

// BuildAngular returns ng build exec command
func BuildAngular(args []string, profile string) []string {
	_profile := "--"
	if profile == "" {
		_profile += "prod"
	} else {
		_profile += profile
	}
	return append([]string{"ng", "build", _profile}, args...)
}

// BuildNpm returns npm build exec command
func BuildNpm(args []string) []string {
	return append([]string{"npm", "run", "build"}, args...)
}

// BuildGradle returns gradle build exec command
func BuildGradle(args []string, useNativeGradle bool, ignoreTest bool, isWindows bool) []string {
	cmd := gradle(!useNativeGradle, isWindows)
	cmd = append(cmd, "clean", "build")
	if ignoreTest {
		cmd = append(cmd, "-x", "test")
	}
	cmd = append(cmd, args...)
	return cmd
}

// BuildGo returns go build exec command and environment variables
func BuildGo(args []string, mugfileArgs []string, target string) ([]string, []string) {
	cmd := []string{"go", "build"}
	cmd = append(cmd, mugfileArgs...)
	cmd = append(cmd, args...)

	var env []string
	switch target {
	case "linux":
		env = append(env, "GOOS=linux", "GOARCH=amd64")
		break
	case "windows":
		env = append(env, "GOOS=windows", "GOARCH=amd64")
		break
	}

	return cmd, env
}

// RunAngular return the ng run exec command
func RunAngular(args []string) []string {
	cmd := []string{"ng", "serve"}
	cmd = append(cmd, args...)
	return cmd
}

// RunNpm returns the npm start exec command
func RunNpm(args []string) []string {
	cmd := []string{"npm", "start"}
	if len(args) > 0 {
		cmd = append(cmd, "--")
		cmd = append(cmd, args...)
	}
	return cmd
}

// RunGradle returns the gradle bootRun exec command
func RunGradle(args []string, useNativeGradle bool, profile string, isWindows bool) []string {
	cmd := gradle(!useNativeGradle, isWindows)
	cmd = append(cmd, "bootRun")
	if profile != "" {
		cmd = append(cmd, "-Pprofile="+profile)
	}
	cmd = append(cmd, args...)
	return cmd
}

// RunGo returns the go run exec command
func RunGo(args []string) []string {
	cmd := []string{"go", "run", "main.go"}
	cmd = append(cmd, args...)
	return cmd
}

// FrdCommit returns the frd beng commit exec commands add and commit
func FrdCommit(args []string, branch string, addAll bool, overrideType string) ([]string, []string) {
	message := composeMessage(branch, overrideType, args)
	var addCmd []string
	if addAll {
		addCmd = []string{"git", "add", "."}
	}
	commitCmd := []string{"git", "commit", "-m", message}
	return addCmd, commitCmd
}

// LogDocker return the docker log exec command
func LogDocker(args []string, containerID string, limit string, follow bool) []string {
	cmd := []string{"docker", "logs", containerID, "--tail", limit}
	if follow {
		cmd = append(cmd, "-f")
	}
	cmd = append(cmd, args...)
	return cmd
}

// LogGit returns the git log exec command
func LogGit(args []string, format string, limit string, fileNames bool) []string {
	cmd := []string{"git", "log", "--format=" + format, "--graph", "-" + limit}
	if fileNames {
		cmd = append(cmd, "--name-only")
	}
	cmd = append(cmd, args...)
	return cmd
}

func composeMessage(branch string, overrideType string, args []string) string {
	if len(args) < 1 {
		ExitWithError("wrong number of arguments. expecting message")
	}
	commitType, id := parseBranch(branch, overrideType)
	var message string
	if strings.ToLower(id) == "x" {
		message = fmt.Sprintf("[%s] %s", commitType, id)
	} else {
		message = fmt.Sprintf("[%s][FRD-%s] %s", commitType, id, args[0])
	}
	return message
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
	case "test", "t":
		return "TEST"
	default:
		ExitWithError("unknown commit type " + commitType)
	}
	return ""
}

// parseBranch parses name of branch and returns the commit type and the jira id
func parseBranch(branch string, overrideType string) (string, string) {
	parts := strings.Split(branch, "-")
	if len(parts) < 3 {
		ExitWithError("branch has wrong format <type>-<id>-cool-feature")
	}
	var returnType string
	if overrideType != "" {
		returnType = overrideType
	} else {
		returnType = parts[0]
	}
	return parseCommitType(returnType), parts[1]
}

// gradle os specific gradle command
func gradle(useScript bool, isWindows bool) []string {
	var cmd []string
	if useScript {
		if isWindows {
			cmd = append(cmd, "cmd.exe", "/C", "gradlew.bat")
		} else {
			cmd = append(cmd, "sh", "gradlew")
		}
	} else {
		cmd = append(cmd, "gradle")
	}
	return cmd
}
