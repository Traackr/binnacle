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

// syncCmd represents the status command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs each release within the given Binnacle configuration with `helm`",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		syncCmdPreRun()
	},
	Run: func(cmd *cobra.Command, args []string) {
		syncCmdRun(args...)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		syncCmdPostRun()
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
}

func syncCmdPreRun() {
	log.Debug("Executing `sync` command.")
}

func syncCmdRun(args ...string) {
	// Load our configuration
	c, err := config.LoadAndValidateFromViper()
	if err != nil {
		log.Fatalf("unable to load configuration: %v", err)
	}

	// Sync repositories
	if err := syncRepositories(c.Repositories, args...); err != nil {
		log.Fatal(err)
	}

	// Sync charts
	if err := syncCharts(c.Charts, args...); err != nil {
		log.Fatal(err)
	}
}

func syncCmdPostRun() {
	log.Debug("Execution of the `sync` command has completed.")
}

func syncCharts(charts []config.ChartConfig, args ...string) error {
	for _, chart := range charts {
		var cmdArgs []string
		var res Result

		log.Debugf("Processing chart: %s", chart.ChartURL())

		if chart.State == config.StatePresent {
			// Create a temp working directory
			dir, err := ioutil.TempDir("", "binnacle-exec")
			if err != nil {
				return err
			}
			defer os.RemoveAll(dir)
			var valuesFile = dir + "/values.yml"

			//
			// Template out the charts values
			//
			if err = chart.WriteValueFile(valuesFile); err != nil {
				return err
			}

			cmdArgs = append(cmdArgs, "upgrade")
			cmdArgs = append(cmdArgs, chart.Release)
			cmdArgs = append(cmdArgs, chart.ChartURL())
			cmdArgs = append(cmdArgs, "-i")

			if IsHelm2() {
				cmdArgs = append(cmdArgs, "--force")
			}

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
		} else {

			// If the release does not exist do not attempt to delete the release
			exists := ReleaseExists(chart.Namespace, chart.Release)
			if !exists {
				log.Infof("Skipping '%s/%s' as the release does not exist.", chart.Namespace, chart.Release)
				continue
			}

			if IsHelm2() {
				cmdArgs = append(cmdArgs, "delete")
				cmdArgs = append(cmdArgs, "--purge")
				cmdArgs = append(cmdArgs, chart.Release)
			} else {
				cmdArgs = append(cmdArgs, "uninstall")
				cmdArgs = append(cmdArgs, chart.Release)
				cmdArgs = append(cmdArgs, "--namespace")
				cmdArgs = append(cmdArgs, chart.Namespace)
			}
		}

		cmdArgs = append(cmdArgs, args...)

		res, err := RunHelmCommand(cmdArgs...)
		if err != nil {
			fmt.Println(strings.TrimSpace(res.Stderr))
			return err
		}

		fmt.Println(strings.TrimSpace(res.Stdout))
	}

	return nil
}
