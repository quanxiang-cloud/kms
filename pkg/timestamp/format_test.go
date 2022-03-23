package timestamp

import (
	"fmt"
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	cs := []struct {
		Data   string
		Expect string
	}{
		{
			"yyyyMMddThhmmss",
			"",
		},
		{
			"HH:mm  yyyy-MM-dd",
			"",
		},
		{
			"yyyy-MM-ddTHH:mm:ss.Szz",
			"",
		},
	}

	tm := time.Now()
	for _, c := range cs {
		ft := Format(tm, c.Data)
		fmt.Printf("%s\n", ft)
	}
}
