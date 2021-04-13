// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64

package ginframe

import (
	"fmt"
	"net/http"
	"path"
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"

	"github.com/TingYunGo/goagent"
)

func getHandlerName(handler gin.HandlerFunc) string {
	handlerPC := reflect.ValueOf(handler).Pointer()
	return runtime.FuncForPC(handlerPC).Name()
}
func preHandler(relativePath, method string) gin.HandlerFunc {
	return func(c *gin.Context) {
		action := tingyun3.GetAction()
		if action == nil {
			return
		}
		action.SetName(relativePath, method)
	}
}
func pushFrontHandler(group *gin.RouterGroup, relativePath string, handlers []gin.HandlerFunc) []gin.HandlerFunc {
	count := len(handlers)
	if count == 0 {
		return handlers
	}
	newHandlers := make([]gin.HandlerFunc, count+1)
	newHandlers[0] = preHandler(path.Join(group.BasePath(), relativePath), getHandlerName(handlers[0]))
	for i := range handlers {
		newHandlers[i+1] = handlers[i]
	}
	return newHandlers
}

//go:noinline
func RouterGrouphandle(group *gin.RouterGroup, httpMethod, relativePath string, handlers gin.HandlersChain) gin.IRoutes {
	fmt.Println(httpMethod, relativePath)
	return nil
}

//go:noinline
func WrapRouterGrouphandle(group *gin.RouterGroup, httpMethod, relativePath string, handlers gin.HandlersChain) gin.IRoutes {
	handlers = pushFrontHandler(group, relativePath, handlers)
	return RouterGrouphandle(group, httpMethod, relativePath, handlers)
}

//go:noinline
func RouterGroupHandle(group *gin.RouterGroup, httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	fmt.Println(httpMethod, relativePath)
	return nil
}

//go:noinline
func WrapRouterGroupHandle(group *gin.RouterGroup, httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {

	handlers = pushFrontHandler(group, relativePath, handlers)
	return RouterGroupHandle(group, httpMethod, relativePath, handlers...)
}

//go:noinline
func RouterGroupPOST(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	fmt.Println(relativePath)
	return nil
}

//go:noinline
func WrapRouterGroupPOST(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	handlers = pushFrontHandler(group, relativePath, handlers)
	return RouterGroupPOST(group, relativePath, handlers...)
}

//go:noinline
func RouterGroupGET(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	fmt.Println(relativePath)
	return nil
}

//go:noinline
func WrapRouterGroupGET(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	handlers = pushFrontHandler(group, relativePath, handlers)
	return RouterGroupGET(group, relativePath, handlers...)
}

//go:noinline
func RouterGroupDELETE(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	fmt.Println(relativePath)
	return nil
}

//go:noinline
func WrapRouterGroupDELETE(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	handlers = pushFrontHandler(group, relativePath, handlers)
	return RouterGroupDELETE(group, relativePath, handlers...)
}

//go:noinline
func RouterGroupPATCH(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	fmt.Println(relativePath)
	return nil
}

//go:noinline
func WrapRouterGroupPATCH(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	handlers = pushFrontHandler(group, relativePath, handlers)
	return RouterGroupPATCH(group, relativePath, handlers...)
}

//go:noinline
func RouterGroupPUT(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	fmt.Println(relativePath)
	return nil
}

//go:noinline
func WrapRouterGroupPUT(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	handlers = pushFrontHandler(group, relativePath, handlers)
	return RouterGroupPUT(group, relativePath, handlers...)
}

//go:noinline
func RouterGroupOPTIONS(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	fmt.Println(relativePath)
	return nil
}

//go:noinline
func WrapRouterGroupOPTIONS(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	handlers = pushFrontHandler(group, relativePath, handlers)
	return RouterGroupOPTIONS(group, relativePath, handlers...)
}

//go:noinline
func RouterGroupHEAD(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	fmt.Println(relativePath)
	return nil
}

//go:noinline
func WrapRouterGroupHEAD(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	handlers = pushFrontHandler(group, relativePath, handlers)
	return RouterGroupHEAD(group, relativePath, handlers...)
}

//go:noinline
func RouterGroupAny(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	fmt.Println(relativePath)
	return nil
}

//go:noinline
func WrapRouterGroupAny(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	handlers = pushFrontHandler(group, relativePath, handlers)
	return RouterGroupAny(group, relativePath, handlers...)
}

//go:noinline
func RouterGroupStaticFile(group *gin.RouterGroup, relativePath, filepath string) gin.IRoutes {
	fmt.Println(relativePath, filepath)
	return nil
}

//go:noinline
func WrapRouterGroupStaticFile(group *gin.RouterGroup, relativePath, filepath string) gin.IRoutes {
	return RouterGroupStaticFile(group, relativePath, filepath)
}

//go:noinline
func RouterGroupStatic(group *gin.RouterGroup, relativePath, root string) gin.IRoutes {
	fmt.Println(relativePath, root)
	return nil
}

//go:noinline
func WrapRouterGroupStatic(group *gin.RouterGroup, relativePath, root string) gin.IRoutes {
	return RouterGroupStatic(group, relativePath, root)
}

//go:noinline
func RouterGroupStaticFS(group *gin.RouterGroup, relativePath string, fs http.FileSystem) gin.IRoutes {
	fmt.Println(relativePath)
	return nil
}

//go:noinline
func WrapRouterGroupStaticFS(group *gin.RouterGroup, relativePath string, fs http.FileSystem) gin.IRoutes {
	return RouterGroupStaticFS(group, relativePath, fs)
}

func init() {
	tingyun3.Register(reflect.ValueOf(WrapRouterGrouphandle).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRouterGroupHandle).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRouterGroupPOST).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRouterGroupGET).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRouterGroupDELETE).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRouterGroupPATCH).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRouterGroupPUT).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRouterGroupOPTIONS).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRouterGroupAny).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRouterGroupStaticFile).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRouterGroupStatic).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRouterGroupStaticFS).Pointer())
}
