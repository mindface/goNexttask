# DBスキーマ（DDL） — Go Nexttask

```sql
-- 生産計画テーブル
CREATE TABLE production_plans (
    id VARCHAR(64) PRIMARY KEY,
    order_id VARCHAR(64) NOT NULL,
    material VARCHAR(128) NOT NULL,
    quantity INT NOT NULL,
    status VARCHAR(32) NOT NULL CHECK (status IN ('planned', 'in_progress', 'completed', 'delayed')),
    scheduled_start_date TIMESTAMP NOT NULL,
    scheduled_end_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- NCプログラムテーブル
CREATE TABLE nc_programs (
    id VARCHAR(64) PRIMARY KEY,
    part_id VARCHAR(64) NOT NULL,
    machine_id VARCHAR(64) NOT NULL,
    version VARCHAR(32) NOT NULL,
    data TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 検査結果テーブル
CREATE TABLE inspections (
    id VARCHAR(64) PRIMARY KEY,
    lot_number VARCHAR(64) NOT NULL,
    machine_id VARCHAR(64) NOT NULL,
    operator_id VARCHAR(64) NOT NULL,
    result VARCHAR(16) NOT NULL CHECK (result IN ('pass', 'fail')),
    measured_values JSON NOT NULL,
    inspection_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ロット入出庫管理テーブル（トレーサビリティ用）
CREATE TABLE lot_inventory (
    id VARCHAR(64) PRIMARY KEY,
    lot_number VARCHAR(64) NOT NULL,
    product_type VARCHAR(64) NOT NULL, -- 原材料、半製品、完成品など
    quantity INT NOT NULL,
    in_out VARCHAR(8) NOT NULL CHECK (in_out IN ('in', 'out')),
    transaction_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```
