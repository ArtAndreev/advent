package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func readDrawNumbers(sc *bufio.Scanner) ([]int, error) {
	var res []int
	if sc.Scan() {
		rawNums := strings.Split(sc.Text(), ",")
		res = make([]int, 0, len(rawNums))
		for i, rn := range rawNums {
			n, err := strconv.Atoi(rn)
			if err != nil {
				return nil, fmt.Errorf("draw number in pos %d is invalid: %v", i, err)
			}
			res = append(res, n)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read first line: %v", err)
	}

	if sc.Scan() && sc.Text() != "" {
		return nil, errors.New("expected empty second line")
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read second line: %v", err)
	}

	return res, nil
}

func readBoards(sc *bufio.Scanner) (boards [][]cell, rowNum, colNum int, err error) {
	var board []cell
	for i := 2; sc.Scan(); i++ {
		t := sc.Text()
		if t == "" {
			if len(board) != 0 {
				currRowNum := len(board) / colNum
				if rowNum == 0 {
					rowNum = currRowNum
				}
				if currRowNum != rowNum {
					return nil, 0, 0, fmt.Errorf("all row nums in boards must be equal to %d, got %d, line %d", rowNum, currRowNum, i)
				}
				boards = append(boards, board)
				board = make([]cell, 0, rowNum*colNum)
			}
			continue
		}
		rawNums := strings.Fields(t)
		currColNum := len(rawNums)
		if colNum == 0 {
			colNum = currColNum
		}
		if currColNum != colNum {
			return nil, 0, 0, fmt.Errorf("all column nums in boards must be equal to %d, got %d, line %d", colNum, currColNum, i)
		}
		for j, rn := range rawNums {
			n, err := strconv.Atoi(rn)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("number on board in line %d, pos %d is invalid: %v", i, j, err)
			}
			board = append(board, cell{number: n})
		}
	}
	if err = sc.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("failed to scan: %s", err)
	}
	if len(board) != 0 {
		currRowNum := len(board) / colNum
		if rowNum == 0 {
			rowNum = currRowNum
		}
		if currRowNum != rowNum {
			return nil, 0, 0, fmt.Errorf("all row nums in boards must be equal to %d, got %d, last board", rowNum, currRowNum)
		}
		boards = append(boards, board)
		board = make([]cell, 0, rowNum*colNum)
	}

	return boards, rowNum, colNum, nil
}

func assertNoDuplicateNumbers(numbers []int) {
	unique := make(map[int]struct{}, len(numbers))
	for _, n := range numbers {
		if _, ok := unique[n]; ok {
			log.Panicf("duplicate number %d", n)
		}
		unique[n] = struct{}{}
	}
}

func assertNoDuplicateCellNumbers(boards [][]cell) {
	for i, b := range boards {
		unique := make(map[int]struct{}, len(b))
		for _, c := range b {
			n := c.number
			if _, ok := unique[n]; ok {
				log.Panicf("duplicate cell number %d, board %d", n, i)
			}
			unique[n] = struct{}{}
		}
	}
}
