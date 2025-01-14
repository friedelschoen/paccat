package parser

import (
	"fmt"

	"friedelschoen.io/paccat/internal/errors"
)

type state string

type testFunc func(string) int

type tokenDefine struct {
	state       state
	name        string
	stateChange stateFunc
	expr        testFunc
}

type Token struct {
	Pos           errors.Position
	Name, Content string
}

func (this Token) GetPosition() errors.Position {
	return this.Pos
}

type TokenizerState struct {
	pos   int
	state []state
	token Token
}

type Tokenizer struct {
	current []state

	File  *errors.ErrorFile
	Pos   int
	Valid bool
	Token Token
}

func (this *Tokenizer) Next() bool {
	if len(this.current) == 0 {
		if len(this.current) == 0 {
			this.Token = Token{
				Pos: errors.Position{
					File:  this.File,
					Start: this.Pos,
					End:   this.Pos + 1,
				},
				Name:    "illegal",
				Content: "empty state",
			}
			this.Valid = false
			return false
		}
	}

	if this.Pos >= len(this.File.Content) {
		if len(this.current) != 1 {
			this.Token = Token{
				Pos: errors.Position{
					File:  this.File,
					Start: this.Pos - 1,
					End:   this.Pos,
				},
				Name:    "illegal",
				Content: fmt.Sprintf("unclosed %s", this.current[0]),
			}
		} else {
			this.Token = Token{
				Pos: errors.Position{
					File:  this.File,
					Start: this.Pos - 1,
					End:   this.Pos,
				},
				Name: "eof",
			}
		}
		this.Valid = false
		return false
	}

	for _, tok := range tokens {
		if tok.state != this.current[0] {
			continue
		}

		length := tok.expr(this.File.Content[this.Pos:])
		if length == 0 {
			continue
		}
		if tok.stateChange != nil {
			this.current = tok.stateChange(this.current)
		}
		content := this.File.Content[this.Pos : this.Pos+length]
		this.Pos += length

		if len(tok.name) == 0 {
			return this.Next()
		}
		this.Token = Token{
			Pos: errors.Position{
				File:  this.File,
				Start: this.Pos - length,
				End:   this.Pos,
			},
			Name:    tok.name,
			Content: string(content),
		}
		this.Valid = true
		return true
	}

	this.Token = Token{
		Pos: errors.Position{
			File:  this.File,
			Start: this.Pos,
			End:   this.Pos + 1,
		},
		Name:    "illegal",
		Content: fmt.Sprintf("illegal character `%c`", this.File.Content[this.Pos]),
	}
	this.Valid = false
	return false
}

func (this *Tokenizer) Reset() {
	this.Pos = 0
	this.current = []state{tokens[0].state}
}

func (this *Tokenizer) Save() TokenizerState {
	return TokenizerState{
		pos:   this.Pos,
		state: this.current,
		token: this.Token,
	}
}

func (this *Tokenizer) Load(save TokenizerState) {
	this.Pos = save.pos
	this.current = save.state
	this.Token = save.token
}
