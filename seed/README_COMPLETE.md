# 完全データベースセットアップツール

## 概要
データベースの削除、スキーマ作成、seedデータ投入を一括で実行する統合ツールです。

## ファイル構成

```
seed/
├── complete_setup.go     # 統合セットアップ処理
├── seed.go              # 基本seedデータ（10件）
├── extended_seed.go     # 拡張seedデータ（30件）
├── feedback.go          # 品質フィードバック制御
├── main_complete/
│   └── main.go          # 実行用メインファイル
└── README_COMPLETE.md   # このファイル
```

## 機能

### 1. 完全リセット&セットアップ
- 既存テーブルの削除（CASCADE）
- 全スキーマの作成
- インデックスの作成
- 基本seedデータ投入（10件）
- 拡張seedデータ投入（30件）

### 2. 作成されるテーブル

#### 認証系
- `users` - ユーザー認証

#### 生産管理系
- `production_orders` - 生産オーダー（DDD設計）
- `production_plans` - 生産計画（seedデータ用）

#### NC加工系
- `nc_programs` - NCプログラム管理
- `machines` - 機械管理

#### 品質管理系
- `inspections` - 検査結果
- `measurement_results` - 測定結果詳細
- `quality_adjustments` - 品質調整履歴
- `quality_alerts` - 品質アラート

#### 在庫管理系
- `lot_inventory` - ロット在庫
- `purchase_orders` - 購入注文

#### システム系
- `schema_migrations` - マイグレーション管理

## 使用方法

### 基本実行（全データリセット&投入）
```bash
docker compose -f docker-compose.dev.yml exec api sh -c "go run -mod=mod seed/main_complete/main.go"
```

### オプション付き実行
```bash
# 詳細ログ表示
docker compose -f docker-compose.dev.yml exec api sh -c "go run -mod=mod seed/main_complete/main.go -v"

# リセットなし（テーブル作成とデータ投入のみ）
docker compose -f docker-compose.dev.yml exec api sh -c "go run -mod=mod seed/main_complete/main.go -reset=false"

# 基本データのみ
docker compose -f docker-compose.dev.yml exec api sh -c "go run -mod=mod seed/main_complete/main.go -basic -extend=false"
```

## 投入されるデータ

### 基本データ（seed.go）
- 生産計画: 10件
- NCプログラム: 5件
- 検査結果: 6件
- 在庫管理: 22件

### 拡張データ（extended_seed.go）
業界別に各5件、合計30件：

1. **自動車部品製造**
   - トランスミッションギア
   - クランクシャフト
   - ブレーキディスク
   - ターボチャージャー
   - EVモーターシャフト

2. **半導体製造装置**
   - ウェハステージ
   - EUVマスクホルダー
   - プラズマチャンバー
   - 真空チャック
   - イオン注入部品

3. **医療機器**
   - 人工股関節
   - 脊椎インプラント
   - 歯科インプラント
   - 血管ステント
   - 手術器具

4. **航空宇宙**
   - タービンブレード
   - 主翼リブ
   - ロケット燃焼室
   - 衛星構造
   - ローターハブ

5. **協働ロボット**
   - ロボットアーム関節
   - 力覚センサー
   - ビジョンマウント
   - 安全グリッパー
   - AIコントローラー

## 制御理論統合

各データには以下の制御要素が含まれます：
- **PID制御パラメータ**: Kp=0.5, Ki=0.1, Kd=0.05
- **適応制御**: リアルタイムフィードバック
- **品質予測**: Cpk値モニタリング
- **安全制御**: ISO/TS 15066準拠（ロボット）

## トラブルシューティング

### エラー: "relation does not exist"
```bash
# 完全セットアップを実行
docker compose -f docker-compose.dev.yml exec api sh -c "go run -mod=mod seed/main_complete/main.go -reset"
```

### エラー: "duplicate key"
```bash
# リセットフラグ付きで実行
docker compose -f docker-compose.dev.yml exec api sh -c "go run -mod=mod seed/main_complete/main.go -reset"
```

### データ確認
```bash
# Adminerで確認
http://localhost:8081

# またはCLIで確認
docker compose -f docker-compose.dev.yml exec postgres psql -U postgres -d gonexttask -c "\dt"
```

## 実行例の出力

```
========================================
  Complete Database Setup Tool
========================================
Database: postgres@postgres:5432/gonexttask
Options: reset=true, basic=true, extend=true
----------------------------------------
✅ Connected to database successfully
----------------------------------------
Dropping all tables...
Creating all tables...
Inserting basic seed data...
Inserting extended seed data...
----------------------------------------

📊 Database Statistics
----------------------
👤 users                : 0 records
📋 production_plans     : 35 records
📦 production_orders    : 0 records
🔧 nc_programs          : 8 records
🏭 machines             : 0 records
🔍 inspections          : 8 records
📦 lot_inventory        : 22 records

🏭 Industry Breakdown
---------------------
🚗 Automotive          : 6 records
🔬 Semiconductor       : 5 records
🏥 Medical Device      : 5 records
✈️ Aerospace           : 5 records
🤖 Robotics            : 5 records

📈 Quality Metrics
-----------------
Total Inspections: 8
Pass Rate: 87.5% (7 passed, 1 failed)
Average Cpk: 1.73

========================================
✅ Complete setup finished in 2.3s
========================================
```