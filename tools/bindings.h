#ifndef CLIENT_GO_BINDINGS_H_
#define CLIENT_GO_BINDINGS_H_

#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

#include "uast.h"

extern char* goGetInternalType(uintptr_t);
extern int goGetPropertiesSize(uintptr_t);
extern char* goGetToken(uintptr_t);
extern int goGetChildrenSize(uintptr_t);
extern uintptr_t goGetChild(uintptr_t, int);
extern int goGetRolesSize(uintptr_t);
extern uint16_t goGetRole(uintptr_t, int);

static const char *InternalType(const void *node) {
  return goGetInternalType((uintptr_t)node);
}

static const char *Token(const void *node) {
  return goGetToken((uintptr_t)node);
}

static int ChildrenSize(const void *node) {
  return goGetChildrenSize((uintptr_t)node);
}

static void *ChildAt(const void *data, int index) {
  return (void*)goGetChild((uintptr_t)data, index);
}

static int RolesSize(const void *node) {
  return goGetRolesSize((uintptr_t)node);
}

static uint16_t RoleAt(const void *node, int index) {
  return goGetRole((uintptr_t)node, index);
}

static Uast *ctx;
static Nodes *nodes;

static void CreateUast() {
  ctx = UastNew((NodeIface){
      .InternalType = InternalType,
      .Token = Token,
      .ChildrenSize = ChildrenSize,
      .ChildAt = ChildAt,
      .RolesSize = RolesSize,
      .RoleAt = RoleAt,
  });
}

static bool Filter(uintptr_t node_ptr, const char *query) {
  nodes = UastFilter(ctx, (void*)node_ptr, query);
  return nodes != NULL;
}

static int Size() {
  return NodesSize(nodes);
}

static uintptr_t At(int i) {
  return (uintptr_t)NodeAt(nodes, i);
}

#endif // CLIENT_GO_BINDINGS_H_
