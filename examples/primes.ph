let

tail = func(l)
    slice(l, 1, -1),

all = func(f, l)
    if equals(len(l), 0) then
        true
    else
        let this = at(l, 0) in
            if f(this) then
                all(f, tail(l))
            else
                false,

takeWhile = func(f, l)
    if equals(len(l), 0) then
        []
    else
        let this = at(l, 0) in
            if f(this) then
                prepend(this, takeWhile(f, tail(l)))
            else
                [],

isPrime = func(n)
    if lessThan(n, 2) then
        false
    else if equals(n, 2) then
        true
    else all(
        func(x) not(equals(0, mod(n, x))),
        takeWhile(
            func(x) not(greaterThan(x * x, n)),
            range(2, n - 1)
        )
    ),

lim = 1000

in

[
    "primes:", filter(isPrime, range(2, lim))
]