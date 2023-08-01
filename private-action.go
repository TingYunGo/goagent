// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/TingYunGo/goagent/libs/tystring"
)

func (a *Action) setError(e interface{}, errType string, skipStack int, isError bool) {
	if a == nil || a.stateUsed != actionUsing {
		return
	} //errorTrace 聚合,以 callstack + message
	errTime := time.Now()
	errorStack := []string{}
	if readServerConfigBool(ServerConfigBoolErrorCollectorStackEnabled, false) {
		errorStack = callStack(skipStack)
	}
	a.errors.Put(&errInfo{errTime, patchSize(toString(e), 4000), errorStack, errType, isError})
}
func (a *Action) makeTracerID() int32 {
	return atomic.AddInt32(&a.tracerIDMaker, 1) - 1
}

func (a *Action) prefixName() string {
	prefix := "WebAction/"
	if a.isTask {
		prefix = "TaskAction/"
	} else if !strings.Contains(a.trackID, ";n=") {
		prefix = "Transaction/"
	}
	return prefix
}

func (a *Action) unicID() string {
	if len(a.actionID) == 0 {
		a.actionID = unicID(a.time.begin, a)
	}
	return a.actionID
}
func (a *Action) getTransactionID() string {
	if _, transactionID := a.parseTrackID(); len(transactionID) > 0 {
		return transactionID
	}
	return a.unicID()
}

func (a *Action) parseTrackID() (callList, transactionID string) {
	callList, transactionID = "", ""
	if parts := strings.Split(a.trackID, ";"); len(parts) > 0 {
		for _, v := range parts {
			if tystring.SubString(v, 0, 2) == "c=" {
				callList = tystring.SubString(v, 2, len(v)-2)
			} else if tystring.SubString(v, 0, 2) == "x=" {
				transactionID = tystring.SubString(v, 2, len(v)-2)
			}
		}
	}
	return
}

func namingMatchTargetValue(matchType int, matchRule string, matchString string) bool {
	if matchType == 1 {

		if matchRule != matchString {
			return false
		}
	} else if matchType == 2 {

		if subString(matchString, 0, len(matchRule)) != matchRule {
			return false
		}
	} else if matchType == 3 {

		if len(matchString) < len(matchRule) {
			return false
		}
		begin := len(matchString) - len(matchRule)
		if subString(matchString, begin, len(matchRule)) != matchRule {
			return false
		}
	} else if matchType == 4 {

		if !strings.Contains(matchString, matchRule) {
			return false
		}
	} else if matchType == 5 {

		found, err := regexp.MatchString(matchRule, matchString)
		if err != nil || !found {
			return false
		}
	} else {
		return false
	}
	return true
}
func namingMatchParamValue(matchType int, matchRule, matchString string) bool {
	if len(matchString) == 0 {
		return false
	}
	if matchType == 0 || len(matchRule) == 0 {
		return true
	}
	return namingMatchTargetValue(matchType, matchRule, matchString)
}
func namingParseURIPartName(uriRule, uri string) string {
	uriRule = trimSpace(uriRule)
	if len(uriRule) == 0 {
		return ""
	}

	ruleParts := strings.Split(uriRule, ",")
	if len(ruleParts) == 1 {
		partCount, err := strconv.Atoi(trimSpace(ruleParts[0]))
		if err != nil {
			return ""
		}
		return getStringParts(uri, '/', partCount)
	}
	uriParts := strings.Split(uri, "/")
	res := ""
	for _, v := range ruleParts {
		value := trimSpace(v)
		if len(value) > 0 {

			partId, err := strconv.Atoi(value)
			if err == nil {
				if partId > 0 && partId < len(uriParts) {
					res = res + "/" + uriParts[partId]
				}
			}
		}
	}
	return res
}
func namingCustomizeNameByRule(rule *nameRule, r *http.Request) string {

	if !rule.Valid() {
		fmt.Println("No rule")
		return ""
	}
	//匹配 HTTP Method
	if rule.matchRule.methodMatch != 0 {
		if caseCMP(rule.matchRule.MethodEnabled(), r.Method) != 0 {
			fmt.Println("Not Match method")
			return ""
		}
	}
	//匹配 URI
	if len(rule.matchRule.uriMatchTarget) > 0 {

		if !namingMatchTargetValue(rule.matchRule.howToMatchURI, rule.matchRule.uriMatchTarget, r.URL.Path) {
			fmt.Println("uri match failed")
			return ""
		}
	}
	var urlQuery url.Values = nil

	getQuery := func() url.Values {
		if urlQuery == nil {
			urlQuery = r.URL.Query()
		}
		return urlQuery
	}

	var cookies []*http.Cookie = nil
	getCookies := func(name string) string {
		if cookies == nil {
			cookies = r.Cookies()
		}
		if cookies == nil {
			return ""
		}
		for _, cookie := range cookies {
			if cookie != nil {
				if cookie.Name == name {
					return cookie.Value
				}
			}
		}
		return ""
	}

	//匹配参数
	for i := 0; i < len(rule.matchRule.params); i++ {
		paramMatch := rule.matchRule.params[i]
		if len(paramMatch.paramName) == 0 {
			continue
		}
		if paramMatch.matchType == 1 {
			v := getQuery().Get(paramMatch.paramName)

			if !namingMatchParamValue(paramMatch.howToMatchParam, paramMatch.matchTarget, v) {
				fmt.Println("param match failed")
				return ""
			}
		} else if paramMatch.matchType == 2 {
			v := r.Header.Get(paramMatch.paramName)
			if !namingMatchParamValue(paramMatch.howToMatchParam, paramMatch.matchTarget, v) {
				fmt.Println("header match failed")
				return ""
			}
		} else {
			fmt.Println("unsupported matchType", paramMatch.matchType)
			return ""
		}
	}

	//生成名字

	name := ""
	if rule.naming.method {
		name = r.Method + " "
	}
	name = name + namingParseURIPartName(rule.naming.uri, r.URL.Path)
	sepAdded := false
	appendKV := func(key, value string) {
		if !sepAdded {
			name = name + "?" + key + "=" + value
			sepAdded = true
		} else {
			name = name + "&" + key + "=" + value
		}
	}
	paramNaming := trimSpace(rule.naming.param)
	if len(paramNaming) > 0 {
		names := strings.Split(paramNaming, ",")
		for _, v := range names {
			name := trimSpace(v)
			if len(name) > 0 {
				value := getQuery().Get(name)
				if len(value) > 0 {
					appendKV(name, value)
				}
			}
		}
	}
	headerNaming := trimSpace(rule.naming.header)
	if len(headerNaming) > 0 {
		names := strings.Split(headerNaming, ",")
		for _, v := range names {
			name := trimSpace(v)
			if len(name) > 0 {
				value := r.Header.Get(name)
				if len(value) > 0 {
					appendKV(name, value)
				}
			}
		}
	}

	cookieNaming := trimSpace(rule.naming.cookie)
	if len(cookieNaming) > 0 {
		names := strings.Split(cookieNaming, ",")
		for _, v := range names {
			name := trimSpace(v)
			if len(name) > 0 {
				value := getCookies(name)
				if len(value) > 0 {
					appendKV(name, value)
				}
			}
		}
	}
	name = trimSpace(name)

	if len(name) > 0 {
		name = rule.name + "/" + name
	}

	return name
}
func namingCustomizeName(r *http.Request) string {
	rules := getServerNamingRules()
	if rules == nil {
		return ""
	}
	for _, rule := range rules.rules {
		name := namingCustomizeNameByRule(rule, r)
		if len(name) > 0 {
			return name
		}
	}
	return ""
}
