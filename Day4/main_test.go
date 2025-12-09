package main

import (
	"strings"
	"testing"
)

const sampleInput = "..@@.@@@@.\n@@@.@.@.@@\n@@@@@.@.@@\n@.@@@@..@.\n@@.@@@@.@@\n.@@@@@@@.@\n.@.@.@.@@@\n@.@@@.@@@@\n.@@@@@@@@.\n@.@.@@@.@.\n"

func TestSolveSample(t *testing.T) {
	part1, part2, err := Solve(strings.NewReader(sampleInput))
	if err != nil {
		t.Fatalf("Solve() error = %v", err)
	}

	if part1 != 13 {
		t.Fatalf("part1 = %d, want 13", part1)
	}
	if part2 != 43 {
		t.Fatalf("part2 = %d, want 43", part2)
	}
}
