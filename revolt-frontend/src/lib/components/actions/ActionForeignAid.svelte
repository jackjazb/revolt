<script lang="ts">
    import { global } from "../../state.svelte";
    import { ActionType } from "../../types";
    import Button from "../atoms/Button.svelte";

    const foreignAid = () => {
        global.client.attemptAction({
            type: ActionType.ForeignAid,
        });
    };
</script>

<Button onclick={foreignAid}>foreign aid</Button>
