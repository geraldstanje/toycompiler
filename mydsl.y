%{
package dsl
%}

%union {
  s string
  //expr *Expr
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
       | statement END_LINE         { expr, err := NewProgramNode(1,$1)
                                      if err != nil { panic(err); } 
                                      cast(yylex).AppendExpr(expr)
                                    }
       | statement END_LINE program { expr, err := NewProgramNode(1,$1)
                                      if err != nil { panic(err); } 
                                      cast(yylex).AppendExpr(expr)
                                    }

statement: assignation

assignation: IDENTIFIER ASSIGN expression { expr, err := NewAssignNode(1,$1)
                                            if err != nil { panic(err); } 
                                            cast(yylex).AppendExpr(expr)
                                          }

expression: NUMBER { expr, err := NewTokenNode(1,$1)
                     if err != nil { panic(err); } 
                     cast(yylex).AppendExpr(expr)
                   }

%%
func cast(y yyLexer) *MyDsl { return y.(*Lexer).p }
