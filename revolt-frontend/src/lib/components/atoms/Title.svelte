<script lang="ts">
    import type { Snippet } from "svelte";

    const {
        children,
    }: {
        children: Snippet;
    } = $props();
</script>

<h1 class="text-xl my-2 font-bold">
    {@render children()}
</h1>
