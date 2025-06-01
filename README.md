# タスク管理アプリ

Go言語で作成されたWebベースのタスク管理アプリケーションです。

## 機能

- ✅ **タスクの登録・管理**: ブラウザからタスクを簡単に登録
- ✅ **SQLiteデータベース**: 軽量で高速なローカルデータベース
- ✅ **完了ステータス管理**: タスクの完了・未完了状態を管理
- ✅ **期限設定**: カレンダーから期限を選択可能
- ✅ **詳細表示**: モーダルウィンドウでタスクの詳細を表示
- ✅ **レスポンシブデザイン**: スマートフォンにも対応

## 必要な環境

- Go 1.21以上
- SQLite3 (通常、Goのsqlite3ドライバに含まれます)

## インストール方法

### 1. リポジトリのクローン

```bash
git clone https://github.com/tossyyukky/claude-code-action-test1.git
cd claude-code-action-test1
```

### 2. 依存関係のインストール

```bash
go mod tidy
```

## 実行方法

### 開発環境での実行

```bash
go run main.go
```

### バイナリをビルドして実行

```bash
# ビルド
go build -o task-manager main.go

# 実行 (Windows)
./task-manager.exe

# 実行 (Linux/macOS)
./task-manager
```

## 使用方法

1. アプリケーションを起動します
2. ブラウザで http://localhost:8080 にアクセスします
3. 「新しいタスクを追加」ボタンをクリックしてタスクを作成します
4. タスク一覧で各タスクの状態を管理できます
5. 「詳細」ボタンでタスクの詳細情報を確認できます

## プロジェクト構成

```
.
├── main.go              # メインアプリケーション
├── go.mod              # Go依存関係管理
├── templates/          # HTMLテンプレート
│   ├── index.html      # メイン画面
│   └── create.html     # タスク作成画面
├── static/             # 静的ファイル
│   ├── style.css       # スタイルシート
│   └── script.js       # JavaScript
├── tasks.db            # SQLiteデータベース（実行時に自動生成）
└── README.md           # このファイル
```

## 技術仕様

- **バックエンド**: Go (Gorilla Mux, SQLite3)
- **フロントエンド**: HTML5, CSS3, JavaScript, Bootstrap 5
- **データベース**: SQLite3
- **ポート**: 8080

## データベーススキーマ

```sql
CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    deadline DATETIME,
    completed BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## トラブルシューティング

### ポート8080が使用中の場合

main.goの最後の行を編集してポート番号を変更してください：

```go
log.Fatal(http.ListenAndServe(":8081", r)) // 8080 -> 8081
```

### データベースファイルの場所

SQLiteデータベースファイル（`tasks.db`）は、アプリケーションの実行ディレクトリに作成されます。

## ライセンス

MIT License