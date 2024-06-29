package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type TokenType uint8

const repeatInfinity = -1
const (
	group           TokenType = iota
	bracket         TokenType = iota
	or              TokenType = iota
	repeat          TokenType = iota
	literal         TokenType = iota
	groupUncaptured TokenType = iota
)

type Token struct {
	TokenType TokenType
	value     interface{}
}

type ParseContext struct {
	pos    int
	tokens []Token
}

type RepeatPayload struct {
	min   int
	max   int
	token Token
}

func Parse(regex string) *ParseContext {
	ctx := &ParseContext{
		pos:    0,
		tokens: []Token{},
	}
	for ctx.pos < len(regex) {
		Process(regex, ctx)
		ctx.pos++
	}

	return ctx
}

func Process(regex string, ctx *ParseContext) {
	ch := regex[ctx.pos]
	if ch == '(' {
		groupCtx := &ParseContext{
			pos:    ctx.pos,
			tokens: []Token{},
		}
		ParseGroup(regex, groupCtx)
		ctx.tokens = append(ctx.tokens, Token{
			TokenType: group,
			value:     groupCtx.tokens,
		})
	} else if ch == '[' {
		ParseBracket(regex, ctx)
	} else if ch == '|' {
		ParseOr(regex, ctx)
	} else if ch == '*' || ch == '?' || ch == '+' {
		ParseRepeat(regex, ctx)
	} else if ch == '{' {
		ParseRepeatSpecified(regex, ctx)
	} else {
		t := Token{
			TokenType: literal,
			value:     ch,
		}
		ctx.tokens = append(ctx.tokens, t)
	}
}

func ParseGroup(regex string, ctx *ParseContext) {
	ctx.pos += 1
	for regex[ctx.pos] != ')' {
		Process(regex, ctx)
		ctx.pos += 1
	}
}

func ParseBracket(regex string, ctx *ParseContext) {
	ctx.pos++
	var literals []string
	for regex[ctx.pos] != ']' {
		ch := regex[ctx.pos]

		if ch == '-' {
			next := regex[ctx.pos+1]
			prev := literals[len(literals)-1][0]
			literals[len(literals)-1] = fmt.Sprintf("%c%c", prev, next)
			ctx.pos++
		} else {
			literals = append(literals, fmt.Sprintf("%c", ch))
		}
		ctx.pos++
	}

	literalsSet := map[uint8]bool{}

	for _, l := range literals {
		for i := l[0]; i <= l[len(l)-1]; i++ {
			literalsSet[i] = true
		}
	}
	ctx.tokens = append(ctx.tokens, Token{
		TokenType: bracket,
		value:     literalsSet,
	})
}

func ParseOr(regex string, ctx *ParseContext) {

	rhsContext := &ParseContext{
		pos:    ctx.pos,
		tokens: []Token{},
	}
	rhsContext.pos += 1
	for rhsContext.pos < len(regex) && regex[rhsContext.pos] != ')' {
		Process(regex, rhsContext)
		rhsContext.pos += 1
	}
	left := Token{
		TokenType: groupUncaptured,
		value:     ctx.tokens,
	}
	right := Token{
		TokenType: groupUncaptured,
		value:     rhsContext.tokens,
	}
	ctx.pos = rhsContext.pos
	ctx.tokens = []Token{{
		TokenType: or,
		value:     []Token{left, right},
	}}
}

func ParseRepeat(regex string, ctx *ParseContext) {
	ch := regex[ctx.pos]
	var min, max int
	if ch == '*' {
		min = 0
		max = repeatInfinity
	} else if ch == '?' {
		min = 0
		max = 1
	} else {
		min = 1
		max = repeatInfinity
	}
	lastToken := ctx.tokens[len(ctx.tokens)-1]
	ctx.tokens[len(ctx.tokens)-1] = Token{
		TokenType: repeat,
		value: RepeatPayload{
			min:   min,
			max:   max,
			token: lastToken,
		},
	}
}

func ParseRepeatSpecified(regex string, ctx *ParseContext) {
	start := ctx.pos + 1
	for regex[ctx.pos] != '}' {
		ctx.pos++
	}
	boundariesStr := regex[start:ctx.pos]
	pieces := strings.Split(boundariesStr, ",")
	var min, max int
	if len(pieces) == 1 {
		if value, err := strconv.Atoi(pieces[0]); err != nil {
			panic(err.Error())
		} else {
			min = value
			max = value
		}
	} else if len(pieces) == 2 {
		if value, err := strconv.Atoi(pieces[0]); err != nil {
			panic(err.Error())
		} else {
			min = value
		}

		if pieces[1] == "" {
			max = repeatInfinity
		} else if value, err := strconv.Atoi(pieces[1]); err != nil {
			panic(err.Error())
		} else {
			max = value
		}
	} else {
		panic(fmt.Sprintf("Debe existir uno o dos valores especificados por el cuantificador: dados '%s'", boundariesStr))
	}

	lastToken := ctx.tokens[len(ctx.tokens)-1]
	ctx.tokens[len(ctx.tokens)-1] = Token{
		TokenType: repeat,
		value: RepeatPayload{
			min:   min,
			max:   max,
			token: lastToken,
		},
	}
}
