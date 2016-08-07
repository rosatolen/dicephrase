package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	gen "github.com/rosatolen/diceware/generator"
)

func main() {
	var wc int
	defaultWc := gen.WordCount
	wcUsage := fmt.Sprintf("Specify the number of words you want. Must be greater than %v", gen.WordCount)
	flag.IntVar(&wc, "words", defaultWc, wcUsage)
	flag.IntVar(&wc, "w", defaultWc, "(words shorthand) "+wcUsage)

	var sep string
	defaultSep := gen.Separator
	sepUsage := "Specify the type of separator you want between the words. Cannot be alphanumeric and cannot exist as a word in the wordlist"
	flag.StringVar(&sep, "sep", defaultSep, sepUsage)
	flag.StringVar(&sep, "s", defaultSep, "(sep shorthand) "+sepUsage)

	curdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Error: cannot find absolute path to the executable %s\n", os.Args[0])
		return
	}
	var wl string
	defaultWl := curdir + "/wordlists/original-diceware.wordlist.asc"
	wlUsage := "Specify the path to the wordlist you prefer to use."
	flag.StringVar(&wl, "wordlist", defaultWl, wlUsage)
	flag.StringVar(&wl, "wl", defaultWl, "(wordlist shorthand) "+wlUsage)

	flag.Parse()

	p, err := gen.PassphraseWith(wl, sep, wc)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println(p)
	return
}
