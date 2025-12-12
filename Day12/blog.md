# Day 12: Reactor Packing Log

Today’s puzzle is basically a tiling SAT problem wrapped in flavor text about stuffing reactor modules into rectangular rooms. The input describes:

- A gallery of polyomino-like shapes (drawn with `#` and `.`) that represent present bundles.
- A set of region requests such as `12x5: 2 1 0`, meaning “attempt to fit two copies of shape 0, one of shape 1, none of shape 2, inside a 12×5 grid”.
- Only Part 1 has an answer: report how many regions are satisfiable. Part 2 is intentionally blank for this day, so the solver returns 0 for it.

## Parsing and Canonical Shapes

I read the entire file once, separating the shape gallery from the later region lines (`WxH:`). Each shape is normalized into `(x, y)` coordinates for every filled cell, then I generate all rotations and reflections, normalizing and deduplicating them so the solver only evaluates unique orientations. That keeps the later search from reconsidering equivalent placements.

## Exact Cover Modeling

Fitting shapes into a rectangle is a textbook exact-cover problem, so I leaned on Algorithm X with a DLX-inspired implementation:

1. Columns represent every grid cell plus one column for each individual requested shape instance.
2. Rows represent feasible placements of a single shape instance—when a shape can sit at `(ox, oy)` in a given orientation, I emit a row that touches the covered cells and that instance’s column.
3. If the requested area doesn’t fill the whole region (total shape area < `W*H`), I add “empty” rows that claim a single uncovered cell so the algorithm can satisfy every column exactly once.
4. A custom cover/uncover routine, favoring instance columns first, dramatically prunes the branching factor on larger grids.

If the DLX search reports success, the region is counted as placeable.

## Safety Net Backtracking

Modeling mistakes around exact cover can be subtle, so after DLX fails I fall back to a straightforward recursive placement search: sort the remaining shape instances by area, try every orientation and origin, and backtrack when tiles collide. It’s slower (`O(b^n)` in the worst case) but only runs on DLX failures, so the correctness guard is worth the cost.

## Complexity Notes

- Shape parsing plus orientation generation is `O(s * k)` where `s` is the number of shapes and `k` the cells per shape.
- Region evaluation builds `O(W*H + totalPlacements)` rows; DLX runs in time proportional to the explored branches, which in practice stays reasonable for the puzzle sizes.
- Memory holds the row/column incidence lists and temporary stacks, all within a few megabytes even on the largest inputs.

## Testing

`Day12/main_test.go` locks in the sample region count (expecting 2 valid layouts) to ensure both the DLX pathway and the backtracking fallback remain consistent. Running `go test ./...` before submission covers this day along with the rest of the Advent progress.

Even though Part 2 never materializes, the solver now confidently tells me exactly which reactor rooms can accommodate every requested present bundle.
