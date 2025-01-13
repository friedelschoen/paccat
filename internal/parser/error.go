package parser

import (
	"fmt"
	"slices"
	"strings"

	"friedelschoen.io/paccat/internal/errors"
)

type parseError struct {
	got    Token
	expect []string /* expected ... */
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

func (this *parseError) Error() string {
	this.expect = unique(this.expect)
	slices.Sort(this.expect)

	message := &strings.Builder{}
	message.WriteString("expected token ")
	for i, token := range this.expect {
		switch {
		case i == len(this.expect)-1:
			message.WriteString(" or ")
		case i >= 1:
			message.WriteString(", ")
		}
		message.WriteString(token)
	}
	message.WriteString(" but got ")
	if this.got.Name == "illegal" {
		message.WriteString(this.got.Content)
	} else {
		fmt.Fprintf(message, "`%s` (%s)", this.got.Content, this.got.Name)
	}
	return message.String()
}

func (this *parseError) GetPosition() errors.Position {
	return this.got.Pos
}
