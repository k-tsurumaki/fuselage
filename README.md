# Fuselage

軽量でミニマルなGo REST APIフレームワーク

## 特徴

- **標準ライブラリのみ**: 外部依存ゼロ
- **高速**: net/httpベースで高いパフォーマンス
- **シンプル**: 最小限の機能でREST APIを構築
- **拡張可能**: ミドルウェアによる機能拡張

## クイックスタート

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
)

func main() {
    router := New()
    
    // ミドルウェアの適用
    router.Use(Logger)
    router.Use(Recover)
    router.Use(Timeout(30 * time.Second))
    
    // ルートの定義
    router.GET("/users/:id", func(w http.ResponseWriter, r *http.Request) {
        id := GetParam(r, "id")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"id": id})
    })
    
    server := NewServer(":8080", router)
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
    id := GetParam(r, "id")
    // ...
})
```

### 組み込みミドルウェア

- `Logger` - リクエストログ出力
- `Recover` - パニック回復
- `Timeout(duration)` - リクエストタイムアウト

## テスト実行

```bash
go test
```

## サンプル実行

```bash
go run main.go
```

サーバーが起動したら以下のエンドポイントにアクセス可能:

- `GET /users` - 全ユーザー取得
- `GET /users/:id` - 特定ユーザー取得
- `POST /users` - ユーザー作成
- `PUT /users/:id` - ユーザー更新
- `DELETE /users/:id` - ユーザー削除