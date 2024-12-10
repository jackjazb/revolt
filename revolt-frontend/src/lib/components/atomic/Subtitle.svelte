<script lang="ts">
    import type { Snippet } from "svelte";

    const {
        children,
    }: {
        children: Snippet;
    } = $props();
</script>

<h1 class="text-lg font-bold">
    {@render children()}
</h1>
