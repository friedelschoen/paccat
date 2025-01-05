package parser

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"friedelschoen.io/paccat/internal/ast"
	"friedelschoen.io/paccat/internal/errors"
)

type parseState struct {
	filename string
	lex      *Tokenizer
}

type parseError struct {
	got    Token
	expect []string /* expected ... */
}

func (this *parseState) choice(choices ...func() (ast.Node, *parseError)) (ast.Node, *parseError) {
	expect := []string{}
	got := this.lex.Token

	for _, ch := range choices {
		this.lex.Save()
		res, err := ch()
		if err == nil {
			return res, nil
		}
		if err.got.Start > got.Start {
			expect = err.expect
			got = err.got
		} else if err.got == got {
			expect = append(expect, err.expect...)
		}
		this.lex.Load()
	}
	return nil, &parseError{got, expect}
}

func (this *parseState) newPos(from, to Token) errors.Position {
	return errors.Position{
		Filename: this.filename,
		Content:  &this.lex.Text, Start: from.Start, End: to.End,
	}
}

func (this *parseState) expectToken(name string) (Token, *parseError) {
	current := this.lex.Token
	if current.Name == name {
		this.lex.Next()
		return current, nil
	}
	return Token{}, &parseError{this.lex.Token, []string{name}}
}

func (this *parseState) expectTokenContent(content string) (Token, *parseError) {
	current := this.lex.Token
	// fmt.Printf("`%s` == `%s`: %v\n", current.Content, content, current.Content == content)
	if current.Content == content {
		this.lex.Next()
		return current, nil
	}
	return Token{}, &parseError{this.lex.Token, []string{"`" + content + "`"}}
}

func (this *parseState) asLiteral(tok Token) *ast.LiteralNode {
	return &ast.LiteralNode{
		Pos:     this.newPos(tok, tok),
		Content: tok.Content,
	}
}

func (this *parseState) parseLambda() (ast.Node, *parseError) {
	begin, err := this.expectTokenContent("(")
	if err != nil {
		return nil, err
	}
	args := ast.LiteralMap{}
tokenLoop:
	for this.lex.Valid {
		if len(args) > 0 {
			_, err := this.expectTokenContent(",")
			if err != nil {
				break tokenLoop
			}
		}
		var def ast.Node
		ident, err := this.expectToken("ident")
		if err != nil {
			break tokenLoop
		}
		_, err = this.expectTokenContent("=")
		if err == nil {
			def, err = this.parseValue()
			if err != nil {
				return nil, err
			}
		}
		args[ident.Content] = ast.LiteralMapPair{
			Key:   this.asLiteral(ident),
			Value: def,
		}
	}

	end, err := this.expectTokenContent(")")
	if err != nil {
		return nil, err
	}

	_, err = this.expectTokenContent("->")
	if err != nil {
		return nil, err
	}

	target, err := this.parseValue()
	if err != nil {
		return nil, err
	}

	return &ast.LambdaNode{
		Pos:    this.newPos(begin, end),
		Target: target,
		Args:   args,
	}, nil
}

func (this *parseState) parseDict() (ast.Node, *parseError) {
	begin, err := this.expectTokenContent("{")
	if err != nil {
		return nil, err
	}
	items := ast.LiteralMap{}
tokenLoop:
	for this.lex.Valid {
		if len(items) > 0 {
			_, err := this.expectTokenContent(",")
			if err != nil {
				break tokenLoop
			}
		}
		ident, err := this.expectToken("ident")
		if err != nil {
			break tokenLoop
		}
		_, err = this.expectTokenContent("=")
		if err != nil {
			return nil, err
		}
		value, err := this.parseValue()
		if err != nil {
			return nil, err
		}
		items[ident.Content] = ast.LiteralMapPair{
			Key:   this.asLiteral(ident),
			Value: value,
		}
	}

	end, err := this.expectTokenContent("}")
	if err != nil {
		return nil, err
	}

	return &ast.DictNode{
		Pos:   this.newPos(begin, end),
		Items: items,
	}, nil
}

