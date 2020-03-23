#ifdef __cplusplus
#include <cstddef>
extern "C" {
#else
#include <stddef.h>
#endif

void marisa_read_all(int iid, char ***out_wd, size_t *out_wd_sz, char **out_err);
void marisa_write_all(int iid, const char** in_wd, size_t in_wd_sz, char **out_err);
void marisa_wd_free(char **wd, size_t wd_sz);
void marisa_go_test_error_helper(int at, int val);

#ifdef __cplusplus
}
#endif