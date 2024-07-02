let
    // general utils
    tail = func(l)
        if len(l) <= 1 then
            []
        else
            l[1:],
    
    filter = func(l, f)
        if len(l) is 0 then
            []
        else
            let
                first = l[0],
                rest = filter(tail(l), f)
            in
                if f(first) then
                    prepend(first, rest)
                else
                    rest,

    min = func(l)
        if len(l) is 0 then
            // false indicates no minumum
            false
        else
            let
                first = l[0],
                minRest = min(tail(l))
            in
                if minRest is false or first <= minRest.min then {
                    min: first,
                    idx: 0
                } else {
                    min: minRest.min,
                    idx: minRest.idx + 1
                },

    find = func(f, l)
        if len(l) is 0 then
            // false indicates not found
            false
        else if f(l[0]) then
            0
        else
            let
                res = find(f, tail(l))
            in
                if res is false then
                    false
                else
                    res + 1,

    concatAll = func(l)
        if len(l) is 0 then
            ""
        else
            concatStr(l[0], concatAll(tail(l))),
    
    withIdxAs = func(l, i, v)
        if i >= len(l) then l else
            concat(append(l[:i], v), l[i+1:]),
    
    abs = func(a) if a < 0 then -1 * a else a,

    dist = func(a, b) abs(a.x - b.x) + abs(a.y - b.y),

    push = func(queue, item, cmp)
        //let _ = print(concatAll(["pushing ", toString(item), " onto ", toString(queue)])) in
        if len(queue) is 0 then
            //let _ = print(concatAll(["queue is empty; returning [", toString(item), "]"])) in
            [item]
        else let
            idx = len(queue) / 2,
            cmpRes = cmp(item, queue[idx])
            //_ = print(concatAll(["comparison with ", toString(queue[idx]), " at index ", toString(idx), " is ", toString(cmpRes)]))
        in if cmpRes then //let _ = print("pushing onto left") in
            concat(push(queue[:idx], item, cmp), queue[idx:])
        else //let _ = print("pushing onto right") in
            concat(queue[:idx+1], push(queue[idx+1:], item, cmp))
in {
    tail: tail,
    filter: filter,
    min: min,
    find: find,
    concatAll: concatAll,
    withIdxAs: withIdxAs,
    abs: abs,
    dist: dist,
    push: push
}