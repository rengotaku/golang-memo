package main

import (
	"regexp"
	"strings"
)

// AbleLine はエイブルの沿線の構造体
type AbleLine struct {
	AbleLineID   int64
	AbleLineName string
	AbleLineType int64
	DoorLineID   int64
	DoorLineName string
}

// AbleLines is AbleLine 配列データ構造
type AbleLines []AbleLine

// Remove is remove from array specifing position
func (l *AbleLines) Remove(index int) {
	res := AbleLines{}
	for i, v := range *l {
		if i == index {
			continue
		}
		res = append(res, v)
	}
	*l = res
}

// NormalizedStationName は正規化した比較用の駅名の構造体
type NormalizedStationName struct {
	Name    string
	SubName string
}

// AbleStation is ct_mst_pref_eki データ構造
type AbleStation struct {
	AbleLine              *AbleLine
	DoorStation           *DoorStation
	NormalizedStationName *NormalizedStationName

	Prefcd      string
	Ensencd     string
	Ekicd       string
	Ekiname     string
	SortKey     int64
	StopFlg     string
	BukkenCnt   int64
	ShopCnt     int64
	Ekino       int64
	StartekiFlg string
	Ekiseq      int64
	NyBukkenCnt int64
}

// Equal は駅名を比較して等しいか比べる
func (e *AbleStation) Equal(l DoorStation) bool {
	if e.AbleLine.DoorLineID != int64(l.LineID) {
		return false
	}

	return e.NormalizedStationName.Name == l.NormalizedStationName.Name ||
		e.NormalizedStationName.SubName == l.NormalizedStationName.Name ||
		e.NormalizedStationName.Name == l.NormalizedStationName.SubName ||
		e.NormalizedStationName.SubName == l.NormalizedStationName.Name
}

// Include は駅名をどちらかが含んでいないか検証
func (e *AbleStation) Include(l DoorStation) bool {
	if e.AbleLine.DoorLineID != int64(l.LineID) {
		return false
	}

	return (e.Compare(e.NormalizedStationName.Name, l.NormalizedStationName.Name) ||
		e.Compare(e.NormalizedStationName.SubName, l.NormalizedStationName.Name) ||
		e.Compare(e.NormalizedStationName.Name, l.NormalizedStationName.SubName) ||
		e.Compare(e.NormalizedStationName.SubName, l.NormalizedStationName.Name))
}

// Compare は沿線名をどちらかが含んでいないか検証
func (e *AbleStation) Compare(n1 string, n2 string) bool {
	if n1 == "" || n2 == "" {
		return false
	}

	var x1 string
	var x2 string
	if len(n1) > len(n2) {
		x1 = n1
		x2 = n2
	} else {
		x1 = n2
		x2 = n1
	}

	return strings.Contains(x1, x2)
}

// NormalizeAbleStation はエイブルの駅の正規化
func (e *AbleStation) NormalizeAbleStation() {
	name := SymbolToHan(e.Ekiname)

	rep := regexp.MustCompile(`(.+)(?:\()(.+)(?:\))`)
	if rep.MatchString(name) {
		mainN := rep.ReplaceAllString(name, "$1")
		otherN := rep.ReplaceAllString(name, "$2")
		e.NormalizedStationName = &NormalizedStationName{e.RemoveChar(mainN), e.RemoveChar(otherN)}

		return
	}

	rep = regexp.MustCompile(`(.+)・(.+)`)
	if rep.MatchString(name) {
		mainN := rep.ReplaceAllString(name, "$1")
		otherN := rep.ReplaceAllString(name, "$2")
		e.NormalizedStationName = &NormalizedStationName{e.RemoveChar(mainN), e.RemoveChar(otherN)}

		return
	}

	e.NormalizedStationName = &NormalizedStationName{e.RemoveChar(name), ""}
}

// RemoveChar は指定の文字を削除する
func (e *AbleStation) RemoveChar(rawName string) string {
	return Trim(rawName)
	// name := RemoveWord(rawName)
	// return RemoveSufix(name)
}

// DoorStation はDOORの駅のデータ構造
type DoorStation struct {
	NormalizedStationName *NormalizedStationName

	ID     int64
	Name   string
	LineID int64

	// ID                int64
	// CityID            int64
	// PrefID            int64
	// Name              string
	// NameRuby          string
	// Latitude          string
	// Longitude         string
	// NameDuplicateFlag bool
	// TmpID             int64
	// EkitanEkiCode     int64
	// Status            int64
}

// NormalizeDoorStation はDOORの駅の正規化
func (d *DoorStation) NormalizeDoorStation() {
	name := SymbolToHan(d.Name)

	rep := regexp.MustCompile(`(.+)(?:\()(.+)(?:\))`)
	if rep.MatchString(name) {
		mainN := rep.ReplaceAllString(name, "$1")
		// otherN := rep.ReplaceAllString(name, "$2")

		// rep = regexp.MustCompile(`(.+)・(.+)`)
		// if rep.MatchString(otherN) {
		// 	otherN1 := rep.ReplaceAllString(otherN, "$1")
		// 	// otherN2 := rep.ReplaceAllString(otherN, "$2")
		// 	d.NormalizedStationName = &NormalizedStationName{d.RemoveChar(mainN), "", d.RemoveChar(otherN1)}

		// 	return
		// }

		d.NormalizedStationName = &NormalizedStationName{d.RemoveChar(mainN), ""}

		return
	}

	rep = regexp.MustCompile(`(.+)・(.+)`)
	if rep.MatchString(name) {
		mainN := rep.ReplaceAllString(name, "$1")
		// otherN := rep.ReplaceAllString(name, "$2")
		d.NormalizedStationName = &NormalizedStationName{d.RemoveChar(mainN), ""}

		return
	}

	d.NormalizedStationName = &NormalizedStationName{d.RemoveChar(name), ""}
}

// RemoveChar は指定の文字を削除する
func (d *DoorStation) RemoveChar(rawName string) string {
	return Trim(rawName)
	// name := RemoveWord(rawName)
	// name = RemovePrefix(name)
	// return RemoveSufix(name)
}

// DoorStations はDOORの駅の構造配列
type DoorStations []DoorStation

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
	rel = strings.Replace(str, "）", ")", -1)
	rel = strings.Replace(str, "＜", "<", -1)
	rel = strings.Replace(str, "＞", ">", -1)
	return rel
}

// RemoveWord は特定のキーワードを削除
func RemoveWord(name string) string {
	rep := regexp.MustCompile(`(都市)`)
	return rep.ReplaceAllString(name, "")
}

// RemovePrefix は接頭語を削除
func RemovePrefix(name string) string {
	rep := regexp.MustCompile(`^(ＪＲ|JR)`)
	return rep.ReplaceAllString(name, "")
}

// RemoveSufix は接尾語を削除
func RemoveSufix(name string) string {
	rep := regexp.MustCompile(`(鉄道線|鉄道|鐵道|線)$`)
	return rep.ReplaceAllString(name, "")
}
