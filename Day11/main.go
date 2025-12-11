package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
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
	if _, err := os.Stat("Day11/input.txt"); err == nil {
		return "Day11/input.txt"
	}
	return "input.txt"
}

// Solve reads a directed graph specification and returns:
// - number of distinct simple paths from "you" to "out"
// - number of distinct simple paths from "svr" to "out" that visit both "dac" and "fft"
func Solve(r io.Reader) (int64, int64, error) {
	graph, err := parseGraph(r)
	if err != nil {
		return 0, 0, err
	}

	// Part 1: paths from "you" to "out"
	var p1 int64
	if _, ok := graph["you"]; ok {
		p1 = countPathsSimple(graph, "you", "out")
	} else {
		p1 = 0
	}

	// Part 2: paths from "svr" to "out" that visit both dac and fft
	var p2 int64
	if _, ok := graph["svr"]; ok {
		p2 = countPathsWithMustVisit(graph, "svr", "out", "dac", "fft")
	} else {
		p2 = 0
	}

	return p1, p2, nil
}

func parseGraph(r io.Reader) (map[string][]string, error) {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 1024)
	scanner.Buffer(buf, 1<<20)
	graph := make(map[string][]string)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// allow comments in input
		if strings.HasPrefix(line, "#") {
			continue
		}
		// Expect format: name: a b c
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			// If the line has no colon, skip it gracefully
			continue
		}
		src := strings.TrimSpace(parts[0])
		targets := strings.Fields(strings.TrimSpace(parts[1]))
		// ensure node exists even with no targets
		if _, exists := graph[src]; !exists {
			graph[src] = nil
		}
		if len(targets) > 0 {
			graph[src] = append(graph[src], targets...)
			// ensure target nodes appear in map too (even if no outgoing list given elsewhere)
			for _, t := range targets {
				if _, ok := graph[t]; !ok {
					graph[t] = nil
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return graph, nil
}

func countPathsSimple(graph map[string][]string, start, goal string) int64 {
	var total int64
	visited := make(map[string]bool)
	var dfs func(string)
	dfs = func(u string) {
		if u == goal {
			total++
			return
		}
		visited[u] = true
		for _, v := range graph[u] {
			if !visited[v] {
				dfs(v)
			}
		}
		visited[u] = false
	}
	dfs(start)
	return total
}

func countPathsWithMustVisit(graph map[string][]string, start, goal, mustA, mustB string) int64 {
	var total int64
	visited := make(map[string]bool)
	var dfs func(string, bool, bool)
	dfs = func(u string, seenA, seenB bool) {
		if u == mustA {
			seenA = true
		}
		if u == mustB {
			seenB = true
		}
		if u == goal {
			if seenA && seenB {
				total++
			}
			return
		}
		visited[u] = true
		for _, v := range graph[u] {
			if !visited[v] {
				dfs(v, seenA, seenB)
			}
		}
		visited[u] = false
	}
	dfs(start, false, false)
	return total
}
