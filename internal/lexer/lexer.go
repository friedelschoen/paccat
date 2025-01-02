package lexer

import (
	"fmt"
	"regexp"
)

//go:generate python3 gentokens.py tokens.txt tokens.go

type state string

const (
	stateRoot   = ""
	stateMulti  = "multiline-string"
	stateString = "string"
)

type stateFunc func([]state) []state

type token struct {
	state       state
	name        string
	stateChange stateFunc
	skip        bool
	expr        *regexp.Regexp
}

func statePop() stateFunc {
	return func(in []state) []state {
		return in[1:]
	}
}

func statePush(s state) stateFunc {
	return func(in []state) []state {
		return append([]state{s}, in...)
	}
}

type Token struct {
	Start, End    int
	Name, Content string
}

type Tokenizer struct {
	current   []state
	save      int
	savestate []state

	Text  string
	Pos   int
	Valid bool
	Token Token
}

func NewTokenizer(text string) *Tokenizer {
	return &Tokenizer{
		Text:    text,
		current: []state{stateRoot},
	}
}

func (this *Tokenizer) Next() bool {
	if len(this.current) == 0 {
		if len(this.current) == 0 {
			this.Token = Token{
				Start:   this.Pos,
				End:     this.Pos + 1,
				Name:    "illegal",
				Content: "empty state",
			}
			this.Valid = false
			return false
		}
	}

	if this.Pos >= len(this.Text) {
		if len(this.current) != 1 {
			this.Token = Token{
				Start:   this.Pos - 1,
				End:     this.Pos,
				Name:    "illegal",
				Content: fmt.Sprintf("unclosed %s", this.current[0]),
			}
		} else {
			this.Token = Token{
				Start: this.Pos - 1,
				End:   this.Pos,
				Name:  "eof",
			}
		}
		this.Valid = false
		return false
	}

	for _, tok := range tokens {
		if tok.state != this.current[0] {
			continue
		}

		loc := tok.expr.FindStringIndex(this.Text[this.Pos:])
		if loc == nil {
			continue
		}
		length := loc[1] // loc[0] is always 0
		if tok.stateChange != nil {
			this.current = tok.stateChange(this.current)
		}
		content := this.Text[this.Pos : this.Pos+length]
		this.Pos += length

		if tok.skip {
			return this.Next()
		}
		this.Token = Token{
			Start:   this.Pos - length,
			End:     this.Pos,
			Name:    tok.name,
			Content: string(content),
		}
		this.Valid = true
		return true
	}

	this.Token = Token{
		Start:   this.Pos,
		End:     this.Pos + 1,
		Name:    "illegal",
		Content: fmt.Sprintf("illegal character %c", this.Text[this.Pos]),
	}
	this.Valid = false
	return false
}

func (this *Tokenizer) Reset() {
	this.Pos = 0
	this.current = []state{stateRoot}
}

func (this *Tokenizer) Save() {
	this.save = this.Pos
	this.savestate = this.current
}

func (this *Tokenizer) Load() {
	this.Pos = this.save
	this.current = this.savestate
}
