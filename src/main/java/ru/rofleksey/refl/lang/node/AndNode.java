package ru.rofleksey.refl.lang.node;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.NumberValue;

public class AndNode implements Node {
    private final Node left;
    private final Node right;

    public AndNode(Node left, Node right) {
        this.left = left;
        this.right = right;
    }


    @Override
    public @NotNull Value evaluate(ReflContext ctx) throws EvalError {
        var leftValue = left.evaluate(ctx);
        if (!leftValue.isTruthy()) {
            return NumberValue.FALSE;
        }
        return leftValue.and(right.evaluate(ctx));
    }
}
