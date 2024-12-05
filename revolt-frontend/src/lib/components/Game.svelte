<script lang="ts">
    import { type GameComponentContext } from "../types";
    import LeaderView from "./LeaderView.svelte";
    import PeerView from "./PeerView.svelte";

    let props: GameComponentContext = $props();
</script>

{#if props.gameState.isLeader}
    <LeaderView {...props} />
{:else}
    <PeerView {...props} />
{/if}
