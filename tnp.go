// Copyright ieee0824

/*
This package parses Japanese phone numbers.
*/
package tnp

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type TelType int

func (t TelType) String() string {
	if t < 0 || len(typeNames) <= int(t) {
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
	InformationCharge
	FMC
)

var typeNames = []string{
	"fixed line phone",
	"m2m",
	"pocket bell",
	"ip phone",
	"mobile phone",
	"incoming charge",
	"unified number",
	"information charge",
	"fmc",
}

var (
	ignoreTypes   []TelType
	ignoreTypesMu sync.RWMutex
)

// SetIgnoreTypes sets the types to exclude from detection. Calling it without
// arguments clears the exclusions.
func SetIgnoreTypes(t ...TelType) {
	ignoreTypesMu.Lock()
	ignoreTypes = append([]TelType(nil), t...)
	ignoreTypesMu.Unlock()
}

func isIgnore(t TelType) bool {
	ignoreTypesMu.RLock()
	defer ignoreTypesMu.RUnlock()
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
		regexp.MustCompile(`0200-\d{5}-\d{5}`),

		// Draft implementation
		regexp.MustCompile(`020([1-3]|[5-9])\d{2}\d{5}`),
		regexp.MustCompile(`0200\d{10}`),
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
		regexp.MustCompile(`0[6-9]0-[1-9]\d{2}-\d{5}`),
		regexp.MustCompile(`0[6-9]0-[1-9]\d{3}-\d{4}`),

		// Draft implementation
		regexp.MustCompile(`0[6-9]0[1-9]\d{2}\d{5}`),
		regexp.MustCompile(`0[6-9]0[1-9]\d{3}\d{4}`),
	},
	[]*regexp.Regexp{
		regexp.MustCompile(`0120-\d{3}-\d{3}`),
		regexp.MustCompile(`0800-\d{3}-\d{4}`),

		// Draft implementation
		regexp.MustCompile(`0120\d{3}\d{3}`),
		regexp.MustCompile(`0800\d{3}\d{4}`),
	},
	[]*regexp.Regexp{
		regexp.MustCompile(`0570-\d{3}-\d{3}`),

		// Draft implementation
		regexp.MustCompile(`0570\d{3}\d{3}`),
	},
	[]*regexp.Regexp{
		regexp.MustCompile(`0990-\d{3}-\d{3}`),
		regexp.MustCompile(`0990\d{6}`),
	},
	[]*regexp.Regexp{
		regexp.MustCompile(`0600-\d{3}-\d{4}`),
		regexp.MustCompile(`0600\d{7}`),
	},
}

// Check specific number types before fixed lines, whose patterns can also
// match some unhyphenated service numbers.
var detectionOrder = []TelType{
	M2M, PocketBell, IPPhone, MobilePhone, IncomingCharge, UnifiedNumber,
	InformationCharge, FMC, FixedLinePhone,
}

// IsTelNumber reports whether the entire string is a phone number.
func IsTelNumber(s string) (bool, TelType) {
	_, kind, found := findTelNumber(s, true)
	return found, kind
}

// CropTelNumber returns the first phone number found in a string.
func CropTelNumber(s string) (string, error) {
	number, _, found := findTelNumber(s, false)
	if !found {
		return "", fmt.Errorf("not tel number")
	}
	return number, nil
}

func findTelNumber(s string, entire bool) (string, TelType, bool) {
	if hasParenthesis(s) {
		s = replaceParenthesis(s)
	}

	bestStart, bestEnd := len(s)+1, -1
	var bestType TelType
	for _, kind := range detectionOrder {
		if isIgnore(kind) {
			continue
		}
		for _, r := range telPatternRegs[kind] {
			for _, loc := range r.FindAllStringIndex(s, -1) {
				start, end := loc[0], loc[1]
				if entire && (start != 0 || end != len(s)) {
					continue
				}
				if start > 0 && isDigit(s[start-1]) || end < len(s) && isDigit(s[end]) {
					continue
				}
				if kind == FixedLinePhone && isServicePrefix(s[start:end]) {
					continue
				}
				if start < bestStart || start == bestStart && end > bestEnd {
					bestStart, bestEnd, bestType = start, end, kind
				}
			}
		}
	}
	if bestEnd < 0 {
		return "", -1, false
	}
	return s[bestStart:bestEnd], bestType, true
}

func isServicePrefix(s string) bool {
	for _, prefix := range []string{"0120", "0170", "0180", "0570", "0800", "0990"} {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
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
