package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func expectNoError(err error, msg string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, msg)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type AttackType uint8

const (
	Rock AttackType = iota
	Paper
	Scissor
)

func (a AttackType) String() string {
	switch a {
	case Rock:
		return "Rock"
	case Paper:
		return "Paper"
	case Scissor:
		return "Scissor"
	default:
		return "Undefined"
	}
}

func (a AttackType) Points() uint64 {
	switch a {
	case Rock:
		return 1
	case Paper:
		return 2
	case Scissor:
		return 3
	default:
		return 0
	}
}

func (a AttackType) Battle(b AttackType) Outcome {
	switch a {
	case Rock:
		switch b {
		case Paper:
			return Defeat
		case Scissor:
			return Victory
		default:
			return Draw
		}
	case Paper:
		switch b {
		case Rock:
			return Victory
		case Scissor:
			return Defeat
		default:
			return Draw
		}
	case Scissor:
		switch b {
		case Rock:
			return Defeat
		case Paper:
			return Victory
		default:
			return Draw
		}
	default:
		return Draw
	}
}

func (a AttackType) GuessAttackType(outcome Outcome) AttackType {
	switch a {
	case Rock:
		switch outcome {
		case Victory:
			return Paper
		case Defeat:
			return Scissor
		default:
			return a
		}
	case Paper:
		switch outcome {
		case Victory:
			return Scissor
		case Defeat:
			return Rock
		default:
			return a
		}
	case Scissor:
		switch outcome {
		case Victory:
			return Rock
		case Defeat:
			return Paper
		default:
			return a
		}
	default:
		return a
	}
}

type Outcome uint8

const (
	Victory Outcome = iota
	Defeat
	Draw
)

func (a Outcome) String() string {
	switch a {
	case Victory:
		return "Victory"
	case Defeat:
		return "Defeat"
	case Draw:
		return "Draw"
	default:
		return "Undefined"
	}
}

func (a Outcome) Points() uint64 {
	switch a {
	case Victory:
		return 6
	case Defeat:
		return 0
	case Draw:
		return 3
	default:
		return 0
	}
}

type MatchDetail struct {
	EnemyAttack      AttackType
	ExpectedResponse AttackType
	Outcome          Outcome
}

func (m MatchDetail) ComputeScore() uint64 {
	return m.Outcome.Points() + m.ExpectedResponse.Points()
}

type MatchDetailSlice []MatchDetail

func parseAttackType(baseChar, rawType byte) (AttackType, error) {
	switch rawType {
	case baseChar + 0:
		return Rock, nil
	case baseChar + 1:
		return Paper, nil
	case baseChar + 2:
		return Scissor, nil
	}

	return Scissor, errors.New("Invalid AttackType!")
}

func parseOutcome(rawType byte) (Outcome, error) {
	switch rawType {
	case 'X':
		return Defeat, nil
	case 'Y':
		return Draw, nil
	case 'Z':
		return Victory, nil
	}

	return Draw, errors.New("Invalid Outcome!")
}

func parseMatchDetailList(scanner *bufio.Scanner, secondArgumentIsOutcome bool) (MatchDetailSlice, error) {
	data := MatchDetailSlice{}

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")

		if len(parts) != 2 || len(parts[0]) != 1 || len(parts[1]) != 1 {
			return nil, errors.New("Invalid line found!")
		}

		enemyAttack, err := parseAttackType('A', parts[0][0])

		if err != nil {
			return nil, err
		}

		var outcome Outcome
		var expectedResponse AttackType

		if secondArgumentIsOutcome {
			outcome, err = parseOutcome(parts[1][0])
			if err != nil {
				return nil, err
			}

			expectedResponse = enemyAttack.GuessAttackType(outcome)
		} else {
			expectedResponse, err = parseAttackType('X', parts[1][0])
			if err != nil {
				return nil, err
			}
			outcome = expectedResponse.Battle(enemyAttack)
		}

		data = append(data, MatchDetail{EnemyAttack: enemyAttack, ExpectedResponse: expectedResponse, Outcome: outcome})
	}

	return data, nil
}

func partGeneric(scanner *bufio.Scanner, isPart2 bool) {
	data, err := parseMatchDetailList(scanner, isPart2)
	expectNoError(err, "Parsing error")

	var totalScore uint64

	for _, match := range data {
		totalScore += match.ComputeScore()
	}

	fmt.Println(totalScore)
}

func main() {
	var err error

	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage ", os.Args[0], "<1|2> <file>")
		os.Exit(1)
	}

	partNum, err := strconv.ParseUint(os.Args[1], 10, 8)
	expectNoError(err, "Part Number must be an integer")

	file, err := os.Open(os.Args[2])
	defer file.Close()
	expectNoError(err, "Cannot open file")

	var data = bufio.NewScanner(file)

	switch partNum {
	case 1:
		partGeneric(data, false)
		break
	case 2:
		partGeneric(data, true)
		break
	default:
		fmt.Fprintln(os.Stderr, "Part number must be between 1 and 2")
		os.Exit(1)
		break
	}

	expectNoError(data.Err(), "Error during scanner read")
}
