package readfile

import (
	"regexp"
)

func MatchFileName(filename string) bool {
	match, err := regexp.MatchString("[a-zA-z0-9_-]+.bem$", filename)
	if err != nil {
		println("Girdiğiniz dosya uzantısı veya adı geçerli değil", err)
		return false
	}

	if match == true {
		println("başarılıyla compile edildi")
		return match
	}

	return match
}
