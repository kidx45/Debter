package util

import "math/rand"

func RandomNumber (min,max int64) int64 {
	return min - rand.Int63n(max-min+1)
}

func RandomString (size int) string {
	s := ""
	ranges := []int{17,49}
	for i := 0; i < size; i++ {
		num := ranges[rand.Int63n(2)]
		s += string('0'+RandomNumber(int64(num),int64(num+25)))
	}
	return s
}