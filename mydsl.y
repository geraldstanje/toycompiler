%{
package dsl
%}

%union {
  s string
  node *Node
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
  $$.node = $1.node
}
| statement END_LINE
{ 
  $$.node = $1.node
}
| statement END_LINE program 
{ 
  programNode := newProgramNode($2)
  programNode.Left = $1.node
  programNode.Right = $3.node
  $$.node = programNode

  cast(yylex).InitAST($$.node)
}

statement: assignation 
{
  $$.node = $1.node
}

assignation: IDENTIFIER ASSIGN expression 
{   
  identifierNode := newIdentifierNode($1)
  assignNode := newAssignNode($2)
  assignNode.Left = identifierNode
  // the expression is already a node, so we just assign it directly
  assignNode.Right = $3.node

  $$.node = assignNode
}

expression: NUMBER 
{ 
  $$.node = newNumberNode($1)
}

%%
func cast(y yyLexer) *MyDsl { return y.(*Lexer).p }
