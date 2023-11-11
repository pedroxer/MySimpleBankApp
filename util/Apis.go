package util

import rand2 "math/rand"

func FinReport() float64 {
	rand := rand2.Float64()
	return rand
}

func SecurityCheck() bool {
	v := RandomInt(0, 1)
	if v == 0 {
		return false
	}
	return true
}
