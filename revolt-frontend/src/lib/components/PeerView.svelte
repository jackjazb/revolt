<!--
The base game view for the leader, primarily allowing action choice.
-->
<script lang="ts">
    import { global } from "../state.svelte";
    import type { Block } from "../types";
    import { Card, TurnState } from "../types";
    import { getCurrentActionBlockers, getPlayerById, stateIn } from "../utils";
    import ActionBlock from "./actions/ActionBlock.svelte";
    import ActionChallenge from "./actions/ActionChallenge.svelte";
    import Subtitle from "./atoms/Subtitle.svelte";
    import ResolveDeath from "./ResolveDeath.svelte";

    const resolveDeath = (card: number) => {
        global.client.resolveDeath(card);
    };

    const block = () => {
        const block: Block = { card: Card.Duke };
        global.client.attemptBlock(block);
    };
</script>

<div class="flex flex-col gap-2">
    {#if global.state.turnState === TurnState.Default}
        <Subtitle>Waiting for leader to play.</Subtitle>
    {:else if stateIn(global.state, TurnState.ActionPending, TurnState.BlockPending) && global.state.pendingAction}
        <h1>leader attempted {global.state.pendingAction.type}</h1>
        {#if global.state.pendingAction.target}
            <h2>
                targeting {getPlayerById(
                    global.state,
                    global.state.pendingAction.target,
                )}!
            </h2>
        {/if}

        {#each getCurrentActionBlockers(global.state) as card}
            <ActionBlock {card} />
        {/each}
        <ActionChallenge />
    {:else if stateIn(global.state, TurnState.PlayerKilled, TurnState.LeaderLostChallenge, TurnState.PlayerLostChallenge)}
        <ResolveDeath />
    {:else if global.state.turnState === TurnState.Finished}
        <Subtitle>Waiting for leader to end turn</Subtitle>
    {/if}
</div>
