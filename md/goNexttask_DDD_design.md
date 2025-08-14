# goNexttask — DDD 設計書

## 1. 概要

本設計書は、ベアリング製造／金属加工工場向けの統合情報システム（以下、本システム）を、ドメイン駆動設計（DDD）を中心に設計するための詳細設計書です。目的は生産管理・NC加工連携・品質管理を統合し、トレーサビリティ、可用性、リアルタイム性、安全性を確保することです。

### 目的と背景（要約）

- 現行は情報分断と手作業が多く、効率低下・誤り・トレーサビリティ不備を招いている。
- 新システムは生産計画の自動生成、NC加工機との双方向連携、測定データの自動収集・紐付け等により業務改善を目指す。

## 2. 境界づけられたコンテキスト（Bounded Contexts）

システムは業務上の関心ごとで3つの主要サブドメインに分離します。

1. **Production（生産管理）**
   - 生産計画作成、進捗管理、工程実行指示、在庫/ロット管理。
2. **NCIntegration（NC加工連携）**
   - NCプログラム管理、バージョン管理、機械への転送、稼働・エラー情報の収集。
3. **Quality（品質管理）**
   - 測定データの取り込み、検査実行、合否判定、不良原因の分析、トレーサビリティ照会。

各コンテキストは明確な境界（API契約、イベント）で接続します。責務は極力分離し、ドメインモデル中心で設計します。

## 3. 高レベルアーキテクチャ

- マイクロサービス化（将来のスケーラビリティを考慮）またはモノリシック内部に明確なパッケージ境界を確保。
- 通信: HTTP/REST（同期API） + メッセージバス（イベント駆動、例: Kafka/RabbitMQ）による非同期連携。
- データベース: 各コンテキストは独立した永続化を持つ（DB per Bounded Context）。
- インフラ: コンテナ (Docker) + オーケストレーション (Kubernetes)。

## 4. パッケージ構成（Go / internalレイアウト）

```
goNexttask/
├── cmd/api/main.go
├── internal/
│   ├── production/
│   │   ├── domain/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   └── interface/
│   ├── nc/
│   │   └── (同様)
│   └── quality/
│       └── (同様)
├── pkg/ (共通ライブラリ)
└── configs/
```

- `domain/`: エンティティ、値オブジェクト、ドメインサービス、リポジトリインターフェース、ドメインイベント
- `application/`: ユースケース（アプリケーションサービス）、DTO、トランザクションテンプレート
- `infrastructure/`: DBリポジトリ実装、外部APIクライアント、NCデバイスコネクタ
- `interface/`: HTTPハンドラー、GRPC、CLI、ジョブ（バッチ）

## 5. ドメインモデル（主要エンティティ・VO・イベント）

### Production（生産管理）

- **Entity: ProductionOrder**
  - id: ProductionOrderID
  - orderNumber: string (受注番号)
  - partId: PartID
  - quantity: int
  - status: enum (Planned, InProgress, Completed, Cancelled)
  - schedule: Schedule (VO)
  - createdAt, updatedAt
- **ValueObject: Schedule**
  - plannedStart, plannedEnd, assignedMachines []MachineID
- **DomainService: ProductionSchedulingService**
  - 生産計画自動生成ロジック（受注 + リソース可用性 + NC機の稼働状況）
- **DomainEvent**
  - `ProductionOrderCreated`, `ProductionOrderStarted`, `ProductionOrderDelayed`, `ProductionOrderCompleted`

### NCIntegration（NC加工連携）

- **Entity: NCProgram**
  - id, name, version, fileHash, machineCompatibility, createdBy, createdAt
- **Entity: Machine**
  - id, name, ip, type, capabilities, status: MachineStatus (VO)
- **ValueObject: MachineStatus**
  - runningState (Running/Stopped/Error), currentJobID, lastHeartbeat
- **DomainService: NCTransferService**
  - 適切なプログラム選定、転送制御、結果受信の整合性確保
- **DomainEvent**
  - `NCProgramDeployed`, `MachineStatusChanged`, `NCJobCompleted`, `NCJobError`

### Quality（品質管理）

- **Entity: Inspection**
  - id, productionOrderId, lotNumber, inspectorId, results[], status
- **ValueObject: Measurement**
  - dimensions map[string]float64, instrumentId, measuredAt
- **DomainService: DefectAnalysisService**
  - 不良原因の関連付け（NCプログラム版、工具、機械ログ）
- **DomainEvent**
  - `InspectionCompleted`, `DefectDetected`, `MeasurementRecorded`

## 6. ユースケース（代表的）

- 受注 -> 生産計画作成（Production）
- 生産開始（オペレーター操作 or 自動） -> NCプログラム選定・転送（Production ⇄ NCIntegration）
- 加工完了通知 -> 進捗更新 + 測定ジョブ発行（NCIntegration → Quality）
- 測定データ取り込み -> 判定・トレーサビリティ更新（Quality）
- 遅延検知 -> アラート発信（Production）

## 7. API設計（代表的エンドポイント）

### Production HTTP API

- `POST /api/v1/production/orders` - 生産オーダー作成
- `GET /api/v1/production/orders/{id}` - オーダー取得
- `POST /api/v1/production/orders/{id}/start` - 開始指示
- `GET /api/v1/production/orders/{id}/progress` - 進捗取得

