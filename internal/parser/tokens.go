package parser

import "regexp"

var tokens = []token{
	{state: stateRoot, name: "interp-end", stateChange: statePop(), skip: false, expr: regexp.MustCompile("^(}})")},
	{state: stateRoot, name: "path", stateChange: nil, skip: false, expr: regexp.MustCompile("^(\\.?/[a-zA-Z0-9._-]*)")},
	{state: stateRoot, name: "arrow", stateChange: nil, skip: false, expr: regexp.MustCompile("^(->)")},
	{state: stateRoot, name: "symbol", stateChange: nil, skip: false, expr: regexp.MustCompile("^([(){}[\\].=,$\\\\;])")},
	{state: stateRoot, name: "multiline-begin", stateChange: statePush(stateMulti), skip: false, expr: regexp.MustCompile("^('')")},
	{state: stateRoot, name: "string-begin", stateChange: statePush(stateString), skip: false, expr: regexp.MustCompile("^(\")")},
	{state: stateRoot, name: "keyword", stateChange: nil, skip: false, expr: regexp.MustCompile("^(panic|output|import)")},
	{state: stateRoot, name: "ident", stateChange: nil, skip: false, expr: regexp.MustCompile("^([a-zA-Z0-9_]+)")},
	{state: stateRoot, name: "comment", stateChange: nil, skip: true, expr: regexp.MustCompile("^(#[^\\n\\r]*)")},
	{state: stateRoot, name: "space", stateChange: nil, skip: true, expr: regexp.MustCompile("^([ \\t\\n\\r])")},
	{state: stateString, name: "interp-begin", stateChange: statePush(stateRoot), skip: false, expr: regexp.MustCompile("^({{)")},
	{state: stateString, name: "string-end", stateChange: statePop(), skip: false, expr: regexp.MustCompile("^(\")")},
	{state: stateString, name: "char", stateChange: nil, skip: false, expr: regexp.MustCompile("^(.)")},
	{state: stateMulti, name: "interp-begin", stateChange: statePush(stateRoot), skip: false, expr: regexp.MustCompile("^({{)")},
	{state: stateMulti, name: "multi-end", stateChange: statePop(), skip: false, expr: regexp.MustCompile("^('')")},
	{state: stateMulti, name: "char", stateChange: nil, skip: false, expr: regexp.MustCompile("^(.|\\s)")},
}
