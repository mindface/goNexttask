# ER図（Markdown表現）

```
[users] ----------------< [user_roles] >---------------- [roles]
   |
   v
[organization_users] >------- [organizations]

[users] ---< [auth_tokens]

[organizations] ---< [production_orders]
[organizations] ---< [nc_programs]
[production_orders] ---< [inspection_results]
```

**テーブル定義（簡易）**
- `users`: id, name, email, password_hash, created_at, updated_at
- `roles`: id, name
- `user_roles`: user_id, role_id
- `organizations`: id, name, address, created_at
- `organization_users`: organization_id, user_id
- `auth_tokens`: id, user_id, token, expires_at
- `production_orders`: id, order_no, status, created_at
- `nc_programs`: id, name, version, machine_id
- `inspection_results`: id, lot_no, status, measured_at
