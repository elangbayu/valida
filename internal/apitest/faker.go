package apitest

import "github.com/brianvoe/gofakeit/v7"

func FakeString() string {
	return gofakeit.Word()
}

func FakeInt() int {
	return gofakeit.Number(1, 99999)
}

func FakeEmail() string {
	return gofakeit.Email()
}
