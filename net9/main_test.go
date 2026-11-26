package main // 跟 net8 的 main.go 同一个包

import (
	"strings"
	"testing"
)

// ══════════════════════════════════════════════════════════
// TestGenerateToken：验证格式
// ══════════════════════════════════════════════════════════
func TestGenerateToken(t *testing.T) {
	token := generateToken()

	// ① 长度必须是 16
	if len(token) != 16 {
		t.Errorf("token 长度 = %d，期望 16（内容: %q）", len(token), token)
	}

	// ② 必须全是十六进制字符（0-9, a-f）
	const hexChars = "0123456789abcdef"
	for i, ch := range token {
		if !strings.ContainsRune(hexChars, ch) {
			t.Errorf("第 %d 位是非法字符: %q", i, ch)
		}
	}
}

// ══════════════════════════════════════════════════════════
// TestGenerateTokenUnique：验证唯一性
// ══════════════════════════════════════════════════════════
func TestGenerateTokenUnique(t *testing.T) {
	seen := make(map[string]bool) // 记录已经生成过的

	for i := 0; i < 1000; i++ { // 生成 1000 次
		token := generateToken()

		if seen[token] {
			// 重复了 → 严重问题
			t.Fatalf("第 %d 次生成了重复 token: %s", i, token)
			// Fatalf：立刻中断（后面的不用测了）
		}
		seen[token] = true
	}
}

func TestParseMove(t *testing.T) {
	cases := []struct {
		name    string
		line    string
		wantX   int
		wantY   int
		wantErr bool
	}{
		{"正常", "move 100 200", 100, 200, false},
		{"边界下限", "move 0 0", 0, 0, false},
		{"边界上限", "move 1000 1000", 1000, 1000, false},
		{"参数不足", "move 100", 0, 0, true},
		{"只有 move", "move", 0, 0, true},
		{"空字符串", "", 0, 0, true},
		{"x 不是数字", "move abc 200", 0, 0, true},
		{"y 不是数字", "move 100 abc", 0, 0, true},
		{"x 超上限", "move 1001 0", 0, 0, true},
		{"y 超上限", "move 0 1001", 0, 0, true},
		{"负数", "move -5 -10", 0, 0, true},
		{"参数过多", "move 100 200 300", 100, 200, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			x, y, err := parseMove(c.line)

			if (err != nil) != c.wantErr {
				t.Fatalf("parseMove(%q) 出错 = %v，期望出错 = %v", c.line, err, c.wantErr)
			}

			if !c.wantErr {
				if x != c.wantX || y != c.wantY {
					t.Errorf("parseMove(%q) = (%d,%d)，期望 (%d,%d)", c.line, x, y, c.wantX, c.wantY)
				}
			}
		})
	}
}
