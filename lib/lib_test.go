// Copyright 2016 Lennart Espe. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package lib

import "testing"

func TestCompare(t *testing.T) {
	SetVerbose(false)
	testCases := []struct {
		Local, Global, Missing HashSet
	}{
		{
			HashSet{
				HashItem{"a", "b"},
				HashItem{"b", "c"},
			}, HashSet{
				HashItem{"a", "b"},
				HashItem{"b", "c"},
				HashItem{"c", "d"},
			}, HashSet{
				HashItem{"c", "d"},
			},
		},
	}

	for _, tc := range testCases {
		result := Compare(tc.Local, tc.Global)
		assertMap := make(map[string]bool)
		for _, e := range tc.Missing {
			assertMap[e.String()] = true
		}
		for _, e := range result {
			if _, ok := assertMap[e.String()]; ok {
				assertMap[e.String()] = false
			}
		}
		for _, v := range assertMap {
			if v {
				t.Fatal("Missing difference:", tc.Local, tc.Global, "->", tc.Missing, "is not equal to", result)
			}
		}
	}
}
