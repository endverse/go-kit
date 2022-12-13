package tmpl

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/zeus-techs/go-kit/term"

	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
)

var helpTmpl = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

var usageTmpl = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}
`

var useCommandTmpl = `{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

var templateFuncs = template.FuncMap{
	"trim":                    strings.TrimSpace,
	"trimRightSpace":          trimRightSpace,
	"trimTrailingWhitespaces": trimRightSpace,
	"appendIfNotPresent":      appendIfNotPresent,
	"rpad":                    rpad,
	"gt":                      Gt,
	"eq":                      Eq,
}

func HelpTmpl(cmd *cobra.Command) {
	tmpl(cmd.OutOrStderr(), helpTmpl, cmd)
}

func UsageTmpl(cmd *cobra.Command) {
	tmpl(cmd.OutOrStderr(), usageTmpl, cmd)
}

func UseCommandTmpl(cmd *cobra.Command) {
	tmpl(cmd.OutOrStderr(), useCommandTmpl, cmd)
}

func tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("top")
	t.Funcs(templateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// FIXME appendIfNotPresent is unused by cobra and should be removed in a version 2. It exists only for compatibility with users of cobra.

// appendIfNotPresent will append stringToAppend to the end of s, but only if it's not yet present in s.
func appendIfNotPresent(s, stringToAppend string) string {
	if strings.Contains(s, stringToAppend) {
		return s
	}
	return s + " " + stringToAppend
}

// rpad adds padding to the right of a string.
func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}

// FIXME Gt is unused by cobra and should be removed in a version 2. It exists only for compatibility with users of cobra.

// Gt takes two types and checks whether the first type is greater than the second. In case of types Arrays, Chans,
// Maps and Slices, Gt will compare their lengths. Ints are compared directly while strings are first parsed as
// ints and then compared.
func Gt(a interface{}, b interface{}) bool {
	var left, right int64
	av := reflect.ValueOf(a)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		left = int64(av.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		left = av.Int()
	case reflect.String:
		left, _ = strconv.ParseInt(av.String(), 10, 64)
	}

	bv := reflect.ValueOf(b)

	switch bv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		right = int64(bv.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		right = bv.Int()
	case reflect.String:
		right, _ = strconv.ParseInt(bv.String(), 10, 64)
	}

	return left > right
}

// FIXME Eq is unused by cobra and should be removed in a version 2. It exists only for compatibility with users of cobra.

// Eq takes two types and checks whether they are equal. Supported types are int and string. Unsupported types will panic.
func Eq(a interface{}, b interface{}) bool {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		panic("Eq called on unsupported type")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return av.Int() == bv.Int()
	case reflect.String:
		return av.String() == bv.String()
	}
	return false
}

func UsageFunc(cmd *cobra.Command, fn func(c *cobra.Command) error) {

}

func globalFlags(cmd *cobra.Command, namedFlagSets cliflag.NamedFlagSets) cliflag.NamedFlagSets {
	namedFlagSets.FlagSet("global").AddFlagSet(cmd.PersistentFlags())
	cmd.VisitParents(func(c *cobra.Command) {
		namedFlagSets.FlagSet("global").AddFlagSet(c.PersistentFlags())
	})

	return namedFlagSets
}

func SetHelpAndUsageFunc(cmd *cobra.Command, namedFlagSets cliflag.NamedFlagSets) {
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())

	cmd.SetUsageFunc(func(c *cobra.Command) error {
		namedFlagSets = globalFlags(c, namedFlagSets)
		UsageTmpl(c)
		cliflag.PrintSections(c.OutOrStderr(), namedFlagSets, cols)
		UseCommandTmpl(c)
		return nil
	})

	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		HelpTmpl(c)
	})
}
