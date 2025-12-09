package main

import (
	"strings"
	"testing"
)

func TestSolveSample(t *testing.T) {
	const input = "11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124"

	got, err := Solve(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Solve() error = %v", err)
	}

	const want int64 = 1227775554
	if got != want {
		t.Fatalf("Solve() = %d, want %d", got, want)
	}
}
