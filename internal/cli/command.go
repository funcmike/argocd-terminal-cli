package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/funcmike/argocd-terminal-cli/internal/argocd"
	"os"
	"path"
)

var errHelpWanted = errors.New("help wanted")

type command interface {
	Init([]string) error
	Name() string
	Run(ctx context.Context) error
}

type optionsDefault struct {
	ArgoConfigFilepath string
	argocd.OptionsAuth
	argocd.OptionsApp
}

func (o *optionsDefault) SetDefaults() error {
	o.ArgoCDServer = valueOrEnv(o.ArgoCDServer, "ARGOCD_SERVER")
	o.ArgoCDAuthToken = valueOrEnv(o.ArgoCDAuthToken, "ARGOCD_AUTH_TOKEN")

	if o.ArgoCDServer == "" || o.ArgoCDAuthToken == "" {
		auth, err := argocd.ParseArgoCDConfigFile(o.ArgoConfigFilepath)
		if err != nil {
			return fmt.Errorf("invalid argo configuration file: %w", err)
		}
		o.OptionsAuth = auth
	}
	return nil
}

func setDefaultArgs(flagSet *flag.FlagSet, od *optionsDefault) {
	flagSet.StringVar(&od.ArgoConfigFilepath, "config", valueOrDefault(os.Getenv("ARGOCD_CONFIG_FILE"), path.Join(os.Getenv("HOME"), ".config/argocd/config")), "Path to ArgoCD config file")
	flagSet.StringVar(&od.ArgoCDServer, "server", "", "ArgoCD server URL")
	flagSet.StringVar(&od.ArgoCDAuthToken, "auth-token", "", "ArgoCD Auth token")

	flagSet.StringVar(&od.AppName, "app-name", "", "ArgoCD app name")
	flagSet.StringVar(&od.AppNamespace, "app-namespace", "core-argocd", "ArgoCD app namespace name")
}

func setK8sArgs(flagSet *flag.FlagSet, ok *argocd.OptionsK8s) {
	flagSet.StringVar(&ok.Namespace, "namespace", "", "Kubernetes app namespace name")
}

func helpWanted(args []string) bool {
	if (len(args) == 0) || args[0] == "--help" || args[0] == "-h" {
		return true
	}
	return false
}
