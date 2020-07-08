package main

import (
	"log"

	"github.com/spf13/viper"
)

type configuration struct {
	allowedProjects []string
}

var config configuration

func initConfig() {
	viper.SetConfigName("autosigner")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/puppetlabs/puppet/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Config file not found. Please create the autosigner.yaml config.")
		} else {
			log.Fatalf("Error in reading config file: %v", err)
		}
	}

	config.allowedProjects = viper.GetStringSlice("allowed_projects")
}
