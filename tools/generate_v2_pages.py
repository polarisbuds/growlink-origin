#!/usr/bin/env python3
"""Generate all samples/v2 static HTML pages."""
from __future__ import annotations

import os
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
OUT = ROOT / "samples" / "v2" / "pages"

TAILWIND = '<script src="https://cdn.tailwindcss.com"></script>'
DAISY = '<link href="https://cdn.jsdelivr.net/npm/daisyui@4.12.14/dist/full.min.css" rel="stylesheet" />'
LUCIDE = '<script src="https://unpkg.com/lucide@latest/dist/umd/lucide.min.js"></script>'
ALPINE = '<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.14.1/dist/cdn.min.js"></script>'


def js_depth(rel: str) -> str:
    depth = rel.count("/")
    return "../" * (depth + 1) + "static/js/"


def head_block(rel: str, title: str, app: bool = True) -> str:
    js = js_depth(rel)
    if not app:
        return f"""<!DOCTYPE html>
<html lang="ja" data-theme="light">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>{title} | Growlink</title>
  {TAILWIND}
  {DAISY}
  <script src="{js}root.js"></script>
  <script src="{js}sample-data.js"></script>
  <script src="{js}paths.js"></script>
  {LUCIDE}
  <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.14.1/dist/cdn.min.js"></script>
  <style>[x-cloak]{{display:none!important}}</style>
</head>"""
    return f"""<!DOCTYPE html>
<html lang="ja" data-theme="light">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>{title} | Growlink</title>
  {TAILWIND}
  {DAISY}
  <script src="{js}root.js"></script>
  <script src="{js}sample-data.js"></script>
  <script src="{js}paths.js"></script>
  <script src="{js}page-utils.js"></script>
  <script src="{js}store.js"></script>
  <script src="{js}app-shell.js"></script>
  {LUCIDE}
  <script src="{js}lucide-boot.js"></script>
  <style>[x-cloak]{{display:none!important}}</style>
</head>"""


def app_body(attrs: str, content: str, page_script: str) -> str:
    return f"""<body {attrs} x-cloak>
  <div id="page-content">
{content}
  </div>
{page_script}
  <script>SampleShell.mount();</script>
  {ALPINE}
</body>
</html>"""


def write_page(rel: str, html: str) -> None:
    path = OUT / rel
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(html, encoding="utf-8")
    print(f"  wrote {rel}")


def login_page() -> str:
    rel = "login.html"
    return head_block(rel, "ログイン", app=False) + f"""
<body class="min-h-screen bg-base-200" x-data="loginPage()" x-cloak>
  <div class="flex min-h-screen items-center justify-center p-4">
    <div class="w-full max-w-md rounded-3xl border border-base-300 bg-base-100 p-8 shadow-sm">
      <div class="mb-6 flex items-center gap-3">
        <div class="grid h-10 w-10 place-items-center rounded-xl bg-rose-900 text-white"><i data-lucide="sprout" class="h-5 w-5"></i></div>
        <div>
          <h1 class="font-display text-xl font-semibold">Growlink</h1>
          <p class="text-sm text-base-content/60">1on1 ミーティング</p>
        </div>
      </div>
      <div class="tabs tabs-boxed mb-4">
        <button type="button" class="tab" :class="mode === 'ldap' && 'tab-active'" @click="mode='ldap'">LDAP</button>
        <button type="button" class="tab" :class="mode === 'basic' && 'tab-active'" @click="mode='basic'">BASIC</button>
      </div>
      <template x-if="error"><p class="mb-4 text-sm text-error" x-text="error"></p></template>
      <form class="space-y-4" @submit.prevent="submit()">
        <label class="form-control w-full">
          <span class="label-text">メールアドレス</span>
          <input type="email" class="input input-bordered w-full" x-model="email" required />
        </label>
        <label class="form-control w-full">
          <span class="label-text">パスワード</span>
          <input type="password" class="input input-bordered w-full" x-model="password" required />
        </label>
        <button type="submit" class="btn w-full border-rose-900 bg-rose-900 text-white">ログイン</button>
      </form>
    </div>
  </div>
  <script>
    function loginPage() {{
      return {{
        mode: 'ldap',
        email: '',
        password: '',
        error: '',
        submit() {{
          const ok = (this.mode === 'ldap' && this.email && this.password) || (this.mode === 'basic' && this.email === 'admin@example.com' && this.password === 'owner');
          if (!ok) {{ this.error = '認証に失敗しました。'; return; }}
          window.location.href = window.SAMPLE_PATHS.home;
        }},
      }};
    }}
    document.addEventListener('alpine:initialized', () => lucide.createIcons());
  </script>
  {ALPINE}
</body>
</html>"""


def logout_page() -> str:
    rel = "logout.html"
    return head_block(rel, "ログアウト", app=False) + f"""
<body class="flex min-h-screen items-center justify-center bg-base-200" x-data="{{}}" x-init="setTimeout(() => window.location.href = window.SAMPLE_PATHS.login, 3000)" x-cloak>
  <div class="rounded-3xl border border-base-300 bg-base-100 p-8 text-center shadow-sm">
    <i data-lucide="log-out" class="mx-auto h-10 w-10 text-rose-900"></i>
    <h1 class="mt-4 font-display text-xl font-semibold">ログアウトしました</h1>
    <p class="mt-2 text-sm text-base-content/60">3秒後にログイン画面へ移動します...</p>
  </div>
  <script>document.addEventListener('alpine:initialized', () => lucide.createIcons());</script>
  {ALPINE}
</body>
</html>"""


def error_page() -> str:
    rel = "error.html"
    return head_block(rel, "エラー", app=False) + f"""
<body class="flex min-h-screen items-center justify-center bg-base-200" x-cloak>
  <div class="rounded-3xl border border-base-300 bg-base-100 p-8 text-center shadow-sm">
    <i data-lucide="alert-triangle" class="mx-auto h-10 w-10 text-error"></i>
    <h1 class="mt-4 font-display text-xl font-semibold">エラーが発生しました</h1>
    <p class="mt-2 text-sm text-base-content/60">しばらくしてから再度お試しください。</p>
    <a href="#" onclick="window.location.href=window.SAMPLE_PATHS.login;return false;" class="btn mt-6 border-rose-900 bg-rose-900 text-white">ログイン画面に戻る</a>
  </div>
  <script>lucide.createIcons();</script>
</body>
</html>"""


