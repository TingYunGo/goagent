// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
#include <stdio.h>
typedef struct {
    const char *first;
    const char *second;
    const char *third;
} WrapItem;

typedef struct {
    const char *targetmethod; 
    const char *replacemethod;
    const char *wrapPacks;
    const char *callback;
} ReplaceItem;
extern const ReplaceItem* TingyunReplaceItems() {
    static ReplaceItem targets[] = {
     {"net/http.NotFound", ".WraphttpNotFound", "net/http", ".httpNotFound"},
     {"net/http.(*Server).Serve", ".WrapHttpServerServe", "", ".HttpServerServe"},
     {"net/http.(*ServeMux).Handle", ".WrapServerMuxHandle", "", ".ServerMuxHandle"},
     {"net/http.(*ServeMux).Handler", ".WrapServerMuxHandler", "", ".ServerMuxHandler"},
     {"net/http.(*Client).do", ".WrapHttpClientDo", "", ".HttpClientDo"},

        {"database/sql.(*Rows).Close", "/database.WrapRowsClose", "", "/database.RowsClose"},

        {"github.com/gin-gonic/gin.(*RouterGroup).handle", "/frameworks/gin.WrapRouterGrouphandle", "", "/frameworks/gin.RouterGrouphandle"},

        {"github.com/gorilla/websocket.(*Conn).NextReader", "/frameworks/websocket/gorilla.WrapConnNextReader", "", "/frameworks/websocket/gorilla.ConnNextReader"},

        {"github.com/labstack/echo.(*Echo).add", "/frameworks/echo.WrapechoEchoadd", "", "/frameworks/echo.echoEchoadd"},
        {"github.com/labstack/echo.(*Echo).Add", "/frameworks/echo.WrapechoEchoAdd", "","/frameworks/echo.echoEchoAdd"},
        {"github.com/labstack/echo.(*Router).Add", "/frameworks/echo.WrapechoRouterAdd", "","/frameworks/echo.echoRouterAdd"},
        {"github.com/labstack/echo/v4.(*Echo).add", "/frameworks/echo/v4.WrapechoEchoadd", "","/frameworks/echo/v4.echoEchoadd"},
        {"github.com/labstack/echo/v4.(*Echo).Add", "/frameworks/echo/v4.WrapechoEchoAdd", "", "/frameworks/echo/v4.echoEchoAdd"},
        {"github.com/labstack/echo/v4.(*Router).Add", "/frameworks/echo/v4.WrapechoRouterAdd", "", "/frameworks/echo/v4.echoRouterAdd"},


        {"github.com/beego/beego/v2/server/web.(*ControllerRegister).addToRouter", "/frameworks/beego/v2.WrapbeegoaddToRouter", "","/frameworks/beego/v2.beegoaddToRouter"},
        {"github.com/beego/beego/v2/server/web.(*ControllerRegister).AddMethod", "/frameworks/beego/v2.WrapbeegoAddMethod", "","/frameworks/beego/v2.beegoAddMethod"},
        {"github.com/beego/beego/v2/server/web.(*ControllerRegister).Handler", "/frameworks/beego/v2.WrapbeegoHandler", "","/frameworks/beego/v2.beegoHandler"},
        {"github.com/beego/beego/v2/server/web.(*ControllerRegister).AddAutoPrefix", "/frameworks/beego/v2.WrapbeegoAddAutoPrefix", "","/frameworks/beego/v2.beegoAddAutoPrefix"},
        {"github.com/beego/beego/v2/server/web.(*ControllerRegister).addWithMethodParams", "/frameworks/beego/v2.WrapbeegoaddWithMethodParams", "","/frameworks/beego/v2.beegoaddWithMethodParams"},
        //go mod 
        {"github.com/astaxie/beego.(*ControllerRegister).addToRouter", "/frameworks/beego.WrapbeegoaddToRouter", "","/frameworks/beego.beegoaddToRouter"},
        {"github.com/astaxie/beego.(*ControllerRegister).AddMethod", "/frameworks/beego.WrapbeegoAddMethod", "","/frameworks/beego.beegoAddMethod"},
        {"github.com/astaxie/beego.(*ControllerRegister).Handler", "/frameworks/beego.WrapbeegoHandler", "","/frameworks/beego.beegoHandler"},
        {"github.com/astaxie/beego.(*ControllerRegister).AddAutoPrefix", "/frameworks/beego.WrapbeegoAddAutoPrefix", "","/frameworks/beego.beegoAddAutoPrefix"},
        {"github.com/astaxie/beego.(*ControllerRegister).addWithMethodParams", "/frameworks/beego.WrapbeegoaddWithMethodParams", "","/frameworks/beego.beegoaddWithMethodParams"},
        //go path
        {"github.com/beego/beego/server/web.(*ControllerRegister).addToRouter", "/frameworks/beego/path.WrapbeegoaddToRouter", "","/frameworks/beego/path.beegoaddToRouter"},
        {"github.com/beego/beego/server/web.(*ControllerRegister).AddMethod", "/frameworks/beego/path.WrapbeegoAddMethod", "","/frameworks/beego/path.beegoAddMethod"},
        {"github.com/beego/beego/server/web.(*ControllerRegister).Handler", "/frameworks/beego/path.WrapbeegoHandler", "","/frameworks/beego/path.beegoHandler"},
        {"github.com/beego/beego/server/web.(*ControllerRegister).AddAutoPrefix", "/frameworks/beego/path.WrapbeegoAddAutoPrefix", "","/frameworks/beego/path.beegoAddAutoPrefix"},
        {"github.com/beego/beego/server/web.(*ControllerRegister).addWithMethodParams", "/frameworks/beego/path.WrapbeegoaddWithMethodParams", "","/frameworks/beego/path.beegoaddWithMethodParams"},
        //go path astaxie
        {"github.com/astaxie/beego/server/web.(*ControllerRegister).addToRouter", "/frameworks/beego/path/astaxie.WrapbeegoaddToRouter", "","/frameworks/beego/path/astaxie.beegoaddToRouter"},
        {"github.com/astaxie/beego/server/web.(*ControllerRegister).AddMethod", "/frameworks/beego/path/astaxie.WrapbeegoAddMethod", "","/frameworks/beego/path/astaxie.beegoAddMethod"},
        {"github.com/astaxie/beego/server/web.(*ControllerRegister).Handler", "/frameworks/beego/path/astaxie.WrapbeegoHandler", "","/frameworks/beego/path/astaxie.beegoHandler"},
        {"github.com/astaxie/beego/server/web.(*ControllerRegister).AddAutoPrefix", "/frameworks/beego/path/astaxie.WrapbeegoAddAutoPrefix", "", "/frameworks/beego/path/astaxie.beegoAddAutoPrefix"},
        {"github.com/astaxie/beego/server/web.(*ControllerRegister).addWithMethodParams", "/frameworks/beego/path/astaxie.WrapbeegoaddWithMethodParams", "","/frameworks/beego/path/astaxie.beegoaddWithMethodParams"},

       {"github.com/kataras/iris/v12/websocket.Upgrade", "/frameworks/iris/v12.WrapwebsocketUpgrade", "","/frameworks/iris/v12.websocketUpgrade"},
       {"github.com/kataras/iris/v12/mvc.(*ControllerActivator).handleMany",  "/frameworks/iris/v12.WrapirishandleMany", "","/frameworks/iris/v12.irishandleMany"},
       {"github.com/kataras/neffos.(*Conn).handleMessage",  "/frameworks/iris/v12.WrapneffosConnhandleMessage", "","/frameworks/iris/v12.neffosConnhandleMessage"},
       {"github.com/kataras/neffos.makeEventFromMethod",  "/frameworks/iris/v12.WrapneffosmakeEventFromMethod", "","/frameworks/iris/v12.neffosmakeEventFromMethod"},



        {"github.com/gomodule/redigo/redis.DialContext", "/nosql/redigo.WrapRedigoDialContext", "","/nosql/redigo.RedigoDialContext"},
        {"github.com/gomodule/redigo/redis.Dial", "/nosql/redigo.WrapredigoDial", "","/nosql/redigo.redigoDial"},
        {"github.com/gomodule/redigo/redis.(*conn).Close", "/nosql/redigo.WrapRedigoConnClose", "","/nosql/redigo.RedigoConnClose"},
        {"github.com/gomodule/redigo/redis.(*conn).DoWithTimeout", "/nosql/redigo.WrapRedigoDoWithTimeout", "","/nosql/redigo.RedigoDoWithTimeout"},

        {"github.com/go-redis/redis.(*baseClient).process", "/nosql/go-redis.WrapbaseClientprocess", "","/nosql/go-redis.baseClientprocess"},
        {"github.com/go-redis/redis.(*baseClient).Process", "/nosql/go-redis.WrapbaseClientProcess", "","/nosql/go-redis.baseClientProcess"},
        {"github.com/go-redis/redis.(*baseClient).processPipeline", "/nosql/go-redis.WrapbaseClientprocessPipeline", "","/nosql/go-redis.baseClientprocessPipeline"},
        {"github.com/go-redis/redis.(*baseClient).generalProcessPipeline", "/nosql/go-redis.WrapbaseClientgeneralProcessPipeline", "","/nosql/go-redis.baseClientgeneralProcessPipeline"},

        {"github.com/go-redis/redis.NewClient", "/nosql/go-redis.WrapredisNewClient", "","/nosql/go-redis.redisNewClient"},
        {"github.com/go-redis/redis.NewClusterClient", "/nosql/go-redis.WrapredisNewClusterClient", "","/nosql/go-redis.redisNewClusterClient"},

        {"github.com/go-redis/redis.(*Client).WrapProcess", "/nosql/go-redis.WrapredisClientWrapProcess", "","/nosql/go-redis.redisClientWrapProcess"},
        {"github.com/go-redis/redis.(*Client).WrapProcessPipeline", "/nosql/go-redis.WrapredisClientWrapProcessPipeline", "","/nosql/go-redis.redisClientWrapProcessPipeline"},
        {"github.com/go-redis/redis.(*ClusterClient).WrapProcess", "/nosql/go-redis.WrapredisClusterClientWrapProcess", "","/nosql/go-redis.redisClusterClientWrapProcess"},
        {"github.com/go-redis/redis.(*ClusterClient).WrapProcessPipeline", "/nosql/go-redis.WrapredisClusterClientWrapProcessPipeline", "","/nosql/go-redis.redisClusterClientWrapProcessPipeline"},

        {"github.com/go-redis/redis/v7.(*baseClient).process", "/nosql/go-redis/v7.WrapbaseClientprocess", "","/nosql/go-redis/v7.baseClientprocess"},
        {"github.com/go-redis/redis/v7.(*baseClient).processPipeline", "/nosql/go-redis/v7.WrapbaseClientprocessPipeline", "","/nosql/go-redis/v7.baseClientprocessPipeline"},
        {"github.com/go-redis/redis/v7.(*baseClient).generalProcessPipeline", "/nosql/go-redis/v7.WrapbaseClientgeneralProcessPipeline", "","/nosql/go-redis/v7.baseClientgeneralProcessPipeline"},
        {"github.com/go-redis/redis/v8.(*baseClient).process", "/nosql/go-redis/v8.WrapbaseClientprocess", "","/nosql/go-redis/v8.baseClientprocess"},
        {"github.com/go-redis/redis/v8.(*baseClient).processPipeline", "/nosql/go-redis/v8.WrapbaseClientprocessPipeline", "","/nosql/go-redis/v8.baseClientprocessPipeline"},
        {"github.com/go-redis/redis/v8.(*baseClient).generalProcessPipeline", "/nosql/go-redis/v8.WrapbaseClientgeneralProcessPipeline", "","/nosql/go-redis/v8.baseClientgeneralProcessPipeline"},

        {"go.mongodb.org/mongo-driver/mongo.(*Collection).BulkWrite", "/nosql/mongodb.WrapmongodbBulkWrite", "","/nosql/mongodb.mongodbBulkWrite"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).InsertOne", "/nosql/mongodb.WrapmongodbInsertOne", "","/nosql/mongodb.mongodbInsertOne"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).InsertMany", "/nosql/mongodb.WrapmongodbInsertMany", "","/nosql/mongodb.mongodbInsertMany"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).DeleteOne", "/nosql/mongodb.WrapmongodbDeleteOne", "","/nosql/mongodb.mongodbDeleteOne"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).DeleteMany", "/nosql/mongodb.WrapmongodbDeleteMany", "","/nosql/mongodb.mongodbDeleteMany"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).UpdateMany", "/nosql/mongodb.WrapmongodbUpdateMany", "","/nosql/mongodb.mongodbUpdateMany"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).UpdateByID", "/nosql/mongodb.WrapmongodbUpdateByID", "","/nosql/mongodb.mongodbUpdateByID"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).UpdateOne", "/nosql/mongodb.WrapmongodbUpdateOne", "","/nosql/mongodb.mongodbUpdateOne"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).ReplaceOne", "/nosql/mongodb.WrapmongodbReplaceOne", "","/nosql/mongodb.mongodbReplaceOne"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).CountDocuments", "/nosql/mongodb.WrapmongodbCountDocuments", "","/nosql/mongodb.mongodbCountDocuments"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).EstimatedDocumentCount", "/nosql/mongodb.WrapmongodbEstimatedDocumentCount", "","/nosql/mongodb.mongodbEstimatedDocumentCount"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).Distinct", "/nosql/mongodb.WrapmongodbDistinct", "","/nosql/mongodb.mongodbDistinct"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).Find", "/nosql/mongodb.WrapmongodbFind", "","/nosql/mongodb.mongodbFind"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).FindOne", "/nosql/mongodb.WrapmongodbFindOne", "","/nosql/mongodb.mongodbFindOne"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).FindOneAndDelete", "/nosql/mongodb.WrapmongodbFindOneAndDelete", "","/nosql/mongodb.mongodbFindOneAndDelete"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).FindOneAndReplace", "/nosql/mongodb.WrapmongodbFindOneAndReplace", "","/nosql/mongodb.mongodbFindOneAndReplace"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).FindOneAndUpdate", "/nosql/mongodb.WrapmongodbFindOneAndUpdate", "","/nosql/mongodb.mongodbFindOneAndUpdate"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).Watch", "/nosql/mongodb.WrapmongodbWatch", "","/nosql/mongodb.mongodbWatch"},
        {"go.mongodb.org/mongo-driver/mongo.(*Collection).Drop", "/nosql/mongodb.WrapmongodbDrop", "","/nosql/mongodb.mongodbDrop"},
        {"go.mongodb.org/mongo-driver/mongo.NewClient", "/nosql/mongodb.WrapmongodbNewClient", "","/nosql/mongodb.mongodbNewClient"},
        {"go.mongodb.org/mongo-driver/mongo.(*Client).Disconnect", "/nosql/mongodb.WrapmongodbDisconnect", "","/nosql/mongodb.mongodbDisconnect"},

        {0, 0, 0, 0}
    };
    return targets;
}

