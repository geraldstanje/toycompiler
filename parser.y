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
  programNode := newProgramNode()
  programNode.Left = $1.node
  $$.node = programNode
  
  cast(yylex).SetAstRoot($$.node)
}
| statement END_LINE program 
{
  programNode := newProgramNode()
  programNode.Left = $1.node
  programNode.Right = $3.node
  $$.node = programNode
  
  cast(yylex).SetAstRoot($$.node)
}

statement: assignation 
{
  $$.node = $1.node
}
| structure
{
  $$.node = $1.node
}
| PRINT expression
{
  printNode := newPrintNode()
  printNode.Left = $2.node
  $$.node = printNode
}

assignation: IDENTIFIER ASSIGN expression 
{   
  tokenNode := newTokenNode($1)
  assignNode := newAssignNode($2)
  assignNode.Left = tokenNode
  // the expression is already a node, so we just assign it directly
  assignNode.Right = $3.node
  $$.node = assignNode
}

structure: WHILE expression BEGIN_BLOCK program END_BLOCK
{
   whileNode := newWhileNode()
   whileNode.Left = $2.node
   whileNode.Right = $4.node
   $$.node = whileNode
}

expression: NUMBER 
{ 
  $$.node = newTokenNode($1)
}
| BEGIN_EXPRESSION expression END_EXPRESSION
{
  programNode := newProgramNode()
  programNode.Left = $2.node
  $$.node = programNode
}
| expression ADD_OP expression
{
  opNode := newOpNode($2)
  opNode.Left = $1.node
  opNode.Right = $3.node
  $$.node = opNode
}
| expression MUL_OP expression
{
  opNode := newOpNode($2)
  opNode.Left = $1.node
  opNode.Right = $3.node
  $$.node = opNode
}
| IDENTIFIER
{
  $$.node = newTokenNode($1)
}
               
%%
func cast(y yyLexer) *Compiler { return y.(*Lexer).p }