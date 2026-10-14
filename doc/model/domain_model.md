# ドメインモデル図

簡易クラス図。属性は代表のみ、メソッドは書かない。sudo モデリングの D。

和英対訳はクラス名・属性の英語が実装名（`internal/domain/entity`, `valueobject`）。

```mermaid
classDiagram
    direction TB

    class User {
        <<集約ルート User>>
        UUID uuid
        氏名 name
        メール email
        パスワード password
        アイコンURI imageURI
        権限 authority
        役職ID roleID
        部署ID departmentID
    }

    class Authority {
        <<値オブジェクト Authority>>
        general
        boss
        owner
    }

    class Role {
        名称 name
    }

    class Department {
        名称 name
    }

    class LDAPUser {
        <<外部読取モデル LDAPUser>>
        氏名 name
        メール mail
    }

    class Team {
        名称 name
        部署ID departmentID
    }

    class TeamMember {
        リーダーか isReader
    }

    class DepartmentHierarchy {
        親部署ID parentDepartmentID
        子部署ID childDepartmentID
    }

    class DepartmentRole {
        表示用役職名 name
    }

    class Pairing {
        <<集約ルート Pairing>>
        メンターID mentorID
        メンティID menteeID
    }

    class UserRelationship {
        関係種別 relationType
        起点ユーザID sourceUserID
        相手ユーザID targetUserID
    }

    class RelationshipType {
        <<値オブジェクト RelationshipType>>
        boss
        subordinate
    }

    class Template {
        <<集約ルート Template>>
        タイトル title
        作成者ID userID
    }

    class Topic {
        表示順 displayNo
        質問 question
        必須か isRequire
    }

    class Meeting {
        <<集約ルート Meeting>>
        実施日 implementedAt
        ステータス status
        テンプレートID templateID
        ペアID pairingID
    }

    class MeetingStatus {
        <<値オブジェクト MeetingStatus>>
        pending
        draft
        done
    }

    class Answer {
        本文 content
        回答者ID userID
    }

    class AnswerComment {
        本文 content
    }

    class MeetingComment {
        本文 content
    }

    class UserReferenceUser {
        参照可能ユーザID referableUserID
    }

    class UserReferenceTeam {
        参照可能チームID referableTeamID
    }

    class UserReferenceDepartment {
        参照可能部署ID referableDepartmentID
    }

    class SystemSettings {
        <<集約ルート SystemSettings>>
        定期実施週 meetingWeek
        定期実施曜日 meetingWeekday
        デフォルトテンプレートID defaultTemplateID
        ログ出力先 logOutputPath
        ログPrefix logFilePrefix
    }

    User "1" --> "1" Authority : 持つ
    User "1" --> "0..1" Role : 所属
    User "1" --> "1" Department : 所属
    LDAPUser ..> User : ログイン時同期
    Department "1" --> "*" Team : 含む
    Team "1" --> "*" TeamMember : 含む
    User "1" --> "*" TeamMember : メンバー
    Department "1" --> "*" DepartmentHierarchy : 親子
    User "1" --> "*" DepartmentRole : 役職表示
    Department "1" --> "*" DepartmentRole : ""

    User "1" --> "*" Pairing : mentor
    User "1" --> "*" Pairing : mentee
    User "1" --> "*" UserRelationship : source
    User "1" --> "*" UserRelationship : target
    UserRelationship "1" --> "1" RelationshipType : 持つ

    User "1" --> "*" Template : 作成
    Template "1" --> "1..*" Topic : 含む
    SystemSettings "0..1" --> "0..1" Template : default

    Pairing "1" --> "*" Meeting : 実施
    Template "1" --> "*" Meeting : 形式
    Meeting "1" --> "1" MeetingStatus : 持つ
    Meeting "1" --> "*" Answer : 含む
    Topic "1" --> "*" Answer : 対応
    User "1" --> "*" Answer : メンティ回答
    Answer "1" --> "*" AnswerComment : コメント
    Meeting "1" --> "*" MeetingComment : コメント
    User "1" --> "*" AnswerComment : 投稿
    User "1" --> "*" MeetingComment : 投稿

    User "1" --> "*" UserReferenceUser : 閲覧許可
    User "1" --> "*" UserReferenceTeam : 閲覧許可
    User "1" --> "*" UserReferenceDepartment : 閲覧許可
```

%% --- ルール / 制約（ドメイン知識） ---
%%
%% User
%% - authority=owner は Basic ログインに制限する運用。password は nullable（LDAP ユーザは NULL 可）
%% - 閲覧スコープ: general は自分のミーティング。メンターは pairings。上司は user_relationships の部下。owner は全件
%%
%% Pairing vs UserRelationship
%% - Pairing: mentor–mentee（1on1 参加ペア）。ミーティング作成・メンター閲覧の基準
%% - UserRelationship: boss|subordinate（組織）。上司の部下ミーティング閲覧・boss のペアリング操作範囲の基準
%% - 両者の役割を混ぜない
%%
%% Meeting 集約
%% - status: pending → draft → done。提出時は submitted を経由せず done にする
%% - done 以外は権限者以外から見えない
%% - Answer.user_id はメンティ（pairings.mentee_id）
%% - meeting_comments.content は not null
%%
%% Template 集約
%% - topics は 1 件以上必須
%%
%% SystemSettings
%% - 単一行想定。定期ミーティング自動作成（cron）とデフォルトテンプレート指定
%%
%% 永続化しない / 設計外
%% - histories, read_commentables は持たない
%% - セッションはインフラ層
%%
%% 集約と Repository（実装の目安）
%% - User: UserRepository
%% - Department / Team / DepartmentHierarchy / DepartmentRole: 各 Repository
%% - Pairing: PairingRepository
%% - UserRelationship: UserRelationshipRepository
%% - Template + Topic: TemplateRepository, TopicRepository
%% - Meeting + Answer + AnswerComment + MeetingComment: 各 Repository
%% - SystemSettings: SystemSettingsRepository
%% - UserReference*: 各 Repository
