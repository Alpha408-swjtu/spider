package utils

import (
	"regexp"
	"strconv"
)

func SpliteInfo(info string) (string, string, int) {
	dirRe, _ := regexp.Compile(`导演: (.*)主演`)
	dir := string(dirRe.Find([]byte(info)))

	actRe, _ := regexp.Compile(`主演:(.*)`)
	actor := string(actRe.Find([]byte(info)))

	yearRe, _ := regexp.Compile(`(\d+)`)
	year := string(yearRe.Find([]byte(info)))
	y, _ := strconv.Atoi(year)
	return dir, actor, y
}
