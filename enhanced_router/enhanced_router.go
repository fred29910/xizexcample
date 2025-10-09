package enhanced_router

import (
	"context"
	"sync"
)

// HandleFunc defines the request handler function signature.
type HandleFunc func(ctx context.Context, req interface{}) (resp interface{}, err error)

// MiddlewareFunc defines the middleware function signature.
type MiddlewareFunc func(next HandleFunc) HandleFunc

// EnhancedRouter is an interface that extends the basic router capabilities.
type EnhancedRouter interface {
	Use(mw ...MiddlewareFunc)
	AddRoute(path string, handler HandleFunc)
	Serve(ctx context.Context, path string, req interface{}) (interface{}, error)
}

type enhancedRouter struct {
	mu          sync.RWMutex
	middlewares []MiddlewareFunc
	routes      map[string]HandleFunc
}

// NewEnhancedRouter creates a new instance of EnhancedRouter.
func NewEnhancedRouter() EnhancedRouter {
	return &enhancedRouter{
		routes: make(map[string]HandleFunc),
	}
}

// Use adds middleware to the router.
func (e *enhancedRouter) Use(mw ...MiddlewareFunc) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.middlewares = append(e.middlewares, mw...)
}

// AddRoute adds a new route to the router.
func (e *enhancedRouter) AddRoute(path string, handler HandleFunc) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Wrap handler with middlewares in reverse order to create an onion-like chain.
	h := handler
	for i := len(e.middlewares) - 1; i >= 0; i-- {
		h = e.middlewares[i](h)
	}
	e.routes[path] = h
}

// Serve executes the handler for a given path.
func (e *enhancedRouter) Serve(ctx context.Context, path string, req interface{}) (interface{}, error) {
	e.mu.RLock()
	handler, ok := e.routes[path]
	e.mu.RUnlock()

	if !ok {
		return nil, nil // Or return an error
	}

	return handler(ctx, req)
}
