let
// tree utils
    emptyTree = [false, false, [], []],

    treeWithVal = func(tree, val) [true, val, at(tree, 2), at(tree, 3)],
            
    treeWithLeft = func(tree, left) [at(tree, 0), at(tree, 1), left, at(tree, 3)],

    treeWithRight = func(tree, right) [at(tree, 0), at(tree, 1), at(tree, 2), right],

    treeHasVal = func(tree) at(tree, 0),

    treeVal = func(tree) at(tree, 1),

    treeHasLeft = func(tree) greaterThan(len(at(tree, 2)), 0),

    treeLeft = func(tree) at(tree, 2),

    treeHasRight = func(tree) greaterThan(len(at(tree, 3)), 0),

    treeRight = func(tree) at(tree, 3),

// binary search tree functions
    bstPush = func(tree, val)
        if not(treeHasVal(tree)) then
            treeWithVal(tree, val)
        else if lessThan(val, treeVal(tree)) then
            if treeHasLeft(tree) then
                treeWithLeft(tree, bstPush(treeLeft(tree), val))
            else
                treeWithLeft(tree, treeWithVal(emptyTree, val))
        else
            if treeHasRight(tree) then
                treeWithRight(tree, bstPush(treeRight(tree), val))
            else
                treeWithRight(tree, treeWithVal(emptyTree, val)),
    
    bstFind = func(tree, val)
        if not(treeHasVal(tree)) then
            false
        else if equals(val, treeVal(tree)) then
            true
        else if lessThan(val, treeVal(tree)) then
            if not(treeHasLeft(tree)) then
                false
            else
                bstFind(treeLeft(tree), val)
        else
            if not(treeHasRight(tree)) then
                false
            else
                bstFind(treeRight(tree), val),
    
// sample tree
    tree = bstPush(bstPush(bstPush(bstPush(emptyTree, 5), 3), 7), 1)
in
    [i, bstFind(tree, i)] for i in range(0, 10)