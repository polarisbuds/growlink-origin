# システム関連図（System Context）

開発するシステムと、関わりのあるアクター・外部システムとの関連。sudo モデリングの S。

```mermaid
graph TB
    mentee((メンティ\nLDAP / authority=general))
    mentor((メンター\nLDAP / authority=general))
    boss((上司\nLDAP / authority=boss))
    owner((管理者\nauthority=owner))

    subgraph growlink["1on1ミーティングシステム Growlink"]
        mpa[MPA + Alpine.js]
        api[JSON API / セッション Cookie]
        files[ユーザアイコンファイル]
        cron[定期ミーティング生成ジョブ]
    end

    ldap[LDAPサーバー]
    pg[(PostgreSQL)]
    logfiles[監査ログファイル\nsystem_settings.log_output_path]

    mentee --> mpa
    mentor --> mpa
    boss --> mpa
    owner --> mpa
    mpa --> api

    mentee -.->|ログイン| ldap
    mentor -.->|ログイン| ldap
    boss -.->|ログイン| ldap
    owner -.->|Basic: users.password 照合\nLDAP も利用可| ldap

    api --> pg
    api --> files
    api --> logfiles
    cron --> pg
    cron -.->|meeting_week / meeting_weekday / default_template_id| pg
```

%% 境界メモ
%% - mentee / mentor / boss は LDAP 認証。owner は Basic 認証（password を持つユーザ）。Owner ロールは Basic に制限する運用
%% - セッションはインフラ層（histories テーブルは持たない）
%% - pairings = mentor–mentee（1on1）。user_relationships = 上司–部下（組織）。閲覧範囲は別途 user_reference_* で拡張可
%% - meetings.status が done 以外のミーティングは、権限者以外から見えない
%% - 画面遷移の詳細は doc/model/old/statedisplaytransition.md を参考（SUDO 正本ではない）
