/*
shellescape provides the Quote to escape arbitrary
strings for a safe use as command line arguments in the most common
POSIX shells.
The Python package which this work was inspired by can be found
at https://pypi.python.org/pypi/shellescape.
*/
package main

/*
The functionality provided by Quote could be helpful
in those cases where it is known that the output of a Go program will
be appended to/used in the context of shell programs' command line arguments.
*/

import (
	"regexp"
	"strings"
)

var pattern *regexp.Regexp

func init() {
	pattern = regexp.MustCompile(`[^\w@%+=:,./-]`)
}

// Quote returns a shell-escaped version of the string s. The returned value
// is a string that can safely be used as one token in a shell command line.
func Quote(s string) string {
	if len(s) == 0 {
		return "''"
	}

	if pattern.MatchString(s) {
		return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
	}

	return s
}

// QuoteCommand returns a shell-escaped version of the slice of strings.
// The returned value is a string that can safely be used as shell command arguments.
func QuoteCommand(args []string) string {
	l := make([]string, len(args))

	for i, s := range args {
		l[i] = Quote(s)
	}

	return strings.Join(l, " ")
}
