package lib

import "testing"

func TestComparePatch(t *testing.T) {
	testCases := []struct {
		Local, Global, Missing HashSet
	}{
		{
			HashSet{
				GetItem("a:b"),
				GetItem("b:c"),
			}, HashSet{
				GetItem("a:b"),
				GetItem("b:c"),
				GetItem("c:d"),
			}, HashSet{
				GetItem("c:d"),
			},
		},
	}

	for _, tc := range testCases {
		result := ComparePatch(tc.Local, tc.Global)
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
