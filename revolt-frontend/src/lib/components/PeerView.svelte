<!--
The base game view for the leader, primarily allowing action choice.
-->
<script lang="ts">
    import { TurnState, type GameComponentContext } from "../types";

    let { gameState, client }: GameComponentContext = $props();
</script>

{#if gameState.turnState === TurnState.Default}
    <h1>leader choosing action...</h1>
{:else}
    <h1>leader attempted {gameState.pendingAction.type}</h1>
    {#if gameState.pendingAction.target}
        <h2>targeting {gameState.pendingAction.target}!</h2>
    {/if}

    <button>block</button>
    <button>challenge</button>
{/if}
