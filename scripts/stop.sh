#!/bin/bash

set -e

echo "================================================"
echo "  GoNexttask - Docker停止スクリプト"
echo "================================================"

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 実行中のコンテナ確認
echo ""
echo "実行中のコンテナを確認しています..."

# 開発環境のコンテナ確認
DEV_RUNNING=$(docker-compose -f docker-compose.dev.yml ps -q 2>/dev/null | wc -l)

# 本番環境のコンテナ確認  
PROD_RUNNING=$(docker-compose ps -q 2>/dev/null | wc -l)

if [ $DEV_RUNNING -eq 0 ] && [ $PROD_RUNNING -eq 0 ]; then
    echo -e "${YELLOW}実行中のコンテナはありません${NC}"
    exit 0
fi

echo ""
echo "停止するサービスを選択してください:"
echo "  1) 開発環境"
echo "  2) 本番環境"
echo "  3) すべて"
echo "  4) キャンセル"
echo ""
read -p "選択 [1-4]: " choice

case $choice in
    1)
        echo -e "${YELLOW}開発環境を停止します...${NC}"
        docker-compose -f docker-compose.dev.yml down
        echo -e "${GREEN}✅ 開発環境を停止しました${NC}"
        ;;
        
    2)
        echo -e "${YELLOW}本番環境を停止します...${NC}"
        docker-compose down
        echo -e "${GREEN}✅ 本番環境を停止しました${NC}"
        ;;
        
    3)
        echo -e "${YELLOW}すべての環境を停止します...${NC}"
        docker-compose -f docker-compose.dev.yml down 2>/dev/null || true
        docker-compose down 2>/dev/null || true
        echo -e "${GREEN}✅ すべての環境を停止しました${NC}"
        ;;
        
    4)
        echo "キャンセルしました"
        exit 0
        ;;
        
    *)
        echo -e "${RED}無効な選択です${NC}"
        exit 1
        ;;
esac

echo ""
echo "データを削除しますか？"
echo "  1) いいえ (データを保持)"
echo "  2) はい (ボリュームも削除)"
echo ""
read -p "選択 [1-2]: " clean_choice

case $clean_choice in
    2)
        echo -e "${YELLOW}ボリュームを削除します...${NC}"
        docker-compose -f docker-compose.dev.yml down -v 2>/dev/null || true
        docker-compose down -v 2>/dev/null || true
        echo -e "${GREEN}✅ ボリュームを削除しました${NC}"
        ;;
esac

echo ""
echo "================================================"
echo "  停止完了"
echo "================================================"