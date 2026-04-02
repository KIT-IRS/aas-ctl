package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	aasjsonization "github.com/aas-core-works/aas-core3.0-golang/jsonization"
	aasstringification "github.com/aas-core-works/aas-core3.0-golang/stringification"
	aastypes "github.com/aas-core-works/aas-core3.0-golang/types"
)

// Print function prints a Referable to the Stdout
func Print(r aastypes.IReferable) {
	Fprint(os.Stdout, r)
}

// PrintVerbose function prints detailed informaiton to a Referable to the Stdout
func PrintVerbose(r aastypes.IReferable) {
	FprintVerbose(os.Stdout, r)
}

// PrintValue function prints the value of a Referable to the Stdout
func PrintValue(r aastypes.IReferable) {
	FprintValue(os.Stdout, r)
}

// PrintIdentifiable function prints a Identifiable according to the provided flags
func PrintIdentifiable(i aastypes.IIdentifiable, flags *Flags, verbose bool) {
	if flags.OnlyID {
		PrintID(i)
	} else if flags.OnlyURL {
		PrintURL(i)
	} else if flags.OnlyJSON {
		PrintJSON(i)
	} else if verbose {
		PrintVerbose(i)
	} else {
		Print(i)
	}
}

// PrintSubmodelElement function prints a SubmodelElement according to the provided flags
func PrintSubmodelElement(sm aastypes.ISubmodel, e aastypes.ISubmodelElement, flags *FlagsSMShow) {
	if flags.OnlyURL {
		PrintElementURL(sm, e)
	} else if flags.OnlyJSON {
		PrintJSON(e)
	} else if flags.OnlyValue {
		PrintValue(e)
	} else {
		PrintVerbose(e)
	}
}

// PrintID function prints the ID of a Identifiable to the Stdout
func PrintID(i aastypes.IIdentifiable) {
	FprintID(os.Stdout, i)
}

// PrintURL function prints the URL of a Identifiable to the Stdout
func PrintURL(i aastypes.IIdentifiable) {
	FprintURL(os.Stdout, i)
}

// PrintElementURL function prints the URL of a SubmodelElement to the Stdout
func PrintElementURL(sm aastypes.ISubmodel, e aastypes.ISubmodelElement) {
	FprintElementURL(os.Stdout, sm, e)
}

// PrintJSON function prints the JSON representation of an AAS element to the Stdout
func PrintJSON(i aastypes.IClass) {
	FprintJSON(os.Stdout, i)
}

// Fprint function prints a Referable to the given writer for better testability
func Fprint(w io.Writer, r aastypes.IReferable) {
	var str string
	switch r := r.(type) {
	case aastypes.IIdentifiable:
		str = formatIdentifiable(r)
	case aastypes.ISubmodelElement:
		str = formatSubmodelElement(r)
	default:
		str = fmt.Sprintf("Print function is not defined for %t\n", r)
	}
	fmt.Fprint(w, str)
}

// print function prints detailed informaiton to a Referable to the given writer for better testability
func FprintVerbose(w io.Writer, r aastypes.IReferable) {
	var str string
	switch r := r.(type) {
	case aastypes.IAssetAdministrationShell:
		str = formatAssetAdministrationShellVerbose(r)
	case aastypes.ISubmodel:
		str = formatSubmodelVerbose(r)
	case aastypes.ISubmodelElement:
		str = formatSubmodelElementVerbose(r)
	default:
		str = fmt.Sprintf("Verbose print function is not defined for %t\n", r)
	}
	fmt.Fprint(w, str)
}

// FprintValue function prints the value of a Referable to the given writer for better testability
func FprintValue(w io.Writer, r aastypes.IReferable) {
	var str string
	switch r := r.(type) {
	case aastypes.IAssetAdministrationShell, aastypes.ISubmodel:
		str = "The --value flag is only usable with the --elementId or the --elementIdx flag"
	case aastypes.IMultiLanguageProperty:
		str = r.Value()[0].Text()
	case aastypes.IProperty:
		str = *r.Value()
	default:
		str = "Value of this element is not printable or the element of this type does not have a value attribute"
	}
	fmt.Fprintln(w, str)
}

