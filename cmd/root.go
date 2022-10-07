/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"log"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "weasel",
	Short: "Cli to discover wallet with activity and from randomly private key",
	Long: `Allow search ethereum accounts from a randomly private key.
Each account with some activity is stored to make subsequent actions.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("[FAIL]: %v", err)
		os.Exit(1)
	}
}

func SetVersion(version string) {
	rootCmd.Flags().BoolP("version", "v", false, "Number of version of weasel")
	rootCmd.Version = version
}

func init() {

	rootCmd.Version = "unknown"

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flag("version").Value.String() == "true" {
			cmd.Printf(cmd.Version)
			return nil
		}

		cmd.HelpFunc()(cmd, args)
		return nil
	}
}
