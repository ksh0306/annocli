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
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type User struct {
	// UserId   int    `json:userid` // auto generated from the server
	Username string `json:username`
	Password string `json:password`
}

var signUpEndpoint = "/v1/signup"

func addUser(cmd *cobra.Command, args []string) {

	sURL := viper.GetString(viperServerURL)
	if _, err := url.ParseRequestURI(sURL); err != nil {
		fmt.Println("server URL need to config")
		return
	}

	signUpURL := sURL + signUpEndpoint

	userInfo := strings.Split(account, ":")
	username := userInfo[0]
	password := userInfo[1]
	fmt.Println("add new user\n--username:", username)

	postBody, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	requestBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(http.MethodPost, signUpURL, requestBody)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", viper.GetString(viperToken)))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Debug().Str("response status", resp.Status).Send()
	log.Debug().Str("response body", string(body)).Send()
}

// adduserCmd represents the adduser command
var adduserCmd = &cobra.Command{
	Use:   "adduser",
	Short: "Sign up new user",
	Long: `Sign up new user
Usage:

$ annocli adduser --account=username:password
`,
	Run: addUser,
}

// flags
var (
	account string
)

func init() {
	rootCmd.AddCommand(adduserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adduserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adduserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	adduserCmd.Flags().StringVarP(&account, "account", "a", "", "user account to add. it should be '{username}:{password}' format")
	adduserCmd.MarkFlagRequired("account")
}
