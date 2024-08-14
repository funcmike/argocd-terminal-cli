package argocd

import (
	"errors"
	"github.com/goccy/go-yaml"
	"os"
)

type ArgoCDContext struct {
	Name   string `yaml:"name"`
	Server string `yaml:"server"`
	User   string `yaml:"user"`
}

type ArgoCDConfig struct {
	Contexts       []ArgoCDContext
	CurrentContext string `yaml:"current-context"`
	Servers        []struct {
		Server string `yaml:"server"`
	} `json:"servers" yaml:"servers"`
	Users []struct {
		Name      string `yaml:"name"`
		AuthToken string `yaml:"auth-token"`
	}
}

func ParseArgoCDConfigFile(filepath string) (OptionsAuth, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return OptionsAuth{}, err
	}
	defer f.Close()

	var config ArgoCDConfig
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return OptionsAuth{}, err
	}

	var foundContext *ArgoCDContext
	for _, context := range config.Contexts {
		if context.Name == config.CurrentContext {
			foundContext = &context
		}
	}
	if foundContext == nil {
		return OptionsAuth{}, errors.New("context not found")
	}

	var oa OptionsAuth

	for _, server := range config.Servers {
		if foundContext.Server == server.Server {
			oa.ArgoCDServer = server.Server
		}
	}
	if oa.ArgoCDServer == "" {
		return OptionsAuth{}, errors.New("server not found")
	}

	for _, user := range config.Users {
		if foundContext.User == user.Name {
			oa.ArgoCDAuthToken = user.AuthToken
		}
	}
	if oa.ArgoCDAuthToken == "" {
		return OptionsAuth{}, errors.New("auth token not found")
	}

	return oa, nil
}
