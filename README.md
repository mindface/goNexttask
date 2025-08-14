# GoNexttask - ベアリング製造・金属加工向け統合管理システム

[![CI/CD Pipeline](https://github.com/YOUR_USERNAME/goNexttask/workflows/CI/CD%20Pipeline/badge.svg)](https://github.com/YOUR_USERNAME/goNexttask/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_USERNAME/goNexttask)](https://goreportcard.com/report/github.com/YOUR_USERNAME/goNexttask)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](https://www.docker.com/)

DDDアーキテクチャに基づいた生産管理、NC加工連携、品質管理システム

## 🚀 特徴

- **Domain Driven Design**: 明確に分離された3つのコンテキスト
- **マイクロサービス対応**: 独立してスケール可能な設計
- **リアルタイム連携**: NC機器との双方向通信
- **完全トレーサビリティ**: 原材料から完成品まで追跡
- **高可用性**: 99.5%稼働率を目標とした設計
- **セキュア**: JWT認証、RBAC、監査ログ

## プロジェクト構造

```
goNexttask/
├── cmd/api/              # アプリケーションエントリーポイント
├── internal/             # プライベートアプリケーションコード
│   ├── production/       # 生産管理コンテキスト
│   ├── nc/              # NC加工連携コンテキスト
│   ├── quality/         # 品質管理コンテキスト
│   └── auth/            # 認証
├── pkg/                 # 共有ライブラリ
├── migrations/          # データベースマイグレーション
└── configs/             # 設定ファイル
```

## セットアップ

1. 環境変数の設定
```bash
cp .env.example .env
# .envファイルを編集して環境に合わせた値を設定
```

2. データベースのセットアップ
```bash
# PostgreSQLを起動
# マイグレーションを実行
psql -U postgres -d gonexttask < migrations/001_create_tables.sql
```

3. 依存関係のインストール
```bash
go mod download
```

4. アプリケーションの起動
```bash
go run cmd/api/main.go
```

## API エンドポイント

### 認証
- `POST /api/v1/auth/register` - ユーザー登録
- `POST /api/v1/auth/login` - ログイン

### 生産管理
- `POST /api/v1/production/orders` - 生産オーダー作成
- `GET /api/v1/production/orders` - 生産オーダー一覧
- `GET /api/v1/production/orders/{id}` - 生産オーダー詳細
- `POST /api/v1/production/orders/{id}/start` - 生産開始
- `POST /api/v1/production/orders/{id}/complete` - 生産完了

### NC加工連携
- `POST /api/v1/nc/programs` - NCプログラム登録
- `GET /api/v1/nc/programs` - NCプログラム一覧
- `POST /api/v1/nc/machines/{id}/deploy` - プログラム転送
- `GET /api/v1/nc/machines/{id}/status` - マシンステータス取得

### 品質管理
- `POST /api/v1/quality/inspections` - 検査結果登録
- `GET /api/v1/quality/inspections/{id}` - 検査結果詳細
- `GET /api/v1/quality/traceability?lot={lotNumber}` - トレーサビリティ照会
- `GET /api/v1/quality/defect-analysis?lot={lotNumber}` - 不良分析

## 開発

### テスト実行
```bash
go test ./...
```

### ビルド
```bash
go build -o bin/api cmd/api/main.go
```