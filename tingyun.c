// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
#include <stdio.h>
typedef struct {
    const char *first;
    const char *second;
    const char *third;
} WrapItem;

extern const WrapItem* TingyunWraps() {
    
    static WrapItem targets[] = {
        {"net/http.NotFound", ".", "httpNotFound"},
        {"net/http.(*ServeMux).Handle", ".", "ServerMuxHandle"},
        {"net/http.(*Server).Serve", ".", "ServerServe"},
        {"net/http.(*Client).do", ".", "HttpClientDo"},

        {"database/sql.Open", "/database.", "DBOpen"},
        {"database/sql.(*DB).Close", "/database.", "DBClose"},
        {"database/sql.(*Conn).PrepareContext", "/database.", "ConnPrepareContext"},
        {"database/sql.(*Conn).QueryContext", "/database.", "ConnQueryContext"},
        {"database/sql.(*Conn).ExecContext", "/database.", "ConnExecContext"},
        {"database/sql.(*DB).PrepareContext", "/database.", "DBPrepareContext"},
        {"database/sql.(*DB).QueryContext", "/database.", "DBQueryContext"},
        {"database/sql.(*DB).ExecContext", "/database.", "DBExecContext"},
        {"database/sql.(*Tx).PrepareContext", "/database.", "TxPrepareContext"},
        {"database/sql.(*Tx).QueryContext", "/database.", "TxQueryContext"},
        {"database/sql.(*Tx).ExecContext", "/database.", "TxExecContext"},
        {"database/sql.(*Stmt).QueryContext", "/database.", "StmtQueryContext"},
        {"database/sql.(*Stmt).ExecContext", "/database.", "StmtExecContext"},
        {"database/sql.(*Stmt).Close", "/database.", "StmtClose"},
        {"database/sql.(*Rows).Close", "/database.", "RowsClose"},

        {"github.com/gin-gonic/gin.(*RouterGroup).handle", "/frameworks/gin.", "RouterGrouphandle"},
        {"github.com/labstack/echo.(*Echo).add", "/frameworks/echo.", "echoEchoadd"},
        {"github.com/labstack/echo.(*Echo).Add", "/frameworks/echo.", "echoEchoAdd"},
        {"github.com/labstack/echo.(*Router).Add", "/frameworks/echo.", "echoRouterAdd"},
        {"github.com/labstack/echo/v4.(*Echo).add", "/frameworks/echo/v4.", "echoEchoadd"},
        {"github.com/labstack/echo/v4.(*Echo).Add", "/frameworks/echo/v4.", "echoEchoAdd"},
        {"github.com/labstack/echo/v4.(*Router).Add", "/frameworks/echo/v4.", "echoRouterAdd"},
        //go mod v2
        {"github.com/beego/beego/v2/server/web.(*ControllerRegister).addToRouter", "/frameworks/beego/v2.", "beegoaddToRouter"},
        {"github.com/beego/beego/v2/server/web.(*ControllerRegister).AddMethod", "/frameworks/beego/v2.", "beegoAddMethod"},
        {"github.com/beego/beego/v2/server/web.(*ControllerRegister).Handler", "/frameworks/beego/v2.", "beegoHandler"},
        {"github.com/beego/beego/v2/server/web.(*ControllerRegister).AddAutoPrefix", "/frameworks/beego/v2.", "beegoAddAutoPrefix"},
        {"github.com/beego/beego/v2/server/web.(*ControllerRegister).addWithMethodParams", "/frameworks/beego/v2.", "beegoaddWithMethodParams"},
        //go mod 
        {"github.com/astaxie/beego.(*ControllerRegister).addToRouter", "/frameworks/beego.", "beegoaddToRouter"},
        {"github.com/astaxie/beego.(*ControllerRegister).AddMethod", "/frameworks/beego.", "beegoAddMethod"},
        {"github.com/astaxie/beego.(*ControllerRegister).Handler", "/frameworks/beego.", "beegoHandler"},
        {"github.com/astaxie/beego.(*ControllerRegister).AddAutoPrefix", "/frameworks/beego.", "beegoAddAutoPrefix"},
        {"github.com/astaxie/beego.(*ControllerRegister).addWithMethodParams", "/frameworks/beego.", "beegoaddWithMethodParams"},
        //go path
        {"github.com/beego/beego/server/web.(*ControllerRegister).addToRouter", "/frameworks/beego/path.", "beegoaddToRouter"},
        {"github.com/beego/beego/server/web.(*ControllerRegister).AddMethod", "/frameworks/beego/path.", "beegoAddMethod"},
        {"github.com/beego/beego/server/web.(*ControllerRegister).Handler", "/frameworks/beego/path.", "beegoHandler"},
        {"github.com/beego/beego/server/web.(*ControllerRegister).AddAutoPrefix", "/frameworks/beego/path.", "beegoAddAutoPrefix"},
        {"github.com/beego/beego/server/web.(*ControllerRegister).addWithMethodParams", "/frameworks/beego/path.", "beegoaddWithMethodParams"},
        //go path astaxie
        {"github.com/astaxie/beego/server/web.(*ControllerRegister).addToRouter", "/frameworks/beego/path/astaxie.", "beegoaddToRouter"},
        {"github.com/astaxie/beego/server/web.(*ControllerRegister).AddMethod", "/frameworks/beego/path/astaxie.", "beegoAddMethod"},
        {"github.com/astaxie/beego/server/web.(*ControllerRegister).Handler", "/frameworks/beego/path/astaxie.", "beegoHandler"},
        {"github.com/astaxie/beego/server/web.(*ControllerRegister).AddAutoPrefix", "/frameworks/beego/path/astaxie.", "beegoAddAutoPrefix"},
        {"github.com/astaxie/beego/server/web.(*ControllerRegister).addWithMethodParams", "/frameworks/beego/path/astaxie.", "beegoaddWithMethodParams"},
        
        {"github.com/kataras/iris/v12/core/router.(*APIBuilder).CreateRoutes", "/frameworks/iris/v12.", "irisCreateRoutes"},
        {"github.com/kataras/iris/v12/core/router.FileServer", "/frameworks/iris/v12.", "routerFileServer"},
        {"github.com/kataras/iris/v12/websocket.Upgrade", "/frameworks/iris/v12.", "websocketUpgrade"},
        {"github.com/kataras/iris/v12/mvc.(*ControllerActivator).handleMany",  "/frameworks/iris/v12.", "irishandleMany"},
        {"github.com/kataras/neffos.(*Conn).handleMessage",  "/frameworks/iris/v12.", "neffosConnhandleMessage"},
        {"github.com/kataras/neffos.makeEventFromMethod",  "/frameworks/iris/v12.", "neffosmakeEventFromMethod"},

        {"github.com/kataras/iris/v12/core/router.(*APIBuilder).CreateRoutes", "/frameworks/iris/v12/2.", "irisCreateRoutes"},
        {"github.com/kataras/iris/v12/websocket.Upgrade", "/frameworks/iris/v12/2.", "websocketUpgrade"},
        {"github.com/kataras/iris/v12/mvc.(*ControllerActivator).handleMany",  "/frameworks/iris/v12/2.", "irishandleMany"},

        {"github.com/kataras/neffos.(*Conn).handleMessage",  "/frameworks/iris/v12/2.", "neffosConnhandleMessage"},
        {"github.com/kataras/neffos.Events.fireEvent",  "/frameworks/iris/v12/2.", "neffosfireEvent"},
        {"github.com/kataras/neffos.makeEventFromMethod",  "/frameworks/iris/v12/2.", "neffosmakeEventFromMethod"},

        {"github.com/gomodule/redigo/redis.DialContext", "/nosql/redigo.", "RedigoDialContext"},
        {"github.com/gomodule/redigo/redis.Dial", "/nosql/redigo.", "redigoDial"},
        {"github.com/gomodule/redigo/redis.(*conn).Close", "/nosql/redigo.", "RedigoConnClose"},
        {"github.com/gomodule/redigo/redis.(*conn).DoWithTimeout", "/nosql/redigo.", "RedigoDoWithTimeout"},

        {"github.com/go-redis/redis.(*baseClient).process", "/nosql/go-redis.", "baseClientprocess"},
        {"github.com/go-redis/redis.(*baseClient).Process", "/nosql/go-redis.", "baseClientProcess"},
        {"github.com/go-redis/redis.(*baseClient).processPipeline", "/nosql/go-redis.", "baseClientprocessPipeline"},
        {"github.com/go-redis/redis/v7.(*baseClient).process", "/nosql/go-redis/v7.", "baseClientprocess"},
        {"github.com/go-redis/redis/v7.(*baseClient).processPipeline", "/nosql/go-redis/v7.", "baseClientprocessPipeline"},
        {"github.com/go-redis/redis/v7.(*baseClient).generalProcessPipeline", "/nosql/go-redis/v7.", "baseClientgeneralProcessPipeline"},
        {"github.com/go-redis/redis/v8.(*baseClient).process", "/nosql/go-redis/v8.", "baseClientprocess"},
        {"github.com/go-redis/redis/v8.(*baseClient).processPipeline", "/nosql/go-redis/v8.", "baseClientprocessPipeline"},
        {"github.com/go-redis/redis/v8.(*baseClient).generalProcessPipeline", "/nosql/go-redis/v8.", "baseClientgeneralProcessPipeline"},

        {"go.mongodb.org/mongo-driver/mongo.(*Collection).BulkWrite", "/nosql/mongodb.", "mongodbBulkWrite"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).InsertOne", "/nosql/mongodb.", "mongodbInsertOne"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).InsertMany", "/nosql/mongodb.", "mongodbInsertMany"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).DeleteOne", "/nosql/mongodb.", "mongodbDeleteOne"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).DeleteMany", "/nosql/mongodb.", "mongodbDeleteMany"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).UpdateMany", "/nosql/mongodb.", "mongodbUpdateMany"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).UpdateByID", "/nosql/mongodb.", "mongodbUpdateByID"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).UpdateOne", "/nosql/mongodb.", "mongodbUpdateOne"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).ReplaceOne", "/nosql/mongodb.", "mongodbReplaceOne"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).CountDocuments", "/nosql/mongodb.", "mongodbCountDocuments"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).EstimatedDocumentCount", "/nosql/mongodb.", "mongodbEstimatedDocumentCount"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).Distinct", "/nosql/mongodb.", "mongodbDistinct"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).Find", "/nosql/mongodb.", "mongodbFind"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).FindOne", "/nosql/mongodb.", "mongodbFindOne"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).FindOneAndDelete", "/nosql/mongodb.", "mongodbFindOneAndDelete"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).FindOneAndReplace", "/nosql/mongodb.", "mongodbFindOneAndReplace"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).FindOneAndUpdate", "/nosql/mongodb.", "mongodbFindOneAndUpdate"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).Watch", "/nosql/mongodb.", "mongodbWatch"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).Drop", "/nosql/mongodb.", "mongodbDrop"},
        {"go.mongodb.org/mongo-driver/mongo.NewClient", "/nosql/mongodb.", "mongodbNewClient"},
        {"go.mongodb.org/mongo-driver/mongo.(*Client).Disconnect", "/nosql/mongodb.", "mongodbDisconnect"},
        {0, 0, 0}
    };
    return targets;
}
