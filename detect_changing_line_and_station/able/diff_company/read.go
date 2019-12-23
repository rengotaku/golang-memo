package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

// ReadAbleCompany はCSVを読み込み構造体配列にする
func ReadAbleCompany(fileName string) ([]AbleCompany, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file) //utf8
	reader.Comma = '\t'

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var companies []AbleCompany
	for _, line := range lines {
		company := new(AbleCompany)

		company.RailwaycoName = line[1] // string
		company.RailwaycoKana = line[2] // string
		company.DelKbn = line[4]        // string

		tmp, _ := strconv.ParseInt(line[0], 10, 64)
		company.RailwaycoNo = tmp //  int
		tmp, _ = strconv.ParseInt(line[3], 10, 64)
		company.RailwaycoSeq = tmp // int

		company.NormalizedRailwaycoName = NormalizeAbleCompany(company.RailwaycoName)

		companies = append(companies, *company)
	}

	return companies, nil
}

// ReadDoorCompany はCSVを読み込み構造体配列にする
func ReadDoorCompany(fileName string) (doorCompanies DoorCompanies, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file) //utf8
	reader.Comma = '\t'

	lins, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, lin := range lins {
		company := new(DoorCompany)

		company.Name = lin[1]
		company.EkitanName = lin[2]

		tmp, _ := strconv.ParseInt(lin[0], 10, 64)
		company.ID = tmp
		tmp, _ = strconv.ParseInt(lin[3], 10, 64)
		company.Status = tmp

		company.NormalizedName = NormalizeDoorCompany(company.Name)
		company.NormalizedEkitanName = NormalizeDoorCompany(company.EkitanName)

		doorCompanies = append(doorCompanies, *company)
	}

	return doorCompanies, nil
}
