// Copyright © 2020 Anthony Spring <anthonyspring@gmail.com>
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
	"github.com/spf13/viper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChartURL_WithRepo(t *testing.T) {
	viper.SetConfigFile("../testdata/demo.yml")
	viper.ReadInConfig()
	c, _ := LoadAndValidateFromViper()

	assert.Equal(t, c.Charts[0].ChartURL(), "stable/concourse")
}

func TestChartURL_WithoutRepo(t *testing.T) {
	viper.SetConfigFile("../testdata/without-repo.yml")
	viper.ReadInConfig()
	c, _ := LoadAndValidateFromViper()

	assert.Equal(t, c.Charts[0].ChartURL(), "https://github.com/pantsel/konga/blob/master/charts/konga/konga-1.0.0.tgz?raw=true")
}
