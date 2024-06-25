let
    // takes a list and returns everything after the first element
    tail = func(l)
        l[1:],

    // takes a list and returns all elements that meet the given condition
    filter = func(l, f)
        if len(l) is 0 then
            []
        else
            let
                this = l[0],
                rest = filter(tail(l), f)
            in
                if f(this) then
                    prepend(this, rest)
                else
                    rest,

    // returns the least element in a list
    min = func(l)
        if len(l) is 0 then
            [-1, 0]
        else
            let
                first = l[0],
                rest = tail(l),
                restMin = min(rest),
                restMinIdx = restMin[0],
                restMinVal = restMin[1]
            in if
                restMinIdx is -1 or first < restMinVal
            then
                [0, first]
            else
                [restMinIdx+1, restMinVal],

    // all possible coin values (in pence)
    coinVals = [1, 2, 5, 10, 20, 25, 50, 100, 200],

    // recursive helper function to make change
    makeChangeHelper = func(trg, cur, prevBests)
        let
            bestWithEachCoinValAsLast = append(
                prevBests[cur - coinVal],
                coinVal
            ) for coinVal in filter(
                coinVals,
                func(coinVal) coinVal <= cur
            ),

            curBestIdxAndVal = min(
                len(coinSeq) for coinSeq in bestWithEachCoinValAsLast
            ),
            curBestIdx = curBestIdxAndVal[0],

            curBestVal = if curBestIdx < 0 then
                []
            else
                bestWithEachCoinValAsLast[curBestIdx]
        in
            if cur is trg then
                curBestVal
            else
                makeChangeHelper(trg, cur + 1, append(prevBests, curBestVal)),

    // wrapper function to make change
    makeChange = func(n)
        if n <= 0 then
            []
        else
            makeChangeHelper(n, 1, [[]]),
    
    makeChangeGreedy = func(n)
        if n <= 0 then
            []
        else
            let
                coinVal = filter(coinVals, func(coinVal) coinVal <= n),
                rest = makeChangeGreedy(n - coinVal[len(coinVal)-1])
            in
                prepend(coinVal[len(coinVal)-1], rest),
    
    valStr = input("Enter a value in pence: "),
    val = parseInt(valStr)
in
    [
        "\n\tgreedy:", makeChangeGreedy(val),
        "\n\toptimal:", makeChange(val), "\n"
    ]
