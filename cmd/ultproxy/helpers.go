package main

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type boolFlagParamSet struct {
	Name      string
	Shorthand string
	Value     bool
	Usage     string
	ConfigKey string
}

type intSliceFlagParamSet struct {
	Name      string
	Shorthand string
	Value     []int
	Usage     string
	ConfigKey string
}

type stringFlagParamSet struct {
	Name      string
	Shorthand string
	Value     string
	Usage     string
	ConfigKey string
}

type stringSliceFlagParamSet struct {
	Name      string
	Shorthand string
	Value     []string
	Usage     string
	ConfigKey string
}

func addAndBindPFlagsBoolP(c *cobra.Command, v *viper.Viper, sets []boolFlagParamSet) {
	for _, s := range sets {
		c.PersistentFlags().BoolP(s.Name, s.Shorthand, s.Value, s.Usage)
		v.BindPFlag(s.ConfigKey, c.PersistentFlags().Lookup(s.Name))
	}
}

func addAndBindPFlagsIntSliceP(c *cobra.Command, v *viper.Viper, sets []intSliceFlagParamSet) {
	for _, s := range sets {
		// Convert []int to []string because viper.Unmarshal(&config) outputs parse error with IntSlice flag
		value := make([]string, len(s.Value))
		for i, n := range s.Value {
			value[i] = strconv.Itoa(n)
		}

		c.PersistentFlags().StringSliceP(s.Name, s.Shorthand, value, s.Usage)
		v.BindPFlag(s.ConfigKey, c.PersistentFlags().Lookup(s.Name))
	}
}

func addAndBindPFlagsStringP(c *cobra.Command, v *viper.Viper, sets []stringFlagParamSet) {
	for _, s := range sets {
		c.PersistentFlags().StringP(s.Name, s.Shorthand, s.Value, s.Usage)
		v.BindPFlag(s.ConfigKey, c.PersistentFlags().Lookup(s.Name))
	}
}

func addAndBindPFlagsStringSliceP(c *cobra.Command, v *viper.Viper, sets []stringSliceFlagParamSet) {
	for _, s := range sets {
		c.PersistentFlags().StringSliceP(s.Name, s.Shorthand, s.Value, s.Usage)
		v.BindPFlag(s.ConfigKey, c.PersistentFlags().Lookup(s.Name))
	}
}
