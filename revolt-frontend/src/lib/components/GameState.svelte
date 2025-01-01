<script lang="ts">
    import { global } from "../state.svelte";
    import { TurnState } from "../types";
    import { getCurrentActionBlockers, getPlayerById, stateIn } from "../utils";
    import ActionBlock from "./actions/ActionBlock.svelte";

    // Leaders can end their turn after a timeout.
    let timeout = $state(0);

    $effect(() => {
        if (global.state.turnState !== TurnState.ActionPending) {
            return;
        }
        let last_time = performance.now();

        let frame = requestAnimationFrame(function update(time) {
            frame = requestAnimationFrame(update);
            let delta = time - last_time;
            timeout = Math.max(timeout - delta, 0);
            last_time = time;
        });

        // Called when component unmounts
        return () => {
            cancelAnimationFrame(frame);
        };
    });

    const commit = () => {
        global.client.commitTurn();
    };
    const end = () => {
        global.client.endTurn();
    };
</script>

{#if stateIn(global.state, TurnState.Default)}
    {#if !global.state.self.leading}
        <h2>Waiting for the leader to play.</h2>
    {/if}
{:else if stateIn(global.state, TurnState.ActionPending) && global.state.pendingAction}
    <!-- The leader can end their turn after a timeout. Peers can block or challenge. -->
    {#if global.state.self.leading}
        <button disabled={timeout > 0} onclick={commit}>
            Finish ({(timeout / 1000).toFixed(2)}s)
        </button>
    {:else}
        {#each getCurrentActionBlockers(global.state) as card}
            <ActionBlock {card} />
        {/each}
    {/if}
{:else if global.state.turnState === TurnState.Finished}
    {#if global.state.self.leading}
        <h1>Finished</h1>
        <button onclick={end}>End your turn</button>
    {:else}
        <h2>Waiting for the leader to end the turn</h2>
    {/if}
{:else if global.state.turnState === TurnState.PlayerWon}
    <h1>{getPlayerById(global.state, global.state.winner)} won!</h1>
{/if}
