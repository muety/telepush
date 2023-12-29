package util

// TruncateInRunes truncates a string to fit the given size in Runes.
// https://github.com/prometheus/alertmanager/blob/dc1466487a537ea7cfce366b69429ac79a04a221/notify/util.go#L89
func TruncateInRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n <= 3 {
		return string(r[:n])
	}
	return string(r[:n-3]) + "..."
}
