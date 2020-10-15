package ginserver

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var (
	Bind         = (*gin.Context).ShouldBind
	BindJSON     = (*gin.Context).ShouldBindJSON
	BindXML      = (*gin.Context).ShouldBindXML
	BindQuery    = (*gin.Context).ShouldBindQuery
	BindYAML     = (*gin.Context).ShouldBindYAML
	BindHeader   = (*gin.Context).ShouldBindHeader
	BindURI      = (*gin.Context).ShouldBindUri
	BindForm     = func(c *gin.Context, obj interface{}) error { return c.ShouldBindWith(obj, binding.Form) }
	BindFormPost = func(c *gin.Context, obj interface{}) error { return c.ShouldBindWith(obj, binding.FormPost) }
)
