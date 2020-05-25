package vkapi

import (
	"strconv"
	"strings"
)

func sliceToStr(a []int) string {
	var s []string
	for _, num := range a {
		s = append(s, strconv.Itoa(num))
	}
	return strings.Join(s, ",")
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
