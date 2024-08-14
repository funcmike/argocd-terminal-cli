package cli

import (
	"context"
	"flag"
	"fmt"
	"github.com/funcmike/argocd-terminal-cli/internal/argocd"
	"net/http"
	"unicode"
)

var _ command = &commandGetResource{}

func newCommandGetResource() *commandGetResource {
	cg := &commandGetResource{flagSet: flag.NewFlagSet("", flag.ContinueOnError)}
	setDefaultArgs(cg.flagSet, &cg.optionsDefault)
	setK8sArgs(cg.flagSet, &cg.OptionsK8s)
	cg.flagSet.StringVar(&cg.outputFormat, "output", "yaml", "Output format: yaml or json")
	cg.flagSet.StringVar(&cg.version, "version", "v1", "Resource API version")
	cg.flagSet.StringVar(&cg.group, "group", "", "Resource group")
	return cg
}

type commandGetResource struct {
	optionsDefault
	argocd.OptionsK8s

	flagSet *flag.FlagSet

	outputFormat string
	version      string

	kind  string
	name  string
	group string
}

func (c *commandGetResource) Init(args []string) error {
	c.kind = args[0]

	if helpWanted(args[1:]) {
		fmt.Printf("Usage of %s:\n %s <%s-name> [arguments...]\n", c.kind, c.kind, c.kind)
		c.flagSet.PrintDefaults()
		return errHelpWanted
	}
	c.name = args[1]
	if err := c.flagSet.Parse(args[2:]); err != nil {
		return err
	}
	return c.SetDefaults()
}

func (c *commandGetResource) Name() string {
	return c.flagSet.Name()
}

func (c *commandGetResource) Run(ctx context.Context) error {
	client, err := argocd.NewAPIClient(c.OptionsAuth, http.DefaultClient)
	if err != nil {
		return err
	}
	resource, err := client.GetResourceManifest(ctx, c.AppName, c.AppNamespace, c.Namespace, toUpperFirst(c.kind), c.version, c.name, c.group)
	if err != nil {
		return err
	}
	return output(resource, outFormatType(c.outputFormat))
}

func toUpperFirst(v string) string {
	r := []rune(v)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
