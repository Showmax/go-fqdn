# go-fqdn

Simple wrapper around `net` and `os` golang standard libraries providing Fully
Qualified Domain Name of the machine.

## Usage

This package uses go modules, so just writing code that uses it should be
enough.

```
package main

import (
	"fmt"
	"github.com/Showmax/go-fqdn"
)

func main() {
	fmt.Println(fqdn.Get())
}
```

We can then run it:

```
+   $ go run fqdn.go
go: finding module for package github.com/Showmax/go-fqdn
go: found github.com/Showmax/go-fqdn in github.com/Showmax/go-fqdn v0.0.0-20180501083314-6f60894d629f
localhost
```

`fqdn.Get()` returns:
- machine's FQDN if found.
- hostname if FQDN is not found.
- return "unknown" if nothing is found.

## Supported go versions

Current and current - 1 versions of go are supported.
