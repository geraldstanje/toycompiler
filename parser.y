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
%type <node> program assignation structure expression statement print statement_list
%%
program: statement_list 
{
  $$ = newProgramNode($1)
  cast(yylex).SetAstRoot($$)
}

statement_list: statement
{
  $$ = newStatementNode($1, nil)
}
| statement statement_list
{
  $$ = newStatementNode($1, $2)
}

statement: assignation END_LINE
{
  $$ = $1
}
| structure
{
  $$ = $1
}
| print END_LINE
{
  $$ = $1
}

print: PRINT expression
{
  $$ = newPrintNode($2)
}

assignation: IDENTIFIER ASSIGN expression 
{   
  tokenNode := newTokenNode($1)
  $$ = newAssignNode(tokenNode, $3) // assign.Right = $3 ... the expression is already a node, so we just assign it directly
}

structure: WHILE BEGIN_EXPRESSION expression END_EXPRESSION BEGIN_BLOCK statement_list END_BLOCK
{
  $$ = newWhileNode($3, $6)
}

expression: NUMBER 
{ 
  $$ = newTokenNode($1)
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