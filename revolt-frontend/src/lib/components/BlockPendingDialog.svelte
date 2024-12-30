<script lang="ts">
    import { global } from "../state.svelte";
    import { TurnState } from "../types";
    import { formatCurrentBlock, stateIn } from "../utils";
    import Button from "./atoms/Button.svelte";

    let dialog: HTMLDialogElement;

    const commit = () => {
        global.client.commitTurn();
    };
    const challenge = () => {
        global.client.challenge();
    };

    $effect(() => {
        if (
            stateIn(global.state, TurnState.BlockPending) &&
            global.state.self.leading
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
        <h1>{formatCurrentBlock(global.state)}</h1>
        <Button onclick={challenge}>Challenge</Button>
        <Button onclick={commit}>Accept</Button>
    </div>
</dialog>
