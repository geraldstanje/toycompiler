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
{
  $$ = $1
}
| statement END_LINE
{ 
  $$ = $1
}
| statement END_LINE program 
{ 
  programNode = newProgramNode($2)
  programNode.Left = $1
  programNode.Right = $3
}

statement: assignation {
  $$ = $1
}

assignation: IDENTIFIER ASSIGN expression 
{   
  identifierNode = newIdentifierNode($1)
  assignNode = newAssignNode($2)
  assignNode.Left = identifierNode
  // the expression is already a node, so we just assign it directly
  assignNode.Right = $3
}

expression: NUMBER 
{ 
  $$ = newNumberNode($1)
}

%%
func cast(y yyLexer) *MyDsl { return y.(*Lexer).p }
