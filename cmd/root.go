////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Package cmd initializes the CLI and config parsers as well as the logger.
package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/user-discovery-bot/udb"
	"io/ioutil"
	"os"
)

var cfgFile string
var verbose bool
var showVer bool
var ndfPath string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "user-discovery-bot",
	Short: "Runs a user discovery bot for cMix",
	Long:  `This bot provides user lookup and search functions on cMix`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if showVer {
			printVersion()
			return
		}

		sess := viper.GetString("sessionfile")

		ndfBytes, err := ioutil.ReadFile(ndfPath)
		if err != nil {
			globals.Log.FATAL.Panicf("Could not read network definition file: %v", err)
		}

		ndfJSON := api.VerifyNDF(string(ndfBytes), "")

		StartBot(sess, ndfJSON)
	},
}

// Execute adds all child commands to the root command and sets flags
// appropriately.  This is called by main.main(). It only needs to
// happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		udb.Log.ERROR.Println(err)
		os.Exit(1)
	}
}

// init is the initialization function for Cobra which defines commands
// and flags.
func init() {
	udb.Log.DEBUG.Print("Printing log from init")
	// NOTE: The point of init() is to be declarative.
	// There is one init in each sub command. Do not put variable declarations
	// here, and ensure all the Flags are of the *P variety, unless there's a
	// very good reason not to have them as local params to sub command."
	cobra.OnInitialize(initConfig, initLog)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.Flags().StringVarP(&cfgFile, "config", "", "",
		"config file (default is $PWD/udb.yaml)")
	RootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false,
		"Verbose mode for debugging")
	RootCmd.Flags().BoolVarP(&showVer, "version", "V", false,
		"Show the server version information.")
	RootCmd.PersistentFlags().StringVarP(&ndfPath,
		"ndf",
		"n",
		"ndf.json",
		"Path to the network definition JSON file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile == "" {
		// Default search paths
		var searchDirs []string
		searchDirs = append(searchDirs, "./") // $PWD
		// $HOME
		home, _ := homedir.Dir()
		searchDirs = append(searchDirs, home+"/.elixxir/")
		// /etc/elixxir
		searchDirs = append(searchDirs, "/etc/elixxir")
		jww.DEBUG.Printf("Configuration search directories: %v", searchDirs)

		for i := range searchDirs {
			cfgFile = searchDirs[i] + "gateway.yaml"
			_, err := os.Stat(cfgFile)
			if !os.IsNotExist(err) {
				break
			}
		}
	}
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Unable to read config file (%s): %+v", cfgFile, err.Error())
	}

}

// initLog initializes logging thresholds and the log path.
func initLog() {
	if viper.Get("logPath") != nil {
		// If verbose flag set then log more info for debugging
		if verbose || viper.GetBool("verbose") {
			fmt.Printf("Logging verbosely\n")
			udb.Log.SetLogThreshold(jww.LevelDebug)
			udb.Log.SetStdoutThreshold(jww.LevelDebug)
		} else {
			udb.Log.SetLogThreshold(jww.LevelInfo)
			udb.Log.SetStdoutThreshold(jww.LevelInfo)
		}
		// Create log file, overwrites if existing
		logPath := viper.GetString("logPath")
		logFile, err := os.Create(logPath)
		if err != nil {
			udb.Log.WARN.Println("Invalid or missing log path, default path used.")
		} else {
			udb.Log.SetLogOutput(logFile)
		}
	}
}
