package parser

import (
	"fmt"
	"regexp"
	"strings"
)

type state string

type tokenDefine struct {
	state       state
	name        string
	stateChange stateFunc
	skip        bool
	expr        *regexp.Regexp
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

func parseTokens(content string) ([]tokenDefine, state) {
	firstState := state("")
	tokens := []tokenDefine{}

	for _, line := range strings.Split(content, "\n") {
		tok := tokenDefine{}
		elems := strings.SplitN(line, " ", 3)
		if len(elems) != 3 {
			continue
		}

		tok.state = state(strings.TrimSpace(elems[0]))
		tok.name = strings.TrimSpace(elems[1])
		expr := strings.TrimSpace(elems[2])

		if len(firstState) == 0 {
			firstState = tok.state
		}

		if subs := strings.Split(tok.name, "->"); len(subs) != 1 {
			tok.name = subs[0]
			tok.stateChange = statePush(state(subs[1]))
		} else if strings.HasSuffix(tok.name, "<-") {
			tok.stateChange = statePop()
			tok.name = tok.name[:len(tok.name)-2]
		}

		if tok.name[0] == '.' {
			tok.skip = true
			tok.name = tok.name[1:]
		}

		tok.expr = regexp.MustCompile(fmt.Sprintf("^(%s)", expr))
		tokens = append(tokens, tok)
	}
	return tokens, firstState
}

func NewTokenizer(text string) *Tokenizer {
	return &Tokenizer{
		Text:    text,
		current: []state{tokens[0].state},
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
	this.current = []state{tokens[0].state}
}

func (this *Tokenizer) Save() {
	this.save = this.Pos
	this.savestate = this.current
}

func (this *Tokenizer) Load() {
	this.Pos = this.save
	this.current = this.savestate
}
