#ifndef _GO_CLIENT_BINDINGS_H_
#define _GO_CLIENT_BINDINGS_H_

#include <libuast/uast.h>


extern char* goGetInternalType(uintptr_t);
extern int goGetPropertiesSize(uintptr_t);
extern char* goGetToken(uintptr_t);
extern int goGetChildrenSize(uintptr_t);
extern uintptr_t goGetChild(uintptr_t, int);
extern int goGetRolesSize(uintptr_t);
extern uint16_t goGetRole(uintptr_t, int);

static const char *get_internal_type(const void *node)
{
  return goGetInternalType((uintptr_t)node);
}

static const char *get_token(const void *node)
{
  return goGetToken((uintptr_t)node);
}

static int get_children_size(const void *node)
{
  return goGetChildrenSize((uintptr_t)node);
}

static void *get_child(const void *data, int index)
{
  return (void*)goGetChild((uintptr_t)data, index);
}

static int get_roles_size(const void *node)
{
  return goGetRolesSize((uintptr_t)node);
}

static uint16_t get_role(const void *node, int index)
{
  return goGetRole((uintptr_t)node, index);
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
  return node_api_find(api, ctx, (void*)node_ptr, query);
}

static int _api_get_nu_results() {
  return find_ctx_get_len(ctx);
}

static uintptr_t _api_get_result(unsigned int i) {
  return (uintptr_t)find_ctx_get(ctx, i);
}

#endif
