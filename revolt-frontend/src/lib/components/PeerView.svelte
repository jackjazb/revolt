<!--
The base game view for the leader, primarily allowing action choice.
-->
<script lang="ts">
    import { global } from "../state.svelte";
    import { TurnState } from "../types";
    import Button from "./atomic/Button.svelte";
</script>

<div class="flex flex-col gap-2">
    {#if global.state.turnState === TurnState.Default}
        <h1>leader choosing action...</h1>
    {:else if global.state.turnState === TurnState.ActionPending && global.state.pendingAction}
        <h1>leader attempted {global.state.pendingAction.type}</h1>
        {#if global.state.pendingAction.target}
            <h2>targeting {global.state.pendingAction.target}!</h2>
        {/if}

        <Button>block</Button>
        <Button>challenge</Button>
    {/if}
</div>
