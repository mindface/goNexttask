# GitHub セットアップガイド

## CI/CDパイプライン設定

### 1. GitHub Container Registry (推奨)
デフォルトの`ci.yml`を使用。MY_GITHUB_TOKENのみ必要：

```bash
# Personal Access Tokenを作成（ghcr.ioへの書き込み権限）
gh auth token | pbcopy
```

GitHubリポジトリ設定:
- Settings → Secrets and variables → Actions
- `MY_GITHUB_TOKEN`: Personal Access Token

### 2. Docker Hub
`ci-dockerhub.yml`を使用する場合：

GitHubリポジトリ設定:
- Settings → Secrets and variables → Actions  
- `DOCKER_USERNAME`: Docker Hubユーザー名
- `DOCKER_PASSWORD`: Docker Hubパスワード/トークン

### 3. ローカルテスト用
認証不要の`ci-local.yml`でテスト可能。

## ワークフロー説明

### メインワークフロー (ci.yml)
- **test**: PostgreSQLサービスでテスト実行
- **build**: Goアプリケーションビルド  
- **docker**: Dockerイメージビルド・プッシュ
- **deploy**: ステージング/本番デプロイ（プレースホルダー）

### 対象ブランチ
- `main`: 本番環境デプロイ
- `develop`: ステージング環境デプロイ
- Pull Request: テスト・ビルドのみ

## セキュリティ
- gosecスキャンは一時無効化（act実行時の互換性のため）
- MY_GITHUB_TOKENは必要最小限の権限で設定
- Dockerイメージは非rootユーザーで実行

## 手動デプロイ
```bash
# ローカルテスト
act --job test

# 手動ワークフロー実行
gh workflow run ci-local.yml
```