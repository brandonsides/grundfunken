let
    tail = func(list) if len(list) <= 1 then [] else list[1:],

    splitHelper = func(string) if lenStr(string) is 0 then {
        first: "",
        rest: ""
    } else let
            this = atStr(string, 0),
            rest = if lenStr(string) is 1 then "" else sliceStr(string, 1, -1)
        in if this is " " then {
            first: "",
            rest: rest
        } else let
            splitRest = splitHelper(rest) in {
            first: concatStr(this, splitRest.first),
            rest: splitRest.rest
        },

    split = func(string)
        if lenStr(string) is 0 then
            []
        else let firstAndRest = splitHelper(string),
                first = firstAndRest.first,
                rest = firstAndRest.rest
            in
                prepend(first, split(rest)),

    merge = func(left, right)
        if len(left) is 0 then
            right
        else if len(right) is 0 then
            left
        else if left[0] < right[0] then
            prepend(left[0], merge(tail(left), right))
        else
            prepend(right[0], merge(left, tail(right))),

    sort = func(nums)
        if len(nums) <= 1 then
            nums
        else let
            mid = len(nums) / 2,
            left = nums[:mid],
            right = nums[mid:],
            sortedLeft = sort(left),
            sortedRight = sort(right)
        in
            merge(sortedLeft, sortedRight),


    asNums = func(strings) parseInt(string) for string in strings,

    nums = asNums(split(input("Enter a list of numbers separated by spaces: ")))
in
    sort(nums)