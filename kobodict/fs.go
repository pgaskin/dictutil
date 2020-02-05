package kobodict

// Unpack is a helper function to unpack the contents of a Reader to a folder
// on-disk. The provided dir must be non-existent.
func Unpack(r *Reader, dir string) error {
	panic("not implemented") // TODO(v0)
}

// Pack is a helper function to pack the contents a folder unpacked using Unpack
// into a Writer. It is assumed that the writer has not been used. The provided
// file will be overwritten if it exists and is a regular file, or created if it
// doesn't exist.
func Pack(r *Writer, file string) error {
	panic("not implemented") // TODO(v0)
	// remember to sort filenames to make it deterministic
}
