package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("23456789ABCDEFGHJKLMNPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var nameChars = []rune("ABCDEFGHIJKLMNPQRSTUVWXYZ23456789")

// 私有生成随机名方法
func GenerateRandomName(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, length)
	for i := range b {
		b[i] = nameChars[r.Intn(len(nameChars))]
	}
	return string(b)
}