def home_page() -> str:
    rel = "index.html"
    content = """    <div x-data="homePage()">
      <h1 class="font-display text-2xl font-semibold">ホーム</h1>
      <p class="mt-2 text-sm text-base-content/60">面談サマリ一覧</p>
      <div class="mt-6 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        <template x-for="item in cards" :key="item.uuid">
          <div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
            <div class="flex items-start justify-between gap-3">
              <div>
                <p class="text-sm text-base-content/60" x-text="item.mentor + ' / ' + item.mentee"></p>
                <h2 class="mt-1 font-semibold" x-text="item.template"></h2>
              </div>
              <span class="badge" :class="statusBadge(item.status)" x-text="statusLabel(item.status)"></span>
            </div>
            <p class="mt-3 text-sm">実施予定日: <span x-text="item.implementedAt"></span></p>
            <div class="mt-4 flex flex-wrap gap-2">
              <a :href="paths.meetingAnswer(item.uuid)" class="btn btn-sm btn-outline">記入</a>
              <a :href="paths.meetingView(item.uuid)" class="btn btn-sm btn-ghost">閲覧</a>
              <a :href="paths.meetingSetting(item.uuid)" class="btn btn-sm btn-ghost">設定</a>
            </div>
          </div>
        </template>
      </div>
    </div>"""
    script = """  <script>
    function homePage() {
      return {
        cards: window.SAMPLE_DATA.meetings,
        paths: window.SAMPLE_PATHS,
        statusBadge,
        statusLabel,
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="ホーム" data-nav-active="home" data-breadcrumbs=\'[{"label":"ホーム","key":"home"}]\' data-management-tabs="false" data-show-upload="false"'
    return head_block(rel, "ホーム") + app_body(attrs, content, script)


def setting_page() -> str:
    rel = "setting.html"
    content = """    <div x-data="settingPage()">
      <h1 class="font-display text-2xl font-semibold">個人設定</h1>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
        <div class="flex items-center gap-4">
          <div class="grid h-16 w-16 place-items-center rounded-full bg-rose-900 text-xl font-semibold text-white">田</div>
          <div>
            <button type="button" class="btn btn-sm btn-outline" @click="openIconUpload()">アイコンを変更</button>
            <button type="button" class="btn btn-sm btn-ghost" @click="deleteIcon()">アイコンを削除</button>
          </div>
        </div>
        <div class="mt-6 grid gap-4 md:grid-cols-2">
          <label class="form-control"><span class="label-text">氏名</span><input class="input input-bordered" x-model="form.name" /></label>
          <label class="form-control"><span class="label-text">メール</span><input class="input input-bordered" x-model="form.email" /></label>
        </div>
        <div class="mt-6"><button type="button" class="btn border-rose-900 bg-rose-900 text-white" @click="save()">保存</button></div>
      </div>
    </div>"""
    script = """  <script>
    function settingPage() {
      const user = window.SAMPLE_DATA.users[0];
      return {
        form: { name: user.name, email: user.email },
        save() { apiAlert('PUT', '/api/v1/setting', this.form); },
        deleteIcon() { apiAlert('DELETE', '/api/v1/setting/icon', null); },
        openIconUpload() {
          const stack = document.body._x_dataStack;
          if (stack && stack[0]) stack[0].openUploadModal(true);
        },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="設定" data-nav-active="setting" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"設定"}]\' data-management-tabs="false" data-show-upload="true" data-icon-upload="true"'
    return head_block(rel, "設定") + app_body(attrs, content, script)


def meeting_index() -> str:
    rel = "meeting/index.html"
    content = """    <div x-data="meetingListPage()">
      <h1 class="font-display text-2xl font-semibold">ミーティング</h1>
      <div class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
        <table class="table">
          <thead><tr><th>テンプレート</th><th>メンター</th><th>メンティ</th><th>実施予定日</th><th>状態</th><th></th></tr></thead>
          <tbody>
            <template x-for="m in meetings" :key="m.uuid">
              <tr>
                <td x-text="m.template"></td><td x-text="m.mentor"></td><td x-text="m.mentee"></td><td x-text="m.implementedAt"></td>
                <td><span class="badge" :class="statusBadge(m.status)" x-text="statusLabel(m.status)"></span></td>
                <td class="space-x-1">
                  <a :href="paths.meetingAnswer(m.uuid)" class="btn btn-xs btn-outline">記入</a>
                  <a :href="paths.meetingView(m.uuid)" class="btn btn-xs btn-ghost">閲覧</a>
                  <a :href="paths.meetingSetting(m.uuid)" class="btn btn-xs btn-ghost">設定</a>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>"""
    script = """  <script>
    function meetingListPage() {
      return { meetings: window.SAMPLE_DATA.meetings, paths: window.SAMPLE_PATHS, statusBadge, statusLabel };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="ミーティング" data-nav-active="meeting" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"ミーティング","key":"meeting"}]\''
    return head_block(rel, "ミーティング") + app_body(attrs, content, script)


def meeting_answer() -> str:
    rel = "meeting/answer.html"
    content = """    <div x-data="meetingAnswerPage()">
      <div class="flex items-center justify-between gap-4">
        <h1 class="font-display text-2xl font-semibold">ミーティング記入</h1>
        <a href="index.html" class="btn btn-ghost btn-sm">戻る</a>
      </div>
      <div class="mt-6 space-y-4">
        <template x-for="topic in topics" :key="topic.uuid">
          <div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
            <p class="text-sm text-base-content/60">Q<span x-text="topic.displayNo"></span> <span x-show="topic.isRequire" class="text-error">*</span></p>
            <p class="mt-1 font-medium" x-text="topic.question"></p>
            <textarea class="textarea textarea-bordered mt-3 w-full" rows="3" x-model="answers[topic.uuid]" :disabled="!editable"></textarea>
          </div>
        </template>
      </div>
      <div class="mt-6 flex flex-wrap gap-2">
        <button type="button" class="btn btn-outline" :disabled="!editable" @click="draft()">下書き保存</button>
        <button type="button" class="btn border-rose-900 bg-rose-900 text-white" :disabled="!editable" @click="submit()">提出</button>
        <button type="button" class="btn btn-ghost" :disabled="!editable" @click="cancelSubmit()">提出キャンセル</button>
      </div>
    </div>"""
    script = """  <script>
    function meetingAnswerPage() {
      const uuid = getUUID();
      const meeting = window.SAMPLE_DATA.meetings.find(m => m.uuid === uuid) || window.SAMPLE_DATA.meetings[0];
      const tpl = window.SAMPLE_DATA.templates[0];
      const answers = {};
      tpl.topics.forEach(t => { answers[t.uuid] = ''; });
      return {
        uuid, topics: tpl.topics, answers,
        editable: meeting.status === 'pending' || meeting.status === 'draft',
        draft() { apiAlert('POST', `/api/v1/meeting/${uuid}/draft`, { answers: this.answers }); },
        submit() { apiAlert('POST', `/api/v1/meeting/${uuid}`, { answers: this.answers }); },
        cancelSubmit() { apiAlert('POST', `/api/v1/meeting/${uuid}/draft/cancel`, null); },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="記入" data-nav-active="meeting" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"ミーティング","key":"meeting"},{"label":"記入"}]\''
    return head_block(rel, "記入") + app_body(attrs, content, script)


