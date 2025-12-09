package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type interval struct {
	start int64
	end   int64
}

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
	if _, err := os.Stat("Day5/input.txt"); err == nil {
		return "Day5/input.txt"
	}
	return "input.txt"
}

func Solve(r io.Reader) (int, int64, error) {
	ranges, ids, err := parseInput(r)
	if err != nil {
		return 0, 0, err
	}
	if len(ranges) == 0 {
		return 0, 0, fmt.Errorf("no ranges provided")
	}

	merged := mergeIntervals(ranges)
	countFresh := countFreshIDs(ids, merged)
	totalFresh := totalCovered(merged)

	return countFresh, totalFresh, nil
}

func parseInput(r io.Reader) ([]interval, []int64, error) {
	scanner := bufio.NewScanner(r)
	var ranges []interval
	var ids []int64
	section := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			section++
			continue
		}
		if section == 0 {
			iv, err := parseInterval(line)
			if err != nil {
				return nil, nil, err
			}
			ranges = append(ranges, iv)
		} else {
			id, err := parseInt64(line)
			if err != nil {
				return nil, nil, err
			}
			ids = append(ids, id)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	return ranges, ids, nil
}

func parseInterval(line string) (interval, error) {
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return interval{}, fmt.Errorf("invalid range %q", line)
	}
	start, err := parseInt64(strings.TrimSpace(parts[0]))
	if err != nil {
		return interval{}, err
	}
	end, err := parseInt64(strings.TrimSpace(parts[1]))
	if err != nil {
		return interval{}, err
	}
	if start > end {
		start, end = end, start
	}
	return interval{start: start, end: end}, nil
}

func parseInt64(s string) (int64, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %w", s, err)
	}
	return value, nil
}

func mergeIntervals(ranges []interval) []interval {
	if len(ranges) == 0 {
		return nil
	}
	sort.Slice(ranges, func(i, j int) bool {
		if ranges[i].start == ranges[j].start {
			return ranges[i].end < ranges[j].end
		}
		return ranges[i].start < ranges[j].start
	})
	merged := make([]interval, 0, len(ranges))
	current := ranges[0]
	for i := 1; i < len(ranges); i++ {
		if ranges[i].start <= current.end+1 {
			if ranges[i].end > current.end {
				current.end = ranges[i].end
			}
			continue
		}
		merged = append(merged, current)
		current = ranges[i]
	}
	merged = append(merged, current)
	return merged
}

func countFreshIDs(ids []int64, merged []interval) int {
	count := 0
	for _, id := range ids {
		if idInIntervals(id, merged) {
			count++
		}
	}
	return count
}

func idInIntervals(id int64, merged []interval) bool {
	idx := sort.Search(len(merged), func(i int) bool {
		return merged[i].end >= id
	})
	if idx == len(merged) {
		return false
	}
	iv := merged[idx]
	return iv.start <= id && id <= iv.end
}

func totalCovered(merged []interval) int64 {
	var total int64
	for _, iv := range merged {
		total += iv.end - iv.start + 1
	}
	return total
}
