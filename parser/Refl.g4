grammar Refl;

// Parser rules
program: (statement)* EOF;

statement
    : varDeclaration
    | expressionStatement
    | ifStatement
    | whileStatement
    | forStatement
    | blockStatement
    | breakStatement
    | continueStatement
    | returnStatement
    ;

varDeclaration: 'var' IDENTIFIER '=' expression;
expressionStatement: expression;

ifStatement: 'if' expression block ('elif' expression block)* ('else' block)?;
whileStatement: 'while' expression block;
forStatement: 'for' IDENTIFIER ',' IDENTIFIER 'in' expression block;
breakStatement: 'break';
continueStatement: 'continue';
returnStatement: 'return' expression?;

block: '{' statement* '}';
blockStatement: block;

expression
    : expression '.' IDENTIFIER                                   # memberDot
    | expression ':' IDENTIFIER '(' expressionList? ')'           # methodCall
    | expression '(' expressionList? ')'                          # functionCall
    | expression '[' expression ']'                               # memberBracket
    | op='-' expression                                           # unary
    | op='!' expression                                           # unary
    | expression op=('*' | '/' | '%') expression                  # binary
    | expression op=('+' | '-') expression                        # binary
    | expression op=('<' | '>' | '<=' | '>=') expression          # binary
    | expression op=('==' | '!=') expression                      # binary
    | expression op='&&' expression                               # binary
    | expression op='||' expression                               # binary
    | expression '=' expression                                   # assignment
    | primary                                                     # primaryExpr
    ;

expressionList: expression (',' expression)*;

primary
    : literal                                                     # literalPrimary
    | IDENTIFIER                                                  # identifierPrimary
    | '(' expression ')'                                          # parenPrimary
    | 'fun' '(' parameters? ')' block                             # functionLiteral
    | objectLiteral                                               # objectLiteralPrimary
    | arrayLiteral                                                # arrayLiteralPrimary
    ;

parameters: IDENTIFIER (',' IDENTIFIER)*;

objectLiteral: '{' (property (',' property)*)? '}';
property: (STRING | IDENTIFIER) ':' expression;

arrayLiteral: '{' (expression (',' expression)*)? '}';

literal
    : NUMBER                                                      # numberLiteral
    | STRING                                                      # stringLiteral
    | RAW_STRING                                                  # rawStringLiteral
    | 'nil'                                                       # nilLiteral
    ;

// Keywords
VAR: 'var';
IF: 'if';
ELIF: 'elif';
ELSE: 'else';
WHILE: 'while';
FOR: 'for';
IN: 'in';
BREAK: 'break';
CONTINUE: 'continue';
RETURN: 'return';
FUN: 'fun';
NIL: 'nil';

// Lexer rules
STRING: '"' (~["\\\r\n] | '\\' ["\\nrt])* '"';
RAW_STRING: '`' ~[`]* '`';
NUMBER: [0-9]+ ('.' [0-9]+)?;

IDENTIFIER: [a-zA-Z_$][a-zA-Z0-9_$]*;

// Comments
LINE_COMMENT: '#' ~[\r\n]* -> skip;
WHITESPACE: [ \t\r\n]+ -> skip;

// Operators
DOT: '.';
COLON: ':';
LPAREN: '(';
RPAREN: ')';
LBRACKET: '[';
RBRACKET: ']';
LBRACE: '{';
RBRACE: '}';
COMMA: ',';
SEMICOLON: ';'; // For error checking - we don't allow semicolons
ASSIGN: '=';
PLUS: '+';
MINUS: '-';
ASTERISK: '*';
SLASH: '/';
PERCENT: '%';
BANG: '!';
LT: '<';
GT: '>';
LE: '<=';
GE: '>=';
EQ: '==';
NE: '!=';
AND: '&&';
OR: '||';