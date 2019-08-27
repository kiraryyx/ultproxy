package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var envSetsToBind [][]string

func main() {
	cobra.OnInitialize(func() {
		// TODO: Support upper-case env
		if shouldUseEnv {
			for _, set := range envSetsToBind {
				if err := configViper.BindEnv(set...); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}

		if configFile == "" {
			configViper.SetConfigName("ultproxy")
			configViper.AddConfigPath(".")
			configViper.AddConfigPath("$HOME/.ultproxy")
			configViper.AddConfigPath("/etc/ultproxy")
		} else {
			configViper.SetConfigFile(configFile)
		}

		if err := configViper.ReadInConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Using config file:", configViper.ConfigFileUsed())

		if err := configViper.Unmarshal(&config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	})

	// Hide all flags which "cybozu-go/well" sets
	pflag.VisitAll(func(f *pflag.Flag) {
		f.Hidden = true
		f.Changed = true
	})

	execute()
}
