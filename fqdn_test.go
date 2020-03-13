package fqdn

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// This package is hard to reasonably test in isolation, so take a shortcut and
// assume that no one will set their hostname to localhost.
func TestFqdnHostname(t *testing.T) {
	fqdnHost, err := FqdnHostname()
	if err != nil {
		t.Fatalf("Could not fqdn hostname: %v", err)
	}

	if fqdnHost == "localhost" {
		t.Fatalf("Unexpected fqdn, got: %s", fqdnHost)
	}

	if net.ParseIP(fqdnHost) != nil {
		t.Fatalf("Got IP address: %s", fqdnHost)
	}
}

func TestFromLookup(t *testing.T) {
	testCases := []struct {
		host string
		err  error
		fqdn string
	}{
		// I mean, these 2 are probably the most static IPs I can get
		{"ipv4.google.com", nil, "ipv4.l.google.com"},
		{"ipv6.google.com", nil, "ipv6.l.google.com"},
		{"makwjefalurgaf8", ErrFqdnNotFound, ""},
	}

	for _, tc := range testCases {
		fqdn, err := fromLookup(tc.host)
		if !errors.Is(err, tc.err) {
			t.Fatalf("Unexpected error.\n"+
				"\tExpected: %T\n"+
				"\tActual  : %T\n",
				tc.err, err)
		}
		if fqdn != tc.fqdn {
			t.Fatalf("Fqdn does not match.\n"+
				"\tExpected: %q\n"+
				"\tActual  : %q\n",
				tc.fqdn, fqdn)
		}
	}
}

func cat(file string) {
	content, err := exec.Command("cat", file).Output()
	if err != nil {
		// This probably means we are on windows
		debug("Could not cat %q: %v", file, err)
		return
	}

	debug("%s:\n", file)
	debug("------------------\n")
	debug("%q\n", content)
	debug("------------------\n")
	debug("%s\n", content)
	debug("------------------\n")

}

// In order to behave in expected way, we should verify that we are producing
// same output has hostname utility.
func TestMatchHostname(t *testing.T) {
	cat("/etc/hosts")
	cat("/etc/resolv.conf")

	out, err := exec.Command(hostnameBin, hostnameArgs...).Output()
	if err != nil {
		t.Fatalf("Could not run hostname: %v", err)
	}
	outS := chomp(string(out))

	fqdn, err := FqdnHostname()
	if err != nil {
		t.Fatalf("Could not fqdn hostname: %v", err)
	}

	// Since hostnames (domains) are case-insensitive and mac's hostname
	// returns it with uppercased first letter causing test to fail
	//
	//         	Us  : "mac-1271.local"
	//         	Them: "Mac-1271.local"
	//
	// we should compare lower-cased versions.
	outS = strings.ToLower(outS)
	fqdn = strings.ToLower(fqdn)

	if outS != fqdn {
		t.Fatalf("Output from hostname does not match.\n"+
			"\tUs  : %q\n"+
			"\tThem: %q\n",
			fqdn, outS)
	}
}

func TestParseHosts(t *testing.T) {
	testCases := []struct {
		hosts string
		host  string
		fqdn  string
		err   error
	}{
		{
			`# Static table lookup for hostnames.
# See hosts(5) for details.
127.0.0.1       foo`, "foo", "foo", nil,
		},
		{
			`# Static table lookup for hostnames.
# See hosts(5) for details.
127.0.0.1       bar.foo foo`, "foo", "bar.foo", nil,
		},
		{
			`# Static table lookup for hostnames.
# See hosts(5) for details.
127.0.0.1       yy bar
127.0.0.1       bar.foo foo
127.0.0.1       xx bar`, "foo", "bar.foo", nil,
		},
		{
			// This one is interesting, since it hostname -f with
			// this /etc/hosts gives you different results on musl-c
			// and glibc. I've picked the glibc behaviour, since we
			// can stop on first match.
			`# Static table lookup for hostnames.
# See hosts(5) for details.
127.0.0.1       bar.foo foo
127.0.0.1       foo.bar foo`, "foo", "bar.foo", nil,
		},
	}

	for _, tc := range testCases {
		hosts, err := ioutil.TempFile("", "go-fqdn.hosts.")
		if err != nil {
			panic(err)
		}
		defer os.Remove(hosts.Name())

		if _, err = hosts.Write([]byte(tc.hosts)); err != nil {
			panic(err)
		}
		hostsPath = hosts.Name()

		fqdn, err := fromHosts(tc.host)

		if !errors.Is(err, tc.err) {
			t.Fatalf("Unexpected error.\n"+
				"\tExpected: %T\n"+
				"\tActual  : %T\n",
				tc.err, err)
		}

		if fqdn != tc.fqdn {
			t.Fatalf("Fqdn does not match.\n"+
				"\tExpected: %q\n"+
				"\tActual  : %q\n",
				tc.fqdn, fqdn)
		}
	}
}

func TestParseHostLine(t *testing.T) {
	testCases := []struct {
		host string
		line string
		fqdn string
		ok   bool
	}{
		{"foo", "::1 foo bar", "foo", true},
		{"foo", "127.0.0.1 foo bar", "foo", true},
		{"bar", "::1 foo bar", "foo", true},
		{"bar", "::1 \t foo  \t\t\t  bar  \t\t", "foo", true},
		{"bar", "127.0.0.1 foo bar", "foo", true},
		{"bar", "127.0.0.1 foo.full bar", "foo.full", true},
		{"foo", "::1 foo", "foo", true},
		{"foo", "::1 bar", "", false},
		{"::1", "::1", "", false},
		{"127.0.0.1", "127.0.0.1", "", false},
		{"bar", "127.0.0.1 foo # bar", "", false},
		{"bar", "127.0.0.1 foo#bar", "", false},
		{"bar", "127.0.0.1\tfoo#bar  asdawdf a#", "", false},
		{"b", "127.0.0.1 a b", "a", true},
		{"a", "127.0.0.1 a b", "a", true},
		{"c", "127.0.0.1 a b", "", false},
		{"b", "127.0.0.1 _invalid_ b", "", false},
		{"b", "127.0.0.1 今日は b", "今日は", true},
	}

	for _, tc := range testCases {
		fqdn, ok := parseHostLine(tc.host, tc.line)

		if ok != tc.ok {
			t.Fatalf("Wrong ok value.\n"+
				"\tExpected: %t\n"+
				"\tActual  : %t\n",
				tc.ok, ok)
		}

		if fqdn != tc.fqdn {
			t.Fatalf("Wrong fqdn value.\n"+
				"\tExpected: %q\n"+
				"\tActual  : %q\n",
				tc.fqdn, fqdn)
		}
	}
}
