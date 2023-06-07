package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

public class SubtractNode implements Node {
    private final Node left;
    private final Node right;

    public SubtractNode(Node left, Node right) {
        this.left = left;
        this.right = right;
    }


    @Override
    public  Value evaluate(ReflContext ctx) throws EvalError {
        return left.evaluate(ctx).subtract(right.evaluate(ctx));
    }

    @Override
    public String toString() {
        return "(" + left.toString() + "-" + right.toString() + ")";
    }
}
