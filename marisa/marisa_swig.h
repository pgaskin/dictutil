#ifndef MARISA_SWIG_H_
#define MARISA_SWIG_H_

#include <libmarisa.h>

namespace marisa_swig {

enum ErrorCode {
    OK = MARISA_OK,
    STATE_ERROR = MARISA_STATE_ERROR,
    NULL_ERROR = MARISA_NULL_ERROR,
    BOUND_ERROR = MARISA_BOUND_ERROR,
    RANGE_ERROR = MARISA_RANGE_ERROR,
    CODE_ERROR = MARISA_CODE_ERROR,
    RESET_ERROR = MARISA_RESET_ERROR,
    SIZE_ERROR = MARISA_SIZE_ERROR,
    MEMORY_ERROR = MARISA_MEMORY_ERROR,
    IO_ERROR = MARISA_IO_ERROR,
    FORMAT_ERROR = MARISA_FORMAT_ERROR
};

enum NumTries {
    MIN_NUM_TRIES = MARISA_MIN_NUM_TRIES,
    MAX_NUM_TRIES = MARISA_MAX_NUM_TRIES,
    DEFAULT_NUM_TRIES = MARISA_DEFAULT_NUM_TRIES
};

enum CacheLevel {
    HUGE_CACHE = MARISA_HUGE_CACHE,
    LARGE_CACHE = MARISA_LARGE_CACHE,
    NORMAL_CACHE = MARISA_NORMAL_CACHE,
    SMALL_CACHE = MARISA_SMALL_CACHE,
    TINY_CACHE = MARISA_TINY_CACHE,
    DEFAULT_CACHE = MARISA_DEFAULT_CACHE
};

enum TailMode {
    TEXT_TAIL = MARISA_TEXT_TAIL,
    BINARY_TAIL = MARISA_BINARY_TAIL,
    DEFAULT_TAIL = MARISA_DEFAULT_TAIL
};

enum NodeOrder {
    LABEL_ORDER = MARISA_LABEL_ORDER,
    WEIGHT_ORDER = MARISA_WEIGHT_ORDER,
    DEFAULT_ORDER = MARISA_DEFAULT_ORDER
};

class Key {
public:
    std::string str();
    std::size_t id() const;
    float weight() const;

private:
    const marisa::Key key_;

    Key();
    Key(const Key &key);
    Key &operator=(const Key &);
};

class Query {
public:
    void str(const char **ptr_out, std::size_t *length_out) const;
    std::size_t id() const;

private:
    const marisa::Query query_;

    Query();
    Query(const Query &query);
    Query &operator=(const Query &);
};

class Keyset {
 friend class Trie;

public:
    Keyset();
    ~Keyset();

    void pushBack(const marisa::Key &key);
    void pushBack(const char *ptr, size_t length, float weight = 1.0);

    const Key &key(std::size_t i) const;

    void keyStr(std::size_t i,
            const char **ptr_out, std::size_t *length_out) const;
    std::size_t keyId(std::size_t i) const;

    std::size_t numKeys() const;

    bool empty() const;
    std::size_t size() const;
    std::size_t totalLength() const;

    void reset();
    void clear();

private:
    marisa::Keyset *keyset_;

    Keyset(const Keyset &);
    Keyset &operator=(const Keyset &);
};

class Agent {
 friend class Trie;

public:
    Agent();
    ~Agent();

    void setQuery(const char *ptr, size_t length);
    void setQuery(std::size_t id);

    const Key &key() const;
    const Query &query() const;

    void keyStr(const char **ptr_out, std::size_t *length_out) const;
    std::size_t keyId() const;

    void queryStr(const char **ptr_out, std::size_t *length_out) const;
    std::size_t queryId() const;

private:
    marisa::Agent *agent_;
    char *buf_;
    std::size_t buf_size_;

    Agent(const Agent &);
    Agent &operator=(const Agent &);
};

class Trie {
public:
    Trie();
    ~Trie();

    void build(Keyset &keyset, int config_flags = 0);

    void mmap(const char *ptr, size_t length);
    void load(const char *ptr, size_t length);
    void save(const char *ptr, size_t length) const;

    bool lookup(Agent &agent) const;
    void reverseLookup(Agent &agent) const;
    bool commonPrefixSearch(Agent &agent) const;
    bool predictiveSearch(Agent &agent) const;

    std::size_t numTries() const;
    std::size_t numKeys() const;
    std::size_t numNodes() const;

    TailMode tailNode() const;
    NodeOrder nodeOrder() const;

    bool empty() const;
    std::size_t size() const;
    std::size_t totalSize() const;
    std::size_t ioSize() const;

    void clear();

private:
    marisa::Trie *trie_;

    Trie(const Trie &);
    Trie &operator=(const Trie &);
};

}

#endif
