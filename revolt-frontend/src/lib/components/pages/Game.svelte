<script lang="ts">
    import { global } from "../../state.svelte";
    import Panel from "../atoms/Panel.svelte";
    import LeaderView from "../GameState.svelte";
    import LeaveGame from "../LeaveGame.svelte";
    import Peer from "../Peer.svelte";
    import StateSummary from "../StateSummary.svelte";
</script>

<div class="grid grid-cols-12 gap-2 w-full h-full">
    <Panel row class="col-span-12"
        ><StateSummary /><span class="ml-auto"><LeaveGame /></span></Panel
    >
    <Panel class="col-span-4 grid grid-cols-2 mb-auto">
        <Peer peer={global.state.self} />
        <div></div>
        {#each global.state.peers as peer}
            <Peer {peer} />
        {/each}
    </Panel>
    <Panel class="col-span-8 mb-auto">
        <LeaderView />
    </Panel>
</div>
