```mermaid
graph LR;

user((ユーザ\n（全権限共通）))

mentee((メンティーユーザー))
mentor((メンターユーザー))
boss((上司ユーザー))
owner((管理者ユーザー))

subgraph growlink["1on1ミーティングシステム"]
    subgraph ldap["LDAP認証"]
        ldap_login[ログイン]
    end

    subgraph basic["BASIC認証"]
        basic_login[ログイン]
    end

    subgraph meetings["ミーティング"]
        answer[回答]
        read[閲覧]
        draft[下書き]
        submit[提出]
        topic_comment[トピックコメント]
        meeting_comment[ミーティングコメント]
    end

    subgraph setting["設定"]
        subgraph setting_icon["アイコン設定"]
            upload_icon["ユーザアイコンのアップロード/セット"]
            delete_icon["ユーザアイコンの削除"]
        end
    end

    subgraph management["管理者機能"]
        subgraph mng_peering[ペアリング管理]
            read_peering[ペアリング閲覧] %% 一覧表示
            create_peering[ペアリング作成]
            update_peering[ペアリング更新]
            delete_peering[ペアリング削除]
        end

        subgraph mng_template[テンプレート管理]
            read_template[テンプレート閲覧] %% 一覧表示や詳細表示
            create_template[テンプレート作成]
            update_template[テンプレート更新]
            delete_template[テンプレート削除]
        end

        subgraph mng_meeting[ミーティング管理]
            create_meeting[ミーティング作成]
            update_meeting[ミーティング更新]
            delete_meeting[ミーティング削除]
        end

        subgraph read_history[履歴管理]
            meeting_history[ミーティング履歴閲覧]
            user_history[ユーザ履歴閲覧] %% ログイン/操作(簡易)
        end

        subgraph mng_progress[進捗管理]
            read_progress[進捗閲覧]
        end
    end
end

subgraph LDAPServer["LDAPサーバー"]
    search[ユーザー検索]
    ad_login[ADアカウント認証]
end

mentee ~~~ mentor
mentor ~~~ boss
boss ~~~ owner

mentee --> answer
mentee --> read
mentee --> draft
mentee --> submit
mentee --> topic_comment
mentee --> meeting_comment

mentor -->|peeringを参照しmentorに指定されている場合| read
mentor -->|peeringを参照しmentorに指定されている場合| topic_comment
mentor -->|peeringを参照しmentorに指定されている場合| meeting_comment

boss -->|authority_idがbossかつ、mentee.team_idのチームリーダーに指定されている場合| read
boss -->|authority_idがbossかつ、mentee.team_idのチームリーダーに指定されている場合| topic_comment
boss -->|authority_idがbossかつ、mentee.team_idのチームリーダーに指定されている場合| meeting_comment
boss --> mng_progress

owner --> read
owner --> topic_comment
owner --> meeting_comment
owner --> mng_peering
owner --> mng_meeting
owner --> mng_template
owner --> mng_progress

mentee --> ldap_login
mentor --> ldap_login
boss --> ldap_login
ldap_login --> search
ldap_login --> ad_login

owner --> basic_login

user --> upload_icon
user --> delete_icon
```