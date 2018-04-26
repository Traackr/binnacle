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

		log.Debugf("Processing chart: %s", chart.ChartLongName())

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
			cmdArgs = append(cmdArgs, chart.ChartShortName())
			cmdArgs = append(cmdArgs, "-i")
			cmdArgs = append(cmdArgs, "--force")
			cmdArgs = append(cmdArgs, "--namespace")
			cmdArgs = append(cmdArgs, chart.Namespace)
			cmdArgs = append(cmdArgs, "--values")
			cmdArgs = append(cmdArgs, valuesFile)
			if len(chart.Version) > 0 {
				cmdArgs = append(cmdArgs, "--version")
				cmdArgs = append(cmdArgs, chart.Version)
			}
			cmdArgs = append(cmdArgs, args...)
		} else {
			cmdArgs = append(cmdArgs, "delete")
			cmdArgs = append(cmdArgs, "--purge")
			cmdArgs = append(cmdArgs, chart.Release)
			cmdArgs = append(cmdArgs, args...)
		}

		res, err := RunHelmCommand(cmdArgs...)
		if err != nil {
			fmt.Println(strings.TrimSpace(res.Stderr))
			return err
		}

		fmt.Println(strings.TrimSpace(res.Stdout))
	}

	return nil
}

func syncRepositories(repos []config.RepositoryConfig, args ...string) error {
	for _, repo := range repos {
		var cmdArgs []string
		var err error
		var res Result
		var currentRepos []config.RepositoryConfig

		log.Debugf("Processing repo: %s", repo.Name)

		currentRepos, err = getCurrentRepositories()
		if err != nil {
			return err
		}

		repoExists, repoFullMatch := repoExists(repo, currentRepos)

		// If the repo exists and is not a full match, or we are deleting the repo - we need to delete the repo
		if repoExists && (!repoFullMatch || repo.State != config.StatePresent) {
			cmdArgs = append(cmdArgs, "repo")
			cmdArgs = append(cmdArgs, "remove")
			cmdArgs = append(cmdArgs, args...)
			cmdArgs = append(cmdArgs, repo.Name)

			res, err = RunHelmCommand(cmdArgs...)
			if err != nil {
				fmt.Println(strings.TrimSpace(res.Stderr))
				return err
			}
		}
		cmdArgs = nil

		if repo.State == config.StatePresent {
			cmdArgs = append(cmdArgs, "repo")
			cmdArgs = append(cmdArgs, "add")
			cmdArgs = append(cmdArgs, args...)
			cmdArgs = append(cmdArgs, repo.Name)
			cmdArgs = append(cmdArgs, repo.URL)

			res, err = RunHelmCommand(cmdArgs...)
			if err != nil {
				fmt.Println(strings.TrimSpace(res.Stderr))
				return err
			}

			fmt.Println(strings.TrimSpace(res.Stdout))
		}
	}

	return nil
}

func getCurrentRepositories() ([]config.RepositoryConfig, error) {
	var err error
	var output []string
	var repos []config.RepositoryConfig
	var res Result

	// Get a list of currently configured repositories
	res, err = RunHelmCommand("repo", "list")
	if err != nil {
		fmt.Println(strings.TrimSpace(res.Stderr))
		return nil, err
	}

	// Split the output on the new line
	output = strings.Split(res.Stdout, "\n")

	// Remove the column titles
	if len(output) > 0 {
		output = output[1:]
	}

	// Populate the repos
	for _, line := range output {
		var repo config.RepositoryConfig

		if len(line) == 0 {
			continue
		}

		// Split the string by a space
		split := strings.Fields(line)

		// Build the repository config
		repo.Name = split[0]
		repo.URL = split[1]

		repos = append(repos, repo)
	}

	return repos, nil
}

func repoExists(repo config.RepositoryConfig, repos []config.RepositoryConfig) (bool, bool) {
	var exists = false
	var fullMatch = false

	// Check if this repo already exists
	for _, checkRepo := range repos {
		// Check if the repos are equal
		if repo == checkRepo {
			exists = true
			fullMatch = true
			break
		} else {
			// Check if their names match but URLs are different
			if repo.Name == checkRepo.Name {
				exists = true
				if repo.URL != checkRepo.URL {
					fullMatch = true
				}
				break
			}
		}

	}

	return exists, fullMatch
}