def meeting_view() -> str:
    rel = "meeting/view.html"
    content = """    <div x-data="meetingViewPage()">
      <div class="flex items-center justify-between gap-4">
        <h1 class="font-display text-2xl font-semibold">ミーティング閲覧</h1>
        <div class="flex gap-2">
          <a :href="paths.meetingAnswer(uuid)" class="btn btn-outline btn-sm">更新</a>
          <a href="index.html" class="btn btn-ghost btn-sm">戻る</a>
        </div>
      </div>
      <div class="mt-6 space-y-4">
        <template x-for="topic in topics" :key="topic.uuid">
          <div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
            <p class="font-medium" x-text="topic.question"></p>
            <textarea class="textarea textarea-bordered mt-3 w-full" rows="3" disabled x-text="answers[topic.uuid]"></textarea>
            <div class="mt-3 rounded-2xl bg-base-200/50 p-3 text-sm">
              <p class="font-medium">コメント</p>
              <template x-for="c in comments.filter(x => x.topicUuid === topic.uuid)" :key="c.uuid">
                <div class="mt-2 flex items-start justify-between gap-2">
                  <div><span class="font-medium" x-text="c.user"></span>: <span x-text="c.content"></span></div>
                  <div class="flex gap-1">
                    <button type="button" class="btn btn-ghost btn-xs" @click="updateComment(topic.uuid, c.uuid)">編集</button>
                    <button type="button" class="btn btn-ghost btn-xs" @click="deleteComment(topic.uuid, c.uuid)">削除</button>
                  </div>
                </div>
              </template>
              <button type="button" class="btn btn-xs btn-outline mt-2" @click="postComment(topic.uuid)">コメント投稿</button>
            </div>
          </div>
        </template>
      </div>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
        <h2 class="font-semibold">ミーティングコメント</h2>
        <template x-for="mc in meetingComments" :key="mc.uuid">
          <div class="mt-3 flex items-start justify-between gap-2 text-sm">
            <div><span class="font-medium" x-text="mc.user"></span>: <span x-text="mc.content"></span></div>
            <div class="flex gap-1">
              <button type="button" class="btn btn-ghost btn-xs" @click="updateMeetingComment(mc.uuid)">編集</button>
              <button type="button" class="btn btn-ghost btn-xs" @click="deleteMeetingComment(mc.uuid)">削除</button>
            </div>
          </div>
        </template>
        <button type="button" class="btn btn-sm btn-outline mt-4" @click="postMeetingComment()">コメント投稿</button>
      </div>
    </div>"""
    script = """  <script>
    function meetingViewPage() {
      const uuid = getUUID();
      const tpl = window.SAMPLE_DATA.templates[0];
      const answers = {};
      tpl.topics.forEach(t => { answers[t.uuid] = '（回答サンプル）'; });
      return {
        uuid, paths: window.SAMPLE_PATHS, topics: tpl.topics, answers,
        comments: window.SAMPLE_DATA.answerComments,
        meetingComments: window.SAMPLE_DATA.meetingComments,
        postComment(topicUuid) { apiAlert('POST', `/api/v1/meeting/${uuid}/comment/${topicUuid}`, { content: '新規コメント' }); },
        updateComment(topicUuid, commentUuid) { apiAlert('PUT', `/api/v1/meeting/${uuid}/comment/${topicUuid}`, { commentUuid, content: '更新コメント' }); },
        deleteComment(topicUuid, commentUuid) { apiAlert('DELETE', `/api/v1/meeting/${uuid}/comment/${topicUuid}`, { commentUuid }); },
        postMeetingComment() { apiAlert('POST', `/api/v1/meeting/${uuid}/meeting-comment`, { content: '新規' }); },
        updateMeetingComment(commentUuid) { apiAlert('PUT', `/api/v1/meeting/${uuid}/meeting-comment/${commentUuid}`, { content: '更新' }); },
        deleteMeetingComment(commentUuid) { apiAlert('DELETE', `/api/v1/meeting/${uuid}/meeting-comment/${commentUuid}`, null); },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="閲覧" data-nav-active="meeting" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"ミーティング","key":"meeting"},{"label":"閲覧"}]\''
    return head_block(rel, "閲覧") + app_body(attrs, content, script)


def meeting_setting() -> str:
    rel = "meeting/setting.html"
    content = """    <div x-data="meetingSettingPage()">
      <div class="flex items-center justify-between gap-4">
        <h1 class="font-display text-2xl font-semibold">ミーティング設定</h1>
        <a href="index.html" class="btn btn-ghost btn-sm">戻る</a>
      </div>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
        <label class="form-control max-w-xs"><span class="label-text">実施予定日</span><input type="date" class="input input-bordered" x-model="implementedAt" /></label>
        <button type="button" class="btn mt-6 border-rose-900 bg-rose-900 text-white" @click="save()">保存</button>
      </div>
    </div>"""
    script = """  <script>
    function meetingSettingPage() {
      const uuid = getUUID();
      const meeting = window.SAMPLE_DATA.meetings.find(m => m.uuid === uuid) || window.SAMPLE_DATA.meetings[0];
      return {
        implementedAt: meeting.implementedAt,
        save() { apiAlert('PUT', `/api/v1/meeting/${uuid}/setting`, { implementedAt: this.implementedAt }); },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="ミーティング設定" data-nav-active="meeting" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"ミーティング","key":"meeting"},{"label":"設定"}]\''
    return head_block(rel, "ミーティング設定") + app_body(attrs, content, script)


def mgmt_hub() -> str:
    rel = "management/index.html"
    content = """    <div>
      <h1 class="font-display text-2xl font-semibold">管理</h1>
      <p class="mt-2 text-sm text-base-content/60">上部タブから各管理画面へ移動できます。</p>
      <div class="mt-6 grid gap-4 md:grid-cols-2">
        <a href="meeting/index.html" class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm hover:border-rose-900">ミーティング管理</a>
        <a href="pairing/index.html" class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm hover:border-rose-900">ペアリング管理</a>
        <a href="department/index.html" class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm hover:border-rose-900" x-show="$store.ui.role === 'owner'">部署管理</a>
        <a href="team/index.html" class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm hover:border-rose-900" x-show="$store.ui.role === 'owner'">チーム管理</a>
      </div>
    </div>"""
    attrs = 'data-layout="app" data-title="管理" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"}]\' data-management-tabs="true" data-management-tab=""'
    return head_block(rel, "管理") + app_body(attrs, content, "  <script></script>")


def progress_page() -> str:
    rel = "management/progress.html"
    content = """    <div x-data="progressPage()">
      <h1 class="font-display text-2xl font-semibold">進捗管理</h1>
      <div class="mt-4 flex flex-wrap gap-2">
        <button type="button" class="btn btn-sm" :class="unit==='department' && 'btn-active'" @click="unit='department'">部署</button>
        <button type="button" class="btn btn-sm" :class="unit==='team' && 'btn-active'" @click="unit='team'">チーム</button>
        <button type="button" class="btn btn-sm" :class="unit==='pairing' && 'btn-active'" @click="unit='pairing'">ペア</button>
      </div>
      <div class="mt-6 grid gap-4 md:grid-cols-3">
        <template x-for="card in cards" :key="card.label">
          <div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
            <p class="text-sm text-base-content/60" x-text="card.label"></p>
            <p class="mt-2 text-3xl font-semibold text-rose-900" x-text="card.value"></p>
            <p class="mt-1 text-sm" x-text="card.sub"></p>
          </div>
        </template>
      </div>
    </div>"""
    script = """  <script>
    function progressPage() {
      return {
        unit: 'department',
        get cards() {
          const data = {
            department: [{ label: '開発部 完了率', value: '78%', sub: '完了 7 / 未記入 2' }, { label: '営業部 完了率', value: '65%', sub: '完了 5 / 未記入 3' }, { label: '全社平均', value: '72%', sub: '前月比 +4%' }],
            team: [{ label: 'SM-A 完了率', value: '80%', sub: '完了 4 / 未記入 1' }, { label: 'SM-B 完了率', value: '75%', sub: '完了 3 / 未記入 1' }, { label: '平均', value: '77%', sub: '' }],
            pairing: [{ label: '佐藤-田中', value: '100%', sub: '完了 3' }, { label: '田中-山本', value: '50%', sub: '下書き 1' }, { label: '平均', value: '75%', sub: '' }],
          };
          return data[this.unit];
        },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="進捗" data-nav-active="progress" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"進捗","key":"managementProgress"}]\''
    return head_block(rel, "進捗") + app_body(attrs, content, script)


def mgmt_setting() -> str:
    rel = "management/setting.html"
    content = """    <div x-data="mgmtSettingPage()">
      <h1 class="font-display text-2xl font-semibold">管理設定</h1>
      <div class="tabs tabs-boxed mt-4">
        <button type="button" class="tab" :class="tab==='meeting' && 'tab-active'" @click="tab='meeting'">ミーティング</button>
        <button type="button" class="tab" :class="tab==='log' && 'tab-active'" @click="tab='log'">ログ</button>
      </div>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm" x-show="tab==='meeting'">
        <div class="grid gap-4 md:grid-cols-2">
          <label class="form-control"><span class="label-text">実施週（第N週）</span><input type="number" class="input input-bordered" x-model.number="meeting.week" /></label>
          <label class="form-control"><span class="label-text">曜日</span><select class="select select-bordered" x-model.number="meeting.weekday"><option value="1">月</option><option value="2">火</option><option value="3">水</option><option value="4">木</option><option value="5">金</option></select></label>
          <label class="form-control md:col-span-2"><span class="label-text">デフォルトテンプレート</span><select class="select select-bordered" x-model="meeting.defaultTemplateUuid"><template x-for="t in templates" :key="t.uuid"><option :value="t.uuid" x-text="t.title"></option></template></select></label>
        </div>
        <button type="button" class="btn mt-6 border-rose-900 bg-rose-900 text-white" @click="saveMeeting()">保存</button>
      </div>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm" x-show="tab==='log'">
        <div class="grid gap-4 md:grid-cols-2">
          <label class="form-control md:col-span-2"><span class="label-text">出力先</span><input class="input input-bordered" x-model="log.outputPath" /></label>
          <label class="form-control"><span class="label-text">Prefix</span><input class="input input-bordered" x-model="log.filePrefix" /></label>
        </div>
        <button type="button" class="btn mt-6 border-rose-900 bg-rose-900 text-white" @click="saveLog()">保存</button>
      </div>
    </div>"""
    script = """  <script>
    function mgmtSettingPage() {
      const s = window.SAMPLE_DATA.systemSettings;
      return {
        tab: 'meeting',
        templates: window.SAMPLE_DATA.templates,
        meeting: { week: s.meetingWeek, weekday: s.meetingWeekday, defaultTemplateUuid: s.defaultTemplateUuid },
        log: { outputPath: s.logOutputPath, filePrefix: s.logFilePrefix },
        saveMeeting() { apiAlert('PUT', '/api/v1/management/setting/meeting', this.meeting); },
        saveLog() { apiAlert('PUT', '/api/v1/management/setting/log', this.log); },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="管理設定" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"設定"}]\' data-management-tabs="true" data-management-tab="setting"'
    return head_block(rel, "管理設定") + app_body(attrs, content, script)


