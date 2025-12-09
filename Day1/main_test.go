package main

import (
	"strings"
	"testing"
)

func TestSolveSample(t *testing.T) {
	const input = `L68
L30
R48
L5
R60
L55
L1
L99
R14
L82
`

	got, err := Solve(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Solve() error = %v", err)
	}

	const want = 3
	if got != want {
		t.Fatalf("Solve() = %d, want %d", got, want)
	}
}
