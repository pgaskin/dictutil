// AUTOMATICALLY GENERATED, DO NOT EDIT!
// merged from marisa-trie 970b20c.

// ### COPYING
// 
// libmarisa and its command line tools are dual-licensed under the BSD 2-clause license and the LGPL.
// 
// #### The BSD 2-clause license
// 
// Copyright (c) 2010-2019, Susumu Yata
// All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
// 
// - Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
// - Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
// 
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// 
// #### The LGPL 2.1 or any later version
// 
// marisa-trie - A static and space-efficient trie data structure.
// Copyright (C) 2010-2019  Susumu Yata
// 
// This library is free software; you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 2.1 of the License, or (at your option) any later version.
// 
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Lesser General Public License for more details.
// 
// You should have received a copy of the GNU Lesser General Public
// License along with this library; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA
// 

#ifndef MARISA_H_
#define MARISA_H_

// "marisa/stdio.h" includes <cstdio> for I/O using std::FILE.


#ifndef MARISA_MYSTDIO_H_
#define MARISA_MYSTDIO_H_

#include <cstdio>

namespace marisa {

class Trie;

void fread(std::FILE *file, Trie *trie);
void fwrite(std::FILE *file, const Trie &trie);

}  // namespace marisa

#endif  // MARISA_MYSTDIO_H_




// "marisa/iostream.h" includes <iosfwd> for I/O using std::iostream.


#ifndef MARISA_IOSTREAM_H_
#define MARISA_IOSTREAM_H_

#include <iosfwd>

namespace marisa {

class Trie;

std::istream &read(std::istream &stream, Trie *trie);
std::ostream &write(std::ostream &stream, const Trie &trie);

std::istream &operator>>(std::istream &stream, Trie &trie);
std::ostream &operator<<(std::ostream &stream, const Trie &trie);

}  // namespace marisa

#endif  // MARISA_IOSTREAM_H_




// You can use <marisa/trie.h> instead of <marisa.h> if you don't need the
// above I/O interfaces and don't want to include the above I/O headers.


#ifndef MARISA_TRIE_H_
#define MARISA_TRIE_H_


#ifndef MARISA_KEYSET_H_
#define MARISA_KEYSET_H_


#ifndef MARISA_KEY_H_
#define MARISA_KEY_H_


#ifndef MARISA_BASE_H_
#define MARISA_BASE_H_

// Old Visual C++ does not provide stdint.h.
#ifndef _MSC_VER
 #include <stdint.h>
#endif  // _MSC_VER

#ifdef __cplusplus
 #include <cstddef>
#else  // __cplusplus
 #include <stddef.h>
#endif  // __cplusplus

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

#ifdef _MSC_VER
typedef unsigned __int8  marisa_uint8;
typedef unsigned __int16 marisa_uint16;
typedef unsigned __int32 marisa_uint32;
typedef unsigned __int64 marisa_uint64;
#else  // _MSC_VER
typedef uint8_t  marisa_uint8;
typedef uint16_t marisa_uint16;
typedef uint32_t marisa_uint32;
typedef uint64_t marisa_uint64;
#endif  // _MSC_VER

#if defined(_WIN64) || defined(__amd64__) || defined(__x86_64__) || \
    defined(__ia64__) || defined(__ppc64__) || defined(__powerpc64__) || \
    defined(__sparc64__) || defined(__mips64__) || defined(__aarch64__) || \
    defined(__s390x__)
 #define MARISA_WORD_SIZE 64
#else  // defined(_WIN64), etc.
 #define MARISA_WORD_SIZE 32
#endif  // defined(_WIN64), etc.

//#define MARISA_WORD_SIZE  (sizeof(void *) * 8)

#define MARISA_UINT8_MAX  ((marisa_uint8)~(marisa_uint8)0)
#define MARISA_UINT16_MAX ((marisa_uint16)~(marisa_uint16)0)
#define MARISA_UINT32_MAX ((marisa_uint32)~(marisa_uint32)0)
#define MARISA_UINT64_MAX ((marisa_uint64)~(marisa_uint64)0)
#define MARISA_SIZE_MAX   ((size_t)~(size_t)0)

#define MARISA_INVALID_LINK_ID MARISA_UINT32_MAX
#define MARISA_INVALID_KEY_ID  MARISA_UINT32_MAX
#define MARISA_INVALID_EXTRA   (MARISA_UINT32_MAX >> 8)

// Error codes are defined as members of marisa_error_code. This library throws
// an exception with one of the error codes when an error occurs.
typedef enum marisa_error_code_ {
  // MARISA_OK means that a requested operation has succeeded. In practice, an
  // exception never has MARISA_OK because it is not an error.
  MARISA_OK           = 0,

  // MARISA_STATE_ERROR means that an object was not ready for a requested
  // operation. For example, an operation to modify a fixed vector throws an
  // exception with MARISA_STATE_ERROR.
  MARISA_STATE_ERROR  = 1,

  // MARISA_NULL_ERROR means that an invalid NULL pointer has been given.
  MARISA_NULL_ERROR   = 2,

  // MARISA_BOUND_ERROR means that an operation has tried to access an out of
  // range address.
  MARISA_BOUND_ERROR  = 3,

  // MARISA_RANGE_ERROR means that an out of range value has appeared in
  // operation.
  MARISA_RANGE_ERROR  = 4,

  // MARISA_CODE_ERROR means that an undefined code has appeared in operation.
  MARISA_CODE_ERROR   = 5,

  // MARISA_RESET_ERROR means that a smart pointer has tried to reset itself.
  MARISA_RESET_ERROR  = 6,

  // MARISA_SIZE_ERROR means that a size has exceeded a library limitation.
  MARISA_SIZE_ERROR   = 7,

  // MARISA_MEMORY_ERROR means that a memory allocation has failed.
  MARISA_MEMORY_ERROR = 8,

  // MARISA_IO_ERROR means that an I/O operation has failed.
  MARISA_IO_ERROR     = 9,

  // MARISA_FORMAT_ERROR means that input was in invalid format.
  MARISA_FORMAT_ERROR = 10,
} marisa_error_code;

// Min/max values, flags and masks for dictionary settings are defined below.
// Please note that unspecified settings will be replaced with the default
// settings. For example, 0 is equivalent to (MARISA_DEFAULT_NUM_TRIES |
// MARISA_DEFAULT_TRIE | MARISA_DEFAULT_TAIL | MARISA_DEFAULT_ORDER).

// A dictionary consists of 3 tries in default. Usually more tries make a
// dictionary space-efficient but time-inefficient.
typedef enum marisa_num_tries_ {
  MARISA_MIN_NUM_TRIES     = 0x00001,
  MARISA_MAX_NUM_TRIES     = 0x0007F,
  MARISA_DEFAULT_NUM_TRIES = 0x00003,
} marisa_num_tries;

