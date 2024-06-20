let
    // takes a list and returns everything after the first element
    tail = func(l)
        slice(l, 1, -1),

    // takes a list and returns all elements that meet the given condition
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

    // returns the least element in a list
    min = func(l)
        if equals(len(l), 0) then
            [-1, 0]
        else
            let
                first = at(l, 0),
                rest = tail(l),
                restMin = min(rest),
                restMinIdx = at(restMin, 0),
                restMinVal = at(restMin, 1)
            in if or(
                equals(restMinIdx, -1),
                lessThan(first, restMinVal)
            ) then
                [0, first]
            else
                [restMinIdx+1, restMinVal],

    // all possible coin values (in pence)
    coinVals = [1, 2, 5, 10, 20, 25, 50, 100, 200],

    // recursive helper function to make change
    makeChangeHelper = func(trg, cur, prevBests)
        let
            bestWithEachCoinValAsLast = append(
                at(prevBests, cur - coinVal),
                coinVal
            ) for coinVal in filter(
                coinVals,
                func(coinVal) not(greaterThan(coinVal, cur))
            ),

            curBestIdxAndVal = min(
                len(coinSeq) for coinSeq in bestWithEachCoinValAsLast
            ),
            curBestIdx = at(curBestIdxAndVal, 0),

            curBestVal = if lessThan(curBestIdx, 0)
                then []
            else
                at(bestWithEachCoinValAsLast, curBestIdx)
        in
            if equals(cur, trg) then
                curBestVal
            else
                makeChangeHelper(trg, cur + 1, append(prevBests, curBestVal)),

    // wrapper function to make change
    makeChange = func(n)
        if equals(n, 0) then
            []
        else
            makeChangeHelper(n, 1, [[]]),
    
    makeChangeGreedy = func(n)
        if equals(n, 0) then
            []
        else
            let
                coinVal = filter(coinVals, func(coinVal) not(greaterThan(coinVal, n))),
                rest = makeChangeGreedy(n - at(coinVal, len(coinVal)-1))
            in
                prepend(at(coinVal, len(coinVal)-1), rest),
    
    valStr = input("Enter a value in pence: "),
    val = parseInt(valStr)
in
    ["\n\tgreedy:", makeChangeGreedy(val), "\n\toptimal:", makeChange(val), "\n"]
