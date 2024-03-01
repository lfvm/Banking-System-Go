package utils

import (
	"fmt"
	"math/rand"
	"time"
)


func init(){
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min,max int64) int64{
	return min + rand.Int63() % (max - min + 1)
}

func RandomString(n int) string{
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune,n)
	for i := range s{
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func RandomOwner() string{
	return RandomString(6)
}

func RandomMoney() int64{
	return RandomInt(0,1000)
}

func RandomCurrency() string{
	currencies := []string{USD, MXN, EUR}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string{ 
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
