stateDiagram-v2
[*] --> login: 初期表示(未認証時やエラー後)

login: ログイン[/login]
logout: ログアウト[/logout]
logout --> login: ３秒経過

%% ログイン判定（一般ユーザ: LDAP、Owner: BASIC）
state login_method <<choice>>
login --> login_method: ログインボタン押下
login_method --> USE_LDAP: LDAP認証（mentee/mentor/boss）
login_method --> USE_BASIC: BASIC認証（owner）
USE_LDAP --> growlink: 成功（認証OK）
USE_LDAP --> login: 失敗（認証NG）
USE_BASIC --> growlink: 成功（認証OK）
USE_BASIC --> login: 失敗（認証NG）

growlink: Growlink
state growlink {
    [*] --> home: 初期表示
    home: ホーム[/]
    note left of home: 面談サマリ一覧（meetings.status + pairings）
    state switch_change_tab <<choice>>
    home --> switch_change_tab: タブ操作
    switch_change_tab --> meeting: ミーティングタブ押下
    switch_change_tab --> setting: 設定ボタン押下
    switch_change_tab --> progress: 進捗ボタン押下（boss/owner のみ表示）
    switch_change_tab --> management: 管理ボタン押下（owner/boss/mentor。タブ内容は権限で制限）

    meeting --> [*]: 重大なエラー発生
    setting --> [*]: 重大なエラー発生
    progress --> [*]: 重大なエラー発生
    management --> [*]: 重大なエラー発生
    home --> [*]: ページ遷移

    --

    meeting: ミーティング画面
    state meeting {
        [*] --> meeting_init
        meeting_init: 初期画面[/meeting]
        note right of meeting_answer: 本人かつ未記入のものは入力欄を活性
        meeting_answer: 入力画面[/meeting/:uuid]
        note right of meeting_view: 閲覧権限保持者(メンター、上司、管理者は入力欄を非活性)
        meeting_view: 閲覧画面[/meeting/:uuid]
        note right of meeting_setting: 実施予定日の更新
        meeting_setting: ミーティング設定画面[/meeting/:uuid/setting]

        note right of toast_meeting_controll: 操作結果に対しての通知を右上に表示（成功/失敗など
        note right of toast_meeting_controll: ユーザの × ボタンの押下か、３秒経過でFadeout。複数あればStack
        toast_meeting_controll: トースト通知

        meeting_init --> meeting_answer: 記入ボタン押下
        meeting_init --> meeting_view: 閲覧ボタン押下
        meeting_init --> meeting_setting: 設定タン押下

        meeting_init --> meeting_answer: 更新ボタン押下
        meeting_view --> meeting_answer: 更新ボタン押下
        meeting_answer --> meeting_init: 戻るボタンの押下
        meeting_view --> meeting_init: 戻るボタンの押下
        meeting_setting --> meeting_init: 戻るボタンの押下

        meeting_answer --> toast_meeting_controll: 下書き保存ボタン押下[POST /api/v1/meeting/:uuid/draft]（status: draft）
        meeting_answer --> toast_meeting_controll: 提出ボタン押下[POST /api/v1/meeting/:uuid]（status: done）
        meeting_answer --> toast_meeting_controll: 提出キャンセルボタン押下[POST /api/v1/meeting/:uuid/draft/cancel]（status: draft）
        meeting_view --> toast_meeting_controll: コメント投稿[POST /api/v1/meeting/:uuid/comment/:topic_uuid]
        meeting_view --> toast_meeting_controll: コメント更新[PUT /api/v1/meeting/:uuid/comment/:topic_uuid]
        meeting_view --> toast_meeting_controll: コメント削除[DELETE /api/v1/meeting/:uuid/comment/:topic_uuid]
        meeting_view --> toast_meeting_controll: ミーティングコメント投稿[POST /api/v1/meeting/:uuid/meeting-comment]
        meeting_view --> toast_meeting_controll: ミーティングコメント更新[PUT /api/v1/meeting/:uuid/meeting-comment/:comment_uuid]
        meeting_view --> toast_meeting_controll: ミーティングコメント削除[DELETE /api/v1/meeting/:uuid/meeting-comment/:comment_uuid]

        meeting_setting --> toast_meeting_controll: 設定保存[PUT /api/v1/meeting/:uuid/setting]
        meeting_init--> [*]: ページ遷移
    }

    --

    setting: 設定画面
    state setting {
        [*] --> setting_init
        setting_init: 設定初期画面[/setting]
        toast_setting_controll: トースト通知
        setting_init --> toast_setting_controll: 保存ボタンの押下[PUT /api/v1/setting]
        setting_init --> toast_setting_controll: アイコンアップロード[POST /api/v1/setting/icon]
        setting_init --> toast_setting_controll: アイコン削除[DELETE /api/v1/setting/icon]
        setting_init --> [*]: ページ遷移
    }

    --

    management: 管理画面
    state management {
        [*] --> management_init
        management_init --> [*]: ページ遷移

        management_init: 管理初期画面[/management]
        note left of management_init: Owner=全タブ、Boss=進捗+ログ+部下ペア設定、Mentor=メンティのログ閲覧
        mng_meeting_tab: ミーティング管理画面
        mng_pairing_tab: ペアリング管理画面
        mng_department_tab: 部署管理画面
        mng_template_tab: テンプレート管理画面
        mng_setting_tab: 管理設定画面

        state change_mng_tab <<choice>>
        management_init --> change_mng_tab: タブ操作
        change_mng_tab --> mng_meeting_tab: ミーティング管理タブ押下
        change_mng_tab --> mng_pairing_tab: ペアリング管理タブ押下
        change_mng_tab --> mng_department_tab: 部署管理タブ押下
        change_mng_tab --> mng_template_tab: テンプレート管理タブ押下
        change_mng_tab --> mng_setting_tab: 管理設定タブ押下

        note left of toast_management_controll: 操作結果のトースト通知
        toast_management_controll: トースト通知

        state mng_meeting_tab {
            [*] --> mng_meeting_tab_init
            mng_meeting_tab_init --> [*]
            note left of mng_meeting_tab_init: 今月分のミーティングのリスト表示、次回のミーティング自動登録の年月日表示（system_settingsで指定した 第2週の木曜日などの指定日を表示）
            mng_meeting_tab_init: タブ初期画面[/management/meeting]
            note left of mng_meeting_detail: 過去に行われたミーティングの内容を閲覧（ミーティング詳細の情報を表示）
            mng_meeting_detail: ミーティングログ詳細閲覧画面[/management/meeting/:uuid]
            note left of mng_meeting_log: 過去に行われたミーティングのリストを閲覧
            note left of mng_meeting_log: 部署名、チーム(SMやリーダー)や従業員名指定でのフィルタリング
            mng_meeting_log: ミーティングログ閲覧画面[/management/meeting/log]

            mng_meeting_tab_init --> mng_meeting_detail: 詳細ボタン押下[GET /api/v1/management/meeting/:uuid]
            mng_meeting_tab_init --> mng_meeting_log: ログ閲覧ボタン押下[GET /api/v1/management/meeting/log]
            mng_meeting_tab_init --> toast_management_controll: ミーティング作成[POST /api/v1/management/meeting]
            mng_meeting_tab_init --> toast_management_controll: ミーティング更新[PUT /api/v1/management/meeting/:uuid]
            mng_meeting_tab_init --> toast_management_controll: ミーティング削除[DELETE /api/v1/management/meeting/:uuid]
        }

        state mng_pairing_tab {
            [*] --> mng_pairing_tab_init
            mng_pairing_tab_init --> [*]
            mng_pairing_tab_init: タブ初期画面[/management/pairing]
            mng_pairing_create: ペアリング新規作成画面[/management/pairing/create]
            mng_pairing_edit: ペアリング編集画面[/management/pairing/:uuid/edit]

            mng_pairing_department: 部単位表示画面[/management/pairing/department]
            mng_pairing_department_detail: 部署詳細画面[/management/pairing/department/:uuid]
            note left of mng_pairing_department_edit: 部署単位でペアリングを決定する画面（部署管理タブとは無関係）。チーム(SM)単位で一気に従業員登録を実施可能
            mng_pairing_department_edit: 部署編集画面[/management/pairing/department/:uuid/edit]

            note left of modal_mng_pairing_detail: ペアリングの詳細（メンター-メンティーの情報）を表示するモーダル
            modal_mng_pairing_detail: ペアリング詳細モーダル
            note left of modal_mng_pairing_delete: ペアリング削除の最終確認をおこなうモーダル
            modal_mng_pairing_delete: ペアリング削除最終確認モーダル

            mng_pairing_tab_init --> mng_pairing_create: 新規作成ボタン押下
            mng_pairing_tab_init --> mng_pairing_edit: 編集ボタン押下
            mng_pairing_tab_init --> modal_mng_pairing_detail: 詳細ボタン押下
            mng_pairing_tab_init --> modal_mng_pairing_delete: 削除ボタン押下
            mng_pairing_tab_init --> mng_pairing_department: 表示ボタンを押下
            mng_pairing_tab_init --> mng_pairing_department_edit: 編集ボタンを押下
            mng_pairing_create --> toast_management_controll: 作成ボタン押下[POST /api/v1/management/pairing]
            mng_pairing_edit --> toast_management_controll: 保存ボタン押下[PUT /api/v1/management/pairing/:uuid]

            modal_mng_pairing_detail --> mng_pairing_edit: 編集ボタン押下
            modal_mng_pairing_detail --> modal_mng_pairing_delete: 削除ボタン押下
            modal_mng_pairing_delete --> toast_management_controll: 削除確定ボタン押下[DELETE /api/v1/management/pairing/:uuid]

            mng_pairing_department --> mng_pairing_department_detail: 詳細ボタンを押下
            mng_pairing_department --> mng_pairing_department_edit: 編集ボタンを押下
            mng_pairing_department_detail --> mng_pairing_department_edit: 編集ボタンを押下
            mng_pairing_department_edit --> modal_mng_pairing_delete: 削除ボタン押下
            mng_pairing_department_edit --> toast_management_controll: 更新ボタン押下[PUT /api/v1/management/pairing/department/:uuid]

        }

        state mng_department_tab {
            [*] --> mng_department_tab_init
            mng_department_tab_init --> [*]
            mng_department_tab_init: タブ初期画面[/management/department]
            mng_department_detail: 部署詳細表示画面[/management/department/:uuid]
            mng_department_create: 部署登録画面[/management/department/create]
            mng_department_edit: 部署情報登録画面[/management/department/edit]

            note left of modal_mng_department_delete: 部署削除の最終確認をおこなうモーダル
            modal_mng_department_delete: 部署削除最終確認モーダル

            mng_department_tab_init --> mng_department_detail: 詳細ボタンの押下
            mng_department_tab_init --> mng_department_create: 新規作成ボタンの押下
            mng_department_tab_init --> mng_department_edit: 編集ボタンの押下
            mng_department_detail --> mng_department_edit: 編集ボタンの押下
            mng_department_detail --> modal_mng_department_delete: 削除ボタンの押下

            mng_department_create --> toast_management_controll: 作成ボタンの押下[POST /api/v1/management/department]
            mng_department_edit --> toast_management_controll: 保存ボタンの押下[PUT /api/v1/management/department/:uuid]
            modal_mng_department_delete --> toast_management_controll: 完全に削除ボタンの押下[DELETE /api/v1/management/department/:uuid]
        }

        state mng_template_tab {
            [*] --> mng_template_tab_init
            mng_template_tab_init --> [*]

            mng_template_tab_init: タブ初期画面[/management/template]
            mng_template_edit: テンプレート編集画面[/management/template/:uuid/edit]
            mng_template_create: テンプレート新規作成画面[/management/template/create]

            note left of modal_mng_template_detail: テンプレートの内容を表示するモーダル
            modal_mng_template_detail: 詳細表示モーダル
            note left of modal_mng_template_delete: 本当に削除してもいいのかを質問するモーダル
            modal_mng_template_delete: 削除最終確認モーダル

            mng_template_tab_init --> mng_template_create: 新規作成ボタン押下
            mng_template_tab_init --> mng_template_edit: 更新ボタン押下
            mng_template_tab_init --> modal_mng_template_delete: 削除ボタン押下
            mng_template_tab_init --> modal_mng_template_detail: 詳細ボタン押下[GET /api/v1/management/template/:uuid]
            mng_template_create --> toast_management_controll: 保存ボタン押下[POST /api/v1/management/template]
            mng_template_edit --> toast_management_controll: 更新ボタン押下[PUT /api/v1/management/template/:uuid]

            modal_mng_template_detail --> mng_template_edit: 更新ボタン押下
            modal_mng_template_detail --> modal_mng_template_delete: 削除ボタン押下
            modal_mng_template_delete --> toast_management_controll: 削除確定ボタン押下[DELETE /api/v1/management/template/:uuid]
        }

        state mng_setting_tab {
            [*] --> mng_setting_tab_init
            mng_setting_tab_init --> [*]
            mng_setting_tab_init: タブ初期画面[/management/setting]

            note left of mng_setting_meeting_tab: 定期的な面談実施予定日の入力、デフォルトテンプレートの設定（system_settings）
            mng_setting_meeting_tab: ミーティング設定タブ
            mng_setting_meeting_tab --> toast_management_controll: 保存ボタンを押下[PUT /api/v1/management/setting/meeting]

            note left of mng_setting_log_tab: ログの出力先やログファイルの取得、ログファイルのPrefixの設定など
            mng_setting_log_tab: ログ設定タブ
            mng_setting_log_tab --> toast_management_controll: 保存ボタンを押下[PUT /api/v1/management/setting/log]
        }

    }

    --

    progress: 進捗管理画面[/management/progress]
    state progress {
        [*] --> progress_init
        note left of progress_init: boss/owner のみアクセス可。department/team/pairing 単位で集計・表示切替
        progress_init: 進捗管理初期画面
        progress_init-->[*]: ページ遷移
    }
}
error: エラー画面
growlink--> error: 重大なエラー発生
error-->login: ログイン画面に戻るボタン押下
growlink --> logout: ログアウトボタン押下
