package tagfs

import (
	"code.google.com/p/go-uuid/uuid"
	"reflect"
	"testing"
)

var u uuid.UUID
var t1 = tagset{"A": u, "B": u, "C": u, "D": u}
var t2 = tagset{"A": u, "C": u, "D": u, "E": u}
var t3 = tagset{"A": u, "D": u, "F": u}

var tagsetTests = []struct {
	name     string
	setFunc  func(...tagset) tagset
	testSets []tagset
	expected tagset
}{
	{"intersection", intersection, []tagset{t1, t2}, tagset{"A": u, "C": u, "D": u}},
	{"intersection", intersection, []tagset{t1, t2, t3}, tagset{"A": u, "D": u}},
	{"union", union, []tagset{t1, t2}, tagset{"A": u, "B": u, "C": u, "D": u, "E": u}},
	{"union", union, []tagset{t1, t2, t3}, tagset{"A": u, "B": u, "C": u, "D": u, "E": u, "F": u}},
	{"difference", difference, []tagset{t1, t2}, tagset{"B": u, "E": u}},
	{"difference", difference, []tagset{t1, t2, t3}, tagset{"A": u, "B": u, "D": u, "E": u, "F": u}},
}

func TestTagset(t *testing.T) {
	for _, test := range tagsetTests {
		result := test.setFunc(test.testSets...)
		if !reflect.DeepEqual(result, test.expected) {
			t.Fatalf("Expected %s to be %v got %v", test.name, test.expected, result)
		}
	}
}
