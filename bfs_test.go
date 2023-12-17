package main

import (
	"testing"
	"time"
)

const cubeSide = 400
var edges = initCubicGraph(cubeSide)

// func TestSeqBfs(t *testing.T) {
//   start := time.Now()
// 	ans := seqBFS(edges, 0, cubeSide)
//   elapsed := time.Since(start)

//   t.Logf("Sequential BFS-only execution time: %v", elapsed)

// 	for i := 0; i < cubeSide; i++ {
// 		for j := 0; j < cubeSide; j++ {
// 			for k := 0; k < cubeSide; k++ {
// 				vNum := idFromIJK(i, j, k, cubeSide)
// 				if ans[vNum] != i + j + k {
// 					t.Errorf("At position (%v, %v, %v) dist = %v; want %v", i, j, k, ans[vNum], i + j + k)
//           t.FailNow()
// 				}
// 			}
// 		}
// 	}
// }

func TestParBfs(t *testing.T) {
  start := time.Now()
	ans := parBFS(edges, 0, cubeSide)
  elapsed := time.Since(start)

  t.Logf("Parallel BFS-only execution time: %v", elapsed)

	for i := 0; i < cubeSide; i++ {
		for j := 0; j < cubeSide; j++ {
			for k := 0; k < cubeSide; k++ {
				vNum := idFromIJK(i, j, k, cubeSide)
        actual := ans[vNum]
				if actual != i + j + k {
					t.Errorf("At position (%v, %v, %v) dist = %v; want %v", i, j, k, actual, i + j + k)
				}
			}
		}
	}
}
