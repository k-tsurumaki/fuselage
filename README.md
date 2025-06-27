# Fuselage

軽量でミニマルなGo REST APIフレームワーク

## 特徴

- **標準ライブラリベース**: http.ServeMuxを使用した高速ルーティング
- **YAML設定**: Kubernetes風の設定ファイル対応
- **シンプル**: 最小限の機能でREST APIを構築
- **拡張可能**: ミドルウェアによる機能拡張
- **テンプレート**: Service/Domain/Adapterレイヤーのひな形提供

## クイックスタート

### 1. 基本的な使用方法

```go
package main

import (
    "log"
    "github.com/k-tsurumaki/fuselage"
)

func main() {
    router := fuselage.New()
    
    router.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    server := fuselage.NewServer(":8080", router)
    log.Fatal(server.ListenAndServe())
}
```

### 2. YAML設定ファイルを使用

`config.yaml`:
```yaml
server:
  host: "localhost"
  port: 8080
  readTimeout: 15s
  writeTimeout: 15s
  idleTimeout: 60s

middleware:
  - logger
  - recover
  - timeout
```

```go
func main() {
    config, err := fuselage.LoadConfig("config.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    router := fuselage.New()
    router.GET("/users/:id", getUserHandler)
    
    server := fuselage.NewServerFromConfig(config, router)
    log.Fatal(server.ListenAndServe())
}
```

## API

### Router

- `New()` - 新しいRouterインスタンスを作成
- `GET(path, handler)` - GETルートを登録
- `POST(path, handler)` - POSTルートを登録
- `PUT(path, handler)` - PUTルートを登録
- `DELETE(path, handler)` - DELETEルートを登録
- `Use(middleware)` - ミドルウェアを追加

### パラメータ抽出

```go
router.GET("/users/:id", func(w http.ResponseWriter, r *http.Request) {
    id := fuselage.GetParam(r, "id")
    // ...
})
```

### 組み込みミドルウェア

- `Logger` - リクエストログ出力
- `Recover` - パニック回復
- `Timeout(duration)` - リクエストタイムアウト

### 設定

```yaml
server:
  host: "0.0.0.0"      # サーバーホスト
  port: 8080           # ポート番号
  readTimeout: 15s     # 読み取りタイムアウト
  writeTimeout: 15s    # 書き込みタイムアウト
  idleTimeout: 60s     # アイドルタイムアウト

middleware:            # 適用するミドルウェア
  - logger
  - recover
  - timeout
```

## アーキテクチャテンプレート

`templates/`ディレクトリにレイヤードアーキテクチャのひな形を提供:

- `domain.go` - ドメインエンティティとリポジトリインターフェース
- `service.go` - ビジネスロジック層
- `adapter.go` - 外部サービス連携層

## 開発

### テスト実行

```bash
go test -v -cover
```

### CI/CD

GitHub Actionsによる自動化:
- **CI**: テスト、フォーマットチェック、リンター、カバレッジレポート

## サンプル

`example/`ディレクトリにCRUD操作を持つユーザー管理APIのサンプルを提供。

```bash
cd example
go run main.go
```

## ライセンス

MIT License