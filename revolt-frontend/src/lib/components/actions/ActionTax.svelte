<script lang="ts">
    import { global } from "../../state.svelte";
    import { ActionType } from "../../types";
    import { isAllowedAction } from "../../utils";
    import Button from "../atoms/Button.svelte";

    const tax = () => {
        global.client.attemptAction({
            type: ActionType.Tax,
        });
    };
</script>

<Button onclick={tax}>
    tax
    {#if !isAllowedAction(global.state, ActionType.Tax)}
        (lie)
    {/if}
</Button>
