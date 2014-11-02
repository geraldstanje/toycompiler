# toycompiler

## Goal:

This project ports the following compiler from Python to Go:
(http://www.dimitrifourny.com/2014/04/17/write-your-first-compiler/)

The compiler uses:

- Nex lexer: https://crypto.stanford.edu/~blynn/nex/
- Yacc parser generator: https://golang.org/cmd/yacc/

## Installation:

To try the release you'll need:

* OSX 10.7 (Lion) or greater, x86 and x86-64
* Go >= 1.3
* Graphviz: brew install graphviz
* Nex: https://github.com/blynn/nex/

After the requirements are satisfied:

  go get github.com/geraldstanje/toycompiler