// This library uses a cache technique to accelerate search functions. The
// following enumerated type marisa_cache_level gives a list of available cache
// size options. A larger cache enables faster search but takes a more space.
typedef enum marisa_cache_level_ {
  MARISA_HUGE_CACHE        = 0x00080,
  MARISA_LARGE_CACHE       = 0x00100,
  MARISA_NORMAL_CACHE      = 0x00200,
  MARISA_SMALL_CACHE       = 0x00400,
  MARISA_TINY_CACHE        = 0x00800,
  MARISA_DEFAULT_CACHE     = MARISA_NORMAL_CACHE
} marisa_cache_level;

// This library provides 2 kinds of TAIL implementations.
typedef enum marisa_tail_mode_ {
  // MARISA_TEXT_TAIL merges last labels as zero-terminated strings. So, it is
  // available if and only if the last labels do not contain a NULL character.
  // If MARISA_TEXT_TAIL is specified and a NULL character exists in the last
  // labels, the setting is automatically switched to MARISA_BINARY_TAIL.
  MARISA_TEXT_TAIL         = 0x01000,

  // MARISA_BINARY_TAIL also merges last labels but as byte sequences. It uses
  // a bit vector to detect the end of a sequence, instead of NULL characters.
  // So, MARISA_BINARY_TAIL requires a larger space if the average length of
  // labels is greater than 8.
  MARISA_BINARY_TAIL       = 0x02000,

  MARISA_DEFAULT_TAIL      = MARISA_TEXT_TAIL,
} marisa_tail_mode;

// The arrangement of nodes affects the time cost of matching and the order of
// predictive search.
typedef enum marisa_node_order_ {
  // MARISA_LABEL_ORDER arranges nodes in ascending label order.
  // MARISA_LABEL_ORDER is useful if an application needs to predict keys in
  // label order.
  MARISA_LABEL_ORDER       = 0x10000,

  // MARISA_WEIGHT_ORDER arranges nodes in descending weight order.
  // MARISA_WEIGHT_ORDER is generally a better choice because it enables faster
  // matching.
  MARISA_WEIGHT_ORDER      = 0x20000,

  MARISA_DEFAULT_ORDER     = MARISA_WEIGHT_ORDER,
} marisa_node_order;

typedef enum marisa_config_mask_ {
  MARISA_NUM_TRIES_MASK    = 0x0007F,
  MARISA_CACHE_LEVEL_MASK  = 0x00F80,
  MARISA_TAIL_MODE_MASK    = 0x0F000,
  MARISA_NODE_ORDER_MASK   = 0xF0000,
  MARISA_CONFIG_MASK       = 0xFFFFF
} marisa_config_mask;

#ifdef __cplusplus
}  // extern "C"
#endif  // __cplusplus

#ifdef __cplusplus

// `std::swap` is in <utility> since C++ 11 but in <algorithm> in C++ 98:
#if __cplusplus >= 201103L
 #include <utility>
#else
 #include <algorithm>
#endif
namespace marisa {

typedef ::marisa_uint8  UInt8;
typedef ::marisa_uint16 UInt16;
typedef ::marisa_uint32 UInt32;
typedef ::marisa_uint64 UInt64;

typedef ::marisa_error_code ErrorCode;

typedef ::marisa_cache_level CacheLevel;
typedef ::marisa_tail_mode TailMode;
typedef ::marisa_node_order NodeOrder;

using std::swap;

}  // namespace marisa
#endif  // __cplusplus

#ifdef __cplusplus


#ifndef MARISA_EXCEPTION_H_
#define MARISA_EXCEPTION_H_

#include <exception>


namespace marisa {

// An exception object keeps a filename, a line number, an error code and an
// error message. The message format is as follows:
//  "__FILE__:__LINE__: error_code: error_message"
class Exception : public std::exception {
 public:
  Exception(const char *filename, int line,
      ErrorCode error_code, const char *error_message)
      : std::exception(), filename_(filename), line_(line),
        error_code_(error_code), error_message_(error_message) {}
  Exception(const Exception &ex)
      : std::exception(), filename_(ex.filename_), line_(ex.line_),
        error_code_(ex.error_code_), error_message_(ex.error_message_) {}
  virtual ~Exception() throw() {}

  Exception &operator=(const Exception &rhs) {
    filename_ = rhs.filename_;
    line_ = rhs.line_;
    error_code_ = rhs.error_code_;
    error_message_ = rhs.error_message_;
    return *this;
  }

  const char *filename() const {
    return filename_;
  }
  int line() const {
    return line_;
  }
  ErrorCode error_code() const {
    return error_code_;
  }
  const char *error_message() const {
    return error_message_;
  }

  virtual const char *what() const throw() {
    return error_message_;
  }

