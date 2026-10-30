# Repository テスト実装手順書

この文書は、`internal/infrastructure/repository` 配下の Repository に対する統合テストの書き方をまとめたものです。  
新しい Repository メソッドを追加したとき、または既存メソッドのテストを足すときに、この手順に沿って実装してください。

---

## 1. このテストで検証すること

Repository テストは **本物の PostgreSQL（Testcontainers）** に対して、次を確認します。

| 確認すること | 例 |
| --- | --- |
| 正常系の読み書き | `GetByID` で保存した行が取れる、`UpdateByID` 後に再取得すると値が変わる |
| 異常系・境界 | 存在しない ID、`nil` 入力、一意制約違反、外部キー違反 |
| SQL / GORM の契約 | 並び順、日付範囲、ソフトデリート、preload の有無 |

**やらないこと**

- Usecase の分岐や画面の挙動は、ここでは扱わない（Usecase テスト側）
- モックで SQL を差し替えない。DB の制約・型・デフォルト値まで含めて検証する

---

## 2. 前提

- Docker が起動していること（PostgreSQL コンテナを起動するため）
- テストパッケージは `internal/infrastructure/repository` に既に `TestMain` がある（コンテナの後始末）
- 共有ヘルパーは `internal/shared/testutils` 配下にある

ローカルでの実行例:

```bash
# パッケージ全体
go test -v -count=1 ./internal/infrastructure/repository/

# 対象 Repository だけ
go test -v -count=1 ./internal/infrastructure/repository/ -run "TestUserRepository_"
```

初回はコンテナ起動とマイグレーションのため、数十秒かかることがあります。

---

## 3. 全体の流れ（毎回この順）

```text
1. 対象メソッドの契約を決める（成功時 / 未存在時 / nil など）
2. テストファイルを repository_test パッケージで用意する
3. コンストラクタ用のヘルパーを書く
4. fixture が無ければ公開関数を足す（GORM で INSERT。Repository は使わない）
5. テスト関数を「1ケース = 1関数」で書く
6. 実行して失敗したら、テストではなく実装の契約を疑う
```

各ステップの詳細を以降に書きます。

---

## 4. ファイルとパッケージの約束

| 項目 | ルール |
| --- | --- |
| ファイル名 | `<対象>_repository_test.go`（例: `user_repository_test.go`） |
| パッケージ名 | **`repository_test`**（外部テスト。実装と同じ `repository` にはしない） |
| テスト関数名 | `Test<型名>_<メソッド名>` または `Test<型名>_<メソッド名>_<観点>` |
| データ作成 | `fixtures` 経由。テスト内で生 SQL を書かない |

`repository_test` にする理由は、公開 API（`NewXxxRepository` とメソッド）だけを叩くためです。未公開フィールドに依存すると、実装変更でテストが壊れやすくなります。

---

## 5. テストの骨格（必ずこの形）

どのテストも、次の 3 行から始めます。

```go
tx := testutils.Begin(t)          // 共有 DB 上でトランザクション開始。終了時に Rollback
repo := newXxxRepo(t, tx)         // その tx を渡した Repository
want := fixtures.Xxx(t, tx, ...)  // 期待データの準備（Repository ではなく GORM）
```

### 5.1 なぜトランザクションか

`testutils.Begin` は共有 PostgreSQL に接続し、テスト専用のトランザクションを開始します。`t.Cleanup` で必ず Rollback するため、他テストに行が残りません。

```go
func TestUserRepository_GetByID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserRepo(t, tx)

	want := fixtures.User(t, tx, "get")

	got, err := repo.GetByID(int64(want.ID))
	assertions.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, want.ID, got.ID)
}
```

**注意:** PostgreSQL のシーケンス（`id` の採番）は Rollback しても戻りません。  
「存在しない ID」には `1` ではなく、衝突しにくい大きな値（例: `9_999_999`）を使ってください。

### 5.2 Repository 生成ヘルパー

コンストラクタ呼び出しはヘルパーにまとめます。`t.Helper()` を付け、コンテキストは `context.Background()`、DB はテストの `tx` を渡します。

```go
func newUserRepo(t *testing.T, tx *gorm.DB) *repository.UserRepository {
	t.Helper()
	return repository.NewUserRepository(context.Background(), tx)
}
```

本番コードと同じコンストラクタを使うことで、「テスト専用の組み立て」を避けます。

---

## 6. Fixture の作り方

Fixture は **テストデータ工場** です。Repository の `Insert` ではなく、`tx.Create`（GORM）で依存関係ごと投入します。

これにより「Repository の Insert が壊れていても Get のテストが書ける」状態を保てます。

### 6.1 公開関数を足す

未公開の `createXxx` だけがある場合は、テストから呼べる公開関数を追加します。

