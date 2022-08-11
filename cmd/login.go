/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var signInURL = "http://localhost:1323/api/signin"

func loginUser(cmd *cobra.Command, args []string) {
	fmt.Println("login called")

	userInfo := strings.Split(account, ":")
	username := userInfo[0]
	password := userInfo[1]

	postBody, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(signInURL, "application/json", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	fmt.Println("body:", string(body))
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in user",
	Long: `Log in user
Usage:

wip
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
