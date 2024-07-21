package random

import "math/rand/v2"

func RandomString(size int) string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	res := make([]rune, size)
	for i := 0; i < size; i++ {
		res[i] = rune(letters[rand.IntN(len(letters))])
	}
	return string(res)
}
