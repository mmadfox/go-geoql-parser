package geoqlparser

import (
	"fmt"
	"testing"
)

func TestParseTriggerStmtWhen(t *testing.T) {
	testCases := []struct {
		str string
		err bool
	}{
		{
			str: `
TRIGGER
SET 
	myDevices = index{"127af1cb-ccf3-4b48-874e-4762786d3488", "127af1cb-ccf3-4b48-874e-4762786d3489"}
WHEN 
	(speed between 40kph .. 60kph and somekey between 1 .. 1000) 
or @myDevices in [1,2,3]
or abs not in [40, 60]
repeat 34 times 4h interval 
reset after 48h
`,
		},
	}
	for _, tc := range testCases {
		stmt, err := Parse(tc.str)
		fmt.Println(stmt, err)
		if tc.err {
			if err == nil {
				t.Fatalf("got nil, expected error")
			} else {
				continue
			}
		} else if !tc.err && err != nil {
			t.Fatal(err)
		}
		trigger := stmt.(*TriggerStmt)
		_ = trigger
	}
}
