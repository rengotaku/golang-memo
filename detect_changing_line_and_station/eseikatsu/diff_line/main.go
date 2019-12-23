package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	minScore   int = 60
	pickMaxNum int = 3
)

type MatchScore struct {
	DoorLine DoorLine

	Score int64
}

func main() {
	// -hオプション用文言
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s XXXX_ENSEN_MASTER_utf8.txt line.tsv`+"\n", os.Args[0], os.Args[0])
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

	doorLines, err := ReadDoorLine(args[1])
	if err != nil {
		panic(err)
	}
	// fmt.Println(doorLines)

	eseikatsuLines = DiffLine(eseikatsuLines, doorLines)

	var hLine []string
	hLine = append(hLine, "eseikatsu_line_id")
	hLine = append(hLine, "eseikatsu_line_name")
	hLine = append(hLine, "door_line_id")
	hLine = append(hLine, "door_line_name")
	hLine = append(hLine, "score")
	hLine = append(hLine, "candidates")

	fmt.Println(strings.Join(hLine, "\t"))

	for _, eseikatsuLine := range eseikatsuLines {
		var line []string
		line = append(line, strconv.Itoa(int(eseikatsuLine.LineID)))
		line = append(line, "\""+eseikatsuLine.LineName+"\"")

		if eseikatsuLine.DoorLine != nil {
			line = append(line, strconv.Itoa(int(eseikatsuLine.DoorLine.ID)))
			line = append(line, "\""+eseikatsuLine.DoorLine.Name+"\"")
			line = append(line, strconv.Itoa(int(eseikatsuLine.Score)))
			line = append(line, "\""+eseikatsuLine.Candidates+"\"")
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
// func DiffLine(eseikatsuLines EseikatsuLines, doorLines DoorLines) (results EseikatsuLines) {
// 	for _, eseikatsuLine := range eseikatsuLines {
// 		for _, doorLine := range doorLines {
// 			// 完全一致
// 			if eseikatsuLine.Equal(doorLine) {
// 				eseikatsuLine.Score = 100
// 				eseikatsuLine.DoorLine = &doorLine
// 				results = append(results, eseikatsuLine)
// 				break
// 			}
// 			// 片方の名称が長く合致しないパターンを考慮
// 			// if eseikatsuLine.Include(doorLine) {
// 			// 	score := eseikatsuLine.CalScore(doorLine)
// 			// 	// 一致率が低い
// 			// 	if score < minScore {
// 			// 		// fmt.Println(eseikatsuLine.LineName, doorLine.Name)
// 			// 		continue
// 			// 	}
// 			// 	eseikatsuLine.Score = int64(score)
// 			// 	eseikatsuLine.DoorLine = &doorLine
// 			// 	results = append(results, eseikatsuLine)
// 			// 	break
// 			// }
// 		}

// 		if eseikatsuLine.DoorLine != nil {
// 			continue
// 		}

// 		var matches []map[string]string
// 		// 一致しそうな沿線を検索
// 		for _, doorLine := range doorLines {
// 			score := eseikatsuLine.CalScore(doorLine)
// 			if score > 0 {
// 				m := map[string]string{
// 					"id":    strconv.Itoa(int(doorLine.ID)),
// 					"name":  doorLine.Name,
// 					"score": strconv.Itoa(int(score)),
// 				}
// 				matches = append(matches, m)
// 			}
// 		}

// 		// 降順
// 		sort.SliceStable(matches, func(i, j int) bool {
// 			x1, _ := strconv.Atoi(matches[i]["score"])
// 			x2, _ := strconv.Atoi(matches[j]["score"])
// 			return x1 > x2
// 		})
// 		var chars []string
// 		// 指定数分の上位をピック
// 		for _, match := range matches[0:pickItemNum] {
// 			chars = append(chars, match["name"]+"("+match["id"]+")/"+match["score"])
// 		}

// 		eseikatsuLine.Candidates = strings.Join(chars, ",")

// 		results = append(results, eseikatsuLine)
// 	}

// 	return results
// }

func DiffLine(eseikatsuLines EseikatsuLines, doorLines DoorLines) (results EseikatsuLines) {
	for _, eseikatsuLine := range eseikatsuLines {
		var matches []MatchScore
		for _, doorLine := range doorLines {
			// 完全一致
			if eseikatsuLine.Equal(doorLine) {
				eseikatsuLine.Score = 100
				eseikatsuLine.DoorLine = &doorLine
				results = append(results, eseikatsuLine)
				break
			}
			// 片方の名称が長く合致しないパターンを考慮
			score := eseikatsuLine.CalScore(doorLine)
			if score > 0 {
				m := MatchScore{doorLine, int64(score)}
				matches = append(matches, m)
			}
		}
		if eseikatsuLine.DoorLine != nil {
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
				chars = append(chars, match.DoorLine.Name+"("+strconv.Itoa(int(match.DoorLine.ID))+")/"+strconv.Itoa(int(match.Score)))
			}

			top := matches[0]
			eseikatsuLine.Score = top.Score
			eseikatsuLine.DoorLine = &top.DoorLine
			eseikatsuLine.Candidates = strings.Join(chars, ",")
		}

		results = append(results, eseikatsuLine)
	}

	return results
}

func Than(n1 string, n2 string) (x1 string, x2 string) {
	if len(n1) > len(n2) {
		x1 = n1
		x2 = n2
	} else {
		x1 = n2
		x2 = n1
	}
	return
}
