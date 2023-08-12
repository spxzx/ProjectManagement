package test

import (
	"github.com/spxzx/project-common/jwts"
	"testing"
)

func TestParseToken(t *testing.T) {
	tokenString1 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzQyMTU4MDQsInRva2VuIjoiNCJ9.YMc9hXekNPXnTMBUD1Yawm64jDSlyJfqJ77iPbfYlfo"
	tokenString2 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzQ4MjA2MDQsInRva2VuIjoiNCJ9.DsBFztXswS7gynqwAeTgZd7RY4OZiVsbrmZdF4JXUpc"
	jwts.ParseToken(tokenString1, "g2r3fa")
	jwts.ParseToken(tokenString2, "13f3ah")
}
