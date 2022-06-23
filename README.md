# go-geoql-parser
Package declares the types used to represent syntax trees for GeoQL rules.

# Trigger example
```text
TRIGGER
SET
	someplace = multipoint[
		point[1.1, 1.1]:400m, 
		point[-2.1, 2.1]:5Km
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
| Data type            | Example                                                                                         |
|----------------------|-------------------------------------------------------------------------------------------------|
| Selector             | tracker_speed, coord, some_selector                                                             |
| Wildcard             | *                                                                                               |
| Speed                | 20Kph, 45Mph                                                                                    |
| Integer              | 100, 1, -1, 0, 5000                                                                             |
| Float                | -2.300, 5.5, 3000.00                                                                            |
| String               | "some string"                                                                                   |
| Duration             | 1h, 20s, 7h3m45s, 7h3m, 3m                                                                      |
| Distance             | 100M, 5Km                                                                                       |
| Temperature          | 19C, 30F                                                                                        |
| Pressure             | 2.2Bar, 4Psi                                                                                    |
| GeometryPoint        | point[-1.1, 1.1]                                                                                |
| GeometryMultiPoint   | multipoint[point[-1.1, 1.1], point[-1.1, 1.1]]                                                  |
| GeometryLine         | line[[1.1, 1.1], [2.1, 3.1], [3.1, 5.5], [5.5, 5.5]]                                            |
| GeometryMultiLine    | multiline[line[[1.1, 1.1], [2.1, 3.1]], [[1.1, 1.1], [2.1, 3.1]], line[[1.1, 1.1], [2.1, 3.1]]] |
| GeometryPolygon      | polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]]                   |
| GeometryMultiPolygon | multipolygon[polygon[...], polygon[...]]                                                        |
| GeometryCircle       | point[-1.1, 1.1]:12km, point[-1.1, 1.1]:500m                                                    |
| GeometryCollection   | collection[point[...], line[...], polygon[...], ...]                                            |
| Date                 | 2030-10-02                                                                                      |
| Time                 | 11:11:11, 11:11, 9:11AM, 3:04Pm                                                                 |
| DateTime             | 2030-10-02T11:11:11                                                                             |
| Percent              | 100%                                                                                            |
| Calendar weekday     | Sun, Mon, Tue, Wed, Thu, Fri, Sat                                                               |
| Calendar month       | Jan, Feb, Mar, Apr, May, Jun, Jul, Aug, Sep, Oct, Nov, Dec                                      |
| Variable             | @somevar                                                                                        |
| Boolean              | true, false                                                                                     |
| Array                | [1, 2, 3]                                                                                       |
| Range                | 1 .. 1                                                                                          |


## Selector
This data type is used to match values from INPUT data.

```go
 *geoqlparser.SelectorExpr
