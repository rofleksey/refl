package ru.rofleksey.refl.lang.node;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

public class AddNode implements Node {
    private final Node left;
    private final Node right;

    public AddNode(Node left, Node right) {
        this.left = left;
        this.right = right;
    }


    @Override
    public @NotNull Value evaluate(ReflContext ctx) throws EvalError {
        return left.evaluate(ctx).add(right.evaluate(ctx));
    }
}
