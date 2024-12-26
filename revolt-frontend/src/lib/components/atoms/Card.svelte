<script lang="ts">
    import type { Snippet } from "svelte";

    const {
        children,
        row,
        ...rest
    }: {
        children: Snippet;
        class?: string;
        row?: boolean;
    } = $props();
</script>

<div
    class={`rounded-sm bg-neutral-700 p-2 gap-2 flex ${row ? "flex-row items-center" : "flex-col"} ${rest.class}`}
>
    {@render children()}
</div>
