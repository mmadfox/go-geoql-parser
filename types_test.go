package geoqlparser

import "testing"

func TestTypes(t *testing.T) {
	for _, tok := range keywords {
		if !isSelector(tok) {
			continue
		}
		typ, found := kwdTypList[tok]
		if !found {
			t.Fatalf("type for selector %s not found", KeywordString(tok))
		}
		if typ == 0 {
			t.Fatalf("type is not defined for slector %s", KeywordString(tok))
		}
	}
}
