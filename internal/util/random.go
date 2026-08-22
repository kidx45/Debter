package util

import "math/rand"

func RandomNumber(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(size int64) string {
	s := ""
	ranges := []int{17, 49}
	for i := 0; i < int(size); i++ {
		num := ranges[rand.Int63n(2)]
		randomNum := RandomNumber(int64(num), int64(num+25))
		s += string(rune('0' + randomNum))
	}
	return s
}

func RandomUserName(min, max int) string {
	size := RandomNumber(int64(min), int64(max))
	return RandomString(size)
}

func RandomFullName(min, max int) string {
	size := RandomNumber(int64(min), int64(max))
	return RandomString(size)
}

func RandomEmail(min, max int) string {
	domain := RandomNumber(int64(min), int64(max))
	size := RandomNumber(int64(min), int64(max))
	return RandomString(size) + "@" + RandomString(domain) + ".com"
}

func RandomPassword(min, max int) string {
	size := RandomNumber(int64(min), int64(max))
	return RandomString(size)
}
