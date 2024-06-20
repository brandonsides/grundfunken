let
    // takes a list and returns everything after the first element
    tail = func(l)
        slice(l, 1, -1),

    filter = func(l, f)
        if equals(len(l), 0) then
            []
        else
            let
                this = at(l, 0),
                rest = filter(tail(l), f)
            in
                if f(this) then
                    prepend(this, rest)
                else
                    rest,

    min = func(l)
        if len(l) == 0 then
            [false, 0]
        else
            let
                first = at(l, 0),
                rest = tail(l)
                restMin = min(rest),
                restHasMin = at(restMin, 0),
                restMinVal = at(restMin, 1)
            in if or(
                not(restHasMin),
                lessThan(first, restMinVal)
            ) then
                [true, first]
            else
                restMinVal,

    coinVals = [1, 2, 5, 10, 20, 50, 100, 200],

    makeChangeHelper = func(trg, cur, prevBests)
        let curBest = min(
            (at(prevWays, cur - coinVal) + 1) for coinVal in filter(
                coinVals, func(coinVal) not(greaterThan(coinVal, cur))
            )
        ) in
            if equals(cur, trg) then
                curWays
            else
                makeChangeHelper(trg, cur + 1, append(prevWays, curWays)),

    makeChange = func(n)
        if equals(n, 0) then
            []
        else
            makeChangeHelper(n, 1, [[]])
in
    makeChange(i) for i in range(0, 20)

//

// 1

// 1 1
// 2

// 1 1 1
// 2 1
// 1 2 <-- wrong

// 1 1 1 1
// 2 1 1
// 1 2 1 <-- wrong
// 1 1 2 <-- wrong
// 2 2

// 1 1 1 1 1
// 2 1 1 1
// 1 1 1 2 <-- wrong
// 2 2 1
// 1 2 2 <-- wrong
// 2 1 2 <-- wrong
// 5

// 1 1 1 1 1 1