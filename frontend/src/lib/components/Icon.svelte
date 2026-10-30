<script lang="ts">
    let {
        width,
        height,
        base64,
        url,
        text,
        alt = "",
    }: {
        width: number | string;
        height: number | string;
        base64?: string | null;
        url?: string | null;
        text?: string | null;
        alt?: string;
    } = $props();

    let src = $derived.by(() => {
        if (base64) {
            return base64.startsWith("data:")
                ? base64
                : `data:image/png;base64,${base64}`;
        }
        if (url) {
            return url;
        }
        return null;
    });

    let widthCss = $derived(typeof width === "number" ? `${width}px` : width);
    let heightCss = $derived(
        typeof height === "number" ? `${height}px` : height,
    );

    let initial = $derived(
        !src && text !== undefined ? (text || "?").slice(0, 1) : null,
    );
</script>

{#if src}
    <img
        {src}
        {alt}
        class="rounded-full object-cover"
        style:width={widthCss}
        style:height={heightCss} />
{:else if initial}
    <span
        class="grid place-items-center rounded-full bg-primary text-xs font-semibold text-primary-content"
        style:width={widthCss}
        style:height={heightCss}
        aria-hidden="true">{initial}</span
    >
{/if}
