package mux

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// ErrNilHandler ...
var ErrNilHandler = errors.New("handler can't be empty")

// ErrNotStartSlash ...
var ErrNotStartSlash = errors.New("pattern must start with /")

// ErrAmbigiousMapping ...
var ErrAmbigiousMapping = errors.New("ambigious mapping")

type contextKey string

var pathParamsKey = contextKey("params")

type exactMuxEntry struct {
	pattern string
	handler http.Handler
}

// ExactMux ...
type ExactMux struct {
	mutex            sync.RWMutex
	exactRoutes      map[string]map[string]exactMuxEntry
	paramRoutes      map[string][]MuxParams
	notFoundHandler  http.Handler
	MethodNotAllowed http.Handler
}

// Middleware ...
type Middleware func(handler http.HandlerFunc) http.HandlerFunc

// NewExactMux ...
func NewExactMux() *ExactMux {
	return &ExactMux{}
}

// ServerHTTP ...
func (m *ExactMux) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request, handler, err := m.handler(request.Method, request.URL.Path, request); err == nil {
		handler.ServeHTTP(writer, request)
		return
	}

	if m.notFoundHandler != nil {
		m.notFoundHandler.ServeHTTP(writer, request)
		return
	}

	writer.WriteHeader(404)
}

// GET ...
func (m *ExactMux) GET(
	pattern string,
	handlerFunc http.HandlerFunc,
	middlewares ...Middleware,
) {
	m.HandleFuncWithMiddlewares(
		http.MethodGet,
		pattern,
		handlerFunc,
		middlewares...,
	)
}

// POST ...
func (m *ExactMux) POST(
	pattern string,
	handlerFunc http.HandlerFunc,
	middlewares ...Middleware,
) {
	m.HandleFuncWithMiddlewares(
		http.MethodPost,
		pattern,
		handlerFunc,
		middlewares...,
	)
}

// DELETE ...
func (m *ExactMux) DELETE(
	pattern string,
	handlerFunc http.HandlerFunc,
	middlewares ...Middleware,
) {
	m.HandleFuncWithMiddlewares(
		http.MethodDelete,
		pattern,
		handlerFunc,
		middlewares...,
	)
}

// PUT ...
func (m *ExactMux) PUT(
	pattern string,
	handlerFunc http.HandlerFunc,
	middlewares ...Middleware,
) {
	m.HandleFuncWithMiddlewares(
		http.MethodPut,
		pattern,
		handlerFunc,
		middlewares...,
	)
}

// HandleFuncWithMiddlewares ...
func (m *ExactMux) HandleFuncWithMiddlewares(
	method string,
	pattern string,
	handlerFunc http.HandlerFunc,
	middlewares ...Middleware,
) {
	for _, middleware := range middlewares {
		handlerFunc = middleware(handlerFunc)
	}
	m.HandleFunc(method, pattern, handlerFunc)
}

// HandleFunc ...
func (m *ExactMux) HandleFunc(method string, pattern string, handlerFunc http.HandlerFunc) {
	if !strings.HasPrefix(pattern, "/") {
		panic(ErrNotStartSlash)
	}

	if handlerFunc == nil {
		panic(ErrNilHandler)
	}

	if isExact(pattern) {
		m.AddExact(method, pattern, handlerFunc)
		return
	}
	m.AddParams(method, pattern, handlerFunc)
}

// AddExact ...
func (m *ExactMux) AddExact(method string, pattern string, handlerFunc http.Handler) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	entry := exactMuxEntry{
		pattern: pattern,
		handler: handlerFunc,
	}

	if _, exists := m.exactRoutes[method][pattern]; exists {
		panic(ErrAmbigiousMapping)
	}

	if m.exactRoutes == nil {
		m.exactRoutes = make(map[string]map[string]exactMuxEntry)
	}

	if m.exactRoutes[method] == nil {
		m.exactRoutes[method] = make(map[string]exactMuxEntry)
	}

	m.exactRoutes[method][pattern] = entry
}

// AddParams ...
func (m *ExactMux) AddParams(method string, pattern string, handlerFunc http.Handler) {
	entry := parsePathParams(pattern)
	entry.handler = handlerFunc

	if m.paramRoutes == nil {
		m.paramRoutes = make(map[string][]MuxParams)

	}

	if m.paramRoutes[method] == nil {
		m.paramRoutes[method] = make([]MuxParams, 0)
	}

	m.paramRoutes[method] = append(m.paramRoutes[method], entry)
}

func (m *ExactMux) handler(method string, path string, original *http.Request) (result *http.Request, handler http.Handler, err error) {
	exactEntries, exactExists := m.exactRoutes[method]
	if exactExists {
		if entry, ok := exactEntries[path]; ok {
			return original, entry.handler, nil
		}
	}
	paramEntries, paramExists := m.paramRoutes[method]
	if !paramExists {
		return nil, nil, fmt.Errorf("no handlers for %s, %s", method, path)
	}
	weight := calculateWeight(path)
	for _, paramEntry := range paramEntries {
		if weight != paramEntry.weight {
			continue
		}

		if params, ok := paramEntry.Match(path); ok {
			ctx := context.WithValue(original.Context(), pathParamsKey, params)
			result = original.WithContext(ctx)
			return result, paramEntry.handler, nil
		}
	}

	return nil, nil, fmt.Errorf("can't find handler for: %s, %s", method, path)
}

// FromContext ...
func FromContext(ctx context.Context, key string) (value string, ok bool) {
	params, ok := ctx.Value(pathParamsKey).(map[string]string)
	if !ok {
		return "", false
	}
	param, exists := params[key]
	return param, exists
}

// isExact ...
func isExact(pattern string) bool {
	return !strings.Contains(pattern, "{")
}
