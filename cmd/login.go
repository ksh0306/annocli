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
	fmt.Println("login called")
	loginURL := viper.GetString(viperServerURL) + loginEndpoint

	userInfo := strings.Split(account, ":")
	username := userInfo[0]
	password := userInfo[1]

	postBody, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(loginURL, "application/json", responseBody)

	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	fmt.Println("body:", string(body))

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal().Err(err).Send()
	}

	msg := data.(map[string]interface{})
	token := msg["token"]
	fmt.Println("token: ", token)

	// 여기서 username, password, token 을 저장하자
	viper.Set(viperUsername, username)
	viper.Set(viperPassword, password)
	viper.Set(viperToken, token)
	if err := viper.WriteConfig(); err != nil {
		log.Fatal().Err(err).Send()
	}

	// test
	fmt.Println(viperUsername, viper.GetString(viperUsername))
	fmt.Println(viperPassword, viper.GetString(viperPassword))
	fmt.Println(viperServerURL, viper.GetString(viperServerURL))
	fmt.Println(viperToken, viper.GetString(viperToken))
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
