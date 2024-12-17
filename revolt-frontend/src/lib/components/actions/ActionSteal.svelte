<script lang="ts">
    import { global } from "../../state.svelte";
    import { ActionType } from "../../types";
    import { isAllowedAction } from "../../utils";
    import Button from "../atoms/Button.svelte";

    let { target }: { target: string } = $props();

    const steal = () => {
        global.client.attemptAction({
            type: ActionType.Steal,
            target,
        });
    };
</script>

<Button onclick={steal}>
    steal
    {#if !isAllowedAction(global.state, ActionType.Steal)}
        (lie)
    {/if}
</Button>
