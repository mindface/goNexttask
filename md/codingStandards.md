# Go Backend コーディング規約

## 1. 基本原則

### 1.1 Go言語の慣習に従う
- [Effective Go](https://golang.org/doc/effective_go.html)のガイドラインに準拠
- Go標準ライブラリのコードスタイルを参考にする
- `gofmt`と`golint`を必ず使用する

### 1.2 シンプルさと明確さ
- **KISS原則**: Keep It Simple, Stupid
- **DRY原則**: Don't Repeat Yourself
- **YAGNI原則**: You Aren't Gonna Need It

## 2. ファイル構造とパッケージ

### 2.1 ディレクトリ構造
```
goNexttask/
├── cmd/                  # エントリーポイント（main.go）
│   └── api/
│       └── main.go
├── internal/
│   ├── production/       # 生産管理 (サブドメイン)
│   │   ├── domain/
│   │   │   ├── entity.go
│   │   │   ├── value_object.go
│   │   │   ├── repository.go      # インターフェース
│   │   │   └── service.go         # ドメインサービス
│   │   ├── application/
│   │   │   ├── dto.go
│   │   │   └── usecase.go         # アプリケーションサービス
│   │   ├── infrastructure/
│   │   │   └── repository_impl.go # DBアクセス実装
│   │   └── interface/
│   │       └── rest_handler.go    # HTTPハンドラー
│   ├── nc/              # NC加工連携 (サブドメイン)
│   │   └── ...（同構成）
│   └── quality/         # 品質管理 (サブドメイン)
│       └── ...（同構成）
├── pkg/                  # 共通ユーティリティ
│   ├── middleware/
│   ├── config/
│   └── validator/
└── docs/
```

### 2.2 パッケージ設計
```go
// パッケージコメントは必須
// Package controller handles HTTP requests and responses
package controller

// インポートは標準ライブラリ、サードパーティ、プロジェクト内の順
import (
    "encoding/json"
    "net/http"
    
    "github.com/gin-gonic/gin"
    
    "betaTasker/model"
    "betaTasker/service"
)
```

## 3. 命名規則

### 3.1 変数・関数名
```go
// 良い例
userID      // キャメルケース
getUserByID // 動詞から始まる関数名
MaxRetries  // 定数は大文字始まり

// 悪い例
userid      // 単語の区切りが不明
GetUser     // IDの指定が不明確
MAX_RETRIES // スネークケースは使わない
```

### 3.2 構造体とインターフェース
```go
// 構造体は名詞、単数形
type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
}

// インターフェースは~erで終わる
type UserRepository interface {
    FindByID(id int) (*User, error)
    Save(user *User) error
}

// 小さなインターフェースを推奨
type Reader interface {
    Read([]byte) (int, error)
}
```

### 3.3 ファイル名
```
task_controller.go      // スネークケース
task_controller_test.go // テストファイル
task_repository.go      // 機能別に分割
```

## 4. エラーハンドリング

### 4.1 エラーの返却
```go
// エラーは最後の戻り値として返す
func GetUser(id int) (*User, error) {
    user, err := repository.FindByID(id)
    if err != nil {
        // エラーをラップして文脈を追加
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }
    return user, nil
}
```

### 4.2 カスタムエラー
```go
// エラー型の定義
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

// センチネルエラー
var (
    ErrUserNotFound = errors.New("user not found")
    ErrUnauthorized = errors.New("unauthorized")
)
```

### 4.3 エラーのチェック
```go
// errors.Isを使用
if errors.Is(err, ErrUserNotFound) {
    return c.JSON(404, gin.H{"error": "User not found"})
}

// 型アサーション
var valErr *ValidationError
if errors.As(err, &valErr) {
    return c.JSON(400, gin.H{"field": valErr.Field, "error": valErr.Message})
}
```

## 5. 関数設計

### 5.1 関数の長さと責任
```go
// 1つの関数は1つの責任のみ（単一責任の原則）
// 30行以内を目安にする
func (s *TaskService) CreateTask(task *model.Task) error {
    // バリデーション
    if err := s.validateTask(task); err != nil {
        return err
    }
    
    // ビジネスロジック
    task.Status = "pending"
    task.CreatedAt = time.Now()
    
    // 永続化
    return s.repository.Save(task)
}

// バリデーションは別関数に分離
func (s *TaskService) validateTask(task *model.Task) error {
    if task.Title == "" {
        return ValidationError{Field: "title", Message: "required"}
    }
    if task.Priority < 1 || task.Priority > 5 {
        return ValidationError{Field: "priority", Message: "must be between 1 and 5"}
    }
    return nil
}
```

### 5.2 引数と戻り値
```go
// 引数は3つ以内を推奨
// それ以上の場合は構造体を使用
type CreateTaskRequest struct {
    Title       string
    Description string
    Priority    int
    DueDate     *time.Time
}

func CreateTask(req CreateTaskRequest) (*Task, error) {
    // ...
}

// オプショナルなパラメータにはポインタを使用
func UpdateTask(id int, title *string, status *string) error {
    if title != nil {
        // titleを更新
    }
    if status != nil {
        // statusを更新
    }
    return nil
}
```

## 6. 並行処理

### 6.1 ゴルーチン
```go
// ゴルーチンリークを防ぐ
func ProcessTasks(tasks []Task) {
    var wg sync.WaitGroup
    // チャネルでゴルーチン数を制限
    sem := make(chan struct{}, 10)
    
    for _, task := range tasks {
        wg.Add(1)
        sem <- struct{}{}
        
        go func(t Task) {
            defer wg.Done()
            defer func() { <-sem }()
            
            processTask(t)
        }(task)
    }
    
    wg.Wait()
}
```

### 6.2 チャネル
```go
// チャネルは必ずクローズする
func Producer() <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 0; i < 10; i++ {
            ch <- i
        }
    }()
    return ch
}

// select文でタイムアウト処理
func ProcessWithTimeout(ch <-chan int) error {
    select {
    case val := <-ch:
        return process(val)
    case <-time.After(5 * time.Second):
        return errors.New("timeout")
    }
}
```

## 7. テスト

### 7.1 テストファイル構造
```go
package task_test

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    "betaTasker/service"
)

// テストケースは明確な名前を付ける
func TestTaskService_CreateTask_Success(t *testing.T) {
    // Arrange
    mockRepo := new(MockTaskRepository)
    svc := service.NewTaskService(mockRepo)
    task := &model.Task{Title: "Test Task"}
    
    mockRepo.On("Save", task).Return(nil)
    
    // Act
    err := svc.CreateTask(task)
    
    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### 7.2 テーブルドリブンテスト
```go
func TestValidateTask(t *testing.T) {
    tests := []struct {
        name    string
        task    *Task
        wantErr bool
    }{
        {
            name:    "valid task",
            task:    &Task{Title: "Test", Priority: 3},
            wantErr: false,
        },
        {
            name:    "empty title",
            task:    &Task{Title: "", Priority: 3},
            wantErr: true,
        },
        {
            name:    "invalid priority",
            task:    &Task{Title: "Test", Priority: 10},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateTask(tt.task)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## 8. データベースアクセス

### 8.1 Repository パターン
```go
// インターフェース定義
type TaskRepository interface {
    FindByID(id int) (*Task, error)
    FindAll(userID int) ([]*Task, error)
    Save(task *Task) error
    Update(task *Task) error
    Delete(id int) error
}

// GORM実装
type GormTaskRepository struct {
    db *gorm.DB
}

func (r *GormTaskRepository) FindByID(id int) (*Task, error) {
    var task Task
    if err := r.db.First(&task, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrTaskNotFound
        }
        return nil, fmt.Errorf("failed to find task: %w", err)
    }
    return &task, nil
}
```

### 8.2 トランザクション
```go
func (s *TaskService) TransferTask(fromUserID, toUserID, taskID int) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // トランザクション内での処理
        var task Task
        if err := tx.First(&task, taskID).Error; err != nil {
            return err
        }
        
        if task.UserID != fromUserID {
            return ErrUnauthorized
        }
        
        task.UserID = toUserID
        return tx.Save(&task).Error
    })
}
```

## 9. ログ

### 9.1 構造化ログ
```go
import "github.com/sirupsen/logrus"

// ログレベルの使い分け
func ProcessTask(taskID int) {
    log := logrus.WithFields(logrus.Fields{
        "task_id": taskID,
        "action":  "process",
    })
    
    log.Debug("Starting task processing")
    
    if err := doSomething(); err != nil {
        log.WithError(err).Error("Failed to process task")
        return
    }
    
    log.Info("Task processed successfully")
}
```

### 9.2 ログレベル
- **Debug**: 開発時のデバッグ情報
- **Info**: 通常の処理フロー
- **Warn**: 潜在的な問題
- **Error**: エラー発生（処理は継続）
- **Fatal**: 致命的エラー（プログラム終了）

## 10. セキュリティ

### 10.1 入力値検証
```go
// 必ず入力値を検証する
func ValidateEmail(email string) error {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    if !emailRegex.MatchString(email) {
        return errors.New("invalid email format")
    }
    return nil
}

// SQLインジェクション対策（GORMのプレースホルダを使用）
func GetUserByEmail(email string) (*User, error) {
    var user User
    // 良い例
    err := db.Where("email = ?", email).First(&user).Error
    
    // 悪い例（絶対に避ける）
    // err := db.Where(fmt.Sprintf("email = '%s'", email)).First(&user).Error
    
    return &user, err
}
```

### 10.2 秘密情報の管理
```go
// 環境変数から読み込む
type Config struct {
    DatabaseURL string
    JWTSecret   string
    APIKey      string
}

func LoadConfig() *Config {
    return &Config{
        DatabaseURL: os.Getenv("DATABASE_URL"),
        JWTSecret:   os.Getenv("JWT_SECRET"),
        APIKey:      os.Getenv("API_KEY"),
    }
}

// 秘密情報をログに出力しない
func (c *Config) String() string {
    return fmt.Sprintf("Config{DatabaseURL: ***, JWTSecret: ***, APIKey: ***}")
}
```

## 11. パフォーマンス

### 11.1 データベースクエリ最適化
```go
// N+1問題を避ける（Preloadを使用）
func GetTasksWithUsers() ([]*Task, error) {
    var tasks []*Task
    // 良い例
    err := db.Preload("User").Find(&tasks).Error
    
    // 悪い例
    // for _, task := range tasks {
    //     db.First(&task.User, task.UserID)
    // }
    
    return tasks, err
}

// 必要なカラムのみ取得
func GetTaskTitles() ([]string, error) {
    var titles []string
    err := db.Model(&Task{}).Pluck("title", &titles).Error
    return titles, err
}
```

### 11.2 メモリ管理
```go
// スライスの容量を事前に確保
func ProcessLargeData(count int) []Result {
    results := make([]Result, 0, count)
    for i := 0; i < count; i++ {
        results = append(results, processItem(i))
    }
    return results
}

// defer文は関数の最後で実行される
func ReadFile(path string) ([]byte, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close() // 必ずクローズする
    
    return ioutil.ReadAll(file)
}
```

## 12. コードレビューチェックリスト

### 必須確認項目
- [ ] `gofmt`でフォーマット済み
- [ ] `golint`でリント済み
- [ ] `go vet`でチェック済み
- [ ] テストが書かれている
- [ ] テストがパスする
- [ ] エラーハンドリングが適切
- [ ] ログが適切に出力される
- [ ] ドキュメントコメントがある
- [ ] セキュリティ考慮されている
- [ ] パフォーマンスが考慮されている

### 推奨確認項目
- [ ] 関数が30行以内
- [ ] 循環的複雑度が10以下
- [ ] DRY原則に従っている
- [ ] 命名が明確で一貫性がある
- [ ] インターフェースが小さい
- [ ] 並行処理が安全

## 13. コミットメッセージ

### フォーマット
```
<type>(<scope>): <subject>

<body>

<footer>
```

### タイプ
- **feat**: 新機能
- **fix**: バグ修正
- **docs**: ドキュメント変更
- **style**: コードフォーマット
- **refactor**: リファクタリング
- **test**: テスト追加・修正
- **chore**: ビルドプロセスやツールの変更

### 例
```
feat(task): add task priority feature

- Add priority field to task model
- Update task service to handle priority
- Add validation for priority range (1-5)

Closes #123
```

---

最終更新日: 2024年
バージョン: 1.0.0