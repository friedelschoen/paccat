package parser

import (
	"regexp"

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

var tokens = []tokenDefine{
	{state: "root", name: "interp-end", stateChange: statePop(), skip: false, expr: regexp.MustCompile("^(}})")},
	{state: "root", name: "path", stateChange: nil, skip: false, expr: regexp.MustCompile("^(\\.{0,2}/[a-zA-Z0-9._-]*)")},
	{state: "root", name: "arrow", stateChange: nil, skip: false, expr: regexp.MustCompile("^(->)")},
	{state: "root", name: "symbol", stateChange: nil, skip: false, expr: regexp.MustCompile("^([(){}[\\].=,\\\\;])")},
	{state: "root", name: "multiline-begin", stateChange: statePush("multi"), skip: false, expr: regexp.MustCompile("^('')")},
	{state: "root", name: "string-begin", stateChange: statePush("string"), skip: false, expr: regexp.MustCompile("^(\")")},
	{state: "root", name: "keyword", stateChange: nil, skip: false, expr: regexp.MustCompile("^(panic|output|import)")},
	{state: "root", name: "ident", stateChange: nil, skip: false, expr: regexp.MustCompile("^([a-zA-Z0-9_]+)")},
	{state: "root", name: "comment", stateChange: nil, skip: true, expr: regexp.MustCompile("^(#[^\\n\\r]*)")},
	{state: "root", name: "space", stateChange: nil, skip: true, expr: regexp.MustCompile("^([ \\t\\n\\r])")},
	{state: "string", name: "interp-begin", stateChange: statePush("root"), skip: false, expr: regexp.MustCompile("^({{)")},
	{state: "string", name: "string-end", stateChange: statePop(), skip: false, expr: regexp.MustCompile("^(\")")},
	{state: "string", name: "char", stateChange: nil, skip: false, expr: regexp.MustCompile("^(.)")},
	{state: "multi", name: "interp-begin", stateChange: statePush("root"), skip: false, expr: regexp.MustCompile("^({{)")},
	{state: "multi", name: "multi-end", stateChange: statePop(), skip: false, expr: regexp.MustCompile("^('')")},
	{state: "multi", name: "char", stateChange: nil, skip: false, expr: regexp.MustCompile("^(.|\\s)")},
}
