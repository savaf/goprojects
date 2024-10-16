package utils

import (
	"fmt"
	"time"
)

func TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Second:
		return "a few seconds ago"
	case diff < time.Minute:
		return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
	case diff < 2*time.Minute:
		return "a minute ago"
	case diff < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	case diff < 2*time.Hour:
		return "an hour ago"
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	case diff < 48*time.Hour:
		return "yesterday"
	default:
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	}
}
