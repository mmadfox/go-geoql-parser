# go-geoql-parser
GeoQL parser

# Example
```text
TRIGGER
SET
	someplace = multipoint[
		[1.1, 1.1], 
		[-2.1, 2.1]
	]
WHEN
	tracker_point3%2 == 0 
	and tracker_cords intersects @someplace 
	and tracker_point1/tracker_point2*100 > 20% 
	and tracker_week in Sun .. Fri 
	and tracker_time in 9:01AM .. 12:12PM 
	and tracker_temperature in 12Bar .. 44Psi 
	and (
		tracker_speed in 10Kph .. 40Kph 
		or tracker_speed in [10Kph .. 40Kph, 10Kph .. 40Kph, 10Kph .. 40Kph]
	)
REPEAT 5 times 10s interval
RESET after 1h0m0s
```

# Table of contents
- [Operators](#operators)
- [Data Types](#data-types)
  + [Selector](#selector)
  + [Wildcard](#wildcard)
  + [Boolean](#boolean)
  + [Speed](#speed)
  + [Integer](#integer)
  + [Float](#float)
  + [String](#string)
  + [Duration](#duration)
  + [Distance](#distance)
  + [Temperature](#temperature)
  + [Pressure](#pressure)
  + [GeometryPoint](#geometrypoint)
  + [GeometryMultiPoint](#geometrymultipoint)
  + [GeometryLine](#geometryline)
  + [GeometryMultiLine](#geometrymultiline)
  + [GeometryPolygon](#geometrypolygon)
  + [GeometryMultiPolygon](#geometrymultipolygon)
  + [GeometryCircle](#geometrycircle)
  + [GeometryCollection](#geometrycollection)
  + [Date](#date)
  + [Time](#time)
  + [DateTime](#datetime)
  + [Percent](#percent)
  + [Calendar](#calendar)
  + [Array](#array)
  + [Range](#range)
  + [Variable](#variable)

# Operators
| Operator       | Precedence | Literal        |
|----------------|------------|----------------|
| OR             | 1          | or             |
| AND            | 2          | and            |
| LSS            | 3          | <              |
| LEQ            | 3          | <=             |
| GTR            | 3          | >              |
| GEQ            | 3          | >=             |
| EQL            | 3          | ==, eq         |
| NOT EQL        | 3          | !=, not eq     |
| IN             | 3          | in             |
| NOT IN         | 3          | not in         |
| INTERSECTS     | 3          | intersects     |
| NOT INTERSECTS | 3          | not intersects |
| NEARBY         | 3          | nearby         |
| NOT NEARBY     | 3          | not nearby     |
| ADD            | 4          | +              |
| SUB            | 4          | -              |
| MUL            | 5          | *              |
| QUO            | 5          | /              |
| REM            | 5          | %              |

# Data Types
| Data type            | Example                                                                                    |
|----------------------|--------------------------------------------------------------------------------------------|
| Selector             | tracker_speed, coord, some_selector                                                        |
| Wildcard             | *                                                                                          |
| Speed                | 20Kph, 45Mph                                                                               |
| Integer              | 100, 1, -1, 0, 5000                                                                        |
| Float                | -2.300, 5.5, 3000.00                                                                       |
| String               | "some string"                                                                              |
| Duration             | 1h, 20s, 7h3m45s, 7h3m, 3m                                                                 |
| Distance             | 100m, 5Km                                                                                  |
| Temperature          | 19C, 30F                                                                                   |
| Pressure             | 2.2Bar, 4Psi                                                                               |
| GeometryPoint        | point[-1.1, 1.1]                                                                           |
| GeometryMultiPoint   | multipoint[[1.1,1.1],[-2.1, 2.1]]                                                          |
| GeometryLine         | line[[1.1, 1.1], [2.1, 3.1], [3.1, 5.5], [5.5, 5.5]]                                       |
| GeometryMultiLine    | multiline[[[1.1, 1.1], [2.1, 3.1]], [[1.1, 1.1], [2.1, 3.1]], [[1.1, 1.1], [2.1, 3.1]]]    |
| GeometryPolygon      | polygon[[1.1, 1.1], [2.1, 3.1], [3.1, 5.5], [5.5, 5.5], [1.1, 1.1]]                        |
| GeometryMultiPolygon | multipolygon[[[1.1, 1.1], [2.1, 3.1]], [[1.1, 1.1], [2.1, 3.1]], [[1.1, 1.1], [2.1, 3.1]]] |
| GeometryCircle       | circle[-1.1, 1.1]:12km, circle[-1.1, 1.1]:500m                                             |
| GeometryCollection   | collection[point[...], line[...], polygon[...], ...]                                       |
| Date                 | 2030-10-02                                                                                 |
| Time                 | 11:11:11, 11:11, 9:11AM, 3:04Pm                                                            |
| DateTime             | 2030-10-02T11:11:11                                                                        |
| Percent              | 100%                                                                                       |
| Calendar weekday     | Sun, Mon, Tue, Wed, The, Sri, Sat                                                          |
| Calendar month       | Jan, Feb, Mar, Apr, May, Jun, Jul, Aug, Sep, Oct, Nov, Dec                                 |
| Variable             | @somevar                                                                                   |
| Boolean              | true, false                                                                                |
| Array                | [1, 2, 3]                                                                                  |
| Range                | 1 .. 1                                                                                     |


## Selector
## Wildcard
## Boolean
## Speed
## Integer
## Float
## String
## Duration
## Distance
## Temperature
## Pressure
## GeometryPoint
## GeometryMultiPoint
## GeometryLine
## GeometryMultiLine
## GeometryPolygon
## GeometryMultiPolygon
## GeometryCircle
## GeometryCollection
## Date
## Time
## DateTime
## Percent
## Calendar
## Array
## Range
## Variable