// FprintID function prints the ID of a Identifiable to the specified Writer
// for better testability
func FprintID(w io.Writer, i aastypes.IIdentifiable) {
	fmt.Fprint(w, i.ID())
}

// FprintURL function prints the URL of a Identifiable to the specified Writer
// for better testability
func FprintURL(w io.Writer, i aastypes.IIdentifiable) {
	endpoint, err := getEndpoint(i)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, endpoint)
}

// FprintElementURL function prints the URL of a SubmodelElement to the specified Writer
// for better testability
func FprintElementURL(w io.Writer, sm aastypes.ISubmodel, e aastypes.ISubmodelElement) {
	endpoint, err := getElementEndpoint(sm, e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, endpoint)
}

// FprintJSON function prints the JSON representation of an AAS element to the specified Writer
// for better testability
func FprintJSON(w io.Writer, i aastypes.IClass) {
	jsonable, err := aasjsonization.ToJsonable(i)
	if err != nil {
		log.Fatal(err)
	}
	jsonBytes, err := json.Marshal(jsonable)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(jsonBytes))
}

// formatReferable function returns the base information of the given Referable as string
func formatReferable(r aastypes.IReferable) string {
	str := ""
	if r.IDShort() != nil {
		str = str + *r.IDShort()
	}
	modelType, ok := aasstringification.ModelTypeToString(r.ModelType())
	if !ok {
		modelType = "UNKNOWN"
	}
	return str + "<" + modelType + ">: "
}

// formatIdentifiable function returns the base information of the given Identifiable as string
func formatIdentifiable(i aastypes.IIdentifiable) string {
	str := i.ID()
	modelType, ok := aasstringification.ModelTypeToString(i.ModelType())
	if !ok {
		modelType = "UNKNOWN"
	}
	str = str + "<" + modelType + ">"
	if i.IDShort() != nil {
		str = str + " " + *i.IDShort()
	}
	return str
}

// formatAssetAdministrationShellVerbose function formats the information of the given AAS into a string
// All Submodels of the shell are formated as Identifiable in a list
func formatAssetAdministrationShellVerbose(aas aastypes.IAssetAdministrationShell) string {
	str := formatIdentifiable(aas)
	for _, ref := range aas.Submodels() {
		sm, err := getSubmodelFromReference(ref)
		if err != nil {
			log.Printf("WARNING: Unable to get Submodel from reference %v; %v", ref.Keys()[0].Value(), err)
		}
		str = str + "\n" + formatIdentifiable(sm)
	}
	return str
}

// formatSubmodelVerbose function formats the information of the given Submodel into a string
// All SubmodelElements of the Submodel are formatted in a list
func formatSubmodelVerbose(sm aastypes.ISubmodel) string {
	str := formatIdentifiable(sm)
	for i, element := range sm.SubmodelElements() {
		str = fmt.Sprintf("%v\n%2d\t%v", str, i, formatSubmodelElement(element))
	}
	return str
}

// formatSubmodelElement function determines the type of the submodel element and formats it into a corresponding string
func formatSubmodelElement(e aastypes.ISubmodelElement) string {
	switch e := e.(type) { // TODO: All other element types
	case aastypes.IMultiLanguageProperty:
		return formatMultiLanguageProperty(e)
	case aastypes.IProperty:
		return formatProperty(e)
	case aastypes.IRange:
		return formatRange(e)
	case aastypes.ISubmodelElementList:
		return formatSubmodelElementList(e)
	case aastypes.ISubmodelElementCollection: // Order is relevant, SEC must be after SEL
		return formatSubmodelElementCollection(e)
	default:
		return fmt.Sprintf("%v Formatting is not implemented for SubmodelElement of that type", formatReferable(e))
	}
}

