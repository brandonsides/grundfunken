let
    // functions

    // takes a list and returns everything after the first element
    tail = func(l)
        slice(l, 1, -1),

    // takes a function and a list and returns
    // true if all elements in the list satisfy the function
    all = func(f, l)
        if equals(len(l), 0) then
            true
        else
            let this = at(l, 0) in
                and(
                    f(this),
                    all(f, tail(l))
                ),

    // takes a function and a list and returns a list
    // containing all the elements of the given list up to
    // the first element that does not satisfy the function
    // unlike filter, takeWhile stops at the first element
    // that does not satisfy the function
    takeWhile = func(f, l)
        if equals(len(l), 0) then
            []
        else
            let this = at(l, 0) in
                if f(this) then
                    prepend(this, takeWhile(f, tail(l)))
                else
                    [],

    isFactor = func(n, x) equals(0, mod(n, x)),

    // takes a number and returns true if it is prime
    isPrime = func(n)
        if lessThan(n, 2) then
            false
        else if equals(n, 2) then
            true
        else all(
            func(x) not(isFactor(n, x)),
            takeWhile(
                func(x) not(greaterThan(x * x, n)),
                range(2, n - 1)
            )
        ),

    // takes a function and a list and returns a list
    // containing all the elements of the given list that
    // satisfy the function
    // unlike takeWhile, filter does not stop at the
    // first element that does not satisfy the function
    filter = func(f, l)
        if equals(len(l), 0) then
            []
        else
            let
                this = at(l, 0),
                rest = filter(f, tail(l))
            in
                if f(this) then
                    prepend(this, rest)
                else
                    rest,

    // helper function for fib
    fibHelper = func(n, a, b)
        if equals(n, 0) then
            a
        else
            fibHelper(n - 1, b, a + b),

    // takes a number n and returns the nth Fibonacci number
    fib = func(n) fibHelper(n, 1, 1),

    // takes a number n and returns a list of the first n Fibonacci numbers
    firstNFibs = func(n) fib(x) for x in range(0, n)
in let
    // variables
    lim = parseInt(input("Enter a limit: "))
in if lessThan(lim, 0) then
    "Limit must be non-negative"
else if greaterThan(lim, 40) then
    "Limit must be less than or equal to 40"
else let fibs = firstNFibs(lim)
in [
        "\n", "fibs:", fibs, "\n",
        "prime fibs:", filter(isPrime, fibs), "\n"
]