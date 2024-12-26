<script lang="ts">
    import { global } from "../../state.svelte";
    import Button from "../atoms/Button.svelte";
    import Card from "../atoms/Card.svelte";
    import Subtitle from "../atoms/Subtitle.svelte";
    import Table from "../atoms/Table.svelte";
    import Title from "../atoms/Title.svelte";
    import LeaveGame from "../LeaveGame.svelte";
</script>

<Card class="flex-col">
    <Title>Lobby</Title>
    <div class="flex gap-4 items-center">
        <Subtitle>
            Connected to: <span class="font-mono">{global.state.gameId}</span>
        </Subtitle>
    </div>
    <a href={window.location.href}>{window.location.href}</a>

    <Table
        headers={["peers"]}
        data={[
            global.state.self.name,
            ...global.state.peers.map((p) => p.name),
        ]}
    >
        {#snippet row(peer)}
            <tr>
                <td class="p-1 border-neutral-500 border">
                    {peer}
                </td>
            </tr>
        {/snippet}</Table
    >

    <div class="ml-auto">
        <Button
            onclick={() => global.client.startGame()}
            class="ml-auto"
            disabled={!global.state.self.leading}
        >
            Start Game
        </Button>
        <LeaveGame />
    </div>
</Card>
