<script lang="ts">
    import {
        themeOptions,
        applyTheme,
        THEME_KEY,
        DEFAULT_THEME,
    } from "$lib/themes";
    import Palette from "lucide-svelte/icons/palette";
    import X from "lucide-svelte/icons/x";

    let dialogEl = $state<HTMLDialogElement | null>(null);
    let theme = $state(
        typeof localStorage !== "undefined"
            ? localStorage.getItem(THEME_KEY) || DEFAULT_THEME
            : DEFAULT_THEME,
    );

    function openModal() {
        dialogEl?.showModal();
    }

    function closeModal() {
        dialogEl?.close();
    }

    function pick(value: string) {
        theme = value;
        applyTheme(value);
        closeModal();
    }

    function onDialogClick(e: MouseEvent) {
        const rect = dialogEl?.getBoundingClientRect();
        if (!rect) return;
        const inDialog =
            e.clientX >= rect.left &&
            e.clientX <= rect.right &&
            e.clientY >= rect.top &&
            e.clientY <= rect.bottom;
        if (!inDialog) closeModal();
    }
</script>

<button
    type="button"
    class="btn btn-ghost btn-sm gap-1"
    onclick={openModal}
    aria-label="テーマを選択">
    <Palette class="h-4 w-4" />
</button>

<dialog
    bind:this={dialogEl}
    class="modal"
    onclick={onDialogClick}
    onkeydown={(e: KeyboardEvent) => e.key === "Escape" && closeModal()}>
    <div class="modal-box max-w-md" role="document">
        <div class="flex items-center justify-between gap-2">
            <h2 class="font-display text-lg font-semibold">テーマ</h2>
            <button
                type="button"
                class="btn btn-ghost btn-sm btn-square"
                onclick={closeModal}
                aria-label="閉じる">
                <X class="h-4 w-4" />
            </button>
        </div>
        <p class="mt-1 text-sm text-base-content/60">
            表示テーマを選んでください
        </p>
        <ul
            class="mt-4 grid max-h-[min(60vh,24rem)] gap-1 overflow-y-auto pr-1">
            {#each themeOptions as opt}
                <li>
                    <button
                        type="button"
                        class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-left text-sm transition hover:bg-base-200 {theme ===
                        opt.value
                            ? 'bg-primary/10 font-medium text-primary'
                            : ''}"
                        data-theme={opt.value}
                        onclick={() => pick(opt.value)}>
                        <span
                            class="h-8 w-8 shrink-0 rounded-lg border border-base-300 bg-primary"
                            aria-hidden="true"></span>
                        {opt.label}
                    </button>
                </li>
            {/each}
        </ul>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="submit" aria-label="閉じる">閉じる</button>
    </form>
</dialog>
