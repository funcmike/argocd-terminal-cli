package cli

import (
	"context"
	"fmt"
	"os"
)

var _ command = &commandSection{}

type commandSection struct {
	SectionName string
	Subcommands []command
	ExtraHelp   func(string)

	cmd command
}

func (c *commandSection) Init(args []string) error {
	if helpWanted(args) {
		c.Usage()
		return errHelpWanted
	}

	for _, cmd := range c.Subcommands {
		if cmd.Name() == args[0] {
			if err := cmd.Init(args[1:]); err != nil {
				return err
			}
			c.cmd = cmd
			return nil
		}
		if cmd.Name() == "" {
			if err := cmd.Init(args); err != nil {
				return err
			}
			c.cmd = cmd
			return nil
		}
	}

	_, _ = fmt.Fprintf(os.Stderr, "command not found: %s\n\n", args[0])
	c.Usage()
	return errHelpWanted
}

func (c *commandSection) Name() string {
	return c.SectionName
}

func (c *commandSection) Run(ctx context.Context) error {
	if c.cmd == nil {
		return fmt.Errorf("command not initialized: %s", c.Name())
	}
	return c.cmd.Run(ctx)
}

func (c *commandSection) Usage() {
	fmt.Printf("Usage of %s:\n", c.Name())
	for _, cmd := range c.Subcommands {
		if cmd.Name() != "" {
			fmt.Println(" ", c.Name(), cmd.Name(), "[arguments...]")
		}
	}

	if c.ExtraHelp != nil {
		c.ExtraHelp(c.Name())
	}
}
