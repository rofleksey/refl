parser grammar ReflParser;

options { tokenVocab=ReflLexer; }

root
  : declarationList EOL* EOF
  ;

declarationList
    : (declaration EOL)* declaration
    ;

declaration
    : controlExpression
    | functionDeclaration
    | structDeclaration
    | simpleExpression
    ;

controlExpression
    : ifExpression
    | whileExpression
    ;

ifExpression
    : IF mainCondition=simpleExpression EOL mainBody=declarationList (EOL ELIF simpleExpression EOL declarationList)* (EOL ELSE EOL declarationList)? EOL END
    ;

whileExpression
    : WHILE condition=simpleExpression EOL declarationList EOL END
    ;

functionDeclaration
    : normalFunction
    ;

normalFunction
    : FUN name=IDENTIFIER EOL declarationList EOL END
    ;

structDeclaration
    : scope
    ;

scope
    : SCOPE name=IDENTIFIER EOL declarationList EOL END
    ;

simpleExpression
    : primaryExpression #primary0
    | left=simpleExpression LBRACK right=simpleExpression RBRACK #arrayCall
    | simpleExpression bop=DOT
      (
         name=IDENTIFIER
       | methodCall
      ) #dot
    | methodCall #methodCall0
    | simpleExpression postfix=(INC | DEC) #incDec
    | prefix=(ADD | SUB | INC | DEC | BANG) simpleExpression #unary
    | left=simpleExpression bop=(MUL | DIV | MOD) right=simpleExpression #mulDiv
    | left=simpleExpression bop=(ADD | SUB) right=simpleExpression #addSub
    | left=simpleExpression bop=(LT | GT | LE | GE | EQUAL | NOTEQUAL) right=simpleExpression #compare
    | left=simpleExpression bop=AND right=simpleExpression #and
    | left=simpleExpression bop=(OR | ELVIS) right=simpleExpression #or
    | <assoc=right> left=simpleExpression
      bop=(ASSIGN | ADD_ASSIGN | SUB_ASSIGN | MUL_ASSIGN | DIV_ASSIGN | MOD_ASSIGN | AND_ASSIGN | OR_ASSIGN)
      right=simpleExpression #assign
    | bop=(CHAN_ARROW | RETURN_ARROW) simpleExpression #arrows
    ;

methodCall
    : name=IDENTIFIER LPAREN argument* RPAREN
    ;

argument
    : simpleExpression
    | name=IDENTIFIER TILDE simpleExpression
    ;

primaryExpression
    : LPAREN simpleExpression RPAREN #parenExpression
    | IDENTIFIER #identifier
    | literal #literal0
    ;

literal
    : integerLiteral #numberLiteral0
    | floatLiteral #numberLiteral0
    | CHAR_LITERAL #stringLiteral
    | STRING_LITERAL #stringLiteral
    | NIL_LITERAL #nilLiteral
    ;

integerLiteral
    : DECIMAL_LITERAL
    | HEX_LITERAL
    | OCT_LITERAL
    | BINARY_LITERAL
    ;

floatLiteral
    : FLOAT_LITERAL
    | HEX_FLOAT_LITERAL
    ;
