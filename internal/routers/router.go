package routers

import (
	"net/http"
	"strings"
)

const (
	FILE = "FILE"
)

type HandleFunc func(w http.ResponseWriter, r *http.Request)
type Route struct {
	method  string
	pattern string
	Handle  HandleFunc
}

type Router struct {
	routes []Route
}

func NewRouter() *Router {
	return &Router{routes: make([]Route, 0)}
}

func (t *Router) addRoute(pattern string, method string, handleFunc HandleFunc) {
	t.routes = append(t.routes, Route{
		method:  method,
		pattern: pattern,
		Handle:  handleFunc,
	})
}

func (t *Router) match(pattern string, method string) (HandleFunc, bool) {
	for _, route := range t.routes {
		if route.pattern == pattern && route.method == FILE {
			return route.Handle, true
		}
		if route.method == method && route.pattern == pattern {
			return route.Handle, true
		}
	}
	return nil, false
}

func (t *Router) Get(pattern string, handleFunc HandleFunc) {
	t.addRoute(pattern, http.MethodGet, handleFunc)
}

func (t *Router) Post(pattern string, handleFunc HandleFunc) {
	t.addRoute(pattern, http.MethodPost, handleFunc)
}

func (t *Router) File(pattern string, handleFunc HandleFunc) {
	t.addRoute(pattern, FILE, handleFunc)
}

func (t *Router) Delete(pattern string, handleFunc HandleFunc) {
	t.addRoute(pattern, http.MethodDelete, handleFunc)
}

func (t *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	lastLen := len(r.URL.Path) - 4
	if lastLen < 0 {
		lastLen = 0
	}
	last := r.URL.Path[lastLen:]
	pattern := r.URL.Path
	method := r.Method

	if strings.Contains(last, ".css") ||
		strings.Contains(last, ".js") ||
		strings.Contains(last, ".map") ||
		strings.Contains(last, ".ico") ||
		strings.Contains(last, ".vue") {
		pattern = "/"
		method = FILE
	}

	handle, b := t.match(pattern, method)
	if b {
		handle(w, r)
		return
	}
	http.NotFound(w, r)
	return
}
