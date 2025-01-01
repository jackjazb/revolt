<script lang="ts">
    import type { Snippet } from "svelte";
    import Icon from "../atoms/Icon.svelte";

    let {
        onclick,
        children,
        lie,
    }: {
        onclick: () => void;
        children: Snippet;
        lie: boolean;
    } = $props();
</script>

<button
    {onclick}
    class="flex flex-row gap-2 items-center content justify-center"
>
    {@render children()}

    {#if lie}
        <span class="text-red-500">
            <Icon type="lie" size={16} />
        </span>
    {/if}
</button>
