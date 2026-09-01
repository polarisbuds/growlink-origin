# Usecase テスト実装手順書

この文書は、`internal/application/usecase` 配下の Usecase に対する **単体テスト** の書き方をまとめたものです。  
新しい Usecase メソッドを追加したとき、または既存メソッドのテストを足すときに、この手順に沿って実装してください。

Repository の統合テスト（本物の PostgreSQL）とは役割が違います。DB やファイル実体は使わず、依存は **gomock で生成したモック** で差し替えます。

Repository 側の手順は [Repository テスト実装手順書](../repository/README.md) を参照してください。

---

## 1. このテストで検証すること

Usecase テストは **アプリケーションの分岐と呼び出し順** を確認します。

| 確認すること | 例 |
| --- | --- |
| 通常系 | 依存が成功したとき、期待するリポジトリメソッドが呼ばれ、エラーが返らない |
| 異常系 | 依存がエラーを返したとき、ラップしたメッセージで失敗し、後続処理が走らない |
| 不正値 | `userID <= 0`、空文字、`nil` の Reader など、入力が不正なときの契約 |

**やらないこと**

- SQL・GORM・PostgreSQL の制約は、ここでは扱わない（Repository テスト側）
- HTTP / Echo / テンプレートの描画は、ここでは扱わない（Adapter テスト側）
- モックの中で「本物の Repository」を呼び出さない

---

## 2. 前提

- Docker は不要（DB を起動しない）
- 依存のインターフェースは `internal/domain/interface/repository` などにある
- モックは `mockgen` で生成し、各インターフェース配下の `mock/` に置く
- テストは実装パッケージの **外部テスト**（`xxx_usecase_test`）で書く
- ライブラリは **`go.uber.org/mock`**（`mockgen` + `gomock`）

ローカルでの実行例:

```bash
# モック再生成
go generate ./internal/domain/interface/... ./internal/adapter/interface/...

# 対象パッケージ全体
go test -v -count=1 ./internal/application/usecase/setting_usecase/

# メソッド単位
go test -v -count=1 ./internal/application/usecase/setting_usecase/ -run "TestRegisterIconUsecase_UploadIcon"
```

---

## 3. 全体の流れ（毎回この順）

```text
1. 対象メソッドの契約を読む（成功時 / 依存エラー時 / 不正値）
2. インターフェースに go:generate があることを確認し、必要なら go generate する
3. テストファイルを <パッケージ名>_test で用意する
4. テストファイル内に newXxxUsecase を置き、モックを注入する
5. テーブル駆動で「通常系・異常系・不正値」を並べ、setup で EXPECT を積む
6. エラーだけでなく、DoAndReturn でキャプチャした引数も断言する
7. 実行して失敗したら、テストの期待を実装のバグに合わせない
```

---

## 4. ファイルとパッケージの約束

| 項目 | ルール |
| --- | --- |
| テストファイル | 実装と同じディレクトリの `<対象>_usecase_test.go` |
| 共通ヘルパー | 同パッケージの `helper_test.go`（テストデータ・アサーションのみ） |
| パッケージ名 | **`<usecaseパッケージ>_test`**（例: `setting_usecase_test`） |
| テスト関数名 | `Test<型名>_<メソッド名>` |
| ケース名 | `通常系_...` / `異常系_...` / `不正値_...` |
| モック注入 | **テストファイルごと** に `newXxxUsecase(t)` を置く |
| 依存の差し替え | `internal/domain/interface/**/mock` の生成モックを使う |

外部テストにする理由は、公開コンストラクタ（`NewXxxUsecase`）と公開メソッドだけを叩くためです。

---

## 5. mockgen の設定

### 5.1 go:generate コメント

各インターフェースファイルの package 宣言直下に付けます。

```go
package repository_interface

//go:generate mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import "growlink/internal/domain/entity"

type UserRepository interface {
    GetByID(id int64) (*entity.User, error)
    UpdateByID(user *entity.User) error
}
```

| オプション | 意味 |
| --- | --- |
| `-typed` | 型付き EXPECT。引数型の誤りをコンパイルで検知できる |
| `-source=$GOFILE` | 同ファイル内のインターフェースを対象にする |
| `-destination=mock/mock_$GOFILE` | 同階層の `mock/` サブパッケージへ出力 |
| `-package=mock` | 生成コードのパッケージ名 |

