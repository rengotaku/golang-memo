package main

import (
	"regexp"
	"strconv"
	"strings"
)

// AbleCompanies はエイブル鉄道会社の構造体
type AbleCompany struct {
	AbleCompanyID   int64
	AbleCompanyName string
	DoorCompanyID   int64
	DoorCompanyName string
}

// AbleCompanies is AbleCompany テーブル 配列データ構造
type AbleCompanies []AbleCompany

// Remove is remove from array specifing position
func (l *AbleCompanies) Remove(index int) {
	res := AbleCompanies{}
	for i, v := range *l {
		if i == index {
			continue
		}
		res = append(res, v)
	}
	*l = res
}

type NormalizedLineName struct {
	Name    string
	SubName string
}

// CtMstPrefEnsen is ct_mst_pref_ensen データ構造
type CtMstPrefEnsen struct {
	AbleCompany        *AbleCompany
	Line               *Line
	NormalizedLineName *NormalizedLineName

	Prefcd      string
	Ensencd     string
	Ensenname   string
	Ensentype   string
	Areacd      string
	SortKey     uint32
	BukkenCnt   uint32
	ShopCnt     uint32
	Ensenkana   string
	RailwaycoNo uint32
	Ensenseq    uint32
	NyBukkenCnt uint32
}

// Equal は沿線名を比較して等しいか比べる
func (e *CtMstPrefEnsen) Equal(l Line) bool {
	prefcd, _ := strconv.Atoi(e.Prefcd)
	if prefcd != int(l.PrefID) {
		return false
	}

	if e.AbleCompany.DoorCompanyID != int64(l.CompanyID) {
		return false
	}

	return e.NormalizedLineName.Name == l.NormalizedLineName.Name ||
		e.NormalizedLineName.SubName == l.NormalizedLineName.Name ||
		e.NormalizedLineName.Name == l.NormalizedLineName.SubName ||
		e.NormalizedLineName.SubName == l.NormalizedLineName.Name
}

// Include は沿線名をどちらかが含んでいないか検証
func (e *CtMstPrefEnsen) Include(l Line) bool {
	prefcd, _ := strconv.Atoi(e.Prefcd)
	if prefcd != int(l.PrefID) {
		return false
	}

	if e.AbleCompany.DoorCompanyID != int64(l.CompanyID) {
		return false
	}

	return (e.Compare(e.NormalizedLineName.Name, l.NormalizedLineName.Name) ||
		e.Compare(e.NormalizedLineName.SubName, l.NormalizedLineName.Name) ||
		e.Compare(e.NormalizedLineName.Name, l.NormalizedLineName.SubName) ||
		e.Compare(e.NormalizedLineName.SubName, l.NormalizedLineName.Name))
}

// Include は沿線名をどちらかが含んでいないか検証
func (e *CtMstPrefEnsen) Compare(n1 string, n2 string) bool {
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

	return strings.Index(x1, x2) > -1
}

// NormalizeAbleLine はエイブルの沿線の正規化
func (e *CtMstPrefEnsen) NormalizeAbleLine() {
	name := e.Ensenname

	rep := regexp.MustCompile(`(.+)(?:<|\(|（)(.+)(?:）|\)|\>)`)
	if rep.MatchString(name) {
		mainN := rep.ReplaceAllString(name, "$1")
		otherN := rep.ReplaceAllString(name, "$2")
		e.NormalizedLineName = &NormalizedLineName{e.RemoveChar(mainN), e.RemoveChar(otherN)}

		return
	}

	rep = regexp.MustCompile(`(.+)・(.+)`)
	if rep.MatchString(name) {
		mainN := rep.ReplaceAllString(name, "$1")
		otherN := rep.ReplaceAllString(name, "$2")
		e.NormalizedLineName = &NormalizedLineName{e.RemoveChar(mainN), e.RemoveChar(otherN)}

		return
	}

	e.NormalizedLineName = &NormalizedLineName{e.RemoveChar(name), ""}
}

func (l *CtMstPrefEnsen) RemoveChar(rawName string) string {
	name := RemoveWord(rawName)
	name = RemovePrefix(name)
	return RemoveSufix(name)
}

// Line is line テーブル データ構造
type Line struct {
	NormalizedLineName *NormalizedLineName

	ID              uint32
	PrefID          uint32
	Name            string
	CompanyID       uint32
	TmpID           uint32
	EkitanRosenCode uint32
	Status          uint32
}

// NormalizeDoorLine はDOORの沿線の正規化
func (l *Line) NormalizeDoorLine() {
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

func (l *Line) RemoveChar(rawName string) string {
	name := RemoveWord(rawName)
	name = RemovePrefix(name)
	return RemoveSufix(name)
}

// Lines is line テーブル 配列データ構造
type Lines []Line

// Remove is remove from array specifing position
func (l *Lines) Remove(index int) {
	res := Lines{}
	for i, v := range *l {
		if i == index {
			continue
		}
		res = append(res, v)
	}
	*l = res
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
