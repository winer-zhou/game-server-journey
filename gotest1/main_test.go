package main

import "testing"

// TestAdd：测试 Add 函数
// 名字必须以 Test 开头
// 参数必须是 *testing.T
func TestAd(t *testing.T) {
	got := Add(1, 2)
	want := 3

	if got != want {
		t.Errorf("Add(1,2) = %d,期望 %d", got, want)

	}
	// Errorf = 记录错误并继续（不中断）
	// Fatalf = 记录错误并立即中断
}

func TestIsEven(t *testing.T) {
	//情况 1：4 是偶数
	if !IsEven(4) {
		t.Errorf("4应该是偶数")
	}
	//情况 2：3 不是偶数
	if IsEven(3) {
		t.Error("3不该是偶数")

	}
}

// ══════════════════════════════════════════════════════════
// TestAddTable：表驱动测试
// 把所有测试用例写成一张"表"，循环跑
// ══════════════════════════════════════════════════════════
func TestAddTable(t *testing.T) {

	// ① 定义"表格"：每个元素是一组测试数据
	cases := []struct {
		name string // 用例名（显示在输出里）
		a    int    // 输入 1
		b    int    // 输入 2
		want int    // 期望结果
	}{
		{"正数", 1, 2, 3},
		{"负数", -1, -2, -3},
		{"零", 0, 0, 0},
		{"一正一负", -5, 5, 0},
		{"大数", 1000000, 2000000, 3000000},
	}

	// ② 循环跑每一组
	for _, c := range cases {
		// t.Run(名字, 函数)：跑一个"子测试"
		// 好处：输出里能看到每个用例的独立结果
		t.Run(c.name, func(t *testing.T) {
			// := 在 if 里声明并判断（Go 常用写法）
			if got := Add(c.a, c.b); got != c.want {
				t.Errorf("Add(%d,%d) = %d,期望 %d", c.a, c.b, got, c.want)
			}
		})
	}
}
