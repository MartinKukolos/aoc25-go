package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	maxDigits = 18
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
	if _, err := os.Stat("Day2/input.txt"); err == nil {
		return "Day2/input.txt"
	}
	return "input.txt"
}

func Solve(r io.Reader) (int64, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}

	text := strings.TrimSpace(string(data))
	if text == "" {
		return 0, fmt.Errorf("input is empty")
	}

	ranges, err := parseRanges(text)
	if err != nil {
		return 0, err
	}

	var total int64
	for _, rg := range ranges {
		if rg.start > rg.end {
			return 0, fmt.Errorf("range start %d greater than end %d", rg.start, rg.end)
		}
		total += sumInvalidInRange(rg)
	}

	return total, nil
}

type idRange struct {
	start int64
	end   int64
}

func parseRanges(text string) ([]idRange, error) {
	parts := strings.Split(text, ",")
	ranges := make([]idRange, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		dash := strings.IndexByte(part, '-')
		if dash == -1 {
			return nil, fmt.Errorf("missing '-' in %q", part)
		}

		start, err := strconv.ParseInt(part[:dash], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid start %q: %w", part[:dash], err)
		}

		end, err := strconv.ParseInt(part[dash+1:], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid end %q: %w", part[dash+1:], err)
		}

		ranges = append(ranges, idRange{start: start, end: end})
	}

	if len(ranges) == 0 {
		return nil, fmt.Errorf("no ranges found")
	}

	return ranges, nil
}

var pow10 = func() []int64 {
	vals := make([]int64, maxDigits+1)
	vals[0] = 1
	for i := 1; i <= maxDigits; i++ {
		vals[i] = vals[i-1] * 10
	}
	return vals
}()

func sumInvalidInRange(rg idRange) int64 {
	var sum int64
	for k := 1; k <= maxDigits/2; k++ {
		baseMin := pow10[k-1]
		multiplier := pow10[k] + 1

		if baseMin*multiplier > rg.end {
			break
		}

		baseMax := pow10[k] - 1

		loBase := baseMin
		if candidate := ceilDiv(rg.start, multiplier); candidate > loBase {
			loBase = candidate
		}

		hiBase := baseMax
		if candidate := rg.end / multiplier; candidate < hiBase {
			hiBase = candidate
		}

		if loBase > hiBase {
			continue
		}

		count := hiBase - loBase + 1
		sumBases := (loBase + hiBase) * count / 2
		sum += sumBases * multiplier
	}

	return sum
}

func ceilDiv(num, denom int64) int64 {
	if denom <= 0 {
		panic("denominator must be positive")
	}
	if num >= 0 {
		return (num + denom - 1) / denom
	}
	return num / denom
}
