package util

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gitaepark/pha/util/validator"
)

const HANGUEL = "가나다라마바사아자차카타파하"

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// 임의의 휴대폰 번호 생성 함수
func CreateRandomPhoneNumber() string {
	randomNumber := rand.Intn(1e8)
	if randomNumber < 10000000 {
		return fmt.Sprintf("0100%d", randomNumber)
	} else {
		return fmt.Sprintf("010%d", randomNumber)
	}
}

// 임의의 한글 문자열 생성 함수
func CreateRandomString(n int) string {
	start := 0xAC00
	end := 0xD7A3

	result := make([]rune, n)

	for i := 0; i < n; i++ {
		randomCodePoint := start + rand.Intn(end-start+1)
		result[i] = rune(randomCodePoint)
	}

	return string(result)
}

// 임의의 int32 생성 함수
func CreateRandomInt32(min, max int32) int32 {
	return min + rand.Int31n(max-min+1)
}

// 임의의 int64 생성 함수
func CreateRandomInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// 임의의 상품 사이즈 생성 함수
func CreateRandomProductSize() string {
	sizes := []string{validator.Small, validator.Large}
	n := len(sizes)

	return sizes[r.Intn(n)]
}
