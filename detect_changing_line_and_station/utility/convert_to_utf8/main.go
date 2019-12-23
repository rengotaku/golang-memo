package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var (
	delimiter string
	encode    string
)

func init() {
	// -hオプション用文言
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
	 %s file(csv or tsv)
What is this?
   For encoding CSV(or TSV) was encoded Japanese(like as ShiftJIS) to UTF-8.
Options`+"\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&delimiter, "d", "\t", "set using delimiter.")
	flag.StringVar(&encode, "e", "SHIFTJIS", "choose from EUCJP, ISO2022JP, SHIFTJIS.")

	flag.Parse()
}

func main() {
	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	args := flag.Args()

	file, err := os.Open(args[0])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var reader *csv.Reader
	if encode == "SHIFTJIS" {
		reader = csv.NewReader(transform.NewReader(file, japanese.ShiftJIS.NewDecoder()))
	} else if encode == "ISO2022JP" {
		reader = csv.NewReader(transform.NewReader(file, japanese.ISO2022JP.NewDecoder()))
	} else if encode == "EUCJP" {
		reader = csv.NewReader(transform.NewReader(file, japanese.EUCJP.NewDecoder()))
	} else { // utf8
		reader = csv.NewReader(file)
	}

	// // reader.LazyQuotes = true       // a quote may appear in an unquoted field and a non-doubled quote may appear in a quoted field.
	// // reader.TrimLeading   = true // leading white    in a field is ignored.

	for {
		line, err := reader.Read()
		// 最終行でerrが返る
		if err != nil {
			break
		}
		fmt.Println(strings.Join(line, delimiter))
	}
}
