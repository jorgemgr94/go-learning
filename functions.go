package main

import "fmt"

func concat(s1 string, s2 string) string {
	return s1 + s2
}

func concatOneDeclaration(s1, s2 string) string {
	return s1 + s2
}

func test(s1 string, s2 string) {
	fmt.Println(concat(s1, s2))
	fmt.Println(concatOneDeclaration(s1, s2))
}

func functions() {
	test("Lane,", " happy birthday!")
	test("Elon,", " hope that Tesla thing works out")
	test("Go", " is fantastic")
}
