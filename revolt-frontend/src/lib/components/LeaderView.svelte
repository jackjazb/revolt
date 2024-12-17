<!--
The base game view for the leader, primarily allowing action choice.
-->
<script lang="ts">
    import { global } from "../state.svelte";
    import { TurnState } from "../types";
    import { getPlayerById, stateIn } from "../utils";
    import ActionAssassinate from "./actions/ActionAssassinate.svelte";
    import ActionChallenge from "./actions/ActionChallenge.svelte";
    import ActionCommit from "./actions/ActionCommit.svelte";
    import ActionForeignAid from "./actions/ActionForeignAid.svelte";
    import ActionIncome from "./actions/ActionIncome.svelte";
    import ActionRevolt from "./actions/ActionRevolt.svelte";
    import ActionSteal from "./actions/ActionSteal.svelte";
    import ActionTax from "./actions/ActionTax.svelte";
    import Button from "./atoms/Button.svelte";
    import Subtitle from "./atoms/Subtitle.svelte";
    import ResolveDeath from "./ResolveDeath.svelte";

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
    <div class="flex flex-col gap-2">
        <ActionIncome />
        <ActionForeignAid />
        <ActionTax />
        <Button disabled>exchange</Button>
    </div>
    <div class={`grid grid-cols-5`}>
        {#each global.state.peers as peer}
            <div class="flex flex-col gap-2">
                <Subtitle>{peer.name}</Subtitle>

                <ActionAssassinate target={peer.id} />
                <ActionRevolt target={peer.id} />
                <ActionSteal target={peer.id} />
            </div>
        {/each}
    </div>
{:else if stateIn(global.state, TurnState.ActionPending) && global.state.pendingAction}
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
{:else if stateIn(global.state, TurnState.BlockPending)}
    <h1>block pending</h1>
    <ActionChallenge />
    <ActionCommit />
{:else if stateIn(global.state, TurnState.PlayerKilled, TurnState.LeaderLostChallenge, TurnState.PlayerLostChallenge)}
    <ResolveDeath />
{:else if global.state.turnState === TurnState.Finished}
    {#if global.state.self.leading}
        <h1>finished</h1>
        <Button onclick={end}>end turn</Button>
    {:else}
        <Subtitle>Waiting for leader to end turn</Subtitle>
    {/if}
{/if}