def mgmt_meeting_index() -> str:
    rel = "management/meeting/index.html"
    content = """    <div x-data="mgmtMeetingPage()">
      <div class="flex flex-wrap items-center justify-between gap-4">
        <h1 class="font-display text-2xl font-semibold">ミーティング管理</h1>
        <div class="flex gap-2">
          <a href="log.html" class="btn btn-outline btn-sm">ログ閲覧</a>
          <button type="button" class="btn btn-sm border-rose-900 bg-rose-900 text-white" @click="create()">作成</button>
        </div>
      </div>
      <p class="mt-2 text-sm text-base-content/60">次回自動登録日: 第2週 木曜日</p>
      <div class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
        <table class="table">
          <thead><tr><th>メンター</th><th>メンティ</th><th>実施日</th><th>状態</th><th></th></tr></thead>
          <tbody>
            <template x-for="m in meetings" :key="m.uuid">
              <tr>
                <td x-text="m.mentor"></td><td x-text="m.mentee"></td><td x-text="m.implementedAt"></td>
                <td><span class="badge" :class="statusBadge(m.status)" x-text="statusLabel(m.status)"></span></td>
                <td class="space-x-1">
                  <a :href="'detail.html?uuid=' + m.uuid" class="btn btn-xs btn-ghost">詳細</a>
                  <button type="button" class="btn btn-xs btn-outline" @click="update(m.uuid)">更新</button>
                  <button type="button" class="btn btn-xs btn-ghost" @click="remove(m.uuid)">削除</button>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>"""
    script = """  <script>
    function mgmtMeetingPage() {
      return {
        meetings: window.SAMPLE_DATA.meetings,
        statusBadge, statusLabel,
        create() { apiAlert('POST', '/api/v1/management/meeting', { pairingUuid: 'pair-001' }); },
        update(uuid) { apiAlert('PUT', `/api/v1/management/meeting/${uuid}`, { implementedAt: '2026-09-11' }); },
        remove(uuid) { apiAlert('DELETE', `/api/v1/management/meeting/${uuid}`, null); },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="ミーティング管理" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"ミーティング"}]\' data-management-tabs="true" data-management-tab="meeting"'
    return head_block(rel, "ミーティング管理") + app_body(attrs, content, script)


def mgmt_meeting_detail() -> str:
    rel = "management/meeting/detail.html"
    content = """    <div x-data="mgmtMeetingDetailPage()">
      <h1 class="font-display text-2xl font-semibold">ミーティング詳細</h1>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
        <p><span class="text-base-content/60">メンター:</span> <span x-text="meeting.mentor"></span></p>
        <p class="mt-2"><span class="text-base-content/60">メンティ:</span> <span x-text="meeting.mentee"></span></p>
        <p class="mt-2"><span class="text-base-content/60">実施日:</span> <span x-text="meeting.implementedAt"></span></p>
        <a href="index.html" class="btn btn-ghost btn-sm mt-6">戻る</a>
      </div>
    </div>"""
    script = """  <script>
    function mgmtMeetingDetailPage() {
      const uuid = getUUID();
      return { meeting: window.SAMPLE_DATA.meetings.find(m => m.uuid === uuid) || window.SAMPLE_DATA.meetings[0] };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="ミーティング詳細" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"詳細"}]\' data-management-tabs="true" data-management-tab="meeting"'
    return head_block(rel, "ミーティング詳細") + app_body(attrs, content, script)


