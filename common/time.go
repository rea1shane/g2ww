package common

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const (
	TimeLayout = "2006-01-02 Mon 15:04:05"
)

func FormatDuration(d time.Duration) string {
	var symbol, day, hour, min, sec string

	dString := d.String()
	dRegexp := regexp.MustCompile(`(-)?(\d*h)?(\d*m)?(\d*\.?\d*s)?`)
	params := dRegexp.FindStringSubmatch(dString)

	symbol = params[1]
	min = params[3]

	h := params[2]
	s := params[4]

	hRegexp := regexp.MustCompile(`(\d*)h?`)
	hParams := hRegexp.FindStringSubmatch(h)
	hCount, _ := strconv.Atoi(hParams[1])
	if hCount >= 24 {
		day = strconv.Itoa(hCount/24) + "d"
		hour = strconv.Itoa(hCount%24) + "h"
	} else {
		hour = h
	}

	sRegexp := regexp.MustCompile(`(\d*)\.?\d*s?`)
	sParams := sRegexp.FindStringSubmatch(s)
	sCount, _ := strconv.Atoi(sParams[1])
	sec = strconv.Itoa(sCount) + "s"

	if day != "" {
		day += " "
	}
	if hour != "" {
		hour += " "
	}
	if min != "" {
		min += " "
	}

	return fmt.Sprintf("%s%s%s%s%s", symbol, day, hour, min, sec)
}
