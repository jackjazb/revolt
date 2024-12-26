<script lang="ts">
    import { global } from "../../state.svelte";
    import { ActionType } from "../../types";
    import { isAllowedAction } from "../../utils";
    import Button from "../atoms/Button.svelte";

    let { target }: { target: string } = $props();

    const assassinate = () => {
        global.client.attemptAction({
            type: ActionType.Assassinate,
            target,
        });
    };
</script>

<Button onclick={assassinate}>
    Assassinate
    {#if !isAllowedAction(global.state, ActionType.Assassinate)}
        (lie)
    {/if}
</Button>
