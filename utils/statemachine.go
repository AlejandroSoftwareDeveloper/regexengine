package utils

type State struct {
	start       bool
	terminal    bool
	transitions map[uint8][]*State
}

const epsilonChar uint8 = 0

func ToNfa(ctx *ParseContext) *State {
	startState, endState := TokenToNfa(&ctx.tokens[0])
	for i := 1; i < len(ctx.tokens); i++ {
		startNext, endNext := TokenToNfa(&ctx.tokens[i])
		endState.transitions[epsilonChar] = append(
			endState.transitions[epsilonChar],
			startNext,
		)
		endState = endNext
	}

	start := &State{
		transitions: map[uint8][]*State{
			epsilonChar: {startState},
		},
		start: true,
	}
	end := &State{
		transitions: map[uint8][]*State{},
		terminal:    true,
	}

	endState.transitions[epsilonChar] = append(
		endState.transitions[epsilonChar],
		end,
	)

	return start
}

func TokenToNfa(t *Token) (*State, *State) {
	start := &State{
		transitions: map[uint8][]*State{},
	}
	end := &State{
		transitions: map[uint8][]*State{},
	}

	switch t.TokenType {
	case literal:
		ch := t.value.(uint8)
		start.transitions[ch] = []*State{end}
	case or:

		values := t.value.([]Token)
		left := values[0]
		right := values[1]

		s1, e1 := TokenToNfa(&left)
		s2, e2 := TokenToNfa(&right)

		start.transitions[epsilonChar] = []*State{s1, s2}
		e1.transitions[epsilonChar] = []*State{end}
		e2.transitions[epsilonChar] = []*State{end}

	case bracket:

		literals := t.value.(map[uint8]bool)

		for l := range literals {
			start.transitions[l] = []*State{end}
		}
	case group, groupUncaptured:
		tokens := t.value.([]Token)
		start, end = TokenToNfa(&tokens[0])
		for i := 1; i < len(tokens); i++ {
			ts, te := TokenToNfa(&tokens[i])
			end.transitions[epsilonChar] = append(
				end.transitions[epsilonChar],
				ts,
			)
			end = te
		}
	case repeat:
		p := t.value.(RepeatPayload)

		if p.min == 0 {
			start.transitions[epsilonChar] = []*State{end}
		}

		var copyCount int

		if p.max == repeatInfinity {
			if p.min == 0 {
				copyCount = 1
			} else {
				copyCount = p.min
			}
		} else {
			copyCount = p.max
		}

		from, to := TokenToNfa(&p.token)
		start.transitions[epsilonChar] = append(
			start.transitions[epsilonChar],
			from,
		)

		for i := 2; i <= copyCount; i++ {
			s, e := TokenToNfa(&p.token)

			to.transitions[epsilonChar] = append(
				to.transitions[epsilonChar],
				s,
			)

			from = s
			to = e

			if i > p.min {
				s.transitions[epsilonChar] = append(
					s.transitions[epsilonChar],
					end,
				)
			}
		}

		to.transitions[epsilonChar] = append(
			to.transitions[epsilonChar],
			end,
		)

		if p.max == repeatInfinity {
			end.transitions[epsilonChar] = append(
				end.transitions[epsilonChar],
				from,
			)
		}
	default:
		panic("Tipo de token desconocido.")
	}

	return start, end
}