 private:
  const char *filename_;
  int line_;
  ErrorCode error_code_;
  const char *error_message_;
};

// These macros are used to convert a line number to a string constant.
#define MARISA_INT_TO_STR(value) #value
#define MARISA_LINE_TO_STR(line) MARISA_INT_TO_STR(line)
#define MARISA_LINE_STR MARISA_LINE_TO_STR(__LINE__)

// MARISA_THROW throws an exception with a filename, a line number, an error
// code and an error message. The message format is as follows:
//  "__FILE__:__LINE__: error_code: error_message"
#define MARISA_THROW(error_code, error_message) \
  (throw marisa::Exception(__FILE__, __LINE__, error_code, \
       __FILE__ ":" MARISA_LINE_STR ": " #error_code ": " error_message))

// MARISA_THROW_IF throws an exception if `condition' is true.
#define MARISA_THROW_IF(condition, error_code) \
  (void)((!(condition)) || (MARISA_THROW(error_code, #condition), 0))

// MARISA_DEBUG_IF is ignored if _DEBUG is undefined. So, it is useful for
// debugging time-critical codes.
#ifdef _DEBUG
 #define MARISA_DEBUG_IF(cond, error_code) MARISA_THROW_IF(cond, error_code)
#else
 #define MARISA_DEBUG_IF(cond, error_code)
#endif

}  // namespace marisa

#endif  // MARISA_EXCEPTION_H_





#ifndef MARISA_SCOPED_PTR_H_
#define MARISA_SCOPED_PTR_H_


namespace marisa {

template <typename T>
class scoped_ptr {
 public:
  scoped_ptr() : ptr_(NULL) {}
  explicit scoped_ptr(T *ptr) : ptr_(ptr) {}

  ~scoped_ptr() {
    delete ptr_;
  }

  void reset(T *ptr = NULL) {
    MARISA_THROW_IF((ptr != NULL) && (ptr == ptr_), MARISA_RESET_ERROR);
    scoped_ptr(ptr).swap(*this);
  }

  T &operator*() const {
    MARISA_DEBUG_IF(ptr_ == NULL, MARISA_STATE_ERROR);
    return *ptr_;
  }
  T *operator->() const {
    MARISA_DEBUG_IF(ptr_ == NULL, MARISA_STATE_ERROR);
    return ptr_;
  }
  T *get() const {
    return ptr_;
  }

  void clear() {
    scoped_ptr().swap(*this);
  }
  void swap(scoped_ptr &rhs) {
    marisa::swap(ptr_, rhs.ptr_);
  }

 private:
  T *ptr_;

  // Disallows copy and assignment.
  scoped_ptr(const scoped_ptr &);
  scoped_ptr &operator=(const scoped_ptr &);
};

}  // namespace marisa

#endif  // MARISA_SCOPED_PTR_H_





#ifndef MARISA_SCOPED_ARRAY_H_
#define MARISA_SCOPED_ARRAY_H_


namespace marisa {

template <typename T>
class scoped_array {
 public:
  scoped_array() : array_(NULL) {}
  explicit scoped_array(T *array) : array_(array) {}

  ~scoped_array() {
    delete [] array_;
  }

  void reset(T *array = NULL) {
    MARISA_THROW_IF((array != NULL) && (array == array_), MARISA_RESET_ERROR);
    scoped_array(array).swap(*this);
  }

  T &operator[](std::size_t i) const {
    MARISA_DEBUG_IF(array_ == NULL, MARISA_STATE_ERROR);
    return array_[i];
  }
  T *get() const {
    return array_;
  }

  void clear() {
    scoped_array().swap(*this);
  }
  void swap(scoped_array &rhs) {
    marisa::swap(array_, rhs.array_);
  }

 private:
  T *array_;

  // Disallows copy and assignment.
  scoped_array(const scoped_array &);
  scoped_array &operator=(const scoped_array &);
};

}  // namespace marisa

#endif  // MARISA_SCOPED_ARRAY_H_



#endif  // __cplusplus

#endif  // MARISA_BASE_H_




namespace marisa {

class Key {
 public:
  Key() : ptr_(NULL), length_(0), union_() {
    union_.id = 0;
  }
  Key(const Key &key)
      : ptr_(key.ptr_), length_(key.length_), union_(key.union_) {}

  Key &operator=(const Key &key) {
    ptr_ = key.ptr_;
    length_ = key.length_;
    union_ = key.union_;
    return *this;
  }

  char operator[](std::size_t i) const {
    MARISA_DEBUG_IF(i >= length_, MARISA_BOUND_ERROR);
    return ptr_[i];
  }

  void set_str(const char *str) {
    MARISA_DEBUG_IF(str == NULL, MARISA_NULL_ERROR);
    std::size_t length = 0;
    while (str[length] != '\0') {
      ++length;
    }
    MARISA_DEBUG_IF(length > MARISA_UINT32_MAX, MARISA_SIZE_ERROR);
    ptr_ = str;
    length_ = (UInt32)length;
  }
  void set_str(const char *ptr, std::size_t length) {
    MARISA_DEBUG_IF((ptr == NULL) && (length != 0), MARISA_NULL_ERROR);
    MARISA_DEBUG_IF(length > MARISA_UINT32_MAX, MARISA_SIZE_ERROR);
    ptr_ = ptr;
    length_ = (UInt32)length;
  }
  void set_id(std::size_t id) {
    MARISA_DEBUG_IF(id > MARISA_UINT32_MAX, MARISA_SIZE_ERROR);
    union_.id = (UInt32)id;
  }
  void set_weight(float weight) {
    union_.weight = weight;
  }

  const char *ptr() const {
    return ptr_;
  }
  std::size_t length() const {
    return length_;
  }
  std::size_t id() const {
    return union_.id;
  }
  float weight() const {
    return union_.weight;
  }

  void clear() {
    Key().swap(*this);
  }
  void swap(Key &rhs) {
    marisa::swap(ptr_, rhs.ptr_);
    marisa::swap(length_, rhs.length_);
    marisa::swap(union_.id, rhs.union_.id);
  }

 private:
  const char *ptr_;
  UInt32 length_;
  union Union {
    UInt32 id;
    float weight;
  } union_;
};

}  // namespace marisa

#endif  // MARISA_KEY_H_




namespace marisa {

class Keyset {
 public:
  enum {
    BASE_BLOCK_SIZE  = 4096,
    EXTRA_BLOCK_SIZE = 1024,
    KEY_BLOCK_SIZE   = 256
  };

  Keyset();

  void push_back(const Key &key);
  void push_back(const Key &key, char end_marker);

  void push_back(const char *str);
  void push_back(const char *ptr, std::size_t length, float weight = 1.0);

  const Key &operator[](std::size_t i) const {
    MARISA_DEBUG_IF(i >= size_, MARISA_BOUND_ERROR);
    return key_blocks_[i / KEY_BLOCK_SIZE][i % KEY_BLOCK_SIZE];
  }
  Key &operator[](std::size_t i) {
    MARISA_DEBUG_IF(i >= size_, MARISA_BOUND_ERROR);
    return key_blocks_[i / KEY_BLOCK_SIZE][i % KEY_BLOCK_SIZE];
  }

  std::size_t num_keys() const {
    return size_;
  }

  bool empty() const {
    return size_ == 0;
  }
  std::size_t size() const {
    return size_;
  }
  std::size_t total_length() const {
    return total_length_;
  }

  void reset();

  void clear();
  void swap(Keyset &rhs);

 private:
  scoped_array<scoped_array<char> > base_blocks_;
  std::size_t base_blocks_size_;
  std::size_t base_blocks_capacity_;
  scoped_array<scoped_array<char> > extra_blocks_;
  std::size_t extra_blocks_size_;
  std::size_t extra_blocks_capacity_;
  scoped_array<scoped_array<Key> > key_blocks_;
  std::size_t key_blocks_size_;
  std::size_t key_blocks_capacity_;
  char *ptr_;
  std::size_t avail_;
  std::size_t size_;
  std::size_t total_length_;

  char *reserve(std::size_t size);

  void append_base_block();
  void append_extra_block(std::size_t size);
  void append_key_block();

  // Disallows copy and assignment.
  Keyset(const Keyset &);
  Keyset &operator=(const Keyset &);
};

}  // namespace marisa

#endif  // MARISA_KEYSET_H_





#ifndef MARISA_AGENT_H_
#define MARISA_AGENT_H_


#ifndef MARISA_KEY_H_
#define MARISA_KEY_H_


#ifndef MARISA_BASE_H_
#define MARISA_BASE_H_

// Old Visual C++ does not provide stdint.h.
#ifndef _MSC_VER
 #include <stdint.h>
#endif  // _MSC_VER

#ifdef __cplusplus
 #include <cstddef>
#else  // __cplusplus
 #include <stddef.h>
#endif  // __cplusplus

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

#ifdef _MSC_VER
typedef unsigned __int8  marisa_uint8;
typedef unsigned __int16 marisa_uint16;
typedef unsigned __int32 marisa_uint32;
typedef unsigned __int64 marisa_uint64;
#else  // _MSC_VER
typedef uint8_t  marisa_uint8;
typedef uint16_t marisa_uint16;
typedef uint32_t marisa_uint32;
typedef uint64_t marisa_uint64;
#endif  // _MSC_VER

#if defined(_WIN64) || defined(__amd64__) || defined(__x86_64__) || \
    defined(__ia64__) || defined(__ppc64__) || defined(__powerpc64__) || \
    defined(__sparc64__) || defined(__mips64__) || defined(__aarch64__) || \
    defined(__s390x__)
 #define MARISA_WORD_SIZE 64
#else  // defined(_WIN64), etc.
 #define MARISA_WORD_SIZE 32
#endif  // defined(_WIN64), etc.

//#define MARISA_WORD_SIZE  (sizeof(void *) * 8)

#define MARISA_UINT8_MAX  ((marisa_uint8)~(marisa_uint8)0)
#define MARISA_UINT16_MAX ((marisa_uint16)~(marisa_uint16)0)
#define MARISA_UINT32_MAX ((marisa_uint32)~(marisa_uint32)0)
#define MARISA_UINT64_MAX ((marisa_uint64)~(marisa_uint64)0)
#define MARISA_SIZE_MAX   ((size_t)~(size_t)0)

#define MARISA_INVALID_LINK_ID MARISA_UINT32_MAX
#define MARISA_INVALID_KEY_ID  MARISA_UINT32_MAX
#define MARISA_INVALID_EXTRA   (MARISA_UINT32_MAX >> 8)

// Error codes are defined as members of marisa_error_code. This library throws
// an exception with one of the error codes when an error occurs.
typedef enum marisa_error_code_ {
  // MARISA_OK means that a requested operation has succeeded. In practice, an
  // exception never has MARISA_OK because it is not an error.
  MARISA_OK           = 0,

  // MARISA_STATE_ERROR means that an object was not ready for a requested
  // operation. For example, an operation to modify a fixed vector throws an
  // exception with MARISA_STATE_ERROR.
  MARISA_STATE_ERROR  = 1,

  // MARISA_NULL_ERROR means that an invalid NULL pointer has been given.
  MARISA_NULL_ERROR   = 2,

  // MARISA_BOUND_ERROR means that an operation has tried to access an out of
  // range address.
  MARISA_BOUND_ERROR  = 3,

  // MARISA_RANGE_ERROR means that an out of range value has appeared in
  // operation.
  MARISA_RANGE_ERROR  = 4,

  // MARISA_CODE_ERROR means that an undefined code has appeared in operation.
  MARISA_CODE_ERROR   = 5,

  // MARISA_RESET_ERROR means that a smart pointer has tried to reset itself.
  MARISA_RESET_ERROR  = 6,

  // MARISA_SIZE_ERROR means that a size has exceeded a library limitation.
  MARISA_SIZE_ERROR   = 7,

  // MARISA_MEMORY_ERROR means that a memory allocation has failed.
  MARISA_MEMORY_ERROR = 8,

  // MARISA_IO_ERROR means that an I/O operation has failed.
  MARISA_IO_ERROR     = 9,

  // MARISA_FORMAT_ERROR means that input was in invalid format.
  MARISA_FORMAT_ERROR = 10,
} marisa_error_code;

// Min/max values, flags and masks for dictionary settings are defined below.
// Please note that unspecified settings will be replaced with the default
// settings. For example, 0 is equivalent to (MARISA_DEFAULT_NUM_TRIES |
// MARISA_DEFAULT_TRIE | MARISA_DEFAULT_TAIL | MARISA_DEFAULT_ORDER).

// A dictionary consists of 3 tries in default. Usually more tries make a
// dictionary space-efficient but time-inefficient.
typedef enum marisa_num_tries_ {
  MARISA_MIN_NUM_TRIES     = 0x00001,
  MARISA_MAX_NUM_TRIES     = 0x0007F,
  MARISA_DEFAULT_NUM_TRIES = 0x00003,
} marisa_num_tries;

// This library uses a cache technique to accelerate search functions. The
// following enumerated type marisa_cache_level gives a list of available cache
// size options. A larger cache enables faster search but takes a more space.
typedef enum marisa_cache_level_ {
  MARISA_HUGE_CACHE        = 0x00080,
  MARISA_LARGE_CACHE       = 0x00100,
  MARISA_NORMAL_CACHE      = 0x00200,
  MARISA_SMALL_CACHE       = 0x00400,
  MARISA_TINY_CACHE        = 0x00800,
  MARISA_DEFAULT_CACHE     = MARISA_NORMAL_CACHE
} marisa_cache_level;

// This library provides 2 kinds of TAIL implementations.
typedef enum marisa_tail_mode_ {
  // MARISA_TEXT_TAIL merges last labels as zero-terminated strings. So, it is
  // available if and only if the last labels do not contain a NULL character.
  // If MARISA_TEXT_TAIL is specified and a NULL character exists in the last
  // labels, the setting is automatically switched to MARISA_BINARY_TAIL.
  MARISA_TEXT_TAIL         = 0x01000,

  // MARISA_BINARY_TAIL also merges last labels but as byte sequences. It uses
  // a bit vector to detect the end of a sequence, instead of NULL characters.
  // So, MARISA_BINARY_TAIL requires a larger space if the average length of
  // labels is greater than 8.
  MARISA_BINARY_TAIL       = 0x02000,

  MARISA_DEFAULT_TAIL      = MARISA_TEXT_TAIL,
} marisa_tail_mode;

// The arrangement of nodes affects the time cost of matching and the order of
// predictive search.
typedef enum marisa_node_order_ {
  // MARISA_LABEL_ORDER arranges nodes in ascending label order.
  // MARISA_LABEL_ORDER is useful if an application needs to predict keys in
  // label order.
  MARISA_LABEL_ORDER       = 0x10000,

  // MARISA_WEIGHT_ORDER arranges nodes in descending weight order.
  // MARISA_WEIGHT_ORDER is generally a better choice because it enables faster
  // matching.
  MARISA_WEIGHT_ORDER      = 0x20000,

  MARISA_DEFAULT_ORDER     = MARISA_WEIGHT_ORDER,
} marisa_node_order;

typedef enum marisa_config_mask_ {
  MARISA_NUM_TRIES_MASK    = 0x0007F,
  MARISA_CACHE_LEVEL_MASK  = 0x00F80,
  MARISA_TAIL_MODE_MASK    = 0x0F000,
  MARISA_NODE_ORDER_MASK   = 0xF0000,
  MARISA_CONFIG_MASK       = 0xFFFFF
} marisa_config_mask;

#ifdef __cplusplus
}  // extern "C"
#endif  // __cplusplus

#ifdef __cplusplus

// `std::swap` is in <utility> since C++ 11 but in <algorithm> in C++ 98:
#if __cplusplus >= 201103L
 #include <utility>
#else
 #include <algorithm>
#endif
namespace marisa {

typedef ::marisa_uint8  UInt8;
typedef ::marisa_uint16 UInt16;
typedef ::marisa_uint32 UInt32;
typedef ::marisa_uint64 UInt64;

typedef ::marisa_error_code ErrorCode;

typedef ::marisa_cache_level CacheLevel;
typedef ::marisa_tail_mode TailMode;
typedef ::marisa_node_order NodeOrder;

using std::swap;

}  // namespace marisa
#endif  // __cplusplus

#ifdef __cplusplus


#ifndef MARISA_EXCEPTION_H_
#define MARISA_EXCEPTION_H_

#include <exception>


namespace marisa {

// An exception object keeps a filename, a line number, an error code and an
// error message. The message format is as follows:
//  "__FILE__:__LINE__: error_code: error_message"
class Exception : public std::exception {
 public:
  Exception(const char *filename, int line,
      ErrorCode error_code, const char *error_message)
      : std::exception(), filename_(filename), line_(line),
        error_code_(error_code), error_message_(error_message) {}
  Exception(const Exception &ex)
      : std::exception(), filename_(ex.filename_), line_(ex.line_),
        error_code_(ex.error_code_), error_message_(ex.error_message_) {}
  virtual ~Exception() throw() {}

  Exception &operator=(const Exception &rhs) {
    filename_ = rhs.filename_;
    line_ = rhs.line_;
    error_code_ = rhs.error_code_;
    error_message_ = rhs.error_message_;
    return *this;
  }

  const char *filename() const {
    return filename_;
  }
  int line() const {
    return line_;
  }
  ErrorCode error_code() const {
    return error_code_;
  }
  const char *error_message() const {
    return error_message_;
  }

  virtual const char *what() const throw() {
    return error_message_;
  }

 private:
  const char *filename_;
  int line_;
  ErrorCode error_code_;
  const char *error_message_;
};

// These macros are used to convert a line number to a string constant.
#define MARISA_INT_TO_STR(value) #value
#define MARISA_LINE_TO_STR(line) MARISA_INT_TO_STR(line)
#define MARISA_LINE_STR MARISA_LINE_TO_STR(__LINE__)

// MARISA_THROW throws an exception with a filename, a line number, an error
// code and an error message. The message format is as follows:
//  "__FILE__:__LINE__: error_code: error_message"
#define MARISA_THROW(error_code, error_message) \
  (throw marisa::Exception(__FILE__, __LINE__, error_code, \
       __FILE__ ":" MARISA_LINE_STR ": " #error_code ": " error_message))

// MARISA_THROW_IF throws an exception if `condition' is true.
#define MARISA_THROW_IF(condition, error_code) \
  (void)((!(condition)) || (MARISA_THROW(error_code, #condition), 0))

// MARISA_DEBUG_IF is ignored if _DEBUG is undefined. So, it is useful for
// debugging time-critical codes.
#ifdef _DEBUG
 #define MARISA_DEBUG_IF(cond, error_code) MARISA_THROW_IF(cond, error_code)
#else
 #define MARISA_DEBUG_IF(cond, error_code)
#endif

}  // namespace marisa

#endif  // MARISA_EXCEPTION_H_





#ifndef MARISA_SCOPED_PTR_H_
#define MARISA_SCOPED_PTR_H_


namespace marisa {

template <typename T>
class scoped_ptr {
 public:
  scoped_ptr() : ptr_(NULL) {}
  explicit scoped_ptr(T *ptr) : ptr_(ptr) {}

  ~scoped_ptr() {
    delete ptr_;
  }

  void reset(T *ptr = NULL) {
    MARISA_THROW_IF((ptr != NULL) && (ptr == ptr_), MARISA_RESET_ERROR);
    scoped_ptr(ptr).swap(*this);
  }

  T &operator*() const {
    MARISA_DEBUG_IF(ptr_ == NULL, MARISA_STATE_ERROR);
    return *ptr_;
  }
  T *operator->() const {
    MARISA_DEBUG_IF(ptr_ == NULL, MARISA_STATE_ERROR);
    return ptr_;
  }
  T *get() const {
    return ptr_;
  }

  void clear() {
    scoped_ptr().swap(*this);
  }
  void swap(scoped_ptr &rhs) {
    marisa::swap(ptr_, rhs.ptr_);
  }

 private:
  T *ptr_;

  // Disallows copy and assignment.
  scoped_ptr(const scoped_ptr &);
  scoped_ptr &operator=(const scoped_ptr &);
};

}  // namespace marisa

#endif  // MARISA_SCOPED_PTR_H_





#ifndef MARISA_SCOPED_ARRAY_H_
#define MARISA_SCOPED_ARRAY_H_


namespace marisa {

template <typename T>
class scoped_array {
 public:
  scoped_array() : array_(NULL) {}
  explicit scoped_array(T *array) : array_(array) {}

  ~scoped_array() {
    delete [] array_;
  }

  void reset(T *array = NULL) {
    MARISA_THROW_IF((array != NULL) && (array == array_), MARISA_RESET_ERROR);
    scoped_array(array).swap(*this);
  }

  T &operator[](std::size_t i) const {
    MARISA_DEBUG_IF(array_ == NULL, MARISA_STATE_ERROR);
    return array_[i];
  }
  T *get() const {
    return array_;
  }

  void clear() {
    scoped_array().swap(*this);
  }
  void swap(scoped_array &rhs) {
    marisa::swap(array_, rhs.array_);
  }

 private:
  T *array_;

  // Disallows copy and assignment.
  scoped_array(const scoped_array &);
  scoped_array &operator=(const scoped_array &);
};

}  // namespace marisa

#endif  // MARISA_SCOPED_ARRAY_H_



#endif  // __cplusplus

#endif  // MARISA_BASE_H_




namespace marisa {

class Key {
 public:
  Key() : ptr_(NULL), length_(0), union_() {
    union_.id = 0;
  }
  Key(const Key &key)
      : ptr_(key.ptr_), length_(key.length_), union_(key.union_) {}

  Key &operator=(const Key &key) {
    ptr_ = key.ptr_;
    length_ = key.length_;
    union_ = key.union_;
    return *this;
  }

  char operator[](std::size_t i) const {
    MARISA_DEBUG_IF(i >= length_, MARISA_BOUND_ERROR);
    return ptr_[i];
  }

  void set_str(const char *str) {
    MARISA_DEBUG_IF(str == NULL, MARISA_NULL_ERROR);
    std::size_t length = 0;
    while (str[length] != '\0') {
      ++length;
    }
    MARISA_DEBUG_IF(length > MARISA_UINT32_MAX, MARISA_SIZE_ERROR);
    ptr_ = str;
    length_ = (UInt32)length;
  }
  void set_str(const char *ptr, std::size_t length) {
    MARISA_DEBUG_IF((ptr == NULL) && (length != 0), MARISA_NULL_ERROR);
    MARISA_DEBUG_IF(length > MARISA_UINT32_MAX, MARISA_SIZE_ERROR);
    ptr_ = ptr;
    length_ = (UInt32)length;
  }
  void set_id(std::size_t id) {
    MARISA_DEBUG_IF(id > MARISA_UINT32_MAX, MARISA_SIZE_ERROR);
    union_.id = (UInt32)id;
  }
  void set_weight(float weight) {
    union_.weight = weight;
  }

  const char *ptr() const {
    return ptr_;
  }
  std::size_t length() const {
    return length_;
  }
  std::size_t id() const {
    return union_.id;
  }
  float weight() const {
    return union_.weight;
  }

  void clear() {
    Key().swap(*this);
  }
  void swap(Key &rhs) {
    marisa::swap(ptr_, rhs.ptr_);
    marisa::swap(length_, rhs.length_);
    marisa::swap(union_.id, rhs.union_.id);
  }

 private:
  const char *ptr_;
  UInt32 length_;
  union Union {
    UInt32 id;
    float weight;
  } union_;
};

}  // namespace marisa

#endif  // MARISA_KEY_H_





#ifndef MARISA_QUERY_H_
#define MARISA_QUERY_H_


#ifndef MARISA_BASE_H_
#define MARISA_BASE_H_

// Old Visual C++ does not provide stdint.h.
#ifndef _MSC_VER
 #include <stdint.h>
#endif  // _MSC_VER

#ifdef __cplusplus
 #include <cstddef>
#else  // __cplusplus
 #include <stddef.h>
#endif  // __cplusplus

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

#ifdef _MSC_VER
typedef unsigned __int8  marisa_uint8;
typedef unsigned __int16 marisa_uint16;
typedef unsigned __int32 marisa_uint32;
typedef unsigned __int64 marisa_uint64;
#else  // _MSC_VER
typedef uint8_t  marisa_uint8;
typedef uint16_t marisa_uint16;
typedef uint32_t marisa_uint32;
typedef uint64_t marisa_uint64;
#endif  // _MSC_VER

#if defined(_WIN64) || defined(__amd64__) || defined(__x86_64__) || \
    defined(__ia64__) || defined(__ppc64__) || defined(__powerpc64__) || \
    defined(__sparc64__) || defined(__mips64__) || defined(__aarch64__) || \
    defined(__s390x__)
 #define MARISA_WORD_SIZE 64
#else  // defined(_WIN64), etc.
 #define MARISA_WORD_SIZE 32
#endif  // defined(_WIN64), etc.

//#define MARISA_WORD_SIZE  (sizeof(void *) * 8)

#define MARISA_UINT8_MAX  ((marisa_uint8)~(marisa_uint8)0)
#define MARISA_UINT16_MAX ((marisa_uint16)~(marisa_uint16)0)
#define MARISA_UINT32_MAX ((marisa_uint32)~(marisa_uint32)0)
#define MARISA_UINT64_MAX ((marisa_uint64)~(marisa_uint64)0)
#define MARISA_SIZE_MAX   ((size_t)~(size_t)0)

#define MARISA_INVALID_LINK_ID MARISA_UINT32_MAX
#define MARISA_INVALID_KEY_ID  MARISA_UINT32_MAX
#define MARISA_INVALID_EXTRA   (MARISA_UINT32_MAX >> 8)

// Error codes are defined as members of marisa_error_code. This library throws
// an exception with one of the error codes when an error occurs.
typedef enum marisa_error_code_ {
  // MARISA_OK means that a requested operation has succeeded. In practice, an
  // exception never has MARISA_OK because it is not an error.
  MARISA_OK           = 0,

  // MARISA_STATE_ERROR means that an object was not ready for a requested
  // operation. For example, an operation to modify a fixed vector throws an
  // exception with MARISA_STATE_ERROR.
  MARISA_STATE_ERROR  = 1,

  // MARISA_NULL_ERROR means that an invalid NULL pointer has been given.
  MARISA_NULL_ERROR   = 2,

  // MARISA_BOUND_ERROR means that an operation has tried to access an out of
  // range address.
  MARISA_BOUND_ERROR  = 3,

  // MARISA_RANGE_ERROR means that an out of range value has appeared in
  // operation.
  MARISA_RANGE_ERROR  = 4,

  // MARISA_CODE_ERROR means that an undefined code has appeared in operation.
  MARISA_CODE_ERROR   = 5,

  // MARISA_RESET_ERROR means that a smart pointer has tried to reset itself.
  MARISA_RESET_ERROR  = 6,

  // MARISA_SIZE_ERROR means that a size has exceeded a library limitation.
  MARISA_SIZE_ERROR   = 7,

  // MARISA_MEMORY_ERROR means that a memory allocation has failed.
  MARISA_MEMORY_ERROR = 8,

  // MARISA_IO_ERROR means that an I/O operation has failed.
  MARISA_IO_ERROR     = 9,

  // MARISA_FORMAT_ERROR means that input was in invalid format.
  MARISA_FORMAT_ERROR = 10,
} marisa_error_code;

// Min/max values, flags and masks for dictionary settings are defined below.
// Please note that unspecified settings will be replaced with the default
// settings. For example, 0 is equivalent to (MARISA_DEFAULT_NUM_TRIES |
// MARISA_DEFAULT_TRIE | MARISA_DEFAULT_TAIL | MARISA_DEFAULT_ORDER).

// A dictionary consists of 3 tries in default. Usually more tries make a
// dictionary space-efficient but time-inefficient.
typedef enum marisa_num_tries_ {
  MARISA_MIN_NUM_TRIES     = 0x00001,
  MARISA_MAX_NUM_TRIES     = 0x0007F,
  MARISA_DEFAULT_NUM_TRIES = 0x00003,
} marisa_num_tries;

// This library uses a cache technique to accelerate search functions. The
// following enumerated type marisa_cache_level gives a list of available cache
// size options. A larger cache enables faster search but takes a more space.
typedef enum marisa_cache_level_ {
  MARISA_HUGE_CACHE        = 0x00080,
  MARISA_LARGE_CACHE       = 0x00100,
  MARISA_NORMAL_CACHE      = 0x00200,
  MARISA_SMALL_CACHE       = 0x00400,
  MARISA_TINY_CACHE        = 0x00800,
  MARISA_DEFAULT_CACHE     = MARISA_NORMAL_CACHE
} marisa_cache_level;

// This library provides 2 kinds of TAIL implementations.
typedef enum marisa_tail_mode_ {
  // MARISA_TEXT_TAIL merges last labels as zero-terminated strings. So, it is
  // available if and only if the last labels do not contain a NULL character.
  // If MARISA_TEXT_TAIL is specified and a NULL character exists in the last
  // labels, the setting is automatically switched to MARISA_BINARY_TAIL.
  MARISA_TEXT_TAIL         = 0x01000,

  // MARISA_BINARY_TAIL also merges last labels but as byte sequences. It uses
  // a bit vector to detect the end of a sequence, instead of NULL characters.
  // So, MARISA_BINARY_TAIL requires a larger space if the average length of
  // labels is greater than 8.
  MARISA_BINARY_TAIL       = 0x02000,

  MARISA_DEFAULT_TAIL      = MARISA_TEXT_TAIL,
} marisa_tail_mode;

// The arrangement of nodes affects the time cost of matching and the order of
// predictive search.
typedef enum marisa_node_order_ {
  // MARISA_LABEL_ORDER arranges nodes in ascending label order.
  // MARISA_LABEL_ORDER is useful if an application needs to predict keys in
  // label order.
  MARISA_LABEL_ORDER       = 0x10000,

  // MARISA_WEIGHT_ORDER arranges nodes in descending weight order.
  // MARISA_WEIGHT_ORDER is generally a better choice because it enables faster
  // matching.
  MARISA_WEIGHT_ORDER      = 0x20000,

  MARISA_DEFAULT_ORDER     = MARISA_WEIGHT_ORDER,
} marisa_node_order;

typedef enum marisa_config_mask_ {
  MARISA_NUM_TRIES_MASK    = 0x0007F,
  MARISA_CACHE_LEVEL_MASK  = 0x00F80,
  MARISA_TAIL_MODE_MASK    = 0x0F000,
  MARISA_NODE_ORDER_MASK   = 0xF0000,
  MARISA_CONFIG_MASK       = 0xFFFFF
} marisa_config_mask;

#ifdef __cplusplus
}  // extern "C"
#endif  // __cplusplus

#ifdef __cplusplus

// `std::swap` is in <utility> since C++ 11 but in <algorithm> in C++ 98:
#if __cplusplus >= 201103L
 #include <utility>
#else
 #include <algorithm>
#endif
namespace marisa {

typedef ::marisa_uint8  UInt8;
typedef ::marisa_uint16 UInt16;
typedef ::marisa_uint32 UInt32;
typedef ::marisa_uint64 UInt64;

typedef ::marisa_error_code ErrorCode;

typedef ::marisa_cache_level CacheLevel;
typedef ::marisa_tail_mode TailMode;
typedef ::marisa_node_order NodeOrder;

using std::swap;

}  // namespace marisa
#endif  // __cplusplus

#ifdef __cplusplus


#ifndef MARISA_EXCEPTION_H_
#define MARISA_EXCEPTION_H_

#include <exception>


namespace marisa {

// An exception object keeps a filename, a line number, an error code and an
// error message. The message format is as follows:
//  "__FILE__:__LINE__: error_code: error_message"
class Exception : public std::exception {
 public:
  Exception(const char *filename, int line,
      ErrorCode error_code, const char *error_message)
      : std::exception(), filename_(filename), line_(line),
        error_code_(error_code), error_message_(error_message) {}
  Exception(const Exception &ex)
      : std::exception(), filename_(ex.filename_), line_(ex.line_),
        error_code_(ex.error_code_), error_message_(ex.error_message_) {}
  virtual ~Exception() throw() {}

  Exception &operator=(const Exception &rhs) {
    filename_ = rhs.filename_;
    line_ = rhs.line_;
    error_code_ = rhs.error_code_;
    error_message_ = rhs.error_message_;
    return *this;
  }

  const char *filename() const {
    return filename_;
  }
  int line() const {
    return line_;
  }
  ErrorCode error_code() const {
    return error_code_;
  }
  const char *error_message() const {
    return error_message_;
  }

  virtual const char *what() const throw() {
    return error_message_;
  }

 private:
  const char *filename_;
  int line_;
  ErrorCode error_code_;
  const char *error_message_;
};

// These macros are used to convert a line number to a string constant.
#define MARISA_INT_TO_STR(value) #value
#define MARISA_LINE_TO_STR(line) MARISA_INT_TO_STR(line)
#define MARISA_LINE_STR MARISA_LINE_TO_STR(__LINE__)

// MARISA_THROW throws an exception with a filename, a line number, an error
// code and an error message. The message format is as follows:
//  "__FILE__:__LINE__: error_code: error_message"
#define MARISA_THROW(error_code, error_message) \
  (throw marisa::Exception(__FILE__, __LINE__, error_code, \
       __FILE__ ":" MARISA_LINE_STR ": " #error_code ": " error_message))

// MARISA_THROW_IF throws an exception if `condition' is true.
#define MARISA_THROW_IF(condition, error_code) \
  (void)((!(condition)) || (MARISA_THROW(error_code, #condition), 0))

// MARISA_DEBUG_IF is ignored if _DEBUG is undefined. So, it is useful for
// debugging time-critical codes.
#ifdef _DEBUG
 #define MARISA_DEBUG_IF(cond, error_code) MARISA_THROW_IF(cond, error_code)
#else
 #define MARISA_DEBUG_IF(cond, error_code)
#endif

}  // namespace marisa

#endif  // MARISA_EXCEPTION_H_





#ifndef MARISA_SCOPED_PTR_H_
#define MARISA_SCOPED_PTR_H_


namespace marisa {

template <typename T>
class scoped_ptr {
 public:
  scoped_ptr() : ptr_(NULL) {}
  explicit scoped_ptr(T *ptr) : ptr_(ptr) {}

  ~scoped_ptr() {
    delete ptr_;
  }

  void reset(T *ptr = NULL) {
    MARISA_THROW_IF((ptr != NULL) && (ptr == ptr_), MARISA_RESET_ERROR);
    scoped_ptr(ptr).swap(*this);
  }

  T &operator*() const {
    MARISA_DEBUG_IF(ptr_ == NULL, MARISA_STATE_ERROR);
    return *ptr_;
  }
  T *operator->() const {
    MARISA_DEBUG_IF(ptr_ == NULL, MARISA_STATE_ERROR);
    return ptr_;
  }
  T *get() const {
    return ptr_;
  }

  void clear() {
    scoped_ptr().swap(*this);
  }
  void swap(scoped_ptr &rhs) {
    marisa::swap(ptr_, rhs.ptr_);
  }

 private:
  T *ptr_;

  // Disallows copy and assignment.
  scoped_ptr(const scoped_ptr &);
  scoped_ptr &operator=(const scoped_ptr &);
};

}  // namespace marisa

#endif  // MARISA_SCOPED_PTR_H_





#ifndef MARISA_SCOPED_ARRAY_H_
#define MARISA_SCOPED_ARRAY_H_


namespace marisa {

template <typename T>
class scoped_array {
 public:
  scoped_array() : array_(NULL) {}
  explicit scoped_array(T *array) : array_(array) {}

  ~scoped_array() {
    delete [] array_;
  }

  void reset(T *array = NULL) {
    MARISA_THROW_IF((array != NULL) && (array == array_), MARISA_RESET_ERROR);
    scoped_array(array).swap(*this);
  }

  T &operator[](std::size_t i) const {
    MARISA_DEBUG_IF(array_ == NULL, MARISA_STATE_ERROR);
    return array_[i];
  }
  T *get() const {
    return array_;
  }

  void clear() {
    scoped_array().swap(*this);
  }
  void swap(scoped_array &rhs) {
    marisa::swap(array_, rhs.array_);
  }

 private:
  T *array_;

  // Disallows copy and assignment.
  scoped_array(const scoped_array &);
  scoped_array &operator=(const scoped_array &);
};

}  // namespace marisa

#endif  // MARISA_SCOPED_ARRAY_H_



#endif  // __cplusplus

#endif  // MARISA_BASE_H_




namespace marisa {

class Query {
 public:
  Query() : ptr_(NULL), length_(0), id_(0) {}
  Query(const Query &query)
      : ptr_(query.ptr_), length_(query.length_), id_(query.id_) {}

  Query &operator=(const Query &query) {
    ptr_ = query.ptr_;
    length_ = query.length_;
    id_ = query.id_;
    return *this;
  }

  char operator[](std::size_t i) const {
    MARISA_DEBUG_IF(i >= length_, MARISA_BOUND_ERROR);
    return ptr_[i];
  }

  void set_str(const char *str) {
    MARISA_DEBUG_IF(str == NULL, MARISA_NULL_ERROR);
    std::size_t length = 0;
    while (str[length] != '\0') {
      ++length;
    }
    ptr_ = str;
    length_ = length;
  }
  void set_str(const char *ptr, std::size_t length) {
    MARISA_DEBUG_IF((ptr == NULL) && (length != 0), MARISA_NULL_ERROR);
    ptr_ = ptr;
    length_ = length;
  }
  void set_id(std::size_t id) {
    id_ = id;
  }

  const char *ptr() const {
    return ptr_;
  }
  std::size_t length() const {
    return length_;
  }
  std::size_t id() const {
    return id_;
  }

  void clear() {
    Query().swap(*this);
  }
  void swap(Query &rhs) {
    marisa::swap(ptr_, rhs.ptr_);
    marisa::swap(length_, rhs.length_);
    marisa::swap(id_, rhs.id_);
  }

 private:
  const char *ptr_;
  std::size_t length_;
  std::size_t id_;
};

}  // namespace marisa

#endif  // MARISA_QUERY_H_




namespace marisa {
namespace grimoire {
namespace trie {

class State;

}  // namespace trie
}  // namespace grimoire

class Agent {
 public:
  Agent();
  ~Agent();

  const Query &query() const {
    return query_;
  }
  const Key &key() const {
    return key_;
  }

  void set_query(const char *str);
  void set_query(const char *ptr, std::size_t length);
  void set_query(std::size_t key_id);

  const grimoire::trie::State &state() const {
    return *state_;
  }
  grimoire::trie::State &state() {
    return *state_;
  }

  void set_key(const char *str) {
    MARISA_DEBUG_IF(str == NULL, MARISA_NULL_ERROR);
    key_.set_str(str);
  }
  void set_key(const char *ptr, std::size_t length) {
    MARISA_DEBUG_IF((ptr == NULL) && (length != 0), MARISA_NULL_ERROR);
    MARISA_DEBUG_IF(length > MARISA_UINT32_MAX, MARISA_SIZE_ERROR);
    key_.set_str(ptr, length);
  }
  void set_key(std::size_t id) {
    MARISA_DEBUG_IF(id > MARISA_UINT32_MAX, MARISA_SIZE_ERROR);
    key_.set_id(id);
  }

  bool has_state() const {
    return state_.get() != NULL;
  }
  void init_state();

  void clear();
  void swap(Agent &rhs);

 private:
  Query query_;
  Key key_;
  scoped_ptr<grimoire::trie::State> state_;

  // Disallows copy and assignment.
  Agent(const Agent &);
  Agent &operator=(const Agent &);
};

}  // namespace marisa

#endif  // MARISA_AGENT_H_




namespace marisa {
namespace grimoire {
namespace trie {

class LoudsTrie;

}  // namespace trie
}  // namespace grimoire

class Trie {
  friend class TrieIO;

 public:
  Trie();
  ~Trie();

  void build(Keyset &keyset, int config_flags = 0);

  void mmap(const char *filename);
  void map(const void *ptr, std::size_t size);

  void load(const char *filename);
  void read(int fd);

  void save(const char *filename) const;
  void write(int fd) const;

  bool lookup(Agent &agent) const;
  void reverse_lookup(Agent &agent) const;
  bool common_prefix_search(Agent &agent) const;
  bool predictive_search(Agent &agent) const;

  std::size_t num_tries() const;
  std::size_t num_keys() const;
  std::size_t num_nodes() const;

  TailMode tail_mode() const;
  NodeOrder node_order() const;

  bool empty() const;
  std::size_t size() const;
  std::size_t total_size() const;
  std::size_t io_size() const;

  void clear();
  void swap(Trie &rhs);

 private:
  scoped_ptr<grimoire::trie::LoudsTrie> trie_;

  // Disallows copy and assignment.
  Trie(const Trie &);
  Trie &operator=(const Trie &);
};

}  // namespace marisa

#endif  // MARISA_TRIE_H_




#endif  // MARISA_H_


