all:
	nex -e=true lexer.nex
	# Could use nex instead of ed, but that'd be a little gratuitous.
	printf '/NEX_END_OF_LEXER_STRUCT/i\np *Compiler\n.\nw\nq\n' | ed -s lexer.nn.go
	go tool yacc -o=parser.yacc.go parser.y
	go fmt 
	go build
test:
	go test
clean:
	-rm *.output *.yacc.go *.nn.go *.pdf *.dot generated.txt