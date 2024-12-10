<script lang="ts">
    import type { Snippet } from "svelte";

    const {
        children,
        ...rest
    }: {
        children: Snippet;
        class?: string;
    } = $props();
</script>

<div class={`rounded-sm bg-neutral-700 p-2 flex flex-col gap-2 ${rest.class}`}>
    {@render children()}
</div>
