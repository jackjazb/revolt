<script lang="ts">
    import { global } from "../state.svelte";
    import { TurnState } from "../types";
    import { formatCurrentBlock, stateIn } from "../utils";

    let dialog: HTMLDialogElement;

    let open = $derived(
        stateIn(global.state, TurnState.BlockPending) &&
            global.state.self.leading,
    );

    const commit = () => {
        global.client.commitTurn();
    };

    const challenge = () => {
        global.client.challenge();
    };

    $effect(() => {
        if (open) {
            dialog.showModal();
        } else {
            dialog.close();
        }
    });
</script>

<dialog bind:this={dialog}>
    <div class="panel flex-col">
        <h1>{formatCurrentBlock(global.state)}</h1>
        <button onclick={challenge}>Challenge</button>
        <button onclick={commit}>Accept</button>
    </div>
</dialog>
