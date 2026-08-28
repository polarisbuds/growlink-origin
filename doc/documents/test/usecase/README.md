# Usecase テスト実装手順書

この文書は、`internal/application/usecase` 配下の Usecase に対する **単体テスト** の書き方をまとめたものです。  
新しい Usecase メソッドを追加したとき、または既存メソッドのテストを足すときに、この手順に沿って実装してください。

Repository の統合テスト（本物の PostgreSQL）とは役割が違います。DB やファイル実体は使わず、依存は **stab（スタブ）** で差し替えます。

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
- stab の中で「本物の Repository」を呼び出さない

---

## 2. 前提

- Docker は不要（DB を起動しない）
- 依存のインターフェースは `internal/domain/interface/repository` にある
- 差し替え用の stab は `internal/shared/testutils/stab` にある
- テストは実装パッケージの **外部テスト**（`xxx_usecase_test`）で書く

ローカルでの実行例:

```bash
# 対象パッケージ全体
go test -v -count=1 ./internal/application/usecase/setting_usecase/

# メソッド単位
go test -v -count=1 ./internal/application/usecase/setting_usecase/ -run "TestRegisterIconUsecase_UploadIcon"
```

---

## 3. 全体の流れ（毎回この順）

```text
1. 対象メソッドの契約を読む（成功時 / 依存エラー時 / 不正値）
2. 必要な Repository stab が無ければ追加する
3. テストファイルを <パッケージ名>_test で用意する
4. テーブル駆動で「通常系・異常系・不正値」を並べる
5. エラーだけでなく、呼ばれた / 呼ばれていないメソッドも断言する
6. 実行して失敗したら、テストの期待を実装のバグに合わせない
```

---

## 4. ファイルとパッケージの約束

| 項目 | ルール |
| --- | --- |
| テストファイル | 実装と同じディレクトリの `<対象>_usecase_test.go` |
| パッケージ名 | **`<usecaseパッケージ>_test`**（例: `setting_usecase_test`） |
| テスト関数名 | `Test<型名>_<メソッド名>` |
| ケース名 | `通常系_...` / `異常系_...` / `不正値_...` |
| 依存の差し替え | `internal/shared/testutils/stab`。テストファイル内に stab を再定義しない |

外部テストにする理由は、公開コンストラクタ（`NewXxxUsecase`）と公開メソッドだけを叩くためです。

---

## 5. なぜ stab を使うか

Usecase は Repository インターフェースに依存しています。テストではその実装を差し替え、次を自由に制御します。

- 成功を返す
- 決まったエラーを返す
- 引数（保存されたファイル名、更新されたエンティティ）を記録する
- 「呼ばれてはいけないメソッド」を `t.Fatal` で検知する

gomock / mockery は使いません。プロジェクト共通の stab 構造体に、関数フィールドを渡す方式です。

---

## 6. stab の置き方

### 6.1 ファイル名と型名

| 項目 | ルール |
| --- | --- |
| ディレクトリ | `internal/shared/testutils/stab` |
| ファイル名 | ドメイン側と同じベース名 + `_stab.go`（例: `user_repository.go` → `user_repository_stab.go`） |
| パッケージ名 | `stab` |
| 型名 | インターフェースと同じ（`stab.UserRepository`） |

テストからは次のように使います。

```go
userRepo := &stab.UserRepository{
    GetByIDFunc: func(id int64) (*entity.User, error) {
        return user, nil
    },
}
fileRepo := &stab.UploadFileRepository{
    NewFileFunc: func(string, io.Reader) error { return nil },
}
uc := setting_usecase.NewRegisterIconUsecase(fileRepo, userRepo)
```

### 6.2 stab の骨格

インターフェースを満たすこと、呼び出しを記録すること、未設定の関数は安全に失敗することをセットにします。

```go
package stab

var _ repository_interface.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
    GetByIDFunc    func(id int64) (*entity.User, error)
    UpdateByIDFunc func(user *entity.User) error

    UpdatedUser *entity.User // 呼び出し結果の記録
}

func (s *UserRepository) GetByID(id int64) (*entity.User, error) {
    if s.GetByIDFunc == nil {
        return nil, errors.New("unexpected GetByID")
    }
    return s.GetByIDFunc(id)
}
```

ポイントは次のとおりです。

