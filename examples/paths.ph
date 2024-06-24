let
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

    // general utils
    concatAll = func(l)
        if len(l) is 0 then
            ""
        else
            concatStr(l[0], concatAll(tail(l))),
    
    zip = func(l1, l2)
        if len(l1) is 0 or len(l2) is 0 then
            []
        else
            concat([l1[0], l2[0]], zip(tail(l1), tail(l2))),

    // tile types
    TILE_TYPE_EMPTY = 0,
    TILE_TYPE_WALL = 1,
    TILE_TYPE_START = 2,
    TILE_TYPE_END = 3,
    TILE_TYPE_PATH = 4,

    // directions
    DIR_LEFT = 0,
    DIR_UP = 1,
    DIR_RIGHT = 2,
    DIR_DOWN = 3,

    // maze utils
    mazeTileString = func(tile)
        if tile is TILE_TYPE_EMPTY then
            " "
        else if tile is TILE_TYPE_WALL then
            "#"
        else if tile is TILE_TYPE_START then
            "S"
        else if tile is TILE_TYPE_END then
            "E"
        else if tile is TILE_TYPE_PATH then
            "."
        else
            "?",

    dirString = func(dir)
        if dir is DIR_LEFT then
            "L"
        else if dir is DIR_UP then
            "U"
        else if dir is DIR_RIGHT then
            "R"
        else if dir is DIR_DOWN then
            "D"
        else
            "?",

    mazeRowString = func(mazeRow)
        concatAll(mazeTileString(tile) for tile in mazeRow),

    mazeString = func(maze)
        concatAll(concatStr(mazeRowString(mazeRow), "\n") for mazeRow in maze),
    
    mazeRowWithCoordAs = func(mazeRow, x, tile)
        if len(mazeRow) is 0 then
            mazeRow
        else
            let
                curTile = mazeRow[0],
                mazeRowRest = tail(mazeRow)
            in
                if x is 0 then
                    prepend(tile, mazeRowRest)
                else
                    prepend(curTile, mazeRowWithCoordAs(mazeRowRest, x-1, tile)),

    mazeWithCoordsAs = func(maze, x, y, tile)
        if len(maze) is 0 then
            maze
        else
            let
                curRow = maze[0],
                mazeRest = tail(maze)
            in
                if y is 0 then
                    prepend(mazeRowWithCoordAs(curRow, x, tile), mazeRest)
                else
                    prepend(curRow, mazeWithCoordsAs(mazeRest, x, y-1, tile)),

    // maze
    defaultMaze = [
        [0, 1, 1, 0, 1, 1, 1, 1, 0, 0],
        [2, 0, 0, 0, 0, 0, 0, 0, 0, 1],
        [0, 1, 1, 1, 0, 1, 1, 1, 0, 1],
        [1, 0, 0, 0, 0, 1, 0, 0, 0, 0],
        [0, 1, 0, 1, 1, 1, 1, 1, 1, 0],
        [0, 0, 0, 0, 0, 0, 0, 0, 1, 0],
        [1, 0, 1, 0, 1, 1, 1, 1, 3, 1],
        [0, 0, 1, 0, 0, 0, 0, 0, 0, 1]
    ],

    findStart = func(maze)
        let
            startIdxForEachRow = find(func(tile) tile is TILE_TYPE_START, row) for row in maze,
            startRow = find(func(startIdxResult) startIdxResult is not false, startIdxForEachRow),
            startCol = startIdxForEachRow[startRow]
        in {
            x: startCol,
            y: startRow
        },
        

    solveMazeHelper = func(maze, x, y, pathSoFar)
        if x < 0 or y < 0 or x >= len(maze[0]) or y >= len(maze) then
            false
        else
            let
                tile = maze[y][x]
            in
                if tile is TILE_TYPE_WALL or tile is TILE_TYPE_PATH then
                    false
                else if tile is TILE_TYPE_END then
                    pathSoFar
                else
                    let
                        maze = mazeWithCoordsAs(maze, x, y, TILE_TYPE_PATH),
                        unused = print(concatStr("\n", mazeString(maze))),
                        pathLeft = solveMazeHelper(maze, x - 1, y, append(pathSoFar, DIR_LEFT)),
                        pathUp = solveMazeHelper(maze, x, y - 1, append(pathSoFar, DIR_UP)),
                        pathRight = solveMazeHelper(maze, x + 1, y, append(pathSoFar, DIR_RIGHT)),
                        pathDown = solveMazeHelper(maze, x, y + 1, append(pathSoFar, DIR_DOWN)),
                        paths = [pathLeft, pathUp, pathRight, pathDown],
                        goodPaths = filter(paths, func(path) path is not false),
                        minPathLenAndIdx = min(len(path) for path in goodPaths),
                        res = if minPathLenAndIdx is false then
                                false
                            else
                                goodPaths[minPathLenAndIdx.idx]
                    in
                        res,
    
    solveMaze = func(maze)
        let
            startCoords = findStart(maze)
        in
            solveMazeHelper(maze, startCoords.x, startCoords.y, []),

    unused = print(mazeString(defaultMaze)),
    
    res = solveMaze(defaultMaze)
in
    if res is false then
        false
    else
        dirString(dir) for dir in res