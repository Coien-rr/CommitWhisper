package errors

import "regexp"

func removeRequestID(msg string) string {
	re := regexp.MustCompile(`(requestId|Request id):\s*.+`)
	resMsg := re.ReplaceAllString(msg, "")
	return resMsg
}