生成例:

- `internal/domain/interface/repository/user_repository.go` → `internal/domain/interface/repository/mock/mock_user_repository.go`
- `internal/adapter/interface/usecase/setting_usecase.go` → `internal/adapter/interface/usecase/mock/mock_setting_usecase.go`

### 5.2 生成コマンド

```bash
go generate ./internal/domain/interface/... ./internal/adapter/interface/...
```

`tools.go`（`//go:build tools`）で `go.uber.org/mock/mockgen` をツール依存としてピン留めしています。

### 5.3 インターフェースを変更したとき

1. インターフェースにメソッドを追加 / 変更する
2. `go generate` を再実行する
3. 生成された `mock/` をコミットする
4. Usecase テストの EXPECT を更新する

---

## 6. モック注入（ファイル内共通関数）

**1 つのテストファイルに複数テストがある場合**、そのファイル内に `newXxxUsecase` を置き、Controller とモック生成を共通化します。  
ファイルをまたぐ注入ヘルパーは作りません。

```go
func newRegisterIconUsecase(t *testing.T) (
    *setting_usecase.RegisterIconUsecase,
    *repomock.MockUserRepository,
    *repomock.MockUploadFileRepository,
) {
    t.Helper()
    ctrl := gomock.NewController(t)
    userRepo := repomock.NewMockUserRepository(ctrl)
    fileRepo := repomock.NewMockUploadFileRepository(ctrl)
    return setting_usecase.NewRegisterIconUsecase(fileRepo, userRepo), userRepo, fileRepo
}
```

ポイント:

- `gomock.NewController(t)` は Go 1.22 以降、`t.Cleanup` で自動 Finish する
- **各 `t.Run` ごと** に `newRegisterIconUsecase(t)` を呼ぶ（Controller をケース間で共有しない）
- `helper_test.go` には `newTestUser` や `assertUsecaseErr` など、モック非依存の共通部品だけ置く

---

## 7. テストの骨格（テーブル駆動 + setup）

1 メソッドに対して、入力と setup を変えたケースをスライスに並べ、`t.Run` で回します。

```go
type registerIconCaptures struct {
    newFileNames []string
    updatedUser  *entity.User
}

func TestRegisterIconUsecase_UploadIcon(t *testing.T) {
    tests := []struct {
        name        string
        userID      int64
        filename    string
        content     io.Reader
        wantErr     bool
        errContains string
        setup       func(t *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *registerIconCaptures)
        assertFn    func(t *testing.T, cap *registerIconCaptures)
    }{
        // 通常系 / 異常系 / 不正値
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            uc, userRepo, fileRepo := newRegisterIconUsecase(t)
            cap := &registerIconCaptures{}
            if tt.setup != nil {
                tt.setup(t, userRepo, fileRepo, cap)
            }

            gotErr := uc.UploadIcon(tt.userID, tt.filename, tt.content)
            assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "UploadIcon")
            if t.Failed() {
                return
            }
            if tt.assertFn != nil {
                tt.assertFn(t, cap)
            }
        })
    }
}
```

### 7.1 フィールドの意味

| フィールド | 役割 |
| --- | --- |
| `name` | サブテスト名。`通常系_` / `異常系_` / `不正値_` で始める |
| `setup` | ケースごとに `EXPECT()` を積む。`t.Run` 内で呼ぶ |
| 入力 | Usecase メソッドの引数 |
| `wantErr` | エラーを期待するか |
| `errContains` | `fmt.Errorf("...: %w", err)` のメッセージ断片 |
| `assertFn` | 成功 / 失敗のあと、キャプチャした引数を見る（省略可） |

エラー判定だけだと「成功したが間違ったファイルを消した」を見逃します。通常系では `assertFn` で副作用を見てください。

---

## 8. ケースの切り方

メソッドの分岐を上からなぞり、次の 3 種を検討します。契約に無いケースは書きません。

### 8.1 通常系

依存をすべて成功させ、Usecase が組み立てた値を断言します。

