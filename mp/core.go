/*
Package mp : core helper functions

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
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// IsProjectType determines the type of project mug is executed in
func IsProjectType(projectType string) bool {
	var filename string

	switch projectType {
	case "angular":
		filename = "angular.json"
	case "npm":
		filename = "package.json"
	case "go":
		filename = "main.go"
	case "gradle":
		filename = "gradlew"
	default:
		filename = ""
	}
	if filename == "" {
		fmt.Println("unknown project type")
		os.Exit(1)
	}

	filename = filepath.FromSlash(filepath.ToSlash(WorkingDirectory()) + "/" + filename)
	return FileExists(filename)
}

// FilePath returns a mugFile os compliant mugFile path in the current working directory
func FilePath(paths ...string) string {
	var path = filepath.ToSlash(WorkingDirectory())
	for _, p := range paths {
		path += "/" + p
	}
	return filepath.FromSlash(path)
}

// FileExists checks wether the mugFile with the given name exists
func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// ReadFile returns the contents of a file as string
func ReadFile(name string) string {
	data, err := ioutil.ReadFile(name)
	CheckErrorExit(err)
	return string(data)
}

// WorkingDirectory determines the directory mug is executed in
func WorkingDirectory() string {
	wd, err := os.Getwd()
	CheckErrorExit(err)
	return wd
}

// Exec executed an os command without environment params
func Exec(command []string) {
	ExecEnv(command, []string{})
}

// ExecEnvResult executes an os command with environment params
func ExecEnvResult(command []string, env []string) []byte {
	lookupCommand, err := exec.LookPath(command[0])
	CheckErrorExit(err)
	args := append([]string{lookupCommand}, command[1:]...)
	cmd := &exec.Cmd{
		Path: lookupCommand,
		Args: args,
		Env:  append(os.Environ(), env...),
	}
	// fmt.Println(cmd.String())
	output, err := cmd.CombinedOutput()
	CheckErrorExit(err)
	return output
}

// ExecEnv executes an os command with environment params
func ExecEnv(command []string, env []string) {
	lookupCommand, err := exec.LookPath(command[0])
	CheckErrorExit(err)
	args := append([]string{lookupCommand}, command[1:]...)
	cmd := &exec.Cmd{
		Path:   lookupCommand,
		Args:   args,
		Env:    append(os.Environ(), env...),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	// fmt.Println(cmd.String())
	err = cmd.Run()
	CheckErrorExit(err)
}

// CurrentGitBranch returns the name of the current git branch in the working directory
func CurrentGitBranch() string {
	output := ExecEnvResult([]string{"git", "branch", "--show-current"}, []string{})
	return fmt.Sprintf("%s", output)
}

// IsWindows determines whether mug runs on windows or not
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// CreateConfigFile creates the mug config mugFile in the given directory
func CreateConfigFile(dir string) {
	filename := filepath.FromSlash(filepath.ToSlash(dir) + "/.mug.yaml")
	if !FileExists(filename) {
		err := ioutil.WriteFile(filename, []byte(""), 0644)
		CheckErrorExit(err)
		fmt.Println("init done. created config mugFile")
	} else {
		fmt.Println("init done. mugFile exists")
	}
}
