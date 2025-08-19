# 製造業ナレッジ情報化システム - Seedデータ

## 概要
製造業の生産管理システムに制御理論を統合したサンプルデータとフィードバック制御システムの実装。

## ファイル構成

```
seed/
├── seed.go          # メインのseedデータ投入処理
├── feedback.go      # 品質フィードバック制御システム
├── main/
│   └── main.go      # 実行用メインファイル
└── README.md        # このファイル
```

## 実装内容

### 1. Seedデータ (seed.go)
4つのテーブルに対する実データベースのseed処理：

- **production_plans**: 10件の生産計画（自動車、半導体、医療、航空宇宙、ロボット）
- **nc_programs**: 5件のNCプログラム（実際のGコード付き）
- **inspections**: 6件の検査結果（Cpk値、業界規格準拠）
- **lot_inventory**: 22件の在庫管理（トレーサビリティ対応）

### 2. フィードバック制御 (feedback.go)

#### QualityFeedbackController
- PID制御による品質誤差の自動補正
- 検査結果からのリアルタイムフィードバック
- NCプログラムの自動調整

#### AdaptiveManufacturingSystem
- モデル予測制御(MPC)による生産最適化
- OEE向上のための適応学習
- 動的スケジューリング

## 使用方法

### 環境変数設定
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=gonexttask
```

### 実行
```bash
cd seed/main
go run main.go
```

## 制御理論の実装

### 制御ループ
```
目標値(G) → [誤差計算] → PID制御 → 補正値 → NCプログラム更新
    ↑                                           ↓
    └────────── 実測値(result) ←────────────────┘
```

### PIDゲイン設定
- Kp (比例): 0.5
- Ki (積分): 0.1
- Kd (微分): 0.05

### 目標値
- Cpk: 1.67以上
- OEE: 85%以上
- 品質合格率: 99.5%以上

## 業界標準準拠

- **自動車**: IATF 16949
- **医療機器**: FDA 21 CFR Part 820
- **航空宇宙**: AS9100D, NADCAP
- **半導体**: ISO 14644-1 (クリーンルーム)
- **ロボット**: ISO/TS 15066 (協働ロボット)

## 特徴

1. **リアルタイム適応制御**: 検査結果から即座にパラメータ調整
2. **予測保全**: 工具摩耗や熱変位の自動補正
3. **トレーサビリティ**: 完全な製造履歴追跡
4. **JIT対応**: かんばん方式での在庫最適化
5. **サステナビリティ**: リサイクル材料の管理