package vkapi

import (
	"strconv"
	"strings"
)

func sliceToStr(a []int64) string {
	var s []string
	for _, num := range a {
		s = append(s, strconv.FormatInt(num, 10))
	}
	return strings.Join(s, ",")
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
