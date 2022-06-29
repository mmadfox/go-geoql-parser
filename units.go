package geoqlparser

const (
	Plus  Sign = 1
	Minus Sign = 2
)

const (
	Unknown Unit = iota
	Kph
	Mph
	Celsius
	Fahrenheit
	Kilometer
	Meter
	Bar
	Psi
	Percent
	AM
	PM
)

type (
	Unit uint8
	Sign uint8
)

var unitSizes = map[Unit]Pos{
	Kph:        3,
	Mph:        3,
	Celsius:    1,
	Fahrenheit: 1,
	Kilometer:  2,
	Meter:      1,
	Bar:        3,
	Psi:        3,
	AM:         2,
	PM:         2,
}

func (u Unit) size() Pos {
	p, ok := unitSizes[u]
	if !ok {
		return 0
	}
	return p
}

func (u Unit) String() (s string) {
	switch u {
	default:
		s = "?"
	case Kph:
		s = "Kph"
	case Mph:
		s = "Mph"
	case Celsius:
		s = "C"
	case Fahrenheit:
		s = "F"
	case Kilometer:
		s = "Km"
	case Meter:
		s = "M"
	case Bar:
		s = "Bar"
	case Psi:
		s = "Psi"
	case AM:
		s = "AM"
	case PM:
		s = "PM"
	}
	return
}

func isTimeUnit(s string) (ok bool) {
	switch s {
	case "am", "pm", "Am", "Pm", "PM", "AM":
		ok = true
	}
	return
}

func isDistanceUnit(s string) (ok bool) {
	switch s {
	case "rm", "rkm", "rM", "rKM", "Rm", "Rkm", "Km", "km", "M":
		ok = true
	}
	return
}

func isTemperatureUnit(s string) (ok bool) {
	switch s {
	case "f", "c", "F", "C":
		ok = true
	}
	return
}

func isPressureUnit(s string) (ok bool) {
	switch s {
	case "bar", "Bar", "Psi", "BAR", "PSI", "psi":
		ok = true
	}
	return
}

func isPercentUnit(s string) (ok bool) {
	switch s {
	case "%", "PCT", "Pct":
		ok = true
	}
	return
}

func isSpeedUnit(s string) (ok bool) {
	switch s {
	case "kph", "mph", "KPH", "Kph", "Mph", "MPH":
		ok = true
	}
	return
}

func unitFromString(in string) (out Unit) {
	switch in {
	default:
		out = Unknown
	case "%", "PCT", "Pct":
		out = Percent
	case "rkm", "rKM", "Rkm", "km", "Km":
		out = Kilometer
	case "rm", "rM", "Rm", "M", "met":
		out = Meter
	case "kph", "KPH", "Kph":
		out = Kph
	case "mph", "Mph", "MPH":
		out = Mph
	case "c", "C":
		out = Celsius
	case "f", "F":
		out = Fahrenheit
	case "bar", "Bar", "BAR":
		out = Bar
	case "Psi", "PSI", "si", "psi":
		out = Psi
	case "am", "Am", "AM":
		out = AM
	case "pm", "Pm", "PM":
		out = PM
	}
	return
}

func (v Sign) String() (s string) {
	switch v {
	case Plus:
		s = "+"
	case Minus:
		s = "-"
	}
	return
}
