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

func (this *parseState) parseLambda() (ast.Node, *parseError) {
	begin, err := this.expectTokenContent("(")
	if err != nil {
		return nil, err
	}
	args := map[string]ast.Node{}
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
		args[ident.Content] = def
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
	items := map[string]ast.Node{}
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
		items[ident.Content] = value
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

	return &ast.ListNode{this.newPos(begin, end), items}, nil
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
	return &ast.ReferenceNode{this.newPos(begin, ident), ident.Content}, nil
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
	pos := errors.Position{this.filename, &this.lex.Text, begin.Start, options.GetPosition().End}
	return &ast.OutputNode{pos, options.(*ast.DictNode).Items}, nil
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
	pos := errors.Position{this.filename, &this.lex.Text, begin.Start, source.GetPosition().End}
	return &ast.ImportNode{pos, source}, nil
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
	pos := errors.Position{this.filename, &this.lex.Text, begin.Start, message.GetPosition().End}
	return &ast.PanicNode{pos, message}, nil
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
					pos := errors.Position{this.filename, &this.lex.Text, currentPos, currentPos + builder.Len()}
					result = append(result, &ast.LiteralNode{pos, builder.String()}, value)
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
			pos := errors.Position{this.filename, &this.lex.Text, currentPos, currentPos + builder.Len()}
			result = append(result, &ast.LiteralNode{pos, builder.String()})
		}

		return &ast.StringNode{this.newPos(begin, end), result}, nil
	}
}

func (this *parseState) parsePath() (ast.Node, *parseError) {
	val, err := this.expectToken("path")
	if err != nil {
		return nil, err
	}
	return &ast.LiteralNode{this.newPos(val, val), val.Content}, nil
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
			val = &ast.GetterNode{this.newPos(begin, ident), val, ident.Content}
		case "(":
			this.lex.Next()
			args := map[string]ast.Node{}
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
				args[ident.Content] = value
			}
			end, err := this.expectTokenContent(")")
			if err != nil {
				return nil, err
			}
			val = &ast.CallNode{this.newPos(begin, end), val, args}
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
