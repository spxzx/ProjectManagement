package common

import "regexp"

func VerifyEmail(email string) bool {
	if email == "" {
		return false
	}
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z].){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func VerifyMobile(mobile string) bool {
	if mobile == "" {
		return false
	}
	regulation := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regulation)
	return reg.MatchString(mobile)
}
