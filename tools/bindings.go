package tools

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"unsafe"

	"gopkg.in/bblfsh/sdk.v1/uast"
)

// libuast can be linked in two modes on UNIX platforms: hosted and embedded.
// Hosted mode - libuast is installed globally in the system.
// Embedded mode - libuast source is inside "tools" directory and we compile it with cgo.
// This is what happens during `make dependencies`. It is the default.
//
// Build tags:
// custom_libuast - disables all the default CXXFLAGS and LDFLAGS.
// host_libuast - forces hosted mode.
//
// !unix defaults:
// CFLAGS: -Iinclude -DLIBUAST_STATIC
// CXXFLAGS: -Iinclude -DLIBUAST_STATIC
// LDFLAGS: -luast -lxml2 -Llib -static -lstdc++ -static-libgcc
// Notes: static linkage, libuast installation prefix is expected
// to be extracted into . ("tools΅ directory). Windows requires *both*
// CFLAGS and CXXFLAGS be set.
//
// unix defaults:
// CXXFLAGS: -I/usr/local/include -I/usr/local/include/libxml2 -I/usr/include -I/usr/include/libxml2
// LDFLAGS: -lxml2
// Notes: expects the embedded mode. "host_libuast" tag prepends -luast to LDFLAGS.
//
// Final notes:
// Cannot actually use "unix" tag until this is resolved: https://github.com/golang/go/issues/20322
// So inverted the condition: unix == !windows here.

// #cgo !custom_libuast,windows CFLAGS: -Iinclude -DLIBUAST_STATIC
// #cgo !custom_libuast,windows CXXFLAGS: -Iinclude -DLIBUAST_STATIC
// #cgo !custom_libuast,!windows CXXFLAGS: -I/usr/local/include -I/usr/local/include/libxml2 -I/usr/include -I/usr/include/libxml2
// #cgo !custom_libuast,host_libuast !custom_libuast,windows LDFLAGS: -luast
// #cgo !custom_libuast LDFLAGS: -lxml2
// #cgo !custom_libuast,windows LDFLAGS: -Llib -static -lstdc++ -static-libgcc
// #cgo !custom_libuast CXXFLAGS: -std=c++14
// #include "bindings.h"
import "C"

var (
	findMutex sync.Mutex
	spool     cstringPool
	kpool     = make(map[*uast.Node][]string)

	lastHandle handle
	nodes      = make(map[handle]*uast.Node)
)

type handle = uintptr

type ErrInvalidArgument struct {
	Message string
}

func (e *ErrInvalidArgument) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "invalid argument"
}

type errInternal struct {
	Method  string
	Message string
}

func (e *errInternal) Error() string {
	if e.Method == "" {
		if e.Message == "" {
			return "internal error"
		}
		return e.Message
	}
	return fmt.Sprintf("%s() failed: %s", e.Method, e.Message)
}

var itMutex sync.Mutex

// TreeOrder represents the traversal strategy for UAST trees
type TreeOrder int

const (
	// PreOrder traversal
	PreOrder TreeOrder = iota
	// PostOrder traversal
	PostOrder
	// LevelOrder (aka breadth-first) traversal
	LevelOrder
	// PositionOrder by node position in the source file
	PositionOrder
)

// Iterator allows for traversal over a UAST tree.
type Iterator struct {
	root     *uast.Node
	iterPtr  *C.UastIterator
	finished bool
}

var uastCtx *C.Uast

func init() {
	uastCtx = C.CreateUast()
}

func nodeToHandleC(node *uast.Node) C.NodeHandle {
	return C.NodeHandle(nodeToHandle(node))
}
func nodeToHandle(node *uast.Node) handle {
	if node == nil {
		return 0
	}
	lastHandle++
	h := lastHandle
	nodes[h] = node
	return h
}

func handleToNodeC(h C.NodeHandle) *uast.Node {
	return handleToNode(handle(h))
}
func handleToNode(h handle) *uast.Node {
	if h == 0 {
		return nil
	}
	n, ok := nodes[h]
	if !ok {
		panic(fmt.Errorf("unknown handle: %x", h))
	}
	return n
}

