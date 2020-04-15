package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var yesFlg bool

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(string(bytes))
}

func init() {
	// -hオプション用文言
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
  %s [original folder]
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	// HACK: not working
	flag.BoolVar(&yesFlg, "y", true, "not confirm dir name.")
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	flag.Parse()
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		return
	}
	args := flag.Args()
	genDir := args[0]

	if _, err := os.Stat(genDir); os.IsNotExist(err) {
		panic(fmt.Errorf("Not exist folder: %s", genDir))
	}

	dir, _ := os.Getwd()
	log.Print(fmt.Sprint(`PWD: ` + dir))

	now := time.Now().Format("20060102")
	log.Print(fmt.Sprint(`DATE: ` + now))

	cDir := dir + "/" + now
	log.Print(fmt.Sprint(`
Create dir will be ` + cDir + `
Continue? Yy):`))

	reader := bufio.NewReader(os.Stdin) //create new reader, assuming bufio imported
	inp, _ := reader.ReadString('\n')
	fmt.Print(yesFlg)
	if !yesFlg {
		linp := strings.ToLower(strings.ReplaceAll(inp, "\n", ""))

		if linp != "y" {
			os.Exit(1)
		}
	}

	if err := CopyDir(genDir, cDir); err != nil {
		panic(err)
	}

	todoF, err := os.Create(cDir + "/" + now + "todo.md")
	if err != nil {
		log.Fatal(err)
	}
	todoF.Close()

	noteF, err := os.Create(cDir + "/" + now + ".md")
	if err != nil {
		log.Fatal(err)
	}
	noteF.Close()

	log.Print(fmt.Sprint(`Create dir is done!`))

	os.Exit(0)
}
