# Docker使用ガイド

## クイックスタート

### 1. 最も簡単な起動方法

```bash
# 開発環境を起動（推奨）
./scripts/start.sh

# または Makefileを使用
make quick-start
```

### 2. 本番環境での起動

```bash
# 本番環境を起動
make quick-start-prod

# またはdocker-composeコマンド直接
docker-compose up -d
```

## 主要なコマンド

### Makefileを使用した操作

#### 開発環境
```bash
make dev-up          # 開発環境起動（ホットリロード付き）
make dev-down        # 開発環境停止
make dev-logs        # ログ表示
make dev-shell       # コンテナ内シェルアクセス
make dev-restart     # 開発環境再起動
```

#### 本番環境
```bash
make docker-up       # 本番環境起動
make docker-down     # 本番環境停止
make docker-logs     # ログ表示
make docker-restart  # 本番環境再起動
make status          # ステータス確認
```

#### データベース操作
```bash
make db-shell        # PostgreSQLシェルアクセス
make db-migrate      # マイグレーション実行
make db-backup       # DBバックアップ
make db-restore      # DBリストア
```

### docker-composeコマンド直接使用

#### 開発環境
```bash
# 起動
docker-compose -f docker-compose.dev.yml up -d

# 停止
docker-compose -f docker-compose.dev.yml down

# ログ確認
docker-compose -f docker-compose.dev.yml logs -f

# 再ビルド
docker-compose -f docker-compose.dev.yml build
```

#### 本番環境
```bash
# 起動
docker-compose up -d

# 停止
docker-compose down

# ログ確認
docker-compose logs -f api

# ステータス確認
docker-compose ps
```

## サービスアクセス

### 開発環境
- API: http://localhost:8080
- Adminer (DB管理): http://localhost:8081
  - Server: `postgres`
  - Username: `postgres`
  - Password: `password`
  - Database: `gonexttask`

### 本番環境
- API: http://localhost:8080
- Nginx (リバースプロキシ): http://localhost:80

## トラブルシューティング

### ポートが使用中の場合
```bash
# 使用中のポートを確認
lsof -i :8080
lsof -i :5432

# プロセスを終了
kill -9 <PID>
```

### コンテナが起動しない場合
```bash
# ログを確認
docker-compose logs api
docker-compose logs postgres

# コンテナを削除して再起動
docker-compose down -v
docker-compose up --build
```

### データベース接続エラー
```bash
# PostgreSQLコンテナの状態確認
docker-compose ps postgres

# データベースを再初期化
./scripts/init-db.sh
```

### ディスク容量不足
```bash
# 未使用のDockerリソースをクリーンアップ
docker system prune -a

# ボリュームも含めてクリーンアップ（データが消えるので注意）
docker system prune -a --volumes
```

## 環境変数の設定

1. `.env.example`をコピー
```bash
cp .env.example .env
```

2. `.env`ファイルを編集
```bash
# 本番環境では必ず変更
JWT_SECRET=your-production-secret-key
DB_PASSWORD=strong-password
```

## データのバックアップとリストア

### バックアップ
```bash
make db-backup
# backups/backup_YYYYMMDD_HHMMSS.sql が作成される
```

### リストア
```bash
make db-restore
# プロンプトでバックアップファイル名を入力
```

## セキュリティ注意事項

1. **本番環境では必ず環境変数を変更**
   - `JWT_SECRET`
   - `DB_PASSWORD`
   - その他の認証情報

2. **HTTPSの設定**
   - 本番環境ではNginxにSSL証明書を設定することを推奨

3. **ファイアウォール設定**
   - 必要なポートのみを開放
   - データベースポート(5432)は外部に公開しない

## よく使うワークフロー

### 1. 新機能開発時
```bash
make dev-up          # 開発環境起動
make dev-logs        # ログ監視
# コード編集（ホットリロードで自動反映）
make dev-shell       # 必要に応じてコンテナ内で作業
make dev-down        # 作業終了時
```

### 2. 本番デプロイ時
```bash
make db-backup       # 既存データのバックアップ
make docker-build    # イメージビルド
make docker-up       # サービス起動
make status          # ステータス確認
make docker-logs     # ログ確認
```

### 3. トラブル対応時
```bash
make docker-logs     # エラーログ確認
make db-shell        # DB状態確認
make docker-restart  # サービス再起動
```