package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	lsdp "github.com/deltam/go-lsd-parametrized"
)

const (
	minScore   int = 99 // 許容スコア
	pickMaxNum int = 3  // 候補数
)

var (
	wd lsdp.Weights
	nd lsdp.DistanceMeasurer
)

type MatchScore struct {
	EseikatsuStation EseikatsuStation

	Score int64
}

func init() {
	wd = lsdp.Weights{Insert: 0.1, Delete: 1, Replace: 1}
	nd = lsdp.Normalized(wd)
}

func main() {
	// -hオプション用文言
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s compare_line.tsv XXXX_ENSEN_MASTER_utf8.txt station_and_line.tsv`+"\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		return
	}
	args := flag.Args()

	eseikatsuLines, err := ReadEseikatsuLine(args[0])
	if err != nil {
		panic(err)
	}
	// fmt.Println(eseikatsuLines)

	eseikatsuStations, err := ReadEseikatsuStation(args[1], eseikatsuLines)
	if err != nil {
		panic(err)
	}
	// fmt.Println(eseikatsuStations)

	doorStations, err := ReadDoorStation(args[2])
	if err != nil {
		panic(err)
	}
	// fmt.Println(doorStations)

	doorStations = DiffLine(eseikatsuStations, doorStations)
	// fmt.Println(eseikatsuStations)

	var hLine []string
	hLine = append(hLine, "eseikatsu_line_id")
	hLine = append(hLine, "eseikatsu_station_id")
	hLine = append(hLine, "door_station_id")
	hLine = append(hLine, "station_name")
	hLine = append(hLine, "score")
	hLine = append(hLine, "candidates")

	fmt.Println(strings.Join(hLine, "\t"))

	for _, doorStation := range doorStations {
		eseikatsuStation := doorStation.EseikatsuStation

		var line []string
		line = append(line, strconv.Itoa(int(eseikatsuStation.StationID)))
		line = append(line, "\""+eseikatsuStation.StationName+"\"")

		if doorStation.EseikatsuStation != nil {
			eseikatsuStation := doorStation.EseikatsuStation
			line = append(line, strconv.Itoa(int(doorStation.ID)))
			line = append(line, "\""+doorStation.Name+"\"")
			line = append(line, strconv.Itoa(int(doorStation.Score)))
			line = append(line, "\""+eseikatsuStation.Candidates+"\"")
		} else {
			line = append(line, "")
			line = append(line, "\"\"")
			line = append(line, "")
			line = append(line, "\"\"")
		}

		fmt.Println(strings.Join(line, "\t"))
	}
}

// DiffLine is
func DiffLine(eseikatsuStations EseikatsuStations, doorStations DoorStations) (results DoorStations) {
	// for _, eseikatsuStation := range eseikatsuStations {
	// 	fmt.Println("ESEIKATSU_NAME: ", eseikatsuStation.StationName)
	// }

	for _, doorStation := range doorStations {
		var matches []MatchScore
		for _, eseikatsuStation := range eseikatsuStations {
			// FIXME: なぜか関係ない値で合致してしまう。
			// 完全一致でもいいような比較がなぜかスコア判定されている原因を調査する
			if strings.Compare(doorStation.Name, "弘明寺") > -1 && strings.Compare(eseikatsuStation.StationName, "弘明寺") > -1 {
				fmt.Println("debug")
			}

			// 完全一致
			if eseikatsuStation.Equal(doorStation) {
				doorStation.EseikatsuStation = &eseikatsuStation
				doorStation.Score = 100
				results = append(results, doorStation)
				break
			}
			// 片方の名称が長く合致しないパターンを考慮
			score := eseikatsuStation.CalScore(doorStation)
			if score > 0 {
				m := MatchScore{eseikatsuStation, int64(score)}
				matches = append(matches, m)
			}
		}
		if doorStation.EseikatsuStation != nil {
			continue
		}

		if len(matches) > 0 {
			pickNum := pickMaxNum + 1
			if len(matches) < pickNum {
				pickNum = len(matches)
			}

			// 降順
			sort.SliceStable(matches, func(i, j int) bool {
				return matches[i].Score > matches[j].Score
			})

			var chars []string
			// 指定数分の上位をピック
			for _, match := range matches[1:pickNum] {
				chars = append(chars, match.EseikatsuStation.StationName+"("+strconv.Itoa(int(match.EseikatsuStation.StationID))+")/"+strconv.Itoa(int(match.Score)))
			}

			top := matches[0]

			doorStation.Score = top.Score
			doorStation.EseikatsuStation = &top.EseikatsuStation
			doorStation.EseikatsuStation.Candidates = strings.Join(chars, ",")
		}

		results = append(results, doorStation)
	}

	return
}

func Than(n1 string, n2 string) (string, string) {
	if len(n1) > len(n2) {
		return n1, n2
	}

	return n2, n1
}

// Trim は先頭、後尾の全角、半角を削除
func Trim(name string) string {
	rep := regexp.MustCompile(`(^　+|　+$)`)
	name = rep.ReplaceAllString(name, "")
	return strings.TrimSpace(name)
}
