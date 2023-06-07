package ru.rofleksey.refl.lang.node;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

public class AssignNode implements Node {
    private final String name;
    private final Node node;

    public AssignNode(String name, Node node) {
        this.name = name;
        this.node = node;
    }


    @Override
    public @NotNull Value evaluate(ReflContext ctx) throws EvalError {
        var result = node.evaluate(ctx);
        ctx.setVar(name, result);
        return result;
    }
}
