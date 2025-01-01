<!-- Displays for all players when an action has been taken. Allows challenging -->
<script lang="ts">
    import { global } from "../state.svelte";
    import { TurnState } from "../types";
    import {
        formatCurrentAction,
        getCurrentActionBlockers,
        stateIn,
    } from "../utils";
    import ActionBlock from "./actions/ActionBlock.svelte";
    import Modal from "./atoms/Modal.svelte";

    let dialog: HTMLDialogElement;
    let closed = $state(false);

    const challenge = () => {
        global.client.challenge();
    };

    // Show for non leading peers after the leader has chosen an action.
    let open = $derived(
        !global.state.self.leading &&
            stateIn(global.state, TurnState.ActionPending),
    );

    $effect(() => {
        if (open) {
            dialog.showModal();
        } else {
            dialog.close();
        }
    });
</script>

<Modal {open}>
    <div class="flex gap-2 flex-col">
        <h1>
            {formatCurrentAction(global.state)}
        </h1>
        <!-- This will show if a) the current action is foreign aid or b) the current player is targeted. -->
        {#each getCurrentActionBlockers(global.state) as card}
            <ActionBlock {card} />
        {/each}
        <button onclick={challenge}>Challenge</button>
        <button onclick={() => dialog.close()}>Pass</button>
    </div>
</Modal>
