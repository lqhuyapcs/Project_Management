package utils

import (
	"regexp"
)

var regexMail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var regexPhone = regexp.MustCompile(`^\(?(09|03|07|08|05)+([0-9]{1})\)?[-. ]?([0-9]{3})[-. ]?([0-9]{4})$`)

// CheckValidMail - Check valid input mail
func CheckValidMail(mail string) (bool, string) {
	if mail == "" {
		return false, "empty mail"
	}
	if !regexMail.MatchString(mail) {
		return false, "Invalid mail type"
	}
	return true, "Valid mail type"
}

// CheckValidPhone - Check valid input phone number
func CheckValidPhone(phone string) (bool, string) {
	if phone == "" {
		return false, "Phone is required"
	}
	if !regexPhone.MatchString(phone) {
		return false, "Invalid phone number"
	}
	return true, "Valid phone number"
}
