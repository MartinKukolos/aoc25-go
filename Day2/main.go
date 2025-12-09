package main

import (
	"fmt"
	"io"
	"os"
	"sort"
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

	part1, part2, err := Solve(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "solve error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
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

func Solve(r io.Reader) (int64, int64, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, 0, err
	}

	text := strings.TrimSpace(string(data))
	if text == "" {
		return 0, 0, fmt.Errorf("input is empty")
	}

	ranges, err := parseRanges(text)
	if err != nil {
		return 0, 0, err
	}

	var totalTwice int64
	var totalAny int64
	for _, rg := range ranges {
		if rg.start > rg.end {
			return 0, 0, fmt.Errorf("range start %d greater than end %d", rg.start, rg.end)
		}
		totalTwice += sumInvalidRepeatedTwice(rg)
		totalAny += sumInvalidAnyRepeat(rg)
	}

	return totalTwice, totalAny, nil
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

func sumInvalidRepeatedTwice(rg idRange) int64 {
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

func sumInvalidAnyRepeat(rg idRange) int64 {
	var sum int64
	for length := 2; length <= maxDigits && pow10[length-1] <= rg.end; length++ {
		segmentStart := maxInt64(rg.start, pow10[length-1])
		segmentEnd := minInt64(rg.end, pow10[length]-1)
		if segmentStart > segmentEnd {
			continue
		}
		sum += sumInvalidLengthSegment(segmentStart, segmentEnd, length)
	}
	return sum
}

func sumInvalidLengthSegment(start, end int64, length int) int64 {
	divisors := properDivisors(length)
	if len(divisors) == 0 {
		return 0
	}
	sums := make(map[int]int64, len(divisors))
	var total int64

	for idx, d := range divisors {
		repeats := length / d
		multiplier := repeatMultiplier(d, repeats)

		loBase := pow10[d-1]
		hiBase := pow10[d] - 1

		if candidate := ceilDiv(start, multiplier); candidate > loBase {
			loBase = candidate
		}
		if candidate := end / multiplier; candidate < hiBase {
			hiBase = candidate
		}

		var g int64
		if loBase <= hiBase {
			count := hiBase - loBase + 1
			sumBases := (loBase + hiBase) * count / 2
			g = sumBases * multiplier
		}

		for j := 0; j < idx; j++ {
			smaller := divisors[j]
			if d%smaller == 0 {
				g -= sums[smaller]
			}
		}

		sums[d] = g
		total += g
	}

	return total
}

func properDivisors(n int) []int {
	if n <= 1 {
		return nil
	}
	var divs []int
	for d := 1; d*d <= n; d++ {
		if n%d != 0 {
			continue
		}
		if d < n {
			divs = append(divs, d)
		}
		other := n / d
		if other != d && other < n {
			divs = append(divs, other)
		}
	}
	sort.Ints(divs)
	return divs
}

func repeatMultiplier(blockLen, repeats int) int64 {
	totalDigits := blockLen * repeats
	return (pow10[totalDigits] - 1) / (pow10[blockLen] - 1)
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
