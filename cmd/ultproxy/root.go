package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config       Config
	configFile   string
	shouldUseEnv bool
)

var configViper = viper.New()

var rootCmd = &cobra.Command{
	Use: "ultproxy",
	Run: func(c *cobra.Command, args []string) {
		c.Help()
		os.Exit(1)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to the configuration file")
	rootCmd.PersistentFlags().BoolVarP(&shouldUseEnv, "use-env", "", false, "Use environment variables (http_proxy and no_proxy) to override configuration")

	addAndBindPFlagsStringP(rootCmd, configViper, []stringFlagParamSet{
		stringFlagParamSet{"loglevel", "", "info", "Log level [debug,info,warn,error,fatal,panic]", "logging.level"},
	})
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
