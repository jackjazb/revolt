<!--
The base game view for the leader, primarily allowing action choice.
-->
<script lang="ts">
    import { global } from "../state.svelte";
    import { TurnState } from "../types";
    import { getPlayerById } from "../utils";
    import Button from "./atomic/Button.svelte";
    import Subtitle from "./atomic/Subtitle.svelte";

    const resolveDeath = (card: number) => {
        global.client.resolveDeath(card);
    };
</script>

<div class="flex flex-col gap-2">
    {#if global.state.turnState === TurnState.Default}
        <Subtitle>Waiting for leader to play.</Subtitle>
    {:else if global.state.turnState === TurnState.ActionPending && global.state.pendingAction}
        <h1>leader attempted {global.state.pendingAction.type}</h1>
        {#if global.state.pendingAction.target}
            <h2>
                targeting {getPlayerById(
                    global.state,
                    global.state.pendingAction.target,
                )}!
            </h2>
        {/if}

        <Button>block</Button>
        <Button>challenge</Button>
    {:else if global.state.turnState === TurnState.PlayerKilled}
        {#if global.state.nextDeath === global.state.self.id}
            <Subtitle>you lose a card!</Subtitle>
            {#each global.state.self.cards as card, i}
                <Button onclick={() => resolveDeath(i)}>kill {card.card}</Button
                >
            {/each}
        {:else}
            <Subtitle
                >{getPlayerById(global.state, global.state.nextDeath)} loses a card!</Subtitle
            >
        {/if}
    {:else if global.state.turnState === TurnState.Finished}
        <Subtitle>Waiting for leader to end turn</Subtitle>
    {/if}
</div>
