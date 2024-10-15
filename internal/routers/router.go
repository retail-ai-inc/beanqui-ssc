package routers

import (
	"log"
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
	handle  HandleFunc
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
		handle:  handleFunc,
	})
}

func (t *Router) match(pattern string, method string) (HandleFunc, bool) {
	for _, route := range t.routes {
		if route.pattern == pattern && route.method == FILE {
			return route.handle, true
		}
		if route.method == method && route.pattern == pattern {
			return route.handle, true
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

func (t *Router) Put(pattern string, handleFunc HandleFunc) {
	t.addRoute(pattern, http.MethodPut, handleFunc)
}

func (t *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if re := recover(); re != nil {
			log.Printf("panic error: %+v \n", re)
		}
	}()

	pattern := r.URL.Path
	method := r.Method

	if strings.HasSuffix(pattern, ".css") ||
		strings.HasSuffix(pattern, ".js") ||
		strings.HasSuffix(pattern, ".map") ||
		strings.HasSuffix(pattern, ".ico") ||
		strings.HasSuffix(pattern, ".vue") {
		pattern = "/"
		method = FILE
	}
	log.Printf("Method:%+v,UserAgent:%+v,URI:%+v \n", r.Method, r.UserAgent(), r.RequestURI)
	handle, b := t.match(pattern, method)
	if b {
		handle(w, r)
		return
	}
	http.NotFound(w, r)
	return
}
