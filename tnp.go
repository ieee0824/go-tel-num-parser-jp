package tnp

import (
	"fmt"
	"regexp"
	"strings"
)

var ignoreTypes = []string{}

func SetIgnoreTypes(s ...string) {
	ignoreTypes = append(ignoreTypes, s...)
}

func isIgnore(s string) bool {
	for _, v := range ignoreTypes {
		if v == s {
			return true
		}
	}
	return false
}

var telPatternRegs = map[string][]*regexp.Regexp{
	"fixed line phone": []*regexp.Regexp{
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
		if isIgnore(k) {
			continue
		}
		for _, r := range rs {
			if r.Copy().MatchString(s) {
				return true, k
			}
		}
	}

	if hasParenthesis(s) {
		return IsTelNumber(replaceParenthesis(s))
	}

	return false, ""
}

func CropTelNumber(s string) (string, error) {
	for k, rs := range telPatternRegs {
		if isIgnore(k) {
			continue
		}
		for _, r := range rs {
			if r.Copy().MatchString(s) {
				return r.Copy().FindString(s), nil
			}
		}
	}

	if hasParenthesis(s) {
		return CropTelNumber(replaceParenthesis(s))
	}

	return "", fmt.Errorf("not tel number")
}

func hasParenthesis(s string) bool {
	return strings.Contains(s, "(") && strings.Contains(s, ")")
}

func replaceParenthesis(s string) string {
	var ret string
	ret = strings.Replace(s, "(", "-", -1)
	ret = strings.Replace(ret, ")", "-", -1)
	return ret
}
