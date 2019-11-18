package main

import (
  "fmt"
  "log"
  "net/http"
  "github.com/PuerkitoBio/goquery"
  "encoding/json"
  "strconv"
  "time"
  "os"
  "io/ioutil"
  "strings"
)

const tmpFileName string = "yokohama-weather.html"
const port string = "8080"
const url string = "https://www.jma.go.jp/jp/amedas_h/today-46106.html"

type ErrorMessage struct {
  Message string
}

// [時刻 気温 降水量 風向 風速 日照時間 湿度 気圧]
type WeatherUnit struct {
  // 時刻
  Time string
  // 気温
  Temperature string
  // 降水量
  Precipitation string
  // 風向
  WindowDirection string
  // 風速
  WindowSpeed string
  // 日照時間
  SunshineHours string
  // 湿度
  Humidity string
  // 気圧
  BarometricPressure string
}

type WeatherItem struct {
  WeatherUnit `json:"unit"`

  // 時刻
  Time string
  // 気温
  Temperature string
  // 降水量
  Precipitation string
  // 風向
  WindowDirection string
  // 風速
  WindowSpeed string
  // 日照時間
  SunshineHours string
  // 湿度
  Humidity string
  // 気圧
  BarometricPressure string
}
type WeatherItems []WeatherItem

func main() {
  http.HandleFunc("/yokohama", YokohamaHandler) // ハンドラを登録してウェブページを表示させる

  fmt.Println("application starting on localhost:" + port)
  http.ListenAndServe(":" + port, nil)
}

func YokohamaHandler(w http.ResponseWriter, r *http.Request) {
  rawSpecificTime := r.URL.Query().Get("time")

  specificTime := -1
  if rawSpecificTime == "" {
    fmt.Println(time.Now())
    specificTime = time.Now().Hour()
  } else {
    specificTime, _ = strconv.Atoi(rawSpecificTime)
  }

  if !(specificTime > 0 && specificTime < 25) {
    message := ErrorMessage{Message: "時間は、1~24のみ選択が可能です。"}
    res, _ := json.Marshal(message)

    w.Header().Set("Content-Type", "application/json")
    w.Write(res)

    return
  }

  weatherItems := FetchYokohamaWeather()
  weatherItem := weatherItems[specificTime - 1]

  // 配列をjsonに変換する
  res, err := json.Marshal(weatherItem)

  if err != nil {
    fmt.Println(err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(res)
}

func FetchYokohamaWeather() (weatherItems WeatherItems) {
  // キャッシュがなければ取りに行く
  if !Exists(CacheFilePath(tmpFileName)) {
    res, _ := http.Get(url)
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
      log.Fatal("Status code isn't OK. It was " + string(res.StatusCode))
      // FIXME: errorで止めた方がよさげ
    }

    file, err := os.Create(CacheFilePath(tmpFileName))
    if err != nil {
      log.Panic(err)
    }
    defer file.Close()

    bodyBytes, err := ioutil.ReadAll(res.Body)
    if err != nil {
      log.Fatal(err)
    }

    file.Write(bodyBytes)
  }

  fileInfos, _ := ioutil.ReadFile(CacheFilePath(tmpFileName))
  stringReader := strings.NewReader(string(fileInfos))

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(stringReader)
  if err != nil {
    log.Fatal(err)
  }

  units := []string{}
  doc.Find("#tbl_list tr:nth-child(2) td").Each(func(i int, td *goquery.Selection) {
    units = append(units, td.Text())
  })
  weatherUnit := WeatherUnit{units[0], units[1], units[2], units[3], units[4], units[5], units[6], units[7]}

  // Find the review items
  doc.Find("#tbl_list tr").Each(func(i int, tr *goquery.Selection) {
    if i > 1 {
      elements := []string{}
      // For each item found, get the band and title
      tr.Find("td").Each(func(i int, td *goquery.Selection) {
        elements = append(elements, td.Text())
      })

      weatherItem := WeatherItem{weatherUnit, elements[0], elements[1], elements[2], elements[3], elements[4], elements[5], elements[6], elements[7]}
      weatherItems = append(weatherItems, weatherItem)
    }
  })

  return
}

func Exists(filename string) bool {
    _, err := os.Stat(filename)
    return err == nil
}

func CacheFilePath(name string) string {
  t := time.Now()
  const tmpFileLayout = "20000101"

  return "./cache" + "/" + name + t.Format(tmpFileLayout)
}