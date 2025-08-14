  # 1. 完全リセットして起動
  make dev-down
  docker volume prune -f
  make dev-up

  # 2. サーバーが起動するまで待機（30秒程度）
  sleep 30

  # 3. ヘルスチェック
  curl http://localhost:8080/health

  # 4. ユーザー登録テスト
  curl -X POST http://localhost:8080/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d '{
      "email": "admin@test.com",
      "password": "password123",
      "name": "Admin User",
      "role": "admin"
    }'