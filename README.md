# AquesTalk-module-Golang

AquesTalk（Linux版の共有ライブラリ）を Go から呼び出すための薄いラッパーです。

- ライブラリ本体: `package aquestalk`
- デモCLI: `cmd/aquestalk-demo`

## 機能

このモジュールが提供する機能は次のとおりです。

- AquesTalk共有ライブラリのロード
- 音声合成（WAVバイト列の取得）
- 声種（ボイス）一覧の列挙（`xxx/<voice>/libAquesTalk.so` を探索）
- 入力テキストのひらがな寄せ（kagome IPA辞書ベース）
- WAVヘッダのサンプルレート/バイトレートを書き換える再生速度変更

公開API（`package aquestalk`）

メイン（推奨）

- `func Open(libPath, devKey string) (*AquesTalk, error)`
- `func (a *AquesTalk) Speak(text string, speechSpeed int) ([]byte, error)`
- `func (a *AquesTalk) SpeakWithPlayback(text string, speechSpeed int, playbackSpeed int) ([]byte, error)`
- `func Voices(baseDir string) ([]Voice, error)`
- `func ToHiragana(text string) string`
- `func ChangePlaybackSpeed(wav []byte, playbackSpeed int) ([]byte, error)`

互換（旧名のラッパー）

- `func LoadAquesTalk(libPath, devKey string) (*AquesTalk, error)`
- `func (a *AquesTalk) Synthe(text string, speechSpeed int) ([]byte, error)`
- `func (a *AquesTalk) SyntheWithPlaybackSpeed(text string, speechSpeed int, playbackSpeed int) ([]byte, error)`
- `func ListVoices(baseDir string) ([]Voice, error)`
- `func Normalize(text string) string`
- `func ApplyPlaybackSpeed(wav []byte, playbackSpeed int) ([]byte, error)`

## 動作環境

- Linux
- Go（cgo有効）
- AquesTalk の `libAquesTalk.so`

補足

- このモジュールは Linux + cgo のときのみ AquesTalk を実際に呼び出します。
- Windows/macOS（または cgo 無効）でも `go get` や `go doc` が破綻しないよう、スタブ実装を含めています（その環境で `LoadAquesTalk` / `Synthe` を呼ぶとエラーになります）。

## インストール

次のように利用します。

```bash
go get github.com/retdayo/AquesTalk-module-Golang
```

## 使い方（コード例）

```go
package main

import (
	"fmt"
	"log"

	"github.com/retdayo/AquesTalk-module-Golang"
)

func main() {
	voices, err := aquestalk.Voices("/path/to/AquesTalk/lib64")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(voices)
}
```

## デモCLI

```bash
go run ./cmd/aquestalk-demo
```

`cmd/aquestalk-demo/config.go` の `AQUESTALK_BASE_DIR` と `DEV_KEY` を環境に合わせて設定してください。

デモの動作

- `Voices` で声種一覧を取得
- 対話で声種を選択
- `ToHiragana` で入力テキストを正規化
- `SpeakWithPlayback` でWAV生成＋再生速度変更
- `out.wav` に保存

## 各APIの説明

### Open

`Open(libPath, devKey) (*AquesTalk, error)`

AquesTalk の共有ライブラリ（`libAquesTalk.so`）をロードして、呼び出し用のハンドルを作ります。

- `libPath`
	- 例: `/path/to/AquesTalk/lib64/f1/libAquesTalk.so`
	- 実体ファイルが無い/読めない場合は失敗します
- `devKey`
	- AquesTalk の開発キー文字列
	- 無効な場合の挙動は AquesTalk 側仕様に従います（このラッパーでは検証しません）

失敗しやすいポイント

- Linux + cgo 以外では常にエラー（スタブ実装）
- `dlopen` 失敗（パス違い、依存ライブラリ不足など）

### Speak

`(*AquesTalk) Speak(text, speechSpeed) ([]byte, error)`

テキストを音声合成して WAV バイト列を返します。

- `text`
	- AquesTalk に渡す文字列
	- このリポジトリの想定は「`ToHiragana` 済み」を推奨（もちろん任意）
- `speechSpeed`
	- AquesTalk 側の速度パラメータ
	- 有効範囲は AquesTalk 側仕様に従います（このラッパーは値域チェックをしません）

戻り値

- 成功時: WAV（ヘッダ含む）
- 失敗時: error

### SpeakWithPlayback

`(*AquesTalk) SpeakWithPlayback(text, speechSpeed, playbackSpeed) ([]byte, error)`

`Speak` で生成した WAV に対して、`ChangePlaybackSpeed` を適用して返します。

- `playbackSpeed`
	- 100 で等速、200 で2倍など
	- 0 以下はエラー

注意

- この処理は WAV データをリサンプルしません
- WAV ヘッダの `sampleRate` と `byteRate` を倍率で書き換えるだけなので、プレイヤー側では「再生が速く（短く）なり、音程も上がる」方向になります

### Voices

`Voices(baseDir) ([]Voice, error)`

声種（ボイス）をディレクトリ探索で列挙します。

- `baseDir`
	- `lib64` ディレクトリを指す想定
	- `baseDir/<voice>/libAquesTalk.so` が存在するものだけを返します

戻り値

- `[]Voice`
	- `Voice.Name` はフォルダ名
	- `Voice.Path` は `libAquesTalk.so` のフルパス

### ToHiragana

`ToHiragana(text) string`

日本語テキストを解析して、読み（カタカナ）をひらがなへ寄せます。

挙動

- kagome（IPA辞書）でトークナイズ
- トークンの読み（カタカナ）をひらがなに変換
- 読みが取れないトークンは表層（原文）をそのまま使います
- 記号は表層をそのまま使います
- AquesTalk 向けに簡易置換をします（例: `!`/`?` を `。` に置換、長音を `う` に置換）

### ChangePlaybackSpeed

`ChangePlaybackSpeed(wav, playbackSpeed) ([]byte, error)`

WAV ヘッダの `sampleRate` と `byteRate` を倍率で書き換えて、再生速度を変えます。

入力

- `wav`: WAV バイト列（44バイト以上のヘッダを前提）
- `playbackSpeed`: 100=等速、200=2倍

戻り値と注意

- `wav` の同じスライスを破壊的に書き換えます（コピーは作りません）
- `len(wav) < 44` はエラー
- `playbackSpeed <= 0` はエラー

## 注意

- AquesTalk本体（共有ライブラリ）や辞書・鍵の配布条件は、このリポジトリとは別に確認してください。