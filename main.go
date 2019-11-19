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
  "io"
  "io/ioutil"
  "strings"
  "errors"
)

const tmpFileName string = "yokohama-weather.html"
const port string = "8080"
const url string = "https://www.jma.go.jp/jp/amedas_h/today-46106.html"
const adminMessage string = "予期しないエラーが発生しました。管理者にご連絡下さい。"

type ErrorItem struct {
  Resource string
  Field string
  Message string
}
type ErrorItems []ErrorItem

type ErrorMessage struct {
  // FIXME: 項目名が型名になる
  ErrorItems `json:errors`

  Overview string
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

  // HACK: もっといい初期化の方法がありそう。
  specificTime := -1
  if rawSpecificTime == "" {
    fmt.Println(time.Now())
    specificTime = time.Now().Hour()
  } else {
    specificTime, _ = strconv.Atoi(rawSpecificTime)
  }

  if !(specificTime > 0 && specificTime < 25) {
    // TODO: どの項目が？をレスポンスに設定する
    item := ErrorItem{Resource: "", Field: "time", Message: "時間は、1~24のみ選択が可能です。"}
    var items ErrorItems
    items = append(items, item)
    // item
    message := ErrorMessage{items, "バリデーションエラー"}

    ErrorResponseJson(message, w)

    return
  }

  weatherItems, err := FetchYokohamaWeather()

  if err != nil {
    var items ErrorItems
    message := ErrorMessage{items, err.Error()}

    ErrorResponseJson(message, w)
    return
  }

  weatherItem := weatherItems[specificTime - 1]
  // 値が設定されていない場合は、1時間前に戻る
  if rawSpecificTime == "" && specificTime - 2 > 1 && weatherItem.Temperature == "" {
    weatherItem = weatherItems[specificTime - 2]
  }

  // HACK: 関数化したいけど引数の型を指定しない方法がわからない
  // 配列をjsonに変換する
  res, err := json.Marshal(weatherItem)

  if err != nil {
    log.Panic(err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(res)
}

func ErrorResponseJson(message ErrorMessage, w http.ResponseWriter) error {
  res, err := json.Marshal(message)

  if err != nil {
    return err
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(res)

  return nil
}

func FetchYokohamaWeather() (WeatherItems, error) {
  // キャッシュがなければ取りに行く
  if !Exists(CacheFilePath(tmpFileName)) {
    result := CreateCache()

    if !result {
      return nil, errors.New(adminMessage)
    }
  }

  fileInfos, _ := ioutil.ReadFile(CacheFilePath(tmpFileName))
  stringReader := strings.NewReader(string(fileInfos))

  return AnalyzeHtml(stringReader)
}

func CreateCache() bool {
  res, _ := http.Get(url)
  defer res.Body.Close()

  if res.StatusCode != http.StatusOK {
    log.Fatal("Status code isn't OK. It was " + string(res.StatusCode))
    return false
  }

  file, err := os.Create(CacheFilePath(tmpFileName))
  if err != nil {
    log.Fatal(err)
    return false
  }
  defer file.Close()

  bodyBytes, err := ioutil.ReadAll(res.Body)
  if err != nil {
    log.Fatal(err)
    return false
  }

  file.Write(bodyBytes)

  return true
}

func AnalyzeHtml(r io.Reader) (weatherItems WeatherItems, err error) {
  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(r)
  if err != nil {
    log.Fatal("HTMLの解析が正しく行えませんでした。")
    return nil, err
  }

  // HACKME: 配列の初期化の方法がわからない
  // weatherItems := make([]WeatherItems, 0)
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
  // https://qiita.com/unbabel/items/c8782420391c108e3cac
  const tmpFileLayout = "2006010203"

  return "./.cache" + "/" + name + t.Format(tmpFileLayout)
}