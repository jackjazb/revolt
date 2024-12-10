<!--
The base game view for the leader, primarily allowing action choice.
-->
<script lang="ts">
    import { global } from "../state.svelte";
    import { ActionType, TurnState } from "../types";
    import Button from "./atomic/Button.svelte";
    import Title from "./atomic/Title.svelte";

    // Leaders can end their turn after a timeout.
    let timeout = $state(3000);

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

    const income = () => {
        global.client.attemptAction({
            type: ActionType.Income,
        });
    };
    const commit = () => {
        global.client.commitTurn();
    };
    const end = () => {
        global.client.endTurn();
    };
</script>

{#if global.state.turnState === TurnState.Default}
    <Title>actions</Title>
    <div class="flex flex-col gap-2">
        <Button onclick={income}>income</Button>
        <Button>foreign aid</Button>
        <Button>tax</Button>
        <Button>assassinate</Button>
        <Button>revolt</Button>
        <Button>exchange</Button>
        <Button>steal</Button>
    </div>
{:else if global.state.turnState === TurnState.ActionPending && global.state.pendingAction}
    <h1>you've attempted {global.state.pendingAction.type}</h1>
    {#if global.state.pendingAction.target}
        <h2>targeting {global.state.pendingAction.target}!</h2>
    {/if}
    <p>end turn in {(timeout / 1000).toFixed(2)}s</p>
    <Button disabled={timeout > 0} onclick={commit}>commit</Button>
{:else if global.state.turnState === TurnState.BlockPending}
    <h1>block pending</h1>
{:else if global.state.turnState === TurnState.Finished}
    <h1>finished</h1>
    <Button onclick={end}>end turn</Button>
{/if}
