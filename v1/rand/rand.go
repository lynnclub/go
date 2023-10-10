package rand

import (
	"math/rand"
	"time"
)

// Range 区间随机 nanoseconds种子，全开区间[min, max]
func Range(min, max int) int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(max+1-min) + min
}

// Range64 区间随机 nanoseconds种子，全开区间[min, max]
func Range64(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())

	return rand.Int63n(max+1-min) + min
}

// Probability 根据概率随机
func Probability(list map[int]int64) int {
	hit := 0
	if len(list) <= 0 {
		return hit
	}

	var total int64
	total = 0
	for _, num := range list {
		total += num
	}
	if total < 1 {
		return hit
	}

	randNum := Range64(1, total)

	var increase int64 = 0
	for k, num := range list {
		increase += num
		if randNum <= increase {
			hit = k
			break
		}
	}

	return hit
}
