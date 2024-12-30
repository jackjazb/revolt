<script lang="ts">
    import { global } from "../../state.svelte";
    import { ActionType } from "../../types";
    import {
        formatActionType,
        formatCurrency,
        isAllowedAction,
    } from "../../utils";
    import Button from "../atoms/Button.svelte";
    import Icon from "../atoms/Icon.svelte";

    let {
        type,
        target,
        disabled,
    }: {
        type: ActionType;
        target?: string;
        disabled?: boolean;
    } = $props();

    const actionCredits: Record<ActionType, number> = {
        [ActionType.Empty]: 0,
        [ActionType.Income]: 1,
        [ActionType.ForeignAid]: 2,
        [ActionType.Tax]: 3,
        [ActionType.Assassinate]: -3,
        [ActionType.Revolt]: -7,
        [ActionType.Steal]: 2,
        [ActionType.Exchange]: 0,
    };

    const profit = actionCredits[type];
    const lie = !isAllowedAction(global.state, type);

    const canAfford = global.state.self.credits + profit >= 0;
    const attempt = () => {
        global.client.attemptAction({
            type,
            target,
        });
    };
</script>

<Button
    class="flex flex-row gap-2 items-center content justify-center"
    onclick={attempt}
    disabled={disabled || !canAfford || !global.state.self.leading}
>
    {formatActionType(type)}
    {#if profit}
        ({profit >= 0 ? "+" : ""}{formatCurrency(profit)})
    {/if}
    {#if lie}
        <span class="text-red-500">
            <Icon type="lie" size={16} />
        </span>
    {/if}
</Button>
