package main

import (
	"math"
	"sync"
)

const parBlockSize = 1000
const scanBlockSize = 1000

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sum(a, b int) int {
	return a + b
}

func parFor(n int, f func(int)) {
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(curBlock int) {
			f(curBlock)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func parFor2(arr []int, f func(int, int)) {
	blocks := int(math.Ceil(float64(len(arr)) / parBlockSize))
	var wg sync.WaitGroup

	wg.Add(blocks)
	for i := 0; i < blocks; i++ {
		if i != blocks - 1 {
			go func(curBlock int) {
				for k := curBlock*parBlockSize; k < min((curBlock+1)*parBlockSize, len(arr)); k++ {
					val := arr[k]
					f(k, val)
				}
				wg.Done()
			}(i)
		} else {
			curBlock := i
			for k := curBlock*parBlockSize; k < min((curBlock+1)*parBlockSize, len(arr)); k++ {
				val := arr[k]
				f(k, val)
			}
			wg.Done()
		}
	}
	wg.Wait()
}

func parScan(a []int, l, r int, f func(int, int) int, startVal int) []int {
	if r-l < scanBlockSize {
		ans := make([]int, r-l)

		curVal := startVal
		for i := l; i < r; i++ {
			curVal = f(curVal, a[i])
			ans[i-l] = curVal
		}

		return ans
	}

	blocks := int(math.Ceil(float64(r-l) / scanBlockSize))
	sums := make([]int, blocks)

	parFor(blocks, func(curBlock int) {
		curBlockVal := 0
		for k := l + curBlock*scanBlockSize; k < min(l+(curBlock+1)*scanBlockSize, r); k++ {
			curBlockVal = f(curBlockVal, a[k])
		}
		sums[curBlock] = curBlockVal
	})

	sums = parScan(sums, 0, len(sums), sum, 0)
	ans := make([]int, r-l)

	parFor(blocks, func(curBlock int) {
		curBlockVal := 0
		if curBlock > 0 {
			curBlockVal = sums[curBlock-1]
		}

		for k := l + curBlock*scanBlockSize; k < min(l+(curBlock+1)*scanBlockSize, r); k++ {
			curBlockVal = f(curBlockVal, a[k])
			ans[k-l] = curBlockVal
		}
	})

	return ans
}

func parMap(a, b []int, l, r int, f func(int) int) {
	if r-l < parBlockSize {
		for i := l; i < r; i++ {
			b[i] = f(a[i])
		}
		return
	}

	m := l + (r-l)/2

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		parMap(a, b, l, m, f)
		wg.Done()
	}()
	parMap(a, b, m, r, f)
	wg.Wait()
}

func parFilter(a []int, l, r int, f func(int) bool) []int {
	if r-l < parBlockSize {
		var ans []int
		for i := l; i < r; i++ {
			if f(a[i]) {
				ans = append(ans, a[i])
			}
		}
		return ans
	}

	flags := make([]int, r-l)
	parMap(a, flags, l, r, func(x int) int {
		if f(x) {
			return 1
		}
		return 0
	})

	blocks := int(math.Ceil(float64(r-l) / parBlockSize))
	sums := make([]int, blocks)

	parFor(blocks, func(curBlock int) {
		curBlockVal := 0
		for k := l + curBlock*parBlockSize; k < min(l+(curBlock+1)*parBlockSize, r); k++ {
			curBlockVal = curBlockVal + flags[k]
		}
		sums[curBlock] = curBlockVal
	})

	sums = parScan(sums, 0, len(sums), sum, 0)
	ans := make([]int, sums[len(sums)-1])

	parFor(blocks, func(curBlock int) {
		shift := 0
		if curBlock > 0 {
			shift = sums[curBlock-1]
		}
		lastWritten := shift
		for k := l + curBlock*parBlockSize; k < min(l+(curBlock+1)*parBlockSize, r); k++ {
			if flags[k] == 1 {
				ans[lastWritten] = a[k]
				lastWritten += 1
			}
		}
	})

	return ans
}
