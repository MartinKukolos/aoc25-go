package main

import (
	"strings"
	"testing"
)

const sampleInput = "3-5\n10-14\n16-20\n12-18\n\n1\n5\n8\n11\n17\n32\n"

func TestSolveSample(t *testing.T) {
	part1, part2, err := Solve(strings.NewReader(sampleInput))
	if err != nil {
		t.Fatalf("Solve() error = %v", err)
	}
	if part1 != 3 {
		t.Fatalf("part1 = %d, want 3", part1)
	}
	if part2 != 14 {
		t.Fatalf("part2 = %d, want 14", part2)
	}
}
