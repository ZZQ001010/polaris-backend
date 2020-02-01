package format

import "regexp"

func VerifyEmailFormat(emails ...string) bool {
	reg := regexp.MustCompile(EmailPattern)
	for _, email := range emails{
		suc := reg.MatchString(email)
		if ! suc {
			return false
		}
	}
	return true
}

