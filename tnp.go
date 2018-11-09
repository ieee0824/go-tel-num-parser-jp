package tnp

import (
	"fmt"
	"regexp"
)

var telPatternRegs = map[string][]*regexp.Regexp{
	"fixed ;ine phone": []*regexp.Regexp{
		regexp.MustCompile(`0[0-9]-[2-9]\d{3}-\d{4}`),
		regexp.MustCompile(`0\d{2}-[2-9]\d{2}-\d{4}`),
		regexp.MustCompile(`0\d{3}-[2-9]\d{1}-\d{4}`),
		regexp.MustCompile(`0\d{4}-[2-9]-\d{4}`),
	},
	"m2m": []*regexp.Regexp{
		regexp.MustCompile(`020-([1-3]|[5-9])\d{2}-\d{5}`),
	},
	"pocket bell": []*regexp.Regexp{
		regexp.MustCompile(`020-4\d{2}-\d{5}`),
	},
	"ip phone": []*regexp.Regexp{
		regexp.MustCompile(`050-[1-9]\d{3}-\d{4}`),
	},
	"mobile phone": []*regexp.Regexp{
		regexp.MustCompile(`0[7-9]0-[1-9]\d{2}-\d{5}`),
		regexp.MustCompile(`0[7-9]0-[1-9]\d{3}-\d{4}`),
	},
	"incoming charge": []*regexp.Regexp{
		regexp.MustCompile(`0120-\d{3}-\d{3}`),
		regexp.MustCompile(`0800-\d{3}-\d{3}`),
	},
	"unified number": []*regexp.Regexp{
		regexp.MustCompile(`0570-\d{3}-\d{3}`),
	},
}

func IsTelNumber(s string) (bool, string) {
	for k, rs := range telPatternRegs {
		for _, r := range rs {
			if r.Copy().MatchString(s) {
				return true, k
			}
		}
	}
	return false, ""
}

func CropTelNumber(s string) (string, error) {
	for _, rs := range telPatternRegs {
		for _, r := range rs {
			if r.Copy().MatchString(s) {
				return r.Copy().FindString(s), nil
			}
		}
	}

	return "", fmt.Errorf("not tel number")
}
