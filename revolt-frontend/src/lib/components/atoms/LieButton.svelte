<script lang="ts">
    import type { Snippet } from "svelte";
    import Button from "../atoms/Button.svelte";
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

<Button
    {onclick}
    class="flex flex-row gap-2 items-center content justify-center"
>
    {@render children()}

    {#if lie}
        <Icon type="lie" size={16} />
    {/if}
</Button>
