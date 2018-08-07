#ifndef CLIENT_GO_BINDINGS_H_
#define CLIENT_GO_BINDINGS_H_

#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

#if __has_include("uast.h") // std C++17, GCC 5.x || Clang || VSC++ 2015u2+
// Embedded mode on UNIX, MSVC build on Windows.
#include "uast.h"
#else
// Hosted mode on UNIX, MinGW build on Windows.
#include "libuast/uast.h"
#endif

extern char* goGetInternalType(void*);
extern char* goGetToken(void*);
extern int goGetChildrenSize(void*);
extern void* goGetChild(void*, int);
extern int goGetRolesSize(void*);
extern uint16_t goGetRole(void*, int);
extern int goGetPropertiesSize(void*);
extern char* goGetPropertyKey(void*, int);
extern char* goGetPropertyValue(void*, int);
extern bool goHasStartOffset(void*);
extern uint32_t goGetStartOffset(void*);
extern bool goHasStartLine(void*);
extern uint32_t goGetStartLine(void*);
extern bool goHasStartCol(void*);
extern uint32_t goGetStartCol(void*);
extern bool goHasEndOffset(void*);
extern uint32_t goGetEndOffset(void*);
extern bool goHasEndLine(void*);
extern uint32_t goGetEndLine(void*);
extern bool goHasEndCol(void*);
extern uint32_t goGetEndCol(void*);

static const char *InternalType(const void *node) {
  return goGetInternalType((void*)node);
}

static const char *Token(const void *node) {
  return goGetToken((void*)node);
}

static size_t ChildrenSize(const void *node) {
  return goGetChildrenSize((void*)node);
}

static void *ChildAt(const void *data, int index) {
  return (void*)goGetChild((void*)data, index);
}

static size_t RolesSize(const void *node) {
  return goGetRolesSize((void*)node);
}

static uint16_t RoleAt(const void *node, int index) {
  return goGetRole((void*)node, index);
}

static size_t PropertiesSize(const void *node) {
  return goGetPropertiesSize((void*)node);
}

static const char *PropertyKeyAt(const void *node, int index) {
  return goGetPropertyKey((void*)node, index);
}

static const char *PropertyValueAt(const void *node, int index) {
  return goGetPropertyValue((void*)node, index);
}

static bool HasStartOffset(const void *node) {
  return goHasStartOffset((void*)node);
}

static uint32_t StartOffset(const void *node) {
  return goGetStartOffset((void*)node);
}

static bool HasStartLine(const void *node) {
  return goHasStartLine((void*)node);
}

static uint32_t StartLine(const void *node) {
  return goGetStartLine((void*)node);
}

static bool HasStartCol(const void *node) {
  return goHasStartCol((void*)node);
}

static uint32_t StartCol(const void *node) {
  return goGetStartCol((void*)node);
}

static bool HasEndOffset(const void *node) {
  return goHasEndOffset((void*)node);
}

static uint32_t EndOffset(const void *node) {
  return goGetEndOffset((void*)node);
}

static bool HasEndLine(const void *node) {
  return goHasEndLine((void*)node);
}

static uint32_t EndLine(const void *node) {
  return goGetEndLine((void*)node);
}

static bool HasEndCol(const void *node) {
  return goHasEndCol((void*)node);
}

static uint32_t EndCol(const void *node) {
  return goGetEndCol((void*)node);
}

static Uast* CreateUast() {
  return UastNew((NodeIface){
      .InternalType = InternalType,
      .Token = Token,
      .ChildrenSize = ChildrenSize,
      .ChildAt = ChildAt,
      .RolesSize = RolesSize,
      .RoleAt = RoleAt,
      .PropertiesSize = PropertiesSize,
      .PropertyKeyAt = PropertyKeyAt,
      .PropertyValueAt = PropertyValueAt,
      .HasStartOffset = HasStartOffset,
      .StartOffset = StartOffset,
      .HasStartLine = HasStartLine,
      .StartLine = StartLine,
      .HasStartCol = HasStartCol,
      .StartCol = StartCol,
      .HasEndOffset = HasEndOffset,
      .EndOffset = EndOffset,
      .HasEndLine = HasEndLine,
      .EndLine = EndLine,
      .HasEndCol = HasEndCol,
      .EndCol = EndCol,
  });
}

#endif // CLIENT_GO_BINDINGS_H_
