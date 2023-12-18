package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

// 開始

func nVertices(cubeSide int) int {
	return cubeSide * cubeSide * cubeSide
}

func idFromIJK(i, j, k, cubeSide int) int {
	return i*cubeSide*cubeSide + j*cubeSide + k
}

func tryAdd(edgesFromV []int, i, j, k, cubeSide int) []int {
	if i >= cubeSide || j >= cubeSide || k >= cubeSide {
		return edgesFromV
	}

	return append(edgesFromV, idFromIJK(i, j, k, cubeSide))
}

func initCubicGraph(cubeSide int) [][]int {
	edges := make([][]int, nVertices(cubeSide))
	for i := 0; i < cubeSide; i++ {
		for j := 0; j < cubeSide; j++ {
			for k := 0; k < cubeSide; k++ {
				cur := idFromIJK(i, j, k, cubeSide)
				edges[cur] = make([]int, 0, 3)
				edges[cur] = tryAdd(edges[cur], i+1, j, k, cubeSide)
				edges[cur] = tryAdd(edges[cur], i, j+1, k, cubeSide)
				edges[cur] = tryAdd(edges[cur], i, j, k+1, cubeSide)
			}
		}
	}
	return edges
}

func seqBFS(edges [][]int, start, cubeSide int) []int {
	ans := make([]int, nVertices(cubeSide))
	// visited := make([]bool, nVertices(cubeSide))
	for i := range ans {
		ans[i] = -1
	}

	q := []int{}
	q = append(q, start)
	ans[start] = 0

	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, to := range edges[v] {
			if ans[to] == -1 {
			// if !visited[to] {
				ans[to] = ans[v] + 1
				// visited[to] = true
				q = append(q, to)
			}
		}
	}

	return ans
}

func magicFrom500() int {
	return 565 * 1000
}

func parBFS(edges [][]int, start, cubeSide int) []atomic.Int32 {
	front := make([]int, 0, magicFrom500())
	front = append(front, start)

	dist := make([]atomic.Int32, nVertices(cubeSide))
	newFront := make([]int, 0, magicFrom500())

	for i := 0; i < (cubeSide-1)*3; i++ {
		degs := parScan(front, 0, len(front), func(a, b int) int {
			return a + len(edges[b])
		}, 0)

		newFront = newFront[:degs[len(degs) - 1]]
		parFor2(newFront, 0, len(newFront), func(pos, v int) {
			newFront[pos] = 0
		})

		parFor2(front, 0, len(front), func(pos, v int) {
			shift := 0
			if pos > 0 {
				shift = degs[pos-1]
			}

			for _, to := range edges[v] {
				tmp := dist[v].Load() + 1
				if dist[to].CompareAndSwap(0, tmp) {
					newFront[shift] = to
					shift++
				}
			}
		})

		front = parFilter(newFront, 0, len(newFront), func(a int) bool {
			return a != 0
		}, front)
	}
	return dist
}


func doTest[R atomic.Int32 | int](edges [][]int, bfs func([][]int, int, int) []R, testName string, cubeSide int) {
	const nTests = 5
	var totalTime time.Duration

	fmt.Printf("%v, averaged over %v launches\n\n", testName, nTests)

	for i := 0; i < nTests; i++ {
		start := time.Now()

		// f, err := os.Create("cpu.prof")
		// if err != nil {
		// 	log.Fatal("could not create CPU profile: ", err)
		// }
		// defer f.Close() // error handling omitted for example
		// if err := pprof.StartCPUProfile(f); err != nil {
		// 	log.Fatal("could not start CPU profile: ", err)
		// }

		bfs(edges, 0, cubeSide)

		// pprof.StopCPUProfile()

		elapsed := time.Since(start)
		totalTime += elapsed

		fmt.Printf("Launch %v: %v\n", i+1, elapsed)
	}

	avgTime := totalTime / nTests
	fmt.Printf("\nAverage time: %v\n", avgTime)
}

func main() {
	cubeSide := 400
	edges := initCubicGraph(cubeSide)
	doTest(edges, parBFS, "Parallel BFS", cubeSide)
	doTest(edges, seqBFS, "Sequential BFS", cubeSide)
}
