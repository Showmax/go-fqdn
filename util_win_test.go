// +build windows

package fqdn

import (
	"testing"
)

func TestReadlineWin(t *testing.T) {
	testCases := []readlineTestCase{
		{"foo\r\nbar\r\nbaz\r\n", []string{"foo", "bar", "baz"}},
		{"foo\r\nbar\r\nbaz", []string{"foo", "bar", "baz"}},
		{"foo\r\nbar\r\nbz\r\n\r\n", []string{"foo", "bar", "bz", ""}},
		{"foo\nbar\rbz\n\r", []string{"foo", "bar\rbz", ""}},
		{"foo\nbar\r\nbz\r\n\n", []string{"foo", "bar", "bz", ""}},
	}

	testReadline(t, testCases)
}
