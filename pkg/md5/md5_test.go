package md5

import "testing"

func TestSum(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "正确的示例",
			args: args{data: []byte("27933")},
			want: "dccf189cbe63472d0f4f5b00facfd2e1",
		},
		{
			name: "正确的示例",
			args: args{data: []byte("hello")},
			want: "5d41402abc4b2a76b9719d911017c592",
		},
		{
			name: "错误的示例",
			args: args{data: []byte("27933")},
			want: "dccf189cbe63472d0f4f5b00facfd2e1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.data); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}
