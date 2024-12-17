<script lang="ts">
    import { global } from "../state.svelte";
    import { getPlayerById } from "../utils";
    import Button from "./atoms/Button.svelte";
    import Subtitle from "./atoms/Subtitle.svelte";

    const resolveDeath = (card: number) => {
        global.client.resolveDeath(card);
    };
</script>

{#if global.state.nextDeath === global.state.self.id}
    <Subtitle>you lose a card!</Subtitle>
    {#each global.state.self.cards as card, i}
        {#if card.alive}
            <Button onclick={() => resolveDeath(i)}>kill {card.card}</Button>
        {/if}
    {/each}
{:else}
    <Subtitle>
        {getPlayerById(global.state, global.state.nextDeath)} loses a card!
    </Subtitle>
{/if}
