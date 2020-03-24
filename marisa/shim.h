#ifndef GO_SHIM_H
#define GO_SHIM_H

#ifdef __cplusplus
#include <cstddef>
extern "C" {
#else
#include <stdbool.h>
#include <stddef.h>
#endif

// go_iop_type represents interfaces an iid may implement.
enum go_iop_type {
    reader = 1 << 0, // io.Reader
    writer = 1 << 1, // io.Writer
};

// go_iop_reader checks if the specified iid implements the specified ORed type
// flags. Note that it doesn't have to be checked here, as go_iop_* will return
// an error if it doesn't implement the necessary interfaces. If out_err is not
// NULL and the return value is false, it will be set to an error message, which
// must be freed by the caller, if the iid doesn't implement the specified
// flags.
bool go_iop_check(int iid, int t, char **out_err);

// Note: we use ptrdiff_t over ssize_t for portability (and not size_t because
// it will return -1 for EOF). Also, note that C++'s std::streamsize uses this
// internally too, which is a nice advantage.

// go_iop_read reads from the iid's underlying io.Reader. It has the same
// semantics as the Go one, but io.EOF is returned as -1. out_err must be a
// valid pointer to a char pointer. If an error occured, it is set and must be
// freed by the caller.
ptrdiff_t go_iop_read(int iid, const char *p, size_t n, char **out_err);
// go_iop_write writes to the iid's underlying io.Writer. It has the same
// semantics as the Go one, but io.EOF is returned as -1. out_err must be a
// valid pointer to a char pointer. If an error occured, it is set and must be
// freed by the caller.
ptrdiff_t go_iop_write(int iid, const char *p, size_t n, char **out_err);

#ifdef __cplusplus
}

#include <cstdarg>
#include <cstdlib>
#include <iostream>
#include <stdexcept>

// https://golang.org/cmd/cgo/#hdr-C_references_to_Go
// https://en.cppreference.com/w/cpp/io/basic_streambuf <- this describes it better than many of the other sites I found

namespace go {

bool dbg(const char* format, ...) {
    static bool _dbg = getenv("GOSHIMDEBUG") ? getenv("GOSHIMDEBUG")[0] == '1' && getenv("GOSHIMDEBUG")[1] == '\0' : false;
    if (!_dbg)
        return false;
    fprintf(stderr, "GOSHIMDEBUG: ");
    va_list arg;
    va_start(arg, format);
    vfprintf(stderr, format, arg);
    va_end(arg);
    fflush(stderr);
    return true;
}

class error : public std::runtime_error {
public:
    error(const char* what) : std::runtime_error(what) {
        go::dbg("new go::error(%s)\n", what);
    };

    // check checks an output err pointer and frees+throws it if set.
    static void check(char* err) {
        if (!err)
            return;
        go::error ex = go::error(err);
        free(err);
        throw ex;
    }
};

class iopbuf : public std::basic_streambuf<char> {
    int iid_;
    char rbuf_; // single-byte read buffer (i.e. direct access to the io.Reader)
public:
    static_assert((std::is_same<iopbuf::char_type, char>::value && std::is_same<iopbuf::traits_type::char_type, char>::value), "Go shim only supports char"); // just to be safe
    #ifndef __clang__
    static_assert(iopbuf::traits_type::eof() != iopbuf::traits_type::to_int_type((char) 0xFF), "EOF not distinct from 0xFF"); // this is already specified in the spec, but just to make sure
    #endif

    iopbuf(int iid) : iid_(iid) {
        this->setg(&this->rbuf_, &this->rbuf_ + 1, &this->rbuf_ + 1); // set the buffer, but at the end to force the next read to underflow
    }

    iopbuf(int iid, int t) : iopbuf(iid) {
        char* err = NULL;
        go_iop_check(iid, t, &err);
        go::error::check(err);
    }

