/** Matches tailwind max-h-48 (12rem) used on pick panels. */
const PANEL_FLIP_HEIGHT_PX = 192;
const GAP_PX = 4;

export function portal(node: HTMLElement) {
    document.body.appendChild(node);
    return {
        destroy() {
            node.remove();
        },
    };
}

export type FloatingParams = {
    getTrigger: () => HTMLElement | null;
};

export function floating(node: HTMLElement, params: FloatingParams) {
    let current = params;

    function position() {
        const trigger = current.getTrigger();
        if (!trigger) return;

        const rect = trigger.getBoundingClientRect();
        const spaceBelow = window.innerHeight - rect.bottom - GAP_PX;
        const openUp =
            spaceBelow < PANEL_FLIP_HEIGHT_PX &&
            rect.top > spaceBelow;

        node.style.left = `${rect.left}px`;
        node.style.width = `${rect.width}px`;
        node.style.maxWidth = `${rect.width}px`;

        if (openUp) {
            node.style.top = "auto";
            node.style.bottom = `${window.innerHeight - rect.top + GAP_PX}px`;
        } else {
            node.style.bottom = "auto";
            node.style.top = `${rect.bottom + GAP_PX}px`;
        }
    }

    const onScroll = () => position();
    const onResize = () => position();

    position();
    window.addEventListener("resize", onResize);
    window.addEventListener("scroll", onScroll, true);

    return {
        update(next: FloatingParams) {
            current = next;
            position();
        },
        destroy() {
            window.removeEventListener("resize", onResize);
            window.removeEventListener("scroll", onScroll, true);
        },
    };
}