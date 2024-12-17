package security

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(999999-100000) + 100000
	return strconv.Itoa(otp)
}
