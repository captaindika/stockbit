package main

import (
	"fmt"
	"sort"
	"strings"
)

var (
	dummyDict = []string{"kita", "atik", "tika", "aku", "kia", "makan", "kua"}
)

func Anagram() {
	list := make(map[string][]string)

	for _, word := range dummyDict {
		key := sortStr(word)
		list[key] = append(list[key], word)
	}
	for _, words := range list {
		fmt.Println(words)
		
	}
}

func sortStr(k string) string {
	s := strings.Split(k, "")
	sort.Strings(s)

	return strings.Join(s, "")
}


func main() {
	Anagram()
}