```go
setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *registerIconCaptures) {
    user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, nil), nil)
    file.EXPECT().NewFile(gomock.Any(), gomock.Any()).DoAndReturn(func(name string, _ io.Reader) error {
        cap.newFileNames = append(cap.newFileNames, name)
        return nil
    })
    user.EXPECT().UpdateByID(gomock.Any()).DoAndReturn(func(u *entity.User) error {
        cap.updatedUser = u
        return nil
    })
    file.EXPECT().Exists(gomock.Any()).Return(true, nil)
    file.EXPECT().DeleteFile(gomock.Any()).Return(nil)
},
assertFn: func(t *testing.T, cap *registerIconCaptures) {
    t.Helper()
    require.Len(t, cap.newFileNames, 1)
    assertUUIDFilename(t, cap.newFileNames[0], ".png")
    require.NotNil(t, cap.updatedUser.ImageURI)
},
```

UUID など非決定な値は `gomock.Any()` + `DoAndReturn` でキャプチャし、形式（UUID + 拡張子）で見ます。

### 8.2 異常系（依存のエラー）

Usecase が `if err != nil { return fmt.Errorf("固定メッセージ: %w", err) }` している箇所ごとに、1 ケース用意します。

```go
setup: func(_ *testing.T, user *repomock.MockUserRepository, _ *repomock.MockUploadFileRepository, _ *registerIconCaptures) {
    user.EXPECT().GetByID(validUserID).Return(nil, errUserNotFound)
},
wantErr:     true,
errContains: "ユーザアイコンのアップロードで存在しないユーザIDが指定されました",
```

早い段階の失敗では、後続メソッドの **EXPECT を置かない**。予期しない呼び出しがあれば gomock が失敗します。

`errContains` はラップメッセージの日本語部分を使い、内部の `errors.New("...")` だけに依存しないようにします。

### 8.3 不正値

Usecase または最初の依存が拒否する入力です。

| 種類 | 扱い |
| --- | --- |
| `userID` が `0` / 負数 | `GetByID` がエラーを返す EXPECT。後続は EXPECT しない |
| 空のファイル名、拡張子なし | 実装がどう正規化するかを断言する |
| `io.Reader` が `nil` | `NewFile` の `DoAndReturn` でエラーを返し、ラップメッセージを確認 |

「不正だから必ずエラー」と決め打ちせず、**今の契約どおり** に書きます。

---

## 9. テストデータ（エンティティ）

Usecase テストでは DB fixture（`testutils/fixtures`）は使いません。Docker もトランザクションも不要です。

メモリ上のエンティティで足りる場合は、`helper_test.go` の小さなヘルパーで十分です。

```go
func newTestUser(id int64, imageURI *string) *entity.User {
    return &entity.User{
        Model:    gorm.Model{ID: uint(id)},
        Name:     "icon-user",
        Email:    "icon@example.com",
        ImageURI: imageURI,
    }
}
```

- ID や名前はケースが読みやすい固定値でよい
- FK の整合や UNIQUE はここでは問わない
- 同じ形が複数 Usecase で必要になったら `helper_test.go` に共通化する

---

## 10. アサーションの使い分け

| 手段 | 用途 |
| --- | --- |
| `wantErr` + `assert.ErrorContains` | ラップされたエラーメッセージ |
| `DoAndReturn` のキャプチャ | 渡したファイル名、更新エンティティ |
| EXPECT を置かない | 到達してはいけない分岐（予期しない呼び出しで失敗） |
| `gomock.Any()` | 非決定値（UUID ファイル名など） |
| `require` | 後続の比較の前提（nil、件数） |
| `assert` | 値の比較（複数項目を一度に見たいとき） |

Repository テスト用の `assertions.NotFound` などは、GORM エラー向けです。Usecase のラップエラーには `ErrorContains` を使います。

---

## 11. 実装とテストを往復するときの注意

### 11.1 期待値をバグに合わせない

テストが落ちたとき、まずメソッドのコメントと呼び出し順を読みます。  
契約と実装が食い違っているなら **実装を直す** のが正しい対応です。

### 11.2 呼ばれてはいけない処理

早期 return のあとに `NewFile` や `UpdateByID` が走ると、副作用が残ります。  
異常系では後続メソッドの EXPECT を置かず、`assertFn` でキャプチャが空であることも見てください。

### 11.3 Controller はケースごと

`gomock.NewController(t)` は `t.Run` の中で呼びます。  
親テストの Controller を全ケースで共有すると、EXPECT の混線や Finish タイミングの問題が起きます。

---

## 12. 新規 Usecase テストを足すときのチェックリスト

