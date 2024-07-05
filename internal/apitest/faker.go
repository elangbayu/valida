package apitest

import (
	"github.com/brianvoe/gofakeit/v7"
)

func fakeString() string {
	return gofakeit.Word();
}

func fakeEmail() string {
	return gofakeit.Email();
}

func fakeInt() int {
	return gofakeit.Number(0, 100);
}
