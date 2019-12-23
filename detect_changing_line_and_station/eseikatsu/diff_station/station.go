package main

import (
	"math"
	"regexp"
	"strings"
)

type NormalizedStationName struct {
	Name    string
	SubName string
}

// EseikatsuLine はいい生活のライン情報をまとめる
type EseikatsuLine struct {
	EseikatsuLineID   int64
	EseikatsuLineName string
	DoorLineID        int64
	DoorLineName      string
}

// Lines is line テーブル 配列データ構造
type EseikatsuLines []EseikatsuLine

// EseikatsuLine はいい生活のライン情報をまとめる
type EseikatsuStation struct {
	EseikatsuLine         *EseikatsuLine
	DoorStation           *DoorStation
	NormalizedStationName *NormalizedStationName

	LineID      int64  // 路線ID
	LineName    string // 路線名(路線－駅名)
	StationID   int64  // 駅ID
	StationName string // 駅名(路線－駅名)
	PrefName    string // 都道府県

	Candidates string // 名前(ID)
}

// Lines is line テーブル 配列データ構造
type EseikatsuStations []EseikatsuStation

// Equal はレーベンシュタイン距離を用いてスコア化する
func (e *EseikatsuStation) CalScore(l DoorStation) int {
	x1, x2 := Than(e.NormalizedStationName.Name, l.NormalizedStationName.Name)

	return round(100.0 - nd.Distance(x1, x2)*100.0)
}

// Equal は沿線名を比較して等しいか比べる
func (e *EseikatsuStation) Equal(l DoorStation) bool {
	if e.EseikatsuLine.DoorLineID != l.LineID {
		return false
	}

	return e.NormalizedStationName.Name == l.NormalizedStationName.Name ||
		e.NormalizedStationName.SubName == l.NormalizedStationName.Name ||
		e.NormalizedStationName.Name == l.NormalizedStationName.SubName ||
		e.NormalizedStationName.SubName == l.NormalizedStationName.Name
}

// NormalizeLine 沿線の正規化
func (e *EseikatsuStation) NormalizeStation() {
	name := SymbolToHan(e.StationName)
	// TODO: 使用するパターンを検討する
	// // 成田線［我孫子－成田］
	// rep := regexp.MustCompile(`(.+)(?:\[)(.+)(?:\])$`)
	// if rep.MatchString(name) {
	// 	name1 := rep.ReplaceAllString(name, "$1")
	// 	// name2 := rep.ReplaceAllString(name, "$2")
	// 	e.NormalizedLineName = &NormalizedLineName{e.RemoveChar(name1), ""}

	// 	return
	// }

	// // 嵐電（京福）嵐山本線
	// rep = regexp.MustCompile(`(.+)(?:\()(.+)(?:\))(.+)$`)
	// if rep.MatchString(name) {
	// 	name1 := rep.ReplaceAllString(name, "$1")
	// 	name2 := rep.ReplaceAllString(name, "$2")
	// 	name3 := rep.ReplaceAllString(name, "$3")
	// 	e.NormalizedLineName = &NormalizedLineName{
	// 		e.RemoveChar(name1) + e.RemoveChar(name3),
	// 		e.RemoveChar(name2) + e.RemoveChar(name3),
	// 	}

	// 	return
	// }

	// 小林（宮崎）
	reg := regexp.MustCompile(`(.+)(?:\()(.+)(?:\))$`)
	if reg.MatchString(name) {
		name1 := reg.ReplaceAllString(name, "$1")
		// name2 := reg.ReplaceAllString(name, "$2")
		e.NormalizedStationName = &NormalizedStationName{
			e.RemoveChar(name1),
			"",
			// e.RemoveChar(name2),
		}

		return
	}

	// // 富山地鉄不二越・上滝線
	// rep = regexp.MustCompile(`(.+)・(.+)`)
	// if rep.MatchString(name) {
	// 	name1 := rep.ReplaceAllString(name, "$1")
	// 	// name2 := rep.ReplaceAllString(name, "$2")
	// 	e.NormalizedLineName = &NormalizedLineName{
	// 		e.RemoveChar(name1),
	// 		"",
	// 		// e.RemoveChar(name2),
	// 	}

	// 	return
	// }

	e.NormalizedStationName = &NormalizedStationName{e.RemoveChar(name), ""}
}

func (l *EseikatsuStation) RemoveChar(rawName string) string {
	name := Trim(rawName)
	name = RemoveWord(name)
	name = RemovePrefix(name)
	return RemoveSufix(name)
}

// DoorStation はDOORの駅のデータ構造
type DoorStation struct {
	EseikatsuStation      *EseikatsuStation
	NormalizedStationName *NormalizedStationName

	ID     int64
	Name   string
	LineID int64

	Score int64 // 一致スコア
}

// Lines is line テーブル 配列データ構造
type DoorStations []DoorStation

// NormalizeDoorStation はDOORの駅の正規化
func (d *DoorStation) NormalizeStation() {
	name := SymbolToHan(d.Name)

	rep := regexp.MustCompile(`(.+)(?:\()(.+)(?:\))`)
	if rep.MatchString(name) {
		name1 := rep.ReplaceAllString(name, "$1")
		// name2 := rep.ReplaceAllString(name, "$2")

		d.NormalizedStationName = &NormalizedStationName{
			d.RemoveChar(name1),
			"",
			// d.RemoveChar(name2),
		}

		return
	}

	rep = regexp.MustCompile(`(.+)・(.+)`)
	if rep.MatchString(name) {
		name1 := rep.ReplaceAllString(name, "$1")
		// name2 := rep.ReplaceAllString(name, "$2")
		d.NormalizedStationName = &NormalizedStationName{
			d.RemoveChar(name1),
			"",
			// d.RemoveChar(name2),
		}

		return
	}

	d.NormalizedStationName = &NormalizedStationName{d.RemoveChar(name), ""}
}

// RemoveChar は指定の文字を削除する
func (d *DoorStation) RemoveChar(rawName string) string {
	name := Trim(rawName)
	name = RemoveWord(name)
	name = RemovePrefix(name)
	return RemoveSufix(name)
}

// Remove is remove from array specifing position
func (l *DoorStations) Remove(index int) {
	res := DoorStations{}
	for i, v := range *l {
		if i == index {
			continue
		}
		res = append(res, v)
	}
	*l = res
}

// SymbolToHan は特定の記号を半角にする
func SymbolToHan(str string) string {
	rel := strings.Replace(str, "（", "(", -1)
	rel = strings.Replace(rel, "）", ")", -1)
	rel = strings.Replace(rel, "＜", "<", -1)
	rel = strings.Replace(rel, "＞", ">", -1)
	rel = strings.Replace(rel, "［", "[", -1)
	rel = strings.Replace(rel, "］", "]", -1)

	return rel
}

func RemoveWord(name string) string {
	rep := regexp.MustCompile(`(都市)`)
	return rep.ReplaceAllString(name, "")
}

func RemovePrefix(name string) string {
	rep := regexp.MustCompile(`^(ＪＲ|JR)`)
	return rep.ReplaceAllString(name, "")
}

func RemoveSufix(name string) string {
	rep := regexp.MustCompile(`(鉄道線|鉄道|鐵道|線)$`)
	return rep.ReplaceAllString(name, "")
}

// 丸め
func round(f float64) int {
	return int(math.Floor(f + .5))
}
