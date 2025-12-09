package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	dialSize      = 100
	startPosition = 50
)

func main() {
	path := resolveInputPath(os.Args)

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open input %q: %v\n", path, err)
		os.Exit(1)
	}
	defer file.Close()

	answer, err := Solve(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "solve error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(answer)
}

func resolveInputPath(args []string) string {
	if len(args) > 1 {
		return args[1]
	}
	if _, err := os.Stat("Day1/input.txt"); err == nil {
		return "Day1/input.txt"
	}
	return "input.txt"
}

func Solve(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	position := startPosition
	zeroHits := 0
	lineNumber := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lineNumber++

		if len(line) < 2 {
			return 0, fmt.Errorf("line %d: rotation too short", lineNumber)
		}

		dir := line[0]
		steps, err := strconv.Atoi(line[1:])
		if err != nil {
			return 0, fmt.Errorf("line %d: invalid distance: %w", lineNumber, err)
		}

		switch dir {
		case 'L':
			position = mod(position - steps)
		case 'R':
			position = mod(position + steps)
		default:
			return 0, fmt.Errorf("line %d: invalid direction %q", lineNumber, dir)
		}

		if position == 0 {
			zeroHits++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return zeroHits, nil
}

func mod(value int) int {
	value %= dialSize
	if value < 0 {
		value += dialSize
	}
	return value
}
