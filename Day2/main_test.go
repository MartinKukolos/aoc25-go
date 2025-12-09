package main

import (
	"strings"
	"testing"
)

func TestSolveSample(t *testing.T) {
	const input = "11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124"

	part1, part2, err := Solve(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Solve() error = %v", err)
	}

	const wantPart1 int64 = 1227775554
	const wantPart2 int64 = 4174379265
	if part1 != wantPart1 {
		t.Fatalf("Solve() part1 = %d, want %d", part1, wantPart1)
	}
	if part2 != wantPart2 {
		t.Fatalf("Solve() part2 = %d, want %d", part2, wantPart2)
	}
}
