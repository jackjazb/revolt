<script>
    import { global } from "../state.svelte";
    import { TurnState } from "../types";
    import {
        formatActionType,
        formatCard,
        getLeader,
        getPlayerById,
        stateIn,
    } from "../utils";
    import Title from "./atoms/Title.svelte";

    // TODO can base this on pendingChallenge, not the pending block state
    let leaderFailedToChallengeBlock = $derived(
        stateIn(global.state, TurnState.LeaderLostChallenge) &&
            global.state.pendingBlock.initiator !== "",
    );

    let leaderSuccessfullyChallengedBlock = $derived(
        stateIn(global.state, TurnState.PlayerLostChallenge) &&
            global.state.pendingBlock.initiator !== "",
    );

    let playerSuccessfullyChallengedAction = $derived(
        stateIn(global.state, TurnState.LeaderLostChallenge) &&
            global.state.pendingBlock.initiator === "",
    );

    let playerFailedToChallengeAction = $derived(
        stateIn(global.state, TurnState.PlayerLostChallenge) &&
            global.state.pendingBlock.initiator === "",
    );
</script>

<Title>
    {#if stateIn(global.state, TurnState.Default)}
        {getLeader(global.state)} is choosing an action.
    {:else if stateIn(global.state, TurnState.ActionPending)}
        {getLeader(global.state)}
        has attempted
        {formatActionType(global.state.pendingAction.type)}
        {#if global.state.pendingAction.target}
            against
            {getPlayerById(global.state, global.state.pendingAction.target)}
        {/if}
        .
    {:else if stateIn(global.state, TurnState.Finished)}
        The turn is over.
    {:else if stateIn(global.state, TurnState.BlockPending)}
        <!-- 'Jack has blocked Bob's Assassinate using their Contessa'-->
        {getPlayerById(global.state, global.state.pendingBlock.initiator ?? "")}
        has blocked
        {getLeader(global.state)}'s
        {formatActionType(global.state.pendingAction.type)}
        using their
        {formatCard(global.state.pendingBlock.card)}
    {:else if leaderFailedToChallengeBlock}
        The leader failed to challenge a block.
    {:else if leaderSuccessfullyChallengedBlock}
        The leader successfully challenged the block.
    {:else if playerSuccessfullyChallengedAction}
        {getPlayerById(global.state, global.state.pendingChallenge.initiator)}
        successfully challenged the action.
    {:else if playerFailedToChallengeAction}
        {getPlayerById(global.state, global.state.pendingChallenge.initiator)}
        failed to challenge the action.
    {:else if stateIn(global.state, TurnState.PlayerWon)}
        {getPlayerById(global.state, global.state.winner)}
        won!
    {:else if stateIn(global.state, TurnState.PlayerKilled)}
        {getPlayerById(global.state, global.state.nextDeath)}
        loses a card.
    {/if}
</Title>
