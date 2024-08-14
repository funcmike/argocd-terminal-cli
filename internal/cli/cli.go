package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/goccy/go-yaml"
	"os"
)

type outFormatType string

const outFormatJSON outFormatType = "json"
const outFormatYAML outFormatType = "yaml"

func Run(args []string, ctx context.Context) error {
	cs := &commandSection{
		SectionName: "atc",
		Subcommands: []command{
			newCommandTerm(),
			newCommandGet(),
		},
	}

	if err := cs.Init(args); err != nil {
		if errors.Is(err, errHelpWanted) {
			return flag.ErrHelp
		}
		return err
	}
	return cs.Run(ctx)
}

func output(data json.RawMessage, format outFormatType) error {
	if format == outFormatYAML {
		out, err := yaml.JSONToYAML(data)
		if err != nil {

			return fmt.Errorf("marshalling pod manifest: %w", err)
		}

		_, err = fmt.Fprintf(os.Stdout, "%s\n", out)
		return err
	}

	if format == outFormatJSON {
		_, err := fmt.Fprintf(os.Stdout, "%s\n", data)
		return err
	}

	return fmt.Errorf("unknown output format: %s", format)
}

func valueOrEnv(value string, envName string) string {
	if value != "" {
		return value
	}
	return os.Getenv(envName)
}

func valueOrDefault(value string, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}
