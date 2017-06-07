package hamgo

import "strconv"

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func StringToInt(str string, must int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return must
	}
	return value
}

func StringToInt64(str string, must int64) int64 {
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return must
	}
	return value
}

func StringToFloat32(str string, must float32) float32 {
	value, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return must
	}
	return float32(value)
}

func StringToFloat64(str string, must float64) float64 {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return must
	}
	return value
}
