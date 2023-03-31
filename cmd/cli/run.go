package cli

import (
	"errors"
	"github.com/mvazquezc/latency-tester/pkg/commands"
	"github.com/mvazquezc/latency-tester/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	target       string
	numberOfRuns int
	waitInterval string
	outputFormat string
)

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Exec the run command",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate command Args
			err := validateRunCommandArgs()
			if err != nil {
				return err
			}
			// We have the run command logic implemented in our example pkg
			waitIntervalSeconds := utils.IntervalTimeToSeconds(waitInterval)
			latencyTestOutput, err := commands.RunLatencyTestCmd(target, numberOfRuns, waitIntervalSeconds)
			if err != nil {
				return err
			}
			switch {
			case outputFormat == "yaml":
				utils.WriteOutputYaml(latencyTestOutput)
			case outputFormat == "json":
				utils.WriteOutputJson(latencyTestOutput)
			default:
				utils.WriteOutputTable(latencyTestOutput)
			}
			return err
		},
	}
	addRunCommandFlags(cmd)
	return cmd
}

func addRunCommandFlags(cmd *cobra.Command) {

	flags := cmd.Flags()
	flags.StringVarP(&target, "target", "t", "", "The target for the test. Supports http/s and tcp. e.g: https://google.com | tcp://127.0.0.1:3000")
	flags.IntVarP(&numberOfRuns, "runs", "r", 1, "The number of executions.")
	flags.StringVarP(&waitInterval, "interval", "i", "1m", "The amount of waiting time between runs. Allowed values (<num>s, <num>m, <num>h")
	flags.StringVarP(&outputFormat, "output-format", "o", "table", "Output in an specific format. Usage: '-o [ table | yaml | json ]'")
	cmd.MarkFlagRequired("target-url")
}

// validateCommandArgs validates that arguments passed by the user are valid
func validateRunCommandArgs() error {
	validInterval := utils.ValidateIntervalTime(waitInterval)
	if !validInterval {
		return errors.New("Wait interval is not valid")
	}
	// Validate URL is valid
	validTarget := utils.ValidateTarget(target)
	if !validTarget {
		return errors.New("Target is not valid")
	}
	return nil
}
