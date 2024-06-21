let
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
            if not treeHasLeft(tree) then
                false
            else
                bstFind(tree.left, val)
        else
            if not treeHasRight(tree) then
                false
            else
                bstFind(tree.right, val),
    
// sample tree
    tree = bstPush(bstPush(bstPush(bstPush(emptyTree, 5), 3), 7), 1)
in
    [i, bstFind(tree, i)] for i in range(0, 10)