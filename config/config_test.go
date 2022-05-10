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
	"testing"

	"github.com/spf13/viper"
)

func TestBooleanIsNotCoerced(t *testing.T) {
	viper.SetConfigFile("../testdata/demo.yml")
	viper.ReadInConfig()
	c, _ := LoadAndValidateFromViper()

	ingressConfig := c.Charts[0].Values["ingress"].(map[string]interface{})
	want := true
	got := ingressConfig["enabled"]
	if want != got {
		t.Errorf("want `ingress.enabled` to be type=%T value=%v, but got type=%T value=%v", want, want, got, got)
	}
}

func TestLoadAndValidateFromViper_Unmarshallable(t *testing.T) {
	viper.SetConfigFile("../testdata/unmarshallable.yml")
	viper.ReadInConfig()

	_, err := LoadAndValidateFromViper()
	if err == nil {
		t.Errorf("want an error for unmarshallable data, but was nil")
	}
}

func TestLoadAndValidateFromViper_DefaultChartState(t *testing.T) {
	viper.SetConfigFile("../testdata/default-state.yml")
	viper.ReadInConfig()

	c, _ := LoadAndValidateFromViper()
	got := c.Charts[0].State
	want := "present"
	if got != want {
		t.Errorf("want state to be %s, but got %s", want, got)
	}
}

func TestLoadAndValidateFromViper_DefaultRepoState(t *testing.T) {
	viper.SetConfigFile("../testdata/default-state.yml")
	viper.ReadInConfig()

	c, _ := LoadAndValidateFromViper()
	got := c.Repositories[0].State
	want := "present"
	if got != want {
		t.Errorf("want state to be %s, but got %s", want, got)
	}
}
