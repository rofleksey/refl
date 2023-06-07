package ru.rofleksey.refl.lang.node;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

public class NotNode implements Node {
    private final Node child;

    public NotNode(Node child) {
        this.child = child;
    }


    @Override
    public @NotNull Value evaluate(ReflContext ctx) throws EvalError {
        return child.evaluate(ctx).not();
    }
}
