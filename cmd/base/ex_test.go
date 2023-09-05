package ex

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "sum",
			args: args{
				x: 1,
				y: 2,
			},
			want: 3,
		},
		{
			name: "sumnegotive",
			args: args{
				x: -10,
				y: 0,
			},
			want: -10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: ErrEmptyFilename,
		},
		{
			name:    "not exist",
			args:    args{filename: "not_exist"},
			wantErr: ErrNotFound,
		},
		{
			name:    "success",
			args:    args{filename: "ex.go"},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OpenFile(tt.args.filename)
			assert.Error(t, err, ErrEmptyFilename)
			assert.Contains(t, err.Error(), tt.wantErr, "error must in conts")
		})
	}
}

func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Вызываем функцию, которую хотим замерить
		Fibonacci(20)
	}
}

func FuzzReverse(f *testing.F) {
	testcases := []string{"aeghaetha1ux0Gee9ui8", "jod4io0Ya4ieH}ei7ohl", "!ereihei{l,ie`cae0oR"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		rev, err1 := Reverse(orig)
		if err1 != nil {
			return
		}
		doubleRev, err2 := Reverse(rev)
		if err2 != nil {
			return
		}
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}

func TestDivide(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "divide",
			args: args{
				x: 10,
				y: 2,
			},
			want: 5,
		},
		// {
		// 	name: "dividezero",
		// 	args: args{
		// 		x: 10,
		// 		y: 0,
		// 	},
		// 	want: 0,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotZero(t, tt.args.y, "Делитель не может быть нулем")
			result := Divide(tt.args.x, tt.args.y)
			require.Equal(t, 5, result, "Деление работает неверно")
		})
	}
}
