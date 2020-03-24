// Package marisa is imported with _ to enable marisa for the kobodict, if
// supported. It is in a separate package so functions in kobodict which don't
// require marisa can be used without compiling it. As an alternative to
// importing this package, you can provide your own implementation of marisa in
// kobodict.Marisa. If imported, this package will fail to compile unless marisa
// is available for your GOOS/GOARCH.
package marisa

import "github.com/geek1011/dictutil/kobodict"

// This is done so it can still be instantiated even if not implemented for the
// current platform (it will be caught when assigning it to kobodict.Marisa),
// named platform for better error messages.

type platform struct{}

func init() {
	kobodict.Marisa = new(platform) // platform-specific implementation
}