```

- valid characters: *a-zA-Z_0-9*
- min length: *1*
- max length: *64*

### Get some value from current device
wildcard - indicates the current device
```text
tracker_speed{*} in 1mph .. 40mph 
tracker_index1{*} in 1 .. 10
someTextToo{*} in 1mph .. 40mph 
```
the same, but shorter
```text
tracker_speed in 1mph .. 40mph 
tracker_index1 in 1 .. 10
someTextToo in 1mph .. 40mph 
```
### Get some value from other devices by their IDs
```text
tracker_speed{
    "786d9e27-f277-4d5d-b658-3198c133c43d", 
    "786d9e27-f277-4d5d-b658-3198c133c43b"
} in 1mph .. 40mph
```
### Get some values from current device and from other devices by their IDs
```text
tracker_speed{*, "786d9e27-f277-4d5d-b658-3198c133c43d", "786d9e27-f277-4d5d-b658-3198c133c43b"} in 1mph .. 40mph
```

### Properties
##### selector :prop1,...,propN
Example:
```text
tracker_coords:1km // radius 1km
coords{*}:1km,2km  // radius 1km or 2km
someIndex:1,2,3    // some index 
ads:2,6            // some index too
selector{"id1", "id2", *}:1km,6km,12km
h3Index:1,2        // same as calculating index h3 with levels 1 and 2 
```

## Wildcard
This data type is used to indicate the current device
```go
*geoqlparser.WildcardLit
```
Example:
```text
selector{*}
someSelector{*, "786d9e27-f277-4d5d-b658-3198c133c43d"}
```
## Boolean
This data type is used to describe the boolean value TRUE or FALSE
```go
*geoqlparser.BooleanLit
```
Example:
```text
tracker_status eq true
tracker_status not eq false 
tracker_abs != true
tracker_status == false 
```
## Speed
This data type is used to describe the SPEED value 
```go
*geoqlparser.SpeedLit
```
Unit of measurement:
- Kilometer/hour: Kph
- Mile/hour: Mph

Example:
```text
tracker_speed > 1.1Kph
tracker_speed in 0kph .. 120kph
// if the speed falls into one or the second range
tracker_speed in [1kph .. 20kph, 40kph .. 80kph]
```
## Integer
This data type is used to describe the INTEGER value. 
The size of the generic int type is platform dependent. It is 32 bits wide on a 32-bit system and 64-bits wide on a 64-bit system.
```go
*geoqlparser.IntLit
```
Example:
```go
some_selector > 100
1+2 == 3
h3_idx_selector{*, "786d9e27-f277-4d5d-b658-3198c133c43d"}:1,2 in [1, 2]
```

## Float
This data type is used to describe the FLOAT value
```go
*geoqlparser.FloatLit
```
Example:
```text
tracker_pid_value >= 33.3455
tracker_coords intersects point[-74.232423423, 54.455644]
```

## String
This data type is used to describe the STRING value
```go
*geoqlparser.StringLit
```
Example:
```text
tracker_status == "Y" or tracker_status == "Valid"
tracker_model eq "ER54x3"
```

## Duration
This data type is used to describe the DURATION value
```go
*geoqlparser.DurationLit
```

1h, 20s, 7h3m45s, 7h3m, 3m

- hours: 1h, 2h, 24h
- minutes: 7m, 5m, 1m
- seconds: 1s, 10s, 45s

Example:
```text
tracker_some_select in 1h .. 2h
tracker_some_select in [1s .. 30s, 40s .. 60s]
tracker_some_select > 7h3m
tracker_some_select < 7h3m45s 
```

## Distance
This data type is used to describe the DISTANCE value
```go
*geoqlparser.DistanceLit
```

Unit of measurement:
 - Kilometers: Km
 - Meters: M 

Example:
```text
tracker_radius > 40Km
tracker_radius in [1M, 100M, 1000M]
tracker_radius in 40Km .. 45Km
```

## Temperature
This data type is used to describe the TEMPERATURE value
```go
*geoqlparser.TemperatureLit
```
Unit of measurement:
- Celsius: C
- Fahrenheit: F

Example:
```text
tracker_temperature in +34C .. -15C
tracker_temperature >= 0C and tracker_temperature < -5C
```

## Pressure
This data type is used to describe the PRESSURE value
```go
*geoqlparser.PressureLit
```
Unit of measurement:
- Psi
- Bar

Psi and Bar are units of measurement of pressure. The key difference between psi and bar is that psi measures the pressure as the one-pound force applied on an area of one square inch whereas bar measures the pressure as a force applied perpendicularly on a unit area of a surface.

Example:
```text
tracker_some_val in 1.1Bar .. 20Bar
tracker_some_val < 40Psi
```

## GeometryPoint
This data type is used to describe the GEOMETRY POINT value is a single position

Each value represents a float type
```go
*geoqlparser.GeometryPointExpr

type GeometryPointExpr struct {
    Val      [2]float64
    Radius   *DistanceLit
    // ...	
}
```

Example:

tracker_coords - selector for longitude and latitude current device
```text
tracker_coords intersects point[114.60937499999999, 69.90011762668541]
```
with radius 500 meters or 50 kilometers 
```text
tracker_coords intersects point[114.60937499999999, 69.90011762668541]:500m
tracker_coords intersects point[114.60937499999999, 69.90011762668541]:50Km

```

## GeometryMultiPoint
This data type is used to describe the GEOMETRY MULTI POINT values is an array of positions

Each value represents a float type

```go
*geoqlparser.GeometryMultiObject
```

Example:

tracker_coords - selector for longitude and latitude current device
```text
tracker_coords intersects multipoint[
    point[114.60937499999999,69.90011762668541], 
    point[124.1015625, 68.39918004344189], 
    point[ 113.5546875, 66.65297740055279]]
```

## GeometryLine
This data type is used to describe the GEOMETRY LINE values is an array of two or
more positions

Each value represents a float type
```go
*geoqlparser.GeometryLineExpr

type GeometryLineExpr struct {
    Val      [][2]float64
	Margin   *DistanceLit
}
```
Example:

tracker_coords - selector for longitude and latitude current device

```text
tracker_coords intersects line[[38.3203125,60.413852350464914], [69.60937499999999,51.6180165487737]]
tracker_coords intersects line[[38.3203125,60.413852350464914], [69.60937499999999,51.6180165487737], [83.3203125,30.751277776257812]]
```
with margin 400 meters
```go
tracker_coords intersects line[[38.3203125,60.413852350464914], [69.60937499999999,51.6180165487737]]:400m
```

## GeometryMultiLine
This data type is used to describe the GEOMETRY MULTI LINE values is an array of
line coordinate arrays

Each value represents a float type

```go
*geoqlparser.GeometryMultiObject
```

Example:

tracker_coords - selector for longitude and latitude current device

```text
tracker_coords intersects multiline[
    line[[38.3203125,60.413852350464914], [69.60937499999999,51.6180165487737]],
    line[[38.3203125,60.413852350464914], [69.60937499999999,51.6180165487737]]
]
```

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