// initFilter converts the query string and node pointer to C types. It acquires findMutex
// and initializes the string pool. The caller should defer returned function to release
// the resources.
func initFilter(xpath string) (*C.char, func()) {
	findMutex.Lock()
	cquery := spool.getCstring(xpath)

	return cquery, func() {
		spool.release()
		kpool = make(map[*uast.Node][]string)
		nodes = make(map[handle]*uast.Node)
		findMutex.Unlock()
	}
}

func cError(name string) error {
	e := C.LastError()
	msg := strings.TrimSpace(C.GoString(e))
	C.free(unsafe.Pointer(e))
	// TODO: find a way to access this error code or constant
	if strings.HasPrefix(msg, "Invalid expression") {
		return &ErrInvalidArgument{Message: msg}
	}
	return &errInternal{Method: name, Message: msg}
}

// Filter takes a `*uast.Node` and a xpath query and filters the tree,
// returning the list of nodes that satisfy the given query.
// Filter is thread-safe but not concurrent by an internal global lock.
func Filter(node *uast.Node, xpath string) ([]*uast.Node, error) {
	if len(xpath) == 0 || node == nil {
		return nil, nil
	}

	cquery, closer := initFilter(xpath)
	defer closer()

	nodes := C.UastFilter(uastCtx, nodeToHandleC(node), cquery)
	if nodes == nil {
		return nil, cError("UastFilter")
	}
	defer C.NodesFree(nodes)

	nu := int(C.NodesSize(nodes))
	results := make([]*uast.Node, nu)
	for i := 0; i < nu; i++ {
		h := C.NodeAt(nodes, C.int(i))
		results[i] = handleToNodeC(h)
	}
	return results, nil
}

// FilterBool takes a `*uast.Node` and a xpath query with a boolean
// return type (e.g. when using XPath functions returning a boolean type).
// FilterBool is thread-safe but not concurrent by an internal global lock.
func FilterBool(node *uast.Node, xpath string) (bool, error) {
	if len(xpath) == 0 || node == nil {
		return false, nil
	}

	cquery, closer := initFilter(xpath)
	defer closer()

	var ok C.bool
	res := C.UastFilterBool(uastCtx, nodeToHandleC(node), cquery, &ok)
	if !bool(ok) {
		return false, cError("UastFilterBool")
	}
	return bool(res), nil
}

// FilterNumber takes a `*uast.Node` and a xpath query with a float
// return type (e.g. when using XPath functions returning a float type).
// FilterNumber is thread-safe but not concurrent by an internal global lock.
func FilterNumber(node *uast.Node, xpath string) (float64, error) {
	if len(xpath) == 0 || node == nil {
		return 0, nil
	}

	cquery, closer := initFilter(xpath)
	defer closer()

	var ok C.bool
	res := C.UastFilterNumber(uastCtx, nodeToHandleC(node), cquery, &ok)
	if !bool(ok) {
		return 0, cError("UastFilterNumber")
	}
	return float64(res), nil
}

// FilterString takes a `*uast.Node` and a xpath query with a string
// return type (e.g. when using XPath functions returning a string type).
// FilterString is thread-safe but not concurrent by an internal global lock.
func FilterString(node *uast.Node, xpath string) (string, error) {
	if len(xpath) == 0 || node == nil {
		return "", nil
	}

	cquery, closer := initFilter(xpath)
	defer closer()

	var res *C.char
	res = C.UastFilterString(uastCtx, nodeToHandleC(node), cquery)
	if res == nil {
		return "", cError("UastFilterString")
	}
	return C.GoString(res), nil
}

//export goGetInternalType
func goGetInternalType(ctx *C.Uast, ptr C.NodeHandle) *C.char {
	n := handleToNodeC(ptr)
	if n == nil {
		return nil
	}
	return spool.getCstring(n.InternalType)
}

