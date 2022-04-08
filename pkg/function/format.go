package function

import (
	"kms/internal/enums"
	"regexp"
	"strings"
	"time"
)

// Format the current datetime with given format to a string
func Format(dt time.Time, format string) string {
	format = parseFormat(format)
	return dt.Format(format)
}

func joinFormat() string {
	return strings.Join(FormatSet.GetAll(), "|")
}

func parseFormat(format string) string {
	re := regexp.MustCompile(joinFormat())
	b := re.ReplaceAllFunc([]byte(format), func(bytes []byte) []byte {
		return []byte(formatValueMap[enums.Enum(bytes)])
	})
	return string(b)
}
