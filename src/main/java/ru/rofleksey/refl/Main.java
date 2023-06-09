package ru.rofleksey.refl;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.ReflExecutor;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.FunctionValue;
import ru.rofleksey.refl.lang.value.ReflValue;
import ru.rofleksey.refl.lexer.LexerError;
import ru.rofleksey.refl.parser.error.ParserError;

import java.util.List;
import java.util.Map;
import java.util.Scanner;

public class Main {
    public static void main(String[] args) {
        var scanner = new Scanner(System.in);
        var ctx = new ReflContext();
        var compiler = new ReflCompiler();

        System.out.println("Enter commands below:");

        ctx.setVar("print", new FunctionValue("print") {
            @Override
            public Value call(ReflContext ctx, List<Value> args, Map<String, Value> namedArgs) {
                args.forEach(it -> {
                    System.out.println(it.toString());
                });
                return ReflValue.INSTANCE;
            }
        });

        ctx.setVar("help", new FunctionValue("help") {
            @Override
            public Value call(ReflContext ctx, List<Value> args, Map<String, Value> namedArgs) {
                System.out.println("Print exit() to exit.");
                return ReflValue.INSTANCE;
            }
        });

        ctx.setVar("exit", new FunctionValue("exit") {
            @Override
            public Value call(ReflContext ctx, List<Value> args, Map<String, Value> namedArgs) {
                System.exit(0);
                return ReflValue.INSTANCE;
            }
        });

        while (scanner.hasNextLine()) {
            var line = scanner.nextLine();
            var time = System.currentTimeMillis();
            ReflExecutor executor;
            try {
                executor = compiler.compile(line);
            } catch (ParserError | LexerError e) {
                System.err.println(e.getMessage());
                continue;
            }
            try {
                executor.execute(ctx);
            } catch (EvalError e) {
                System.err.println(e.getMessage());
            } finally {
                System.err.println("Processed in " + (System.currentTimeMillis() - time) + "ms");
            }
        }
    }
}