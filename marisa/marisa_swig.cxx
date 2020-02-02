#include <new>
#include <string>

#include <cstdio>
#include <cstdlib>
#include <cstring>

#include "marisa_swig.h"

#ifdef _WIN32

#include <sys/types.h>
#include <malloc.h>

size_t strnlen(const char *string, size_t maxlen) {
    const char* end = (const char*)memchr(string, '\0', maxlen);
    return end ? end - string : maxlen;
}

char* strndup(const char* s, size_t n) {
    size_t len = strnlen(s, n);
    char* new_str = (char*)malloc(len + 1);

    if (new_str == NULL)
        return NULL;

    new_str[len] = '\0';
    return (char*)memcpy(new_str, s, len);
}

#endif

namespace marisa_swig {

std::string Key::str() {
    std::string keyS(key_.ptr(), key_.length());
    return keyS;
}

size_t Key::id() const {
    return key_.id();
}

float Key::weight() const {
    return key_.weight();
}

void Query::str(const char **ptr_out, size_t *length_out) const {
    *ptr_out = query_.ptr();
    *length_out = query_.length();
}

size_t Query::id() const {
    return query_.id();
}

Keyset::Keyset() : keyset_(new (std::nothrow) marisa::Keyset) {
    MARISA_THROW_IF(keyset_ == NULL, ::MARISA_MEMORY_ERROR);
}

Keyset::~Keyset() {
    delete keyset_;
}

void Keyset::pushBack(const marisa::Key &key) {
    keyset_->push_back(key);
}

void Keyset::pushBack(const char *ptr, size_t length, float weight) {
    keyset_->push_back(ptr, length, weight);
}

const Key &Keyset::key(size_t i) const {
    return reinterpret_cast<const Key &>((*keyset_)[i]);
}

void Keyset::keyStr(size_t i,
        const char **ptr_out, size_t *length_out) const {
    *ptr_out = (*keyset_)[i].ptr();
    *length_out = (*keyset_)[i].length();
}

size_t Keyset::keyId(size_t i) const {
    return (*keyset_)[i].id();
}

size_t Keyset::numKeys() const {
    return keyset_->num_keys();
}

bool Keyset::empty() const {
    return keyset_->empty();
}

size_t Keyset::size() const {
    return keyset_->size();
}

size_t Keyset::totalLength() const {
    return keyset_->total_length();
}

void Keyset::reset() {
    keyset_->reset();
}

void Keyset::clear() {
    keyset_->clear();
}

Agent::Agent()
        : agent_(new (std::nothrow) marisa::Agent), buf_(NULL), buf_size_(0) {
    MARISA_THROW_IF(agent_ == NULL, ::MARISA_MEMORY_ERROR);
}

Agent::~Agent() {
    delete agent_;
    delete [] buf_;
}

void Agent::setQuery(const char *ptr, size_t length) {
    if (length > buf_size_) {
        size_t new_buf_size = (buf_size_ != 0) ? buf_size_ : 1;
        if (length >= (MARISA_SIZE_MAX / 2)) {
            new_buf_size = MARISA_SIZE_MAX;
        } else {
            while (new_buf_size < length) {
                new_buf_size *= 2;
            }
        }
        char *new_buf = new (std::nothrow) char[new_buf_size];
        MARISA_THROW_IF(new_buf == NULL, MARISA_MEMORY_ERROR);
        delete [] buf_;
        buf_ = new_buf;
        buf_size_ = new_buf_size;
    }
    std::memcpy(buf_, ptr, length);
    agent_->set_query(buf_, length);
}

void Agent::setQuery(size_t id) {
    agent_->set_query(id);
}

const Key &Agent::key() const {
    return reinterpret_cast<const Key &>(agent_->key());
}

const Query &Agent::query() const {
    return reinterpret_cast<const Query &>(agent_->query());
}

void Agent::keyStr(const char **ptr_out, size_t *length_out) const {
    *ptr_out = agent_->key().ptr();
    *length_out = agent_->key().length();
}

size_t Agent::keyId() const {
    return agent_->key().id();
}

void Agent::queryStr(const char **ptr_out, size_t *length_out) const {
    *ptr_out = agent_->query().ptr();
    *length_out = agent_->query().length();
}

size_t Agent::queryId() const {
    return agent_->query().id();
}

Trie::Trie() : trie_(new (std::nothrow) marisa::Trie) {
    MARISA_THROW_IF(trie_ == NULL, ::MARISA_MEMORY_ERROR);
}

Trie::~Trie() {
    delete trie_;
}

void Trie::build(Keyset &keyset, int config_flags) {
    trie_->build(*keyset.keyset_, config_flags);
}

void Trie::mmap(const char *ptr, size_t length) {
    char* filename = strndup(ptr, length);
    trie_->mmap(filename);
    free(filename);
}

void Trie::load(const char *ptr, size_t length) {
    char* filename = strndup(ptr, length);
    trie_->load(filename);
    free(filename);
}

void Trie::save(const char *ptr, size_t length) const {
    char* filename = strndup(ptr, length);
    trie_->save(filename);
    free(filename);
}

bool Trie::lookup(Agent &agent) const {
    return trie_->lookup(*agent.agent_);
}

void Trie::reverseLookup(Agent &agent) const {
    trie_->reverse_lookup(*agent.agent_);
}

bool Trie::commonPrefixSearch(Agent &agent) const {
    return trie_->common_prefix_search(*agent.agent_);
}

bool Trie::predictiveSearch(Agent &agent) const {
    return trie_->predictive_search(*agent.agent_);
}

size_t Trie::numTries() const {
    return trie_->num_tries();
}

size_t Trie::numKeys() const {
    return trie_->num_keys();
}

size_t Trie::numNodes() const {
    return trie_->num_nodes();
}

TailMode Trie::tailNode() const {
    if (trie_->tail_mode() == ::MARISA_TEXT_TAIL) {
        return TEXT_TAIL;
    } else {
        return BINARY_TAIL;
    }
}

NodeOrder Trie::nodeOrder() const {
    if (trie_->node_order() == ::MARISA_LABEL_ORDER) {
        return LABEL_ORDER;
    } else {
        return WEIGHT_ORDER;
    }
}

bool Trie::empty() const {
    return trie_->empty();
}

size_t Trie::size() const {
    return trie_->size();
}

size_t Trie::totalSize() const {
    return trie_->total_size();
}

size_t Trie::ioSize() const {
    return trie_->io_size();
}

void Trie::clear() {
    trie_->clear();
}

}
