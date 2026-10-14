```mermaid
erDiagram

%% 各テーブルは bigint id PKを所持
%% 各テーブルは created_at, updated_at, deleted_atを所持
%% 各テーブルは uuid uuid UK(Unique Key)を所持
%% tables { bigint id PK, datetime created_at, datetime updated_at, deleted_at } が存在する前提とする。
%% pending|draft|done / general|boss|owner / boss|subordinate は Domain Value Object（列挙体）。マスタテーブルは持たない。
%% VO列は smallint で保持する。

%% ユーザ系
users {
    text name "not null"
    text email "not null"
    text password "nullable。LDAP/Basic認証。Owner（Authority == Owner）はBasic認証のみに制限。"
    text image_uri "nullable。イメージ画像の保存先"
    smallint authority "not null, general|boss|owner。閲覧範囲（メンター：Pairings参照、上司：UserRelationships参照（部下のミーティングのみ）、オーナー：全件）"

    bigint role_id FK
    bigint department_id FK
}
roles ||--o{ users: ""
departments ||--|{ users: ""

team_members {
    boolean is_reader "チームリーダーフラグ"

    bigint team_id FK
    bigint member_id FK
}
teams ||--o{ team_members: "チームID"
users ||--o{ team_members: "メンバーID"

departments {
    text name "not null"
}

roles {
    text name "not null"
}

teams {
    text name "not null"

    bigint department_id FK
}
departments ||--|{ teams: "部署＞チーム（SM以下のチームを想定）"

pairings {

    bigint mentor_id FK
    bigint mentee_id FK
}
%% pairings: 1on1ミーティングの参加ペア（メンター-メンティ）を管理。ミーティング作成・メンター閲覧の基準。
users ||--o{ pairings: "Mentor"
users ||--o{ pairings: "Mentee"

%% ミーティング系
%% ステータス done 以外では他ユーザから見られないようにする
meetings {
    datetime implemented_at "nullable。ミーティング実施日"
    smallint status "not null, pending|draft|done"

    bigint template_id FK
    bigint pairing_id FK
}
templates ||--o{ meetings: "ミーティングにはテンプレートの指定が必須"
pairings ||--o{ meetings: "ミーティング参加ユーザはペア指定（閲覧権限とは別管理）"

templates {
    text title UK

    bigint user_id FK "作成者のユーザID"
}
users ||--o{ templates: "テンプレートには作成者が存在"

topics {
    number display_no "表示順序 1~N"
    text question "not null"
    bool is_require "not null"

    bigint template_id FK
}
templates ||--|{ topics: "テンプレートに対して、１以上のトピックが存在する"

answers {
    text content "not null"

    bigint meeting_id FK
    bigint topic_id FK
    bigint user_id FK "回答者ユーザID（メンティ）"
}
meetings ||--o{ answers: ""
topics ||--o{ answers: ""
users ||--o{ answers: "メンティの回答"

meeting_comments {
    text content "not null"

    bigint meeting_id FK
    bigint user_id FK
}
meetings ||--o{ meeting_comments: ""
users ||--o{ meeting_comments: ""

answer_comments {
    text content

    bigint answer_id FK
    bigint user_id FK
}
answers ||--o{ answer_comments: ""
users ||--o{ answer_comments: ""

%% user_relationships: 上司-部下の組織関係を管理。pairings とは役割が異なる（mentor-mentee は pairings で管理）。
user_relationships {
    smallint relation_type "not null, boss|subordinate"

    bigint source_user_id FK
    bigint target_user_id FK
}
%% Source側からTargetを見た場合のリレーションを記述
%% 例: RelationType=Boss, source_user_id=SubordinateID, target_user_id=BossID
users ||--o{ user_relationships: "source_user_id"
users ||--o{ user_relationships: "target_user_id"

system_settings {
    number meeting_week "定期実施：第N週"
    number meeting_weekday "定期実施：曜日（0=日曜）"
    bigint default_template_id FK "デフォルトテンプレート"
    text log_output_path "ログ出力先"
    text log_file_prefix "ログファイルPrefix"
}
templates ||--o| system_settings: "default_template_id"

department_hierarchies {
    bigint parent_department_id FK
    bigint child_department_id FK
}
departments ||--o{ epartment_hierarchies: "parent_department_id"
departments ||--o{ epartment_hierarchies: "child_department_id"

department_roles {
    text name "ロール名。表示上利用するのみのため自由入力"

    bigint user_id FK
    bigint department_id FK
}
users ||--o{ department_roles: "役職名"
departments ||o{ department_roles: "部署ID"

user_reference_users {
    bigint user_id FK
    bigint referable_user_id FK
}
users ||--o{ user_reference_users: ""
users ||--o{ user_reference_users: "参照可能ユーザID"

user_reference_teams {
    bigint user_id FK
    bigint referable_team_id FK
}
users ||--o{ user_reference_teams: ""
users ||--o{ user_reference_teams: "参照可能チームID"

user_reference_departments {
    bigint user_id FK
    bigint referable_departments_id FK
}
users ||--o{ user_reference_departments: ""
users ||--o{ user_reference_departments: "参照可能部署ID"

```