def mgmt_meeting_log() -> str:
    rel = "management/meeting/log.html"
    content = """    <div x-data="mgmtMeetingLogPage()">
      <h1 class="font-display text-2xl font-semibold">ミーティングログ</h1>
      <div class="mt-4 flex flex-wrap gap-2">
        <select class="select select-bordered select-sm"><option>部署: すべて</option><option>開発部</option></select>
        <select class="select select-bordered select-sm"><option>チーム: すべて</option><option>SM-A</option></select>
        <input class="input input-bordered input-sm" placeholder="氏名で検索" />
      </div>
      <div class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
        <table class="table">
          <thead><tr><th>日付</th><th>メンター</th><th>メンティ</th><th>状態</th><th></th></tr></thead>
          <tbody>
            <template x-for="m in meetings" :key="m.uuid">
              <tr>
                <td x-text="m.implementedAt"></td><td x-text="m.mentor"></td><td x-text="m.mentee"></td>
                <td><span class="badge" :class="statusBadge(m.status)" x-text="statusLabel(m.status)"></span></td>
                <td><a :href="'detail.html?uuid=' + m.uuid" class="btn btn-xs btn-ghost">詳細</a></td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>"""
    script = """  <script>
    function mgmtMeetingLogPage() {
      return { meetings: window.SAMPLE_DATA.meetings, statusBadge, statusLabel };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="ミーティングログ" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"ログ"}]\' data-management-tabs="true" data-management-tab="meeting"'
    return head_block(rel, "ミーティングログ") + app_body(attrs, content, script)


def mgmt_pairing_index() -> str:
    rel = "management/pairing/index.html"
    content = """    <div x-data="mgmtPairingPage()">
      <div class="flex flex-wrap items-center justify-between gap-4">
        <h1 class="font-display text-2xl font-semibold">ペアリング管理</h1>
        <div class="flex gap-2">
          <a href="department/index.html" class="btn btn-outline btn-sm">部単位表示</a>
          <a href="create.html" class="btn btn-sm border-rose-900 bg-rose-900 text-white">新規作成</a>
        </div>
      </div>
      <div class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
        <table class="table">
          <thead><tr><th>メンター</th><th>メンティ</th><th>部署</th><th>チーム</th><th></th></tr></thead>
          <tbody>
            <template x-for="p in pairings" :key="p.uuid">
              <tr>
                <td x-text="p.mentor"></td><td x-text="p.mentee"></td><td x-text="p.department"></td><td x-text="p.team"></td>
                <td class="space-x-1">
                  <button type="button" class="btn btn-xs btn-ghost" @click="openDetail(p)">詳細</button>
                  <a :href="'edit.html?uuid=' + p.uuid" class="btn btn-xs btn-outline">編集</a>
                  <button type="button" class="btn btn-xs btn-ghost" @click="openDelete(p)">削除</button>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
      <dialog class="modal" :class="detailOpen && 'modal-open'">
        <div class="modal-box rounded-3xl" x-show="selected">
          <h3 class="font-semibold">ペアリング詳細</h3>
          <p class="mt-4" x-text="selected ? selected.mentor + ' → ' + selected.mentee : ''"></p>
          <div class="modal-action">
            <a :href="selected ? 'edit.html?uuid=' + selected.uuid : '#'" class="btn btn-outline">編集</a>
            <button type="button" class="btn btn-ghost" @click="detailOpen=false; openDelete(selected)">削除</button>
            <button type="button" class="btn" @click="detailOpen=false">閉じる</button>
          </div>
        </div>
        <form method="dialog" class="modal-backdrop bg-black/60 backdrop-blur-md"><button type="button" @click="detailOpen=false">close</button></form>
      </dialog>
      <dialog class="modal" :class="deleteOpen && 'modal-open'">
        <div class="modal-box rounded-3xl">
          <h3 class="font-semibold">削除の確認</h3>
          <p class="mt-2 text-sm">このペアリングを削除しますか？</p>
          <div class="modal-action">
            <button type="button" class="btn btn-ghost" @click="deleteOpen=false">キャンセル</button>
            <button type="button" class="btn btn-error" @click="confirmDelete()">削除確定</button>
          </div>
        </div>
        <form method="dialog" class="modal-backdrop bg-black/60 backdrop-blur-md"><button type="button" @click="deleteOpen=false">close</button></form>
      </dialog>
    </div>"""
    script = """  <script>
    function mgmtPairingPage() {
      return {
        pairings: window.SAMPLE_DATA.pairings,
        selected: null, detailOpen: false, deleteOpen: false,
        openDetail(p) { this.selected = p; this.detailOpen = true; },
        openDelete(p) { this.selected = p; this.deleteOpen = true; },
        confirmDelete() {
          apiAlert('DELETE', `/api/v1/management/pairing/${this.selected.uuid}`, null);
          this.deleteOpen = false;
        },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="ペアリング" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"ペアリング"}]\' data-management-tabs="true" data-management-tab="pairing"'
    return head_block(rel, "ペアリング") + app_body(attrs, content, script)


def mgmt_pairing_create() -> str:
    rel = "management/pairing/create.html"
    content = """    <div x-data="mgmtPairingFormPage('create')">
      <h1 class="font-display text-2xl font-semibold">ペアリング新規作成</h1>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm space-y-4">
        <label class="form-control"><span class="label-text">メンター</span><input class="input input-bordered" x-model="form.mentor" /></label>
        <label class="form-control"><span class="label-text">メンティ</span><input class="input input-bordered" x-model="form.mentee" /></label>
        <label class="form-control"><span class="label-text">部署</span><input class="input input-bordered" x-model="form.department" /></label>
        <button type="button" class="btn border-rose-900 bg-rose-900 text-white" @click="save()">作成</button>
        <a href="index.html" class="btn btn-ghost">戻る</a>
      </div>
    </div>"""
    script = """  <script>
    function mgmtPairingFormPage(mode) {
      const uuid = getUUID('pair-001');
      const existing = window.SAMPLE_DATA.pairings.find(p => p.uuid === uuid);
      return {
        form: existing ? { ...existing } : { mentor: '', mentee: '', department: '', team: '' },
        save() {
          const path = mode === 'create' ? '/api/v1/management/pairing' : `/api/v1/management/pairing/${uuid}`;
          apiAlert(mode === 'create' ? 'POST' : 'PUT', path, this.form);
        },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="ペアリング作成" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"作成"}]\' data-management-tabs="true" data-management-tab="pairing"'
    return head_block(rel, "ペアリング作成") + app_body(attrs, content, script)


def mgmt_pairing_edit() -> str:
    rel = "management/pairing/edit.html"
    content = """    <div x-data="mgmtPairingFormPage('edit')">
      <h1 class="font-display text-2xl font-semibold">ペアリング編集</h1>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm space-y-4">
        <label class="form-control"><span class="label-text">メンター</span><input class="input input-bordered" x-model="form.mentor" /></label>
        <label class="form-control"><span class="label-text">メンティ</span><input class="input input-bordered" x-model="form.mentee" /></label>
        <button type="button" class="btn border-rose-900 bg-rose-900 text-white" @click="save()">保存</button>
        <a href="index.html" class="btn btn-ghost">戻る</a>
      </div>
    </div>"""
    attrs = 'data-layout="app" data-title="ペアリング編集" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"編集"}]\' data-management-tabs="true" data-management-tab="pairing"'
    return head_block(rel, "ペアリング編集") + app_body(attrs, content, mgmt_pairing_create().split("<script>")[2].split("</script>")[0] if False else """  <script>
    function mgmtPairingFormPage(mode) {
      const uuid = getUUID('pair-001');
      const existing = window.SAMPLE_DATA.pairings.find(p => p.uuid === uuid) || window.SAMPLE_DATA.pairings[0];
      return {
        form: { ...existing },
        save() { apiAlert('PUT', `/api/v1/management/pairing/${uuid}`, this.form); },
      };
    }
  </script>""")


def mgmt_pairing_dept_index() -> str:
    rel = "management/pairing/department/index.html"
    content = """    <div x-data="{ departments: window.SAMPLE_DATA.departments }">
      <h1 class="font-display text-2xl font-semibold">部単位ペアリング</h1>
      <div class="mt-6 grid gap-4 md:grid-cols-2">
        <template x-for="d in departments" :key="d.uuid">
          <div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
            <h2 class="font-semibold" x-text="d.name"></h2>
            <p class="mt-1 text-sm text-base-content/60" x-text="'部長: ' + d.head"></p>
            <div class="mt-4 flex gap-2">
              <a :href="'detail.html?uuid=' + d.uuid" class="btn btn-sm btn-ghost">詳細</a>
              <a :href="'edit.html?uuid=' + d.uuid" class="btn btn-sm btn-outline">編集</a>
            </div>
          </div>
        </template>
      </div>
      <a href="../index.html" class="btn btn-ghost mt-6">戻る</a>
    </div>"""
    attrs = 'data-layout="app" data-title="部単位ペアリング" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"部単位"}]\' data-management-tabs="true" data-management-tab="pairing"'
    return head_block(rel, "部単位ペアリング") + app_body(attrs, content, "  <script></script>")


def mgmt_pairing_dept_detail() -> str:
    rel = "management/pairing/department/detail.html"
    content = """    <div x-data="deptDetailPage()">
      <h1 class="font-display text-2xl font-semibold" x-text="department.name"></h1>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
        <p x-text="'部長: ' + department.head"></p>
        <p class="mt-2" x-text="'メンバー数: ' + department.memberCount"></p>
        <a :href="'edit.html?uuid=' + department.uuid" class="btn btn-outline btn-sm mt-6">編集</a>
        <a href="index.html" class="btn btn-ghost btn-sm mt-6">戻る</a>
      </div>
    </div>"""
    script = """  <script>
    function deptDetailPage() {
      const uuid = getUUID('dept-001');
      return { department: window.SAMPLE_DATA.departments.find(d => d.uuid === uuid) || window.SAMPLE_DATA.departments[0] };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="部署詳細" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"部署詳細"}]\' data-management-tabs="true" data-management-tab="pairing"'
    return head_block(rel, "部署詳細") + app_body(attrs, content, script)


