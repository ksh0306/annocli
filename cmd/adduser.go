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

type User struct {
	// UserId   int    `json:userid` // auto generated from the server
	Username string `json:username`
	Password string `json:password`
}

var signUpURL = "http://localhost:1323/api/signup"

func addUser(cmd *cobra.Command, args []string) {
	fmt.Println("adduser called")

	userInfo := strings.Split(account, ":")
	username := userInfo[0]
	password := userInfo[1]

	postBody, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(signUpURL, "application/json", responseBody)

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

// adduserCmd represents the adduser command
var adduserCmd = &cobra.Command{
	Use:   "adduser",
	Short: "Sign up new user",
	Long: `Sign up new user
Usage:

wip
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