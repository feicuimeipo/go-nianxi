package main

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

func main() {
	s := "hello 中国"
	var len int
	var buf [utf8.UTFMax]byte
	for i, r := range s {
		rl := utf8.RuneLen(r)
		si := rl + i
		copy(buf[:], s[i:si])
		len += rl
		fmt.Sprintf("%2d:%q;codepoint:%#6x; encode bytes: %#v\n", i, r, r, buf[:rl])
	}
	fmt.Printf("len=" + strconv.Itoa(len))

}
