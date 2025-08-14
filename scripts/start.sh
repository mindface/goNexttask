#!/bin/bash

set -e

echo "================================================"
echo "  GoNexttask - Docker起動スクリプト"
echo "================================================"

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 環境選択
echo ""
echo "起動環境を選択してください:"
echo "  1) 開発環境 (ホットリロード付き)"
echo "  2) 本番環境"
echo "  3) 終了"
echo ""
read -p "選択 [1-3]: " choice

case $choice in
    1)
        echo -e "${GREEN}開発環境を起動します...${NC}"
        echo ""
        
        # .envファイルのチェック
        if [ ! -f .env ]; then
            echo -e "${YELLOW}.envファイルが見つかりません。.env.exampleからコピーします...${NC}"
            cp .env.example .env
            echo -e "${GREEN}.envファイルを作成しました。必要に応じて設定を変更してください。${NC}"
        fi
        
        # 開発環境起動
        echo "Dockerコンテナをビルドしています..."
        docker-compose -f docker-compose.dev.yml build
        
        echo ""
        echo "サービスを起動しています..."
        docker-compose -f docker-compose.dev.yml up -d
        
        echo ""
        echo -e "${GREEN}✅ 開発環境が起動しました！${NC}"
        echo ""
        echo "📌 アクセスURL:"
        echo "   API:     http://localhost:8080"
        echo "   Adminer: http://localhost:8081"
        echo ""
        echo "📝 便利なコマンド:"
        echo "   ログ確認:     make dev-logs"
        echo "   コンテナ停止: make dev-down"
        echo "   DB接続:       make db-shell"
        ;;
        
    2)
        echo -e "${GREEN}本番環境を起動します...${NC}"
        echo ""
        
        # .envファイルのチェック
        if [ ! -f .env ]; then
            echo -e "${RED}エラー: .envファイルが見つかりません${NC}"
            echo "本番環境では.envファイルの設定が必要です。"
            echo "cp .env.example .env を実行して、適切な値を設定してください。"
            exit 1
        fi
        
        # 本番環境起動
        echo "Dockerイメージをビルドしています..."
        docker-compose build
        
        echo ""
        echo "サービスを起動しています..."
        docker-compose up -d
        
        echo ""
        echo -e "${GREEN}✅ 本番環境が起動しました！${NC}"
        echo ""
        echo "📌 アクセスURL:"
        echo "   API: http://localhost:8080"
        echo ""
        echo "📝 便利なコマンド:"
        echo "   ログ確認:     make docker-logs"
        echo "   コンテナ停止: make docker-down"
        echo "   ステータス:   make status"
        ;;
        
    3)
        echo "終了します"
        exit 0
        ;;
        
    *)
        echo -e "${RED}無効な選択です${NC}"
        exit 1
        ;;
esac

# ヘルスチェック
echo ""
echo "ヘルスチェック中..."
sleep 5

max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ APIサーバーが正常に起動しました${NC}"
        break
    fi
    
    attempt=$((attempt + 1))
    if [ $attempt -eq $max_attempts ]; then
        echo -e "${RED}⚠️  APIサーバーの起動に失敗した可能性があります${NC}"
        echo "ログを確認してください: docker-compose logs api"
    else
        echo -n "."
        sleep 2
    fi
done

echo ""
echo "================================================"
echo "  起動完了"
echo "================================================"