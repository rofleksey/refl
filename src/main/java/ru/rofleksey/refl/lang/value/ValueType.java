package ru.rofleksey.refl.lang.value;

public enum ValueType {
    NUMBER("number"), STRING("string"), NIL("nil"), OBJECT("object"), FUNCTION("function");
    private final String type;

    ValueType(String type) {
        this.type = type;
    }

    @Override
    public String toString() {
        return type;
    }
}
