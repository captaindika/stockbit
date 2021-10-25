package main
import (
	"fmt"
	"strings"
	"unicode"
)

// func findFirstStringInBracket(str string) string {
//     if (len(str) > 0) {
//         indexFirstBracketFound:= strings.Index(str, "(")
// 		if indexFirstBracketFound >= 0 {
//             runes:= [] rune(str) 
// 			fmt.Println("runes: ",runes)
// 			fmt.Println(runes[indexFirstBracketFound: len(str)])
// 			wordsAfterFirstBracket:= string(runes[indexFirstBracketFound: len(str)])
// 			fmt.Println("wordsAfterFirstBracket: ", wordsAfterFirstBracket)
// 			indexClosingBracketFound:= strings.Index(wordsAfterFirstBracket, ")")
// 			if indexClosingBracketFound >= 0 {
//                 runes:= [] rune(wordsAfterFirstBracket)
// 				return string(runes[1: indexClosingBracketFound])
// 			} else {
//                 return ""
//             }
// 		} else {
// 			return ""
// 		}
//     } else {
//         return ""
//     }
//     return ""
// }
func findFirstStringInBracket(s string)string {
	return (strings.TrimFunc(s, func(r rune)bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r)}))
}
func main () {
	fmt.Println(findFirstStringInBracket("(halo aku)"))
}