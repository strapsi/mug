/*

Package cmd : proxy command

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
	"crypto/tls"
	"encoding/json"
	"mug/mp"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/strapsi/go-docker"
)

var kongURL string
var proxyBackendModule string
var proxyBackendURL string
var angularURL string
var credentials string

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:     "proxy",
	Aliases: mp.Aliases["proxy"],
	Short:   "start the kong proxy container",
	Long:    `a long version of the short version :D`,
	Run: func(cmd *cobra.Command, args []string) {
		proxyCommand()
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
	proxyCmd.Flags().StringVarP(&kongURL, "kong-url", "k", "10.10.227.175", "kong url")
	proxyCmd.Flags().StringVarP(&proxyBackendModule, "backend", "b", "", "name of the backend module")
	proxyCmd.Flags().StringVarP(&proxyBackendURL, "backend-url", "u", "", "url of the backend")
	proxyCmd.Flags().StringVarP(&angularURL, "angular-url", "a", "", "url of the frontend")
	proxyCmd.Flags().StringVarP(&credentials, "credentials", "c", "user1:be32", "<user>:<password>")
}

func proxyCommand() {
	token := login()
	dockerEnv := buildDockerEnv(token)
	options := &docker.RunOptions{
		Image: "dontenwill/dev-kong-proxy:ci",
		Name:  "kong-proxy",
		Force: true,
		Env:   dockerEnv,
	}
	err := docker.Run(options)
	mp.CheckErrorExit(err)
}

func login() string {
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transport}
	url := buildURL()
	response, err := client.Get(url)
	mp.CheckErrorExit(err)
	defer response.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	mp.CheckErrorExit(err)
	return data["token"].(string)
}

func buildDockerEnv(token string) map[string]string {
	dockerEnv := map[string]string{
		"BE_TOKEN": token,
	}
	if kongURL != "" {
		dockerEnv["KONG_URL"] = kongURL
	}
	if proxyBackendModule != "" && proxyBackendURL != "" {
		dockerEnv["BACKEND_MODULE"] = proxyBackendModule
		dockerEnv["BACKEND_URL"] = proxyBackendURL
	}
	if angularURL != "" {
		dockerEnv["NG_URL"] = angularURL
	}
	return dockerEnv
}

func buildURL() string {
	parts := strings.Split(credentials, ":")
	if len(parts) != 2 {
		mp.ExitWithError("wrong credentials format. has to be <user>:<password>")
	}
	return "https://" + kongURL + "/api/auth/auth?login=" + parts[0] + "&password=" + parts[1]
}
