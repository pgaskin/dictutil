package marisa

//go:generate go run -tags libmarisa_generate libmarisa_generate.go
//go:generate swig -go -cgo -intgosize 64 -module marisa -outdir . -o marisa.cc -c++ marisa_swig.swigcxx_

// #cgo CPPFLAGS:
// #cgo LDFLAGS:
import "C"
