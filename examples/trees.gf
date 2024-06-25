let
// misc helpers
    tail = func(l) if len(l) <= 1 then [] else l[1:],

    fold = func(t, f, l)
        if len(l) is 0 then
            t
        else
            let
                this = l[0],
                rest = tail(l),
                nextT = f(t, this)
            in
                fold(nextT, f, tail(l)),

    splitHelper = func(string)
        if lenStr(string) is 0 then {
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
                if this is " " then {
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
    
    // takes a list of strings, parses them as integers, and returns a list of the parsed integers
    asNums = func(strings)
        parseInt(string) for string in strings,
    
    forever = func(f)
        [f(), forever(f)],

// tree utils
    emptyTree = {
        val: false,
        hasVal: false,
        left: false,
        right: false
    },

    treeWithVal = func(tree, val) {
        val: val,
        hasVal: true,
        left: tree.left,
        right: tree.right
    },
            
    treeWithLeft = func(tree, left) {
        val: tree.val,
        hasVal: tree.hasVal,
        left: left,
        right: tree.right
    },

    treeWithRight = func(tree, right) {
        val: tree.val,
        hasVal: tree.hasVal,
        left: tree.left,
        right: right
    },

    treeHasLeft = func(tree) tree.left is not false,

    treeHasRight = func(tree) tree.right is not false,

// binary search tree functions
    bstPush = func(tree, val)
        if not tree.hasVal then
            treeWithVal(tree, val)
        else if val < tree.val then
            if treeHasLeft(tree) then
                treeWithLeft(tree, bstPush(tree.left, val))
            else
                treeWithLeft(tree, treeWithVal(emptyTree, val))
        else
            if treeHasRight(tree) then
                treeWithRight(tree, bstPush(tree.right, val))
            else
                treeWithRight(tree, treeWithVal(emptyTree, val)),
    
    bstFind = func(tree, val)
        if not tree.hasVal then
            false
        else if val is tree.val then
            true
        else if val < tree.val then
//            if not treeHasLeft(tree) then
//                false
//            else
                bstFind(tree.left, val)
        else
//            if not treeHasRight(tree) then
//                false
//            else
                bstFind(tree.right, val),

    nums = asNums(split(input("Enter a list of numbers: "))),

    tree = fold(emptyTree, bstPush, nums)
in
    forever(
        func()
            let
                numToFind = parseInt(input("Enter a number to find: "))
            in
                print(bstFind(tree, numToFind))
    )