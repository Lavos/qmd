package main

import (
	"fmt"
	"os"
	"github.com/Lavos/qmd/lib"
)

func main() {
	q := lib.NewQmd("env")
	n := lib.NewNamespace("TEST")

	q.AppendEnv(lib.E("test", "12345"), lib.E("abc", "123%s", "abc"))
	q.AppendEnv(n.E("BANANAS", "VERY DELICIOUS"))

	q.RedirectFile("12345")
	q.Cmd.Stdout = os.Stdout

	err := q.Start()

	fmt.Printf("Start %s\n", err)

	err = q.Wait()

	fmt.Printf("Wait %s\n", err)
}
