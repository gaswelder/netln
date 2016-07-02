package main

import (
	"fmt"
	"math"
)

func parseSpeed(str string) (int64, error) {
	if str == "none" { // unlimited
		return -1, nil
	}

	var value float64
	var unit string

	n := len(str)
	i := 0

	if n == 0 {
		return -1, fmt.Errorf("Empty speed string")
	}

	// The first character must be a digit or a dot.
	if !isDigit(str[i]) && str[i] != '.' {
		return -1, fmt.Errorf("Can't parse speed value: %s", str)
	}

	// while have chars and str[i] is a digit
	for i < n && isDigit(str[i]) {
		value = value*10 + float64(str[i]-'0')
		i++
	}

	// If a dot follows, read the fraction.
	if i < n && str[i] == '.' {
		i++
		// A digit must follow the dot.
		if i >= n || !isDigit(str[i]) {
			return -1, fmt.Errorf("Can't parse speed value: %s", str)
		}
		var mul float64
		mul = 10
		// While have chars and str[i] is a digit
		for i < n && isDigit(str[i]) {
			value += float64(str[i]-'0') / mul
			mul *= 10
			i++
		}
	}

	// Rest of the string is a unit specifier.
	for i < n {
		unit += string(str[i])
		i++
	}

	switch unit {
	case "K": // kilobits per second
		return round(value * 1000 / 8), nil
	case "M": // megabits per second
		return round(value * 1e6 / 8), nil
	case "": // bytes per second
		return round(value), nil
	default:
		return -1, fmt.Errorf("Unknown speed unit: %s", unit)
	}
}

func isDigit(ch byte) bool {
	return ch <= '9' && ch >= '0'
}

func round(x float64) int64 {
	return int64(math.Floor(x + 0.5))
}
