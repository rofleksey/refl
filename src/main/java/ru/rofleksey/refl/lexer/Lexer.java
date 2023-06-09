package ru.rofleksey.refl.lexer;

import ru.rofleksey.refl.lexer.lexem.*;

import java.util.ArrayList;
import java.util.List;

public final class Lexer {
    private int curPos;
    private char curChar;
    private List<Lexem> result;

    private boolean lookUp(String str, int start, String what) {
        var length = 1;
        for (var i = start; i < str.length() && length <= what.length(); i++, length++) {
            if (str.charAt(i) != what.charAt(length - 1)) {
                return false;
            }
        }
        return length == what.length() + 1;
    }

    private double parseNumber(String text) throws LexerError {
        var startIndex = curPos;
        curPos++;
        while (curPos < text.length() && (Character.isDigit(text.charAt(curPos)) || text.charAt(curPos) == '.')) {
            curPos++;
        }
        curPos--;
        var numberStr = text.substring(startIndex, curPos + 1);
        try {
            return Double.parseDouble(numberStr);
        } catch (NumberFormatException e) {
            throw new LexerError("Failed to parse number '" + numberStr + "' at " + startIndex);
        }
    }

    private String parseVar(String text) {
        var startPos = curPos++;
        while (curPos < text.length() &&
                (Character.isLetter(text.charAt(curPos))
                        || Character.isDigit(text.charAt(curPos)))) {
            curPos++;
        }
        curPos--;
        return text.substring(startPos, curPos + 1);
    }

    private String parseQuotes(String text) throws LexerError {
        var startPos = curPos;
        curPos++;
        var builder = new StringBuilder();
        while (curPos < text.length()) {
            if (text.charAt(curPos) == '\\' && curPos + 1 < text.length()) {
                builder.append(text.charAt(curPos + 1));
                curPos += 2;
                continue;
            }
            if (text.charAt(curPos) == curChar) {
                break;
            }
            builder.append(text.charAt(curPos));
            curPos++;
        }
        if (curPos >= text.length() || text.charAt(curPos) != curChar) {
            throw new LexerError("Closing quote not found, begins at " + startPos);
        }
        return builder.toString();
    }

    private void parseOther(String text) throws LexerError {
        switch (curChar) {
            case '+':
                result.add(PlusLexem.INSTANCE);
                break;

            case '-':
                result.add(MinusLexem.INSTANCE);
                break;

            case '*':
                result.add(MultiplyLexem.INSTANCE);
                break;

            case '/':
                result.add(DivideLexem.INSTANCE);
                break;

            case '&':
                result.add(AndLexem.INSTANCE);
                break;

            case '|':
                result.add(OrLexem.INSTANCE);
                break;

            case ',':
                result.add(CommaLexem.INSTANCE);
                break;

            case ':':
                result.add(ColonLexem.INSTANCE);
                break;

            case ';':
                result.add(SemicolonLexem.INSTANCE);
                break;

            case '!':
                result.add(NotLexem.INSTANCE);
                break;

            case '<':
                result.add(LtLexem.INSTANCE);
                break;

            case '>':
                result.add(GtLexem.INSTANCE);
                break;

            case '(':
                result.add(BracketOpenLexem.INSTANCE);
                break;

            case ')':
                result.add(BracketCloseLexem.INSTANCE);
                break;

            default:
                if (lookUp(text, curPos, "==")) {
                    result.add(EqLexem.INSTANCE);
                    curPos++;
                    return;
                } else if (curChar == '=') {
                    result.add(AssignLexem.INSTANCE);
                    return;
                }
                throw new LexerError("Unexpected symbol '" + curChar + "' at " + curPos);
        }
    }

    public List<Lexem> process(String text) throws LexerError {
        result = new ArrayList<>();

        for (curPos = 0; curPos < text.length(); curPos++) {
            curChar = text.charAt(curPos);

            if (Character.isWhitespace(curChar)) {
                continue;
            }

            if (Character.isDigit(curChar)) {
                var number = parseNumber(text);
                result.add(new NumberLexem(number));
                continue;
            }

            if (Character.isLetter(curChar)) {
                var str = parseVar(text);
                switch (str) {
                    case "if":
                        result.add(IfLexem.INSTANCE);
                        break;

                    case "while":
                        result.add(WhileLexem.INSTANCE);
                        break;

                    case "refl":
                        result.add(ReflLexem.INSTANCE);
                        break;

                    case "end":
                        result.add(EndLexem.INSTANCE);
                        break;

                    default:
                        result.add(new VarLexem(str));
                        break;
                }
                continue;
            }

            if (curChar == '\'' || curChar == '"') {
                var str = parseQuotes(text);
                result.add(new StringLexem(str));
                continue;
            }

            parseOther(text);
        }

        result.add(EofLexem.INSTANCE);

        return result;
    }
}
