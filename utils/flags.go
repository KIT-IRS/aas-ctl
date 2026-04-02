package utils

import (
	"github.com/spf13/cobra"
)

// Flags struct groups all common flags of all commands exept "config" and its sub-commands.
type Flags struct {
	OnlyID   bool
	OnlyURL  bool
	OnlyJSON bool
}

// NewFlags function creates a new Flags struct with default values.
func NewFlags() *Flags {
	return &Flags{
		OnlyID:   false,
		OnlyURL:  false,
		OnlyJSON: false,
	}
}

// FlagsSM struct groups all common flags of all "sm" sub-commands.
type FlagsSM struct {
	*Flags
	Shell string
}

// NewFlagsSM function creates a new FlagsSM struct with default values.
func NewFlagsSM() *FlagsSM {
	return &FlagsSM{
		Flags: NewFlags(),
		Shell: "",
	}
}

// FlagsSMShow struct groups all flags of the "sm show" commannd.
type FlagsSMShow struct {
	*FlagsSM
	ElementID  string
	ElementIdx int
	OnlyValue  bool
}

// NewFlagsSMShow function creates a new FlagsSMShow struct with default values.
func NewFlagsSMShow() *FlagsSMShow {
	return &FlagsSMShow{
		FlagsSM:    NewFlagsSM(),
		ElementID:  "",
		ElementIdx: -1,
		OnlyValue:  false,
	}
}

// FlagsSearch struct groups all flags of the "search" commmand.
type FlagsSearch struct {
	*Flags
	SMID       string
	ElementID  string
	ElementIdx int
	Value      string
}

// NewFlagsSearch function creates a new FlagsSearch struct with default values.
func NewFlagsSearch() *FlagsSearch {
	return &FlagsSearch{
		Flags:      NewFlags(),
		SMID:       "",
		ElementID:  "",
		ElementIdx: -1,
		Value:      "",
	}
}

// FlagsDiscover struct groups all flags of the "discover" command.
type FlagsDiscover struct {
	OnlyJSON bool
	OnlyURL  bool
}

// NewFlagsDiscover function creates a new FlagsDiscover struct with default values.
func NewFlagsDiscover() *FlagsDiscover {
	return &FlagsDiscover{
		OnlyJSON: false,
		OnlyURL:  false,
	}
}

// NewFlagsFromCMD function takes a command and creates a Flags struct from it.
func NewFlagsFromCMD(cmd *cobra.Command) (*Flags, error) {
	f := NewFlags()
	var err error
	f.OnlyID, err = cmd.Flags().GetBool("id")
	if err != nil {
		return nil, err
	}
	f.OnlyURL, err = cmd.Flags().GetBool("url")
	if err != nil {
		return nil, err
	}
	f.OnlyJSON, err = cmd.Flags().GetBool("json")
	if err != nil {
		return nil, err
	}
	return f, nil
}

// NewFlagsSMFromCMD function takes a command and creates a FlagsSM struct from it.
func NewFlagsSMFromCMD(cmd *cobra.Command) (*FlagsSM, error) {
	f := NewFlagsSM()
	var err error
	f.Flags, err = NewFlagsFromCMD(cmd)
	if err != nil {
		return nil, err
	}
	f.OnlyID, err = cmd.Flags().GetBool("id")
	if err != nil {
		return nil, err
	}
	f.OnlyURL, err = cmd.Flags().GetBool("url")
	if err != nil {
		return nil, err
	}
	f.OnlyJSON, err = cmd.Flags().GetBool("json")
	if err != nil {
		return nil, err
	}
	f.Shell, err = cmd.Flags().GetString("aas")
	if err != nil {
		return nil, err
	}
	return f, nil
}

// NewFlagsSMShowFromCMD function takes a command and creates a FlagsSMShow struct from it.
func NewFlagsSMShowFromCMD(cmd *cobra.Command) (*FlagsSMShow, error) {
	f := NewFlagsSMShow()
	var err error
	f.FlagsSM, err = NewFlagsSMFromCMD(cmd)
	if err != nil {
		return nil, err
	}
	f.ElementID, err = cmd.Flags().GetString("elementId")
	if err != nil {
		return nil, err
	}
	f.ElementIdx, err = cmd.Flags().GetInt("elementIdx")
	if err != nil {
		return nil, err
	}
	f.OnlyValue, err = cmd.Flags().GetBool("value")
	if err != nil {
		return nil, err
	}
	return f, nil
}

// NewFlagsSearchFromCMD takes a command and creates a FlagsSearch struct from it.
func NewFlagsSearchFromCMD(cmd *cobra.Command) (*FlagsSearch, error) {
	f := NewFlagsSearch()
	var err error
	f.Flags, err = NewFlagsFromCMD(cmd)
	if err != nil {
		return nil, err
	}
	f.SMID, err = cmd.Flags().GetString("sm")
	if err != nil {
		return nil, err
	}
	f.ElementID, err = cmd.Flags().GetString("elementId")
	if err != nil {
		return nil, err
	}
	f.ElementIdx, err = cmd.Flags().GetInt("elementIdx")
	if err != nil {
		return nil, err
	}
	f.Value, err = cmd.Flags().GetString("value")
	if err != nil {
		return nil, err
	}
	return f, nil
}

// NewFlagsDiscovery function takes a command and creates a FlagsDiscover struct from it.
func NewFlagsDiscoverFromCMD(cmd *cobra.Command) (*FlagsDiscover, error) {
	f := NewFlagsDiscover()
	var err error
	f.OnlyJSON, err = cmd.Flags().GetBool("json")
	if err != nil {
		return nil, err
	}
	f.OnlyURL, err = cmd.Flags().GetBool("url")
	if err != nil {
		return nil, err
	}
	return f, nil
}