def mgmt_pairing_dept_edit() -> str:
    rel = "management/pairing/department/edit.html"
    content = """    <div x-data="deptEditPage()">
      <h1 class="font-display text-2xl font-semibold">部署ペアリング編集</h1>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm space-y-4">
        <p class="text-sm text-base-content/60">チーム単位で一括登録できます</p>
        <template x-for="team in teams" :key="team.uuid">
          <div class="rounded-2xl border border-base-300 p-4">
            <p class="font-medium" x-text="team.name"></p>
            <input class="input input-bordered input-sm mt-2 w-full" placeholder="ペアリング設定" />
          </div>
        </template>
        <div class="flex gap-2">
          <button type="button" class="btn border-rose-900 bg-rose-900 text-white" @click="save()">更新</button>
          <button type="button" class="btn btn-error btn-outline" @click="deleteOpen=true">削除</button>
        </div>
      </div>
      <dialog class="modal" :class="deleteOpen && 'modal-open'">
        <div class="modal-box rounded-3xl">
          <h3 class="font-semibold">削除の確認</h3>
          <div class="modal-action">
            <button type="button" class="btn btn-ghost" @click="deleteOpen=false">キャンセル</button>
            <button type="button" class="btn btn-error" @click="remove()">削除確定</button>
          </div>
        </div>
        <form method="dialog" class="modal-backdrop bg-black/60 backdrop-blur-md"><button type="button" @click="deleteOpen=false">close</button></form>
      </dialog>
    </div>"""
    script = """  <script>
    function deptEditPage() {
      const uuid = getUUID('dept-001');
      return {
        teams: window.SAMPLE_DATA.teams,
        deleteOpen: false,
        save() { apiAlert('PUT', `/api/v1/management/pairing/department/${uuid}`, { teams: this.teams.map(t => t.uuid) }); },
        remove() { apiAlert('DELETE', `/api/v1/management/pairing/${uuid}`, null); this.deleteOpen = false; },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="部署ペアリング編集" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"編集"}]\' data-management-tabs="true" data-management-tab="pairing"'
    return head_block(rel, "部署ペアリング編集") + app_body(attrs, content, script)


def mgmt_department_index() -> str:
    rel = "management/department/index.html"
    content = """    <div x-data="{ departments: window.SAMPLE_DATA.departments }">
      <div class="flex justify-between gap-4">
        <h1 class="font-display text-2xl font-semibold">部署管理</h1>
        <a href="create.html" class="btn btn-sm border-rose-900 bg-rose-900 text-white">新規作成</a>
      </div>
      <div class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
        <table class="table">
          <thead><tr><th>部署名</th><th>部長</th><th>人数</th><th></th></tr></thead>
          <tbody>
            <template x-for="d in departments" :key="d.uuid">
              <tr>
                <td x-text="d.name"></td><td x-text="d.head"></td><td x-text="d.memberCount"></td>
                <td class="space-x-1">
                  <a :href="'detail.html?uuid=' + d.uuid" class="btn btn-xs btn-ghost">詳細</a>
                  <a :href="'edit.html?uuid=' + d.uuid" class="btn btn-xs btn-outline">編集</a>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>"""
    attrs = 'data-layout="app" data-title="部署管理" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"部署"}]\' data-management-tabs="true" data-management-tab="department"'
    return head_block(rel, "部署管理") + app_body(attrs, content, "  <script></script>")


def mgmt_department_create() -> str:
    rel = "management/department/create.html"
    content = """    <div x-data="deptFormPage('create')">
      <h1 class="font-display text-2xl font-semibold">部署登録</h1>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm space-y-4">
        <label class="form-control"><span class="label-text">部署名</span><input class="input input-bordered" x-model="form.name" /></label>
        <button type="button" class="btn border-rose-900 bg-rose-900 text-white" @click="save()">作成</button>
        <a href="index.html" class="btn btn-ghost">戻る</a>
      </div>
    </div>"""
    script = """  <script>
    function deptFormPage(mode) {
      const uuid = getUUID('dept-001');
      const existing = window.SAMPLE_DATA.departments.find(d => d.uuid === uuid);
      return {
        form: existing ? { name: existing.name } : { name: '' },
        save() {
          apiAlert(mode === 'create' ? 'POST' : 'PUT', mode === 'create' ? '/api/v1/management/department' : `/api/v1/management/department/${uuid}`, this.form);
        },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="部署登録" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"登録"}]\' data-management-tabs="true" data-management-tab="department"'
    return head_block(rel, "部署登録") + app_body(attrs, content, script)


