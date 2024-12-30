<script lang="ts">
    import { global } from "../state.svelte";
    import { ActionType, TurnState, type Peer } from "../types";
    import { formatCurrency, getPlayerById, stateIn } from "../utils";
    import ActionButton from "./actions/ActionButton.svelte";
    import Button from "./atoms/Button.svelte";
    import GameCard from "./atoms/GameCard.svelte";

    const {
        peer,
        self,
    }: {
        peer: Peer;
        self?: boolean;
    } = $props();

    const isOut = peer.cards.filter((c) => !c.alive).length === 2;

    let overlayText = $derived.by(() => {
        if (self) {
            return;
        }
        // If this peer is the
        const isChoosingAction =
            global.state.turnState === TurnState.Default && peer.leading;
        if (isChoosingAction) {
            return "Choosing an action";
        }

        const isTargeted =
            stateIn(global.state, TurnState.ActionPending) &&
            global.state.pendingAction.target &&
            global.state.pendingAction.target === peer.id;
        if (isTargeted) {
            return "Choosing a response";
        }
        const isChoosingCardToLose =
            stateIn(
                global.state,
                TurnState.PlayerKilled,
                TurnState.PlayerLostChallenge,
                TurnState.LeaderLostChallenge,
            ) && global.state.nextDeath === peer.id;
        if (isChoosingCardToLose) {
            return "Choosing a card to kill";
        }

        const isBlocked =
            stateIn(global.state, TurnState.BlockPending) && peer.leading;
        if (isBlocked) {
            return `Blocked by ${getPlayerById(global.state, global.state.pendingBlock.initiator!)}`;
        }
    });
</script>

<div
    class={[
        "panel relative",
        { "outline-yellow-300 outline-double": peer.leading },
    ]}
>
    {#if overlayText}
        <div
            class="absolute left-0 h-full w-full flex justify-center items-center z-10"
        >
            <div class="panel">
                <h2>
                    {overlayText}
                </h2>
            </div>
        </div>
    {/if}

    <div
        class={[
            "flex flex-col gap-2 flex-1",
            {
                "brightness-50": isOut || overlayText,
            },
        ]}
    >
        <div class="flex flex-row gap-2 text-lg items-center">
            <div class="font-bold">
                {peer.name}
            </div>
            <div>
                {formatCurrency(peer.credits)}
            </div>
        </div>

        <!-- Always show two cards - for `self`, this is the current hand, and for peers this is any dead cards. -->
        <div class="flex flex-1 gap-1 justify-between">
            {#each { length: 2 } as _, i}
                <GameCard card={peer.cards[i]} />
            {/each}
        </div>

        {#if self}
            <!-- If we're rendering the current player, show independent actions like Income. -->

            <ActionButton disabled={isOut} type={ActionType.Income} />
            <ActionButton disabled={isOut} type={ActionType.ForeignAid} />
            <ActionButton disabled={isOut} type={ActionType.Tax} />
            <Button disabled>Exchange</Button>

            <!-- Otherwise, allow actions like Assassinate targeted at the current peer. -->
        {:else}
            <ActionButton
                disabled={isOut}
                type={ActionType.Revolt}
                target={peer.id}
            />
            <ActionButton
                disabled={isOut}
                type={ActionType.Assassinate}
                target={peer.id}
            />
            <ActionButton
                disabled={isOut}
                type={ActionType.Steal}
                target={peer.id}
            />
        {/if}
    </div>
</div>
