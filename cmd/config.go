/*
Package cmd : config command

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
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"mug/mp"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:               "config [name]",
	ValidArgsFunction: configuredConfigFiles,
	Args:              validConfigFiles,
	Short:             "edit config files",
	Long:              ``,
	Run:               executeConfigCmd,
	Example:           "config fish-config",
}

func executeConfigCmd(cmd *cobra.Command, args []string) {
	mp.Exec(mp.EditConfigFile(viper.GetStringMapString("configs")[args[0]]))
}

func configuredConfigFiles(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return availableConfigs(), cobra.ShellCompDirectiveDefault
}

func validConfigFiles(cmd *cobra.Command, args []string) error {
	for _, config := range availableConfigs() {
		if config == args[0] {
			return nil
		}
	}
	return fmt.Errorf("config with name <%s> not found", args[0])
}

func availableConfigs() []string {
	var configs []string
	for key := range viper.GetStringMapString("configs") {
		configs = append(configs, key)
	}
	return configs
}

func init() {
	rootCmd.AddCommand(configCmd)
}
