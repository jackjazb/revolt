<script lang="ts">
    import { global } from "../state.svelte";
    import { TurnState } from "../types";
    import { getCurrentActionBlockers, getPlayerById, stateIn } from "../utils";
    import ActionAssassinate from "./actions/ActionAssassinate.svelte";
    import ActionBlock from "./actions/ActionBlock.svelte";
    import ActionChallenge from "./actions/ActionChallenge.svelte";
    import ActionCommit from "./actions/ActionCommit.svelte";
    import ActionForeignAid from "./actions/ActionForeignAid.svelte";
    import ActionIncome from "./actions/ActionIncome.svelte";
    import ActionRevolt from "./actions/ActionRevolt.svelte";
    import ActionSteal from "./actions/ActionSteal.svelte";
    import ActionTax from "./actions/ActionTax.svelte";
    import Button from "./atoms/Button.svelte";
    import Subtitle from "./atoms/Subtitle.svelte";
    import Title from "./atoms/Title.svelte";
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
    {#if global.state.self.leading}
        <div class="flex flex-col gap-2">
            <ActionIncome />
            <ActionForeignAid />
            <ActionTax />
            <Button disabled>Exchange</Button>
        </div>
        <div class={`grid grid-cols-5 gap-2`}>
            {#each global.state.peers as peer}
                <div class="flex flex-col gap-2">
                    <Subtitle>{peer.name}</Subtitle>

                    <ActionRevolt target={peer.id} />
                    <ActionAssassinate target={peer.id} />
                    <ActionSteal target={peer.id} />
                </div>
            {/each}
        </div>
    {:else}
        <Subtitle>Waiting for the leader to play.</Subtitle>
    {/if}
{:else if stateIn(global.state, TurnState.ActionPending) && global.state.pendingAction}
    <!-- The leader can end their turn after a timeout. Peers can block or challenge. -->
    {#if global.state.self.leading}
        <Button disabled={timeout > 0} onclick={commit}>
            Finish ({(timeout / 1000).toFixed(2)}s)
        </Button>
    {:else}
        {#each getCurrentActionBlockers(global.state) as card}
            <ActionBlock {card} />
        {/each}
        <ActionChallenge />
    {/if}
{:else if stateIn(global.state, TurnState.BlockPending)}
    {#if global.state.self.leading}
        <ActionChallenge />
        <ActionCommit text="Accept" />
    {:else}
        The leader's move has been blocked
    {/if}
{:else if stateIn(global.state, TurnState.PlayerKilled, TurnState.LeaderLostChallenge, TurnState.PlayerLostChallenge)}
    <ResolveDeath />
{:else if global.state.turnState === TurnState.Finished}
    {#if global.state.self.leading}
        <h1>Finished</h1>
        <Button onclick={end}>End your turn</Button>
    {:else}
        <Subtitle>Waiting for the leader to end the turn</Subtitle>
    {/if}
{:else if global.state.turnState === TurnState.PlayerWon}
    <Title>{getPlayerById(global.state, global.state.winner)} won!</Title>
{/if}
