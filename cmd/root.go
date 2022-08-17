/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
		log.Fatal().Err(err)
	}
	configDir := filepath.Join(homeDir, appName)
	configFile := filepath.Join(configDir, configName+"."+configType)
	log.Info().Msgf("configFile: %s", configFile)

	viper.SetConfigName("config")  // name of config file (without extension)
	viper.SetConfigType("yaml")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configDir) // call multiple times to add many search paths

	// check if error is "file not exists"
	_, error := os.Stat(configFile)
	if os.IsNotExist(error) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			log.Fatal().Err(err)
		}
		if _, err := os.Create(configFile); err != nil {
			log.Fatal().Err(err)
		}
	}

	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		log.Fatal().Err(err).Msg("Fatal error config file")
	}

	// watch file changed and reload
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info().Msgf("Config file changed: %s", e.Name)
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

func zerologConfig() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***%s****", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	log.Logger = zerolog.New(output).With().Caller().Timestamp().Logger()
}

var debug bool

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.annocli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("debug", "d", false, "More logs for debugging")

	zerologConfig()
	configViper()
}
