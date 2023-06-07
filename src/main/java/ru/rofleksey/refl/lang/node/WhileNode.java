package ru.rofleksey.refl.lang.node;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.Refl;
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;

import java.util.List;

public class WhileNode implements Node {
    private final Node condition;
    private final List<Node> body;

    public WhileNode(Node condition, List<Node> body) {
        this.condition = condition;
        this.body = body;
    }

    @Override
    public @NotNull Value evaluate(ReflContext ctx) throws EvalError {
        Value result = Refl.INSTANCE;

        while (condition.evaluate(ctx).isTruthy()) {
            for (var node : body) {
                result = node.evaluate(ctx);
            }
        }

        return result;
    }
}