- 差し替え用は `XxxFunc` という公開フィールドにする（テストから代入するため）
- 引数や更新後エンティティはスライス / ポインタで記録する（後から断言するため）
- `var _ インターフェース = (*型)(nil)` でコンパイル時に実装漏れを検知する
- **同じインターフェースの stab をテストファイルにコピーしない。** 足りなければ `stab` にメソッドを足す

### 6.3 既存 stab を拡張するとき

インターフェースにメソッドが増えたら、対応する `*_stab.go` に `XxxFunc` と記録用フィールドを足します。  
Usecase テスト側にローカル型を作ると、他テストと二重管理になります。

---

## 7. テストの骨格（テーブル駆動）

1 メソッドに対して、入力と依存の戻りを変えたケースをスライスに並べ、`t.Run` で回します。

```go
func TestRegisterIconUsecase_UploadIcon(t *testing.T) {
    tests := []struct {
        name        string
        fileRepo    repository_interface.UploadFileRepository
        userRepo    repository_interface.UserRepository
        userID      int64
        filename    string
        content     io.Reader
        wantErr     bool
        errContains string
        assertFn    func(t *testing.T, fileRepo *stab.UploadFileRepository, userRepo *stab.UserRepository)
    }{
        // 通常系 / 異常系 / 不正値
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            uc := setting_usecase.NewRegisterIconUsecase(tt.fileRepo, tt.userRepo)
            gotErr := uc.UploadIcon(tt.userID, tt.filename, tt.content)

            if gotErr != nil {
                if !tt.wantErr {
                    t.Errorf("UploadIcon() failed: %v", gotErr)
                    return
                }
                if tt.errContains != "" {
                    assert.ErrorContains(t, gotErr, tt.errContains)
                }
            } else if tt.wantErr {
                t.Fatal("UploadIcon() succeeded unexpectedly")
            }

            if tt.assertFn != nil {
                fileRepo, okFile := tt.fileRepo.(*stab.UploadFileRepository)
                userRepo, okUser := tt.userRepo.(*stab.UserRepository)
                require.True(t, okFile && okUser)
                tt.assertFn(t, fileRepo, userRepo)
            }
        })
    }
}
```

### 7.1 フィールドの意味

| フィールド | 役割 |
| --- | --- |
| `name` | サブテスト名。`通常系_` / `異常系_` / `不正値_` で始める |
| 依存 stab | ケースごとに組み立てる。共有ポインタをケース間で使い回さない |
| 入力 | Usecase メソッドの引数 |
| `wantErr` | エラーを期待するか |
| `errContains` | `fmt.Errorf("...: %w", err)` のメッセージ断片 |
| `assertFn` | 成功 / 失敗のあと、呼び出し記録を見る（省略可） |

エラー判定だけだと「成功したが間違ったファイルを消した」を見逃します。通常系では `assertFn` で副作用を見てください。

### 7.2 ケースごとに stab を new する

テーブルの各要素で `&stab.UserRepository{...}` を新たに作ります。  
1 つの stab を全ケースで共有すると、記録用スライスが混ざり、失敗理由が分かりにくくなります。

---

## 8. ケースの切り方

メソッドの分岐を上からなぞり、次の 3 種を検討します。契約に無いケースは書きません。

### 8.1 通常系

依存をすべて成功させ、Usecase が組み立てた値を断言します。

```go
{
    name: "通常系_拡張子ありのアイコン登録",
    fileRepo: &stab.UploadFileRepository{
        ExistsFunc:     func(string) (bool, error) { return true, nil },
        NewFileFunc:    func(string, io.Reader) error { return nil },
        DeleteFileFunc: func(string) error { return nil },
    },
    userRepo: &stab.UserRepository{
        GetByIDFunc: func(id int64) (*entity.User, error) {
            return newTestUser(id, nil), nil
        },
    },
    userID:   10,
    filename: "avatar.png",
    content:  strings.NewReader("png-bytes"),
    wantErr:  false,
    assertFn: func(t *testing.T, fileRepo *stab.UploadFileRepository, userRepo *stab.UserRepository) {
        t.Helper()
        require.Len(t, fileRepo.NewFileNames, 1)
        gotName := fileRepo.NewFileNames[0]
        require.True(t, strings.HasSuffix(gotName, ".png"))
        require.NotNil(t, userRepo.UpdatedUser.ImageURI)
        assert.Equal(t, gotName, *userRepo.UpdatedUser.ImageURI)
    },
}
```

