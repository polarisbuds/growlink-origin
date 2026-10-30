# オブジェクト図

ドメインモデルの具体例。sudo モデリングの O。

シナリオ: 4月の 1on1 でメンティ田中が下書き保存後に提出する。メンター佐藤がトピックにコメントし、上司鈴木が進捗画面で部署単位の完了率を確認する。

```mermaid
classDiagram
    direction TB

    class 田中 {
        <<User>>
        uuid = u-tanaka
        name = 田中 太郎
        email = tanaka@example.com
        authority = general
        departmentID = 営業部
    }

    class 佐藤 {
        <<User>>
        uuid = u-sato
        name = 佐藤 花子
        email = sato@example.com
        authority = general
    }

    class 鈴木 {
        <<User>>
        uuid = u-suzuki
        name = 鈴木 部長
        email = suzuki@example.com
        authority = boss
    }

    class 山田 {
        <<User>>
        uuid = u-yamada
        name = 山田 管理
        email = admin@example.com
        password = ハッシュ済み
        authority = owner
    }

    class P1 {
        <<Pairing>>
        mentorID = 佐藤
        menteeID = 田中
    }

    class R1 {
        <<UserRelationship>>
        relationType = subordinate
        sourceUserID = 田中
        targetUserID = 鈴木
    }

    class T1 {
        <<Template>>
        title = 月次1on1標準
        userID = 山田
    }

    class Topic1 {
        <<Topic>>
        displayNo = 1
        question = 今月の成果は？
        isRequire = true
    }

    class M1 {
        <<Meeting>>
        uuid = mt-2026-04
        implementedAt = 2026-04-10
        status = done
        templateID = T1
        pairingID = P1
    }

    class A1 {
        <<Answer>>
        content = KPIを達成しました
        userID = 田中
    }

    class AC1 {
        <<AnswerComment>>
        content = 具体的な数字も共有してほしい
        userID = 佐藤
    }

    class SS {
        <<SystemSettings>>
        meetingWeek = 2
        meetingWeekday = 4
        defaultTemplateID = T1
    }

    佐藤 --> P1 : mentor
    田中 --> P1 : mentee
    田中 --> R1 : source
    鈴木 --> R1 : target（上司）
    山田 --> T1 : 作成
    T1 --> Topic1
    P1 --> M1
    T1 --> M1
    M1 --> A1
    Topic1 --> A1
    田中 --> A1 : 回答者
    A1 --> AC1
    佐藤 --> AC1 : コメント
    SS --> T1 : デフォルト
    鈴木 ..> M1 : 部下の done を進捗集計で参照
```

%% この時点の状態
%% 1. system_settings により第2週木曜が次回実施日としてホームに表示される（cron で meetings 生成は別タイミング）
%% 2. 田中が draft 保存後 POST /api/v1/meeting/:uuid で status=done に提出
%% 3. 佐藤は pairings 経由で閲覧し、answer_comments を POST
%% 4. 鈴木は user_relationships 経由で田中の done ミーティングを閲覧可。/management/progress で部署集計
%% 5. 山田は owner。全ミーティング・全管理タブ・Basic ログイン
