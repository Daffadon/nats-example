package utils

func CheckCondition(temp float32) string {
	if temp < 20.0 {
		return "cool"
	} else if temp >= 20.0 && temp <= 30.0 {
		return "normal"
	} else {
		return "overheat"
	}
}
