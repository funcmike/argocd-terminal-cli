package cli

import (
	"context"
	"flag"
	"github.com/funcmike/argocd-terminal-cli/internal/argocd"
	"github.com/funcmike/argocd-terminal-cli/internal/term"
)

var _ command = &commandTerm{}

func newCommandTerm() *commandTerm {
	ct := &commandTerm{flagSet: flag.NewFlagSet("term", flag.ContinueOnError)}
	setDefaultArgs(ct.flagSet, &ct.optionsDefault)
	setK8sArgs(ct.flagSet, &ct.OptionsK8s)
	ct.flagSet.StringVar(&ct.projectName, "project-name", "", "ArgoCD app project name")
	ct.flagSet.StringVar(&ct.pod, "pod", "", "POD name")
	ct.flagSet.StringVar(&ct.container, "container", "", "Container name")
	return ct
}

type commandTerm struct {
	flagSet *flag.FlagSet

	optionsDefault
	argocd.OptionsK8s

	projectName string
	container   string
	pod         string
}

func (c *commandTerm) Init(args []string) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}
	c.projectName = valueOrDefault(c.projectName, c.AppName)
	return c.SetDefaults()
}

func (c *commandTerm) Name() string {
	return c.flagSet.Name()
}

func (c *commandTerm) Run(ctx context.Context) error {
	return term.Run(ctx, argocd.TerminalClientOptions{
		ArgoCDServer: c.OptionsAuth.ArgoCDServer,
		OptionsApp:   c.OptionsApp,
		OptionsK8s:   c.OptionsK8s,
		ProjectName:  c.projectName,
		Container:    c.container,
		POD:          c.pod,
	}, c.ArgoCDAuthToken)
}
