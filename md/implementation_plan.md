# 実装計画（依存順序付き）

## ログインイベント
**ステップ1: ユーザーログインサービスの作成**
- 編集対象ファイル: src/domain/user/AuthService.java
- 編集対象メソッド: verifyPassword(), generateJwtToken()
- 目的: パスワード検証とJWT生成のビジネスロジックを作成
- 内容: 既存の暗号化ライブラリを用いてverifyPasswordを実装、JWT生成はJWTユーティリティを利用
- 活用クラス: util.JwtUtil, util.PasswordHasher

**ステップ2: ユーザーログインAPIの作成**
- 編集対象ファイル: src/application/auth/AuthController.java
- 編集対象メソッド: login()
- 目的: APIリクエストを受け付け、AuthServiceを呼び出す
- 内容: リクエストからemailとpasswordを取得、UserRepository#findByEmail()でユーザーを取得し、AuthServiceで検証・JWT発行
- 活用クラス: domain.user.UserRepository, domain.user.AuthService

**ステップ3: 例外処理の追加**
- 編集対象ファイル: src/application/auth/AuthController.java
- 目的: ユーザー未発見・認証失敗時の例外処理
- 内容: UserNotFoundException, InvalidCredentialsException をスローし、ExceptionHandlerでHTTPステータスを制御
