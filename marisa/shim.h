#ifndef GO_SHIM_H
#define GO_SHIM_H

#ifdef __cplusplus
#include <cstddef>
extern "C" {
#else
#include <stddef.h>
#endif

int go_iop_read(int iid, const char *p, size_t n, char **out_err);
int go_iop_write(int iid, const char *p, size_t n, char **out_err);

#ifdef __cplusplus
}

#include <cstdlib>
#include <ios>
#include <iostream>
#include <stdexcept>

// https://accu.org/index.php/journals/264
// https://golang.org/cmd/cgo/#hdr-C_references_to_Go
// http://www.cplusplus.com/reference/streambuf/streambuf/overflow/
// http://www.cplusplus.com/reference/streambuf/streambuf/xsputn/
// http://www.cplusplus.com/reference/streambuf/streambuf/underflow/
// http://www.cplusplus.com/reference/streambuf/streambuf/xsgetn/

namespace go {

class error : public std::runtime_error {
public:
    error(const char* what) : std::runtime_error(what) {};
};

class pbuf : public std::streambuf {
    int iid_;
public:
    pbuf(int iid) : iid_(iid) {}

    // --- reading
    int underflow() override {
        char c;
        return this->xsgetn(&c, 1) == 0
            ? std::char_traits<char>::eof()
            : static_cast<int>(c);
    }
    // the other members don't need to be overridden, as they just call overflow
    // if necessary (see http://www.cplusplus.com/reference/streambuf/streambuf)
    // but, we'll implement xsputn for more efficient bulk writes
    std::streamsize xsgetn(char* buf, std::streamsize buf_n) override {
        char* err = NULL;
        int n = go_iop_read(this->iid_, buf, buf_n, &err);
        if (err) {
            go::error ex = go::error(err);
            free(err);
            throw ex;
        }
        return n == -1
            ? std::char_traits<char>::eof()
            : static_cast<std::streamsize>(n);
    }

    // --- writing
    int overflow(int c = std::char_traits<char>::eof()) override {
        if (std::char_traits<char>::eq_int_type(c, std::char_traits<char>::eof()))
            return 0; // usually, we would flush the buffer here, but we're writing directly, so it's basically a nop

        char c_ = c;
        if (this->xsputn(&c_, 1) != 1)
            throw go::error("short write");
        return c;
    }
    // the other members don't need to be overridden, as they just call overflow
    // if necessary (see http://www.cplusplus.com/reference/streambuf/streambuf)
    // but, we'll implement xsputn for more efficient bulk writes
    std::streamsize xsputn(const char* buf, std::streamsize buf_n) override {
        char* err = NULL;
        int n = go_iop_write(this->iid_, buf, buf_n, &err);
        if (err) {
            go::error ex = go::error(err);
            free(err);
            throw ex;
        }
        if (n == -1)
            throw go::error("EOF while writing to Go writer");
        return static_cast<std::streamsize>(n);
    }
};

struct basic_pstream {
    pbuf sbuf_;
    basic_pstream(int iid) : sbuf_(iid) {}
};

class pstream : virtual basic_pstream, public std::iostream {
public:
    pstream(int iid) : basic_pstream(iid)
        , std::ios(&this->sbuf_)
        , std::iostream(&this->sbuf_) {}
};

}

#endif
#endif