def mgmt_department_detail() -> str:
    rel = "management/department/detail.html"
    content = """    <div x-data="deptDetailMgmtPage()">
      <h1 class="font-display text-2xl font-semibold" x-text="department.name"></h1>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
        <p x-text="'部長: ' + department.head"></p>
        <p class="mt-2" x-text="'副部長: ' + (department.deputyHead || '未設定')"></p>
        <div class="mt-6 flex flex-wrap gap-2">
          <button type="button" class="btn btn-outline btn-sm" @click="searchOpen=true">部長・副部長・部員選択</button>
          <a :href="'edit.html?uuid=' + department.uuid" class="btn btn-outline btn-sm">編集</a>
          <button type="button" class="btn btn-error btn-outline btn-sm" @click="deleteOpen=true">削除</button>
        </div>
      </div>
      <dialog class="modal" :class="searchOpen && 'modal-open'">
        <div class="modal-box rounded-3xl">
          <h3 class="font-semibold">ユーザ検索</h3>
          <input class="input input-bordered mt-4 w-full" placeholder="氏名で検索" />
          <div class="modal-action">
            <button type="button" class="btn btn-ghost" @click="searchOpen=false">閉じる</button>
            <button type="button" class="btn border-rose-900 bg-rose-900 text-white" @click="assign('head')">部長に設定</button>
            <button type="button" class="btn btn-outline" @click="assign('deputyhead')">副部長に設定</button>
            <button type="button" class="btn btn-outline" @click="assign('member')">部員に設定</button>
          </div>
        </div>
        <form method="dialog" class="modal-backdrop bg-black/60 backdrop-blur-md"><button type="button" @click="searchOpen=false">close</button></form>
      </dialog>
      <dialog class="modal" :class="deleteOpen && 'modal-open'">
        <div class="modal-box rounded-3xl">
          <h3 class="font-semibold">部署削除の確認</h3>
          <div class="modal-action">
            <button type="button" class="btn btn-ghost" @click="deleteOpen=false">キャンセル</button>
            <button type="button" class="btn btn-error" @click="remove()">完全に削除</button>
          </div>
        </div>
        <form method="dialog" class="modal-backdrop bg-black/60 backdrop-blur-md"><button type="button" @click="deleteOpen=false">close</button></form>
      </dialog>
    </div>"""
    script = """  <script>
    function deptDetailMgmtPage() {
      const uuid = getUUID('dept-001');
      return {
        department: window.SAMPLE_DATA.departments.find(d => d.uuid === uuid) || window.SAMPLE_DATA.departments[0],
        searchOpen: false, deleteOpen: false,
        assign(role) {
          const paths = { head: '/api/v1/management/department/head', deputyhead: '/api/v1/management/department/deputyhead', member: '/api/v1/management/department/member' };
          apiAlert('POST', paths[role], { departmentUuid: uuid, userUuid: 'user-mentee-001' });
          this.searchOpen = false;
        },
        remove() { apiAlert('DELETE', `/api/v1/management/department/${uuid}`, null); this.deleteOpen = false; },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="部署詳細" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"部署詳細"}]\' data-management-tabs="true" data-management-tab="department"'
    return head_block(rel, "部署詳細") + app_body(attrs, content, script)


def mgmt_department_edit() -> str:
    rel = "management/department/edit.html"
    content = """    <div x-data="deptFormPage('edit')">
      <h1 class="font-display text-2xl font-semibold">部署編集</h1>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm space-y-4">
        <label class="form-control"><span class="label-text">部署名</span><input class="input input-bordered" x-model="form.name" /></label>
        <button type="button" class="btn border-rose-900 bg-rose-900 text-white" @click="save()">保存</button>
        <a href="index.html" class="btn btn-ghost">戻る</a>
      </div>
    </div>"""
    attrs = 'data-layout="app" data-title="部署編集" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"編集"}]\' data-management-tabs="true" data-management-tab="department"'
    return head_block(rel, "部署編集") + app_body(attrs, content, """  <script>
    function deptFormPage(mode) {
      const uuid = getUUID('dept-001');
      const existing = window.SAMPLE_DATA.departments.find(d => d.uuid === uuid) || window.SAMPLE_DATA.departments[0];
      return {
        form: { name: existing.name },
        save() { apiAlert('PUT', `/api/v1/management/department/${uuid}`, this.form); },
      };
    }
  </script>""")


def mgmt_team_index() -> str:
    rel = "management/team/index.html"
    content = """    <div x-data="teamIndexPage()">
      <div class="flex justify-between gap-4">
        <h1 class="font-display text-2xl font-semibold">チーム管理</h1>
        <a href="create.html" class="btn btn-sm border-rose-900 bg-rose-900 text-white">新規作成</a>
      </div>
      <div class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
        <table class="table">
          <thead><tr><th>チーム名</th><th>リーダー</th><th>人数</th><th></th></tr></thead>
          <tbody>
            <template x-for="t in teams" :key="t.uuid">
              <tr>
                <td x-text="t.name"></td><td x-text="t.leader"></td><td x-text="t.memberCount"></td>
                <td class="space-x-1">
                  <a :href="'edit.html?uuid=' + t.uuid" class="btn btn-xs btn-outline">編集</a>
                  <button type="button" class="btn btn-xs btn-ghost" @click="openDelete(t)">削除</button>
                  <button type="button" class="btn btn-xs btn-ghost" @click="openSearch(t)">メンバー選択</button>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
      <dialog class="modal" :class="deleteOpen && 'modal-open'">
        <div class="modal-box rounded-3xl"><h3 class="font-semibold">チーム削除の確認</h3>
          <div class="modal-action"><button type="button" class="btn btn-ghost" @click="deleteOpen=false">キャンセル</button><button type="button" class="btn btn-error" @click="remove()">完全に削除</button></div>
        </div>
        <form method="dialog" class="modal-backdrop bg-black/60 backdrop-blur-md"><button type="button" @click="deleteOpen=false">close</button></form>
      </dialog>
      <dialog class="modal" :class="searchOpen && 'modal-open'">
        <div class="modal-box rounded-3xl"><h3 class="font-semibold">ユーザ検索</h3>
          <div class="modal-action"><button type="button" class="btn btn-ghost" @click="searchOpen=false">閉じる</button><button type="button" class="btn border-rose-900 bg-rose-900 text-white" @click="assignLeader()">リーダーに設定</button><button type="button" class="btn btn-outline" @click="assignMember()">メンバーに設定</button></div>
        </div>
        <form method="dialog" class="modal-backdrop bg-black/60 backdrop-blur-md"><button type="button" @click="searchOpen=false">close</button></form>
      </dialog>
    </div>"""
    script = """  <script>
    function teamIndexPage() {
      return {
        teams: window.SAMPLE_DATA.teams,
        selected: null, deleteOpen: false, searchOpen: false,
        openDelete(t) { this.selected = t; this.deleteOpen = true; },
        openSearch(t) { this.selected = t; this.searchOpen = true; },
        remove() { apiAlert('DELETE', `/api/v1/management/team/${this.selected.uuid}`, null); this.deleteOpen = false; },
        assignLeader() { apiAlert('POST', '/api/v1/management/team/leader', { teamUuid: this.selected.uuid }); this.searchOpen = false; },
        assignMember() { apiAlert('POST', '/api/v1/management/team/member', { teamUuid: this.selected.uuid }); this.searchOpen = false; },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="チーム管理" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"チーム"}]\' data-management-tabs="true" data-management-tab="team"'
    return head_block(rel, "チーム管理") + app_body(attrs, content, script)


def mgmt_team_create() -> str:
    rel = "management/team/create.html"
    content = """    <div x-data="teamFormPage('create')">
      <h1 class="font-display text-2xl font-semibold">チーム新規作成</h1>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm space-y-4">
        <label class="form-control"><span class="label-text">チーム名</span><input class="input input-bordered" x-model="form.name" /></label>
        <button type="button" class="btn border-rose-900 bg-rose-900 text-white" @click="save()">作成</button>
        <a href="index.html" class="btn btn-ghost">戻る</a>
      </div>
    </div>"""
    script = """  <script>
    function teamFormPage(mode) {
      const uuid = getUUID('team-001');
      const existing = window.SAMPLE_DATA.teams.find(t => t.uuid === uuid);
      return {
        form: existing ? { name: existing.name } : { name: '' },
        save() { apiAlert(mode === 'create' ? 'POST' : 'PUT', mode === 'create' ? '/api/v1/management/team' : `/api/v1/management/team/${uuid}`, this.form); },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="チーム作成" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"作成"}]\' data-management-tabs="true" data-management-tab="team"'
    return head_block(rel, "チーム作成") + app_body(attrs, content, script)


def mgmt_team_edit() -> str:
    rel = "management/team/edit.html"
    content = """    <div x-data="teamFormPage('edit')">
      <h1 class="font-display text-2xl font-semibold">チーム編集</h1>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm space-y-4">
        <label class="form-control"><span class="label-text">チーム名</span><input class="input input-bordered" x-model="form.name" /></label>
        <button type="button" class="btn border-rose-900 bg-rose-900 text-white" @click="save()">保存</button>
        <a href="index.html" class="btn btn-ghost">戻る</a>
      </div>
    </div>"""
    attrs = 'data-layout="app" data-title="チーム編集" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"編集"}]\' data-management-tabs="true" data-management-tab="team"'
    return head_block(rel, "チーム編集") + app_body(attrs, content, """  <script>
    function teamFormPage(mode) {
      const uuid = getUUID('team-001');
      const existing = window.SAMPLE_DATA.teams.find(t => t.uuid === uuid) || window.SAMPLE_DATA.teams[0];
      return { form: { name: existing.name }, save() { apiAlert('PUT', `/api/v1/management/team/${uuid}`, this.form); } };
    }
  </script>""")


def mgmt_template_index() -> str:
    rel = "management/template/index.html"
    content = """    <div x-data="templateIndexPage()">
      <div class="flex justify-between gap-4">
        <h1 class="font-display text-2xl font-semibold">テンプレート管理</h1>
        <a href="create.html" class="btn btn-sm border-rose-900 bg-rose-900 text-white">新規作成</a>
      </div>
      <div class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
        <table class="table">
          <thead><tr><th>タイトル</th><th>質問数</th><th></th></tr></thead>
          <tbody>
            <template x-for="t in templates" :key="t.uuid">
              <tr>
                <td x-text="t.title"></td><td x-text="t.topics.length"></td>
                <td class="space-x-1">
                  <button type="button" class="btn btn-xs btn-ghost" @click="openDetail(t)">詳細</button>
                  <a :href="'edit.html?uuid=' + t.uuid" class="btn btn-xs btn-outline">更新</a>
                  <button type="button" class="btn btn-xs btn-ghost" @click="openDelete(t)">削除</button>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
      <dialog class="modal" :class="detailOpen && 'modal-open'">
        <div class="modal-box rounded-3xl max-w-lg" x-show="selected">
          <h3 class="font-semibold" x-text="selected?.title"></h3>
          <ul class="mt-4 space-y-2 text-sm">
            <template x-for="topic in (selected?.topics || [])" :key="topic.uuid">
              <li x-text="topic.displayNo + '. ' + topic.question"></li>
            </template>
          </ul>
          <div class="modal-action">
            <a :href="selected ? 'edit.html?uuid=' + selected.uuid : '#'" class="btn btn-outline">更新</a>
            <button type="button" class="btn btn-ghost" @click="detailOpen=false; openDelete(selected)">削除</button>
            <button type="button" class="btn" @click="detailOpen=false">閉じる</button>
          </div>
        </div>
        <form method="dialog" class="modal-backdrop bg-black/60 backdrop-blur-md"><button type="button" @click="detailOpen=false">close</button></form>
      </dialog>
      <dialog class="modal" :class="deleteOpen && 'modal-open'">
        <div class="modal-box rounded-3xl"><h3 class="font-semibold">削除の確認</h3>
          <div class="modal-action"><button type="button" class="btn btn-ghost" @click="deleteOpen=false">キャンセル</button><button type="button" class="btn btn-error" @click="confirmDelete()">削除確定</button></div>
        </div>
        <form method="dialog" class="modal-backdrop bg-black/60 backdrop-blur-md"><button type="button" @click="deleteOpen=false">close</button></form>
      </dialog>
    </div>"""
    script = """  <script>
    function templateIndexPage() {
      return {
        templates: window.SAMPLE_DATA.templates,
        selected: null, detailOpen: false, deleteOpen: false,
        openDetail(t) { this.selected = t; this.detailOpen = true; },
        openDelete(t) { this.selected = t; this.deleteOpen = true; },
        confirmDelete() { apiAlert('DELETE', `/api/v1/management/template/${this.selected.uuid}`, null); this.deleteOpen = false; },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="テンプレート" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"テンプレート"}]\' data-management-tabs="true" data-management-tab="template"'
    return head_block(rel, "テンプレート") + app_body(attrs, content, script)


