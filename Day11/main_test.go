package main

import (
	"os"
	"testing"
)

func TestSolveSample1(t *testing.T) {
	f, err := os.Open("sample.txt")
	if err != nil {
		t.Fatalf("open sample: %v", err)
	}
	defer f.Close()

	part1, part2, err := Solve(f)
	if err != nil {
		t.Fatalf("Solve(sample) error = %v", err)
	}
	if part1 != 5 {
		t.Fatalf("part1 = %d, want 5", part1)
	}
	if part2 != 0 {
		t.Fatalf("part2 = %d, want 0", part2)
	}
}

func TestSolveSample2(t *testing.T) {
	f, err := os.Open("sample2.txt")
	if err != nil {
		t.Fatalf("open sample2: %v", err)
	}
	defer f.Close()

	part1, part2, err := Solve(f)
	if err != nil {
		t.Fatalf("Solve(sample2) error = %v", err)
	}
	if part1 != 0 {
		t.Fatalf("part1 = %d, want 0", part1)
	}
	if part2 != 2 {
		t.Fatalf("part2 = %d, want 2", part2)
	}
}
