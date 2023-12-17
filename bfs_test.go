package main

import (
	"testing"
)

func TestSeqBfs(t *testing.T) {
	cubeSide := 10
	ans := seqBFS(initCubicGraph(cubeSide), 0, cubeSide)

	for i := 0; i < cubeSide; i++ {
		for j := 0; j < cubeSide; j++ {
			for k := 0; k < cubeSide; k++ {
				vNum := idFromIJK(i, j, k, cubeSide)
				if ans[vNum] != i + j + k {
					t.Errorf("At position (%v, %v, %v) dist = %v; want %v", i, j, k, ans[vNum], i + j + k)
          t.FailNow()
				}
			}
		}
	}
}

// func TestParBfs(t *testing.T) {
// 	cubeSide := 10
// 	ans := parBFS(initCubicGraph(cubeSide), 0, cubeSide)

// 	for i := 0; i < cubeSide; i++ {
// 		for j := 0; j < cubeSide; j++ {
// 			for k := 0; k < cubeSide; k++ {
// 				vNum := idFromIJK(i, j, k, cubeSide)
//         actual := int(ans[vNum].Load())
// 				if actual != i + j + k {
// 					t.Errorf("At position (%v, %v, %v) dist = %v; want %v", i, j, k, actual, i + j + k)
// 				}
// 			}
// 		}
// 	}
// }
