<script lang="ts">
    import type { AvatarMe } from "$lib/types";
    import { buildNavigation, settingsNav } from "$lib/navigation";
    import { ui } from "$lib/stores/ui.svelte";
    import Logo from "./Logo.svelte";
    import Home from "lucide-svelte/icons/home";
    import CalendarDays from "lucide-svelte/icons/calendar-days";
    import Shield from "lucide-svelte/icons/shield";
    import BarChart2 from "lucide-svelte/icons/bar-chart-2";
    import Settings from "lucide-svelte/icons/settings";
    import { goto } from "$app/navigation";
    import { fade } from "svelte/transition";

    let { me, pathname }: { me: AvatarMe; pathname: string } = $props();

    const icons: Record<string, typeof Home> = {
        home: Home,
        "calendar-days": CalendarDays,
        shield: Shield,
        "bar-chart-2": BarChart2,
        settings: Settings,
    };

    let navigation = $derived(buildNavigation(pathname, me));
    let navSettings = $derived(settingsNav(pathname));

    function navLinkClass(active: boolean) {
        return active
            ? "bg-primary/10 text-primary"
            : "text-base-content/60 hover:bg-base-200 hover:text-base-content";
    }

    function onNavClick() {
        ui.closeMobileDrawer();
    }
</script>

{#snippet sidebarBody(expanded: boolean)}
    <div
        class="flex h-16 shrink-0 items-center gap-3 border-b border-base-300 px-5">
        <div class="grid h-8 w-8 place-items-center rounded-xl">
            <Logo />
        </div>
        {#if expanded}
            <div>
                <div
                    class="font-display text-sm font-semibold whitespace-nowrap"
                    transition:fade>
                    Growlink
                </div>
                <div
                    class="mt-0.5 text-[10px] text-base-content/50 whitespace-nowrap"
                    transition:fade>
                    1on1 ミーティング
                </div>
            </div>
        {/if}
    </div>
    <nav class="flex-1 min-h-0 space-y-1 overflow-y-auto p-3">
        {#each navigation as item}
            {@const Icon = icons[item.icon] ?? Home}
            <button
                class="flex h-11 w-full items-center gap-3 rounded-xl px-3 text-sm font-medium transition whitespace-nowrap {navLinkClass(
                    item.active,
                )}"
                onclick={() => {
                    onNavClick();
                    goto(item.href);
                }}
                transition:fade>
                <Icon class="h-5 w-5 shrink-0" />
                {#if expanded}
                    <span class="whitespace-nowrap" transition:fade
                        >{item.label}</span>
                {/if}
            </button>
        {/each}
    </nav>
    <div class="mt-auto shrink-0 border-t border-base-300 p-3">
        <button
            class="flex h-11 w-full items-center gap-3 rounded-xl px-3 text-sm font-medium transition {navLinkClass(
                navSettings.active,
            )}"
            onclick={() => {
                onNavClick();
                goto(navSettings.href);
            }}>
            <Settings class="h-5 w-5 shrink-0" />
            {#if expanded}
                <span transition:fade class="whitespace-nowrap"
                    >{navSettings.label}</span>
            {/if}
        </button>
    </div>
{/snippet}

{#if ui.mobileDrawerOpen}
    <div class="fixed inset-0 z-40 md:hidden" aria-hidden={false}>
        <button
            type="button"
            class="absolute inset-0 bg-black/40"
            aria-label="メニューを閉じる"
            onclick={() => ui.closeMobileDrawer()}></button>
        <aside
            class="relative z-50 flex h-full w-64 max-w-[85vw] flex-col border-r border-base-300 bg-base-100 shadow-xl">
            {@render sidebarBody(true)}
        </aside>
    </div>
{/if}

<aside
    class="hidden shrink-0 border-r border-base-300 bg-base-100 transition-[width] duration-200 md:flex md:flex-col {ui.sidebarOpen
        ? 'w-64'
        : 'w-20'}">
    {@render sidebarBody(ui.sidebarOpen)}
</aside>
