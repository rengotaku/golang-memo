package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// ReadAbleLine は沿線情報を読み取り、構造体として返却
func ReadAbleLine(fileName string) (ableLines AbleLines, err error) {
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
		ableLine := new(AbleLine)

		tmp, _ := strconv.ParseInt(line[0], 10, 64)
		ableLine.AbleLineID = tmp
		ableLine.AbleLineName = Trim(line[1])
		tmp, _ = strconv.ParseInt(line[2], 10, 64)
		ableLine.AbleLineType = tmp
		tmp, _ = strconv.ParseInt(line[3], 10, 64)
		ableLine.DoorLineID = tmp
		ableLine.DoorLineName = line[4]

		ableLines = append(ableLines, *ableLine)
	}

	return ableLines, nil
}

// ReadAbleStation はエイブルの駅情報を構造体として返却
func ReadAbleStation(fileName string, ableLines AbleLines) ([]AbleStation, error) {
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

	var ableStations []AbleStation
	for _, line := range lines {
		ableStation := new(AbleStation)

		var tmp int64
		ableStation.Prefcd = line[0]
		ableStation.Ensencd = line[1]
		ableStation.Ekicd = line[2]
		ableStation.Ekiname = line[3]
		tmp, _ = strconv.ParseInt(line[4], 10, 64)
		ableStation.SortKey = tmp
		ableStation.StopFlg = line[5]
		tmp, _ = strconv.ParseInt(line[6], 10, 64)
		ableStation.BukkenCnt = tmp
		tmp, _ = strconv.ParseInt(line[7], 10, 64)
		ableStation.ShopCnt = tmp
		tmp, _ = strconv.ParseInt(line[8], 10, 64)
		ableStation.Ekino = tmp
		ableStation.StartekiFlg = line[9]
		tmp, _ = strconv.ParseInt(line[10], 10, 64)
		ableStation.Ekiseq = tmp
		tmp, _ = strconv.ParseInt(line[11], 10, 64)
		ableStation.NyBukkenCnt = tmp

		ableStation.NormalizeAbleStation()

		matchFlag := false
		compareEnsenCd, _ := strconv.ParseInt(ableStation.Ensencd, 10, 64)
		for _, ableLine := range ableLines {
			if ableLine.AbleLineID == compareEnsenCd {
				ableStation.AbleLine = &ableLine
				// ableLines.Remove(lineI)
				matchFlag = true
				break
			}
		}
		if matchFlag {
			ableStations = append(ableStations, *ableStation)
		} else {
			fmt.Println("除外:" + ableStation.Ekiname)
		}
	}

	return ableStations, nil
}

// ReadDoorStation はDOORの駅情報を構造体で返却
func ReadDoorStation(fileName string) (doorStations DoorStations, err error) {
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
		station := new(DoorStation)

		var tmp int64
		tmp, _ = strconv.ParseInt(lin[0], 10, 64)
		station.ID = tmp
		station.Name = lin[1]
		tmp, _ = strconv.ParseInt(lin[2], 10, 64)
		station.LineID = tmp

		// tmp, _ = strconv.ParseInt(lin[0], 10, 64)
		// station.ID = tmp
		// tmp, _ = strconv.ParseInt(lin[1], 10, 64)
		// station.CityID = tmp
		// tmp, _ = strconv.ParseInt(lin[2], 10, 64)
		// station.PrefID = tmp
		// Name = lin[3]
		// station.NameRuby = lin[4]
		// station.Latitude = lin[5]
		// station.Longitude = lin[6]
		// toBool, _ := strconv.ParseBool(line[7])
		// station.NameDuplicateFlag = toBool
		// tmp, _ = strconv.ParseUint(lin[8], 10, 64)
		// station.TmpID = tmp
		// tmp, _ = strconv.ParseUint(lin[9], 10, 64)
		// station.EkitanEkiCode = tmp
		// tmp, _ = strconv.ParseUint(lin[10], 10, 64)
		// station.Status = tmp

		station.NormalizeDoorStation()

		doorStations = append(doorStations, *station)
	}

	return doorStations, nil
}
