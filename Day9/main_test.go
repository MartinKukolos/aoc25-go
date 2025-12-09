package main

import (
	"strings"
	"testing"
)

const sampleInput = "7,1\n11,1\n11,7\n9,7\n9,5\n2,5\n2,3\n7,3\n"

func TestSolveSample(t *testing.T) {
	part1, part2, err := Solve(strings.NewReader(sampleInput))
	if err != nil {
		t.Fatalf("Solve() error = %v", err)
	}

	var wantPart1 int64 = 50
	var wantPart2 int64 = 24

	if part1 != wantPart1 {
		t.Fatalf("part1 = %d, want %d", part1, wantPart1)
	}
	if part2 != wantPart2 {
		t.Fatalf("part2 = %d, want %d", part2, wantPart2)
	}
}
