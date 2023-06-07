package ru.rofleksey.refl.util;

import org.jetbrains.annotations.NotNull;

public record Pair<A, B>(@NotNull A first, @NotNull B second) {
}
