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
	 * 参数绑定：&Persion{} (newPerson())
	 * 参数来源：url
	 *
	 * curl 'http://127.0.0.1:8080/get2?name=ch&address=china'
	 * output: => param: type(*main.person), content-type: , val: {"name":"ch","age":0,"address":"china"}
	 *
	 * curl -X GET 'http://127.0.0.1:8080/get2?name=ch' --data '{"name":"sh", "age": 100, "address":"china"}' -H "Content-Type:application/json"
	 * output: => param: type(*main.person), content-type: application/json, val: {"name":"sh","age":100,"address":"china"}
	 */
	svr.GET("/get2", doPerson, newPerson, ginserver.Bind, ginserver.BindJSON)

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
	 * 参数绑定：&Persion{} (newPerson())
	 * 参数来源：url + (post&json || post&x-www-form-urlencoded)
	 *
	 * # 数据为 JSON，绑定方法为 BindJSON
	 * curl 'http://127.0.0.1:8080/post2?age=18' --data '{"name": "ch", "address": "china"}' -H "Content-Type:application/json"
	 * output:  => param: type(*main.person), content-type: application/json, val: {"name":"ch","age":18,"address":"china"}
	 *
	 * # 数据为 x-www-form-urlencoded，绑定方法为 BindFormPost
	 * curl 'http://127.0.0.1:8080/post2?name=ch' --data 'age=18&address=china' -H "Content-Type:application/x-www-form-urlencoded"
	 * output: => param: type(*main.person), content-type: application/x-www-form-urlencoded, val: {"name":"ch","age":18,"address":"china"}
	 */
	svr.POST("/post2", doPerson, newPerson, ginserver.BindForm, ginserver.BindJSON, ginserver.BindFormPost)

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
