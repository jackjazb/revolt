<script lang="ts">
    import { global } from "../../state.svelte";
    import { ActionType } from "../../types";
    import { formatCurrency } from "../../utils";
    import Button from "../atoms/Button.svelte";

    let { target }: { target: string } = $props();

    const revolt = () => {
        global.client.attemptAction({
            type: ActionType.Revolt,
            target,
        });
    };
</script>

<Button onclick={revolt}>
    Revolt (-{formatCurrency(7)})
</Button>