//export goGetToken
func goGetToken(ctx *C.Uast, ptr C.NodeHandle) *C.char {
	n := handleToNodeC(ptr)
	if n == nil {
		return nil
	}
	return spool.getCstring(n.Token)
}

//export goGetChildrenSize
func goGetChildrenSize(ctx *C.Uast, ptr C.NodeHandle) C.int {
	n := handleToNodeC(ptr)
	if n == nil {
		return 0
	}
	return C.int(len(n.Children))
}

//export goGetChild
func goGetChild(ctx *C.Uast, ptr C.NodeHandle, index C.int) C.NodeHandle {
	n := handleToNodeC(ptr)
	if n == nil {
		return 0
	}
	child := n.Children[int(index)]
	return nodeToHandleC(child)
}

//export goGetRolesSize
func goGetRolesSize(ctx *C.Uast, ptr C.NodeHandle) C.int {
	n := handleToNodeC(ptr)
	if n == nil {
		return 0
	}
	return C.int(len(n.Roles))
}

//export goGetRole
func goGetRole(ctx *C.Uast, ptr C.NodeHandle, index C.int) C.uint16_t {
	n := handleToNodeC(ptr)
	if n == nil {
		return 0
	}
	role := n.Roles[int(index)]
	return C.uint16_t(role)
}

//export goGetPropertiesSize
func goGetPropertiesSize(ctx *C.Uast, ptr C.NodeHandle) C.int {
	n := handleToNodeC(ptr)
	if n == nil {
		return 0
	}
	return C.int(len(n.Properties))
}

func getPropertyKeys(ctx *C.Uast, node *uast.Node) []string {
	if keys, ok := kpool[node]; ok {
		return keys
	}
	p := node.Properties
	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	kpool[node] = keys
	return keys
}

//export goGetPropertyKey
func goGetPropertyKey(ctx *C.Uast, ptr C.NodeHandle, index C.int) *C.char {
	n := handleToNodeC(ptr)
	if n == nil {
		return nil
	}
	keys := getPropertyKeys(ctx, n)
	return spool.getCstring(keys[int(index)])
}

//export goGetPropertyValue
func goGetPropertyValue(ctx *C.Uast, ptr C.NodeHandle, index C.int) *C.char {
	n := handleToNodeC(ptr)
	if n == nil {
		return nil
	}
	keys := getPropertyKeys(ctx, n)
	p := n.Properties
	return spool.getCstring(p[keys[int(index)]])
}

//export goHasStartOffset
func goHasStartOffset(ctx *C.Uast, ptr C.NodeHandle) C.bool {
	n := handleToNodeC(ptr)
	if n == nil {
		return false
	}
	return n.StartPosition != nil
}

//export goGetStartOffset
func goGetStartOffset(ctx *C.Uast, ptr C.NodeHandle) C.uint32_t {
	n := handleToNodeC(ptr)
	if n == nil {
		return 0
	}
	p := n.StartPosition
	if p != nil {
		return C.uint32_t(p.Offset)
	}
	return 0
}

//export goHasStartLine
func goHasStartLine(ctx *C.Uast, ptr C.NodeHandle) C.bool {
	n := handleToNodeC(ptr)
	if n == nil {
		return false
	}
	return n.StartPosition != nil
}

//export goGetStartLine
func goGetStartLine(ctx *C.Uast, ptr C.NodeHandle) C.uint32_t {
	n := handleToNodeC(ptr)
	if n == nil {
		return 0
	}
	p := n.StartPosition
	if p != nil {
		return C.uint32_t(p.Line)
	}
	return 0
}

//export goHasStartCol
func goHasStartCol(ctx *C.Uast, ptr C.NodeHandle) C.bool {
	n := handleToNodeC(ptr)
	if n == nil {
		return false
	}
	return n.StartPosition != nil
}

//export goGetStartCol
func goGetStartCol(ctx *C.Uast, ptr C.NodeHandle) C.uint32_t {
	n := handleToNodeC(ptr)
	if n == nil {
		return 0
	}
	p := n.StartPosition
	if p != nil {
		return C.uint32_t(p.Col)
	}
	return 0
}