func (this *parseState) parseList() (ast.Node, *parseError) {
	begin, err := this.expectTokenContent("[")
	if err != nil {
		return nil, err
	}
	items := []ast.Node{}
tokenLoop:
	for this.lex.Valid {
		if len(items) > 0 {
			_, err := this.expectTokenContent(",")
			if err != nil {
				break tokenLoop
			}
		}
		value, err := this.parseValue()
		if err != nil {
			break tokenLoop
		}
		items = append(items, value)
	}

	end, err := this.expectTokenContent("]")
	if err != nil {
		return nil, err
	}

	return &ast.ListNode{
		Pos:   this.newPos(begin, end),
		Items: items,
	}, nil
}

func (this *parseState) parseSurrounded() (ast.Node, *parseError) {
	_, err := this.expectTokenContent("(")
	if err != nil {
		return nil, err
	}
	value, err := this.parseValue()
	if err != nil {
		return nil, err
	}
	_, err = this.expectTokenContent(")")
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (this *parseState) parseReference() (ast.Node, *parseError) {
	begin, err := this.expectTokenContent("$")
	if err != nil {
		return nil, err
	}
	ident, err := this.expectToken("ident")
	if err != nil {
		return nil, err
	}
	return &ast.ReferenceNode{
		Pos:      this.newPos(begin, ident),
		Variable: this.asLiteral(ident),
	}, nil
}

func (this *parseState) parseOutput() (ast.Node, *parseError) {
	begin, err := this.expectTokenContent("output")
	if err != nil {
		return nil, err
	}
	options, err := this.parseDict()
	if err != nil {
		return nil, err
	}
	return &ast.OutputNode{
		Pos: errors.Position{
			Filename: this.filename,
			Content:  &this.lex.Text,
			Start:    begin.Start,
			End:      options.GetPosition().End,
		},
		Options: options.(*ast.DictNode).Items,
	}, nil
}

func (this *parseState) parseImport() (ast.Node, *parseError) {
	begin, err := this.expectTokenContent("import")
	if err != nil {
		return nil, err
	}
	source, err := this.parseValue()
	if err != nil {
		return nil, err
	}
	return &ast.ImportNode{
		Pos: errors.Position{
			Filename: this.filename,
			Content:  &this.lex.Text,
			Start:    begin.Start,
			End:      source.GetPosition().End,
		},
		Source: source,
	}, nil
}

func (this *parseState) parsePanic() (ast.Node, *parseError) {
	begin, err := this.expectTokenContent("panic")
	if err != nil {
		return nil, err
	}
	message, err := this.parseValue()
	if err != nil {
		return nil, err
	}
	return &ast.PanicNode{
		Pos: errors.Position{
			Filename: this.filename,
			Content:  &this.lex.Text,
			Start:    begin.Start,
			End:      message.GetPosition().End,
		},
		Message: message,
	}, nil
}

func (this *parseState) parseString(wrap string) func() (ast.Node, *parseError) {
	return func() (ast.Node, *parseError) {
		begin, err := this.expectTokenContent(wrap)
		if err != nil {
			return nil, err
		}

		builder := strings.Builder{}
		result := make([]ast.Node, 0)
		currentPos := 0

	tokenLoop:
		for this.lex.Valid {
			switch this.lex.Token.Name {
			case "char":
				if builder.Len() == 0 {
					currentPos = this.lex.Token.Start
				}
				builder.WriteString(this.lex.Token.Content)
				this.lex.Next()

			case "interp-begin":
				this.lex.Next()
				value, exp := this.parseValue()
				if exp != nil {
					return nil, exp
				}
				if builder.Len() > 0 {
					node := &ast.LiteralNode{
						Pos: errors.Position{
							Filename: this.filename,
							Content:  &this.lex.Text,
							Start:    currentPos,
							End:      currentPos + builder.Len(),
						},
						Content: builder.String(),
					}
					result = append(result, node, value)
					builder.Reset()
				} else {
					result = append(result, value)
				}
				_, err = this.expectToken("interp-end")
				if err != nil {
					return nil, err
				}

			default:
				break tokenLoop
			}
		}
		end, err := this.expectTokenContent(wrap)
		if err != nil {
			return nil, err
		}

		if builder.Len() > 0 {
			node := &ast.LiteralNode{
				Pos: errors.Position{
					Filename: this.filename,
					Content:  &this.lex.Text,
					Start:    currentPos,
					End:      currentPos + builder.Len(),
				},
				Content: builder.String(),
			}
			result = append(result, node)
		}

		return &ast.StringNode{
			Pos:     this.newPos(begin, end),
			Content: result,
		}, nil
	}
}

func (this *parseState) parsePath() (ast.Node, *parseError) {
	val, err := this.expectToken("path")
	if err != nil {
		return nil, err
	}
	return &ast.LiteralNode{
		Pos:     this.newPos(val, val),
		Content: val.Content,
	}, nil
}

func (this *parseState) parseValue() (ast.Node, *parseError) {
	val, err := this.choice(
		this.parseString("\""),
		this.parseString("''"),
		this.parsePath,
		this.parseLambda,
		this.parseReference,
		this.parseList,
		this.parseDict,
		this.parseSurrounded,
		this.parseOutput,
		this.parseImport,
		this.parsePanic,
	)
	if err != nil {
		return nil, err
	}
tokenLoop:
	for this.lex.Valid {
		begin := this.lex.Token
		switch begin.Content {
		case ".":
			this.lex.Next()
			ident, err := this.expectToken("ident")
			if err != nil {
				return nil, err
			}
			val = &ast.GetterNode{
				Pos:       this.newPos(begin, ident),
				Target:    val,
				Attribute: this.asLiteral(ident),
			}
		case "(":
			this.lex.Next()
			args := ast.LiteralMap{}
		argLoop:
			for this.lex.Valid {
				if len(args) > 0 {
					_, err := this.expectTokenContent(",")
					if err != nil {
						break argLoop
					}
				}
				ident, err := this.expectToken("ident")
				if err != nil {
					return nil, err
				}
				_, err = this.expectTokenContent("=")
				if err != nil {
					return nil, err
				}
				value, err := this.parseValue()
				if err != nil {
					return nil, err
				}
				args[ident.Content] = ast.LiteralMapPair{
					Key:   this.asLiteral(ident),
					Value: value,
				}
			}
			end, err := this.expectTokenContent(")")
			if err != nil {
				return nil, err
			}
			val = &ast.CallNode{
				Pos:    this.newPos(begin, end),
				Target: val,
				Args:   args,
			}
		default:
			break tokenLoop
		}
	}
	return val, nil
}

func (this *parseState) parseFile() (ast.Node, *parseError) {
	val, err := this.parseValue()
	if err != nil {
		return nil, err
	}
	_, err = this.expectToken("eof")
	if err != nil {
		return nil, err
	}
	return val, nil
}

func unique[T comparable](slc []T) []T {
	for i := 0; i < len(slc); i++ {
		for j := i + 1; j < len(slc); j++ {
			if slc[i] == slc[j] {
				l := len(slc)
				slc[i] = slc[l-1] /* move last element to current */
				slc = slc[:l-1]   /* shrink slice by one */
				i--               /* decrement i, we want to re-check this */
				break
			}
		}
	}
	return slc
}

func Parse(filename, content string) (ast.Node, error) {
	lex := NewTokenizer(content)
	parser := parseState{
		filename: filename,
		lex:      lex,
	}

	lex.Next()
	result, exp := parser.parseFile()
	if exp != nil {
		exp.expect = unique(exp.expect)
		slices.Sort(exp.expect)
		pos := parser.newPos(exp.got, exp.got)
		message := strings.Builder{}
		message.WriteString("expected token ")
		for _, token := range exp.expect {
			message.WriteString(token)
			message.WriteString(", ")
		}
		if exp.got.Name == "illegal" {
			message.WriteString("but got ")
			message.WriteString(exp.got.Content)
		} else {
			fmt.Fprintf(&message, "but got `%s`", exp.got.Content)
		}
		return nil, errors.NewRecipeError(pos, message.String())
	}
	return result, nil
}

func ParseFile(filename string) (ast.Node, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return Parse(filename, string(content))
}
