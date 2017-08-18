#ifndef _GO_CLIENT_BINDINGS_H_
#define _GO_CLIENT_BINDINGS_H_

#include <libuast/uast.h>

union void_cast {
  uintptr_t i;
  void *ptr;
};

static void *toPtr(uintptr_t a) {
  union void_cast cast;
  cast.i = a;
  return cast.ptr;
}

extern char* goGetInternalType(void*);
extern int goGetPropertiesSize(void*);
extern char* goGetToken(void*);
extern int goGetChildrenSize(void*);
extern uintptr_t goGetChild(void*, int);
extern int goGetRolesSize(void*);
extern uint16_t goGetRole(void*, int);

static const char *
get_internal_type(const void *node)
{
  const char *name = goGetInternalType((void*)node);
  return name;
}

static const char *get_token(const void *node)
{
  const char *name = goGetToken((void*)node);
  return name;
}

static int get_children_size(const void *node)
{
  return goGetChildrenSize((void*)node);
}

static void *get_child(const void *data, int index)
{
  return toPtr(goGetChild((void *)data, index));
}

static int get_roles_size(const void *node)
{
  return goGetRolesSize((void*)node);
}

static uint16_t get_role(const void *node, int index)
{
  return 2;
}

static node_api *api;
static find_ctx *ctx;

static void create_go_node_api()
{
  api = new_node_api((node_iface){
      .internal_type = get_internal_type,
      .token = get_token,
      .children_size = get_children_size,
      .children = get_child,
      .roles_size = get_roles_size,
      .roles = get_role,
  });
  ctx = new_find_ctx();
}

static int _api_find(uintptr_t node_ptr, const char *query) {
  void *node = toPtr(node_ptr);
  return node_api_find(api, ctx, node, query);
}

static int _api_get_nu_results() {
  return find_ctx_get_len(ctx);
}

static void *_api_get_result(unsigned int i) {
  return find_ctx_get(ctx, i);
}

#endif