UUID など非決定な値は、完全一致ではなく「形式（UUID + 拡張子）」で見ます。

分岐があるなら、通常系を複数にします。例: 「既存ファイルが無いので削除しない」。

```go
assertFn: func(t *testing.T, fileRepo *stab.UploadFileRepository, _ *stab.UserRepository) {
    t.Helper()
    require.Len(t, fileRepo.ExistsNames, 1)
    assert.Empty(t, fileRepo.DeleteFileNames)
},
```

### 8.2 異常系（依存のエラー）

Usecase が `if err != nil { return fmt.Errorf("固定メッセージ: %w", err) }` している箇所ごとに、1 ケース用意します。

```go
{
    name: "異常系_ユーザ取得失敗",
    fileRepo: &stab.UploadFileRepository{
        NewFileFunc: func(string, io.Reader) error {
            t.Fatal("NewFile must not be called")
            return nil
        },
    },
    userRepo: &stab.UserRepository{
        GetByIDFunc: func(int64) (*entity.User, error) {
            return nil, errors.New("user not found")
        },
    },
    userID:      10,
    filename:    "avatar.png",
    content:     strings.NewReader("png-bytes"),
    wantErr:     true,
    errContains: "ユーザアイコンのアップロードで存在しないユーザIDが指定されました",
    assertFn: func(t *testing.T, fileRepo *stab.UploadFileRepository, userRepo *stab.UserRepository) {
        t.Helper()
        assert.Empty(t, fileRepo.NewFileNames)
        assert.Nil(t, userRepo.UpdatedUser)
    },
}
```

早い段階の失敗では、後続メソッドが **呼ばれていないこと** を必ず見ます。  
`errContains` はラップメッセージの日本語部分を使い、内部の `errors.New("...")` だけに依存しないようにします。

### 8.3 不正値

Usecase または最初の依存が拒否する入力です。

| 種類 | 扱い |
| --- | --- |
| `userID` が `0` / 負数 | `GetByID` がエラーを返す想定。後続は呼ばない |
| 空のファイル名、拡張子なし | 実装がどう正規化するかを断言する（拒否しないなら成功ケースとして形式を見る） |
| `io.Reader` が `nil` | Usecase が検査しないなら、stab 側でエラーを返し、ラップメッセージを確認する |

```go
{
    name:     "不正値_ファイル名が空",
    filename: "",
    wantErr:  false,
    assertFn: func(t *testing.T, fileRepo *stab.UploadFileRepository, _ *stab.UserRepository) {
        t.Helper()
        gotName := fileRepo.NewFileNames[0]
        require.True(t, strings.HasSuffix(gotName, "."))
    },
}
```

「不正だから必ずエラー」と決め打ちせず、**今の契約どおり** に書きます。バリデーションを足すなら、そのときに `wantErr: true` へ更新します。

---

## 9. テストデータ（エンティティ）

Usecase テストでは DB fixture（`testutils/fixtures`）は使いません。Docker もトランザクションも不要です。

メモリ上のエンティティで足りる場合は、テストファイル内の小さなヘルパーで十分です。

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
- 同じ形が複数 Usecase で必要になったら、そのときに共通化する（最初から巨大な工場にしない）

---

## 10. アサーションの使い分け

| 手段 | 用途 |
| --- | --- |
| `wantErr` + `assert.ErrorContains` | ラップされたエラーメッセージ |
| stab の記録フィールド | 渡したファイル名、更新エンティティ、呼び出し回数 |
| `t.Fatal` を stab の関数内に置く | 到達してはいけない分岐 |
| `require` | 後続の比較の前提（nil、件数） |
| `assert` | 値の比較（複数項目を一度に見たいとき） |

Repository テスト用の `assertions.NotFound` などは、GORM エラー向けです。Usecase のラップエラーには `ErrorContains` を使います。

---

## 11. 実装とテストを往復するときの注意

### 11.1 期待値をバグに合わせない

テストが落ちたとき、まずメソッドのコメントと呼び出し順を読みます。  
「今の実装がそうなっている」だけで期待を緩めると、誤った副作用（新しいファイルを古いファイルとして消す、など）を固定してしまいます。

契約と実装が食い違っているなら **実装を直す** のが正しい対応です。

