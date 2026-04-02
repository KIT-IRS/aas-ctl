/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"aas-ctl/config"
	"log"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new config profile",
	Long:  `Create a new config profile, providing the name`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("Expected 1 argument, got %d", len(args))
		}
		profile := config.CreateProfileWithName(args[0])
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatal(err)
		}
		err = cfg.AddProfile(profile)
		if err != nil {
			log.Fatal(err)
		}
		err = cfg.Save()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	configCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
