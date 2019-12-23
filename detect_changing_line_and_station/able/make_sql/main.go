package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"text/template"
)

// LineAndStation は駅、沿線のリレーションを表す構造体
type LineAndStation struct {
	AbleStationID   string
	AbleStationName string
	AbleLineID      string
	AbleLineName    string
	AbleEnsentype   string
	DoorStationID   string
	DoorStationName string
}

func main() {
	// -hオプション用文言
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s compare_station.tsv`+"\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		return
	}
	args := flag.Args()

	lineAndStations, err := ReadLineAndStation(args[0])
	if err != nil {
		panic(err)
	}

	for _, lineAndStation := range lineAndStations {
		tmpl, err := template.New("import_able_line").Parse("INSERT IGNORE INTO `import_able_line`( `able_line_id`, `name` ) VALUES ( '{{.AbleLineID}}', '{{.DoorStationName}}' );")
		if err != nil {
			panic(err)
		}
		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, lineAndStation)
		if err != nil {
			panic(err)
		}
		sql1 := tpl.String()

		tpl.Reset()
		tmpl, err = template.New("import_able_door_station").Parse("INSERT INTO `import_able_door_station`( `able_line_id`, `able_station_id`, `door_station_id`, `station_name` ) SELECT * FROM ( SELECT '{{ .AbleLineID }}', '{{ .AbleStationID }}', {{ .DoorStationID }}, '{{ .DoorStationName }}' ) AS tmp WHERE NOT EXISTS( SELECT * FROM import_able_door_station WHERE able_line_id = '{{ .AbleLineID }}' and able_station_id = '{{ .AbleStationID }}' and door_station_id = {{ .DoorStationID }} ) LIMIT 1;")
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(&tpl, lineAndStation)
		if err != nil {
			panic(err)
		}
		sql2 := tpl.String()

		tpl.Reset()
		tmpl, err = template.New("import_able_line_station").Parse("INSERT INTO `import_able_line_station`( `able_line_id`, `able_station_id`, `able_station_name`, `able_line_name`, `able_ensentype` ) SELECT * FROM ( SELECT '{{ .AbleLineID }}', '{{ .AbleStationID }}', '{{ .AbleStationName }}', '{{ .AbleLineName }}', '{{ .AbleEnsentype }}' ) AS tmp WHERE NOT EXISTS( SELECT * FROM import_able_line_station WHERE able_line_id = '{{ .AbleLineID }}' and able_station_id = '{{ .AbleStationID }}' and able_ensentype = '{{ .AbleEnsentype }}' ) LIMIT 1;")
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(&tpl, lineAndStation)
		if err != nil {
			panic(err)
		}
		sql3 := tpl.String()

		fmt.Println(sql1)
		fmt.Println(sql2)
		fmt.Println(sql3)
	}
}

// ReadLineAndStation は駅・沿線を読み取り、構造体として返却
func ReadLineAndStation(filename string) (lineAndStations []LineAndStation, err error) {
	file, err := os.Open(filename)
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
		lineAndStation := new(LineAndStation)

		lineAndStation.AbleStationID = line[0]
		lineAndStation.AbleStationName = line[1]
		lineAndStation.AbleLineID = line[2]
		lineAndStation.AbleLineName = line[3]
		lineAndStation.AbleEnsentype = line[4]
		lineAndStation.DoorStationID = line[5]
		lineAndStation.DoorStationName = line[6]

		lineAndStations = append(lineAndStations, *lineAndStation)
	}

	return lineAndStations, nil
}
