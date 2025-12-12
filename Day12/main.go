package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type shape struct {
	// A list of unique orientations; each orientation is a list of (x,y) cells
	orients [][]pt
	area    int
}

type pt struct{ x, y int }

func main() {
	path := resolveInputPath(os.Args)

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open input %q: %v\n", path, err)
		os.Exit(1)
	}
	defer f.Close()

	p1, p2, err := Solve(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "solve error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Part 1: %d\n", p1)
	fmt.Printf("Part 2: %d\n", p2)
}

func resolveInputPath(args []string) string {
	if len(args) > 1 {
		return args[1]
	}
	if _, err := os.Stat("Day12/input.txt"); err == nil {
		return "Day12/input.txt"
	}
	return "input.txt"
}

// Solve parses shapes and regions; returns how many regions can fit the requested presents (part1).
// There is no Part 2 for this day; it returns 0.
func Solve(r io.Reader) (int64, int64, error) {
	shapes, regions, err := parseInput(r)
	if err != nil {
		return 0, 0, err
	}

	var count int64
	for _, reg := range regions {
		if canFitAll(shapes, reg.w, reg.h, reg.counts) {
			count++
		}
	}
	return count, 0, nil
}

type regionSpec struct {
	w, h   int
	counts []int
}

func parseInput(r io.Reader) ([]shape, []regionSpec, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 1024), 1<<20)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	// Parse shapes until the first region line like "12x5: ..."
	regionRe := regexp.MustCompile(`^\s*(\d+)x(\d+)\s*:`)
	i := 0
	var rawShapes [][][]rune
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			i++
			continue
		}
		if regionRe.MatchString(line) {
			break
		}
		// Expect something like "0:" then grid lines of '#' '.' until blank line
		if !strings.HasSuffix(line, ":") {
			// Skip unexpected lines
			i++
			continue
		}
		i++
		var grid [][]rune
		for i < len(lines) {
			l := strings.TrimRight(lines[i], "\r\n")
			if strings.TrimSpace(l) == "" {
				break
			}
			// grid line: consist of '#' and '.'
			row := []rune(strings.TrimSpace(l))
			grid = append(grid, row)
			i++
		}
		rawShapes = append(rawShapes, grid)
		// consume blank line if present
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
	}

	// Build shapes with orientations
	shapes := make([]shape, 0, len(rawShapes))
	for _, grid := range rawShapes {
		cells := extractCells(grid)
		orients := genOrientations(cells)
		s := shape{orients: orients, area: len(cells)}
		shapes = append(shapes, s)
	}

	// Parse regions
	var regions []regionSpec
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		i++
		if line == "" {
			continue
		}
		m := regionRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		w, _ := strconv.Atoi(m[1])
		h, _ := strconv.Atoi(m[2])
		after := strings.TrimSpace(line[strings.Index(line, ":")+1:])
		var counts []int
		if after != "" {
			parts := strings.Fields(after)
			counts = make([]int, len(parts))
			for j, p := range parts {
				v, err := strconv.Atoi(p)
				if err != nil {
					v = 0
				}
				counts[j] = v
			}
		}
		// Pad or trim counts to number of shapes
		if len(counts) < len(shapes) {
			tmp := make([]int, len(shapes))
			copy(tmp, counts)
			counts = tmp
		} else if len(counts) > len(shapes) {
			counts = counts[:len(shapes)]
		}
		regions = append(regions, regionSpec{w: w, h: h, counts: counts})
	}

	return shapes, regions, nil
}

func extractCells(grid [][]rune) []pt {
	var cells []pt
	for y, row := range grid {
		for x, ch := range row {
			if ch == '#' {
				cells = append(cells, pt{x: x, y: y})
			}
		}
	}
	return normalizeCells(cells)
}

func normalizeCells(cells []pt) []pt {
	if len(cells) == 0 {
		return cells
	}
	minx, miny := cells[0].x, cells[0].y
	for _, c := range cells {
		if c.x < minx {
			minx = c.x
		}
		if c.y < miny {
			miny = c.y
		}
	}
	out := make([]pt, len(cells))
	for i, c := range cells {
		out[i] = pt{c.x - minx, c.y - miny}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].y == out[j].y {
			return out[i].x < out[j].x
		}
		return out[i].y < out[j].y
	})
	return out
}

func genOrientations(base []pt) [][]pt {
	// Generate rotations and flips; deduplicate by normalized key
	type keyT string
	seen := make(map[keyT]bool)
	var res [][]pt

	transforms := []func(pt) pt{
		func(p pt) pt { return pt{p.x, p.y} },   // identity
		func(p pt) pt { return pt{-p.x, p.y} },  // flip X
		func(p pt) pt { return pt{p.x, -p.y} },  // flip Y
		func(p pt) pt { return pt{-p.x, -p.y} }, // flip both
	}
	rot := func(p pt) pt { return pt{-p.y, p.x} } // 90deg

	add := func(points []pt) {
		n := normalizeCells(points)
		var b strings.Builder
		for _, c := range n {
			b.WriteString(strconv.Itoa(c.x))
			b.WriteByte(',')
			b.WriteString(strconv.Itoa(c.y))
			b.WriteByte(';')
		}
		k := keyT(b.String())
		if !seen[k] {
			seen[k] = true
			res = append(res, n)
		}
	}

	for _, tf := range transforms {
		// apply flip transform, then try 4 rotations
		cur := make([]pt, len(base))
		for i, p := range base {
			cur[i] = tf(p)
		}
		r0 := cur
		add(r0)
		r1 := make([]pt, len(r0))
		for i, p := range r0 {
			r1[i] = rot(p)
		}
		add(r1)
		r2 := make([]pt, len(r1))
		for i, p := range r1 {
			r2[i] = rot(p)
		}
		add(r2)
		r3 := make([]pt, len(r2))
		for i, p := range r2 {
			r3[i] = rot(p)
		}
		add(r3)
	}
	return res
}

