package util

import "math/rand"

func Shuffle[T any](array []T) {
	for i := len(array) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
}

func GetDir(n int) (int, int) {
	var dx int
	var dy int
	if n == 0 {
		dx = 0
		dy = 1
	} else if n == 1 {
		dx = 1
		dy = 1
	} else if n == 2 {
		dx = 1
		dy = 0
	} else if n == 3 {
		dx = 1
		dy = -1
	} else if n == 4 {
		dx = 0
		dy = -1
	} else if n == 5 {
		dx = -1
		dy = -1
	} else if n == 6 {
		dx = -1
		dy = 0
	} else {
		dx = -1
		dy = 1
	}
	return dx, dy
}

func GetRandomDir() (int, int) {
	return GetDir(rand.Intn(8))
}

func GetAdjDir(n int) (int, int) {
	var dx int
	var dy int
	if n == 0 {
		dx = 0
		dy = 1
	} else if n == 1 {
		dx = 1
		dy = 0
	} else if n == 2 {
		dx = 0
		dy = -1
	} else {
		dx = -1
		dy = 0
	}
	return dx, dy
}

func GetRandomAdjDir() (int, int) {
	return GetAdjDir(rand.Intn(4))
}
