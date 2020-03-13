package fqdn

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

type readlineTestCase struct {
	in  string
	out []string
}

func testReadline(t *testing.T, testCases []readlineTestCase) {
	for _, tc := range testCases {
		var e error
		var l string

		debug("Testing with: %q\n", tc.in)

		r := bufio.NewReader(strings.NewReader(tc.in))
		i := 0

		for l, e = readline(r); e == nil; l, e = readline(r) {
			if i >= len(tc.out) {
				t.Fatalf("Too many lines received")
			}

			if tc.out[i] != l {
				t.Fatalf("Line does not match.\n"+
					"\tExpected: %q\n"+
					"\tActual  : %q\n",
					tc.out[i], l)
			}
			i += 1
		}

		if e != io.EOF {
			t.Fatalf("Expected EOF, but exception is %T.", e)
		}

		if i != len(tc.out) {
			t.Fatalf("Not enough lines received")
		}
	}
}

func TestReadline(t *testing.T) {
	testCases := []readlineTestCase{
		{"foo\nbar\nbaz\n", []string{"foo", "bar", "baz"}},
		{"foo\nbar\nbaz", []string{"foo", "bar", "baz"}},
		{"foo\nbar\nbaz\n\n", []string{"foo", "bar", "baz", ""}},
		{"foo\nbar\nbaz\n\nx", []string{"foo", "bar", "baz", "", "x"}},
	}

	testReadline(t, testCases)
}