func canFitAll(shapes []shape, W, H int, counts []int) bool {
	// Fast area feasibility checks
	totalArea := 0
	for i, c := range counts {
		if c < 0 {
			return false
		}
		totalArea += c * shapes[i].area
	}
	cells := W * H
	if totalArea > cells {
		return false
	}

	// Build exact cover model (Algorithm X / DLX):
	// Columns:
	//  - one for each grid cell (cells columns)
	//  - one for each individual shape instance (sum(counts) columns)
	// Rows:
	//  - for each feasible placement of a shape instance: set 1s for the occupied cells and for that instance column
	//  - if totalArea < cells, add an "empty" row for each cell that covers only that cell column

	nCellCols := cells
	// allocate instance columns
	instColsByShape := make([][]int, len(shapes))
	nextCol := nCellCols
	for i, c := range counts {
		if c <= 0 {
			continue
		}
		instColsByShape[i] = make([]int, c)
		for k := 0; k < c; k++ {
			instColsByShape[i][k] = nextCol
			nextCol++
		}
	}
	nCols := nextCol

	// Precompute bounding boxes for orientations to limit origins
	type orientInfo struct {
		cells []pt
		maxx  int
		maxy  int
	}
	orientInfos := make([][]orientInfo, len(shapes))
	for i, s := range shapes {
		infos := make([]orientInfo, len(s.orients))
		for j, o := range s.orients {
			mx, my := 0, 0
			for _, c := range o {
				if c.x > mx {
					mx = c.x
				}
				if c.y > my {
					my = c.y
				}
			}
			infos[j] = orientInfo{cells: o, maxx: mx, maxy: my}
		}
		orientInfos[i] = infos
	}

	// Generate rows
	rows := make([][]int, 0, 1024)
	// shape placements
	for i := range shapes {
		if len(instColsByShape[i]) == 0 {
			continue
		}
		for _, info := range orientInfos[i] {
			if info.maxx >= W || info.maxy >= H {
				// Still try origins ensuring fit
			}
			for oy := 0; oy+info.maxy < H; oy++ {
				for ox := 0; ox+info.maxx < W; ox++ {
					// collect cell columns for this placement
					cols := make([]int, 0, len(info.cells)+1)
					ok := true
					for _, c := range info.cells {
						x := ox + c.x
						y := oy + c.y
						if x < 0 || x >= W || y < 0 || y >= H {
							ok = false
							break
						}
						cols = append(cols, y*W+x)
					}
					if !ok {
						continue
					}
					// For each instance column for this shape, create a row
					for _, ic := range instColsByShape[i] {
						row := make([]int, 0, len(cols)+1)
						row = append(row, cols...)
						row = append(row, ic)
						rows = append(rows, row)
					}
				}
			}
		}
	}
	// empty cell rows if needed
	if totalArea < cells {
		for cell := 0; cell < cells; cell++ {
			rows = append(rows, []int{cell})
		}
	}

	// If there are no shape instances to place, it's always possible (we can leave cells empty)
	if nCols == nCellCols && totalArea == 0 {
		return true
	}

	if exactCoverExists(rows, nCols, nCellCols) {
		return true
	}
	// Fallback: use backtracking (safer correctness if DLX modeling missed a case)
	return canFitAllBacktrack(shapes, W, H, counts)
}

