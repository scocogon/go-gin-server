# go-gin-server

## 安装

```
go get -u github.com/scocogon/go-gin-server
```

## 说明

    简易的 gin 封装

    注意 tag
        form - 表单元素形式: a=1&b=a&c=100
        json - 数据json形式: {"a": 1, "b": "a", "c": 100}
        uri  - 在 url 路径中: /a/b/:c

```
type person struct {
	Name    string `form:"name" json:"name" uri:"name"`
	Age     int    `form:"age" json:"age"`
	Address string `form:"address" json:"address"`
}
```

## 解决问题

+ 提供同时多路获取 http 请求内容的语法糖

## 问题

+ 每一次 http 请求都会尝试进行多路绑定

## 使用

### 调用方法

> 与 gin 一样的方法

```
POST
GET
DELETE
PATCH
PUT
OPTIONS
HEAD
Any
```

### 使用说明

> 所有方法的使用原型一样

```
# relativePath - url 请求路径
# exec - 原型: func(c *gin.Context, param interface{})
#        处理数据的功能方法，c 为上游传入的 *gin.Context， param 为绑定后的请求内容
#           可以通过 c 获取参数数据
#           没有提供绑定方法，或没有参数，或参数与绑定方式不一致，或未添加结构中对应 tag，param 将为 nil
# creator - 原型: func() interface{}
#           该方法返回接口接收的数据类型的对象
# binds - 原型: func(*gin.Context, interface{})
#         用于绑定数据的方式，对应于 gin.Context.ShouldBind???(...)
#         可传多个，越往后，优先级越高
#         快捷方式
#           Bind
#           BindJSON
#           BindXML
#           BindQuery
#           BindYAML
#           BindHeader
#           BindURI
#           BindForm
#           BindFormPost
#
# 返回类型: 为了与 gin 保持一致，尽量不要使用
POST(relativePath string, exec FuncExec, creator ParamCreator, binds ...ParamBind) gin.IRoutes
```

### 使用示例

> 具体使用可参见 [main.go](main.go)

```
package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/scocogon/go-gin-server/ginserver"
)

type person struct {
	Name    string `form:"name" json:"name" uri:"name"`
	Age     int    `form:"age" json:"age"`
	Address string `form:"address" json:"address"`
}

func newPerson() interface{} { return &person{} }

func doPerson(c *gin.Context, param interface{}) {
	if param == nil {
		fmt.Println("nil")
		return
	}

	s, _ := json.Marshal(param)
	fmt.Printf(" => param: type(%T), content-type: %s, val: %s\n", param, c.ContentType(), string(s))
	// fmt.Println((*param.(*map[string]interface{}))["name"])
}

func main() {
	egn := gin.Default()
	svr := ginserver.NewServer(egn)

	/**
	 * 参数默认绑定：&map[string]interface{}
	 * 参数来源：url
	 *   url中的参数，无法写入 map[string]interface{}，
	 *   可以在 doPerson 中通过 c.Query("") 单独获取
	 *
	 * curl 'http://127.0.0.1:8080/get1?name=ch&address=china'
	 * output: => param: type(*map[string]interface {}), content-type: , val: {}
	 */
	svr.GET("/get1", doPerson, nil, ginserver.Bind)

	/**
	 * 参数默认绑定：&map[string]interface{}
	 * 参数来源：url + post&json
	 *   url中的参数，无法写入 map[string]interface{}，
	 *   可以在 doPerson 中通过 c.Query("") 单独获取
	 *
	 * age=18 无法写入结果中
	 * curl 'http://127.0.0.1:8080/post1?age=18' --data '{"name": "ch", "address": "china"}' -H "Content-Type:application/json"
	 * output: => param: type(*map[string]interface {}), content-type: application/json, val: {"address":"china","name":"ch"}
	 */
	svr.POST("/post1", doPerson, nil, ginserver.BindForm, ginserver.BindJSON)

	/**
	 * 参数默认绑定：&Persion{} (newPerson())
	 * 参数来源：url
	 *
	 * curl 'http://127.0.0.1:8080/get/sh?age=100&address=china'
	 * output: => param: type(*main.person), content-type: , val: {"name":"sh","age":100,"address":"china"}
	 */
	svr.GET("/get/:name", doPerson, newPerson, ginserver.BindURI, ginserver.BindForm)

	svr.Run(":8080")
}

```