all:
	nex mydsl.nex
	# Could use nex instead of ed, but that'd be a little gratuitous.
	printf '/NEX_END_OF_LEXER_STRUCT/i\np *MyDsl\n.\nw\nq\n' | ed -s mydsl.nn.go
	go tool yacc -o=mydsl.yacc.go mydsl.y
	go fmt 
	go build
test:
	go test
clean:
	-rm *.output *.yacc.go *.nn.go