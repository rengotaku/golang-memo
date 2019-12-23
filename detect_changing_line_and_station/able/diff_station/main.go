package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// -hオプション用文言
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s compare_line.tsv ct_mst_pref_eki_XXXX_utf8.tsv station_and_line.tsv`+"\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() < 3 {
		flag.Usage()
		return
	}
	args := flag.Args()

	ableLines, err := ReadAbleLine(args[0])
	if err != nil {
		panic(err)
	}
	// fmt.Println(ableLines)

	ableStations, err := ReadAbleStation(args[1], ableLines)
	if err != nil {
		panic(err)
	}
	// fmt.Println(len(ableStations))

	stations, err := ReadDoorStation(args[2])
	if err != nil {
		panic(err)
	}
	// fmt.Println(stations)

	ableStations = DiffLine(ableStations, stations)

	var hLine []string
	hLine = append(hLine, "able_station_id")
	hLine = append(hLine, "able_station_name")
	hLine = append(hLine, "able_line_id")
	hLine = append(hLine, "able_line_name")
	hLine = append(hLine, "able_ensentype")
	hLine = append(hLine, "door_station_id")
	hLine = append(hLine, "door_station_name")

	fmt.Println(strings.Join(hLine, "\t"))

	for _, ableStation := range ableStations {
		// fmt.Println(ableStation)
		var line []string
		line = append(line, "\""+ableStation.Ekicd+"\"")
		line = append(line, "\""+ableStation.Ekiname+"\"")
		if ableStation.AbleLine != nil {
			ableLine := ableStation.AbleLine
			line = append(line, "\""+strconv.Itoa(int(ableLine.AbleLineID))+"\"")
			line = append(line, "\""+ableLine.AbleLineName+"\"")
			line = append(line, "\""+strconv.Itoa(int(ableLine.AbleLineType))+"\"")
		} else {
			line = append(line, "")
			line = append(line, "")
			line = append(line, "")
		}
		if ableStation.DoorStation != nil {
			doorStation := ableStation.DoorStation
			line = append(line, strconv.Itoa(int(doorStation.ID)))
			line = append(line, "\""+doorStation.Name+"\"")
		} else {
			line = append(line, "")
			line = append(line, "")
		}

		if ableStation.DoorStation != nil {
			fmt.Println(strings.Join(line, "\t"))
		}
	}
}

// DiffLine はABLEとDOORの駅の差分をとる
func DiffLine(ableStations []AbleStation, doorStations DoorStations) (results []AbleStation) {
	var notSearches []AbleStation
	for _, ableStation := range ableStations {
		for lineI, doorStation := range doorStations {
			if ableStation.Equal(doorStation) {
				ableStation.DoorStation = &doorStation
				doorStations.Remove(lineI)
				results = append(results, ableStation)
				break
			}
		}
		if ableStation.DoorStation == nil {
			notSearches = append(notSearches, ableStation)
		}
	}
	// 片方の名称が長く合致しないパターンを考慮
	// e.g. 門司港レトロ観光線, 平成筑豊鉄道門司港レトロ観光線
	for _, notSearche := range notSearches {
		for lineI, doorStation := range doorStations {
			if notSearche.Include(doorStation) {
				notSearche.DoorStation = &doorStation
				doorStations.Remove(lineI)
				break
			}
		}
		results = append(results, notSearche)
	}

	return results
}

// Trim は先頭、後尾の全角、半角を削除
func Trim(name string) string {
	rep := regexp.MustCompile(`(^　+|　+$)`)
	name = rep.ReplaceAllString(name, "")
	return strings.TrimSpace(name)
}
