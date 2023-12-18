package main

import (
	"math"
	"runtime/debug"
	"testing"
	"time"
)

const cubeSide = 500

var edges = initCubicGraph(cubeSide)

func TestParBfs(t *testing.T) {
	gcpercent := debug.SetGCPercent(-1)
	memlimit := debug.SetMemoryLimit(math.MaxInt64)

	start := time.Now()
	ans := parBFS(edges, 0, cubeSide)
	elapsed := time.Since(start)

	t.Logf("Parallel BFS-only execution time: %v", elapsed)

	debug.SetGCPercent(gcpercent)
	debug.SetMemoryLimit(memlimit)

	for i := 0; i < cubeSide; i++ {
		for j := 0; j < cubeSide; j++ {
			for k := 0; k < cubeSide; k++ {
				vNum := idFromIJK(i, j, k, cubeSide)
				actual := int(ans[vNum].Load())
				if actual != i+j+k {
					t.Errorf("At position (%v, %v, %v) dist = %v; want %v", i, j, k, actual, i+j+k)
				}
			}
		}
	}
}

func TestSeqBfs(t *testing.T) {
	start := time.Now()
	ans := seqBFS(edges, 0, cubeSide)
	elapsed := time.Since(start)

	t.Logf("Sequential BFS-only execution time: %v", elapsed)

	for i := 0; i < cubeSide; i++ {
		for j := 0; j < cubeSide; j++ {
			for k := 0; k < cubeSide; k++ {
				vNum := idFromIJK(i, j, k, cubeSide)
				if ans[vNum] != i+j+k {
					t.Errorf("At position (%v, %v, %v) dist = %v; want %v", i, j, k, ans[vNum], i+j+k)
					t.FailNow()
				}
			}
		}
	}
}


// seq-par-bfs (174329) +45 -41 via 🐹 v1.19.3 ❯ GOMAXPROCS=4 go test -v
// === RUN   TestParBfs
//     bfs_test.go:22: Parallel BFS-only execution time: 18.868629686s
// --- PASS: TestParBfs (19.03s)
// === RUN   TestSeqBfs
//     bfs_test.go:45: Sequential BFS-only execution time: 45.299065821s
// --- PASS: TestSeqBfs (45.46s)
// === RUN   TestQueue
// --- PASS: TestQueue (0.00s)
// PASS
// ok      seq-par-bfs     71.466s
