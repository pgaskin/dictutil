#include <cstdlib>
#include <cstring>
#include <exception>
#include <iostream>
#include <sstream>
#include <stdexcept>
#include <string>

#include "libmarisa.h"
#include "marisa.h"

#define try_cstr(out_err)                                                       \
    *(out_err) = NULL;                                                          \
    try

#define catch_cstr(out_err)                                                     \
    catch (const marisa::Exception &ex) {                                       \
        const char* b = "marisa: ";                                             \
        char* err = reinterpret_cast<char*>(                                    \
            calloc(strlen(b)+strlen(ex.what())+1, sizeof(char)));               \
        strcpy(err, b);                                                         \
        strcat(err, ex.what());                                                 \
        *(out_err) = err;                                                       \
        return;                                                                 \
    } catch (const std::runtime_error &ex) {                                    \
        const char* b = "c++ runtime error: ";                                  \
        char* err = reinterpret_cast<char*>(                                    \
            calloc(strlen(b)+strlen(ex.what())+1, sizeof(char)));               \
        strcpy(err, b);                                                         \
        strcat(err, ex.what());                                                 \
        *(out_err) = err;                                                       \
        return;                                                                 \
    } catch (const std::exception &ex) {                                        \
        const char* b = "c++ error: ";                                          \
        char* err = reinterpret_cast<char*>(                                    \
            calloc(strlen(b)+strlen(ex.what())+1, sizeof(char)));               \
        strcpy(err, b);                                                         \
        strcat(err, ex.what());                                                 \
        *(out_err) = err;                                                       \
        return;                                                                 \
    } catch (...) {                                                             \
        *(out_err) = strdup("marisa: unknown c++ exception");                   \
        return;                                                                 \
    }

class membuf : public std::basic_streambuf<char> {
public: membuf(const uint8_t *p, size_t l) { setg((char*)p, (char*)p, (char*)p + l); }
};

class memstream : public std::istream {
public: memstream(const char *p, size_t l) : std::istream(&_buffer), _buffer(reinterpret_cast<const uint8_t*>(p), l) { rdbuf(&_buffer); }
private: membuf _buffer;
};

extern "C" void marisa_read_all(const char* in_buf, size_t in_buf_sz, char ***out_wd, size_t *out_wd_sz, char **out_err) {
    try_cstr(out_err) {
        marisa_go_test_error_helper(-1, in_buf_sz);
        if ((in_buf_sz && !in_buf) || !out_wd || !out_wd_sz || !out_err)
            throw std::runtime_error("parameter is null");
        memstream b(in_buf, in_buf_sz);
        marisa::Trie t;
        marisa::read(b, &t);
        marisa::Agent a;
        a.set_query("");
        *out_wd_sz = 0;
        *out_wd = reinterpret_cast<char**>(calloc(t.num_keys(), sizeof(char**)));
        while (t.predictive_search(a)) {
            if (*out_wd_sz == t.num_keys())
                throw std::runtime_error("expected " + std::to_string(t.num_keys()) + " keys, got more");
            memcpy((*out_wd)[(*out_wd_sz)++] = reinterpret_cast<char*>(calloc(a.key().length()+1, sizeof(char))), a.key().ptr(), a.key().length());
        }
        if (*out_wd_sz != t.num_keys())
            throw std::runtime_error("expected " + std::to_string(t.num_keys()) + " keys, got " + std::to_string(*out_wd_sz));
    } catch_cstr(out_err)
}

extern "C" void marisa_write_all(const char** in_wd, size_t in_wd_sz, char **out_buf, size_t *out_buf_sz, char **out_err) {
    try_cstr(out_err) {
        marisa_go_test_error_helper(-1, in_wd_sz);
        if ((in_wd_sz && !in_wd) || !out_buf || !out_buf_sz || !out_err)
            throw std::runtime_error("parameter is null");
        marisa::Keyset k;
        for (size_t i = 0; i < in_wd_sz; i++)
            k.push_back(in_wd[i]);
        marisa::Trie t;
        t.build(k);
        std::stringstream b;
        marisa::write(b, t);
        std::string s = b.str();
        *out_buf_sz = s.length();
        *out_buf = reinterpret_cast<char*>(malloc(*out_buf_sz));
        memcpy(*out_buf, s.c_str(), *out_buf_sz);
    } catch_cstr(out_err)
}

extern "C" void marisa_wd_free(char **wd, size_t wd_sz) {
    for (size_t i = 0; i < wd_sz; i++)
        free(wd[i]);
    free(wd);
}

extern "C" void marisa_go_test_error_helper(int at, int val) {
    static int _at = -1;
    if (at == -1 && val != -1) {
        if (_at != -1) {
            if (val == _at)
                throw std::runtime_error("go_test_error");
            return;
        }
    } else if (val == -1) {
        // set the test value
        _at = at;
    }
    return;
}
