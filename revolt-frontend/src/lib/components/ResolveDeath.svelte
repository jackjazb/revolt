<script lang="ts">
    import { global } from "../state.svelte";
    import { getPlayerById } from "../utils";
    import GameCard from "./atoms/GameCard.svelte";
    import Subtitle from "./atoms/Subtitle.svelte";

    const resolveDeath = (card: number) => {
        global.client.resolveDeath(card);
    };
</script>

{#if global.state.nextDeath === global.state.self.id}
    <Subtitle>Select a card to kill:</Subtitle>
    <div class="flex gap-2">
        {#each global.state.self.cards as card, i}
            {#if card.alive}
                <button onclick={() => resolveDeath(i)}>
                    <GameCard {card} />
                </button>
            {/if}
        {/each}
    </div>
{:else}
    <Subtitle>
        {getPlayerById(global.state, global.state.nextDeath)} loses a card!
    </Subtitle>
{/if}
