#!/bin/bash
assert() {
    expected="$1"
    input="$2"

    ./9cc "$input" > tmp.s
    cc -o tmp tmp.s
    ./tmp
    actual="$?"
    if [ "$actual" = "$expected" ]; then
        ehco "$input => $actual"
    else
        echo "$input => $expected expected, but get $actual"
        exit 1
    fi
}

assert 0 0
assert 42 42

echo OK
