package fqdn_examples

import (
	"fmt"

	"github.com/Showmax/go-fqdn"
)

func ExampleFqdnHostname() {
	fqdn, err := fqdn.FqdnHostname()
	if err != nil {
		panic(err)
	}
	fmt.Println(fqdn)
}