```go
// User inserts a User with Layer0 master dependencies (Layer1).
func User(t *testing.T, tx *gorm.DB, prefix string) *entity.User {
	t.Helper()
	return createUser(t, tx, prefix)
}
```

- `t.Helper()` を付ける（失敗時のスタックがテスト本体を指す）
- 一意制約がある列（`email` など）は `time.Now().UnixNano()` などで衝突を避ける
- FK がある場合は、ER の依存順で親（Role / Department など）を先に作る

### 6.2 条件を変えたいエンティティは Builder にする

日付や関連の組み合わせが多い場合は、Meeting のように Builder にします。

```go
want := fixtures.Meeting(t, tx).WithImplementedAt(todayNoon).Create()
```

単純な 1 行投入で足りるなら、`fixtures.User(t, tx, "get")` のような関数で十分です。最初から Builder にする必要はありません。

### 6.3 Fixture でやってはいけないこと

- テスト対象の Repository で INSERT する（循環した検証になる）
- 共有の固定 ID をハードコードする
- トランザクション外の `database.DB` に直接書く（他テストと干渉する）

---

## 7. アサーションの使い分け

共通判定は `internal/shared/testutils/assertions` を使います。同じ判定を各テストにコピーしないでください。

| 関数 | 用途 |
| --- | --- |
| `assertions.NoError` | 成功系の `err` |
| `assertions.NotFound` | `errors.Is(err, gorm.ErrRecordNotFound)` |
| `assertions.EqualUUID` / `NonZeroUUID` | UUID |
| `assertions.UniqueViolation` | UNIQUE 制約 |
| `assertions.ForeignKeyViolation` | FK 制約 |
| `require` / `assert`（testify） | 値比較、nil チェック、件数 |

失敗したら続きを実行したくないときは `require`、複数項目を一度に見たいときは `assert` です。

ドメイン固有の検証（preload 済みか、ID リストが一致するか）が増えたら、`assertions` に関数を足してテストを薄く保ちます。

---

## 8. ケースの切り方

メソッドごとに、少なくとも次を検討します。全部必須ではありません。契約にあるものだけ書きます。

### 8.1 取得（GetByID など）

**成功:** fixture で作った行を、同じ `tx` 上の Repository で読む。比較するのはメソッドが返すフィールドだけ（preload しないなら関連は見ない）。

```go
want := fixtures.User(t, tx, "get")

got, err := repo.GetByID(int64(want.ID))
assertions.NoError(t, err)
require.NotNil(t, got)
assert.Equal(t, want.ID, got.ID)
assertions.EqualUUID(t, want.UUID, got.UUID)
assert.Equal(t, want.Name, got.Name)
```

**未存在:** 存在しない ID を渡し、`assertions.NotFound` で契約を固定する。

```go
_, err := repo.GetByID(9_999_999)
assertions.NotFound(t, err)
```

`BaseRepository.First` は「行が無いとき `(nil, nil)`」を返す実装です。Repository 側で `gorm.ErrRecordNotFound` に包んでいないと、このテストは失敗します。  
その場合はテストを「成功扱い」に緩めず、**実装の契約を直す**のが正しい対応です。

### 8.2 更新（UpdateByID など）

1. fixture で作成  
2. メモリ上のエンティティを書き換える  
3. Repository で更新  
4. **同じ Repository の Get で再読込**して永続化を確認する  

```go
want := fixtures.User(t, tx, "update")
imageURI := "icon.png"
want.Name = "updated-name"
want.ImageURI = &imageURI

err := repo.UpdateByID(want)
assertions.NoError(t, err)

got, err := repo.GetByID(int64(want.ID))
assertions.NoError(t, err)
require.NotNil(t, got.ImageURI)
assert.Equal(t, "updated-name", got.Name)
assert.Equal(t, imageURI, *got.ImageURI)
assert.Equal(t, want.Email, got.Email) // 変えていない列は維持される
```

`nil` を受けて何もしない契約があるなら、別関数にします。

```go
func TestUserRepository_UpdateByID_Nil(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserRepo(t, tx)

	err := repo.UpdateByID(nil)
	assertions.NoError(t, err)
}
```

### 8.3 制約違反

UNIQUE / FK は DB にしか出せないので、Repository テストの出番です。

```go
err := repo.Insert(&entity.Meeting{UUID: existing.UUID, /* 他は有効な FK */})
assertions.UniqueViolation(t, err)
```

### 8.4 トランザクション分離

「Rollback した行は別トランザクションから見えない」ことを確認する例です（Meeting の既存テスト）。  
通常の CRUD では不要です。分離の契約を明示したいときだけ書きます。

---

## 9. 実装とテストを往復するときの注意

テストを先に書くと、実装の誤用が表面化します。よくあるものは次のとおりです。

### 9.1 `Update` と `Updates` を混同しない

