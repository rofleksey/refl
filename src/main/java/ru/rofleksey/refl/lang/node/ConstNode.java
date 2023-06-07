package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

public class ConstNode implements Node {
    private final Value value;

    public ConstNode(Value value) {
        this.value = value;
    }


    @Override
    public  Value evaluate(ReflContext ctx) throws EvalError {
        return value;
    }

    @Override
    public String toString() {
        return value.toString();
    }
}