### 11.2 呼ばれてはいけない処理

早期 return のあとに `NewFile` や `UpdateByID` が走ると、副作用が残ります。  
異常系では `t.Fatal("... must not be called")` か、記録スライスが空であることを断言してください。

### 11.3 テーブル初期化時の `t`

ケース定義のクロージャが `t.Fatal` を呼ぶ場合、その `t` は親テストの `*testing.T` です。  
サブテスト単位で失敗を切りたいときは、`setup func(t *testing.T)` で stab を組み立てても構いません。最初は親の `t` で十分です。

### 11.4 インターフェース型と具象 stab

テーブルのフィールドはインターフェース型にして Usecase に渡し、`assertFn` の直前で具象型へ型アサーションします。  
これにより本番と同じコンストラクタ署名を保てます。

---

## 12. 新規 Usecase テストを足すときのチェックリスト

1. [ ] パッケージは `<usecaseパッケージ>_test`
2. [ ] 依存 stab は `internal/shared/testutils/stab` にある（ローカル再定義していない）
3. [ ] stab ファイル名は `<repositoryファイル名>_stab.go`
4. [ ] `var _ インターフェース = (*Stab型)(nil)` がある
5. [ ] テーブル駆動で `通常系` / `異常系` / `不正値` を分けている
6. [ ] 異常系はラップメッセージ（`errContains`）を見ている
7. [ ] 異常系は後続メソッドが呼ばれていないことを見ている
8. [ ] 通常系はエラー有無だけでなく、stab の記録（引数・更新内容）を見ている
9. [ ] ケース間で stab インスタンスを共有していない
10. [ ] `go test -count=1 -run "TestXxxUsecase_"` が通る（Docker 不要）

---

## 13. ディレクトリ対応表

| 役割 | 場所 |
| --- | --- |
| Usecase 実装 | `internal/application/usecase/<領域>/` |
| Usecase インターフェース | `internal/adapter/interface/usecase/` |
| Repository インターフェース | `internal/domain/interface/repository/` |
| Usecase テスト | `internal/application/usecase/<領域>/*_usecase_test.go` |
| 共有 stab | `internal/shared/testutils/stab/<repository>_stab.go` |
| DB fixture（使わない） | `internal/shared/testutils/fixtures` |
| Repository 統合テスト | `internal/infrastructure/repository/` |

```text
internal/
  application/usecase/setting_usecase/
    register_icon_usecase.go
    register_icon_usecase_test.go      ← テーブル駆動。stab を注入
  domain/interface/repository/
    user_repository.go
    upload_file_repository.go
  shared/testutils/stab/
    user_repository_stab.go
    upload_file_repository_stab.go
```

---

## 引用したデータ

本手順書は、次の実装・テスト・既存手順を根拠に作成しています。

| 種別 | パス | 参照した内容 |
| --- | --- | --- |
| 文書の体裁 | `doc/documents/test/repository/README.md` | 章立て、チェックリスト、最下部の引用表 |
| Usecase 実装 | `internal/application/usecase/setting_usecase/register_icon_usecase.go` | `GetByID` → `NewFile` → `UpdateByID` → `Exists` / `DeleteFile`、エラーの `fmt.Errorf` ラップ |
| Usecase インターフェース | `internal/adapter/interface/usecase/register_icon_usecase.go` | `UploadIcon(userID, name, content)` |
| テスト本体 | `internal/application/usecase/setting_usecase/register_icon_usecase_test.go` | 外部パッケージ、テーブル駆動、`wantErr` / `errContains` / `assertFn`、通常系・異常系・不正値 |
| User stab | `internal/shared/testutils/stab/user_repository_stab.go` | `XxxFunc`、記録フィールド、インターフェース充足チェック |
| UploadFile stab | `internal/shared/testutils/stab/upload_file_repository_stab.go` | `NewFileNames` / `ExistsNames` / `DeleteFileNames` による呼び出し記録 |
| User インターフェース | `internal/domain/interface/repository/user_repository.go` | `GetByID` / `UpdateByID` |
| UploadFile インターフェース | `internal/domain/interface/repository/upload_file_repository.go` | `Exists` / `NewFile` / `DeleteFile` |
| Repository テストとの分担 | `doc/documents/test/repository/README.md` | Usecase の分岐は Repository テストでは扱わない、という境界 |
