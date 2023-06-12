package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

public interface Node {
    Value evaluate(ReflContext ctx) throws EvalError;

    Value getLeftSide(ReflContext ctx) throws EvalError;

    Value getSetterKey(ReflContext ctx) throws EvalError;
}
