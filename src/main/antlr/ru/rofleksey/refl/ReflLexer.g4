lexer grammar ReflLexer;

BREAK:              'break';
CONTINUE:           'continue';
IF:                 'if';
ELIF:               'elif';
ELSE:               'else';
END:                'end';
WHILE:              'while';
STRUCT:             'struct';
SCOPE:              'scope';
FUN:                'fun';
DEF:                'def';

DECIMAL_LITERAL:    ('0' | [1-9] (Digits? | '_'+ Digits)) [lL]?;
HEX_LITERAL:        '0' [xX] [0-9a-fA-F] ([0-9a-fA-F_]* [0-9a-fA-F])? [lL]?;
OCT_LITERAL:        '0' '_'* [0-7] ([0-7_]* [0-7])? [lL]?;
BINARY_LITERAL:     '0' [bB] [01] ([01_]* [01])? [lL]?;

FLOAT_LITERAL:      (Digits '.' Digits? | '.' Digits) ExponentPart? [fFdD]?
             |       Digits (ExponentPart [fFdD]? | [fFdD])
             ;

HEX_FLOAT_LITERAL:  '0' [xX] (HexDigits '.'? | HexDigits? '.' HexDigits) [pP] [+-]? Digits [fFdD]?;

CHAR_LITERAL:       '\'' (~['\\\r\n] | EscapeSequence)* '\'';

STRING_LITERAL:     '"' (~["\\\r\n] | EscapeSequence)* '"';

NIL_LITERAL:       'nil';

LPAREN:             '(';
RPAREN:             ')';
LBRACK:             '[';
RBRACK:             ']';
COMMA:              ',';
DOT:                '.';

ASSIGN:             '=';
GT:                 '>';
LT:                 '<';
BANG:               '!';
TILDE:              '~';
ELVIS:              '??';
COLON:              ':';
EQUAL:              '==';
LE:                 '<=';
GE:                 '>=';
NOTEQUAL:           '!=';
AND:                '&';
OR:                 '|';
INC:                '++';
DEC:                '--';
ADD:                '+';
SUB:                '-';
MUL:                '*';
DIV:                '/';
MOD:                '%';

ADD_ASSIGN:         '+=';
SUB_ASSIGN:         '-=';
MUL_ASSIGN:         '*=';
DIV_ASSIGN:         '/=';
MOD_ASSIGN:         '%=';
AND_ASSIGN:         '&=';
OR_ASSIGN:          '|=';

CHAN_ARROW:         '<-';
RETURN_ARROW:       '<<';

// Whitespace and comments

WS:                 [ \t\r\u000C]+ -> channel(HIDDEN);
LINE_COMMENT:       '#' ~[\r\n]*    -> channel(HIDDEN);
EOL:                '\n'+;

// Identifiers

IDENTIFIER:         Letter LetterOrDigit*;

// Fragment rules

fragment ExponentPart
    : [eE] [+-]? Digits
    ;

fragment EscapeSequence
    : '\\' 'u005c'? [btnfr"'\\]
    | '\\' 'u005c'? ([0-3]? [0-7])? [0-7]
    | '\\' 'u'+ HexDigit HexDigit HexDigit HexDigit
    ;

fragment HexDigits
    : HexDigit ((HexDigit | '_')* HexDigit)?
    ;

fragment HexDigit
    : [0-9a-fA-F]
    ;

fragment Digits
    : [0-9] ([0-9_]* [0-9])?
    ;

fragment LetterOrDigit
    : Letter
    | [0-9]
    ;

fragment Letter
    : [a-zA-Z$_] // these are the "java letters" below 0x7F
    | ~[\u0000-\u007F\uD800-\uDBFF] // covers all characters above 0x7F which are not a surrogate
    | [\uD800-\uDBFF] [\uDC00-\uDFFF] // covers UTF-16 surrogate pairs encodings for U+10000 to U+10FFFF
    ;