package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bytedance/sonic"
)

type Book struct {
	Name   string `json:"title"`
	Price  float64
	Tags   []string
	Press  string
	Author People
}

type People struct {
	Name    string
	Age     int
	School  string
	Company string
	Title   string
}

var (
	people = People{
		Name:    "张三",
		Age:     18,
		School:  "清华大学",
		Company: "大乔乔教育",
		Title:   "开发工程师",
	}
	book = Book{
		Name:   "高性能golang",
		Price:  58.0,
		Tags:   []string{"golang", "编程", "计算机"},
		Press:  "机械工业出版社",
		Author: people,
	}
)

func TestStdJson(t *testing.T) {
	// 用标准库进行json序列化。输出是[]byte
	if bs, err := json.Marshal(book); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(string(bs))

		// 用标准库进行json反序列化。输入是[]byte
		var book2 Book
		if err := json.Unmarshal(bs, &book2); err != nil {
			fmt.Println(err)
			t.Fail()
		} else {
			fmt.Printf("%+v\n", book2)
		}
	}
}

func TestSonic(t *testing.T) {
	// 用Sonic库进行json序列化。输出是[]byte
	if bs, err := sonic.Marshal(book); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(string(bs))

		// 用Sonic库进行json反序列化。输入是[]byte
		var book2 Book
		if err := sonic.Unmarshal(bs, &book2); err != nil {
			fmt.Println(err)
			t.Fail()
		} else {
			fmt.Printf("%+v\n", book2)
		}
	}
}

func BenchmarkStdJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bs, _ := json.Marshal(book)
		var book2 Book
		json.Unmarshal(bs, &book2)
	}
}

func BenchmarkSonic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bs, _ := sonic.Marshal(book)
		var book2 Book
		sonic.Unmarshal(bs, &book2)
	}
}

// go test -v ./test -run=TestStdJson -count=1
// go test -v ./test -run=TestSonic -count=1

// go test ./test -bench=BenchmarkStdJson -run=none -count=1 -benchmem -benchtime=2s
// go test ./test -bench=BenchmarkSonic -run=none -count=1 -benchmem -benchtime=2s
