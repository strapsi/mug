/*
Package mp : mug-file

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
	"time"

	"gopkg.in/yaml.v2"
)

// MugFile represents the known mugfile properties
type MugFile struct {
	_Version string        `yaml:"__version"`
	Build    MugFileBuild  `yaml:"build,omitempty"`
	Docker   MugFileDocker `yaml:"docker,omitempty"`
}

// MugFileBuild mug build properties
type MugFileBuild struct {
	Args []string `yaml:"args,omitempty"`
}

// MugFileDocker mug properties for docker commands
type MugFileDocker struct {
	Image string   `yaml:"image"`
	Tags  []string `yaml:"tags,omitempty"`
}

var mugFile = MugFile{}
var userVariables map[interface{}]interface{}

// ReadMugFile returns the mug mugFile as struct
func ReadMugFile() *MugFile {
	start := time.Now()
	if mugFile._Version != "" {
		return &mugFile
	}
	if FileExists(FilePath(".mugfile")) {
		fmt.Print("reading .mugfile")
		content := []byte(ReadFile(FilePath(".mugfile")))
		readKnownProperties(content)
		readUserVariables(content)
		replaceVariablesInArgs()
		fmt.Printf(": took %s\n", time.Since(start))
	}
	return &mugFile
}

func replaceVariablesInArgs() {
	vars := map[string]string{}
	for varKey, varValue := range userVariables {
		if keyAsString, ok := varKey.(string); ok {
			if keyAsString[0] == '$' {
				vars[keyAsString] = varValue.(string)
			}
		}
	}

	for i, arg := range mugFile.Build.Args {
		// substitute user params
		var replaced string
		for key, value := range vars {
			replaced = strings.ReplaceAll(arg, key, value)
		}

		// substitute program params ($utc)
		if strings.Contains(arg, "$utc") {
			utc := strings.ReplaceAll(time.Now().UTC().String(), " ", "_")
			replaced = strings.ReplaceAll(replaced, "$utc", utc)
		}
		mugFile.Build.Args[i] = replaced
	}
}

func readUserVariables(content []byte) {
	err := yaml.Unmarshal(content, &userVariables)
	CheckErrorExit(err)
}

func readKnownProperties(content []byte) {
	err := yaml.Unmarshal(content, &mugFile)
	CheckErrorExit(err)
}
