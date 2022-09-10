package id

import (
	"testing"
)

func TestGenShortId(t *testing.T) {
	shortId := GenShortId()
	if shortId == "" {
		t.Error("GenShortId failed!")
	}

	t.Log("GenShortId test pass")
}

func BenchmarkGenShortId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenShortId()
	}
}

func BenchmarkGenShortIdTimeConsuming(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数

	shortId := GenShortId()
	if shortId == "" {
		b.Error("error")
	}

	b.StartTimer() //重新开始时间

	for i := 0; i < b.N; i++ {
		GenShortId()
	}
}
