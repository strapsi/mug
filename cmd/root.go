/*

Package cmd : root command

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
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"mug/build"
	"time"
)

var cfgFile string
var beVerbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mug",
	Short: "mug - a mug can do anything",
	Long: `
   _____   ____ ___  ________ 
  /     \ |    |   \/  _____/ 
 /  \ /  \|    |   /   \  ___ 
/    Y    \    |  /\    \_\  \
\____|__  /______/  \______  /
        \/                 \/ 

version: ` + build.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	start := time.Now()
	cobra.CheckErr(rootCmd.Execute())
	if beVerbose {
		fmt.Printf("took %s\n", time.Since(start))
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().BoolVarP(&beVerbose, "verbose", "v", false, "print out more information")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mug" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mug")
		viper.SetConfigType("yml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	//fmt.Println("hallo000")
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Printf("config ist: %v\n", viper.GetStringMapString("config")["fish-config"])
		//fmt.Printf("config: %v\n", viper.Get("config"))
		//fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		//fmt.Fprintf(os.Stderr, "config: %v\n", viper.Get("config"))
	}
}
