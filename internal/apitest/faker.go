package apitest

import (
	"github.com/brianvoe/gofakeit/v7"
	"math/rand"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func FakeString() string {
	return gofakeit.Word()
}

func FakeInt() int {
	return gofakeit.Number(1, 99999)
}

func FakeEmail() string {
	return gofakeit.Email()
}

func randInt(min int, max int) int {
	return seededRand.Intn(max-min+1) + min
}
