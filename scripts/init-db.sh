#!/bin/bash

set -e

echo "================================================"
echo "  データベース初期化スクリプト"
echo "================================================"

# カラー定義
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# PostgreSQLコンテナが起動しているか確認
if ! docker-compose ps | grep -q "postgres.*Up"; then
    echo -e "${YELLOW}PostgreSQLコンテナが起動していません。起動します...${NC}"
    docker-compose up -d postgres
    echo "PostgreSQLの起動を待機中..."
    sleep 10
fi

echo ""
echo "データベースを初期化しています..."

# マイグレーション実行
docker-compose exec -T postgres psql -U postgres -d gonexttask << EOF
-- Drop existing tables if needed (be careful in production!)
DROP TABLE IF EXISTS measurement_results CASCADE;
DROP TABLE IF EXISTS inspections CASCADE;
DROP TABLE IF EXISTS machines CASCADE;
DROP TABLE IF EXISTS nc_programs CASCADE;
DROP TABLE IF EXISTS production_orders CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Run migration
\i /docker-entrypoint-initdb.d/001_create_tables.sql

-- Insert sample data
INSERT INTO users (id, email, password_hash, name, role, created_at, updated_at)
VALUES 
    ('user-001', 'admin@example.com', '\$2a\$10\$YourHashedPasswordHere', 'Admin User', 'admin', NOW(), NOW()),
    ('user-002', 'operator@example.com', '\$2a\$10\$YourHashedPasswordHere', 'Operator User', 'operator', NOW(), NOW());

INSERT INTO machines (id, name, ip_address, machine_type, capabilities, running_state, last_heartbeat, created_at, updated_at)
VALUES
    ('machine-001', 'CNC Machine 1', '192.168.1.101', 'CNC-3AXIS', '["milling", "drilling"]', 'stopped', NOW(), NOW(), NOW()),
    ('machine-002', 'CNC Machine 2', '192.168.1.102', 'CNC-5AXIS', '["milling", "drilling", "turning"]', 'stopped', NOW(), NOW(), NOW()),
    ('machine-003', 'Lathe 1', '192.168.1.103', 'LATHE', '["turning"]', 'stopped', NOW(), NOW(), NOW());

-- Show created tables
\dt

-- Count records
SELECT 'Users:' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'Machines:', COUNT(*) FROM machines;
EOF

echo ""
echo -e "${GREEN}✅ データベースの初期化が完了しました${NC}"
echo ""
echo "初期ユーザー:"
echo "  Admin:    admin@example.com"
echo "  Operator: operator@example.com"
echo ""
echo "※ パスワードはアプリケーションで登録し直してください"