# Fuselage Examples

## with-config

YAML設定ファイルを使用するサンプル

```bash
cd with-config
go run main.go
```

- ポート: 8081
- 設定ファイル: config.yaml
- ミドルウェア: 設定ファイルで指定

## without-config

設定ファイルを使用しないサンプル

```bash
cd without-config
go run main.go
```

- ポート: 8082
- 設定ファイル: なし
- ミドルウェア: コードで直接指定

## API エンドポイント

両方のサンプルで同じAPIを提供:

- `GET /users` - 全ユーザー取得
- `GET /users/:id` - 特定ユーザー取得
- `POST /users` - ユーザー作成
- `PUT /users/:id` - ユーザー更新
- `DELETE /users/:id` - ユーザー削除