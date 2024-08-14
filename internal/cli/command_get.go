package cli

import "fmt"

var _ command = &commandGet{}

func newCommandGet() *commandGet {
	cg := &commandGet{
		commandSection{
			SectionName: "get",
			Subcommands: []command{
				newCommandGetAll(),
				newCommandGetResource(),
			},
		},
	}
	cg.commandSection.ExtraHelp = cg.ExtraHelp
	return cg
}

type commandGet struct {
	commandSection
}

func (c *commandGet) ExtraHelp(parentCommand string) {
	fmt.Println(" ", parentCommand, "<kind> <name> [arguments...]  - example: get pod nginx")
}
