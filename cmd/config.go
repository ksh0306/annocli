/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const httpPrefix = "http://"

var serverURL string

func config(cmd *cobra.Command, args []string) {
	fmt.Println("config called")
	if !strings.HasPrefix(serverURL, httpPrefix) {
		serverURL = httpPrefix + serverURL
	}
	if _, err := url.Parse(serverURL); err != nil {
		log.Fatal(err)
	}
	if serverURL != "" {
		viper.Set(viperServerURL, serverURL)
		viper.WriteConfig()
	}
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config sets server URL and more",
	Long: `config sets server URL and more

usage: 
$ annocli config --server=http://222.110.65.138

if you want to use HTTPS, then use like below
$ annocli config --server=https://222.110.65.138
`,
	Run: config,
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	configCmd.Flags().StringVarP(&serverURL, "server", "s", "http://localhost:1323", "Annotation-AI API server URL")
}
