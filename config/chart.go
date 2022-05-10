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
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

// ChartConfig definition
type ChartConfig struct {
	Kustomize BinnacleKustomization `mapstructure:"kustomize"`
	Name      string                `mapstructure:"name"`
	Namespace string                `mapstructure:"namespace"`
	Release   string                `mapstructure:"release"`
	Repo      string                `mapstructure:"repo"`
	State     string                `mapstructure:"state"`
	URL       string                `mapstructure:"url"`
	Values    map[string]any        `mapstructure:"values"`
	Version   string                `mapstructure:"version"`
}

// Adapted from https://github.com/kubernetes-sigs/kustomize/blob/master/api/types/kustomization.go
type BinnacleKustomization struct {
	// https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/resource/
	Resources []string `mapstructure:"resources,omitempty"`

	// https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/patches/
	Patches []Patch `mapstructure:"patches,omitempty"`
}

type Patch struct {
	Path    string          `mapstructure:"path,omitempty"`
	Patch   string          `mapstructure:"patch,omitempty"`
	Target  *Selector       `mapstructure:"target,omitempty"`
	Options map[string]bool `mapstructure:"options,omitempty"`
}

type Selector struct {
	AnnotationSelector string `mapstructure:"annotationSelector,omitempty"`
	LabelSelector      string `mapstructure:"labelSelector,omitempty"`
}

func (k BinnacleKustomization) Empty() bool {
	return len(k.Resources) == 0 && len(k.Patches) == 0
}

// ChartURL returns a URL related to the given repo and name of the chart based off of
// criteria 1 through 4 of the following documentation on how to specify local and remote charts
//
// 1. By chart reference: helm install mymaria example/mariadb
// 2. By path to a packaged chart: helm install mynginx ./nginx-1.2.3.tgz
// 3. By path to an unpacked chart directory: helm install mynginx ./nginx
// 4. By absolute URL: helm install mynginx https://example.com/charts/nginx-1.2.3.tgz
//
func (c ChartConfig) ChartURL() string {
	// If a repository is given return the c
	if len(c.Repo) > 0 {
		return c.Repo + "/" + c.Name
	}
	return c.Name
}

// WriteValueFile writes the given file containing the Chart's Values
func (c ChartConfig) WriteValueFile(dir string) (string, error) {
	// Marshall the values into a string
	y, err := yaml.Marshal(c.Values)
	if err != nil {
		return "", fmt.Errorf("marshalling chart values: %w", err)
	}

	valuesYml := filepath.Join(dir, "values.yml")
	err = os.WriteFile(valuesYml, y, 0644)
	if err != nil {
		return "", fmt.Errorf("writing temporary values.yml file: %w", err)
	}

	return valuesYml, nil
}
