package ru.rofleksey.refl.util;

import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.ReentrantLock;

public class HandshakeChannel<T> {
    private final ReentrantLock lock = new ReentrantLock();
    private final Condition condition = lock.newCondition();
    private volatile T passedValue = null;

    public T read() throws InterruptedException {
        lock.lockInterruptibly();

        try {
            condition.awaitUninterruptibly();
            var result = passedValue;
            passedValue = null;
            return result;
        } finally {
            lock.unlock();
        }
    }

    public boolean write(T value) {
        lock.lock();

        try {
            if (lock.hasWaiters(condition)) {
                passedValue = value;
                condition.signal();
                return true;
            }

            return false;
        } finally {
            lock.unlock();
        }
    }
}
