package match

import (
	"fmt"
	"regexp"
	"strings"
)

// Best find best match of the input string in the provided list.
func Best(
	inp string, // input string to match
	list []string, // list of strings to be matched
	ignore bool, // true to ignore case
) (
	fltr []int, // indices of 'list' matching inp
	mth string, // result if a single value in 'list' is matched
	typ uint8, // type of match: 2 - exact match; 1 - partial match; 0 - regex characters match
) {
	fltr = make([]int, 0)
	li := len(inp)

	for i, str := range list {
		lc := len(str)
		if li < lc {
			if strings.Contains(str, inp) || (ignore && strings.Contains(strings.ToUpper(str), strings.ToUpper(inp))) {
				fltr = append(fltr, i)
			}
		} else if li > lc {
			if strings.Contains(inp, str) || (ignore && strings.Contains(strings.ToUpper(inp), strings.ToUpper(str))) {
				fltr = append(fltr, i)
			}
		} else {
			if inp == str || (ignore && strings.EqualFold(inp, str)) {
				typ = 2
				fltr = []int{i}
				mth = list[fltr[0]]
				return // use the first exact match found
			}
		}
	}

	// partial match
	lf := len(fltr)
	if lf >= 1 {
		typ = 1
		if lf == 1 {
			mth = list[fltr[0]]
		}
		return
	}

	// regex match
	var pstr string
	for _, r := range inp {
		if r == '-' {
			continue
		}
		pstr = fmt.Sprintf("%v.*%c", pstr, r)
	}
	icase := ""
	if ignore {
		icase = "(?i)"
	}
	var pttn = regexp.MustCompile(fmt.Sprintf("%v%v.*", icase, pstr))

	for i, str := range list {
		if pttn.MatchString(str) {
			fltr = append(fltr, i)
		}
	}
	lf = len(fltr)
	if lf == 1 {
		mth = list[fltr[0]]
	}
	return
}
