#include <cstdlib>
#include <cstring>
#include <exception>
#include <stdexcept>
#include <string>

#include "libmarisa.h"
#include "shim.h"

#define catch_go_ex(t, ctx)                                                     \
    catch (const t &ex) {                                                       \
        const char* b = ctx;                                                    \
        char* err = reinterpret_cast<char*>(                                    \
            calloc(strlen(b)+strlen(ex.what())+1, sizeof(char)));               \
        strcpy(err, b);                                                         \
        strcat(err, ex.what());                                                 \
        return err;                                                             \
    }

#define catch_go                                                                \
    catch_go_ex(marisa::Exception, "marisa: ")                                  \
    catch_go_ex(go::error, "go shim: ")                                         \
    catch_go_ex(std::runtime_error, "c++ runtime: ")                            \
    catch_go_ex(std::exception, "c++ error: ")                                  \
    catch (...) { return strdup("marisa: unknown c++ exception"); }             \
    return NULL;

#define go_func extern "C" const char*

go_func marisa_read_all(int iid, char ***out_wd, size_t *out_wd_sz) {
    try {
        if (!out_wd || !out_wd_sz)
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
    } catch_go
}

go_func marisa_write_all(int iid, const char** wd, size_t wd_sz) {
    try {
        if (wd_sz && !wd)
            throw std::runtime_error("parameter is null");
        marisa::Keyset k;
        for (size_t i = 0; i < wd_sz; i++)
            k.push_back(wd[i]);
        marisa::Trie t;
        t.build(k);
        go::pstream w(iid);
        marisa::write(w, t);
    } catch_go
}
