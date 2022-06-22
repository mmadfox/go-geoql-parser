package geoqlparser

import (
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
	somepoly = polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]]
	somepoint = point[1.1,1.1]
	somecircle = point[1.2, 3.3]:500m
	someline  = line[[1.1,1.1], [1.1,1.1], [1.1,1.1]]
	someline2  = line[[1.1,1.1], [1.1,1.1], [1.1,1.1]]:1km
	somemultiline = multiline[
        line[[1.1,1.1], [1.1,1.1], [1.1,1.1]],
		line[[1.1,1.1], [1.1,1.1], [1.1,1.1]],
		line[[1.1,1.1], [1.1,1.1], [1.1,1.1]],
		line[[1.1,1.1], [1.1,1.1], [1.1,1.1]]
	]
	somemultipoly = multipolygon[
	   	polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]],
		polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]]],
		polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]]
	]
    somemultipoint = multipoint[
point[1.1,1.1],
point[1.1,1.1],
point[1.1,1.1],
point[1.1,1.1]:400m
    ]
	somecoll = collection[
        multipoint[
point[1.1,1.1],
point[1.1,1.1],
point[1.1,1.1],
point[1.1,1.1]:400m
    ],
multipolygon[
	   	polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]],
		polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]]],
		polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]]
	],
multiline[
        line[[1.1,1.1], [1.1,1.1], [1.1,1.1]],
		line[[1.1,1.1], [1.1,1.1], [1.1,1.1]],
		line[[1.1,1.1], [1.1,1.1], [1.1,1.1]],
		line[[1.1,1.1], [1.1,1.1], [1.1,1.1]]
	],
polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]],
point[1.1,1.1],
line[[1.1,1.1], [1.1,1.1], [1.1,1.1]],
line[[1.1,1.1], [1.1,1.1], [1.1,1.1]]:1km
    ]
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
