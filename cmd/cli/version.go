package cli

import (
	"fmt"
	"github.com/mvazquezc/latency-tester/pkg/commands"
	"github.com/spf13/cobra"
)

var (
	short bool
)

func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !short {
				fmt.Printf("Build time: %s\n", commands.GetBuildTime())
				fmt.Printf("Git commit: %s\n", commands.GetGitCommit())
				fmt.Printf("Go commands: %s\n", commands.GetGoVersion())
				fmt.Printf("Go compiler: %s\n", commands.GetGoCompiler())
				fmt.Printf("Go Platform: %s\n", commands.GetGoPlatform())
				fmt.Printf("Version: %s\n", commands.PrintVersion())
			} else {
				fmt.Printf("%s\n", commands.PrintVersion())
			}
			return nil
		},
	}
	addVersionFlags(cmd)
	return cmd
}

func addVersionFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.BoolVar(&short, "short", false, "show only the version number")
}
