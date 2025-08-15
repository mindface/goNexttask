#!/bin/bash

echo "================================================"
echo "  Act クリーンアップスクリプト"
echo "================================================"

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "\n${YELLOW}Act関連のDockerリソースをクリーンアップします...${NC}"

# Act関連のコンテナを停止
echo "1. Act関連のコンテナを停止..."
docker ps -q --filter "name=act-" | xargs -r docker stop 2>/dev/null
echo -e "${GREEN}✅ コンテナ停止完了${NC}"

# Act関連のコンテナを削除
echo "2. Act関連のコンテナを削除..."
docker ps -aq --filter "name=act-" | xargs -r docker rm -f 2>/dev/null
echo -e "${GREEN}✅ コンテナ削除完了${NC}"

# Act関連のネットワークを削除
echo "3. Act関連のネットワークを削除..."
docker network ls --format "{{.Name}}" | grep "^act-" | xargs -r docker network rm 2>/dev/null
echo -e "${GREEN}✅ ネットワーク削除完了${NC}"

# 未使用のリソースをクリーンアップ
echo "4. 未使用のDockerリソースをクリーンアップ..."
docker network prune -f
docker container prune -f
echo -e "${GREEN}✅ クリーンアップ完了${NC}"

echo -e "\n${GREEN}================================================${NC}"
echo -e "  ${GREEN}✅ Actのクリーンアップが完了しました${NC}"
echo -e "${GREEN}================================================${NC}"
echo ""
echo "act を再実行できます:"
echo "  act push -j test"
echo "  act push -j build"