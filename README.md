# Todo App

フロントエンド（Next.js）とバックエンド（Go）を使用したタスク管理アプリケーションです。

## プロジェクト構成

```
todo-app-practice/
├── frontend/          # Next.js フロントエンド
│   ├── src/
│   │   ├── app/       # App Router
│   │   ├── components/ # React コンポーネント
│   │   ├── hooks/     # カスタムフック
│   │   ├── lib/       # ユーティリティ
│   │   ├── types/     # TypeScript 型定義
│   │   └── utils/     # ヘルパー関数
│   └── package.json
├── backend/           # Go バックエンド
│   ├── cmd/
│   │   └── main.go    # アプリケーションエントリーポイント
│   ├── internal/
│   │   ├── config/    # 設定管理
│   │   ├── handlers/  # HTTP ハンドラー
│   │   ├── middleware/ # ミドルウェア
│   │   ├── models/    # データモデル
│   │   ├── repository/ # データアクセス層
│   │   └── services/  # ビジネスロジック
│   └── go.mod
└── README.md
```

## 機能

### 認証機能
- ユーザー登録
- ログイン/ログアウト
- JWT トークン認証
- リフレッシュトークン

### タスク管理機能
- タスクの作成
- タスクの一覧表示
- タスクの更新
- タスクの削除
- 優先度設定（high, medium, low）
- ステータス管理（pending, completed）
- 期限設定

## 技術スタック

### フロントエンド
- **Next.js 14** - React フレームワーク
- **TypeScript** - 型安全性
- **Tailwind CSS** - スタイリング
- **React Query** - データフェッチング

### バックエンド
- **Go 1.25** - プログラミング言語
- **Echo** - Web フレームワーク
- **SQLite** - データベース
- **JWT** - 認証トークン
- **bcrypt** - パスワードハッシュ化

## セットアップ

### 前提条件
- Node.js 18以上
- Go 1.25以上
- Docker & Docker Compose（推奨）

### 開発環境でのセットアップ

#### 方法1: Docker Compose（推奨）
```bash
# 開発環境を起動
make dev

# または直接実行
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

#### 方法2: 手動セットアップ

**バックエンドのセットアップ**
```bash
cd backend
go mod tidy
go run cmd/main.go
```

**フロントエンドのセットアップ**
```bash
cd frontend
npm install
npm run dev
```

### 本番環境でのセットアップ

#### Docker Composeを使用
```bash
# 環境変数を設定
export JWT_SECRET="your-production-secret-key"

# 本番環境を起動
make deploy

# または直接実行
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

#### 手動デプロイ
```bash
# バックエンドをビルド
make build

# 環境変数を設定
export JWT_SECRET="your-production-secret-key"
export ENVIRONMENT="production"
export PORT="8080"

# サーバーを起動
./backend/main
```

## API エンドポイント

### 認証
- `POST /api/auth/register` - ユーザー登録
- `POST /api/auth/login` - ログイン
- `POST /api/auth/refresh` - トークンリフレッシュ
- `POST /api/auth/logout` - ログアウト

### タスク管理
- `GET /api/tasks` - タスク一覧取得
- `POST /api/tasks` - タスク作成
- `PUT /api/tasks/:id` - タスク更新
- `DELETE /api/tasks/:id` - タスク削除

### その他
- `GET /health` - ヘルスチェック

## データベーススキーマ

### users テーブル
- `id` (TEXT, PRIMARY KEY)
- `email` (TEXT, UNIQUE, NOT NULL)
- `password_hash` (TEXT, NOT NULL)
- `name` (TEXT, NOT NULL)
- `created_at` (DATETIME, NOT NULL)
- `updated_at` (DATETIME, NOT NULL)

### tasks テーブル
- `id` (TEXT, PRIMARY KEY)
- `user_id` (TEXT, NOT NULL, FOREIGN KEY)
- `title` (TEXT, NOT NULL)
- `description` (TEXT)
- `deadline` (DATETIME)
- `priority` (TEXT, NOT NULL)
- `status` (TEXT, NOT NULL)
- `created_at` (DATETIME, NOT NULL)
- `updated_at` (DATETIME, NOT NULL)

## 開発

### 利用可能なコマンド（Makefile）

```bash
# ヘルプを表示
make help

# バックエンドをビルド
make build

# バックエンドを実行
make run

# テストを実行
make test

# ビルド成果物をクリーンアップ
make clean

# Dockerイメージをビルド
make docker-build

# Docker Composeで実行
make docker-run

# Docker Composeを停止
make docker-stop

# 開発環境を起動
make dev

# 本番環境にデプロイ
make deploy
```

### 手動開発

**バックエンドの開発**
```bash
cd backend
go run cmd/main.go
```

**フロントエンドの開発**
```bash
cd frontend
npm run dev
```

### テスト
```bash
# バックエンドのテスト
make test

# または手動実行
cd backend && go test ./...

# フロントエンドのテスト
cd frontend && npm test
```

## デプロイ

### Docker Composeを使用（推奨）

**開発環境**
```bash
make dev
```

**本番環境**
```bash
# 環境変数を設定
export JWT_SECRET="your-production-secret-key"

# デプロイ
make deploy
```

### 手動デプロイ

**バックエンド**
```bash
make build
export JWT_SECRET="your-production-secret-key"
export ENVIRONMENT="production"
./backend/main
```

**フロントエンド**
```bash
cd frontend
npm run build
npm start
```

## ライセンス

MIT License
