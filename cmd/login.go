/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginEndpoint = "/login"

func loginUser(cmd *cobra.Command, args []string) {
	loginURL := viper.GetString(viperServerURL) + loginEndpoint

	userInfo := strings.Split(account, ":")
	username := userInfo[0]
	password := userInfo[1]

	postBody, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})

	requestBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(http.MethodPost, loginURL, requestBody)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	fmt.Println("--response status", resp.Status)
	fmt.Println("--response body", string(body))

	if resp.StatusCode != http.StatusOK {
		log.Info().Msg("failed to login")
		return
	}

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal().Err(err).Send()
	}

	msg := data.(map[string]interface{})
	token := msg["token"]

	// 여기서 username, password, token 을 저장하자
	viper.Set(viperUsername, username)
	viper.Set(viperPassword, password)
	viper.Set(viperToken, token)
	if err := viper.WriteConfig(); err != nil {
		log.Fatal().Err(err).Send()
	}

	// display current config
	fmt.Println("current conffig")
	fmt.Printf("--%s: %s\n", viperUsername, viper.GetString(viperUsername))
	fmt.Printf("--%s: %s\n", viperPassword, viper.GetString(viperPassword))
	fmt.Printf("--%s: %s\n", viperServerURL, viper.GetString(viperServerURL))
	fmt.Printf("--%s: %s\n", viperToken, viper.GetString(viperToken))
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in user",
	Long: `Log in user with account which means username and password
Example:

$ annocli login --account={username}:{password}
`,
	Run: loginUser,
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loginCmd.Flags().StringVarP(&account, "account", "a", "", "user account to login. it should be '{username}:{password}' format")
	loginCmd.MarkFlagRequired("account")
}
