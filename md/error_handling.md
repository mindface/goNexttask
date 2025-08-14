# エラーハンドリング方針

## API層
- ユーザー未発見時: 404 Not Found
- 認証失敗時: 401 Unauthorized
- バリデーションエラー時: 400 Bad Request
- サーバー内部エラー: 500 Internal Server Error

## サービス層
- ドメインルール違反は DomainException としてスロー
- 予期せぬ例外は例外ラップして再スロー
