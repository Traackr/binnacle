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

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Templates out each release, with values, from a given Binnacle configuration.",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		templateCmdPreRun()
	},
	Run: func(cmd *cobra.Command, args []string) {
		templateCmdRun(args...)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		templateCmdPostRun()
	},
}

func init() {
	RootCmd.AddCommand(templateCmd)
}

func templateCmdPreRun() {
	log.Debug("Executing `template` command.")
}

func templateCmdRun(args ...string) {
	// Load our configuration
	c, err := config.LoadAndValidateFromViper()

	if err != nil {
		log.Fatalf("unable to load configuration: %v", err)
	}

	// Sync repositories
	if err := syncRepositories(c.Repositories, args...); err != nil {
		log.Fatal(err)
	}

	var charts = c.Charts

	var absentCharts []string

	log.Debugf("Loaded %d charts.", len(charts))

	// Iterate the charts in the config
	for _, chart := range charts {
		var cmdArgs []string
		var res Result

		log.Debugf("Processing chart: %s", chart.ChartURL())

		// If the state is not set to present add the namespace/release to the not rendered list
		if chart.State != config.StatePresent {
			absentCharts = append(absentCharts, chart.Namespace+"/"+chart.Release)
			continue
		}

		// Create a temp working directory
		dir, err := ioutil.TempDir("", "binnacle-exec")
		if err != nil {
			log.Fatalf("error creating temp directory: %v", err)
		}
		defer os.RemoveAll(dir)

		var valuesFile = dir + "/values.yml"

		//
		// Fetch the chart for helm2
		//
		if IsHelm2() {
			cmdArgs = append(cmdArgs, "fetch")
			cmdArgs = append(cmdArgs, chart.ChartURL())
			cmdArgs = append(cmdArgs, "--destination")
			cmdArgs = append(cmdArgs, dir)

			cmdArgs = append(cmdArgs, "--untar")

			if len(chart.Version) > 0 {
				cmdArgs = append(cmdArgs, "--version")
				cmdArgs = append(cmdArgs, chart.Version)
			}

			res, err = RunHelmCommand(cmdArgs...)
			if err != nil {
				log.Errorf("helm fetch failed with the following:")
				log.Fatal(res.Stderr)
			}
		}

		//
		// Template out the charts values
		//
		if err = chart.WriteValueFile(valuesFile); err != nil {
			log.Fatal(err)
		}

		//
		// Template against the chart
		//
		cmdArgs = nil

		cmdArgs = append(cmdArgs, "template")

		if IsHelm2() {
			cmdArgs = append(cmdArgs, dir+"/"+chart.Name)

			// Add the namespace if given
			if len(chart.Namespace) > 0 {
				cmdArgs = append(cmdArgs, "--namespace")
				cmdArgs = append(cmdArgs, chart.Namespace)
			}

			cmdArgs = append(cmdArgs, "--name")
			cmdArgs = append(cmdArgs, chart.Release)

			cmdArgs = append(cmdArgs, "--values")
			cmdArgs = append(cmdArgs, valuesFile)
		} else {
			// NAME
			cmdArgs = append(cmdArgs, chart.Release)

			// CHART
			cmdArgs = append(cmdArgs, chart.ChartURL())

			// Add the namespace if given
			if len(chart.Namespace) > 0 {
				cmdArgs = append(cmdArgs, "--namespace")
				cmdArgs = append(cmdArgs, chart.Namespace)
			}

			cmdArgs = append(cmdArgs, "--values")
			cmdArgs = append(cmdArgs, valuesFile)

			if len(chart.Version) > 0 {
				cmdArgs = append(cmdArgs, "--version")
				cmdArgs = append(cmdArgs, chart.Version)
			}
		}

		cmdArgs = append(cmdArgs, args...)

		res, err = RunHelmCommand(cmdArgs...)
		if err != nil {
			log.Errorf("helm template failed with the following:")
			log.Fatal(res.Stderr)
		}

		fmt.Println(strings.TrimSpace(res.Stdout))
	}

	// Display output about the released that were not rendered
	if len(absentCharts) > 0 {
		log.Info("The following releases were set to absent and were not rendered.")
		for _, chart := range absentCharts {
			log.Infof("  %s", chart)
		}
	}
}

func templateCmdPostRun() {
	log.Debug("Execution of the `template` command has completed.")
}
