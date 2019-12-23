package main

import (
	"encoding/csv"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ReadAbleEnsen is read csv and return struct
func ReadEseikatsuLine(fileName string) (EseikatsuLines, error) {
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

	var eseikatsuLines EseikatsuLines
	for i, line := range lines {
		if i == 0 {
			continue
		}
		eseikatsuLine := new(EseikatsuLine)

		var tmp int64
		tmp, _ = strconv.ParseInt(line[0], 10, 64)
		eseikatsuLine.LineID = tmp
		tmp, _ = strconv.ParseInt(line[1], 10, 64)
		eseikatsuLine.StationID = tmp
		lAndS := strings.Split(line[2], "－")
		if len(lAndS) == 2 {
			eseikatsuLine.LineName = lAndS[0]
			eseikatsuLine.StationName = lAndS[1]
		} else {
			rep := regexp.MustCompile(`(.+)［.+－.+］－(.+)$`)
			eseikatsuLine.LineName = rep.ReplaceAllString(line[2], "$1")
			eseikatsuLine.StationName = rep.ReplaceAllString(line[2], "$2")
		}
		eseikatsuLine.PrefName = line[3]

		eseikatsuLine.NormalizeLine()

		eseikatsuLines = append(eseikatsuLines, *eseikatsuLine)
	}

	// for _, e := range eseikatsuLines {
	// 	fmt.Println(e.LineName, "(", e.LineID, ")")
	// }
	// fmt.Println("========================")
	// for _, e := range Uniq(eseikatsuLines) {
	// 	fmt.Println(e.LineName, "(", e.LineID, ")")
	// }
	// fmt.Println("========================")

	return Uniq(eseikatsuLines), nil
}

func Uniq(eseikatsuLines EseikatsuLines) EseikatsuLines {
	m := make(map[int64]bool)
	var uniq EseikatsuLines

	for _, eseikatsuLine := range eseikatsuLines {
		if !m[eseikatsuLine.LineID] {
			m[eseikatsuLine.LineID] = true
			uniq = append(uniq, eseikatsuLine)
		}
	}
	return uniq
}

// ReadDoorEki is read csv and return struct
func ReadDoorLine(fileName string) (DoorLines, error) {
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

	var lines DoorLines
	for i, lin := range lins {
		if i == 0 {
			continue
		}

		line := new(DoorLine)

		line.Name = lin[2]

		var tmp int64
		tmp, _ = strconv.ParseInt(lin[0], 10, 64)
		line.ID = tmp
		tmp, _ = strconv.ParseInt(lin[1], 10, 64)
		line.PrefID = tmp
		tmp, _ = strconv.ParseInt(lin[3], 10, 64)
		line.CompanyID = tmp
		tmp, _ = strconv.ParseInt(lin[4], 10, 64)
		line.TmpID = tmp
		tmp, _ = strconv.ParseInt(lin[5], 10, 64)
		line.EkitanRosenCode = tmp
		tmp, _ = strconv.ParseInt(lin[6], 10, 64)
		line.Status = tmp

		line.NormalizeLine()

		lines = append(lines, *line)
	}

	return lines, nil
}
