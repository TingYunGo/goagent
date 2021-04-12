// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"
)

func (c *Component) fixEnd(t *timeRange) {
	if c.time.duration == -1 {
		c.time.duration = t.duration - c.time.begin.Sub(t.begin)
		if c.time.duration < 0 {
			c.time.duration = 0
		}
	}
}
func (c *Component) destroy() {
	if c._type == componentUnused {
		return
	}
	c.name = ""
	c.method = ""
	c.txdata = ""
	c.callStack = nil
	c.action = nil
	c._type = componentUnused
	//	app.componentTemps.Put(c)
}

//汇总还要修改
func dbMetricName(name string) string {
	array := strings.Split(name, "://")
	serverDb, dburl := array[0], array[1]
	array = strings.Split(dburl, "/")
	return "Database/" + serverDb + "/" + url.QueryEscape(array[0]) + "%2F" + url.QueryEscape(array[1])
}
func nosqlMetricName(name string) string {
	array := strings.Split(name, "://")
	serverDb, dburl := array[0], array[1]
	array = strings.Split(dburl, "/")
	_, _, table, op := array[0], array[1], array[2], array[3]
	return serverDb + "/" + url.QueryEscape(table) + "/" + url.QueryEscape(op)
}
func (c *Component) getSQL() string {
	return c.sql
}
func (c *Component) isDatabaseComponent() bool {
	return c._type == ComponentMysql || c._type == ComponentPostgreSQL || c._type == ComponentDefaultDB || c._type == ComponentMSSQL || c._type == ComponentSQLite
}
func (c *Component) getURL() string {
	if c._type == ComponentExternal {
		return c.name
	}
	return ""
}
func (c *Component) externalTransaction() (metricName string, duration float64, remoteDuration float64, found bool) {
	if len(c.txdata) == 0 {
		return "", 0, 0, false
	}

	jsonData := map[string]interface{}{}
	if err := json.Unmarshal([]byte(c.txdata), &jsonData); err != nil {
		return "", 0, 0, false
	}
	if secretID, err := jsonReadString(jsonData, "id"); err != nil {
		return "", 0, 0, false
	} else if action, err := jsonReadString(jsonData, "action"); err != nil {
		return "", 0, 0, false
	} else if timeObject, err := jsonReadObjects(jsonData, "time"); err != nil {
		return "", 0, 0, false
	} else if remoteDuration, err := jsonReadFloat(timeObject, "duration"); err != nil {
		return "", 0, 0, false
	} else {
		c.extSecretID = secretID
		c.remoteDuration = remoteDuration
		return "ExternalTransaction/" + strings.Replace(c.name, "/", "%2F", -1) + "/" + secretID + "%2F" + action, float64(c.time.duration / time.Millisecond), remoteDuration, true
	}
}

//FixBegin : 校正事务开始时间
func (c *Component) FixBegin(begin time.Time) {
	c.time.begin = begin
}

func (c *Component) unicID() string {
	if c.exID {
		return unicID(c.time.begin, c)
	}
	return ""
}
