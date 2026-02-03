package utils

func CompareFloat4(a, b float64) bool {
	precision := float64(10000)
	return int64(a*precision+0.5) == int64(b*precision+0.5)
}
