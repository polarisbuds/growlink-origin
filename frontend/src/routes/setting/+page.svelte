<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import Camera from "lucide-svelte/icons/camera";
    import Trash2 from "lucide-svelte/icons/trash-2";

    type Profile = {
        name?: string;
        email?: string;
        lineWorksId?: string;
        imageUri?: string | null;
        remindEnabled?: boolean;
    };

    function apiErrorMessage(e: unknown, fallback: string): string {
        if (e instanceof Error) {
            try {
                const j = JSON.parse(e.message) as { error?: string };
                if (j?.error) return j.error;
            } catch {
                /* not json */
            }
            if (e.message && e.message !== "unauthorized") return e.message;
        }
        return fallback;
    }

    let form = $state({
        name: "",
        email: "",
        lineWorksId: "",
        remindEnabled: false,
    });
    let imageUri = $state<string | null>(null);
    let fieldError = $state("");
    let iconLoading = $state(false);

    let avatarUrl = $derived(imageUri ? `/uploads/${imageUri}` : null);
    let initial = $derived((form.name || "?").slice(0, 1));
    let hasIcon = $derived(!!avatarUrl);

    onMount(async () => {
        const profile = await api<Profile>("GET", "/api/v1/setting/profile");
        form = {
            name: profile.name || "",
            email: profile.email || "",
            lineWorksId: profile.lineWorksId || "",
            remindEnabled: !!profile.remindEnabled,
        };
        imageUri = profile.imageUri ?? null;
    });

    async function saveProfile() {
        fieldError = "";
        try {
            await api("PUT", "/api/v1/setting", form);
            showToast("保存しました");
        } catch (e) {
            fieldError = apiErrorMessage(e, "保存に失敗しました");
            showToast(fieldError, "error");
        }
    }

    async function uploadIcon(e: Event) {
        const input = e.currentTarget as HTMLInputElement;
        const file = input.files?.[0];
        if (!file) return;

        if (!file.type.startsWith("image/")) {
            showToast("画像ファイルを選択してください", "error");
            input.value = "";
            return;
        }

        iconLoading = true;
        const fd = new FormData();
        fd.append("icon", file);
        try {
            await api("POST", "/api/v1/setting/icon", fd);
            const data = await api<Profile>("GET", "/api/v1/setting/profile");
            imageUri = data.imageUri ?? null;
            showToast("アイコンをアップロードしました");
        } catch {
            showToast("アップロードに失敗しました", "error");
        } finally {
            iconLoading = false;
            input.value = "";
        }
    }

    async function deleteIcon() {
        if (iconLoading) return;
        iconLoading = true;
        try {
            await api("DELETE", "/api/v1/setting/icon");
            imageUri = null;
            showToast("アイコンを削除しました");
        } catch {
            showToast("削除に失敗しました", "error");
        } finally {
            iconLoading = false;
        }
    }
</script>

<svelte:head><title>設定 - Growlink</title></svelte:head>

<div class="max-w-5xl w-full">
    <h1 class="font-display text-2xl font-semibold">個人設定</h1>
    <p class="text-sm text-base-content/60 mt-1">プロフィールと通知の設定</p>
    <div class="grid grid-cols-3 gap-4">
        <div
            class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm space-y-5 col-span-2">
            <h2 class="font-semibold">基本情報</h2>

            <div class="flex justify-around gap-2">
                <div class="flex flex-col items-center gap-4 justify-around">
                    <label
                        class="group relative block w-44 cursor-pointer focus-within:outline-none focus-within:ring-2 focus-within:ring-primary/40 focus-within:ring-offset-2 focus-within:ring-offset-base-100"
                        class:pointer-events-none={iconLoading}>
                        <input
                            type="file"
                            class="sr-only"
                            accept="image/*"
                            disabled={iconLoading}
                            onchange={uploadIcon} />
                        <div
                            class="mask mask-squircle relative aspect-square w-full overflow-hidden bg-base-200 border-2 border-base-300 shadow-sm">
                            {#if avatarUrl}
                                <img
                                    src={avatarUrl}
                                    alt=""
                                    class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105" />
                            {:else}
                                <div
                                    class="absolute inset-0 grid place-items-center bg-primary text-4xl font-semibold text-primary-content">
                                    {initial}
                                </div>
                            {/if}

                            <div
                                class="absolute inset-0 flex flex-col items-center justify-center gap-2 bg-black/50 text-white opacity-0 backdrop-blur-[2px] transition-opacity duration-200 group-hover:opacity-100 group-focus-within:opacity-100">
                                <Camera class="h-6 w-6" />
                                <span class="text-sm font-medium"
                                    >アイコンを変更</span>
                            </div>

                            {#if iconLoading}
                                <div
                                    class="absolute inset-0 flex items-center justify-center bg-base-100/80 backdrop-blur-sm">
                                    <span
                                        class="loading loading-spinner loading-md text-primary"
                                    ></span>
                                </div>
                            {/if}
                        </div>
                    </label>

                    {#if hasIcon}
                        <button
                            type="button"
                            class="btn btn-ghost btn-sm text-error hover:bg-error/10"
                            disabled={iconLoading}
                            onclick={deleteIcon}>
                            <Trash2 class="h-4 w-4" />
                            アイコンを削除
                        </button>
                    {/if}
                </div>
                <div class="flex flex-col items-start gap-4">
                    <label
                        class="form-control tooltip w-full flex flex-col"
                        data-tip="登録名の変更は総務人事部へ連絡してください">
                        <span class="text-xs text-base-content/50 self-start"
                            >氏名</span>
                        <input
                            class="input input-bordered mt-1"
                            bind:value={form.name}
                            disabled />
                    </label>

                    <label class="form-control w-full">
                        <span class="text-xs text-base-content/50">メール</span>
                        <input
                            class="input input-bordered mt-1"
                            bind:value={form.email}
                            type="email"
                            disabled />
                    </label>

                    <label class="form-control w-full">
                        <span class="text-xs text-base-content/50"
                            >LINE WORKS ID</span>
                        <div class="join w-full mt-1">
                            <input
                                class="input input-bordered join-item w-full placeholder-current"
                                class:input-error={!!fieldError}
                                bind:value={form.lineWorksId}
                                placeholder="ローカル部のみ" />
                            <span
                                class="btn join-item no-animation pointer-events-none bg-base-200"
                                >@cap</span>
                        </div>
                        <span class="label-text-alt text-base-content/50 mt-1"
                            >ドメイン @cap は固定です。変更できません。</span>
                        {#if fieldError}
                            <span class="label-text-alt text-error mt-1"
                                >{fieldError}</span>
                        {/if}
                    </label>
                </div>
            </div>
            <div class="flex justify-end">
                <button
                    type="button"
                    class="btn btn-primary"
                    onclick={saveProfile}>保存</button>
            </div>
        </div>

        <div
            class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
            <h2 class="font-semibold">通知</h2>
            <label
                class="label cursor-pointer justify-start gap-3 mt-3 rounded-2xl border border-base-300 px-4 py-3">
                <input
                    type="checkbox"
                    class="toggle toggle-primary"
                    bind:checked={form.remindEnabled} />
                <span class="label-text"
                    >トピック未提出のリマインドを受け取る</span>
            </label>
            <p class="text-xs text-base-content/50 mt-3">
                送信処理はバックエンドのジョブで実行されます（設定の保存のみ）
            </p>
            <div class="flex flex-1 justify-end">
                <button
                    type="button"
                    class="btn btn-outline btn-sm mt-4"
                    onclick={saveProfile}>通知設定を保存</button>
            </div>
        </div>
    </div>
</div>
