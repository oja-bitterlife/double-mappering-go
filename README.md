# goのmapをダブルバッファを使ってトランザクション的に更新できるようにするライブラリ

## 特徴

- ネイティブGo
- ダブルバッファの作成にユーザー定義のSerialize/Deserializeを使用する
    - リフレクションを使わないSerializerを使うことができる
- 中身はただのmapや構造体なのでon memoryで高速
- ファイル保存のためにbyte[]出力できる
    - もちろんそれをrestoreできる
- スレッドセーフ
- One Fileで完結。ファイルコピーだけですぐに使える
    - もちろんgo getでも使える
- tinygo対応
    - まだ未確認_(:3」∠)_ﾀﾌﾞﾝﾀﾞｲｼﾞｮｳﾌﾞ...

## インストール
```bash
go get github.com/oja-bitterlife/double-mapparing-go
```
あるいは`double-mapparing.go`をコピーして自分のプロジェクトにおいてください。

## 使い方
```go
type Config struct {
    AppID string `json:"app_id"`
}

// 初期化
dbm := doublemapparing.New(
    func(c *Config) ([]byte, error) { return json.Marshal(c) },
    func(b []byte) (*Config, error) {
        var c Config
        err := json.Unmarshal(b, &c)
        return &c, err
    },
)

// トランザクション的な更新
dbm.Update(func(data *Config) error {
    data.AppID = "New System"
    return nil // nilを返すと更新が確定
})

// 安全な読み取り（クローンを取得）
view, _ := dbm.View()
fmt.Println(view.AppID)
```
exampleも多少置いてますが、NewとUpdateとView以外を使うことってそうそうないかなって。

一応ファイル入出力用のメソッドも用意してあります
```go
// バイト列として書き出し
b, _ := dbm.Bytes()
os.WriteFile("config.json", b, 0644)

// ファイルから復元
data, _ := os.ReadFile("config.json")
dbm.Restore(data)
```
