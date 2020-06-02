package web

import (
	"github.com/gin-gonic/gin"
	jsonutils "github.com/sayuthisobri/goutils/json"
	"net/http"
)

//
// HandlerFunc - handler function
//
type HandlerFunc func(ctx *Ctx)

//
// Router - router struct
//
type Router struct {
	gin.IRouter
	DI *jsonutils.J
}

//
// NewRouter - new router creation
//
func NewRouter(IRouter gin.IRouter) *Router {
	return &Router{IRouter: IRouter}
}

//
// AddHandler - add handlers with context
//
func (r *Router) AddHandler(path string, handler ControllerHandler, handlers ...HandlerFunc) *Router {
	group := NewRouter(r.Group(path, r.wrapHandlers(handlers...)...))
	group.DI = r.DI
	if handler != nil {
		handler(group)
	}
	return group
}

//
// AddHandler - add handlers
//
func (r *Router) AddHandlerRaw(path string, handler ControllerHandler, handlers ...gin.HandlerFunc) *Router {
	group := NewRouter(r.Group(path, handlers...))
	group.DI = r.DI
	if handler != nil {
		handler(group)
	}
	return group
}

func (r *Router) wrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	var ginHandlers []gin.HandlerFunc

	for _, h := range handlers {
		ginHandlers = append(ginHandlers, func(context *gin.Context) {
			h(&Ctx{Context: context, Router: r})
		})
	}
	return ginHandlers
}

func (r *Router) Use(handlers ...HandlerFunc) {
	r.IRouter.Use(r.wrapHandlers(handlers...)...)
}

func (r *Router) Handle(method string, path string, handlers ...HandlerFunc) {
	r.IRouter.Handle(method, path, r.wrapHandlers(handlers...)...)

}
func (r *Router) Any(path string, handlers ...HandlerFunc) {
	r.IRouter.Any(path, r.wrapHandlers(handlers...)...)
}
func (r *Router) GET(path string, handlers ...HandlerFunc) {
	r.IRouter.Handle(http.MethodGet, path, r.wrapHandlers(handlers...)...)
}
func (r *Router) POST(path string, handlers ...HandlerFunc) {
	r.IRouter.Handle(http.MethodPost, path, r.wrapHandlers(handlers...)...)
}
func (r *Router) DELETE(path string, handlers ...HandlerFunc) {
	r.IRouter.Handle(http.MethodDelete, path, r.wrapHandlers(handlers...)...)
}
func (r *Router) PATCH(path string, handlers ...HandlerFunc) {
	r.IRouter.Handle(http.MethodPatch, path, r.wrapHandlers(handlers...)...)
}
func (r *Router) PUT(path string, handlers ...HandlerFunc) {
	r.IRouter.Handle(http.MethodPut, path, r.wrapHandlers(handlers...)...)
}
func (r *Router) OPTIONS(path string, handlers ...HandlerFunc) {
	r.IRouter.Handle(http.MethodOptions, path, r.wrapHandlers(handlers...)...)
}
func (r *Router) HEAD(path string, handlers ...HandlerFunc) {
	r.IRouter.Handle(http.MethodHead, path, r.wrapHandlers(handlers...)...)
}
