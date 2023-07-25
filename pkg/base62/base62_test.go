package base62

import (
	"testing"
)

func TestInt2Srting(t *testing.T) {
	tests := []struct {
		name string
		seq  uint64
		want string
	}{
		// TODO: Add test cases.
		{name: "第一个测试用例", seq: 0, want: "0"},
		{name: "第二个测试用例", seq: 2, want: "2"},
		{name: "第三个测试用例", seq: 62, want: "10"},
		{name: "第四个测试用例", seq: 6347, want: "1En"},
		{name: "第五个测试用例", seq: 63, want: "11"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int2Srting(tt.seq); got != tt.want {
				t.Errorf("Int2Srting() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString2Int(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		wantRes uint64
	}{
		// TODO: Add test cases.
		{name: "第一个测试用例", s: "0", wantRes: 0},
		{name: "第二个测试用例", s: "2", wantRes: 2},
		{name: "第三个测试用例", s: "10", wantRes: 62},
		{name: "第四个测试用例", s: "1En", wantRes: 6347},
		{name: "第五个测试用例", s: "11", wantRes: 63},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := String2Int(tt.s); gotRes != tt.wantRes {
				t.Errorf("String2Int() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
