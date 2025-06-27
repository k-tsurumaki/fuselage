# Fuselage Example

Fuselageフレームワークを使用したサンプルREST APIアプリケーション

## 実行方法

```bash
go run main.go
```

サーバーが起動したら、以下のエンドポイントにアクセス可能:

## API エンドポイント

### ユーザー一覧取得
```bash
curl http://localhost:8080/users
```

### 特定ユーザー取得
```bash
curl http://localhost:8080/users/1
```

### ユーザー作成
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Charlie"}'
```

### ユーザー更新
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated"}'
```

### ユーザー削除
```bash
curl -X DELETE http://localhost:8080/users/1
```

## 機能

- RESTful API設計
- URLパラメータ抽出
- JSON レスポンス
- エラーハンドリング
- ミドルウェア（ログ、パニック回復、タイムアウト）