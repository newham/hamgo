package hamgo

import (
	"fmt"
	"strconv"
)

func int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func intToString(i int) string {
	return strconv.Itoa(i)
}

func stringToInt(str string, must int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return must
	}
	return value
}

func stringToInt64(str string, must int64) int64 {
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return must
	}
	return value
}

func stringToFloat32(str string, must float32) float32 {
	value, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return must
	}
	return float32(value)
}

func stringToFloat64(str string, must float64) float64 {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return must
	}
	return value
}

func printLogo() {
	logo := "     __    __\n" +
		"    /  /  /  / _______    _________    _______        \n" +
		"   /  /__/  / /  ___  \\  /  __ __  \\  /  _____\\    \n" +
		"  /   __   / /  /__/  / /  / // /  / /  / ____  _____ \n" +
		" /  /  /  / /  /  /  / /  / // /  / /  /__/  / / _  / \n" +
		"/__/  /__/ /__/  /__/ /__/ // /__/ /________/ /____/  \n" +
		">>hamgo>> (%s)\n\n"
	fmt.Printf(logo, VERSION)
}
