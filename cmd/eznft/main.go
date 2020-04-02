package main

import (
	"eznft/commands"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"

	"os"
)

var cfgFile string

func main() {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:     "eznft",
		Short:   "Run or orchestrate NFTs",
		Long:    `Run or orchestrate non-functional tests`,
		Version: "alpha",
	}

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is eznft.yaml)")
	commands.AddCommands(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("eznft")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
