<script lang="ts">
    import { global } from "../../state.svelte";
    import { ActionType } from "../../types";
    import Button from "../atoms/Button.svelte";

    const income = () => {
        global.client.attemptAction({
            type: ActionType.Income,
        });
    };
</script>

<Button onclick={income}>Income</Button>
