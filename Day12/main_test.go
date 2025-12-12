package main

import (
	"os"
	"testing"
)

func TestSolveSample(t *testing.T) {
	f, err := os.Open("sample.txt")
	if err != nil {
		t.Fatalf("open sample: %v", err)
	}
	defer f.Close()

	p1, p2, err := Solve(f)
	if err != nil {
		t.Fatalf("Solve(sample) error = %v", err)
	}
	if p1 != 2 {
		t.Fatalf("part1 = %d, want 2", p1)
	}
	if p2 != 0 {
		t.Fatalf("part2 = %d, want 0", p2)
	}
}
