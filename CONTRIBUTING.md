# Contributing to GoNexttask

GoNexttaskプロジェクトへのコントリビュートを歓迎します！このドキュメントでは、プロジェクトへの貢献方法について説明します。

## 開発環境のセットアップ

### 必要な環境

- Go 1.21 以上
- Docker & Docker Compose
- Git
- Make (推奨)

### セットアップ手順

1. **リポジトリをフォーク・クローン**
```bash
git clone https://github.com/YOUR_USERNAME/goNexttask.git
cd goNexttask
```

2. **環境設定**
```bash
cp .env.example .env
# .envファイルを必要に応じて編集
```

3. **開発環境起動**
```bash
make dev-up
# または
docker-compose -f docker-compose.dev.yml up
```

4. **APIテスト**
```bash
./scripts/test-api.sh
```

## 開発ガイドライン

### コーディング規約

- **Go標準**: `gofmt`, `go vet`, `golint`に準拠
- **命名規約**: PascalCase（公開）、camelCase（非公開）
- **パッケージ構成**: DDDアーキテクチャに従う
- **コメント**: 公開関数・型には必ずコメントを記述

### アーキテクチャ

本プロジェクトはDomain Driven Design（DDD）を採用：

```
internal/
├── production/     # 生産管理コンテキスト
├── nc/            # NC加工連携コンテキスト
├── quality/       # 品質管理コンテキスト
└── auth/          # 認証コンテキスト

各コンテキスト内:
├── domain/        # ドメインロジック（エンティティ、値オブジェクト）
├── application/   # ユースケース
├── infrastructure/ # 外部接続（DB、API）
└── interface/     # インターフェース層（HTTP、CLI）
```

### コミットメッセージ

[Conventional Commits](https://www.conventionalcommits.org/)に従う：

```
feat: 生産オーダー作成機能を追加
fix: データベース接続エラーを修正
docs: READMEにAPI仕様を追加
test: 品質管理のユニットテストを追加
refactor: 認証ミドルウェアをリファクタリング
```

## プルリクエストの作成

### 作業フロー

1. **イシューの確認**
   - 既存のissueを確認
   - 新機能の場合は事前にissueを作成して議論

2. **ブランチ作成**
```bash
git checkout -b feature/new-feature-name
# または
git checkout -b fix/bug-description
```

3. **開発・テスト**
```bash
# 開発
# ...

# テスト実行
go test ./...
make test

# APIテスト
./scripts/test-api.sh
```

4. **コミット**
```bash
git add .
git commit -m "feat: 新機能を追加"
```

5. **プッシュ・PR作成**
```bash
git push origin feature/new-feature-name
# GitHubでPRを作成
```

### PRテンプレート

PRには以下を含めてください：

- **概要**: 何を変更したか
- **変更理由**: なぜ変更が必要か
- **テスト**: どのようにテストしたか
- **影響範囲**: 既存機能への影響
- **関連Issue**: 関連するissue番号

## テストガイドライン

### テスト種類

1. **ユニットテスト**
```bash
go test ./internal/production/domain/...
```

2. **統合テスト**
```bash
go test ./internal/production/infrastructure/...
```

3. **APIテスト**
```bash
./scripts/test-api.sh
```

### テストカバレッジ

- 新機能: 80%以上のカバレッジ
- バグ修正: 該当部分のテストを追加

## コードレビュー

### レビュー観点

- **機能性**: 要件を満たしているか
- **品質**: エラーハンドリング、パフォーマンス
- **保守性**: 読みやすく、拡張しやすいか
- **セキュリティ**: セキュリティホールがないか
- **テスト**: 適切にテストされているか

### レビュー後の対応

1. フィードバックを受けたら迅速に対応
2. 不明点は遠慮なく質問
3. 修正後はレビュアーに通知

## リリースプロセス

### バージョニング

[Semantic Versioning](https://semver.org/)を採用：

- **MAJOR**: 互換性のない変更
- **MINOR**: 後方互換性のある機能追加
- **PATCH**: 後方互換性のあるバグ修正

### リリース手順

1. `develop`ブランチで変更をまとめる
2. バージョンタグを作成
3. `main`ブランチにマージ
4. GitHub Actionsで自動デプロイ

## サポートとコミュニティ

### 質問・サポート

- **GitHub Issues**: バグ報告、機能要望
- **GitHub Discussions**: 質問、アイデア共有
- **GitHub Wiki**: 詳細なドキュメント

### 行動規範

- 建設的で敬意のあるコミュニケーション
- 多様性と包括性の尊重
- オープンソース精神の体現

## 開発Tips

### よく使うコマンド

```bash
# 開発環境管理
make dev-up          # 開発環境起動
make dev-down        # 開発環境停止
make dev-logs        # ログ確認

# テスト
make test            # 全テスト実行
make test-unit       # ユニットテスト
make test-integration # 統合テスト

# ビルド
make build          # ローカルビルド
make docker-build   # Dockerビルド

# その他
make lint           # コード品質チェック
make fmt            # コードフォーマット
```

### デバッグ

```bash
# サーバー状態確認
./scripts/debug-server.sh

# データベース確認
make db-shell

# ログ確認
docker-compose -f docker-compose.dev.yml logs -f api
```

## 貢献者へのお礼

すべての貢献者の方々に心から感謝いたします。皆様の貢献により、このプロジェクトがより良いものになっています。

---

何かご不明な点がございましたら、遠慮なくissueを作成してください！