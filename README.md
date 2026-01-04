# AquesTalk-module-Golang

AquesTalk（Linux版の共有ライブラリ）を Go から呼び出すための薄いラッパーです。

- ライブラリ本体: `package aquestalk`
- デモCLI: `cmd/aquestalk-demo`

## 動作環境

- Linux
- Go（cgo有効）
- AquesTalk の `libAquesTalk.so`

> Windows/macOS でも `go get` やドキュメント参照ができるように、非Linux向けにはスタブ実装を含めています。

## インストール

公開先（GitHub等）に置いたあと、次のように利用します。

```bash
go get <your-module-path>
```

## 使い方（コード例）

```go
package main

import (
	"fmt"
	"log"

	"<your-module-path>"
)

func main() {
	voices, err := aquestalk.ListVoices("/path/to/AquesTalk/lib64")
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

## 公開のメモ（GitHub想定）

1. GitHub にリポジトリ作成（例: `github.com/<owner>/<repo>`）
2. `go.mod` の `module` を `github.com/<owner>/<repo>` に変更
3. タグを付ける（例: `v0.1.0`）

```bash
git tag v0.1.0
git push --tags
```

これで利用側は `go get github.com/<owner>/<repo>@v0.1.0` のように取得できます。
