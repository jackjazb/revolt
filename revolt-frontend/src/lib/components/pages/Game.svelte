<script lang="ts">
    import { global } from "../../state.svelte";
    import ActionPendingDialog from "../ActionPendingDialog.svelte";
    import BlockPendingDialog from "../BlockPendingDialog.svelte";
    import GameState from "../GameState.svelte";
    import LeaveGame from "../LeaveGame.svelte";
    import Peer from "../Peer.svelte";
    import StateSummary from "../StateSummary.svelte";
</script>

<div class="grid grid-cols-12 gap-2 w-full h-full">
    <!-- State summary, identical for all players. Rendered as the top panel -->
    <div class="panel flex-row col-span-12">
        <StateSummary />
        <span class="ml-auto">
            <LeaveGame />
        </span>
    </div>

    <div class="col-span-2">
        <Peer self peer={global.state.self} />
    </div>
    <div class="col-span-10 mb-auto">
        <div class=" grid grid-cols-5 mb-auto gap-2">
            {#each global.state.peers as peer}
                <Peer {peer} />
            {/each}
        </div>

        <ActionPendingDialog />
        <BlockPendingDialog />
        <GameState />
    </div>
</div>
