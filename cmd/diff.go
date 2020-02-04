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
	"io/ioutil"
	"os"
	"strings"

	"github.com/Traackr/binnacle/config"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Displays a diff between the current release and new release of a Helm chart.  (Requires helm-diff plugin)",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		diffCmdPreRun()
	},
	Run: func(cmd *cobra.Command, args []string) {
		diffCmdRun(args...)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		diffCmdPostRun()
	},
}

func init() {
	RootCmd.AddCommand(diffCmd)
}

func diffCmdPreRun() {
	log.Debug("Executing `diff` command.")
}

func diffCmdRun(args ...string) {
	var err error

	// Detect if the diff plugin is installed.
	pluginInstalled, err := PluginInstalled("diff")
	if err != nil {
		log.Fatalf("error detecting if diff plugin is installed: %v", err)
	}

	if !pluginInstalled {
		log.Fatalf("helm-diff plugin is required. Please see: https://github.com/databus23/helm-diff")
	}

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

		log.Debugf("Processing chart: %s", chart.ChartURL())

		// Create a temp working directory
		dir, err := ioutil.TempDir("", "binnacle-exec")
		if err != nil {
			log.Fatalf("error creating temp directory: %v", err)
		}
		defer os.RemoveAll(dir)

		var valuesFile = dir + "/values.yml"

		//
		// Template out the charts values
		//
		if err = chart.WriteValueFile(valuesFile); err != nil {
			log.Fatal(err)
		}

		cmdArgs = append(cmdArgs, "diff")
		cmdArgs = append(cmdArgs, "upgrade")
		cmdArgs = append(cmdArgs, chart.Release)
		cmdArgs = append(cmdArgs, chart.ChartURL())
		cmdArgs = append(cmdArgs, "--values")
		cmdArgs = append(cmdArgs, valuesFile)

		if len(chart.Version) > 0 {
			cmdArgs = append(cmdArgs, "--version")
			cmdArgs = append(cmdArgs, chart.Version)
		}

		cmdArgs = append(cmdArgs, args...)
		res, err = RunHelmCommand(cmdArgs...)
		if err != nil {
			log.Errorf("helm diff upgrade for release %s failed with the following:", chart.Release)
			log.Fatal(res.Stderr)
		}

		fmt.Println(strings.TrimSpace(res.Stdout))
	}
}

func diffCmdPostRun() {
	log.Debug("Execution of the `diff` command has completed.")
}
