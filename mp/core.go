package mp

import (
	"fmt"
	"os"
	"runtime"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

func IsProjectType(projectType string) bool {
	var filename string

	switch  projectType {
		case "angular": filename = "angular.json"
		case "npm": filename = "package.json"
		case "go": filename = "main.go"
		case "gradle": filename = "gradlew"
		default: filename = ""
	}
	if filename == "" {
		fmt.Println("unknown project type")
		os.Exit(1)
	}
	
	filename = filepath.FromSlash(filepath.ToSlash(WorkingDirectory()) + "/" + filename)
	return FileExists(filename) 
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func WorkingDirectory() string {
	wd, err := os.Getwd()
	CheckErrorExit(err)
	return wd
}

func Exec(command []string) {
	ExecEnv(command, []string{})
}

func ExecEnv(command []string, env []string) {
	lookupCommand, err := exec.LookPath(command[0])
	CheckErrorExit(err)
	args := append([]string{lookupCommand}, command[1:]...)	
	cmd := &exec.Cmd {
		Path: lookupCommand,
		Args: args,
		Env: append(os.Environ(), env...),
		Stdout: os.Stdout,
		Stderr: os.Stderr,		
	}
	// fmt.Println(cmd.String())
	err = cmd.Run()
	CheckErrorExit(err)
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func CreateConfigFile(dir string) {
	filename := filepath.FromSlash(filepath.ToSlash(dir) + "/.ninja.yaml")
    if !FileExists(filename) {
    	err := ioutil.WriteFile(filename, []byte(""), 0644)
    	CheckErrorExit(err)
    	fmt.Println("init done. created config file")
    } else {
    	fmt.Println("init done. file exists")
    }
}