def mgmt_template_create() -> str:
    rel = "management/template/create.html"
    content = """    <div x-data="templateFormPage('create')">
      <h1 class="font-display text-2xl font-semibold">テンプレート新規作成</h1>
      <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm space-y-4">
        <label class="form-control"><span class="label-text">タイトル</span><input class="input input-bordered" x-model="form.title" /></label>
        <template x-for="(topic, i) in form.topics" :key="i">
          <div class="rounded-2xl border border-base-300 p-4">
            <input class="input input-bordered w-full" x-model="topic.question" placeholder="質問文" />
            <label class="label cursor-pointer justify-start gap-2 mt-2"><input type="checkbox" class="checkbox checkbox-sm" x-model="topic.isRequire" /><span class="label-text">必須</span></label>
          </div>
        </template>
        <button type="button" class="btn btn-outline btn-sm" @click="form.topics.push({ question: '', isRequire: false })">質問を追加</button>
        <button type="button" class="btn border-rose-900 bg-rose-900 text-white" @click="save()">保存</button>
        <a href="index.html" class="btn btn-ghost">戻る</a>
      </div>
    </div>"""
    script = """  <script>
    function templateFormPage(mode) {
      const uuid = getUUID('tpl-001');
      const existing = window.SAMPLE_DATA.templates.find(t => t.uuid === uuid);
      return {
        form: existing ? { title: existing.title, topics: existing.topics.map(t => ({ question: t.question, isRequire: t.isRequire })) } : { title: '', topics: [{ question: '', isRequire: true }] },
        save() { apiAlert(mode === 'create' ? 'POST' : 'PUT', mode === 'create' ? '/api/v1/management/template' : `/api/v1/management/template/${uuid}`, this.form); },
      };
    }
  </script>"""
    attrs = 'data-layout="app" data-title="テンプレート作成" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"作成"}]\' data-management-tabs="true" data-management-tab="template"'
    return head_block(rel, "テンプレート作成") + app_body(attrs, content, script)


def mgmt_template_edit() -> str:
    rel = "management/template/edit.html"
    content = mgmt_template_create().split('<div id="page-content">')[1].split("</div>\n  <script>")[0].replace("新規作成", "編集").replace('templateFormPage(\'create\')', "templateFormPage('edit')")
    attrs = 'data-layout="app" data-title="テンプレート編集" data-nav-active="management" data-breadcrumbs=\'[{"label":"ホーム","key":"home"},{"label":"管理","key":"management"},{"label":"編集"}]\' data-management-tabs="true" data-management-tab="template"'
    script = """  <script>
    function templateFormPage(mode) {
      const uuid = getUUID('tpl-001');
      const existing = window.SAMPLE_DATA.templates.find(t => t.uuid === uuid) || window.SAMPLE_DATA.templates[0];
      return {
        form: { title: existing.title, topics: existing.topics.map(t => ({ question: t.question, isRequire: t.isRequire })) },
        save() { apiAlert('PUT', `/api/v1/management/template/${uuid}`, this.form); },
      };
    }
  </script>"""
    return head_block(rel, "テンプレート編集") + app_body(attrs, content, script)


PAGES = [
    ("login.html", login_page),
    ("logout.html", logout_page),
    ("error.html", error_page),
    ("index.html", home_page),
    ("setting.html", setting_page),
    ("meeting/index.html", meeting_index),
    ("meeting/answer.html", meeting_answer),
    ("meeting/view.html", meeting_view),
    ("meeting/setting.html", meeting_setting),
    ("management/index.html", mgmt_hub),
    ("management/progress.html", progress_page),
    ("management/setting.html", mgmt_setting),
    ("management/meeting/index.html", mgmt_meeting_index),
    ("management/meeting/detail.html", mgmt_meeting_detail),
    ("management/meeting/log.html", mgmt_meeting_log),
    ("management/pairing/index.html", mgmt_pairing_index),
    ("management/pairing/create.html", mgmt_pairing_create),
    ("management/pairing/edit.html", mgmt_pairing_edit),
    ("management/pairing/department/index.html", mgmt_pairing_dept_index),
    ("management/pairing/department/detail.html", mgmt_pairing_dept_detail),
    ("management/pairing/department/edit.html", mgmt_pairing_dept_edit),
    ("management/department/index.html", mgmt_department_index),
    ("management/department/create.html", mgmt_department_create),
    ("management/department/detail.html", mgmt_department_detail),
    ("management/department/edit.html", mgmt_department_edit),
    ("management/team/index.html", mgmt_team_index),
    ("management/team/create.html", mgmt_team_create),
    ("management/team/edit.html", mgmt_team_edit),
    ("management/template/index.html", mgmt_template_index),
    ("management/template/create.html", mgmt_template_create),
    ("management/template/edit.html", mgmt_template_edit),
]


def main() -> None:
    print(f"Generating {len(PAGES)} pages into {OUT}")
    for rel, builder in PAGES:
        write_page(rel, builder())
    print("Done.")


if __name__ == "__main__":
    main()