    iopbuf::int_type underflow() override {
        // This is all that's strictly needed for reading. Note that we can't
        // just return the char, and we must set the buffer to point to it to
        // conform to the expected postconditions and prevent unusual bugs from
        // popping up.

        char* err = NULL;
        ptrdiff_t n = go_iop_read(this->iid_, &this->rbuf_, 1, &err);
        go::dbg("underflow: go_iop_read(%d, 1) = %td %02x err=%s\n", this->iid_, n, this->rbuf_, err); fflush(stdout);
        go::error::check(err);

        this->setg(&this->rbuf_, &this->rbuf_, &this->rbuf_ + (n>0 ? n : 0));   // Update the current byte.
        return this->gptr() == this->egptr()                                    // If the new current pos == past end of buffer, no byte was read (n<=0).
            ? iopbuf::traits_type::eof()                                        // If no byte was read (and no error was thrown earlier), it's an EOF.
            : iopbuf::traits_type::to_int_type(this->rbuf_);                    // Otherwise, return the byte we just read (note: without to_int_type, 0xFF would be sign extended to -1/eof).
    }

    std::streamsize xsgetn(iopbuf::char_type* buf, std::streamsize buf_n) override {
        // We can provide a more efficient bulk read implementation than the
        // default one which gets each byte one-by-one in a loop.
        // Note: Remember to test ::underflow by forcing it to use the default
        // implementation: return std::streambuf::xsgetn(buf, buf_n);

        std::streamsize t = 0;

        ptrdiff_t n = 0;
        char* err = NULL;
        while (t != buf_n && n != -1) {
            n = go_iop_read(this->iid_, buf+t, buf_n-t, &err);
            go::dbg("xsgetn: go_iop_read(%d, %zu) = %td (%td/%td) err=%s\n", this->iid_, buf_n-t, n, t+(n>0 ? n : 0), buf_n, err); fflush(stdout);
            t += n>0 ? n : 0;
            if (t > buf_n)
                throw go::error("read returned too many bytes!");
            go::error::check(err);
        }

        this->rbuf_ = t>0 ? buf[t-1] : 0;                                       // Set the current byte to the last one read, if any.
        this->setg(&this->rbuf_, &this->rbuf_, &this->rbuf_ + (t>0 ? 1 : 0));   // Update the current byte.
        return this->gptr() == this->egptr()                                    // If the new current pos == past end of buffer, no byte was read (n<=0).
            ? iopbuf::traits_type::eof()                                        // If no byte was read (and no error was thrown earlier), it's an EOF
            : t;                                                                // Otherwise, return the number of bytes read.
    }

    iopbuf::int_type overflow(iopbuf::int_type c = iopbuf::traits_type::eof()) override {
        // Unlike for reading, we don't have to use a buffer (you can read a
        // byte advancing, but you can't do that kind of thing when writing),
        // so we'll just write it directly. This makes the implementation much
        // simpler, as we're basically just passing the calls to the Go funcs
        // directly.

        // Usually, we would flush the buffer if given an EOF instead of a char,
        // but we're not using one, so it's a no-op.
        if (iopbuf::traits_type::eq_int_type(c, iopbuf::traits_type::eof()))
            return 0;

        // Since the logic is basically a simplified version of xsputn, just
        // with a single char, it's easier just to call it and implement the
        // bulk of the logic there.
        if (this->xsputn(reinterpret_cast<iopbuf::traits_type::char_type*>(&c), 1) != 1)
            throw go::error("short write"); // we still need to check for a short write
        return c;
    }

    std::streamsize xsputn(const iopbuf::char_type* buf, std::streamsize buf_n) override {
        char* err = NULL;
        ptrdiff_t n = go_iop_write(this->iid_, buf, buf_n, &err);
        go::error::check(err);
        if (n == -1)
            throw go::error("EOF while writing to Go writer");
        return n;
    }
};

class rwstream : private iopbuf, public std::iostream {
public: rwstream(int iid) : iopbuf(iid, go_iop_type::reader|go_iop_type::writer), std::iostream(this) {}
};

class wstream : private iopbuf, public std::ostream {
public: wstream(int iid) : iopbuf(iid, go_iop_type::writer), std::ostream(this) {}
};

class rstream : private iopbuf, public std::istream {
public: rstream(int iid) : iopbuf(iid, go_iop_type::reader), std::istream(this) {}
};

}

#endif
#endif