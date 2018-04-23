// Copyright © 2018 Anthony Spring <aspring@traackr.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/traackr/binnacle/config"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays the Helm status for each release within the given Binnacle configuration",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		statusCmdPreRun()
	},
	Run: func(cmd *cobra.Command, args []string) {
		statusCmdRun(args...)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		statusCmdPostRun()
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}

func statusCmdPreRun() {
	log.Debug("Executing `status` command.")
}

func statusCmdRun(args ...string) {

	// Load our configuration
	c, err := config.LoadAndValidateFromViper()

	if err != nil {
		log.Fatalf("unable to load configuration: %v", err)
	}

	var charts = c.Charts

	log.Debugf("Loaded %d charts.", len(charts))

	// Iterate the charts in the config
	for _, chart := range charts {
		var cmdArgs []string
		var res Result

		log.Debugf("Processing chart: %s", chart.ChartLongName())

		cmdArgs = append(cmdArgs, "status")
		cmdArgs = append(cmdArgs, chart.Release)
		cmdArgs = append(cmdArgs, args...)

		res, err = RunHelmCommand(cmdArgs...)
		if err != nil {
			log.Errorf("Helm status for release %s failed with the following:", chart.Release)
			log.Fatal(res.Stderr)
		}

		fmt.Println(strings.TrimSpace(res.Stdout))
	}
}

func statusCmdPostRun() {
	log.Debug("Execution of the `status` command has completed.")
}
