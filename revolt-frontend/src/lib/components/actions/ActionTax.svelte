<script lang="ts">
    import { global } from "../../state.svelte";
    import { ActionType } from "../../types";
    import { isAllowedAction } from "../../utils";
    import LieButton from "../atoms/LieButton.svelte";

    const tax = () => {
        global.client.attemptAction({
            type: ActionType.Tax,
        });
    };
</script>

<LieButton onclick={tax} lie={!isAllowedAction(global.state, ActionType.Tax)}>
    Tax
</LieButton>
