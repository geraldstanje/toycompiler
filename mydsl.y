%{
package dsl
%}

%union {
  s string
  expr *Expr
}

%token WHILE
%token PRINT
%token END_LINE
%token ADD_OP
%token MUL_OP
%token ASSIGN
%token BEGIN_EXPRESSION
%token END_EXPRESSION
%token BEGIN_BLOCK
%token END_BLOCK
%token NUMBER
%token IDENTIFIER
%%
program: statement
       | statement END_LINE         
       | statement END_LINE program { var err error; if $$.expr, err = NewProgramNode(1,$1); err != nil { panic(err); } }

statement: assignation

assignation: IDENTIFIER ASSIGN NUMBER { /*var err error; if $$.expr, err = NewAssignExpr(1,$1); err != nil { panic(err); }*/ }

%%
func cast(y yyLexer) *MyDsl { return y.(*Lexer).p }
