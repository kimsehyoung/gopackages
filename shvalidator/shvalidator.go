package shvalidator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// The validator is designed to be thread-safe and used as a singleton instance.
// Refer to auth_service.proto
type AccountValidator struct {
	Email       string `validate:"required,email"`
	Password    string `validate:"required,password"`
	Name        string `validate:"required,alphaunicode,max=20"`
	PhoneNumber string `validate:"required,phonenumber"`
	// number,len=11,startswith=010|startswith=011|startswith=016|startswith=017|startswith=018|startswith=019
}

var (
	prefixPhoneNumber = [...]string{"010", "011", "016", "018", "019"}
)

func IsValidPassword(fl validator.FieldLevel) bool {
	pw := fl.Field().String()

	// repeating same characters <3
	cnt := 0
	var old_v rune
	for _, v := range pw {
		if old_v != v {
			cnt = 0
			old_v = v
			continue
		}
		cnt++
		if cnt >= 2 {
			return false
		}
	}
	// The regexp package uses re2syntax. It doesn't support backreference to keep linear time.
	// match, _ := regexp.MatchString(`(.)\1\1+`, pw)
	// if match {
	// 	fmt.Println("same character")
	// 	return false
	// }

	// 8~20
	fmt.Println(len(pw))
	if len(pw) < 8 || len(pw) > 20 {
		return false
	}

	// alphabet|number|special combination >= 2
	// This includes all characters àb, 😀 ...
	// consider changing to below regex.
	// ~`!@#$%^&*()-_+={}[]|\/:;"'<>,.?₩
	cnt = 0
	if match, _ := regexp.MatchString(`[^A-Za-z0-9]`, pw); match {
		cnt++
	}
	if match, _ := regexp.MatchString(`[A-Za-z]`, pw); match {
		cnt++
	}
	if match, _ := regexp.MatchString(`[0-9]`, pw); match {
		cnt++
	}
	if cnt < 2 {
		return false
	}

	return true
}

// phonenumber's length: 11, remove '-'(hyphen), prefix: 010|011|016|017|018|019
func IsValidPhoneNumber(fl validator.FieldLevel) bool {
	number := strings.Replace(fl.Field().String(), "-", "", -1)

	if len(number) != 11 {
		return false
	}

	_, err := strconv.Atoi(number)
	if err != nil {
		return false
	}

	prefix := number[:3]
	for _, val := range prefixPhoneNumber {
		if val == prefix {
			return true
		}
	}
	return false
}
