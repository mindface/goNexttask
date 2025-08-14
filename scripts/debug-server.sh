#!/bin/bash

echo "================================================"
echo "  サーバー状態確認スクリプト"
echo "================================================"

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "\n${YELLOW}1. Docker コンテナ状態確認${NC}"
echo "開発環境のコンテナ:"
docker-compose -f docker-compose.dev.yml ps

echo -e "\n本番環境のコンテナ:"
docker-compose ps

echo -e "\n${YELLOW}2. ポート使用状況確認${NC}"
echo "ポート 8080:"
lsof -i :8080 || echo "ポート 8080 は使用されていません"

echo "ポート 5432:"
lsof -i :5432 || echo "ポート 5432 は使用されていません"

echo -e "\n${YELLOW}3. APIヘルスチェック${NC}"
curl -s http://localhost:8080/health && echo -e "\n${GREEN}✅ API応答OK${NC}" || echo -e "\n${RED}❌ API応答なし${NC}"

echo -e "\n${YELLOW}4. データベース接続確認${NC}"
if docker-compose -f docker-compose.dev.yml ps postgres | grep -q "Up"; then
    echo "PostgreSQL コンテナが起動中..."
    docker-compose -f docker-compose.dev.yml exec postgres pg_isready -U postgres && echo -e "${GREEN}✅ DB接続OK${NC}" || echo -e "${RED}❌ DB接続NG${NC}"
else
    echo -e "${RED}PostgreSQL コンテナが起動していません${NC}"
fi

echo -e "\n${YELLOW}5. API コンテナのログ確認${NC}"
echo "最新10行のAPIログ:"
docker-compose -f docker-compose.dev.yml logs --tail=10 api 2>/dev/null || echo "APIコンテナが見つかりません"

echo -e "\n${YELLOW}6. 解決手順の提案${NC}"
if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${RED}APIサーバーにアクセスできません。以下を試してください:${NC}"
    echo ""
    echo "1. サーバーを起動:"
    echo "   make dev-up"
    echo ""
    echo "2. サーバーを再起動:"
    echo "   make dev-restart"
    echo ""
    echo "3. 完全にリセット:"
    echo "   make dev-down"
    echo "   docker system prune -f"
    echo "   make dev-up"
    echo ""
    echo "4. ログの詳細確認:"
    echo "   make dev-logs"
fi