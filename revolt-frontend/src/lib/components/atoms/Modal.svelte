<script lang="ts">
    import type { Snippet } from "svelte";
    import { fade, scale } from "svelte/transition";

    const { open, children }: { open: boolean; children: Snippet } = $props();
</script>

{#if open}
    <div
        class="fixed top-0 left-0 w-full h-full z-20 flex items-center justify-center text-inherit bg-inherit"
    >
        <div
            transition:fade
            class="w-full h-full
            fixed
            bg-black
            opacity-30
     text-inherit bg-inherit"
        ></div>
        <div class="panel z-30" transition:scale={{ duration: 200 }}>
            {@render children()}
        </div>
    </div>
{/if}
