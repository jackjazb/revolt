<script lang="ts">
    import { global } from "../state.svelte";
    import { TurnState } from "../types";
    import { getPlayerById, stateIn } from "../utils";
    import GameCard from "./atoms/GameCard.svelte";
    import Modal from "./atoms/Modal.svelte";

    let forceClose = $state(false);

    let open = $derived(
        !forceClose &&
            stateIn(
                global.state,
                TurnState.PlayerKilled,
                TurnState.LeaderLostChallenge,
                TurnState.PlayerLostChallenge,
            ),
    );

    let isNextDeath = $derived(global.state.nextDeath === global.state.self.id);

    const resolveDeath = (card: number) => {
        global.client.resolveDeath(card);
    };

    $effect(() => {
        if (open && !isNextDeath) {
            setTimeout(() => {
                forceClose = true;
            }, 1000);
        }
    });
</script>

<Modal {open}>
    {#if isNextDeath}
        <div class="flex flex-col">
            <h1>Select a card to kill</h1>
            <div class="flex gap-2">
                {#each global.state.self.cards as card, i}
                    {#if card.alive}
                        <button onclick={() => resolveDeath(i)}>
                            <GameCard {card} />
                        </button>
                    {/if}
                {/each}
            </div>
        </div>
    {:else}
        <h1>
            {getPlayerById(global.state, global.state.nextDeath)} loses a card!
        </h1>
    {/if}
</Modal>
