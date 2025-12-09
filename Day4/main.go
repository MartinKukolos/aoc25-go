package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var neighborOffsets = [8][2]int{
	{1, 0}, {-1, 0}, {0, 1}, {0, -1},
	{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
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
	if _, err := os.Stat("Day4/input.txt"); err == nil {
		return "Day4/input.txt"
	}
	return "input.txt"
}

func Solve(r io.Reader) (int, int, error) {
	grid, err := readGrid(r)
	if err != nil {
		return 0, 0, err
	}
	if len(grid) == 0 {
		return 0, 0, fmt.Errorf("empty grid")
	}

	part1 := countAccessible(grid)
	part2 := totalRemovable(grid)

	return part1, part2, nil
}

func readGrid(r io.Reader) ([][]bool, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	width := -1
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if width == -1 {
			width = len(line)
		} else if len(line) != width {
			return nil, fmt.Errorf("inconsistent row length: got %d want %d", len(line), width)
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	grid := make([][]bool, len(lines))
	for i, line := range lines {
		row := make([]bool, len(line))
		for j, ch := range line {
			if ch == '@' {
				row[j] = true
			}
		}
		grid[i] = row
	}
	return grid, nil
}

func countAccessible(grid [][]bool) int {
	rows := len(grid)
	cols := len(grid[0])
	count := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if !grid[r][c] {
				continue
			}
			if neighborCount(grid, r, c) < 4 {
				count++
			}
		}
	}
	return count
}

func totalRemovable(grid [][]bool) int {
	rows := len(grid)
	cols := len(grid[0])
	current := make([][]bool, rows)
	for i := range grid {
		row := make([]bool, cols)
		copy(row, grid[i])
		current[i] = row
	}
	total := 0
	for {
		var toRemove [][2]int
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				if !current[r][c] {
					continue
				}
				if neighborCount(current, r, c) < 4 {
					toRemove = append(toRemove, [2]int{r, c})
				}
			}
		}
		if len(toRemove) == 0 {
			break
		}
		for _, rc := range toRemove {
			current[rc[0]][rc[1]] = false
		}
		total += len(toRemove)
	}
	return total
}

func neighborCount(grid [][]bool, r, c int) int {
	rows := len(grid)
	cols := len(grid[0])
	count := 0
	for _, off := range neighborOffsets {
		r2 := r + off[0]
		c2 := c + off[1]
		if r2 < 0 || r2 >= rows || c2 < 0 || c2 >= cols {
			continue
		}
		if grid[r2][c2] {
			count++
		}
	}
	return count
}
