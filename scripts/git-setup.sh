#!/bin/bash

set -e

echo "================================================"
echo "  Git セットアップ & GitHub Push スクリプト"
echo "================================================"

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 現在のディレクトリ確認
if [ ! -f "go.mod" ]; then
  echo -e "${RED}エラー: go.modが見つかりません。プロジェクトルートディレクトリで実行してください。${NC}"
  exit 1
fi

echo -e "${BLUE}[1. Git設定確認]${NC}"
# Git設定確認
if ! git config --global user.name > /dev/null 2>&1; then
  echo -e "${YELLOW}Gitユーザー名が設定されていません${NC}"
  read -p "Git ユーザー名を入力: " git_username
  git config --global user.name "$git_username"
fi

if ! git config --global user.email > /dev/null 2>&1; then
  echo -e "${YELLOW}Gitメールアドレスが設定されていません${NC}"
  read -p "Git メールアドレスを入力: " git_email
  git config --global user.email "$git_email"
fi

echo -e "${GREEN}Git設定:${NC}"
echo "  ユーザー名: $(git config --global user.name)"
echo "  メール: $(git config --global user.email)"

echo -e "\n${BLUE}[2. Gitリポジトリ初期化]${NC}"
# Gitリポジトリが既に存在するかチェック
if [ ! -d ".git" ]; then
  echo "Gitリポジトリを初期化しています..."
  git init
  echo -e "${GREEN}✅ Gitリポジトリを初期化しました${NC}"
else
  echo -e "${GREEN}✅ Gitリポジトリは既に初期化されています${NC}"
fi

echo -e "\n${BLUE}[3. .gitignore確認]${NC}"
if [ ! -f ".gitignore" ]; then
  echo -e "${RED}エラー: .gitignoreファイルが見つかりません${NC}"
  exit 1
fi
echo -e "${GREEN}✅ .gitignoreファイルが存在します${NC}"

echo -e "\n${BLUE}[4. vendorディレクトリのクリーンアップ]${NC}"
if [ -d "vendor" ]; then
  echo -e "${YELLOW}vendorディレクトリを削除しています...${NC}"
  rm -rf vendor
  echo -e "${GREEN}✅ vendorディレクトリを削除しました${NC}"
fi

echo -e "\n${BLUE}[5. 一時ファイルのクリーンアップ]${NC}"
# 一時ファイルやビルドファイルを削除
rm -f build-errors.log
rm -rf tmp/
rm -rf bin/
echo -e "${GREEN}✅ 一時ファイルをクリーンアップしました${NC}"

echo -e "\n${BLUE}[6. ファイル追加とコミット]${NC}"
# ファイルをステージング
git add .

# コミット（変更がある場合のみ）
if git diff --staged --quiet; then
  echo -e "${YELLOW}コミットする変更はありません${NC}"
else
  echo -e "${GREEN}変更をコミットしています...${NC}"
  git commit -m "初回コミット: GoNexttask DDD architecture implementation

- DDDアーキテクチャに基づいた構造
- Production, NC, Quality の3つのコンテキスト
- JWT認証システム
- PostgreSQLデータベース統合
- Docker開発環境
- GitHub Actions CI/CD
- API テストスクリプト

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>"
  echo -e "${GREEN}✅ 変更をコミットしました${NC}"
fi

echo -e "\n${BLUE}[7. GitHubリモートリポジトリ設定]${NC}"
echo "GitHubでリポジトリを作成してから以下を実行してください："
echo ""
echo -e "${YELLOW}手順:${NC}"
echo "1. https://github.com/new でリポジトリを作成"
echo "2. リポジトリ名: goNexttask （推奨）"
echo "3. 作成後、以下のコマンドを実行:"
echo ""
echo -e "${GREEN}git remote add origin https://github.com/YOUR_USERNAME/REPOSITORY_NAME.git${NC}"
echo -e "${GREEN}git branch -M main${NC}"
echo -e "${GREEN}git push -u origin main${NC}"
echo ""
echo "または、このスクリプトを続行してリモートURLを設定できます。"
echo ""

read -p "GitHubリモートURLを設定しますか？ (y/n): " setup_remote

if [ "$setup_remote" = "y" ] || [ "$setup_remote" = "Y" ]; then
  read -p "GitHubユーザー名: " github_username
  read -p "リポジトリ名 [goNexttask]: " repo_name
  repo_name=${repo_name:-goNexttask}
  
  remote_url="https://github.com/${github_username}/${repo_name}.git"
  
  # リモートが既に存在する場合は更新
  if git remote get-url origin > /dev/null 2>&1; then
    git remote set-url origin "$remote_url"
  else
    git remote add origin "$remote_url"
  fi
  
  # メインブランチに設定
  git branch -M main
  
  echo -e "\n${BLUE}[8. GitHubにプッシュ]${NC}"
  echo "GitHubにプッシュしています..."
  
  if git push -u origin main; then
    echo -e "\n${GREEN}🎉 成功！GitHubリポジトリにプッシュしました！${NC}"
    echo ""
    echo "リポジトリURL: https://github.com/${github_username}/${repo_name}"
    echo ""
    echo -e "${BLUE}次のステップ:${NC}"
    echo "1. GitHub Actionsが自動実行されます"
    echo "2. READMEを確認してプロジェクトの概要を把握"
    echo "3. Issuesで機能要望・バグ報告を管理"
    echo "4. Projectsでタスク管理"
    echo ""
    echo -e "${BLUE}開発の続行:${NC}"
    echo "  make dev-up     # 開発環境起動"
    echo "  ./scripts/test-api.sh  # API テスト"
  else
    echo -e "\n${RED}プッシュに失敗しました${NC}"
    echo "以下を確認してください："
    echo "1. GitHubリポジトリが作成されているか"
    echo "2. アクセス権限があるか"
    echo "3. リモートURLが正しいか"
    echo ""
    echo "手動でプッシュする場合："
    echo "  git push -u origin main"
  fi
else
  echo -e "\n${YELLOW}リモートリポジトリの設定をスキップしました${NC}"
  echo "後で設定する場合は以下を実行してください："
  echo ""
  echo "  git remote add origin https://github.com/YOUR_USERNAME/REPOSITORY_NAME.git"
  echo "  git branch -M main"
  echo "  git push -u origin main"
fi

echo -e "\n${BLUE}[9. 最終確認]${NC}"
echo "Git status:"
git status --short

if git remote get-url origin > /dev/null 2>&1; then
  echo ""
  echo "Remote URL: $(git remote get-url origin)"
fi

echo ""
echo "================================================"
echo -e "  ${GREEN}✅ Git セットアップ完了${NC}"
echo "================================================"