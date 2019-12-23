package main

import (
	"regexp"
	"strings"

	"golang.org/x/text/width"
)

// AbleCompany is エイブル鉄道会社
type AbleCompany struct {
	DoorCompany *DoorCompany

	RailwaycoNo   int64
	RailwaycoName string
	RailwaycoKana string
	RailwaycoSeq  int64
	DelKbn        string
	// InsDt         time.Time
	// InsPg         string
	// InsID         string
	// UpdDt         time.Time
	// UpdPg         string
	// UpdID         string
	NormalizedRailwaycoName string
}

// NormalizeAbleCompany はエイブルの会社名の正規化
func NormalizeAbleCompany(rawName string) string {
	rep := regexp.MustCompile(`(.+)(?:＜)(.+)(?:＞)`)
	name := rep.ReplaceAllString(rawName, "$1")
	name = width.Fold.String(name)
	name = strings.ToLower(name)

	return name
}

// DoorCompany is DOOR鉄道会社
type DoorCompany struct {
	ID                   int64
	Name                 string
	EkitanName           string
	Status               int64
	NormalizedName       string
	NormalizedEkitanName string
}

// NormalizeDoorCompany はDOORの会社名の正規化
func NormalizeDoorCompany(rawName string) string {
	name := strings.ReplaceAll(rawName, "鐵道", "鉄道")
	name = width.Fold.String(name)
	name = strings.ToLower(name)

	return name
}

// DoorCompanies is DOOR鉄道会社配列
type DoorCompanies []DoorCompany

// Remove is remove from array specifing position
func (l *DoorCompanies) Remove(index int) {
	res := DoorCompanies{}
	for i, v := range *l {
		if i == index {
			continue
		}
		res = append(res, v)
	}
	*l = res
}
