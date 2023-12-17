package main

import (
	"sync/atomic"
)

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
	for i := range ans {
		ans[i] = -1
	}

	q := NewRingBufferQueue()
	ans[start] = 0
	q.push(start)

	for !q.empty() {
		v := q.pop()
		for _, to := range edges[v] {
			if ans[to] == -1 {
				ans[to] = ans[v] + 1
				q.push(to)
			}
		}
	}

	return ans
}

func parBFS(edges [][]int, start, cubeSide int) []int {
	front := []int{}
	front = append(front, start)

	visited := make([]atomic.Bool, nVertices(cubeSide))
	dist := make([]int, nVertices(cubeSide))

	for i := 0; i < (cubeSide-1)*3; i++ { // TODO better handling on when we are done
		degs := parScan(front, 0, len(front), func(a, b int) int {
			return a + len(edges[b])
		}, 0)

		newFront := make([]int, degs[len(degs)-1])

		parFor2(front, func(pos, v int) {
			curShift := 0
			blockShift := 0
			if pos > 0 {
				blockShift = degs[pos-1]
			}

			for _, to := range edges[v] {
				if visited[to].CompareAndSwap(false, true) {
					newFront[blockShift+curShift] = to
					dist[to] = dist[v] + 1
					curShift++
				}
			}
		})

		front = parFilter(newFront, 0, len(newFront), func(a int) bool {
			return a != 0
		})
	}
	return dist
}

func main() {
}
