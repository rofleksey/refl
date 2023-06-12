package ru.rofleksey.refl.lang.parse;

import ru.rofleksey.refl.lang.node.Node;

import java.util.List;
import java.util.Map;

public class MethodCallParseCtx {
    private final String name;
    private final List<Node> posArgs;
    private final Map<String, Node> namedArgs;

    public MethodCallParseCtx(String name, List<Node> posArgs, Map<String, Node> namedArgs) {
        this.name = name;
        this.posArgs = posArgs;
        this.namedArgs = namedArgs;
    }

    public String getName() {
        return name;
    }

    public List<Node> getPosArgs() {
        return posArgs;
    }

    public Map<String, Node> getNamedArgs() {
        return namedArgs;
    }
}