// exactCoverExists implements Algorithm X with a lightweight DLX-like state.
// It checks if there exists a subset of rows covering every column exactly once.
func exactCoverExists(rows [][]int, nCols int, nCellCols int) bool {
	if nCols == 0 {
		return true
	}
	m := len(rows)
	// Build column to rows index
	colRows := make([][]int, nCols)
	for r, cols := range rows {
		// deduplicate columns within a row to avoid double counting
		sort.Ints(cols)
		uniq := cols[:0]
		last := -1
		for _, c := range cols {
			if c < 0 || c >= nCols {
				continue
			}
			if c != last {
				uniq = append(uniq, c)
				last = c
			}
		}
		rows[r] = uniq
		for _, c := range uniq {
			colRows[c] = append(colRows[c], r)
		}
	}

	activeCol := make([]bool, nCols)
	colCount := make([]int, nCols)
	for c := 0; c < nCols; c++ {
		activeCol[c] = true
		colCount[c] = len(colRows[c])
		if colCount[c] == 0 {
			// No row can cover this column -> impossible
			return false
		}
	}
	rowBlock := make([]int, m) // how many covered columns intersect this row

	// Stacks for backtracking
	type coverOp struct {
		col        int
		affectedRs []int
	}

	// Cover a column: mark it inactive and block intersecting rows
	cover := func(col int, ops *[]coverOp) {
		activeCol[col] = false
		// Track which rows transitioned from 0->1 block due to this cover
		var changedRows []int
		// For rows that include this column
		for _, r := range colRows[col] {
			prev := rowBlock[r]
			rowBlock[r]++
			if prev == 0 {
				changedRows = append(changedRows, r)
				// This row becomes inactive; decrement counts for its other columns that are active
				for _, c2 := range rows[r] {
					if c2 == col || !activeCol[c2] {
						continue
					}
					colCount[c2]--
				}
			}
		}
		(*ops) = append(*ops, coverOp{col: col, affectedRs: changedRows})
	}

	// Uncover in reverse order
	uncover := func(op coverOp) {
		// Reactivate the column first so we also restore its count properly
		activeCol[op.col] = true
		// Reactivate rows that were deactivated by this cover
		for _, r := range op.affectedRs {
			// decrement block count first
			rowBlock[r]--
			if rowBlock[r] == 0 {
				// Row becomes active again: increment counts for its active columns (including the just-uncovered one)
				for _, c2 := range rows[r] {
					if activeCol[c2] { // only active columns maintain counts
						colCount[c2]++
					}
				}
			}
		}
		// Recompute count for the just-uncovered column accurately
		cnt := 0
		for _, r := range colRows[op.col] {
			if rowBlock[r] == 0 {
				cnt++
			}
		}
		colCount[op.col] = cnt
	}

	// Check if all columns are covered (i.e., none active)
	allCovered := func() bool {
		for c := 0; c < nCols; c++ {
			if activeCol[c] {
				return false
			}
		}
		return true
	}

	var dfs func() bool
	dfs = func() bool {
		if allCovered() {
			return true
		}
		// Choose the active column with the smallest number of candidate rows (heuristic)
		// Prioritize instance columns (>= nCellCols) to drastically reduce branching when
		// there are empty-cell options or many grid positions.
		sel := -1
		best := int(1<<31 - 1)
		// pass 1: instance columns
		for c := nCellCols; c < nCols; c++ {
			if !activeCol[c] {
				continue
			}
			cnt := colCount[c]
			if cnt == 0 {
				return false
			}
			if cnt < best {
				best = cnt
				sel = c
				if best == 1 {
					break
				}
			}
		}
		// pass 2: if no instance columns active, pick among cell columns
		if sel == -1 {
			for c := 0; c < nCellCols; c++ {
				if !activeCol[c] {
					continue
				}
				cnt := colCount[c]
				if cnt == 0 {
					return false
				}
				if cnt < best {
					best = cnt
					sel = c
					if best == 1 {
						break
					}
				}
			}
		}
		if sel == -1 { // no active columns
			return true
		}

		// Try each active row that covers sel
		for _, r := range colRows[sel] {
			if rowBlock[r] != 0 {
				continue
			}
			// Choose row r: cover all its columns
			var ops []coverOp
			for _, c2 := range rows[r] {
				if activeCol[c2] {
					cover(c2, &ops)
				}
			}
			if dfs() {
				return true
			}
			// backtrack
			for i := len(ops) - 1; i >= 0; i-- {
				uncover(ops[i])
			}
		}
		return false
	}

	return dfs()
}

// Legacy backtracking placement used as a fallback to ensure correctness
func canFitAllBacktrack(shapes []shape, W, H int, counts []int) bool {
	type item struct{ idx, area int }
	var items []item
	for i, c := range counts {
		for k := 0; k < c; k++ {
			items = append(items, item{idx: i, area: shapes[i].area})
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].area == items[j].area {
			return items[i].idx < items[j].idx
		}
		return items[i].area > items[j].area
	})
	if len(items) == 0 {
		return true
	}

	grid := make([]bool, W*H)

	var place func(pos int) bool
	place = func(pos int) bool {
		if pos == len(items) {
			return true
		}
		sIdx := items[pos].idx
		s := shapes[sIdx]
		for _, orient := range s.orients {
			// determine bounding box of orient to limit origins
			maxx, maxy := 0, 0
			for _, c := range orient {
				if c.x > maxx {
					maxx = c.x
				}
				if c.y > maxy {
					maxy = c.y
				}
			}
			for oy := 0; oy+maxy < H; oy++ {
				for ox := 0; ox+maxx < W; ox++ {
					valid := true
					for _, c := range orient {
						x := ox + c.x
						y := oy + c.y
						if grid[y*W+x] {
							valid = false
							break
						}
					}
					if !valid {
						continue
					}
					for _, c := range orient {
						grid[(oy+c.y)*W+(ox+c.x)] = true
					}
					if place(pos + 1) {
						return true
					}
					for _, c := range orient {
						grid[(oy+c.y)*W+(ox+c.x)] = false
					}
				}
			}
		}
		return false
	}

	return place(0)
}
