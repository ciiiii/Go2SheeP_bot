package utils

import (
	"github.com/lithammer/shortuuid"
)

func UUID() string {
	id := shortuuid.New()
	return id
}