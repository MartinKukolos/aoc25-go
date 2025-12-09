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

type point struct {
	x int
	y int
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
	if _, err := os.Stat("Day9/input.txt"); err == nil {
		return "Day9/input.txt"
	}
	return "input.txt"
}

func Solve(r io.Reader) (int64, int64, error) {
	pts, err := parsePoints(r)
	if err != nil {
		return 0, 0, err
	}
	if len(pts) < 2 {
		return 0, 0, fmt.Errorf("need at least two points")
	}

	part1 := maxRectangleAny(pts)

	xVals := make([]int, len(pts))
	yVals := make([]int, len(pts))
	for i, p := range pts {
		xVals[i] = p.x
		yVals[i] = p.y
	}
	xs := uniqueSorted(xVals)
	ys := uniqueSorted(yVals)

	if len(xs) < 2 || len(ys) < 2 {
		return part1, 0, nil
	}

	xIndex := make(map[int]int, len(xs))
	for i, v := range xs {
		xIndex[v] = i
	}
	yIndex := make(map[int]int, len(ys))
	for i, v := range ys {
		yIndex[v] = i
	}

	inside, err := buildInsideGrid(pts, xs, ys, xIndex)
	if err != nil {
		return 0, 0, err
	}
	prefix := buildPrefix(inside)

	part2 := maxRectangleInside(pts, xIndex, yIndex, prefix, xs, ys)

	return part1, part2, nil
}

func parsePoints(r io.Reader) ([]point, error) {
	scanner := bufio.NewScanner(r)
	var pts []point
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid coordinate %q", line)
		}
		x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid x in %q: %w", line, err)
		}
		y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid y in %q: %w", line, err)
		}
		pts = append(pts, point{x: x, y: y})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return pts, nil
}

func uniqueSorted(values []int) []int {
	sort.Ints(values)
	result := make([]int, 0, len(values))
	prevSet := false
	var prev int
	for _, v := range values {
		if !prevSet || v != prev {
			result = append(result, v)
			prev = v
			prevSet = true
		}
	}
	return result
}

func maxRectangleAny(pts []point) int64 {
	var best int64
	for i := 0; i < len(pts); i++ {
		for j := i + 1; j < len(pts); j++ {
			area := tileArea(pts[i], pts[j])
			if area > best {
				best = area
			}
		}
	}
	return best
}

func buildInsideGrid(pts []point, xs, ys []int, xIndex map[int]int) ([][]bool, error) {
	rows := len(ys) - 1
	cols := len(xs) - 1
	inside := make([][]bool, rows)
	for i := range inside {
		inside[i] = make([]bool, cols)
	}
	var intersections []int
	for row := 0; row < rows; row++ {
		yLow := ys[row]
		yHigh := ys[row+1]
		yMid := (float64(yLow) + float64(yHigh)) * 0.5
		intersections = intersections[:0]
		for i := 0; i < len(pts); i++ {
			a := pts[i]
			b := pts[(i+1)%len(pts)]
			if a.x == b.x {
				yMin := minInt(a.y, b.y)
				yMax := maxInt(a.y, b.y)
				if yMid > float64(yMin) && yMid < float64(yMax) {
					intersections = append(intersections, a.x)
				}
			}
		}
		sort.Ints(intersections)
		if len(intersections)%2 != 0 {
			return nil, fmt.Errorf("row %d has odd number of intersections", row)
		}
		for k := 0; k+1 < len(intersections); k += 2 {
			xL := intersections[k]
			xR := intersections[k+1]
			leftIdx, okL := xIndex[xL]
			rightIdx, okR := xIndex[xR]
			if !okL || !okR {
				return nil, fmt.Errorf("intersection coordinate missing from x list")
			}
			for col := leftIdx; col < rightIdx; col++ {
				inside[row][col] = true
			}
		}
	}
	return inside, nil
}

func buildPrefix(inside [][]bool) [][]int {
	rows := len(inside)
	cols := 0
	if rows > 0 {
		cols = len(inside[0])
	}
	prefix := make([][]int, rows+1)
	for i := range prefix {
		prefix[i] = make([]int, cols+1)
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			val := 0
			if inside[i][j] {
				val = 1
			}
			prefix[i+1][j+1] = val + prefix[i][j+1] + prefix[i+1][j] - prefix[i][j]
		}
	}
	return prefix
}

func maxRectangleInside(pts []point, xIndex, yIndex map[int]int, prefix [][]int, xs, ys []int) int64 {
	var best int64
	for i := 0; i < len(pts); i++ {
		xi := xIndex[pts[i].x]
		yi := yIndex[pts[i].y]
		for j := i + 1; j < len(pts); j++ {
			xj := xIndex[pts[j].x]
			yj := yIndex[pts[j].y]
			xL := minInt(xi, xj)
			xH := maxInt(xi, xj)
			yL := minInt(yi, yj)
			yH := maxInt(yi, yj)
			if xL == xH || yL == yH {
				continue
			}
			area := tileArea(pts[i], pts[j])
			if area <= best {
				continue
			}
			if rectangleInside(prefix, xL, xH, yL, yH) {
				best = area
			}
		}
	}
	return best
}

func rectangleInside(prefix [][]int, xL, xH, yL, yH int) bool {
	if xL >= xH || yL >= yH {
		return false
	}
	areacells := (xH - xL) * (yH - yL)
	sum := prefix[yH][xH] - prefix[yL][xH] - prefix[yH][xL] + prefix[yL][xL]
	return sum == areacells
}

func tileArea(a, b point) int64 {
	width := absInt(a.x-b.x) + 1
	height := absInt(a.y-b.y) + 1
	return int64(width) * int64(height)
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
