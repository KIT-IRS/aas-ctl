package cmd

import (
	"aas-ctl/config"
	"aas-ctl/utils"
	"errors"
	"fmt"
	"log"

	aastypes "github.com/aas-core-works/aas-core3.0-golang/types"
	"github.com/spf13/cobra"
)

var (
	// showCmd represents the "show" command
	showCmd = &cobra.Command{
		Use:   "show",
		Short: "Shows the Shell or the Submodel with the given <Identifier>",
		Long: `show <Identifier>
Shows the Shell or the Submodel with the given <Identifier> and  its contents.
Contents are Submodels for Shells and SubmodelElements for Submodels.
The <Identifier> may be an <ID> or an <IDShort>.
If an <IDShort> is provided, the first match is displayed.
The script first searches the AAS repository and then the SM repository.
This command does only provide the flags for the "aas show" commmand.
To use flags specific to the "sm show" command, use "sm show <Identifier>".`,
		Run: show,
	}
	// showAASCmd represents the "aas show" command
	showAASCmd = &cobra.Command{
		Use:   "show",
		Short: "Shows the Shell with the given <Identifier> and its Submodels",
		Long: `aas show <Identifier>
Shows the Shell with the given <Identifier> and its Submodels
Identifier may be an <id> or <idShort>`,
		Run: showAAS,
	}
	// showSMCmd represents the "sm show" command
	showSMCmd = &cobra.Command{
		Use:   "show",
		Short: "Shows the Submodel with the given <Identifier>",
		Long: `sm show <Identifier> [--aas <Identifier>]
Shows the Submodel with the given <Identifier>
Identifier may be an <id> or <idShort>
		`,
		Run: showSM,
	}
	// showConfigCmd represents the "config show" command
	showConfigCmd = &cobra.Command{
		Use:   "show",
		Short: "Shows the path of the config file",
		Long:  "Shows the path of the config file",
		Run:   showConfig,
	}
)

func init() {
	rootCmd.AddCommand(showCmd)
	aasCmd.AddCommand(showAASCmd)
	smCmd.AddCommand(showSMCmd)
	configCmd.AddCommand(showConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	showCmd.Flags().Bool("id", false, "Shows only the <id> of the element")
	showCmd.Flags().Bool("url", false, "Shows only the URL of the element")
	showCmd.Flags().Bool("json", false, "Shows the elements in JSON format")
	showCmd.MarkFlagsMutuallyExclusive("id", "url", "json")

	showSMCmd.Flags().String("elementId", "", "Shows only the SubmodelElement with the provided <idShort>")
	showSMCmd.Flags().Int("elementIdx", -1, "Shows the SubmodelElement with the provided <Index>")
	showSMCmd.MarkFlagsMutuallyExclusive("elementId", "elementIdx")
	showSMCmd.Flags().Bool("value", false, "Show the value of the submodel element")
}

// show function executes the standalone "show" command.
// It first looks for an AAS to show, if there is none, it looks for a SM.
func show(cmd *cobra.Command, args []string) {
	if err := utils.RequireSingleArg(args); err != nil {
		log.Fatal(err)
	}
	flags, err := utils.NewFlagsFromCMD(cmd)
	if err != nil {
		log.Fatal(err)
	}
	if err := tryShowAAS(args, flags); err == nil {
		return
	} else if !errors.As(err, new(*utils.IdentifiableNotFoundError)) {
		log.Fatal(err)
	}
	flagsSM := utils.NewFlagsSMShow()
	flagsSM.Flags = flags
	if err := tryShowSM(args, flagsSM); err != nil {
		log.Fatal(err)
	}
}

// showAAS function executes the "aas show" command
func showAAS(cmd *cobra.Command, args []string) {
	if err := utils.RequireSingleArg(args); err != nil {
		log.Fatal(err)
	}
	flags, err := utils.NewFlagsFromCMD(cmd)
	if err != nil {
		log.Fatal(err)
	}
	if err := tryShowAAS(args, flags); err != nil {
		log.Fatal(err)
	}
}

// showSM function executes the "sm show" command
func showSM(cmd *cobra.Command, args []string) {
	if err := utils.RequireSingleArg(args); err != nil {
		log.Fatal(err)
	}
	flags, err := utils.NewFlagsSMShowFromCMD(cmd)
	if err != nil {
		log.Fatal(err)
	}
	if err := tryShowSM(args, flags); err != nil {
		log.Fatal(err)
	}
}

// showConfig function exectutes the "config show" command
func showConfig(cmd *cobra.Command, args []string) {
	fmt.Println(config.ConfigFile)
}

// tryShowAAS tries to find the AAS specified in the args.
// If it can't find an AAS to show, it returns an error.
func tryShowAAS(args []string, flags *utils.Flags) error {
	shell, err := utils.GetShell(args[0])
	if err != nil {
		return err
	}
	utils.PrintIdentifiable(shell, flags, true)
	return nil
}

// tryShowSM tries to find the SM specified in the args.
// If it can't find a SM to show, it returns an error.
func tryShowSM(args []string, flags *utils.FlagsSMShow) error {
	var submodel aastypes.ISubmodel
	var err error
	if flags.Shell != "" {
		submodel, err = utils.GetShellSubmodel(flags.Shell, args[0])
		if err != nil {
			return err
		}
	} else {
		submodel, err = utils.GetSubmodel(args[0])
		if err != nil {
			return err
		}
	}
	ref, err := resolveSubmodelRef(submodel, flags.ElementID, flags.ElementIdx)
	if err != nil {
		return err
	}
	switch r := ref.(type) {
	case aastypes.ISubmodel:
		utils.PrintIdentifiable(submodel, flags.Flags, true)
	case aastypes.ISubmodelElement:
		utils.PrintSubmodelElement(submodel, r, flags)
	default:
		return fmt.Errorf("unsupported referable type for URL printing: %T", r)
	}
	return nil
}

// resolveSubmodelRef function is a helper function that returns the referable which is to show according to the flags.
func resolveSubmodelRef(sm aastypes.ISubmodel, elementId string, elementIdx int) (aastypes.IReferable, error) {
	if elementId != "" {
		return utils.FindSubmodelElement(sm, elementId)
	}
	if elementIdx != -1 {
		return utils.GetSubmodelElement(sm, elementIdx)
	}
	return sm, nil
}
