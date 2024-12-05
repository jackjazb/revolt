<script lang="ts">
    import type { State } from "../types";

    let { gameState }: { gameState: State } = $props();
</script>

<h1 class="text-xl">in lobby</h1>

<ul class="list-disc">
    {#each gameState.clients as client}
        <li>{client.name}</li>
    {/each}
</ul>
