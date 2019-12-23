package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// -hオプション用文言
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s compare_company.tsv ct_mst_pref_ensen_XXXX.tsv line.tsv`+"\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() < 3 {
		flag.Usage()
		return
	}
	args := flag.Args()

	ableCompanies, err := ReadAbleCompany(args[0])
	if err != nil {
		panic(err)
	}
	// fmt.Println(ableCompanies)

	ableEnsens, err := ReadAbleEnsen(args[1], ableCompanies)
	if err != nil {
		panic(err)
	}
	// fmt.Println(ableEnsens)

	lines, err := ReadDoorEki(args[2])
	if err != nil {
		panic(err)
	}
	// fmt.Println(lines)

	ableEnsens = DiffLine(ableEnsens, lines)

	for _, ableEnsen := range ableEnsens {
		var line []string
		line = append(line, ableEnsen.Ensencd)
		line = append(line, "\""+ableEnsen.Ensenname+"\"")
		line = append(line, ableEnsen.Ensentype)

		if ableEnsen.Line != nil {
			lineID := strconv.Itoa(int(ableEnsen.Line.ID))

			line = append(line, lineID)
			line = append(line, "\""+ableEnsen.Line.Name+"\"")
		} else {
			line = append(line, "")
			line = append(line, "\"\"")
		}

		fmt.Println(strings.Join(line, "\t"))
	}
}

// DiffLine is
// FIXME: ポインタを渡して処理させる方法が分からない
func DiffLine(ctMstPrefEnsens []CtMstPrefEnsen, lines Lines) (results []CtMstPrefEnsen) {
	var notSearches []CtMstPrefEnsen
	for _, ctMstPrefEnsen := range ctMstPrefEnsens {
		for _, line := range lines {
			if ctMstPrefEnsen.Equal(line) {
				ctMstPrefEnsen.Line = &line
				// lines.Remove(lineI)
				results = append(results, ctMstPrefEnsen)
				break
			}
		}
		if ctMstPrefEnsen.Line == nil {
			notSearches = append(notSearches, ctMstPrefEnsen)
		}
	}
	// 片方の名称が長く合致しないパターンを考慮
	// e.g. 門司港レトロ観光線, 平成筑豊鉄道門司港レトロ観光線
	for _, notSearche := range notSearches {
		for _, line := range lines {
			if notSearche.Include(line) {
				notSearche.Line = &line
				// lines.Remove(lineI)
				break
			}
		}
		results = append(results, notSearche)
	}

	return results
}
