<script lang="ts">
    import type { Card, CardState } from "../types";

    const { cards }: { cards: CardState[] } = $props();
</script>

{#snippet card(card?: Card)}
    <div
        class="flex-1 flex items-center justify-center bg-neutral-800 rounded-md h-12"
    >
        {#if card}
            {card}
        {:else}
            ?
        {/if}
    </div>
{/snippet}

<div class="flex flex-1 gap-1 my-2 min-w-12">
    {#each { length: 2 } as _, i}
        {#if cards[i]}
            {@render card(cards[i].card)}
        {:else}
            {@render card()}
        {/if}
    {/each}
</div>
