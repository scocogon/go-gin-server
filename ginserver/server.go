package ginserver

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ParamBind func(*gin.Context, interface{}) error
type ParamCreator func() interface{}
type FuncExec func(*gin.Context, interface{})

type Server struct {
	egn *gin.Engine
}

func NewServer(egn *gin.Engine) *Server {
	return &Server{
		egn: egn,
	}
}

func (s *Server) Run(addr ...string) (err error) {
	return s.egn.Run(addr...)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (s *Server) POST(relativePath string, exec FuncExec, creator ParamCreator, binds ...ParamBind) gin.IRoutes {
	return s.bindParam(http.MethodPost, relativePath, exec, creator, binds)
}

// GET is a shortcut for router.Handle("GET", path, handle).
func (s *Server) GET(relativePath string, exec FuncExec, creator ParamCreator, binds ...ParamBind) gin.IRoutes {
	return s.bindParam(http.MethodGet, relativePath, exec, creator, binds)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (s *Server) DELETE(relativePath string, exec FuncExec, creator ParamCreator, binds ...ParamBind) gin.IRoutes {
	return s.bindParam(http.MethodDelete, relativePath, exec, creator, binds)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (s *Server) PATCH(relativePath string, exec FuncExec, creator ParamCreator, binds ...ParamBind) gin.IRoutes {
	return s.bindParam(http.MethodPatch, relativePath, exec, creator, binds)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (s *Server) PUT(relativePath string, exec FuncExec, creator ParamCreator, binds ...ParamBind) gin.IRoutes {
	return s.bindParam(http.MethodPut, relativePath, exec, creator, binds)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (s *Server) OPTIONS(relativePath string, exec FuncExec, creator ParamCreator, binds ...ParamBind) gin.IRoutes {
	return s.bindParam(http.MethodOptions, relativePath, exec, creator, binds)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (s *Server) HEAD(relativePath string, exec FuncExec, creator ParamCreator, binds ...ParamBind) gin.IRoutes {
	return s.bindParam(http.MethodHead, relativePath, exec, creator, binds)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (s *Server) Any(relativePath string, exec FuncExec, creator ParamCreator, binds ...ParamBind) gin.IRoutes {
	s.bindParam(http.MethodGet, relativePath, exec, creator, binds)
	s.bindParam(http.MethodPost, relativePath, exec, creator, binds)
	s.bindParam(http.MethodPut, relativePath, exec, creator, binds)
	s.bindParam(http.MethodPatch, relativePath, exec, creator, binds)
	s.bindParam(http.MethodHead, relativePath, exec, creator, binds)
	s.bindParam(http.MethodOptions, relativePath, exec, creator, binds)
	s.bindParam(http.MethodDelete, relativePath, exec, creator, binds)
	s.bindParam(http.MethodConnect, relativePath, exec, creator, binds)
	s.bindParam(http.MethodTrace, relativePath, exec, creator, binds)
	return s.egn
}

// BindParam 参数获取
func (s *Server) BindParam(httpMethod, relativePath string, exec FuncExec, creator ParamCreator, binds ...ParamBind) gin.IRoutes {
	return s.bindParam(httpMethod, relativePath, exec, creator, binds)
}

func (s *Server) bindParam(httpMethod, relativePath string, exec FuncExec, creator ParamCreator, binds []ParamBind) gin.IRoutes {
	if len(binds) == 0 {
		return s.egn.Handle(httpMethod, relativePath, func(c *gin.Context) {
			exec(c, nil)
		})
	}

	return s.egn.Handle(httpMethod, relativePath, func(c *gin.Context) {
		var param interface{}
		if creator != nil {
			param = creator()
		} else {
			param = &map[string]interface{}{}
		}

		for _, f := range binds {
			if err := f(c, param); err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
			}
		}

		exec(c, param)
	})
}
