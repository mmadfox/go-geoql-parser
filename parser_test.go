package geoqlparser

import (
	"log"
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
    someplace = multipoint[[1.1,1.1],[-2.1, 2.1]]  
WHEN
	tracker_point3 % 2 == 0
	and tracker_cords intersects @someplace
	and tracker_point1 / tracker_point2 * 100 > 20%
	and tracker_week in Sun .. Fri
	and tracker_time in 9:01AM .. 12:12PM 
	and tracker_temperature in 12Bar .. 44Psi
	and (tracker_speed in 10kph .. 40kph
	or tracker_speed in [10kph .. 40kph, 10kph .. 40kph, 10kph .. 40kph])
repeat 5 times 10s interval 
reset after 1h 
`,
		},
	}
	// TODO:
	for _, tc := range testCases {
		stmt, err := Parse(tc.str)
		log.Println(stmt)
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
