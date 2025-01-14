package parser

import (
	"regexp"
	"strings"

	"friedelschoen.io/paccat/internal/util"
)

type stateFunc func([]state) []state

func statePop() stateFunc {
	return func(in []state) []state {
		return in[1:]
	}
}

func statePush(s state) stateFunc {
	return func(in []state) []state {
		return util.Prepend(in, s)
	}
}

func regexTest(expr string) func(string) int {
	reg := regexp.MustCompile("^(" + expr + ")")
	return func(s string) int {
		loc := reg.FindStringIndex(s)
		if loc == nil {
			return 0
		}
		return loc[1]
	}
}

func literalTest(literals ...string) func(string) int {
	return func(s string) int {
		for _, lit := range literals {
			if strings.HasPrefix(s, lit) {
				return len(lit)
			}
		}
		return 0
	}
}

var tokens = []tokenDefine{
	{state: "root", name: "interp-end", stateChange: statePop(), expr: literalTest("}}")},
	{state: "root", name: "path", stateChange: nil, expr: regexTest("\\.{0,2}/[a-zA-Z0-9._-]*")},
	{state: "root", name: "arrow", stateChange: nil, expr: literalTest("->")},
	{state: "root", name: "symbol", stateChange: nil, expr: regexTest("[#(){}[\\].=,\\\\;]")},
	{state: "root", name: "number", stateChange: nil, expr: regexTest("[0-9]+")},
	{state: "root", name: "multiline-begin", stateChange: statePush("multi"), expr: literalTest("''")},
	{state: "root", name: "string-begin", stateChange: statePush("string"), expr: literalTest("\"")},
	{state: "root", name: "keyword", stateChange: nil, expr: literalTest("panic", "output", "import")},
	{state: "root", name: "ident", stateChange: nil, expr: regexTest("[a-zA-Z0-9_]+")},
	{state: "root", name: "", stateChange: nil, expr: regexTest("//[^\\n\\r]*")},
	{state: "root", name: "", stateChange: nil, expr: regexTest("/\\*(\\s|.)*\\*/")},
	{state: "root", name: "", stateChange: nil, expr: regexTest("[ \\t\\n\\r]")},
	{state: "string", name: "interp-begin", stateChange: statePush("root"), expr: literalTest("{{")},
	{state: "string", name: "string-end", stateChange: statePop(), expr: literalTest("\"")},
	{state: "string", name: "char", stateChange: nil, expr: regexTest(".")},
	{state: "multi", name: "interp-begin", stateChange: statePush("root"), expr: literalTest("{{")},
	{state: "multi", name: "multi-end", stateChange: statePop(), expr: literalTest("''")},
	{state: "multi", name: "char", stateChange: nil, expr: regexTest(".|\\s")},
}
