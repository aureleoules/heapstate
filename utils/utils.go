package utils

import "encoding/json"

// UnescapeString util
func UnescapeString(str string) (unescapedString string) {
	json.Unmarshal([]byte(`"`+str+`"`), &unescapedString)
	return
}

func SanitizeLog(str string) string {
	str = UnescapeString(str)
	if str == "\n" {
		return ""
	}
	return str
}
