package shvalidator

import (
	"fmt"
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

// TODO
func IsValidPassword(fl validator.FieldLevel) bool {
	pw := fl.Field().String()
	fmt.Println(len(pw))
	if len(pw) < 8 || len(pw) > 20 {
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
