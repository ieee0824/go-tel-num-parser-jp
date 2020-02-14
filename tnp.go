// Copyright ieee0824

/*
This package parses Japanese phone numbers.
*/
package tnp

import (
	"fmt"
	"regexp"
	"strings"
)

type TelType int

func (t TelType) String() string {
	if t == -1 || len(typeNames) <= int(t) {
		return "not tel number"
	}
	return typeNames[t]
}

const (
	FixedLinePhone TelType = iota
	M2M
	PocketBell
	IPPhone
	MobilePhone
	IncomingCharge
	UnifiedNumber
)

var typeNames = []string{
	"fixed line phone",
	"m2m",
	"pocket bell",
	"ip phone",
	"mobile phone",
	"incoming charge",
	"unified number",
}

var ignoreTypes = []TelType{}

// SetIgnoreTypes registers a list to be ignored from the judgment conditions.
func SetIgnoreTypes(t ...TelType) {
	ignoreTypes = append(ignoreTypes, t...)
}

func isIgnore(t TelType) bool {
	for _, v := range ignoreTypes {
		if v == t {
			return true
		}
	}
	return false
}

var telPatternRegs = [][]*regexp.Regexp{
	[]*regexp.Regexp{
		regexp.MustCompile(`0[0-9]-[2-9]\d{3}-\d{4}`),
		regexp.MustCompile(`0\d{2}-[2-9]\d{2}-\d{4}`),
		regexp.MustCompile(`0\d{3}-[2-9]\d{1}-\d{4}`),
		regexp.MustCompile(`0\d{4}-[2-9]-\d{4}`),

		// Draft implementation
		regexp.MustCompile(`0[0-9][2-9]\d{3}\d{4}`),
		regexp.MustCompile(`0\d{2}[2-9]\d{2}\d{4}`),
		regexp.MustCompile(`0\d{3}[2-9]\d{1}\d{4}`),
		regexp.MustCompile(`0\d{4}[2-9]\d{4}`),
	},
	[]*regexp.Regexp{
		regexp.MustCompile(`020-([1-3]|[5-9])\d{2}-\d{5}`),

		// Draft implementation
		regexp.MustCompile(`020([1-3]|[5-9])\d{2}\d{5}`),
	},
	[]*regexp.Regexp{
		regexp.MustCompile(`020-4\d{2}-\d{5}`),

		// Draft implementation
		regexp.MustCompile(`0204\d{2}\d{5}`),
	},
	[]*regexp.Regexp{
		regexp.MustCompile(`050-[1-9]\d{3}-\d{4}`),

		// Draft implementation
		regexp.MustCompile(`050[1-9]\d{3}\d{4}`),
	},
	[]*regexp.Regexp{
		regexp.MustCompile(`0[7-9]0-[1-9]\d{2}-\d{5}`),
		regexp.MustCompile(`0[7-9]0-[1-9]\d{3}-\d{4}`),

		// Draft implementation
		regexp.MustCompile(`0[7-9]0[1-9]\d{2}\d{5}`),
		regexp.MustCompile(`0[7-9]0[1-9]\d{3}\d{4}`),
	},
	[]*regexp.Regexp{
		regexp.MustCompile(`0120-\d{3}-\d{3}`),
		regexp.MustCompile(`0800-\d{3}-\d{3}`),

		// Draft implementation
		regexp.MustCompile(`0120\d{3}\d{3}`),
		regexp.MustCompile(`0800\d{3}\d{3}`),
	},
	[]*regexp.Regexp{
		regexp.MustCompile(`0570-\d{3}-\d{3}`),

		// Draft implementation
		regexp.MustCompile(`0570\d{3}\d{3}`),
	},
}

// IsTelNumber determines that it is a phone number.
func IsTelNumber(s string) (bool, TelType) {
	for k, rs := range telPatternRegs {
		if isIgnore(TelType(k)) {
			continue
		}
		for _, r := range rs {
			if r.Copy().MatchString(s) {
				return true, TelType(k)
			}
		}
	}

	if hasParenthesis(s) {
		return IsTelNumber(replaceParenthesis(s))
	}

	return false, -1
}

// CropTelNumber retrieves when a phone number exists.
func CropTelNumber(s string) (string, error) {
	for k, rs := range telPatternRegs {
		if isIgnore(TelType(k)) {
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
