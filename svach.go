// Package svach file names for consistent naming across platforms.
// Incase of invalid fileName, a md5Sum of provided filename is returned
//
// Options
//
// `replaceStr`, replace invalid characters with given string, default "" (empty string)
//
// `maxLen`, limit length of the output, default `240`
package svach

import (
	"crypto/md5"
	"errors"
	"fmt"
	"html"
	"regexp"
	"strings"
)

// Svach object
type Svach struct {
	replaceStr string
	maxLen     int
}

var (
	// ErrCntrl ,Control characters exist in `replaceStr`
	ErrCntrl = errors.New("Control characters exist in replaceStr")
	// ErrInval ,Invalid characters exist in `replaceStr`
	ErrInval = errors.New("Invalid characters like `., <, >, :, \", /, \\, |, ?, *` exist in replaceStr")

	// ErrLen ,`maxLen` greater than 255
	ErrLen = errors.New("maxLen can't be greater than 255")
)

var iMaxLen = 240

// New returns a Svach object
func New() *Svach {
	return &Svach{"", iMaxLen}
}

var (
	cntrlExp = regexp.MustCompile("[[:cntrl:]]") // control

	// invalid characters - windows
	// <, >, :, ", /, \, |, ?, *
	invCharExp = regexp.MustCompile(`[<>:"/\\|\?\*]+`)

	// trim right spaces and dot
	rightSDExp = regexp.MustCompile("(?s:[[:space:]]|\\.)+$")

	// trim left dot's
	leftdotExpr = regexp.MustCompile("^\\.+")
)

// NewWithOpts returns Svach object, with opts set
// and returns error, if conditions aren't met.
//
// Conditions:
//
// - replaceStr can't contain Control or Invalid characters
//
// - maxlen can't be greater than 255
func NewWithOpts(replaceStr string, maxLen int) (*Svach, error) {
	s := &Svach{"", iMaxLen}

	if err := validOptStr(replaceStr); err != nil {
		return s, err
	}

	if maxLen > 255 {
		return s, ErrLen
	}

	return &Svach{replaceStr, maxLen}, nil
}

func validOptStr(s string) error {
	if cntrlExp.MatchString(s) {
		return ErrCntrl
	}

	if invCharExp.MatchString(s) || strings.Contains(s, ".") {
		return ErrInval
	}
	return nil
}

// Name svachs `fileName`
func (s *Svach) Name(fileName string) string {
	return name(fileName, s.replaceStr, s.maxLen)
}

func name(fileName, replaceStr string, maxLen int) string {

	intrStr := strings.ToValidUTF8(fileName, replaceStr) // intermediate string
	intrStr = html.UnescapeString(intrStr)
	intrStr = cntrlExp.ReplaceAllString(intrStr, replaceStr)
	intrStr = invCharExp.ReplaceAllString(intrStr, replaceStr)
	intrStr = rightSDExp.ReplaceAllString(intrStr, replaceStr)
	intrStr = leftdotExpr.ReplaceAllString(intrStr, ".")

	return validName(fileName, intrStr, replaceStr, maxLen)
}

var (
	// Unicode categories

	// Cc - Control
	// Cf - Format
	unicodeControl = regexp.MustCompile("\\p{Cc}|\\p{Cf}")

	// Zl - Line separator
	// Zp - Paragraph separator
	// Zs - Space separator
	unicodeSpace = regexp.MustCompile("\\p{Zl}|\\p{Zp}|\\p{Zs}")
)

var (
	// if below, characters are repeated more than twice,
	// we replace it with single character from `cleanReplaceWith`
	cleanExpr = repeatedCharsExp([]string{
		`[[:space:]]`, `_`, `-`, `\+`, `\.`, `!`})
	cleanReplaceWith = []string{"", " ", "_", "-", "+", ".", "!"}
)

// Clean svachs `fileName` into more humane format.
//
// Remove invisible and control characters, repeated separators.
// Replace different kinds of spaces with normal space.
func (s *Svach) Clean(fileName string) string {
	return clean(fileName, s.replaceStr, s.maxLen)
}

func clean(fileName, replaceStr string, maxLen int) string {

	intrStr := strings.ToValidUTF8(fileName, replaceStr)

	if intrStr != "" {
		intrStr = html.UnescapeString(intrStr)

		// invisible characters
		intrStr = unicodeControl.ReplaceAllString(intrStr, replaceStr)
		intrStr = unicodeSpace.ReplaceAllString(intrStr, " ")

		intrStr = invCharExp.ReplaceAllString(intrStr, replaceStr)

		var replaceExpr *regexp.Regexp
		if replaceStr != "" {
			replaceExpr = regexp.MustCompile("(" + replaceStr + ")" + "{2,}")
		}

		for startStr := intrStr; ; startStr = intrStr {

			// repeated separators
			intrStr = replaceAllStringSubmatch(
				cleanExpr, intrStr, cleanReplaceWith)

			// remove repeated `replaceStr`
			if replaceExpr != nil {
				intrStr = replaceExpr.ReplaceAllString(intrStr, replaceStr)
			}

			intrStr = rightSDExp.ReplaceAllString(intrStr, replaceStr)
			intrStr = leftdotExpr.ReplaceAllString(intrStr, ".")

			if startStr == intrStr {
				break
			}

		}

	}

	return validName(fileName, intrStr, replaceStr, maxLen)
}

func validName(fileName, intrStr, replaceStr string, maxLen int) string {

	if intrStr != replaceStr && intrStr != "" {
		if valid(intrStr) {
			if len(intrStr) > maxLen {
				return strings.ToValidUTF8(intrStr[:maxLen], "")
			}
			return intrStr
		}
	}

	return fmt.Sprintf("%x", md5.Sum([]byte(fileName)))
}

func valid(name string) bool {
	if name == "" {
		return false
	}

	if len(name) < 5 {
		for _, invName := range invalidNamesMap[len(name)] {
			if strings.ToLower(name) == invName {
				return false
			}
		}
	}

	return true
}

// https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file?redirectedfrom=MSDN#naming-conventions
// invalidNamesMap with len of value as key
var invalidNamesMap = map[int][]string{
	1: {
		".",
	},

	2: {
		"..",
	},

	3: {
		"con", "prn", "aux", "nul",
	},

	4: {
		"com1", "com2", "com3", "com4",
		"com5", "com6", "com7", "com8",
		"com9", "lpt1", "lpt2", "lpt3",
		"lpt4", "lpt5", "lpt6", "lpt7",
		"lpt8", "lpt9",
	},
}

func repeatedCharsExp(vals []string) *regexp.Regexp {
	var s strings.Builder
	for i, val := range vals {
		if i == 0 {
			s.WriteString("(" + val + "{2,})")
		} else {
			s.WriteString("|(" + val + "{2,})")
		}

	}
	return regexp.MustCompile(s.String())
}

// replaceAllStringSubmatch replaces matched groups with `replaceGroup`
func replaceAllStringSubmatch(re *regexp.Regexp, src string, replaceGroup []string) string {

	sms := re.FindAllStringSubmatchIndex(src, -1)

	if sms == nil || len(sms) == 0 {
		return src
	}

	if len(sms[0]) != len(replaceGroup)*2 {
		return ""
	}

	var s strings.Builder

	prevPos := 0
	for _, sm := range sms {

		for i := 2; i < len(sm); i += 2 {
			if sm[i] != -1 {
				start := sm[i]
				end := sm[i+1]

				s.WriteString(src[prevPos:start] + replaceGroup[i/2])

				prevPos = end
			}
		}
	}
	s.WriteString(src[prevPos:])
	return s.String()
}
