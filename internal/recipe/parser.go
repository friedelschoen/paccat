package recipe

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"friedelschoen.io/paccat/internal/lexer"
)

type parseState struct {
	filename string
	lex      *lexer.Tokenizer
}

type parseError struct {
	got    lexer.Token
	expect []string /* expected ... */
}

func (this *parseState) choice(choices ...func() (Evaluable, *parseError)) (Evaluable, *parseError) {
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

func (this *parseState) newPos(from, to lexer.Token) Position {
	return Position{this.filename, &this.lex.Text, from.Start, to.End}
}

func (this *parseState) expectToken(name string) (lexer.Token, *parseError) {
	current := this.lex.Token
	if current.Name == name {
		this.lex.Next()
		return current, nil
	}
	return lexer.Token{}, &parseError{this.lex.Token, []string{name}}
}

func (this *parseState) expectTokenContent(content string) (lexer.Token, *parseError) {
	current := this.lex.Token
	// fmt.Printf("`%s` == `%s`: %v\n", current.Content, content, current.Content == content)
	if current.Content == content {
		this.lex.Next()
		return current, nil
	}
	return lexer.Token{}, &parseError{this.lex.Token, []string{"`" + content + "`"}}
}

func (this *parseState) parseLambda() (Evaluable, *parseError) {
	begin, err := this.expectTokenContent("(")
	if err != nil {
		return nil, err
	}
	args := map[string]Evaluable{}
tokenLoop:
	for this.lex.Valid {
		if len(args) > 0 {
			_, err := this.expectTokenContent(",")
			if err != nil {
				break tokenLoop
			}
		}
		var def Evaluable
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

	return &recipeLambda{this.newPos(begin, end), target, args}, nil
}

func (this *parseState) parseDict() (Evaluable, *parseError) {
	begin, err := this.expectTokenContent("{")
	if err != nil {
		return nil, err
	}
	items := map[string]Evaluable{}
tokenLoop:
	for this.lex.Valid {
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
		_, err = this.expectTokenContent(";")
		if err != nil {
			return nil, err
		}
		items[ident.Content] = value
	}

	end, err := this.expectTokenContent("}")
	if err != nil {
		return nil, err
	}

	return &recipeDict{this.newPos(begin, end), items}, nil
}

func (this *parseState) parseList() (Evaluable, *parseError) {
	begin, err := this.expectTokenContent("[")
	if err != nil {
		return nil, err
	}
	items := []Evaluable{}
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

	return &recipeList{this.newPos(begin, end), items}, nil
}

func (this *parseState) parseSurrounded() (Evaluable, *parseError) {
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

func (this *parseState) parseReference() (Evaluable, *parseError) {
	begin, err := this.expectTokenContent("$")
	if err != nil {
		return nil, err
	}
	ident, err := this.expectToken("ident")
	if err != nil {
		return nil, err
	}
	return &recipeReference{this.newPos(begin, ident), ident.Content}, nil
}

func (this *parseState) parseOutput() (Evaluable, *parseError) {
	begin, err := this.expectTokenContent("output")
	if err != nil {
		return nil, err
	}
	options, err := this.parseDict()
	if err != nil {
		return nil, err
	}
	pos := Position{this.filename, &this.lex.Text, begin.Start, options.GetPosition().End}
	return &recipeOutput{pos, options.(*recipeDict).items}, nil
}

func (this *parseState) parseImport() (Evaluable, *parseError) {
	begin, err := this.expectTokenContent("import")
	if err != nil {
		return nil, err
	}
	source, err := this.parseValue()
	if err != nil {
		return nil, err
	}
	pos := Position{this.filename, &this.lex.Text, begin.Start, source.GetPosition().End}
	return &recipeImport{pos, source}, nil
}

func (this *parseState) parsePanic() (Evaluable, *parseError) {
	begin, err := this.expectTokenContent("panic")
	if err != nil {
		return nil, err
	}
	message, err := this.parseValue()
	if err != nil {
		return nil, err
	}
	pos := Position{this.filename, &this.lex.Text, begin.Start, message.GetPosition().End}
	return &recipePanic{pos, message}, nil
}

func (this *parseState) parseString(wrap string) func() (Evaluable, *parseError) {
	return func() (Evaluable, *parseError) {
		begin, err := this.expectTokenContent(wrap)
		if err != nil {
			return nil, err
		}

		builder := strings.Builder{}
		result := make([]Evaluable, 0)
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
					pos := Position{this.filename, &this.lex.Text, currentPos, currentPos + builder.Len()}
					result = append(result, &recipeStringLiteral{pos, builder.String()}, value)
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
			pos := Position{this.filename, &this.lex.Text, currentPos, currentPos + builder.Len()}
			result = append(result, &recipeStringLiteral{pos, builder.String()})
		}

		return &recipeString{this.newPos(begin, end), result}, nil
	}
}

func (this *parseState) parsePath() (Evaluable, *parseError) {
	val, err := this.expectToken("path")
	if err != nil {
		return nil, err
	}
	return &recipeStringLiteral{this.newPos(val, val), val.Content}, nil
}

func (this *parseState) parseValue() (Evaluable, *parseError) {
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
			val = &recipeGetter{this.newPos(begin, ident), val, ident.Content}
		case "(":
			this.lex.Next()
			args := map[string]Evaluable{}
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
			val = &recipeCall{this.newPos(begin, end), val, args}
		default:
			break tokenLoop
		}
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

func Parse(filename, content string) (Evaluable, error) {
	lex := lexer.NewTokenizer(content)
	parser := parseState{
		filename: filename,
		lex:      lex,
	}

	lex.Next()
	result, exp := parser.parseValue()
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
		fmt.Fprintf(&message, "but got `%s`", exp.got.Content)
		return nil, NewRecipeError(pos, message.String())
	}
	return result, nil
}

func ParseFile(filename string) (Evaluable, error) {
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
