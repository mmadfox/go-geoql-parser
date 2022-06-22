# go-geoql-parser
Package declares the types used to represent syntax trees for GeoQL rules.

# Trigger example
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
This data type is used to match values from INPUT data.

```go
 *geoqlparser.SelectorExpr

type SelectorExpr struct {
    Ident    string              // selector name
    Args     map[string]struct{} // device ids
    Wildcard bool                // indicates the current device
    Props    []Expr              // some props
}
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
### Get some value from devices by their IDs
```text
tracker_speed{"786d9e27-f277-4d5d-b658-3198c133c43d", "786d9e27-f277-4d5d-b658-3198c133c43b"} in 1mph .. 40mph
```
### Get some values from current device and from devices by their IDs
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
This data type is used to describe the speed value with units Kph or Mph
```go
*geoqlparser.SpeedLit

type SpeedLit struct {
    Val float64
    U   Unit    // kph, mph
}
```
Example:
```text
tracker_speed > 1.1Kph
tracker_speed in 0kph .. 120kph
// if the speed falls into one or the second range
tracker_speed in [1kph .. 20kph, 40kph .. 80kph]
```
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



