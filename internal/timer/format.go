package timer

import (
	"fmt"
	"time"
)

func Format(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	secs := int(d.Round(time.Second).Seconds())
	h := secs / 3600
	m := (secs % 3600) / 60
	s := secs % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
}
