let
    // functions

    // takes a list and returns everything after the first element
    tail = func(list []) []
        if len(list) <= 1 then [] else list[1:],

    // takes a function and a list and returns
    // true if all elements in the list satisfy the function
    all = func(condition func(any) bool, list []) bool
        if len(list) is 0 then
            true
        else
            let
                this = list[0]
            in
                condition(this) and all(condition, tail(list)),

    // takes a function and a list and returns a list
    // containing all the elements of the given list up to
    // the first element that does not satisfy the function
    // unlike filter, takeWhile stops at the first element
    // that does not satisfy the function
    takeWhile = func(condition func(any) bool, list []) []
        if len(list) is 0 then
            []
        else
            let this = list[0] in
                if condition(this) then
                    prepend(this, takeWhile(condition, tail(list)))
                else
                    [],

    isFactor = func(n, x) bool
        n % x is 0,

    // takes a number and returns true if it is prime
    isPrime = func(n int) bool
        if n <= 2 then
            n is 2
        else all(
            func(x) not isFactor(n, x),
            takeWhile(
                func(x int) not (x * x > n),
                range(2, n - 1)
            )
        ),

    // takes a function and a list and returns a list
    // containing all the elements of the given list that
    // satisfy the function
    // unlike takeWhile, filter does not stop at the
    // first element that does not satisfy the function
    filter = func(f func(any) bool, l []) []
        if len(l) is 0 then
            []
        else
            let
                this = l[0],
                rest = filter(f, tail(l))
            in
                if f(this) then
                    prepend(this, rest)
                else
                    rest,

    // helper function for fib
    fibHelper = func(n int, a int, b int) int
        if n is 0 then
            a
        else
            fibHelper(n - 1, b, a + b),

    // takes a number n and returns the nth Fibonacci number
    fib = func(n int) int fibHelper(n, 1, 1),

    // takes a number n and returns a list of the first n Fibonacci numbers
    firstNFibs = func(n int) [int] fib(x) for x in range(0, n),
    
    // variables
    lim = parseInt(input("Enter a limit: "))
in if lim < 0 then
    "Limit must be non-negative"
else if lim > 40 then
    "Limit must be less than or equal to 40"
else let
        fibs = firstNFibs(lim)
    in [
        "\n", "fibs:", fibs, "\n",
        "prime fibs:", filter(isPrime, fibs), "\n"
    ]