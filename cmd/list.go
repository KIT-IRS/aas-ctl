package cmd

import (
	"aas-ctl/config"
	"aas-ctl/utils"
	"fmt"
	"log"

	aastypes "github.com/aas-core-works/aas-core3.0-golang/types"
	"github.com/spf13/cobra"
)

var (
	// listAASCmd represents the "list" and the "aas list" command
	listAASCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all Shells in the repository",
		Long:    "List all Shells in the repository",
		Run:     listAAS,
	}
	// listSMCmd represents the "sm list" command
	listSMCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List Submodels in the repository",
		Long: `list [--aas <Identifier>]
List Submodels in the repository`,
		Run: listSM,
	}
	// listConfigCmd represents the "config list" command
	listConfigCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Lists all selectable profiles defined in $HOME/.aas/config.json",
		Long:    fmt.Sprintf("Lists all selectable profiles defined in %v", config.ConfigFile),
		Run:     listConfig,
	}
)

func init() {
	rootCmd.AddCommand(listAASCmd)
	aasCmd.AddCommand(listAASCmd)
	smCmd.AddCommand(listSMCmd)
	configCmd.AddCommand(listConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// listAAS function executes the aas list command
func listAAS(cmd *cobra.Command, args []string) {
	flags, err := utils.NewFlagsFromCMD(cmd)
	if err != nil {
		log.Fatal(err)
	}
	shells, err := utils.GetAllShells()
	if err != nil {
		log.Fatal(err)
	}
	for _, shell := range shells {
		utils.PrintIdentifiable(shell, flags, false)
		fmt.Println()
	}
}

// lsitSM function executes the sm list commmand
func listSM(cmd *cobra.Command, args []string) {
	flags, err := utils.NewFlagsSMFromCMD(cmd)
	if err != nil {
		log.Fatal(err)
	}
	var submodels []aastypes.ISubmodel
	if flags.Shell != "" {
		submodels, err = utils.GetAllSubmodelsOfShell(flags.Shell)
	} else {
		submodels, err = utils.GetAllSubmodels()
	}
	if err != nil {
		log.Fatal(err)
	}
	for _, submodel := range submodels {
		utils.PrintIdentifiable(submodel, flags.Flags, false)
		fmt.Println()
	}
}

// listConfig function executes the config list command
func listConfig(cmd *cobra.Command, args []string) {
	configs, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	for _, profile := range configs.Profiles {
		if profile == *configs.ActiveProfile {
			profile.PrintActive()
		} else {
			profile.Print()
		}
	}
}
