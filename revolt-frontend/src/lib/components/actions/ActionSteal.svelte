<script lang="ts">
    import { global } from "../../state.svelte";
    import { ActionType } from "../../types";
    import { formatCurrency, isAllowedAction } from "../../utils";
    import LieButton from "../atoms/LieButton.svelte";

    let { target }: { target: string } = $props();

    const steal = () => {
        global.client.attemptAction({
            type: ActionType.Steal,
            target,
        });
    };
</script>

<LieButton
    onclick={steal}
    lie={!isAllowedAction(global.state, ActionType.Steal)}
>
    Steal (+{formatCurrency(2)})
</LieButton>
