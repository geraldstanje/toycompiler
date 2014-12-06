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
  $$ = newProgramNode($1, $2)

  cast(yylex).SetAstRoot($$)
}
| statement END_LINE
{
  $$ = newProgramNode($1, nil)

  cast(yylex).SetAstRoot($$)
}
| statement END_LINE program 
{
  $$ = newProgramNode($1, $3)

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
  $$ = newPrintNode($2)
}

assignation: IDENTIFIER ASSIGN expression 
{   
  tokenNode := newTokenNode($1)
  $$ = newAssignNode(tokenNode, $3) // assign.Right = $3 ... the expression is already a node, so we just assign it directly
}

structure: WHILE expression BEGIN_BLOCK program END_BLOCK
{
  $$ = newWhileNode($2, $4)
}

expression: NUMBER 
{ 
  $$ = newTokenNode($1)
}
| BEGIN_EXPRESSION expression END_EXPRESSION
{
  $$ = newProgramNode($2, nil)
}
| expression ADD_OP expression
{
  $$ = newOpNode($2, $1, $3)
}
| expression MUL_OP expression
{
  $$ = newOpNode($2, $1, $3)
}
| IDENTIFIER
{
  $$ = newTokenNode($1)
}
               
%%
func cast(y yyLexer) *Compiler { return y.(*Lexer).p }