1. [ ] パッケージは `<usecaseパッケージ>_test`
2. [ ] 依存インターフェースに `go:generate mockgen` がある
3. [ ] `go generate` 後の `mock/` がコミットされている
4. [ ] テストファイル内に `newXxxUsecase(t)` がある（ファイルをまたがない）
5. [ ] `helper_test.go` にはモック非依存の共通部品だけ置いている
6. [ ] テーブル駆動で `通常系` / `異常系` / `不正値` を分けている
7. [ ] 各ケースの `setup` で EXPECT を積んでいる
8. [ ] 異常系はラップメッセージ（`errContains`）を見ている
9. [ ] 異常系は後続メソッドの EXPECT を置いていない
10. [ ] 通常系は `DoAndReturn` でキャプチャした引数を見ている
11. [ ] `go test -count=1 -run "TestXxxUsecase_"` が通る（Docker 不要）

---

## 13. ディレクトリ対応表

| 役割 | 場所 |
| --- | --- |
| Usecase 実装 | `internal/application/usecase/<領域>/` |
| Usecase インターフェース | `internal/adapter/interface/usecase/` |
| Repository インターフェース | `internal/domain/interface/repository/` |
| 生成モック | `internal/domain/interface/**/mock/` |
| Usecase テスト | `internal/application/usecase/<領域>/*_usecase_test.go` |
| テスト共通ヘルパー | `internal/application/usecase/<領域>/helper_test.go` |
| mockgen ツール依存 | `tools.go` |
| DB fixture（使わない） | `internal/shared/testutils/fixtures` |
| Repository 統合テスト | `internal/infrastructure/repository/` |

```text
internal/
  application/usecase/setting_usecase/
    register_icon_usecase.go
    register_icon_usecase_test.go      ← newRegisterIconUsecase + setup/EXPECT
    remove_icon_usecase_test.go        ← newRemoveIconUsecase + setup/EXPECT
    helper_test.go                     ← newTestUser, assertUsecaseErr など
  domain/interface/repository/
    user_repository.go                 ← go:generate
    mock/mock_user_repository.go       ← 生成物（コミット対象）
  adapter/interface/usecase/
    setting_usecase.go                 ← go:generate
    mock/mock_setting_usecase.go
tools.go                               ← mockgen ピン留め
```

---

## 引用したデータ

本手順書は、次の実装・テスト・設定を根拠に作成しています。

| 種別 | パス | 参照した内容 |
| --- | --- | --- |
| 文書の体裁 | `doc/documents/test/repository/README.md` | 章立て、チェックリスト、最下部の引用表 |
| Usecase 実装 | `internal/application/usecase/setting_usecase/register_icon_usecase.go` | `GetByID` → `NewFile` → `UpdateByID` → `Exists` / `DeleteFile`、エラーの `fmt.Errorf` ラップ |
| Usecase 実装 | `internal/application/usecase/setting_usecase/remove_icon_usecase.go` | `GetByID` → `Exists` / `DeleteFile` → `UpdateByID` |
| Usecase インターフェース | `internal/adapter/interface/usecase/setting_usecase.go` | `UploadIcon` / `RemoveIcon`、`go:generate` |
| テスト本体 | `internal/application/usecase/setting_usecase/register_icon_usecase_test.go` | `newRegisterIconUsecase`、`setup` + EXPECT、`registerIconCaptures` |
| テスト本体 | `internal/application/usecase/setting_usecase/remove_icon_usecase_test.go` | `newRemoveIconUsecase`、`setup` + EXPECT、`removeIconCaptures` |
| 共通ヘルパー | `internal/application/usecase/setting_usecase/helper_test.go` | `newTestUser`、`assertUsecaseErr`、`assertUUIDFilename` |
| 生成モック | `internal/domain/interface/repository/mock/mock_user_repository.go` | `MockUserRepository`、`NewMockUserRepository` |
| 生成モック | `internal/domain/interface/repository/mock/mock_upload_file_repository.go` | `MockUploadFileRepository`、typed EXPECT |
| インターフェース | `internal/domain/interface/repository/user_repository.go` | `go:generate mockgen` コメント |
| ツール依存 | `tools.go` | `go.uber.org/mock/mockgen` のピン留め |
| Repository テストとの分担 | `doc/documents/test/repository/README.md` | Usecase の分岐は Repository テストでは扱わない、という境界 |
