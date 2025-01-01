<script lang="ts">
    import { global } from "../../state.svelte";
    import { getLeader } from "../../utils";
    import LeaveGame from "../LeaveGame.svelte";
</script>

<div class="panel flex-col">
    <h1>Waiting for {getLeader(global.state)} to start the game.</h1>
    <div class="text-base">
        Connected to
        {global.state.gameId}
        as
        {global.state.self.name} // Game link:
        <a class="underline" href={window.location.href}>
            {window.location.href}
        </a>
    </div>

    <h1>Players</h1>
    <ul class="list-disc list-inside">
        {#each global.state.peers as peer}
            <li>
                {peer.name}
            </li>
        {/each}
    </ul>
    <div class="ml-auto">
        <button
            onclick={() => global.client.startGame()}
            class="ml-auto"
            disabled={!global.state.self.leading}
        >
            Start Game
        </button>
        <LeaveGame />
    </div>
</div>
