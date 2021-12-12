package main

type cell struct {
	number int
	marked bool
}

func playGame(boards [][]cell, rowNum, colNum int, drawNumbers []int) (winBoardIdx int, winDrawNumberIdx int, ok bool) {
	for dI, dn := range drawNumbers {
		for bI, b := range boards {
			foundRowIdx, foundColIdx := 0, 0
			found := false
			for i, c := range b {
				if c.number == dn {
					b[i].marked = true
					found = true
					foundRowIdx = i / rowNum
					foundColIdx = i % rowNum
					break // all elements are unique.
				}
			}

			if !found {
				continue
			}

			rowWin := true
			for i := 0; i < colNum; i++ {
				if !b[foundRowIdx*colNum+i].marked {
					rowWin = false
					break
				}
			}
			if rowWin {
				return bI, dI, true
			}

			colWin := true
			for i := 0; i < rowNum; i++ {
				if !b[i*colNum+foundColIdx].marked {
					colWin = false
					break
				}
			}
			if colWin {
				return bI, dI, true
			}
		}
	}

	return 0, 0, false
}

func playSquidGame(boards [][]cell, rowNum, colNum int, drawNumbers []int) (unmarkedSum int, lastWinDrawNumberIdx int, ok bool) {
	var (
		lastFound    = false
		lastWinBoard []cell
	)
	for {
		winBoardIdx, winDrawNumberIdx, gameOK := playGame(boards, rowNum, colNum, drawNumbers[lastWinDrawNumberIdx:])
		if !gameOK {
			break
		}
		lastFound = true
		lastWinBoard = boards[winBoardIdx]
		lastWinDrawNumberIdx += winDrawNumberIdx

		boards = append(boards[:winBoardIdx], boards[winBoardIdx+1:]...)
	}

	for _, c := range lastWinBoard {
		if !c.marked {
			unmarkedSum += c.number
		}
	}

	return unmarkedSum, lastWinDrawNumberIdx, lastFound
}