`BaseRepository.Update` の第 1 引数は **カラム名** です。ログ用の文字列を渡すと、存在しない列への UPDATE になります。

```go
// 誤り: 第1引数がカラム名として使われる
repo.baseRepository.Update(fmt.Sprintf("Update user(%d)", user.ID), user, "id = ?", user.ID)

// 1カラムだけ更新するとき
repo.baseRepository.Update("image_uri", filename, "id = ?", user.ID)

// 構造体で複数カラムを更新するとき
repo.baseRepository.Updates(*user, "id = ?", user.ID)
```

### 9.2 未存在を「成功の nil」にしない

Usecase が `if err != nil` だけで未存在を扱うなら、Repository は `ErrRecordNotFound` を返す必要があります。  
`(nil, nil)` だと Usecase が先に進み、nil 参照になります。

### 9.3 テストが落ちたら

1. SQL ログで、実際に発行された文を見る  
2. 「テストの期待」と「メソッドのコメント / Usecase の使い方」が一致しているか確認する  
3. 一致しているのに落ちるなら **実装を直す**（期待値を実装のバグに合わせない）

---

## 10. 新規 Repository を足すときのチェックリスト

1. [ ] テストファイルのパッケージは `repository_test`
2. [ ] `newXxxRepo(t, tx)` ヘルパーがある
3. [ ] 各テストの先頭で `testutils.Begin(t)` している
4. [ ] データ作成は `fixtures`（GORM）。対象 Repository の Insert に依存していない
5. [ ] 成功系と、契約にある異常系（未存在 / nil / 制約）を分けている
6. [ ] 更新系は Get で再取得して永続化を見ている
7. [ ] 未存在 ID にシーケンスと衝突しうる小さな値を使っていない
8. [ ] `go test -count=1 -run "TestXxxRepository_"` が通る
9. [ ] Docker が無い環境ではこのパッケージのテストは動かない、と認識している

---

## 11. ディレクトリ対応表

| 役割 | 場所 |
| --- | --- |
| テスト本体 | `internal/infrastructure/repository/*_repository_test.go` |
| コンテナ後始末 | `internal/infrastructure/repository/main_test.go` |
| トランザクション開始 | `internal/shared/testutils/setup.go` の `Begin` |
| Rollback 登録 | `internal/shared/testutils/transaction` |
| DB / マイグレーション | `internal/shared/testutils/database` |
| PostgreSQL コンテナ | `internal/shared/testutils/containers` |
| テストデータ | `internal/shared/testutils/fixtures` |
| 共通アサーション | `internal/shared/testutils/assertions` |

---

## 引用したデータ

本手順書は、次の実装・テスト・設定を根拠に作成しています。

| 種別 | パス | 参照した内容 |
| --- | --- | --- |
| 手順の雛形（取得・制約・分離） | `internal/infrastructure/repository/meeting_repository_test.go` | `testutils.Begin`、`newMeetingRepo`、`fixtures.Meeting`、`assertions` の使い方、UNIQUE / FK、Rollback 分離 |
| 手順の雛形（Get / Update） | `internal/infrastructure/repository/user_repository_test.go` | `GetByID` / `GetByID_NotFound` / `UpdateByID` / `UpdateByID_Nil` |
| 実装契約 | `internal/infrastructure/repository/user_repository.go` | 未存在時の `gorm.ErrRecordNotFound`、`Updates` による更新、`nil` はエラーにしない |
| 汎用 CRUD | `internal/infrastructure/repository/base_repository.go` | `First` の未存在時 `(nil, nil)`、`Update`（カラム単位）、`Updates`（構造体） |
| テストエントリ | `internal/infrastructure/repository/main_test.go` | `TestMain` でのコンテナ `Terminate` |
| トランザクション | `internal/shared/testutils/setup.go` | `Begin` が共有 DB + Rollback 付き tx を返す |
| Rollback | `internal/shared/testutils/transaction/transaction.go` | `t.Cleanup` で `Rollback` |
| Fixture（単純） | `internal/shared/testutils/fixtures/user.go` | 公開 `User` と Layer0 依存の `createUser` |
| Fixture（Builder） | `internal/shared/testutils/fixtures/meeting.go` | `Meeting(t, tx).With....Create()` |
| アサーション | `internal/shared/testutils/assertions/assertions.go` | `NoError` / `NotFound` / `UniqueViolation` / `ForeignKeyViolation` など |
| コンテナ定義 | `internal/shared/testutils/containers/postgres.go` | テスト用 PostgreSQL イメージと共有コンテナ |
| CI | `.gitea/workflows/gitea-ci.yml` | `go test -v -race ./...` と `DOCKER_HOST` |
| Usecase からの契約 | `internal/application/usecase/setting_usecase/register_icon_usecase.go` | `GetByID` 失敗をエラーとして扱う、`UpdateByID` で `ImageURI` を更新する |
