package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadAbleCompany(fileName string) (ableCompanies AbleCompanies, err error) {
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

	for _, line := range lines {
		ableCompany := new(AbleCompany)

		ableCompany.AbleCompanyName = line[1]
		ableCompany.DoorCompanyName = line[3]

		tmp, _ := strconv.ParseInt(line[0], 10, 64)
		ableCompany.AbleCompanyID = tmp
		tmp, _ = strconv.ParseInt(line[2], 10, 64)
		ableCompany.DoorCompanyID = tmp

		ableCompanies = append(ableCompanies, *ableCompany)
	}

	return ableCompanies, nil
}

// ReadAbleEnsen is read csv and return struct
func ReadAbleEnsen(fileName string, ableCompanies AbleCompanies) ([]CtMstPrefEnsen, error) {
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

	var ctMstPrefEnsens []CtMstPrefEnsen
	for _, line := range lines {
		ctMstPrefEnsen := new(CtMstPrefEnsen)

		ctMstPrefEnsen.Prefcd = line[0]
		ctMstPrefEnsen.Ensencd = line[1]
		ctMstPrefEnsen.Ensenname = line[2]
		ctMstPrefEnsen.Ensentype = line[3]
		ctMstPrefEnsen.Areacd = line[4]
		ctMstPrefEnsen.Ensenkana = line[8]
		ctMstPrefEnsen.NormalizeAbleLine()

		var tmps uint64
		tmps, _ = strconv.ParseUint(line[5], 10, 32)
		ctMstPrefEnsen.SortKey = uint32(tmps)
		tmps, _ = strconv.ParseUint(line[6], 10, 32)
		ctMstPrefEnsen.BukkenCnt = uint32(tmps)
		tmps, _ = strconv.ParseUint(line[7], 10, 32)
		ctMstPrefEnsen.ShopCnt = uint32(tmps)
		tmps, _ = strconv.ParseUint(line[9], 10, 32)
		ctMstPrefEnsen.RailwaycoNo = uint32(tmps)
		tmps, _ = strconv.ParseUint(line[10], 10, 32)
		ctMstPrefEnsen.Ensenseq = uint32(tmps)
		tmps, _ = strconv.ParseUint(line[11], 10, 32)
		ctMstPrefEnsen.NyBukkenCnt = uint32(tmps)

		matchFlag := false
		compareRailNo := int64(ctMstPrefEnsen.RailwaycoNo)
		for _, ableCompany := range ableCompanies {
			if ableCompany.AbleCompanyID == compareRailNo {
				ctMstPrefEnsen.AbleCompany = &ableCompany
				// ableCompanies.Remove(lineI)
				matchFlag = true
				break
			}
		}
		if matchFlag {
			ctMstPrefEnsens = append(ctMstPrefEnsens, *ctMstPrefEnsen)
		} else {
			arr := []string{ctMstPrefEnsen.Ensenname, strconv.Itoa(int(ctMstPrefEnsen.RailwaycoNo))}
			fmt.Println("除外沿線:" + strings.Join(arr, ","))
		}
	}

	return ctMstPrefEnsens, nil
}

// ReadDoorEki is read csv and return struct
func ReadDoorEki(fileName string) (Lines, error) {
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

	var lines Lines
	for _, lin := range lins {
		line := new(Line)

		line.Name = lin[2]
		line.NormalizeDoorLine()

		var tmps uint64
		tmps, _ = strconv.ParseUint(lin[0], 10, 32)
		line.ID = uint32(tmps)
		tmps, _ = strconv.ParseUint(lin[1], 10, 32)
		line.PrefID = uint32(tmps)
		tmps, _ = strconv.ParseUint(lin[3], 10, 32)
		line.CompanyID = uint32(tmps)
		tmps, _ = strconv.ParseUint(lin[4], 10, 32)
		line.TmpID = uint32(tmps)
		tmps, _ = strconv.ParseUint(lin[5], 10, 32)
		line.EkitanRosenCode = uint32(tmps)
		tmps, _ = strconv.ParseUint(lin[6], 10, 32)
		line.Status = uint32(tmps)

		lines = append(lines, *line)
	}

	return lines, nil
}
