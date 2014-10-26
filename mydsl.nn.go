package dsl

import (
	"bufio"
	"io"
	"strings"
)

type intstring struct {
	i int
	s string
}
type Lexer struct {
	// The lexer runs in its own goroutine, and communicates via channel 'ch'.
	ch chan intstring
	// We record the level of nesting because the action could return, and a
	// subsequent call expects to pick up where it left off. In other words,
	// we're simulating a coroutine.
	// TODO: Support a channel-based variant that compatible with Go's yacc.
	stack []intstring
	stale bool

	// The 'l' and 'c' fields were added for
	// https://github.com/wagerlabs/docker/blob/65694e801a7b80930961d70c69cba9f2465459be/buildfile.nex
	l, c int // line number and character position
	// The following line makes it easy for scripts to insert fields in the
	// generated code.
	p *MyDsl
	// [NEX_END_OF_LEXER_STRUCT]
}

// NewLexerWithInit creates a new Lexer object, runs the given callback on it,
// then returns it.
func NewLexerWithInit(in io.Reader, initFun func(*Lexer)) *Lexer {
	type dfa struct {
		acc          []bool           // Accepting states.
		f            []func(rune) int // Transitions.
		startf, endf []int            // Transitions at start and end of input.
		nest         []dfa
	}
	yylex := new(Lexer)
	if initFun != nil {
		initFun(yylex)
	}
	yylex.ch = make(chan intstring)
	var scan func(in *bufio.Reader, ch chan intstring, family []dfa)
	scan = func(in *bufio.Reader, ch chan intstring, family []dfa) {
		// Index of DFA and length of highest-precedence match so far.
		matchi, matchn := 0, -1
		var buf []rune
		n := 0
		checkAccept := func(i int, st int) bool {
			// Higher precedence match? DFAs are run in parallel, so matchn is at most len(buf), hence we may omit the length equality check.
			if family[i].acc[st] && (matchn < n || matchi > i) {
				matchi, matchn = i, n
				return true
			}
			return false
		}
		var state [][2]int
		for i := 0; i < len(family); i++ {
			mark := make([]bool, len(family[i].startf))
			// Every DFA starts at state 0.
			st := 0
			for {
				state = append(state, [2]int{i, st})
				mark[st] = true
				// As we're at the start of input, follow all ^ transitions and append to our list of start states.
				st = family[i].startf[st]
				if -1 == st || mark[st] {
					break
				}
				// We only check for a match after at least one transition.
				checkAccept(i, st)
			}
		}
		atEOF := false
		for {
			if n == len(buf) && !atEOF {
				r, _, err := in.ReadRune()
				switch err {
				case io.EOF:
					atEOF = true
				case nil:
					buf = append(buf, r)
				default:
					panic(err)
				}
			}
			if !atEOF {
				r := buf[n]
				n++
				var nextState [][2]int
				for _, x := range state {
					x[1] = family[x[0]].f[x[1]](r)
					if -1 == x[1] {
						continue
					}
					nextState = append(nextState, x)
					checkAccept(x[0], x[1])
				}
				state = nextState
			} else {
			dollar: // Handle $.
				for _, x := range state {
					mark := make([]bool, len(family[x[0]].endf))
					for {
						mark[x[1]] = true
						x[1] = family[x[0]].endf[x[1]]
						if -1 == x[1] || mark[x[1]] {
							break
						}
						if checkAccept(x[0], x[1]) {
							// Unlike before, we can break off the search. Now that we're at the end, there's no need to maintain the state of each DFA.
							break dollar
						}
					}
				}
				state = nil
			}

			if state == nil {
				// All DFAs stuck. Return last match if it exists, otherwise advance by one rune and restart all DFAs.
				if matchn == -1 {
					if len(buf) == 0 { // This can only happen at the end of input.
						break
					}
					buf = buf[1:]
				} else {
					text := string(buf[:matchn])
					buf = buf[matchn:]
					matchn = -1
					ch <- intstring{matchi, text}
					if len(family[matchi].nest) > 0 {
						scan(bufio.NewReader(strings.NewReader(text)), ch, family[matchi].nest)
					}
					if atEOF {
						break
					}
				}
				n = 0
				for i := 0; i < len(family); i++ {
					state = append(state, [2]int{i, 0})
				}
			}
		}
		ch <- intstring{-1, ""}
	}
	go scan(bufio.NewReader(in), yylex.ch, []dfa{
		// while
		{[]bool{false, false, false, false, false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 104:
					return -1
				case 105:
					return -1
				case 108:
					return -1
				case 101:
					return -1
				case 119:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 119:
					return -1
				case 104:
					return 2
				case 105:
					return -1
				case 108:
					return -1
				case 101:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 119:
					return -1
				case 104:
					return -1
				case 105:
					return 3
				case 108:
					return -1
				case 101:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 105:
					return -1
				case 108:
					return 4
				case 101:
					return -1
				case 119:
					return -1
				case 104:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 108:
					return -1
				case 101:
					return 5
				case 119:
					return -1
				case 104:
					return -1
				case 105:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 119:
					return -1
				case 104:
					return -1
				case 105:
					return -1
				case 108:
					return -1
				case 101:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1}, nil},

		// print
		{[]bool{false, false, false, false, false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 112:
					return 1
				case 114:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 116:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 105:
					return -1
				case 110:
					return -1
				case 116:
					return -1
				case 112:
					return -1
				case 114:
					return 2
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 110:
					return -1
				case 116:
					return -1
				case 112:
					return -1
				case 114:
					return -1
				case 105:
					return 3
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 114:
					return -1
				case 105:
					return -1
				case 110:
					return 4
				case 116:
					return -1
				case 112:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 105:
					return -1
				case 110:
					return -1
				case 116:
					return 5
				case 112:
					return -1
				case 114:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 116:
					return -1
				case 112:
					return -1
				case 114:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1}, nil},

		// ;
		{[]bool{false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 59:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 59:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

		// \+|-
		{[]bool{false, true, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 45:
					return 1
				case 43:
					return 2
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 43:
					return -1
				case 45:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 43:
					return -1
				case 45:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

		// \*|\/
		{[]bool{false, true, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 42:
					return 1
				case 47:
					return 2
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 42:
					return -1
				case 47:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 47:
					return -1
				case 42:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

		// =
		{[]bool{false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 61:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 61:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

		// \(
		{[]bool{false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 40:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 40:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

		// \)
		{[]bool{false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 41:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 41:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

		// \{
		{[]bool{false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 123:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 123:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

		// \}
		{[]bool{false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 125:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 125:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

		// [0-9]+
		{[]bool{false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch {
				case 48 <= r && r <= 57:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch {
				case 48 <= r && r <= 57:
					return 1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

		// [a-z][a-z0-9]*
		{[]bool{false, true, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch {
				case 48 <= r && r <= 57:
					return -1
				case 97 <= r && r <= 122:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch {
				case 48 <= r && r <= 57:
					return 2
				case 97 <= r && r <= 122:
					return 2
				}
				return -1
			},
			func(r rune) int {
				switch {
				case 48 <= r && r <= 57:
					return 2
				case 97 <= r && r <= 122:
					return 2
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},
	})
	return yylex
}
func NewLexer(in io.Reader) *Lexer {
	return NewLexerWithInit(in, nil)
}
func (yylex *Lexer) Text() string {
	return yylex.stack[len(yylex.stack)-1].s
}
func (yylex *Lexer) next(lvl int) int {
	if lvl == len(yylex.stack) {
		yylex.stack = append(yylex.stack, intstring{0, ""})
	}
	if lvl == len(yylex.stack)-1 {
		p := &yylex.stack[lvl]
		*p = <-yylex.ch
		yylex.stale = false
	} else {
		yylex.stale = true
	}
	return yylex.stack[lvl].i
}
func (yylex *Lexer) pop() {
	yylex.stack = yylex.stack[:len(yylex.stack)-1]
}
func (yylex Lexer) Error(e string) {
	panic(e)
}
func (yylex *Lexer) Lex(lval *yySymType) int {
	for {
		switch yylex.next(0) {
		case 0:
			{
				return WHILE
			}
			continue
		case 1:
			{
				return PRINT
			}
			continue
		case 2:
			{
				return END_LINE
			}
			continue
		case 3:
			{
				lval.s = yylex.Text()
				return ADD_OP
			}
			continue
		case 4:
			{
				lval.s = yylex.Text()
				return MUL_OP
			}
			continue
		case 5:
			{
				return ASSIGN
			}
			continue
		case 6:
			{
				return BEGIN_EXPRESSION
			}
			continue
		case 7:
			{
				return END_EXPRESSION
			}
			continue
		case 8:
			{
				return BEGIN_BLOCK
			}
			continue
		case 9:
			{
				return END_BLOCK
			}
			continue
		case 10:
			{
				lval.s = yylex.Text()
				return NUMBER
			}
			continue
		case 11:
			{
				lval.s = yylex.Text()
				return IDENTIFIER
			}
			continue
		}
		break
	}
	yylex.pop()

	return 0
}
