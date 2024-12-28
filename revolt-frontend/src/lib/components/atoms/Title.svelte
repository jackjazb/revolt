<script lang="ts">
    import type { Snippet } from "svelte";

    const {
        children,
    }: {
        children: Snippet;
    } = $props();
</script>

<h1 class="text-xl my-1 font-bold">
    {@render children()}
</h1>
