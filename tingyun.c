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
        {"github.com/gomodule/redigo/redis.DialContext", "/nosql/redigo.", "RedigoDialContext"},
        {"github.com/gomodule/redigo/redis.(*conn).Close", "/nosql/redigo.", "RedigoConnClose"},
        {"github.com/gomodule/redigo/redis.(*conn).DoWithTimeout", "/nosql/redigo.", "RedigoDoWithTimeout"},
        {"github.com/go-redis/redis.(*baseClient).process", "/nosql/go-redis.", "baseClientprocess"},
        {"github.com/go-redis/redis.(*baseClient).Process", "/nosql/go-redis.", "baseClientProcess"},
        {"github.com/go-redis/redis.(*baseClient).processPipeline", "/nosql/go-redis.", "baseClientprocessPipeline"},
        {"github.com/go-redis/redis/v7.(*baseClient).process", "/nosql/go-redis/v7.", "baseClientprocess"},
        {"github.com/go-redis/redis/v7.(*baseClient).processPipeline", "/nosql/go-redis/v7.", "baseClientprocessPipeline"},
        {"github.com/go-redis/redis/v8.(*baseClient).process", "/nosql/go-redis/v8.", "baseClientprocess"},
        {"github.com/go-redis/redis/v8.(*baseClient).processPipeline", "/nosql/go-redis/v8.", "baseClientprocessPipeline"},
        {0, 0, 0}
    };
    return targets;
}