### NCIntegration HTTP API

- `POST /api/v1/nc/programs` - NCプログラム登録
- `POST /api/v1/nc/machines/{id}/deploy` - プログラム展開（転送）
- `POST /api/v1/nc/machines/{id}/status` - マシンステータス受信（NCプッシュ）

### Quality HTTP API

- `POST /api/v1/quality/inspections` - 検査登録（測定データ含む）
- `GET /api/v1/quality/traceability?lot={}` - ロット追跡

### 非同期イベント（メッセージ）

- `ProductionOrderStarted` -> 発行元: Production
- `NCJobCompleted {machineId, jobId, producedCount}` -> 発行元: NCIntegration
- `MeasurementRecorded {inspectionId, measurements...}` -> 発行元: Quality

APIはOpenAPI(spec)で定義し、契約ドキュメントを用意します。

## 8. 永続化設計（概略）

- 各コンテキストは独自DB（RDBMS: PostgreSQL 推奨）を持つ。
- イベントストアが必要ならば、イベントログを別DBに保存（Event Sourcingを採用する場合）。
- テーブル設計（例: production.orders, nc.programs, quality.inspections）
- 重要: `NCProgram.file` はファイルストレージ（S3等）へ格納し、DBにはメタ情報とハッシュを保持。

## 9. インテグレーション（NC機器との接続）

- NC機はメーカー／モデルが複数想定されるため、**アダプターパターン**で抽象化。
- 通信方式は以下を想定し、機種ごとにコネクタを実装:
  - FTP/SFTP（プログラム転送）
  - MTConnect / OPC-UA（稼働・測定データ）
  - プロプライエタリなTCP/HTTP API
- コネクタは `internal/nc/infrastructure/connector/*` に実装。
- 冪等性と再送ポリシーを設計（ファイルハッシュ、トランザクションIDを利用）。

## 10. 非機能要件対応

- **可用性**: 99.5%目標。冗長構成、DBレプリケーション、自動フェールオーバー、定期バックアップ。RPO/RTOを定義。
- **性能**: リアルタイム表示遅延 ≤5秒、NC転送 ≤5秒（機材依存あり）。インメモリキャッシュ（Redis）、WebSocket/Server-Sent Eventsでリアルタイム配信。
- **セキュリティ**: OAuth2/JWTによる認証、RBACによる権限制御、監査ログの保存、TLS通信、ファイルの署名検証。
- **保守性**: 明確なパッケージ分離、テストカバレッジ（ユニット・集約テスト）、CI/CDパイプライン。
- **拡張性**: プラグイン化されたNCコネクタ、新しい測定器ドライバを追加しやすいアーキテクチャ。

## 11. テスト戦略

- **ユニットテスト**: ドメインロジック中心（mock不要でテスト可能なVO・エンティティ）
- **ドメインテスト**: ドメインサービスとイベント発行を検証
- **インテグレーションテスト**: DB、メッセージバス、外部コネクタの擬似環境
- **E2Eテスト**: 代表的なシナリオ（受注→NC転送→加工→測定→トレーサビリティ）
- **負荷試験**: リアルタイム性要件達成の確認

## 12. 運用・監視

- **ログ**: 構造化ログ（JSON）、トレーシング（OpenTelemetry）
- **監視**: メトリクス収集（Prometheus）、アラート（Alertmanager）、ダッシュボード（Grafana）
- **運用プレイブック**: 障害時の手順（ログ収集、再送、差分復元）

## 13. セキュリティ & ガバナンス

- 最小権限の原則、機密情報はVaultで管理
- 操作ログの保存（全操作のユーザーID、タイムスタンプ、操作内容）
- データ整合性: 署名・ハッシュ検証（NCプログラム、測定データ）

## 14. CI/CD / デプロイ戦略

- GitHub Actions / GitLab CI でビルド・テスト・イメージ作成
- コンテナレジストリへ push → Kubernetes (Helm) へデプロイ
- Blue/Green または Canary リリース戦略

## 15. マイグレーション計画（既存システムからの移行）

1. フェーズ1: Read-only 統合とデータ収集（非破壊）
2. フェーズ2: 部分機能切替（例: NCプログラムの一部転送自動化）
3. フェーズ3: フル移行（生産計画の自動化）
4. ロールバック手順とデータコンシステンシー検証を必ず用意

## 16. ネーミング・命名規約（Go向け）

- パッケージ名: 小文字単数形（`production`, `nc`, `quality`）
- ディレクトリ: `domain`, `application`, `infrastructure`, `interface`
- エンティティ: `ProductionOrder`, `NCProgram` (UpperCamelCase)
- 値オブジェクト: `MachineStatus`, `Schedule`
- リポジトリインターフェース: `ProductionOrderRepository` (packageの中でインターフェース定義)
- 実装: `postgresProductionOrderRepository` または `productionOrderRepositoryPG`（構造体名は lowerCamel ではなく Go の公開/非公開規則に従う）
- ファイル: `production_order.go`, `production_service.go` などスネークケースではなく小文字アンダースコアで分割

## 17. サンプルコードテンプレート（抜粋）

**domain/entity: production_order.go**
