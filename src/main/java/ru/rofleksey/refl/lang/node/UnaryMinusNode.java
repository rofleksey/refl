package ru.rofleksey.refl.lang.node;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.NumberValue;

public class UnaryMinusNode implements Node {
    private final Node child;

    public UnaryMinusNode(Node child) {
        this.child = child;
    }


    @Override
    public @NotNull Value evaluate(ReflContext ctx) throws EvalError {
        return child.evaluate(ctx).multiply(NumberValue.MINUS_ONE);
    }
}
