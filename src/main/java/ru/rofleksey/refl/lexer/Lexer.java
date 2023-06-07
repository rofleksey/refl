package ru.rofleksey.refl.lexer;

import ru.rofleksey.refl.lexer.lexem.*;

import java.util.ArrayList;
import java.util.List;

public class Lexer {
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

    private int parseNumber(String text) throws LexerError {
        var startIndex = curPos;
        curPos++;
        while (curPos < text.length() && Character.isDigit(text.charAt(curPos))) {
            curPos++;
        }
        curPos--;
        var numberStr = text.substring(startIndex, curPos + 1);
        try {
            return Integer.parseInt(numberStr);
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
        if (curPos >= text.length() || curPos - startPos < 2 || text.charAt(curPos) != curChar) {
            throw new LexerError("Closing quote not found, begins at " + startPos);
        }
        return builder.toString();
    }

    private void parseOther(String text) throws LexerError {
        switch (curChar) {
            case '?' -> result.add(QuestionLexem.INSTANCE);
            case '$' -> result.add(DollarLexem.INSTANCE);
            case '+' -> result.add(PlusLexem.INSTANCE);
            case '-' -> result.add(MinusLexem.INSTANCE);
            case '*' -> result.add(MultiplyLexem.INSTANCE);
            case '/' -> result.add(DivideLexem.INSTANCE);
            case '&' -> result.add(AndLexem.INSTANCE);
            case '|' -> result.add(OrLexem.INSTANCE);
            case '.' -> result.add(DotLexem.INSTANCE);
            case ',' -> result.add(CommaLexem.INSTANCE);
            case ':' -> result.add(ColonLexem.INSTANCE);
            case '!' -> result.add(NotLexem.INSTANCE);
            case '<' -> result.add(LtLexem.INSTANCE);
            case '>' -> result.add(GtLexem.INSTANCE);
            case '(' -> result.add(BracketOpenLexem.INSTANCE);
            case ')' -> result.add(BracketCloseLexem.INSTANCE);
            default -> {
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
                    case "if" -> result.add(IfLexem.INSTANCE);
                    case "while" -> result.add(WhileLexem.INSTANCE);
                    case "refl" -> result.add(ReflLexem.INSTANCE);
                    case "end" -> result.add(EndLexem.INSTANCE);
                    default -> result.add(new VarLiteral(str));
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
