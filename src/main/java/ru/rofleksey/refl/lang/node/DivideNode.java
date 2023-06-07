package ru.rofleksey.refl.lang.node;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

public class DivideNode implements Node {
    private final Node left;
    private final Node right;

    public DivideNode(Node left, Node right) {
        this.left = left;
        this.right = right;
    }


    @Override
    public @NotNull Value evaluate(ReflContext ctx) throws EvalError {
        return left.evaluate(ctx).divide(right.evaluate(ctx));
    }
}
