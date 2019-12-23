# ABLEのマッピング表を半自動化で更新するスクリプト

## なんのための？
ABLEから送信される物件が、最寄り駅が紐付けられずエラーとなり取り込めない事を防ぐ。

## 実行タイミング
不定期。
取り込み時に発生する駅に関するエラー件数を減らす際に実行する。

## 注意
このスクリプトは基本的に名前一致により、ABLEとDOORの駅、沿線を紐づけている。
故に、紐付け対象は*必ずしも一致しない*事を念頭に置き、最終的には目視で担保する事を前提とする。

## 更新テーブル
* import_able_door_station
* import_able_line
* import_able_line_station

## 実行手順
### FTPサーバにアクセスする
下記情報を元に、FTPクライアント等を使用してサーバにアクセスする。

chintai.door.ac/apps/backend/config/app.yml
```
able:
server:
host: '***.***.***.***'
user: '***'
password: '***'
```

### FTPよりファイルをダウンロードする
* mst_railway_company_XXXX.tsv
* ct_mst_pref_ensen_XXXX.tsv
* ct_mst_pref_eki_XXXX.tsv

※XXXXには日付が入る。
※エンコードはUTF-8で保存する。(utility/convert_to_utf8を使用)

### DBよりデータを取得する
* s_company(テーブル丸ごと)
* line(テーブル丸ごと)
* 下記SQLの結果(station_and_line.tsv)
```
select
s.id,s.name,sr.line_id
from station as s
inner join station_relation as sr on sr.station_id = s.id
```

※データはタブ区切り(tsv)、エンコードはUTF-8で保存する。

### 鉄道会社の紐付け
#### 対象ディレクトリ
diff_company

#### 設置ファイル
* mst_railway_company_XXXX.tsv
* s_company.tsv

#### 実行
```
$ ./diff_company_exe mst_railway_company.tsv s_company.tsv > compare_company.tsv
```

##### 出力結果例
```
46	"道南いさりび鉄道"	174	"道南いさりび鉄道"
0	"ＪＲ"	1	"JR"
11	"東武鉄道"	2	"東武鉄道"
```

※出力結果で不整合な箇所を修正する

### 沿線の紐付け
diff_line

#### 設置ファイル
* ct_mst_pref_ensen_XXXX.tsv
* line.tsv
* compare_company.tsv(diff_companyの出力結果)

#### 実行
```
$ ./diff_line_exe compare_company.tsv ct_mst_pref_ensen.tsv line.tsv > compare_line.tsv
```

##### 出力結果例
```
501001	"函館本線"	1	3	"ＪＲ函館本線"
501002	"札沼線<学園都市線>"	1	12	"ＪＲ札沼線"
501003	"千歳線"	1	9	"ＪＲ千歳線"
```

※出力結果で不整合な箇所を修正する

### 駅の紐付け
diff_station

#### 設置ファイル
* ct_mst_pref_eki_XXXX.tsv
* station_and_line.tsv
* compare_line.tsv(diff_lineの出力結果)

#### 実行
```
$ ./diff_station_exe compare_line.tsv ct_mst_pref_eki.tsv station_and_line.tsv > compare_station.tsv
```

##### 出力結果例
```
able_station_id	able_station_name	able_line_id	able_line_name	able_ensentype	door_station_id	door_station_name
"004"	"高崎"	"102151"	"上信電鉄上信線"	"2"	10813	"高崎(上信)"
"003"	"浅草"	"102160"	"つくばエクスプレス"	"2"	10866	"浅草(つくばＥＸＰ)"
"042"	"弘明寺"	"103106"	"横浜市営地下鉄ブルーライン"	"3"	10748	"弘明寺(横浜市営)"
```

※出力結果で不整合な箇所を修正する

### SQLの作成
make_sql

#### 設定ファイル
* compare_station.tsv(diff_stationの出力結果)

#### 実行
```
$ ./make_sql_exe compare_station.tsv > insert_to_station_and_line_for_able.sql
```

##### 出力結果例
```
INSERT IGNORE INTO `import_able_line`( `able_line_id`, `name` ) VALUES ( '102151', '高崎(上信)' );
INSERT INTO `import_able_door_station`( `able_line_id`, `able_station_id`, `door_station_id`, `station_name` ) SELECT * FROM ( SELECT '102151', '004', 10813, '高崎(上信)' ) AS tmp WHERE NOT EXISTS( SELECT * FROM import_able_door_station WHERE able_line_id = '102151' and able_station_id = '004' and door_station_id = 10813 ) LIMIT 1;
INSERT INTO `import_able_line_station`( `able_line_id`, `able_station_id`, `able_station_name`, `able_line_name`, `able_ensentype` ) SELECT * FROM ( SELECT '102151', '004', '高崎', '上信電鉄上信線', '2' ) AS tmp WHERE NOT EXISTS( SELECT * FROM import_able_line_station WHERE able_line_id = '102151' and able_station_id = '004' and able_ensentype = '2' ) LIMIT 1;
```

### DBへの反映
insert_to_station_and_line_for_able.sql(make_sqlの出力結果)を実行する。

#### 整合性の確認
```
select
    ials.id,
    ials.able_station_name,
    ials.able_line_name,
    s.id,
    s.name,
    e.id,
    e.rosen_name
from
    import_able_line_station as ials
    left join
        import_able_line as ial
    on  ial.able_line_id = ials.able_line_id
    left join
        import_able_door_station as iads
    on  iads.able_station_id = ials.able_station_id
    and iads.able_line_id = ials.able_line_id
    left join
        station as s
    on  s.id = iads.door_station_id
    left join
        ekitan as e
    on  e.eki_code = s.ekitan_eki_code
where
    ials.id >= XXXX # 今回追加したIDを設定
having ials.able_station_name != s.name
or  s.name is null
```

#### 出力結果例
```
142571	郡元	鹿児島市電１系統	NULL	NULL	NULL	NULL
142571	郡元	鹿児島市電１系統	10956	郡元（南側）	10086	鹿児島市電１系統
142572	すすきの	札幌市電<札幌市交通局>	7187	すすきの(市営)	9155	札幌市営地下鉄南北線
```

上記では、郡元の紐付きが正しくないと見える。compare_station.tsvを手動で修正し再度SQLを作成する。
