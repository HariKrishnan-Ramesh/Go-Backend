package common

import (
	"math/rand"
	"time"
)

func GenerateOTP(length int) string {
	rand.Seed(time.Now().UnixNano())
	const chars = "0123456789"
	result := make([]byte,length)
	for i:=0; i<length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}