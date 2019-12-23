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
   %s mst_railway_company_XXXX.tsv s_company.tsv`+"\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		return
	}
	args := flag.Args()

	ableCompanies, err := ReadAbleCompany(args[0])
	if err != nil {
		panic(err)
	}

	doorCompanies, err := ReadDoorCompany(args[1])
	if err != nil {
		panic(err)
	}

	ableCompanies = DiffLine(ableCompanies, doorCompanies)

	for _, ableCompany := range ableCompanies {
		var line []string
		line = append(line, strconv.Itoa(int(ableCompany.RailwaycoNo)))
		line = append(line, "\""+ableCompany.RailwaycoName+"\"")

		if ableCompany.DoorCompany != nil {
			doorCompany := ableCompany.DoorCompany
			line = append(line, strconv.Itoa(int(doorCompany.ID)))
			line = append(line, "\""+doorCompany.Name+"\"")
		} else {
			line = append(line, "")
			line = append(line, "")
		}
		fmt.Println(strings.Join(line, "\t"))
	}
}

// DiffLine で差分を取り合致してそうな鉄道会社を組み合わせる
func DiffLine(ableCompanies []AbleCompany, doorCompanies DoorCompanies) (results []AbleCompany) {
	var notSearches []AbleCompany
	for _, ableCompany := range ableCompanies {
		for lineI, doorCompany := range doorCompanies {
			if ableCompany.NormalizedRailwaycoName == doorCompany.NormalizedName ||
				ableCompany.NormalizedRailwaycoName == doorCompany.NormalizedEkitanName {
				ableCompany.DoorCompany = &doorCompany
				doorCompanies.Remove(lineI)
				results = append(results, ableCompany)
				break
			}
		}
		if ableCompany.DoorCompany == nil {
			notSearches = append(notSearches, ableCompany)
		}
	}
	// 片方の名称が長く合致しないパターンを考慮
	// e.g. 門司港レトロ観光線, 平成筑豊鉄道門司港レトロ観光線
	for _, notSearche := range notSearches {
		for lineI, doorCompany := range doorCompanies {

			var x1 string
			var x2 string
			if len(notSearche.NormalizedRailwaycoName) > len(doorCompany.NormalizedName) {
				x1 = notSearche.NormalizedRailwaycoName
				x2 = doorCompany.NormalizedName
			} else {
				x1 = doorCompany.NormalizedName
				x2 = notSearche.NormalizedRailwaycoName
			}

			if strings.Index(x1, x2) > -1 {
				notSearche.DoorCompany = &doorCompany
				doorCompanies.Remove(lineI)
				break
			}
		}
		results = append(results, notSearche)
	}

	return results
}
