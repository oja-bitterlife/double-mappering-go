# goのmapをダブルバッファを使ってトランザクション的に更新できるようにするライブラリ

特徴

- ネイティブGo
- ダブルバッファの作成にユーザー定義のSerialize/Deserializeを使用する
    - リフレクションを使わないSerializerを使うことができる
- 中身はただのmapや構造体なのでonmemoryで高速
- ファイル保存のためにbyte[]出力できる
    - もちろんそれをrestoreできる
- goroutineセーフ
- tinygo対応


