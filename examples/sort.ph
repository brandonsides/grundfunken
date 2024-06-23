let
    // takes a list and returns everything after the first (i.e., 0th) element
    tail = func(list)
        if len(list) <= 1 then
            []
        else
            list[1:],

    // helper function for split
    // returns the first word in a string and the rest of the string, omitting the intervening space
    splitHelper = func(string)
        if lenStr(string) is 0 then 
            {
                first: "",
                rest: ""
            }
        else
            let this = atStr(string, 0),
                rest = if lenStr(string) is 1 then
                        ""
                    else
                        sliceStr(string, 1, -1)
            in
                if this is " " then 
                    {
                        first: "",
                        rest: rest
                    }
                else
                    let splitRest = splitHelper(rest) in {
                            first: concatStr(this, splitRest.first),
                            rest: splitRest.rest
                        },

    // takes a string and returns a list of words in the string
    split = func(string)
        if lenStr(string) is 0 then
            []
        else let firstAndRest = splitHelper(string),
                first = firstAndRest.first,
                rest = firstAndRest.rest
            in
                prepend(first, split(rest)),

    // takes two lists, presumed to be sorted, and merges them into a single sorted list
    merge = func(left, right)
        if len(left) is 0 then
            right
        else if len(right) is 0 then
            left
        else if left[0] < right[0] then
            prepend(left[0], merge(tail(left), right))
        else
            prepend(right[0], merge(left, tail(right))),

    // takes a list of integers and returns a sorted list containing the same integers
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

    // takes a list of strings, parses them as integers, and returns a list of the parsed integers
    asNums = func(strings)
        parseInt(string) for string in strings,

    // read user input and parse as a list of integers
    nums = asNums(split(input("Enter a list of numbers separated by spaces: ")))
in
    // sort the list of integers and print the result
    sort(nums)