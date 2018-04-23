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

package config

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// ChartConfig definition
type ChartConfig struct {
	Name      string                 `mapstructure:"name"`
	Namespace string                 `mapstructure:"namespace"`
	Release   string                 `mapstructure:"release"`
	Repo      string                 `mapstructure:"repo"`
	State     string                 `mapstructure:"state"`
	URL       string                 `mapstructure:"url"`
	Values    map[string]interface{} `mapstructure:"values"`
	Version   string                 `mapstructure:"version"`
}

// ChartLongName returns the name for the chart with version
func (c ChartConfig) ChartLongName() string {
	extra := ""
	if len(c.Version) > 0 {
		extra = "-" + c.Version
	}
	return c.Repo + "/" + c.Name + extra
}

// ChartShortName returns the name for the chart without version
func (c ChartConfig) ChartShortName() string {
	return c.Repo + "/" + c.Name
}

// WriteValueFile writes the given file containing the Chart's Values
func (c ChartConfig) WriteValueFile(file string) error {
	// Marshall the values into a string
	y, err := yaml.Marshal(c.Values)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, y, 0644)
}
