package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func main() {
	answer, err := solveDistressSignal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to solve distress signal: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(answer)
}

func solveDistressSignal() (int, error) {
	pairs, err := parsePairs()
	if err != nil {
		return 0, fmt.Errorf("parse pairs: %v", err)
	}

	indexSum := 0
	for pi, p := range pairs {
		switch isListPairInTheRightOrder(p[0], p[1]) {
		case comparisonResultRightOrder, comparisonResultContinue:
			indexSum += pi + 1
		}
	}

	return indexSum, nil
}

type comparisonResult int

const (
	comparisonResultRightOrder = iota - 1
	comparisonResultContinue
	comparisonResultWrongOrder
)

func isPairInTheRightOrder(l, r any) comparisonResult {
	switch typedL := l.(type) {
	case float64:
		switch typedR := r.(type) {
		case float64:
			return isFloatPairInTheRightOrder(typedL, typedR)
		case []any:
			return isListPairInTheRightOrder([]any{l}, typedR)
		default:
			panic(fmt.Sprintf("unknown type %T", r))
		}
	case []any:
		switch typedR := r.(type) {
		case float64:
			return isListPairInTheRightOrder(typedL, []any{r})
		case []any:
			return isListPairInTheRightOrder(typedL, typedR)
		default:
			panic(fmt.Sprintf("unknown type %T", l))
		}
	default:
		panic(fmt.Sprintf("unknown type %T", l))
	}
}

func isListPairInTheRightOrder(l, r []any) comparisonResult {
	for i, lv := range l {
		if i == len(r) {
			return comparisonResultWrongOrder
		}

		switch res := isPairInTheRightOrder(lv, r[i]); res {
		case comparisonResultRightOrder, comparisonResultWrongOrder:
			return res
		}
	}
	if len(l) < len(r) {
		return comparisonResultRightOrder
	}
	return comparisonResultContinue
}

func isFloatPairInTheRightOrder(l, r float64) comparisonResult {
	switch {
	case l < r:
		return comparisonResultRightOrder
	case l == r:
		return comparisonResultContinue
	default:
		return comparisonResultWrongOrder
	}
}

func parsePairs() ([][][]any, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var (
		pairs    [][][]any
		currPair = make([][]any, 0, 2)
	)
	for line := 1; sc.Scan(); line++ {
		b := sc.Bytes()
		if line%3 == 0 {
			if len(b) != 0 {
				return nil, fmt.Errorf("expected empty line after pair on line %d", line)
			}
			pairs = append(pairs, currPair)
			currPair = make([][]any, 0, 2)
			continue
		}

		var list []any
		if err = json.Unmarshal(b, &list); err != nil {
			return nil, fmt.Errorf("unmarshal line %d", line)
		}
		currPair = append(currPair, list)
	}
	if err = sc.Err(); err != nil {
		return nil, fmt.Errorf("scan: %v", err)
	}

	if len(currPair) != 0 {
		if len(currPair) != 2 {
			return nil, errors.New("last pair is uncompleted")
		}
		pairs = append(pairs, currPair)
	}
	for i, p := range pairs {
		for j, l := range p {
			if err = validateList(l); err != nil {
				return nil, fmt.Errorf("validate pair %d, list %d: %v", i, j, err)
			}
		}
	}

	return pairs, nil
}

func validateList(l []any) error {
	for i, el := range l {
		switch el := el.(type) {
		case float64:
		case []any:
			if err := validateList(el); err != nil {
				return fmt.Errorf("list at index %d", i)
			}
		default:
			return fmt.Errorf("list element at index %d has wrong type %T", i, el)
		}
	}
	return nil
}
