// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package database

import (
	"fmt"
	"net/url"
	nurl "net/url"
	"strings"
	"unicode"

	"git.codemonky.net/TingYunGo/goagent/libs/tystring"
)

type scanner struct {
	s []rune
	i int
}

// newScanner returns a new scanner initialized with the option string s.
func newScanner(s string) *scanner {
	return &scanner{[]rune(s), 0}
}

// Next returns the next rune.
// It returns 0, false if the end of the text has been reached.
func (s *scanner) Next() (rune, bool) {
	if s.i >= len(s.s) {
		return 0, false
	}
	r := s.s[s.i]
	s.i++
	return r, true
}

// SkipSpaces returns the next non-whitespace rune.
// It returns 0, false if the end of the text has been reached.
func (s *scanner) SkipSpaces() (rune, bool) {
	r, ok := s.Next()
	for unicode.IsSpace(r) && ok {
		r, ok = s.Next()
	}
	return r, ok
}
func parsePostgreOpts(name string, o map[string]string) error {
	s := newScanner(name)

	for {
		var (
			keyRunes, valRunes []rune
			r                  rune
			ok                 bool
		)

		if r, ok = s.SkipSpaces(); !ok {
			break
		}

		// Scan the key
		for !unicode.IsSpace(r) && r != '=' {
			keyRunes = append(keyRunes, r)
			if r, ok = s.Next(); !ok {
				break
			}
		}

		// Skip any whitespace if we're not at the = yet
		if r != '=' {
			r, ok = s.SkipSpaces()
		}

		// The current character should be =
		if r != '=' || !ok {
			return fmt.Errorf(`missing "=" after %q in connection info string"`, string(keyRunes))
		}

		// Skip any whitespace after the =
		if r, ok = s.SkipSpaces(); !ok {
			// If we reach the end here, the last value is just an empty string as per libpq.
			o[string(keyRunes)] = ""
			break
		}

		if r != '\'' {
			for !unicode.IsSpace(r) {
				if r == '\\' {
					if r, ok = s.Next(); !ok {
						return fmt.Errorf(`missing character after backslash`)
					}
				}
				valRunes = append(valRunes, r)

				if r, ok = s.Next(); !ok {
					break
				}
			}
		} else {
		quote:
			for {
				if r, ok = s.Next(); !ok {
					return fmt.Errorf(`unterminated quoted string literal in connection string`)
				}
				switch r {
				case '\'':
					break quote
				case '\\':
					r, _ = s.Next()
					fallthrough
				default:
					valRunes = append(valRunes, r)
				}
			}
		}

		o[string(keyRunes)] = string(valRunes)
	}

	return nil
}

func parseMyDSN(dsn string) (host, db string) {
	hostSize := strings.LastIndex(dsn, "/")
	if hostSize == -1 {
		return dsn, ""
	}
	host = dsn[0:hostSize]
	db = dsn[hostSize+1:]
	if dbSize := strings.LastIndex(db, "?"); dbSize != -1 {
		db = db[0:dbSize]
	}
	db = trimDBName(db)
	return
}

type databaseInfo struct {
	vender string
	dsn    string
	host   string
	dbname string
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func (i *databaseInfo) init(vender, dsn string) {
	i.vender = vender
	i.dsn = dsn
	if tystring.CaseCMP(tystring.SubString(i.vender, 0, 5), "mysql") == 0 {
		if id := strings.Index(dsn, "@"); id != -1 {
			i.dsn = dsn[id+1:]
		}
	}
	i.host, i.dbname = parseDSN(i.vender, i.dsn)
}
func trimDBName(name string) string {
	begin := 0
	for begin < len(name) {
		if name[begin:begin+1] != "/" {
			break
		}
		begin++
	}
	return name[begin:]
}
func parseDSN(vender, dsn string) (host, db string) {
	if tystring.CaseCMP(vender[0:min(len(vender), 5)], "mysql") == 0 {
		return parseMyDSN(dsn)
	}
	if tystring.CaseCMP(vender[0:min(len(vender), 6)], "sqlite") == 0 {
		return "file", url.PathEscape(dsn)
	}
	u, err := nurl.Parse(dsn)
	if err == nil {
		dbname := u.Path
		if dbname == "" {
			q := u.Query()
			if v, found := q["database"]; found && len(v) > 0 {
				dbname = v[0]
			} else if v, found := q["dbname"]; found && len(v) > 0 {
				dbname = v[0]
			}
		}
		return u.Host, trimDBName(dbname)
	}

	if tystring.CaseCMP(vender[0:min(len(vender), 7)], "postgre") == 0 {
		values := map[string]string{}
		if parsePostgreOpts(dsn, values) == nil {
			if h, found := values["host"]; found {
				host = h
				if p, found := values["port"]; found {
					host = host + ":" + p
				}
			}
			if dbname, found := values["dbname"]; found {
				db = trimDBName(dbname)
			}
			return
		}
	}
	return "", ""
}
