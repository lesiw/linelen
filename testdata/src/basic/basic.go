package basic

import "fmt"

func main() {
	// This is a normal line that should not trigger the analyzer
	fmt.Println("Hello, world!")

	// This line is intentionally very long and should trigger the line length analyzer when checking // want "line is 157 characters long, exceeds 79 limit"
	fmt.Println("This is a very long line that exceeds the default 79 character limit and should be reported by the analyzer") // want "line is 182 characters long, exceeds 79 limit"
}
