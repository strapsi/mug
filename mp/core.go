package mp

import (
	// "bytes"
	"fmt"
	"os"
	// "strings"
	"os/exec"
	"path/filepath"
)

func IsProjectType(projectType string) bool {
	var filename string

	switch  projectType {
		case "angular": filename = "angular.json"
		case "npm": filename = "package.json"
		case "go": filename = "main.go"
		default: filename = ""
	}
	if filename == "" {
		fmt.Println("no projecttype given")
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
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return wd
}

func Exec(command []string) {
	lookupCommand, err := exec.LookPath(command[0])
	if err != nil {
		fmt.Println(err)
	}
	args := []string{lookupCommand}
	for _, arg := range command[1:] {
		args = append(args, arg)
	}
	cmd := &exec.Cmd {
		Path: lookupCommand,
		Args: args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,		
	}
	fmt.Println(cmd.String())
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
