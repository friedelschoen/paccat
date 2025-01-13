package parser

import (
	"strings"

	"friedelschoen.io/paccat/internal/ast"
	"friedelschoen.io/paccat/internal/errors"
)

type parseState struct {
	Tokenizer
}

func stretch(from, to errors.Positioned) errors.Position {
	return from.GetPosition().Stretch(to.GetPosition())
}

func (this *parseState) choice(choices ...func() (ast.Node, *parseError)) (ast.Node, *parseError) {
	expect := []string{}
	got := this.Token

	for _, ch := range choices {
		lexsave := this.Save()
		res, err := ch()
		if err == nil {
			return res, nil
		}
		this.Load(lexsave)

		if err.got.Pos.Start > got.Pos.Start {
			expect = err.expect
			got = err.got
		} else if err.got == got {
			expect = append(expect, err.expect...)
		}
	}
	return nil, &parseError{got, expect}
}

func (this *parseState) expectToken(name string) (Token, *parseError) {
	current := this.Token
	if current.Name == name {
		this.Next()
		return current, nil
	}
	return Token{}, &parseError{this.Token, []string{name}}
}

func (this *parseState) expectTokenContent(content string) (Token, *parseError) {
	current := this.Token
	if current.Content == content {
		this.Next()
		return current, nil
	}
	return Token{}, &parseError{this.Token, []string{"`" + content + "`"}}
}

func (this *parseState) asLiteral(tok Token) *ast.LiteralNode {
	return &ast.LiteralNode{
		Pos:     tok.GetPosition(),
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
	for this.Valid {
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
		Pos:    stretch(begin, end),
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
	for this.Valid {
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
		Pos:   stretch(begin, end),
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
	for this.Valid {
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
		Pos:   stretch(begin, end),
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
	ident, err := this.expectToken("ident")
	if err != nil {
		return nil, err
	}
	return &ast.ReferenceNode{
		Pos:      ident.GetPosition(),
		Variable: this.asLiteral(ident),
	}, nil
}

func (this *parseState) parseNumber() (ast.Node, *parseError) {
	token, err := this.expectToken("number")
	if err != nil {
		return nil, err
	}
	return &ast.NumberNode{
		Pos:     token.GetPosition(),
		Content: this.asLiteral(token),
	}, nil
}

func (this *parseState) parseOutput() (ast.Node, *parseError) {
	begin, err := this.expectTokenContent("output")
	if err != nil {
		return nil, err
	}
	options, err := this.parseValue()
	if err != nil {
		return nil, err
	}
	return &ast.OutputNode{
		Pos:     stretch(begin, options),
		Options: options,
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
		Pos:    stretch(begin, source),
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
		Pos:     stretch(begin, message),
		Message: message,
	}, nil
}

func (this *parseState) parseString() (ast.Node, *parseError) {
	begin := this.Token
	if begin.Content != "\"" && begin.Content != "''" {
		return nil, &parseError{this.Token, []string{"`''`", "`\"`"}}
	}
	this.Next()

	builder := strings.Builder{}
	result := make([]ast.Node, 0)
	currentPos := 0

tokenLoop:
	for this.Valid {
		switch this.Token.Name {
		case "char":
			if builder.Len() == 0 {
				currentPos = this.Token.Pos.Start
			}
			builder.WriteString(this.Token.Content)
			this.Next()

		case "interp-begin":
			this.Next()
			value, exp := this.parseValue()
			if exp != nil {
				return nil, exp
			}
			if builder.Len() > 0 {
				node := &ast.LiteralNode{
					Pos: errors.Position{
						File:  this.File,
						Start: currentPos,
						End:   currentPos + builder.Len(),
					},
					Content: builder.String(),
				}
				result = append(result, node, value)
				builder.Reset()
			} else {
				result = append(result, value)
			}
			_, err := this.expectToken("interp-end")
			if err != nil {
				return nil, err
			}

		default:
			break tokenLoop
		}
	}
	end, err := this.expectTokenContent(begin.Content)
	if err != nil {
		return nil, err
	}

	if builder.Len() > 0 {
		node := &ast.LiteralNode{
			Pos: errors.Position{
				File:  this.File,
				Start: currentPos,
				End:   currentPos + builder.Len(),
			},
			Content: builder.String(),
		}
		result = append(result, node)
	}

	return &ast.StringNode{
		Pos:     stretch(begin, end),
		Content: result,
	}, nil
}

func (this *parseState) parsePath() (ast.Node, *parseError) {
	val, err := this.expectToken("path")
	if err != nil {
		return nil, err
	}
	return &ast.LiteralNode{
		Pos:     val.GetPosition(),
		Content: val.Content,
	}, nil
}

func (this *parseState) parseValue() (ast.Node, *parseError) {
	val, err := this.choice(
		this.parseString,
		this.parseNumber,
		this.parsePath,
		this.parseLambda,
		this.parseSurrounded,
		this.parseList,
		this.parseDict,
		this.parseOutput,
		this.parseImport,
		this.parsePanic,
		this.parseReference,
	)
	if err != nil {
		return nil, err
	}
tokenLoop:
	for this.Valid {
		begin := this.Token
		switch begin.Content {
		case ".":
			this.Next()
			ident, err := this.expectToken("ident")
			if err != nil {
				return nil, err
			}
			val = &ast.GetterNode{
				Pos:       stretch(begin, ident),
				Target:    val,
				Attribute: this.asLiteral(ident),
			}
		case "[":
			this.Next()
			attr, err := this.parseValue()
			if err != nil {
				return nil, err
			}
			end, err := this.expectTokenContent("]")
			if err != nil {
				return nil, err
			}
			val = &ast.GetterNode{
				Pos:       stretch(begin, end),
				Target:    val,
				Attribute: attr,
			}
		case "(":
			this.Next()
			args := ast.LiteralMap{}
		argLoop:
			for this.Valid {
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
				Pos:    stretch(begin, end),
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
