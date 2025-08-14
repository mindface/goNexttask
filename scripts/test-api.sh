#!/bin/bash

# API テストスクリプト
# 使用方法: ./scripts/test-api.sh

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# API URL
API_URL="http://localhost:8080"
TOKEN=""

echo "================================================"
echo "  GoNexttask API テストスクリプト"
echo "================================================"

# ヘルスチェック
echo -e "\n${BLUE}[Health Check]${NC}"
curl -X GET "$API_URL/health"
echo -e "\n"

# 1. ユーザー登録
echo -e "\n${BLUE}[1. User Registration]${NC}"
echo "管理者ユーザーを登録します..."

REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "password123",
    "name": "Admin User",
    "role": "admin"
  }')

echo "$REGISTER_RESPONSE" | jq '.'
TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo -e "${YELLOW}ユーザー登録に失敗した可能性があります。ログインを試みます...${NC}"
fi

# 2. ログイン
echo -e "\n${BLUE}[2. User Login]${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "password123"
  }')

echo "$LOGIN_RESPONSE" | jq '.'
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo -e "${RED}ログインに失敗しました。スクリプトを終了します。${NC}"
  exit 1
fi

echo -e "${GREEN}✅ 認証トークン取得成功${NC}"

# 3. 生産オーダー作成
echo -e "\n${BLUE}[3. Create Production Order]${NC}"
PRODUCTION_ORDER=$(curl -s -X POST "$API_URL/api/v1/production/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "orderNumber": "ORD-2024-001",
    "partId": "PART-BEARING-001",
    "quantity": 100,
    "plannedStartDate": "2024-12-15T09:00:00Z",
    "plannedEndDate": "2024-12-15T17:00:00Z",
    "machineIds": ["machine-001", "machine-002"]
  }')

echo "$PRODUCTION_ORDER" | jq '.'
ORDER_ID=$(echo "$PRODUCTION_ORDER" | jq -r '.id')

# 4. 生産オーダー一覧取得
echo -e "\n${BLUE}[4. Get All Production Orders]${NC}"
curl -s -X GET "$API_URL/api/v1/production/orders" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 5. 生産オーダー詳細取得
if [ ! -z "$ORDER_ID" ] && [ "$ORDER_ID" != "null" ]; then
  echo -e "\n${BLUE}[5. Get Production Order Detail]${NC}"
  curl -s -X GET "$API_URL/api/v1/production/orders/$ORDER_ID" \
    -H "Authorization: Bearer $TOKEN" | jq '.'

  # 6. 生産開始
  echo -e "\n${BLUE}[6. Start Production]${NC}"
  curl -s -X POST "$API_URL/api/v1/production/orders/$ORDER_ID/start" \
    -H "Authorization: Bearer $TOKEN" | jq '.'
fi

# 7. NCプログラム登録
echo -e "\n${BLUE}[7. Register NC Program]${NC}"
NC_PROGRAM=$(curl -s -X POST "$API_URL/api/v1/nc/programs" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "BEARING-001-MILLING",
    "version": "v1.0.0",
    "content": "G00 X0 Y0 Z0\nG01 X10 Y10 Z-5 F100\nG02 X20 Y0 I10 J0\nM30",
    "machineCompatibility": ["CNC-3AXIS", "CNC-5AXIS"],
    "createdBy": "admin@test.com"
  }')

echo "$NC_PROGRAM" | jq '.'
PROGRAM_ID=$(echo "$NC_PROGRAM" | jq -r '.id')

# 8. NCプログラム一覧取得
echo -e "\n${BLUE}[8. Get All NC Programs]${NC}"
curl -s -X GET "$API_URL/api/v1/nc/programs" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 9. マシンへのプログラム配置
if [ ! -z "$PROGRAM_ID" ] && [ "$PROGRAM_ID" != "null" ]; then
  echo -e "\n${BLUE}[9. Deploy NC Program to Machine]${NC}"
  curl -s -X POST "$API_URL/api/v1/nc/machines/machine-001/deploy" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
      \"programId\": \"$PROGRAM_ID\"
    }" | jq '.'
fi

# 10. マシンステータス取得
echo -e "\n${BLUE}[10. Get Machine Status]${NC}"
curl -s -X GET "$API_URL/api/v1/nc/machines/machine-001/status" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 11. 検査結果登録
echo -e "\n${BLUE}[11. Create Inspection]${NC}"
INSPECTION=$(curl -s -X POST "$API_URL/api/v1/quality/inspections" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "productionOrderId": "'"$ORDER_ID"'",
    "lotNumber": "LOT-2024-001",
    "inspectorId": "INSPECTOR-001",
    "measurements": [
      {
        "parameterName": "外径",
        "measuredValue": 50.02,
        "targetValue": 50.00,
        "tolerance": 0.05,
        "unit": "mm"
      },
      {
        "parameterName": "内径",
        "measuredValue": 30.01,
        "targetValue": 30.00,
        "tolerance": 0.03,
        "unit": "mm"
      },
      {
        "parameterName": "厚さ",
        "measuredValue": 10.00,
        "targetValue": 10.00,
        "tolerance": 0.02,
        "unit": "mm"
      }
    ]
  }')

echo "$INSPECTION" | jq '.'
INSPECTION_ID=$(echo "$INSPECTION" | jq -r '.id')

# 12. 検査結果取得
if [ ! -z "$INSPECTION_ID" ] && [ "$INSPECTION_ID" != "null" ]; then
  echo -e "\n${BLUE}[12. Get Inspection Detail]${NC}"
  curl -s -X GET "$API_URL/api/v1/quality/inspections/$INSPECTION_ID" \
    -H "Authorization: Bearer $TOKEN" | jq '.'
fi

# 13. トレーサビリティ取得
echo -e "\n${BLUE}[13. Get Traceability]${NC}"
curl -s -X GET "$API_URL/api/v1/quality/traceability?lot=LOT-2024-001" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 14. 不良分析
echo -e "\n${BLUE}[14. Defect Analysis]${NC}"
curl -s -X GET "$API_URL/api/v1/quality/defect-analysis?lot=LOT-2024-001" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 15. 生産完了
if [ ! -z "$ORDER_ID" ] && [ "$ORDER_ID" != "null" ]; then
  echo -e "\n${BLUE}[15. Complete Production]${NC}"
  curl -s -X POST "$API_URL/api/v1/production/orders/$ORDER_ID/complete" \
    -H "Authorization: Bearer $TOKEN" | jq '.'
fi

echo -e "\n================================================"
echo -e "  ${GREEN}✅ APIテスト完了${NC}"
echo "================================================"
echo ""
echo "取得した情報:"
echo "  認証トークン: ${TOKEN:0:20}..."
echo "  生産オーダーID: $ORDER_ID"
echo "  NCプログラムID: $PROGRAM_ID"
echo "  検査ID: $INSPECTION_ID"