<script lang="ts">
    import { global } from "../../state.svelte";
    import { ActionType } from "../../types";
    import { isAllowedAction } from "../../utils";
    import LieButton from "../atoms/LieButton.svelte";

    let { target }: { target: string } = $props();

    const assassinate = () => {
        global.client.attemptAction({
            type: ActionType.Assassinate,
            target,
        });
    };
</script>

<LieButton
    onclick={assassinate}
    lie={!isAllowedAction(global.state, ActionType.Assassinate)}
>
    Assassinate
</LieButton>
