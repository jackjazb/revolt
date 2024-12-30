<script lang="ts">
    import { global } from "../state.svelte";
    import { TurnState } from "../types";
    import {
        formatCurrentAction,
        getCurrentActionBlockers,
        stateIn,
    } from "../utils";
    import ActionBlock from "./actions/ActionBlock.svelte";
    import ActionChallenge from "./actions/ActionChallenge.svelte";
    import Button from "./atoms/Button.svelte";

    let dialog: HTMLDialogElement;

    const commit = () => {
        global.client.commitTurn();
    };
    const end = () => {
        global.client.endTurn();
    };

    $effect(() => {
        if (
            !global.state.self.leading &&
            stateIn(global.state, TurnState.ActionPending)
        ) {
            dialog.showModal();
        } else {
            dialog.close();
        }
    });
</script>

<dialog
    bind:this={dialog}
    class="backdrop:backdrop-brightness-50 text-inherit bg-inherit"
>
    <div class="panel flex-col">
        <h1>
            {formatCurrentAction(global.state)}
        </h1>
        {#each getCurrentActionBlockers(global.state) as card}
            <ActionBlock {card} />
        {/each}
        <ActionChallenge />
        <Button onclick={() => dialog.close()}>Pass</Button>
    </div>
</dialog>
