#include <cstdlib>
#include <cstring>
#include <exception>
#include <stdexcept>
#include <string>

#include "libmarisa.h"
#include "marisa.h"
#include "shim.h"

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
    } catch (const go::error &ex) {                                             \
        const char* b = "go shim error: ";                                      \
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

extern "C" void marisa_read_all(int iid, char ***out_wd, size_t *out_wd_sz, char **out_err) {
    try_cstr(out_err) {
        if (!out_wd || !out_wd_sz || !out_err)
            throw std::runtime_error("parameter is null");
        go::pstream r(iid);
        marisa::Trie t;
        marisa::read(r, &t);
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

extern "C" void marisa_write_all(int iid, const char** in_wd, size_t in_wd_sz, char **out_err) {
    try_cstr(out_err) {
        if ((in_wd_sz && !in_wd) || !out_err)
            throw std::runtime_error("parameter is null");
        marisa::Keyset k;
        for (size_t i = 0; i < in_wd_sz; i++)
            k.push_back(in_wd[i]);
        marisa::Trie t;
        t.build(k);
        go::pstream w(iid);
        marisa::write(w, t);
    } catch_cstr(out_err)
}

extern "C" void marisa_wd_free(char **wd, size_t wd_sz) {
    for (size_t i = 0; i < wd_sz; i++)
        free(wd[i]);
    free(wd);
}
