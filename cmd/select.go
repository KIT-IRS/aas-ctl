package cmd

import (
	"aas-ctl/config"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select a profile as new active profile",
	Long: fmt.Sprintf(`select <name>
Select the profile with the given <name> as new active profile
Selectable profiles can be listed with "config list"
Profiles get managed in "%v"`, config.ConfigFile),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Printf("Config select needs 1 parameter, %v given", len(args))
			return
		}
		configs, err := config.LoadConfig()
		if err != nil {
			log.Fatal(err)
		}
		if configs.ActiveProfile.GetName() == args[0] {
			return
		}
		for _, profile := range configs.Profiles {
			if profile.GetName() == args[0] {
				configs.ActiveProfile = &profile
				configs.Save()
				return
			}
		}
		fmt.Printf("There is no config with the name %v", args[0])
	},
}

func init() {
	configCmd.AddCommand(selectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// selectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// selectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
