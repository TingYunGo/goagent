#include <stdio.h>
typedef struct {
    const char *first;
    const char *second;
    const char *third;
} WrapItem;

extern const WrapItem* TingyunWraps() {
    
    static WrapItem targets[] = {
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
        {"github.com/gomodule/redigo/redis.DialContext", "/nosql/redigo.", "RedigoDialContext"},
        {"github.com/gomodule/redigo/redis.(*conn).Close", "/nosql/redigo.", "RedigoConnClose"},
        {"github.com/gomodule/redigo/redis.(*conn).DoWithTimeout", "/nosql/redigo.", "RedigoDoWithTimeout"},
        {"github.com/go-redis/redis.NewClient", "/nosql/go-redis.", "GoRedisNewClient"},
        {"github.com/go-redis/redis.(*baseClient).process", "/nosql/go-redis.", "baseClientprocess"},
        {"github.com/go-redis/redis.(*baseClient).Process", "/nosql/go-redis.", "baseClientProcess"},
        {"github.com/go-redis/redis.(*baseClient).processPipeline", "/nosql/go-redis.", "baseClientprocessPipeline"},
        {0, 0}
    };
    return targets;
}