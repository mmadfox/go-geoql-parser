package geoqlparser

import (
	"bytes"
	"testing"
)

func TestFormat(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	stmt, err := Parse(`
		trigger 
		set 
		t4=time[11:11 .. 11:11];
		a=1;
		b=5345345345;
		w1=weekday[Sun .. Sat];
		w2=weekday[Sun, Sat];
		w3=weekday[Sun];
		m1=month[Jan .. Jul];
		m2=month[Jan, Jul];
		m3=month[Jan];
		d1=date[2022-11-11 .. 2022-12-12];
		d2=date[2022-11-11, 2022-12-12];
		d3=date[2022-11-11];
		t1=time[11:11];
		t2=time[9:00AM];
		t3=time[11:11:11];
		t5=time[11:11, 12:01, 13:10];
  		some=345345345345;
		floatval=22.22;
		durationval=7h40m;
		temp1=+40C;
		temp0=0C
		temp2=-30C;
		temp3=0C;
		pointval=point[-1.1, 1.1]:1km;
		lineval3 = line[[2.1, 3.1], [3.1, 5.5], [5.5, 5.5], [5.5, 5.5]]:44M;
		lineval = line[[1.1, 1.1], [2.1, 3.1], [3.1, 5.5], [5.5, 5.5], [5.5, 5.5]];
		linevall = line[
			[1.1, 1.1], [2.1, 3.1], [3.1, 5.5], [5.5, 5.5],  
			[5.5, 5.5], [5.5, 5.5], [5.5, 5.5], [5.5, 5.5],
			[5.5, 5.5], [5.5, 5.5], [5.5, 5.5], [5.5, 5.5]
		]:44km;
		polygonval2 = polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]]]; 
		polygonvla = polygon[
			[
				[1.1, 1.1], [1.1, 1.1], [1.1, 1.1], [1.1, 1.1], 
				[1.1, 1.1], [1.1,1.1]
			],
			[
   				[1.1,1.1], [1.1,1.1]	
			]
		];
		multipolygon1 = multipolygon[
			polygon[
			[
				[1.1, 1.1], [1.1, 1.1], [1.1, 1.1], [1.1, 1.1], 
				[1.1, 1.1], [1.1,1.1]
			],
			[
   				[1.1,1.1], [1.1,1.1]	
			]],
		    polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]]]	
		];
		collection1 = collection[
			polygon[
			[
				[1.1, 1.1], [1.1, 1.1], [1.1, 1.1], [1.1, 1.1], 
				[1.1, 1.1], [1.1,1.1]
			],
			[
   				[1.1,1.1], [1.1,1.1]	
			]],
			line[
			[1.1, 1.1], [2.1, 3.1], [3.1, 5.5], [5.5, 5.5],  
			[5.5, 5.5], [5.5, 5.5], [5.5, 5.5], [5.5, 5.5],
			[5.5, 5.5], [5.5, 5.5], [5.5, 5.5], [5.5, 5.5]
			]:44km
		]
   
		when 1*1 == 2
repeat 1 every 10s
reset after 34h
`)
	if err != nil {
		t.Fatal(err)
	}
	if err := Format(buf, stmt); err != nil {
		t.Fatal(err)
	}
	t.Log(buf.String())
	_, err = Parse(buf.String())
	if err != nil {
		t.Fatal(err)
	}
}
