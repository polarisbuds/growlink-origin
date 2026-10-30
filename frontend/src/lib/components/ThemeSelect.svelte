<script lang="ts">
    import {
        themeOptions,
        applyTheme,
        THEME_KEY,
        DEFAULT_THEME,
    } from "$lib/themes";
    import ChevronDown from "lucide-svelte/icons/chevron-down";

    let open = $state(false);
    let theme = $state(
        typeof localStorage !== "undefined"
            ? localStorage.getItem(THEME_KEY) || DEFAULT_THEME
            : DEFAULT_THEME,
    );

    function openModal() {
        open = true;
    }

    function closeModal() {
        open = false;
    }

    function pick(value: string) {
        theme = value;
        applyTheme(value);
    }

    function onKeydown(e: KeyboardEvent) {
        if (e.key === "Escape" && open) closeModal();
    }
</script>

<svelte:window onkeydown={onKeydown} />

{#snippet themeSwatch()}
    <svg
        class="h-5 w-5 shrink-0"
        viewBox="0 0 20 20"
        xmlns="http://www.w3.org/2000/svg"
        aria-hidden="true">
        <rect
            x="1"
            y="1"
            width="18"
            height="18"
            rx="4"
            class="fill-base-100 stroke-base-content/30"
            stroke-width="1" />
        <circle cx="6.5" cy="6.5" r="2.25" class="fill-primary" />
        <circle cx="13.5" cy="6.5" r="2.25" class="fill-neutral" />
        <circle cx="6.5" cy="13.5" r="2.25" class="fill-secondary" />
        <circle cx="13.5" cy="13.5" r="2.25" class="fill-accent" />
    </svg>
{/snippet}

<button
    type="button"
    class="btn btn-ghost btn-sm gap-1 border border-base-300"
    onclick={openModal}
    aria-label="テーマを選択"
    aria-expanded={open}>
    {@render themeSwatch()}
    <ChevronDown class="h-4 w-4 opacity-60" />
</button>

<div
    class="modal backdrop-blur-sm"
    class:modal-open={open}
    role="dialog"
    aria-modal="true"
    aria-labelledby="theme-modal-title">
    <div
        class="modal-box relative flex max-h-[85vh] max-w-5xl flex-col rounded-2xl bg-base-100 shadow-2xl">
        <button
            type="button"
            class="btn btn-sm btn-circle btn-ghost absolute right-3 top-3"
            onclick={closeModal}
            aria-label="閉じる">✕</button>
        <h3 id="theme-modal-title" class="pr-10 text-lg font-bold">
            テーマのリスト
        </h3>
        <div
            class="mt-4 grid min-h-0 flex-1 grid-cols-2 gap-3 overflow-y-auto px-1 py-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6">
            {#each themeOptions as opt}
                <button
                    type="button"
                    data-theme={opt.value}
                    class="flex flex-col gap-2 rounded-box border border-base-300 bg-base-100 p-3 text-left text-base-content transition hover:border-base-content/20 {theme ===
                    opt.value
                        ? 'outline outline-2 outline-offset-2 outline-base-content'
                        : ''}"
                    onclick={() => pick(opt.value)}
                    aria-label={opt.label}>
                    <span class="text-sm font-medium capitalize"
                        >{opt.label}</span>
                    {@render themeSwatch()}
                </button>
            {/each}
        </div>
        <div class="modal-action mt-4 shrink-0">
            <button type="button" class="btn" onclick={closeModal}
                >閉じる</button>
        </div>
    </div>
    <button
        type="button"
        class="modal-backdrop bg-opacity-100"
        onclick={closeModal}
        aria-label="閉じる"></button>
</div>
