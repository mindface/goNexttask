# API テストコマンド集

## 前提条件
- サーバーが起動していること（`http://localhost:8080`）
- jqコマンドがインストールされていること（`brew install jq`）

## 一括テスト実行
```bash
# スクリプトで全APIをテスト
chmod +x scripts/test-api.sh
./scripts/test-api.sh
```

## 個別APIテストコマンド

### 1. 認証 (Authentication)

#### ユーザー登録
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "password",
    "name": "Admin User",
    "role": "admin"
  }' | jq '.'
```

#### ログイン
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "password"
  }' | jq '.'
```

トークンを変数に保存：
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "password123"
  }' | jq -r '.token')

echo "Token: $TOKEN"
```

### 2. 生産管理 (Production)

#### 生産オーダー作成
```bash
curl -X POST http://localhost:8080/api/v1/production/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "orderNumber": "ORD-2024-001",
    "partId": "PART-BEARING-001",
    "quantity": 100,
    "plannedStartDate": "2024-12-15T09:00:00Z",
    "plannedEndDate": "2024-12-15T17:00:00Z",
    "machineIds": ["machine-001", "machine-002"]
  }' | jq '.'
```

#### 生産オーダー一覧取得
```bash
curl -X GET http://localhost:8080/api/v1/production/orders \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

#### 生産オーダー詳細取得
```bash
ORDER_ID="order-ORD-2024-001"
curl -X GET http://localhost:8080/api/v1/production/orders/$ORDER_ID \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

#### 生産開始
```bash
curl -X POST http://localhost:8080/api/v1/production/orders/$ORDER_ID/start \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

#### 生産完了
```bash
curl -X POST http://localhost:8080/api/v1/production/orders/$ORDER_ID/complete \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

### 3. NC加工連携 (NC Integration)

#### NCプログラム登録
```bash
curl -X POST http://localhost:8080/api/v1/nc/programs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "BEARING-001-MILLING",
    "version": "v1.0.0",
    "content": "G00 X0 Y0 Z0\nG01 X10 Y10 Z-5 F100\nG02 X20 Y0 I10 J0\nM30",
    "machineCompatibility": ["CNC-3AXIS", "CNC-5AXIS"],
    "createdBy": "admin@test.com"
  }' | jq '.'
```

#### NCプログラム一覧取得
```bash
curl -X GET http://localhost:8080/api/v1/nc/programs \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

#### プログラムをマシンに配置
```bash
PROGRAM_ID="ncprog-12345678"
curl -X POST http://localhost:8080/api/v1/nc/machines/machine-001/deploy \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "programId": "'$PROGRAM_ID'"
  }' | jq '.'
```

#### マシンステータス取得
```bash
curl -X GET http://localhost:8080/api/v1/nc/machines/machine-001/status \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

#### マシンステータス更新
```bash
curl -X POST http://localhost:8080/api/v1/nc/machines/machine-001/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "runningState": "running",
    "currentJobId": "job-001"
  }' | jq '.'
```

### 4. 品質管理 (Quality)

#### 検査結果登録
```bash
curl -X POST http://localhost:8080/api/v1/quality/inspections \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "productionOrderId": "order-ORD-2024-001",
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
  }' | jq '.'
```

#### 検査結果取得
```bash
INSPECTION_ID="insp-20241214120000"
curl -X GET http://localhost:8080/api/v1/quality/inspections/$INSPECTION_ID \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

#### トレーサビリティ照会
```bash
curl -X GET "http://localhost:8080/api/v1/quality/traceability?lot=LOT-2024-001" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

#### 不良分析
```bash
curl -X GET "http://localhost:8080/api/v1/quality/defect-analysis?lot=LOT-2024-001" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

## HTTPieを使用したテスト

HTTPieをインストール（`brew install httpie`）していれば、より見やすい形式でテストできます：

```bash
# ログイン
http POST localhost:8080/api/v1/auth/login \
  email=admin@test.com \
  password=password123

# 生産オーダー作成（トークン付き）
http POST localhost:8080/api/v1/production/orders \
  "Authorization: Bearer $TOKEN" \
  orderNumber=ORD-2024-002 \
  partId=PART-002 \
  quantity:=200 \
  plannedStartDate=2024-12-16T09:00:00Z \
  plannedEndDate=2024-12-16T17:00:00Z \
  machineIds:='["machine-001"]'
```

## Postmanコレクション

Postmanを使用する場合は、以下の環境変数を設定：

- `BASE_URL`: `http://localhost:8080`
- `TOKEN`: ログイン後に取得したトークン

## テストデータのリセット

```bash
# データベースをリセット
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up -d

# 初期データを投入
./scripts/init-db.sh
```

## エラーレスポンスのテスト

```bash
# 認証なしでアクセス（401エラー）
curl -X GET http://localhost:8080/api/v1/production/orders

# 不正なデータ（400エラー）
curl -X POST http://localhost:8080/api/v1/production/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "quantity": -1
  }'

# 存在しないリソース（404エラー）
curl -X GET http://localhost:8080/api/v1/production/orders/invalid-id \
  -H "Authorization: Bearer $TOKEN"
```