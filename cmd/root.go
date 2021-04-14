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
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/Traackr/binnacle/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// log The general purpose logging interface available to all commands
var log = logrus.New()

// GITCOMMIT The gitcommit the application was built from
var GITCOMMIT string

// VERSION The version of the application
var VERSION string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "binnacle",
	Short: "An opinionated automation tool for Kubernetes' Helm.",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rootCmdPersistentPreRun(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		rootCmdRun()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// General Flags
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The Binnacle config file (required)")
	RootCmd.MarkFlagRequired("config")

	// Logging Flags
	RootCmd.PersistentFlags().String("loglevel", "info", "The level of logging. Acceptable values: debug, info, warn, error, fatal, panic.")
	viper.BindPFlag("loglevel", RootCmd.PersistentFlags().Lookup("loglevel"))

	// Version Flag
	RootCmd.PersistentFlags().Bool("version", false, "Show the version and exit.")
	viper.BindPFlag("version", RootCmd.PersistentFlags().Lookup("version"))

	// Helm Related Flags
	RootCmd.PersistentFlags().String("helm", "helm", "The path to the Helm binary.")
	viper.BindPFlag("helm", RootCmd.PersistentFlags().Lookup("helm"))
}

func rootCmdPersistentPreRun(cmd *cobra.Command) {

	// Handle the special case of the version
	if cmd.Name() == "binnacle" {
		if viper.IsSet("version") && viper.GetBool("version") {
			fmt.Printf("%s-%s\n", VERSION, GITCOMMIT)
			os.Exit(0)
		}
	}

	if cfgFile == "" {
		fmt.Println(cmd.UsageString())
		os.Exit(-1)
	}

	// Verify the file exists
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		fmt.Println("config file does not exist!")
		os.Exit(-1)
	}

	viper.SetConfigFile(cfgFile)
	viper.AddConfigPath(".") // check current dir
	viper.AutomaticEnv()     // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Loading config file:", viper.ConfigFileUsed())
		fmt.Printf("  %s.  Exiting!\n", err)
		os.Exit(-1)
	}
	fmt.Println("Loading config file:", viper.ConfigFileUsed())

	// Initialize the logger for all commands to use
	logLevel, _ := logrus.ParseLevel(viper.GetString("loglevel"))
	log.Level = logLevel
	log.Debug("Logger initialized.")
}

func rootCmdRun() {
	// This is here as a no-op to allow `binnacle --version` to work correctly
}

// IsHelm2 returns true if the provided helm is version 2
func IsHelm2() bool {
	var err error
	var helm2 bool
	var res Result

	// helm version 2 - Client: &version.Version{SemVer:"v2.16.0", GitCommit:"e13bc94621d4ef666270cfbe734aaabf342a49bb", GitTreeState:"clean"}
	// helm version 3 - version.BuildInfo{Version:"v3.0.2", GitCommit:"19e47ee3283ae98139d98460de796c1be1e3975f", GitTreeState:"clean", GoVersion:"go1.13.5"}

	// Get a list of currently installed plugins
	res, err = RunHelmCommand("version")
	if err != nil {
		fmt.Println(strings.TrimSpace(res.Stderr))
		os.Exit(-1)
	}

	// Check if the output contains the helm2 identifying string
	helm2 = strings.Contains(res.Stdout, "Client: &version.Version{SemVer:\"v2")
	if helm2 {
		log.Debugf("Helm 2 detected!")
	}

	return helm2
}

// PluginInstalled returns if the given plugin is installed
func PluginInstalled(plugin string) (bool, error) {
	var err error
	var output []string
	var res Result

	// Get a list of currently installed plugins
	res, err = RunHelmCommand("plugin", "list")
	if err != nil {
		fmt.Println(strings.TrimSpace(res.Stderr))
		return false, err
	}

	// Split the output on the new line
	output = strings.Split(res.Stdout, "\n")

	// Remove the column titles
	if len(output) > 0 {
		output = output[1:]
	}

	// Iterate the plugins
	for _, line := range output {
		if len(line) == 0 {
			continue
		}

		// Split the string by a space
		split := strings.Fields(line)

		if plugin == split[0] {
			return true, nil
		}
	}

	return false, nil
}

func ReleaseExists(namespace string, release string) bool {
	var exists = true
	var err error
	var res Result

	// Get the status of the release for the namespace
	res, err = RunHelmCommand("status", release, "--namespace", namespace)
	if err != nil {
		if res.Stderr != "Error: release: not found" {
			exists = false
		}
	}

	return exists
}

// RunHelmCommand runs the given command against helm
func RunHelmCommand(args ...string) (Result, error) {
	var result Result
	var outbuf, errbuf bytes.Buffer

	cmd := exec.Command(viper.GetString("helm"), args...)

	log.Debugf("Executing command:  %v", cmd.Args)

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	var err = cmd.Run()

	log.Debugf("Execution complete.")

	result.Stdout = strings.Trim(outbuf.String(), " ")
	result.Stderr = strings.Trim(errbuf.String(), " ")

	if err != nil {
		//
		// This crafty snippet is from https://stackoverflow.com/questions/10385551/get-exit-code-go
		//
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.Sys().(syscall.WaitStatus).ExitStatus()
		} else {
			if result.Stderr == "" {
				result.Stderr = err.Error()
			}
		}
	}

	return result, err
}

func syncRepositories(repos []config.RepositoryConfig, args ...string) error {
	var reposModified = false

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
			reposModified = true

			fmt.Println(strings.TrimSpace(res.Stdout))
		}
	}

	// If any repos have been added during this sync execute a helm repos update to update the cache.
	if reposModified {
		var cmdArgs []string
		var err error
		var res Result

		cmdArgs = append(cmdArgs, "repo")
		cmdArgs = append(cmdArgs, "update")
		res, err = RunHelmCommand(cmdArgs...)
		if err != nil {
			fmt.Println(strings.TrimSpace(res.Stderr))
			return err
		}
		fmt.Println(strings.TrimSpace(res.Stdout))
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
		if repo.Equal(checkRepo) {
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
