package main

import (
	"math"
	"regexp"
	"strings"

	lsdp "github.com/deltam/go-lsd-parametrized"
)

type NormalizedLineName struct {
	Name    string
	SubName string
}

// EseikatsuLine はいい生活のライン情報をまとめる
type EseikatsuLine struct {
	DoorLine           *DoorLine
	NormalizedLineName *NormalizedLineName

	LineID      int64  // 路線ID
	LineName    string // 路線名(路線－駅名)
	StationID   int64  // 駅ID
	StationName string // 駅名(路線－駅名)
	PrefName    string // 都道府県

	Score      int64  // 一致スコア
	Candidates string // 名前(ID)
}

// Lines is line テーブル 配列データ構造
type EseikatsuLines []EseikatsuLine

// Equal はレーベンシュタイン距離を用いてスコア化する
func (e *EseikatsuLine) CalScore(l DoorLine) int {
	wd := lsdp.Weights{Insert: 0.1, Delete: 1, Replace: 1}
	nd := lsdp.Normalized(wd)

	x1, x2 := Than(e.NormalizedLineName.Name, l.NormalizedLineName.Name)

	return round(100.0 - nd.Distance(x1, x2)*100.0)
}

// 丸め
func round(f float64) int {
	return int(math.Floor(f + .5))
}

// Equal は沿線名を比較して等しいか比べる
func (e *EseikatsuLine) Equal(l DoorLine) bool {
	return e.NormalizedLineName.Name == l.NormalizedLineName.Name ||
		e.NormalizedLineName.SubName == l.NormalizedLineName.Name ||
		e.NormalizedLineName.Name == l.NormalizedLineName.SubName ||
		e.NormalizedLineName.SubName == l.NormalizedLineName.Name
}

// NormalizeAbleLine はエイブルの沿線の正規化
func (e *EseikatsuLine) NormalizeLine() {
	name := SymbolToHan(e.LineName)

	// 成田線［我孫子－成田］
	rep := regexp.MustCompile(`(.+)(?:\[)(.+)(?:\])$`)
	if rep.MatchString(name) {
		name1 := rep.ReplaceAllString(name, "$1")
		// name2 := rep.ReplaceAllString(name, "$2")
		e.NormalizedLineName = &NormalizedLineName{e.RemoveChar(name1), ""}

		return
	}

	// 嵐電（京福）嵐山本線
	rep = regexp.MustCompile(`(.+)(?:\()(.+)(?:\))(.+)$`)
	if rep.MatchString(name) {
		name1 := rep.ReplaceAllString(name, "$1")
		name2 := rep.ReplaceAllString(name, "$2")
		name3 := rep.ReplaceAllString(name, "$3")
		e.NormalizedLineName = &NormalizedLineName{
			e.RemoveChar(name1) + e.RemoveChar(name3),
			e.RemoveChar(name2) + e.RemoveChar(name3),
		}

		return
	}

	// 東海道本線（東日本）
	rep = regexp.MustCompile(`(.+)(?:\()(.+)(?:\))$`)
	if rep.MatchString(name) {
		name1 := rep.ReplaceAllString(name, "$1")
		// name2 := rep.ReplaceAllString(name, "$2")
		e.NormalizedLineName = &NormalizedLineName{
			e.RemoveChar(name1),
			"",
			// e.RemoveChar(name2),
		}

		return
	}

	// 富山地鉄不二越・上滝線
	rep = regexp.MustCompile(`(.+)・(.+)`)
	if rep.MatchString(name) {
		name1 := rep.ReplaceAllString(name, "$1")
		// name2 := rep.ReplaceAllString(name, "$2")
		e.NormalizedLineName = &NormalizedLineName{
			e.RemoveChar(name1),
			"",
			// e.RemoveChar(name2),
		}

		return
	}

	e.NormalizedLineName = &NormalizedLineName{e.RemoveChar(name), ""}
}

func (l *EseikatsuLine) RemoveChar(rawName string) string {
	name := RemoveWord(rawName)
	name = RemovePrefix(name)
	return RemoveSufix(name)
}

// DoorLine はDOORの駅情報のデータ構造
type DoorLine struct {
	NormalizedLineName *NormalizedLineName

	ID              int64
	PrefID          int64
	Name            string
	CompanyID       int64
	TmpID           int64
	EkitanRosenCode int64
	Status          int64
}

// NormalizeDoorLine はDOORの沿線の正規化
func (l *DoorLine) NormalizeLine() {
	name := l.Name

	rep := regexp.MustCompile(`(.+)(?:＜)(.+)(?:＞)`)
	if rep.MatchString(name) {
		mainN := rep.ReplaceAllString(name, "$1")
		mainN2 := rep.ReplaceAllString(name, "$2")
		l.NormalizedLineName = &NormalizedLineName{l.RemoveChar(mainN) + l.RemoveChar(mainN2), ""}

		return
	}

	rep = regexp.MustCompile(`(.+)・(.+)`)
	if rep.MatchString(name) {
		mainN := rep.ReplaceAllString(name, "$1")
		otherN := rep.ReplaceAllString(name, "$2")
		l.NormalizedLineName = &NormalizedLineName{l.RemoveChar(mainN), l.RemoveChar(otherN)}

		return
	}

	l.NormalizedLineName = &NormalizedLineName{l.RemoveChar(name), ""}
}

func (l *DoorLine) RemoveChar(rawName string) string {
	name := RemoveWord(rawName)
	name = RemovePrefix(name)
	return RemoveSufix(name)
}

// Lines is line テーブル 配列データ構造
type DoorLines []DoorLine

// Remove is remove from array specifing position
func (l *DoorLines) Remove(index int) {
	res := DoorLines{}
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
	rep := regexp.MustCompile(`(鉄道線|鉄道|鐵道|鉄線|線)$`)
	return rep.ReplaceAllString(name, "")
}
