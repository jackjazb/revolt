<!--
The base game view for the leader, primarily allowing action choice.
-->
<script lang="ts">
    import { global } from "../state.svelte";
    import { ActionType, TurnState } from "../types";
    import { getPlayerById } from "../utils";
    import Button from "./atomic/Button.svelte";
    import Subtitle from "./atomic/Subtitle.svelte";

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

    const income = () => {
        global.client.attemptAction({
            type: ActionType.Income,
        });
    };
    const foreignAid = () => {
        global.client.attemptAction({
            type: ActionType.ForeignAid,
        });
    };
    const tax = () => {
        global.client.attemptAction({
            type: ActionType.Tax,
        });
    };
    const steal = (target: string) => {
        global.client.attemptAction({
            type: ActionType.Steal,
            target,
        });
    };
    const assassinate = (target: string) => {
        global.client.attemptAction({
            type: ActionType.Assassinate,
            target,
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
    <div class="flex flex-col gap-2">
        <Button onclick={income}>income</Button>
        <Button onclick={foreignAid}>foreign aid</Button>
        <Button onclick={tax}>tax</Button>
        <Button disabled>exchange</Button>
    </div>
    <div class={`grid grid-cols-5`}>
        {#each global.state.peers as peer}
            <div class="flex flex-col gap-2">
                <Subtitle>{peer.name}</Subtitle>
                <Button onclick={() => assassinate(peer.id)}>assassinate</Button
                >
                <Button>revolt</Button>
                <Button onclick={() => steal(peer.id)}>steal</Button>
            </div>
        {/each}
    </div>
{:else if global.state.turnState === TurnState.ActionPending && global.state.pendingAction}
    <Subtitle>you've attempted {global.state.pendingAction.type}</Subtitle>
    {#if global.state.pendingAction.target}
        <h2>
            targeting {getPlayerById(
                global.state,
                global.state.pendingAction.target,
            )}!
        </h2>
    {/if}
    <p>end turn in {(timeout / 1000).toFixed(2)}s</p>
    <Button disabled={timeout > 0} onclick={commit}>commit</Button>
{:else if global.state.turnState === TurnState.BlockPending}
    <h1>block pending</h1>
{:else if global.state.turnState === TurnState.PlayerKilled}
    <Subtitle
        >{getPlayerById(global.state, global.state.nextDeath)} loses a card!</Subtitle
    >
{:else if global.state.turnState === TurnState.Finished}
    <h1>finished</h1>
    <Button onclick={end}>end turn</Button>
{/if}