extern const WrapItem* TingyunWraps() {
    
    static WrapItem targets[] = {
        {"database/sql.Open", "/database.", "DBOpen"},
        {"database/sql.(*DB).Close", "/database.", "DBClose"},
        {"database/sql.(*DB).queryDC", "/database.", "DBqueryDC"},
        {"database/sql.(*DB).execDC", "/database.", "DBexecDC"},
        {"database/sql.(*DB).prepareDC", "/database.", "DBprepareDC"},
        {"database/sql.(*Stmt).QueryContext", "/database.", "StmtQueryContext"},
        {"database/sql.(*Stmt).ExecContext", "/database.", "StmtExecContext"},
        {"database/sql.(*Stmt).Close", "/database.", "StmtClose"},

        
        {"github.com/kataras/iris/v12/core/router.(*APIBuilder).CreateRoutes", "/frameworks/iris/v12.", "irisCreateRoutes"},
        {"github.com/kataras/iris/v12/core/router.FileServer", "/frameworks/iris/v12.", "routerFileServer"},

        {"github.com/go-redis/redis/v7.NewClient", "/nosql/go-redis/v7.", "redisNewClient"},
        {"github.com/go-redis/redis/v7.NewClusterClient", "/nosql/go-redis/v7.", "redisNewClusterClient"},

        {"github.com/go-redis/redis/v8.NewClient", "/nosql/go-redis/v8.", "redisNewClient"},
        {"github.com/go-redis/redis/v8.NewClusterClient", "/nosql/go-redis/v8.", "redisNewClusterClient"},

        {0, 0, 0}
    };
    return targets;
}
