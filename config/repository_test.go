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
)

func TestRepositoryEquals_NamesDoNotMatch(t *testing.T) {
	rep1 := RepositoryConfig{Name: "foo1", URL: ""}
	rep2 := RepositoryConfig{Name: "foo2", URL: ""}

	if rep1.Equal(rep2) {
		t.Errorf("want %#v to NOT equal %#v", rep1, rep2)
	}
}

func TestRepositoryEquals_URLsDoNotMatch(t *testing.T) {
	rep1 := RepositoryConfig{Name: "foo", URL: "url1"}
	rep2 := RepositoryConfig{Name: "foo", URL: "url2"}

	if rep1.Equal(rep2) {
		t.Errorf("want %#v to NOT equal %#v", rep1, rep2)
	}
}

func TestRepositoryEquals_IgnoresState(t *testing.T) {
	rep1 := RepositoryConfig{Name: "foo", URL: "url", State: "present"}
	rep2 := RepositoryConfig{Name: "foo", URL: "url", State: "absent"}

	if !rep1.Equal(rep2) {
		t.Errorf("want %#v to equal %#v", rep1, rep2)
	}
}

func TestRepositoryEquals(t *testing.T) {
	rep1 := RepositoryConfig{Name: "foo", URL: "url"}
	rep2 := RepositoryConfig{Name: "foo", URL: "url"}

	if !rep1.Equal(rep2) {
		t.Errorf("want %#v to equal %#v", rep1, rep2)
	}
}
