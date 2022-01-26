//go:build windows
// +build windows

package fqdn

const hostnameBin = "hostname"

var hostnameArgs = []string{} //nolint:gochecknoglobals
