```text
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
        cancel_submit[提出キャンセル]
        answer_comment[トピックコメント]
        meeting_comment[ミーティングコメント]
    end

    subgraph setting["設定"]
        subgraph setting_icon["アイコン設定"]
            upload_icon["ユーザアイコンのアップロード/セット"]
            delete_icon["ユーザアイコンの削除"]
        end
    end

    subgraph management["管理者機能"]
        %% Owner: 全タブ
        %% Boss: 進捗＋ミーティングログ閲覧＋ペアリング設定（部下のみ、user_relationships参照）
        %% Mentor: メンティのミーティングログ閲覧

        subgraph mng_pairing[ペアリング管理]
            read_pairing[ペアリング閲覧]
            create_pairing[ペアリング作成]
            update_pairing[ペアリング更新]
            delete_pairing[ペアリング削除]
        end

        subgraph mng_template[テンプレート管理]
            read_template[テンプレート閲覧]
            create_template[テンプレート作成]
            update_template[テンプレート更新]
            delete_template[テンプレート削除]
        end

        subgraph mng_meeting[ミーティング管理]
            read_meeting_log[ミーティングログ閲覧]
            create_meeting[ミーティング作成]
            update_meeting[ミーティング更新]
            delete_meeting[ミーティング削除]
        end

        subgraph mng_department[部署管理]
            read_department[部署閲覧]
            create_department[部署作成]
            update_department[部署更新]
            delete_department[部署削除]
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
mentee --> cancel_submit
mentee --> answer_comment
mentee --> meeting_comment

mentor -->|pairingsを参照しmentorに指定されている場合| read
mentor -->|pairingsを参照しmentorに指定されている場合| answer_comment
mentor -->|pairingsを参照しmentorに指定されている場合| meeting_comment
mentor --> read_meeting_log

boss -->|user_relationshipsで部下に指定されている場合| read
boss -->|user_relationshipsで部下に指定されている場合| answer_comment
boss -->|user_relationshipsで部下に指定されている場合| meeting_comment
boss --> read_meeting_log
boss --> read_pairing
boss --> create_pairing
boss --> update_pairing
boss --> delete_pairing
boss --> mng_progress

owner --> read
owner --> answer_comment
owner --> meeting_comment
owner --> mng_pairing
owner --> mng_meeting
owner --> mng_template
owner --> mng_department
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