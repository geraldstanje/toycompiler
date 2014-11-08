%{
package dsl
%}

%union {
  s string
  node Node
  funcName string
  funcDecl []Node
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
%token FUNC
%%

program: declarations 
{
}

declarations: declaration
{
  declarationNode := newDeclarationNode($$.funcDecl)
  $$.node = declarationNode
  
  cast(yylex).SetAstRoot($$.node)
}

declaration: fun_declaration
{  
  functionDeclNode := newFunctionDeclNode($1.funcName, $1.node)
  $$.funcDecl = append($$.funcDecl, functionDeclNode)
}
| fun_declaration declaration
{  
  functionDeclNode := newFunctionDeclNode($1.funcName, $1.node)
  $2.funcDecl = append($2.funcDecl, functionDeclNode)
  $$ = $2
}

fun_declaration: FUNC IDENTIFIER BEGIN_EXPRESSION END_EXPRESSION block
{ 
  $$.funcName = $2.s
  $$.node = $5.node
}

block: BEGIN_BLOCK program END_BLOCK
{
  $$.node = $2.node
}

program: statement
{
  $$.node = $1.node
}
| statement END_LINE
{   
  programNode := newProgramNode($1.node, nil)
  $$.node = programNode
}
| statement END_LINE program 
{
  programNode := newProgramNode($1.node, $3.node)
  $$.node = programNode
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
  printNode := newPrintNode($2.node)
  $$.node = printNode
}

assignation: IDENTIFIER ASSIGN expression 
{   
  tokenNode := newTokenNode($1.s)
  assignNode := newAssignNode(tokenNode, $3.node) // assign.Right = $3.node ... the expression is already a node, so we just assign it directly
  $$.node = assignNode
}

structure: WHILE expression BEGIN_BLOCK program END_BLOCK
{
   whileNode := newWhileNode($2.node, $4.node)
   $$.node = whileNode
}

expression: NUMBER 
{ 
  $$.node = newTokenNode($1.s)
}
| BEGIN_EXPRESSION expression END_EXPRESSION
{
  programNode := newProgramNode($2.node, nil)
  $$.node = programNode
}
| expression ADD_OP expression
{
  opNode := newOpNode($2.s, $1.node, $3.node)
  $$.node = opNode
}
| expression MUL_OP expression
{
  opNode := newOpNode($2.s, $1.node, $3.node)
  $$.node = opNode
}
| IDENTIFIER
{
  $$.node = newTokenNode($1.s)
}
              
%%
func cast(y yyLexer) *Compiler { return y.(*Lexer).p }