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
    
    zip = func(l1, l2)
        if len(l1) is 0 or len(l2) is 0 then
            []
        else
            concat([l1[0], l2[0]], zip(tail(l1), tail(l2))),
    
    withIdxAs = func(l, i, v)
        if i >= len(l) then l else
            concat(append(l[:i], v), l[i+1:]),

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

    mazeRowString = func(mazeRow, visitedRow, curX)
        let
            isVisited = func(x) visitedRow[x] is not false
        in
            concatAll((
                let
                    tileStr = if x is curX then
                        "-"
                    else if mazeRow[x] is TILE_TYPE_EMPTY and isVisited(x) then
                        itoa(len(visitedRow[x]))
                    else
                        mazeTileString(mazeRow[x])
                in
                    concatStr(tileStr, concatAll(
                        " " for _ in range(0, 4 - lenStr(tileStr))
                    ))
            ) for x in range(0, len(mazeRow))),

    mazeString = func(maze, visited, curX, curY)
        concatAll(
            concatStr(mazeRowString(maze[i], visited[i], if i is curY then curX else false), "\n\n") for i in range(0, len(maze))
        ),
    
    mazeRowWithCoordAs = func(mazeRow, x, tile)
        withIdxAs(mazeRow, x, tile),

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
    
    findStart = func(maze)
        let
            // in each row, try to find the start tile
            startIdxForEachRow = find(func(tile) tile is TILE_TYPE_START, row) for row in maze,
            // find the row with non-false start index; i.e. the row for which the start location was found
            startRow = find(func(startIdxResult) startIdxResult is not false, startIdxForEachRow),
            startCol = startIdxForEachRow[startRow]
        in {
            x: startCol,
            y: startRow
        },


    findEnd = func(maze)
        let
            // in each row, try to find the start tile
            endIdxForEachRow = find(func(tile) tile is TILE_TYPE_END, row) for row in maze,
            // find the row with non-false start index; i.e. the row for which the start location was found
            endRow = find(func(endIdxResult) endIdxResult is not false, endIdxForEachRow),
            endCol = endIdxForEachRow[endRow]
        in {
            x: endCol,
            y: endRow
        },

    // maze should be a 2D array of tile types
    // visited should be a 2D array of the same dimensions as maze, containing the best known paths
    //      from start to these coordinates, or false if no paths are yet known.
    // x is the x coordinate from which we are solving the maze
    // y is the y coordinate from which we are solving the maze
    // pathSoFar is a list of directions representing the path we took from start to get here
    // 
    // returns a 2D array that is the same as visited, but any coordinates
    // that we found a better path to are replaced by the
    // better paths
    solveMazeHelper = func(maze, visited, queue)
        if
            len(queue) is 0
        then
            visited
        else let
            coordsAndPath = queue[0],
            pathSoFar = coordsAndPath.path,
            coords = coordsAndPath.coords,
            x = coords.x,
            y = coords.y,
            queue = tail(queue)
        in if x < 0 or
                y < 0 or
                x >= len(maze[0]) or
                y >= len(maze)
        then
            solveMazeHelper(maze, visited, queue)
        else let
            tile = maze[y][x],
            isVisited = visited[y][x] is not false,
            bestPathFromStartSoFar = visited[y][x]
        in if tile is TILE_TYPE_WALL or (
                isVisited and
                len(bestPathFromStartSoFar) <= len(pathSoFar)
        ) then
            solveMazeHelper(maze, visited, queue)
        else let
                visited = withIdxAs(visited, y, withIdxAs(visited[y], x, pathSoFar)),
                _ = print(concatStr("\n", mazeString(maze, visited, x, y))),
                // _ = input("press enter to continue")
                _ = sleep(100)
            in
                if tile is TILE_TYPE_END then
                    solveMazeHelper(maze, visited, queue)
                else 
                    let
                        queue = append(queue, {
                            coords: {x: x-1, y: y},
                            path: append(pathSoFar, DIR_LEFT)
                        }),
                        queue = append(queue, {
                            coords: {x: x, y: y-1},
                            path: append(pathSoFar, DIR_UP)
                        }),
                        queue = append(queue, {
                            coords: {x: x+1, y: y},
                            path: append(pathSoFar, DIR_RIGHT)   
                        }),
                        queue = append(queue, {
                            coords: {x: x, y: y+1},
                            path: append(pathSoFar, DIR_DOWN)
                        })
                    in
                        solveMazeHelper(maze, visited, queue),
    
    noneVisited = func(maze) (false for _ in mazeRow) for mazeRow in maze,

    drawFinishedMazeHelper = func(maze, coords, path)
        let
            dottedMaze = withIdxAs(maze, coords.y, withIdxAs(maze[coords.y], coords.x, "."))
        in
            if len(path) is 0 then
                dottedMaze
            else
                let
                    newCoords = if path[0] is DIR_LEFT then {
                        x: coords.x - 1,
                        y: coords.y
                    } else if path[0] is DIR_UP then {
                        x: coords.x,
                        y: coords.y - 1
                    } else if path[0] is DIR_RIGHT then {
                        x: coords.x + 1,
                        y: coords.y
                    } else if path[0] is DIR_DOWN then {
                        x: coords.x,
                        y: coords.y + 1
                    } else false
                in
                    if newCoords is false then [] else drawFinishedMazeHelper(dottedMaze, newCoords, tail(path)),

    drawFinishedMaze = func(maze, path)
        let rawMaze = (mazeTileString(tile) for tile in row) for row in maze,
            coords = findStart(maze)
        in
            concatAll(
                concatAll(
                    append(rowStrs, "\n")
                ) for rowStrs in drawFinishedMazeHelper(rawMaze, coords, path)),

    solveMaze = func(maze)
        let
            startCoords = findStart(maze),
            endCoords = findEnd(maze),
            bestPaths = solveMazeHelper(maze, noneVisited(maze), [{
                coords: startCoords,
                path: []
            }])
        in
            bestPaths[endCoords.y][endCoords.x],

    // maze
    defaultMaze = [
        [2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
        [0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0],
        [0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0],
        [0, 1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1],
        [0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0],
        [0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0],
        [0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1],
        [0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0],
        [0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0],
        [0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 3]
    //    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    //    [0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1],
    //    [0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0],
    //    [0, 0, 2, 1, 0, 1, 0, 1, 0, 0, 0, 1],
    //    [0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0],
    //    [0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 1, 0],
    //    [0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0],
    //    [1, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 1],
    //    [0, 0, 0, 0, 1, 0, 1, 0, 1, 3, 0, 0],
    //    [1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0, 1],
    //    [0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0]
    ],
    
    res = solveMaze(defaultMaze)
in
    if res is false then
        false
    else
        concatStr("\n", drawFinishedMaze(defaultMaze, res))
