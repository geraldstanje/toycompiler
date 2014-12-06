%{
package dsl
%}

%union {
  s string
  node Node
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

%token <s> IDENTIFIER STRING NUMBER ADD_OP MUL_OP
%type <node> program assignation structure expression statement
%%
program: statement
{
  $$ = $1
}
| statement program
{
  programNode := newProgramNode($1, $2)
  $$ = programNode

  cast(yylex).SetAstRoot($$)
}
| statement END_LINE
{   
  programNode := newProgramNode($1, nil)
  $$ = programNode

  cast(yylex).SetAstRoot($$)
}
| statement END_LINE program 
{
  programNode := newProgramNode($1, $3)
  $$ = programNode

  cast(yylex).SetAstRoot($$)
}

statement: assignation 
{
  $$ = $1
}
| structure
{
  $$ = $1
}
| PRINT expression
{
  printNode := newPrintNode($2)
  $$ = printNode
}

assignation: IDENTIFIER ASSIGN expression 
{   
  tokenNode := newTokenNode($1)
  assignNode := newAssignNode(tokenNode, $3) // assign.Right = $3 ... the expression is already a node, so we just assign it directly
  $$ = assignNode
}

structure: WHILE expression BEGIN_BLOCK program END_BLOCK
{
  whileNode := newWhileNode($2, $4)
  $$ = whileNode
}

expression: NUMBER 
{ 
  $$ = newTokenNode($1)
}
| BEGIN_EXPRESSION expression END_EXPRESSION
{
  programNode := newProgramNode($2, nil)
  $$ = programNode
}
| expression ADD_OP expression
{
  opNode := newOpNode($2, $1, $3)
  $$ = opNode
}
| expression MUL_OP expression
{
  opNode := newOpNode($2, $1, $3)
  $$ = opNode
}
| IDENTIFIER
{
  $$ = newTokenNode($1)
}
               
%%
func cast(y yyLexer) *Compiler { return y.(*Lexer).p }