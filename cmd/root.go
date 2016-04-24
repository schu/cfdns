package cmd

import (
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var api *cloudflare.API
var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "cfdns",
	Short: "CloudFlare DNS util",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initAPI)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/cfdns/config.yaml)")
}

func initAPI() {
	api = cloudflare.New(viper.GetString("api_key"), viper.GetString("api_email"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		xdgConfig = fmt.Sprintf("%s/.config/cfdns", os.Getenv("HOME"))
	} else {
		xdgConfig = fmt.Sprintf("%s/cfdns", xdgConfig)
	}
	viper.AddConfigPath(xdgConfig)
	viper.SetEnvPrefix("cfdns")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
