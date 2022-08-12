/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName    = "annocli"
	configName = "config"
	configType = "yaml"

	viperUsername  = "username"
	viperPassword  = "password"
	viperServerURL = "serverURL"
	viperToken     = "token"
)

func configViper() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	configDir := filepath.Join(homeDir, appName)
	configFile := filepath.Join(configDir, configName+"."+configType)
	log.Println("configDir:", configDir)
	log.Println("configFile:", configFile)

	viper.SetConfigName("config")  // name of config file (without extension)
	viper.SetConfigType("yaml")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configDir) // call multiple times to add many search paths

	// check if error is "file not exists"
	_, error := os.Stat(configFile)
	if os.IsNotExist(error) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			log.Fatal(err)
		}
		if _, err := os.Create(configFile); err != nil {
			log.Fatal(err)
		}
	}

	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %v \n", err)
	}

	log.Println(viperUsername, viper.GetString(viperUsername))
	log.Println(viperPassword, viper.GetString(viperPassword))
	log.Println(viperServerURL, viper.GetString(viperServerURL))
	log.Println(viperToken, viper.GetString(viperToken))

	// watch file changed and reload
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "annocli",
	Short: "Annotation AI CLI tool which interact with it's API servers",
	Long:  `Annotation AI CLI tool which interact with it's API servers`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.annocli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	configViper()
}
