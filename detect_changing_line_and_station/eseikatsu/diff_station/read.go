package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ReadAbleEnsen is read csv and return struct
func ReadEseikatsuLine(fileName string) (eseikatsuLines EseikatsuLines, err error) {
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

	for i, line := range lines {
		if i == 0 {
			continue
		}
		eseikatsuLine := new(EseikatsuLine)

		var tmp int64
		tmp, _ = strconv.ParseInt(line[0], 10, 64)
		eseikatsuLine.EseikatsuLineID = tmp
		eseikatsuLine.EseikatsuLineName = line[1]
		tmp, _ = strconv.ParseInt(line[2], 10, 64)
		eseikatsuLine.DoorLineID = tmp
		eseikatsuLine.DoorLineName = line[3]

		eseikatsuLines = append(eseikatsuLines, *eseikatsuLine)
	}

	return
}

// ReadAbleEnsen is read csv and return struct
func ReadEseikatsuStation(fileName string, eseikatsuLines EseikatsuLines) (eseikatsuStations EseikatsuStations, err error) {
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

	for i, line := range lines {
		if i == 0 {
			continue
		}
		eseikatsuStation := new(EseikatsuStation)

		var tmp int64
		tmp, _ = strconv.ParseInt(line[0], 10, 64)
		eseikatsuStation.LineID = tmp
		tmp, _ = strconv.ParseInt(line[1], 10, 64)
		eseikatsuStation.StationID = tmp
		lAndS := strings.Split(line[2], "－")
		if len(lAndS) == 2 {
			eseikatsuStation.LineName = lAndS[0]
			eseikatsuStation.StationName = lAndS[1]
		} else {
			rep := regexp.MustCompile(`(.+)［.+－.+］－(.+)$`)
			eseikatsuStation.LineName = rep.ReplaceAllString(line[2], "$1")
			eseikatsuStation.StationName = rep.ReplaceAllString(line[2], "$2")
		}
		eseikatsuStation.PrefName = line[3]

		eseikatsuStation.NormalizeStation()

		for _, eseikatsuLine := range eseikatsuLines {
			if eseikatsuStation.LineID == eseikatsuLine.EseikatsuLineID {
				eseikatsuStation.EseikatsuLine = &eseikatsuLine
			}
		}

		// 沿線がマッチしない
		if eseikatsuStation.EseikatsuLine == nil {
			fmt.Println("except:", eseikatsuStation.LineName, "(", eseikatsuStation.LineID, ")")
			continue
		}
		eseikatsuStations = append(eseikatsuStations, *eseikatsuStation)
	}

	return
}

// ReadDoorEki is read csv and return struct
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

	for i, lin := range lins {
		if i == 0 {
			continue
		}

		station := new(DoorStation)

		var tmp int64
		tmp, _ = strconv.ParseInt(lin[0], 10, 64)
		station.ID = tmp
		station.Name = lin[1]
		tmp, _ = strconv.ParseInt(lin[2], 10, 64)
		station.LineID = tmp

		station.NormalizeStation()

		doorStations = append(doorStations, *station)
	}

	return
}
