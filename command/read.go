package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

// ReadCommand is a Command that reads data from the Vault.
type ReadCommand struct {
	Meta
	List bool
}

func (c *ReadCommand) Run(args []string) int {
	var format string
	var field string
	var err error
	var secret *api.Secret
	var flags *flag.FlagSet
	if c.List {
		flags = c.Meta.FlagSet("list", FlagSetDefault)
	} else {
		flags = c.Meta.FlagSet("read", FlagSetDefault)
	}
	flags.StringVar(&format, "format", "table", "")
	flags.StringVar(&field, "field", "", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		c.Ui.Error("read expects one argument")
		flags.Usage()
		return 1
	}

	path := args[0]
	if path[0] == '/' {
		path = path[1:]
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	if c.List {
		secret, err = client.Logical().List(path)
	} else {
		secret, err = client.Logical().Read(path)
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error reading %s: %s", path, err))
		return 1
	}
	if secret == nil {
		c.Ui.Error(fmt.Sprintf(
			"No value found at %s", path))
		return 1
	}

	// Handle single field output
	if field != "" {
		if val, ok := secret.Data[field]; ok {
			c.Ui.Output(fmt.Sprintf("%s", val))
			return 0
		} else {
			c.Ui.Error(fmt.Sprintf(
				"Field %s not present in secret", field))
			return 1
		}
	}

	return OutputSecret(c.Ui, format, secret)
}

func (c *ReadCommand) Synopsis() string {
	if c.List {
		return "List data in Vault"
	}
	return "Read data or secrets from Vault"
}

func (c *ReadCommand) Help() string {
	helpText := `
Usage: vault read [options] path

  Read data from Vault.

  Reads data at the given path from Vault. This can be used to read
  secrets and configuration as well as generate dynamic values from
  materialized backends. Please reference the documentation for the
  backends in use to determine key structure.
`

	if c.List {
		helpText =
			`
Usage: vault list [options] path

  List data from Vault.

  Retrieve a listing of available data. The data returned is
  backend-specific, and not all backends implement listing capability.
`
	}

	helpText += `
General Options:

  ` + generalOptionsUsage() + `

Read Options:

  -format=table           The format for output. By default it is a whitespace-
                          delimited table. This can also be json.

  -field=field            If included, the raw value of the specified field
                          will be output raw to stdout.

`
	return strings.TrimSpace(helpText)
}
