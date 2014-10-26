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
{
  $$.expr = $1.expr
}
| statement END_LINE
{ 
  $$.expr = $1.expr
}
| statement END_LINE program 
{ 
  programNode = newProgramNode($2.expr)
  programNode.Left = $1.expr
  programNode.Right = $3.expr
}

statement: assignation 
{
  $$.expr = $1.expr
}

assignation: IDENTIFIER ASSIGN expression 
{   
  identifierNode = newIdentifierNode($1.expr)
  assignNode = newAssignNode($2.expr)
  assignNode.Left = identifierNode
  // the expression is already a node, so we just assign it directly
  assignNode.Right = $3.expr
}

expression: NUMBER 
{ 
  $$.expr = newNumberNode($1.expr)
}

%%
func cast(y yyLexer) *MyDsl { return y.(*Lexer).p }
