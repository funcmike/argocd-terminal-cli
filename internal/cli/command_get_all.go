package cli

import (
	"context"
	"flag"
	"github.com/funcmike/argocd-terminal-cli/internal/argocd"
	"net/http"
)

var _ command = &commandGetAll{}

func newCommandGetAll() *commandGetAll {
	cg := &commandGetAll{flagSet: flag.NewFlagSet("all", flag.ContinueOnError)}
	setDefaultArgs(cg.flagSet, &cg.optionsDefault)
	cg.flagSet.StringVar(&cg.outputFormat, "output", "yaml", "Output format: yaml or json")
	return cg
}

type commandGetAll struct {
	optionsDefault

	flagSet      *flag.FlagSet
	outputFormat string
}

func (c *commandGetAll) Init(args []string) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}
	return c.SetDefaults()
}

func (c *commandGetAll) Name() string {
	return c.flagSet.Name()
}

func (c *commandGetAll) Run(ctx context.Context) error {
	client, err := argocd.NewAPIClient(c.OptionsAuth, http.DefaultClient)
	if err != nil {
		return err
	}
	resources, err := client.GetResources(ctx, c.AppName, c.AppNamespace)
	if err != nil {
		return err
	}
	return output(resources, outFormatType(c.outputFormat))
}