// formatSubmodelElementVerbose function formats the contents of SML and SMC
// if the SME is no collection or list it gets formated corresponding to formatSubmodelElement
func formatSubmodelElementVerbose(e aastypes.ISubmodelElement) string {
	switch e := e.(type) {
	case aastypes.ISubmodelElementList:
		return formatSubmodelElementListVerbose(e)
	case aastypes.ISubmodelElementCollection:
		return formatSubmodelElementCollectionVerbose(e)
	default:
		return formatSubmodelElement(e)
	}
}

// formatMultiLanguageProperty function formats the contents of a MLP to a string
// only the first element of the MLP is considered for the formatting
func formatMultiLanguageProperty(mlp aastypes.IMultiLanguageProperty) string {
	str := formatReferable(mlp)
	if len(mlp.Value()) > 0 {
		return str + mlp.Value()[0].Text() + "<" + mlp.Value()[0].Language() + ">"
	} else {
		return str + "empty"
	}
}

// formatProperty function formats the information of a Property to a string
func formatProperty(p aastypes.IProperty) string {
	str := formatReferable(p)
	if p.Value() != nil {
		str = str + *p.Value()
	}
	dataType, ok := aasstringification.DataTypeDefXSDToString(p.ValueType())
	if !ok {
		dataType = "UNKNOWN"
	}
	return str + "<" + dataType + ">"
}

// formatRange function formats the information of a Range to a string
func formatRange(r aastypes.IRange) string {
	var min, max string
	if r.Min() != nil {
		min = *r.Min()
	} else {
		min = "NaN"
	}
	if r.Max() != nil {
		max = *r.Max()
	} else {
		max = "NaN"
	}
	dataType, ok := aasstringification.DataTypeDefXSDToString(r.ValueType())
	if !ok {
		dataType = "UNKNOWN"
	}
	return formatReferable(r) + min + "-" + max + "<" + dataType + ">"
}

// formatSubmodelElementCollection function formats a SMC to a string
// only number of elements in the collectio are given
func formatSubmodelElementCollection(sec aastypes.ISubmodelElementCollection) string {
	return formatReferable(sec) + fmt.Sprintf("%d Elements", len(sec.Value()))
}

// formatSubmodelElementCollectionVerbose function formats a SMC to a string
// the information of each element in the SMC is added as well
func formatSubmodelElementCollectionVerbose(sec aastypes.ISubmodelElementCollection) string {
	str := formatSubmodelElementCollection(sec)
	for i, e := range sec.Value() {
		str = fmt.Sprintf("%v\n%2d\t%v", str, i, formatSubmodelElement(e))
	}
	return str
}

// formatSubmodelElementList function formats the information of a SML to a string
// only number of elements and type of elements in the list are given
func formatSubmodelElementList(sel aastypes.ISubmodelElementList) string {
	smeType, ok := aasstringification.AASSubmodelElementsToString(sel.TypeValueListElement())
	if !ok {
		smeType = "UNKNOWN"
	}
	return formatReferable(sel) + fmt.Sprintf("%d Elements<%v>", len(sel.Value()), smeType)
}

// formatSubmodelElementListVerbose function formats the infotmation of a SML to a string
// the information of each element in the SML is added as well
func formatSubmodelElementListVerbose(sel aastypes.ISubmodelElementList) string {
	str := formatSubmodelElementList(sel)
	for _, e := range sel.Value() {
		str = str + "\n" + formatSubmodelElement(e)
	}
	return str
}

// formatKeyValuePair function formats the first entries of a map[string]any to a string "key: value"
// Useful for raw MultiLanguageProperties
func formatKeyValuePair(pair map[string]any) string {
	for k, v := range pair {
		return fmt.Sprintf("%v:\t%v", k, v)
	}
	return ""
}