//export goHasEndOffset
func goHasEndOffset(ctx *C.Uast, ptr C.NodeHandle) C.bool {
	n := handleToNodeC(ptr)
	if n == nil {
		return false
	}
	return n.EndPosition != nil
}

//export goGetEndOffset
func goGetEndOffset(ctx *C.Uast, ptr C.NodeHandle) C.uint32_t {
	n := handleToNodeC(ptr)
	if n == nil {
		return 0
	}
	p := n.EndPosition
	if p != nil {
		return C.uint32_t(p.Offset)
	}
	return 0
}

//export goHasEndLine
func goHasEndLine(ctx *C.Uast, ptr C.NodeHandle) C.bool {
	n := handleToNodeC(ptr)
	if n == nil {
		return false
	}
	return n.EndPosition != nil
}

//export goGetEndLine
func goGetEndLine(ctx *C.Uast, ptr C.NodeHandle) C.uint32_t {
	n := handleToNodeC(ptr)
	if n == nil {
		return 0
	}
	p := n.EndPosition
	if p != nil {
		return C.uint32_t(p.Line)
	}
	return 0
}

//export goHasEndCol
func goHasEndCol(ctx *C.Uast, ptr C.NodeHandle) C.bool {
	n := handleToNodeC(ptr)
	if n == nil {
		return false
	}
	return n.EndPosition != nil
}

//export goGetEndCol
func goGetEndCol(ctx *C.Uast, ptr C.NodeHandle) C.uint32_t {
	n := handleToNodeC(ptr)
	if n == nil {
		return 0
	}
	p := n.EndPosition
	if p != nil {
		return C.uint32_t(p.Col)
	}
	return 0
}

// NewIterator constructs a new Iterator starting from the given `Node` and
// iterating with the traversal strategy given by the `order` parameter. Once
// the iteration have finished or you don't need the iterator anymore you must
// dispose it with the Dispose() method (or call it with `defer`).
func NewIterator(node *uast.Node, order TreeOrder) (*Iterator, error) {
	itMutex.Lock()
	defer itMutex.Unlock()

	it := C.UastIteratorNew(uastCtx, nodeToHandleC(node), C.TreeOrder(int(order)))
	if it == nil {
		return nil, cError("UastIteratorNew")
	}

	return &Iterator{
		root:     node,
		iterPtr:  it,
		finished: false,
	}, nil
}

// Next retrieves the next `Node` in the tree's traversal or `nil` if there are no more
// nodes. Calling `Next()` on a finished iterator after the first `nil` will
// return an error.This is thread-safe but not concurrent by an internal global lock.
func (i *Iterator) Next() (*uast.Node, error) {
	itMutex.Lock()
	defer itMutex.Unlock()

	if i.finished {
		return nil, fmt.Errorf("Next() called on finished iterator")
	}

	h := handle(C.UastIteratorNext(i.iterPtr))
	if h == 0 {
		// End of the iteration
		i.finished = true
		return nil, nil
	}
	return handleToNode(h), nil
}

// Iterate function is similar to Next() but returns the `Node`s in a channel. It's mean
// to be used with the `for node := range myIter.Iterate() {}` loop.
func (i *Iterator) Iterate() <-chan *uast.Node {
	c := make(chan *uast.Node)
	if i.finished {
		close(c)
		return c
	}

	go func() {
		for {
			n, err := i.Next()
			if n == nil || err != nil {
				close(c)
				break
			}

			c <- n
		}
	}()

	return c
}

// Dispose must be called once you've finished using the iterator or preventively
// with `defer` to free the iterator resources. Failing to do so would produce
// a memory leak.
func (i *Iterator) Dispose() {
	itMutex.Lock()
	defer itMutex.Unlock()

	if i.iterPtr != nil {
		C.UastIteratorFree(i.iterPtr)
		i.iterPtr = nil
	}
	i.finished = true
	i.root = nil
}
