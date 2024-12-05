<!--
The base game view for the leader, primarily allowing action choice.
-->
<script lang="ts">
    import { ActionType, TurnState, type GameComponentContext } from "../types";

    let { gameState, client }: GameComponentContext = $props();

    // Leaders can end their turn after a timeout.
    let timeout = $state(3000);

    $effect(() => {
        if (gameState.turnState !== TurnState.ActionPending) {
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

    const income = () => {
        client.attemptAction({
            type: ActionType.Income,
        });
    };
    const commit = () => {
        client.commitTurn();
    };
    const end = () => {
        client.endTurn();
    };
</script>

{#if gameState.turnState === TurnState.Default}
    <button onclick={income}>income</button>
    <button>foreign aid</button>
    <button>tax</button>
    <button>assassinate</button>
    <button>revolt</button>
    <button>exchange</button>
    <button>steal</button>
{:else if gameState.turnState === TurnState.ActionPending}
    <h1>you've attempted {gameState.pendingAction.type}</h1>
    {#if gameState.pendingAction.target}
        <h2>targeting {gameState.pendingAction.target}!</h2>
    {/if}
    <p>end turn in {(timeout / 1000).toFixed(2)}s</p>
    <button disabled={timeout > 0} onclick={commit}>commit</button>
{:else if gameState.turnState === TurnState.BlockPending}
    <h1>block pending</h1>
{:else if gameState.turnState === TurnState.Finished}
    <h1>finished</h1>
    <button onclick={end}>end turn</button>
{/if}
