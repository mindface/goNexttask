# 認証フロー図（Markdown表現）

```
Client(Web/Mobile)
        |
        | POST /auth/login (email, password)
        v
API Server
        |
        | Verify credentials
        v
Database
        |
        | User record
        v
API Server
        |
        | Generate JWT & store (optional blacklist)
        v
Client(Web/Mobile)
        |
        | Subsequent requests (Authorization: Bearer token)
        v
API Server
        |
        | Validate token -> Fetch user/org data
        v
Database
```
