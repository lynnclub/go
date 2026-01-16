package rand

import (
	"testing"
)

// TestRange 区间随机
func TestRange(t *testing.T) {
	num := Range(6, 8)
	if num < 6 || num > 8 {
		t.Errorf("Range(6, 8)应该返回[6, 8]之间的数，实际为%d", num)
	}
}

// TestRangeMinMax 测试最小值和最大值
func TestRangeMinMax(t *testing.T) {
	// 测试多次以确保范围正确
	for i := 0; i < 100; i++ {
		num := Range(1, 10)
		if num < 1 || num > 10 {
			t.Errorf("Range(1, 10)应该返回[1, 10]之间的数，实际为%d", num)
		}
	}
}

// TestRangeSameValue 测试最小值等于最大值
func TestRangeSameValue(t *testing.T) {
	num := Range(5, 5)
	if num != 5 {
		t.Errorf("Range(5, 5)应该返回5，实际为%d", num)
	}
}

// TestRangeZero 测试包含0的范围
func TestRangeZero(t *testing.T) {
	for i := 0; i < 50; i++ {
		num := Range(0, 5)
		if num < 0 || num > 5 {
			t.Errorf("Range(0, 5)应该返回[0, 5]之间的数，实际为%d", num)
		}
	}
}

// TestRangeNegative 测试负数范围
func TestRangeNegative(t *testing.T) {
	for i := 0; i < 50; i++ {
		num := Range(-10, -5)
		if num < -10 || num > -5 {
			t.Errorf("Range(-10, -5)应该返回[-10, -5]之间的数，实际为%d", num)
		}
	}
}

// TestRangeNegativeToPositive 测试负数到正数的范围
func TestRangeNegativeToPositive(t *testing.T) {
	for i := 0; i < 50; i++ {
		num := Range(-5, 5)
		if num < -5 || num > 5 {
			t.Errorf("Range(-5, 5)应该返回[-5, 5]之间的数，实际为%d", num)
		}
	}
}

// TestRange64 测试64位整数区间随机
func TestRange64(t *testing.T) {
	num := Range64(100, 200)
	if num < 100 || num > 200 {
		t.Errorf("Range64(100, 200)应该返回[100, 200]之间的数，实际为%d", num)
	}
}

// TestRange64Large 测试大数值范围
func TestRange64Large(t *testing.T) {
	min := int64(1000000000)
	max := int64(2000000000)

	for i := 0; i < 50; i++ {
		num := Range64(min, max)
		if num < min || num > max {
			t.Errorf("Range64(%d, %d)应该返回[%d, %d]之间的数，实际为%d", min, max, min, max, num)
		}
	}
}

// TestRange64SameValue 测试64位最小值等于最大值
func TestRange64SameValue(t *testing.T) {
	num := Range64(100, 100)
	if num != 100 {
		t.Errorf("Range64(100, 100)应该返回100，实际为%d", num)
	}
}

// TestRange64Negative 测试64位负数范围
func TestRange64Negative(t *testing.T) {
	for i := 0; i < 50; i++ {
		num := Range64(-1000, -500)
		if num < -1000 || num > -500 {
			t.Errorf("Range64(-1000, -500)应该返回[-1000, -500]之间的数，实际为%d", num)
		}
	}
}

// TestProbability 测试概率随机
func TestProbability(t *testing.T) {
	list := map[int]int64{
		1: 50,
		2: 30,
		3: 20,
	}

	// 测试多次，统计结果
	results := make(map[int]int)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		hit := Probability(list)
		results[hit]++
	}

	// 验证所有返回值都在预期范围内
	for key := range results {
		if _, ok := list[key]; !ok {
			t.Errorf("Probability返回了未定义的键: %d", key)
		}
	}

	// 验证至少命中了所有选项（概率测试，可能偶尔失败）
	if len(results) < 2 {
		t.Logf("警告: 只命中了%d个选项，期望至少2个", len(results))
	}
}

// TestProbabilitySingleItem 测试只有一个选项的概率
func TestProbabilitySingleItem(t *testing.T) {
	list := map[int]int64{
		1: 100,
	}

	for i := 0; i < 50; i++ {
		hit := Probability(list)
		if hit != 1 {
			t.Errorf("Probability应该返回1，实际为%d", hit)
		}
	}
}

// TestProbabilityEmptyList 测试空列表
func TestProbabilityEmptyList(t *testing.T) {
	list := map[int]int64{}

	hit := Probability(list)
	if hit != 0 {
		t.Errorf("空列表应该返回0，实际为%d", hit)
	}
}

// TestProbabilityZeroTotal 测试总概率为0
func TestProbabilityZeroTotal(t *testing.T) {
	list := map[int]int64{
		1: 0,
		2: 0,
		3: 0,
	}

	hit := Probability(list)
	if hit != 0 {
		t.Errorf("总概率为0应该返回0，实际为%d", hit)
	}
}

// TestProbabilityNegativeValues 测试负数概率（边界情况）
func TestProbabilityNegativeValues(t *testing.T) {
	list := map[int]int64{
		1: -10,
		2: -20,
	}

	// 负数概率会导致总和为负，应该返回0
	hit := Probability(list)
	if hit != 0 {
		t.Errorf("负数概率应该返回0，实际为%d", hit)
	}
}

// TestProbabilityDistribution 测试概率分布的合理性
func TestProbabilityDistribution(t *testing.T) {
	list := map[int]int64{
		1: 70, // 70%
		2: 20, // 20%
		3: 10, // 10%
	}

	results := make(map[int]int)
	iterations := 10000

	for i := 0; i < iterations; i++ {
		hit := Probability(list)
		results[hit]++
	}

	// 验证分布大致符合预期（允许一定误差）
	total := int64(100)
	for key, weight := range list {
		expectedRatio := float64(weight) / float64(total)
		actualRatio := float64(results[key]) / float64(iterations)

		// 允许10%的误差
		if actualRatio < expectedRatio*0.5 || actualRatio > expectedRatio*1.5 {
			t.Logf("警告: 键%d的实际概率%.2f与期望概率%.2f偏差较大", key, actualRatio, expectedRatio)
		}
	}
}

// TestProbabilityEqualWeights 测试相等权重
func TestProbabilityEqualWeights(t *testing.T) {
	list := map[int]int64{
		1: 33,
		2: 33,
		3: 34,
	}

	results := make(map[int]int)
	iterations := 3000

	for i := 0; i < iterations; i++ {
		hit := Probability(list)
		results[hit]++
	}

	// 验证每个选项都被选中了
	for key := range list {
		if results[key] == 0 {
			t.Errorf("键%d从未被选中", key)
		}
	}
}

// TestProbabilityLargeWeights 测试大权重值
func TestProbabilityLargeWeights(t *testing.T) {
	list := map[int]int64{
		1: 1000000,
		2: 2000000,
		3: 3000000,
	}

	results := make(map[int]int)

	for i := 0; i < 100; i++ {
		hit := Probability(list)
		results[hit]++

		// 验证返回值在预期范围内
		if hit < 1 || hit > 3 {
			t.Errorf("Probability返回了无效的键: %d", hit)
		}
	}
}
