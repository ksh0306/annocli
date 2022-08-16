/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset all config",
	Long: `Reset config.yaml file
	
	example:
	$ annocli config reset
	`,
	Run: func(cmd *cobra.Command, args []string) {
		homedir := os.Getenv("HOME")
		if err := os.Remove(homedir + "/" + appName + "/" + configName + "." + configType); err != nil {
			fmt.Println("reset config failed: ", err)
			return
		}
		fmt.Println("config.yaml removed. reset config completed.")

	},
}

func init() {
	configCmd.AddCommand(resetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
