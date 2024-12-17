<script lang="ts">
    import { global } from "../../state.svelte";
    import Button from "../atoms/Button.svelte";
    import Card from "../atoms/Card.svelte";
    import Subtitle from "../atoms/Subtitle.svelte";
    import Table from "../atoms/Table.svelte";
    import Title from "../atoms/Title.svelte";
</script>

<Card class="flex-col">
    <Title>lobby</Title>

    <div class="flex gap-2">
        <Subtitle>
            Game ID: <span class="font-mono">{global.state.gameId}</span>
        </Subtitle>
    </div>

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

    {#if global.state.ownerId === global.state.self.id}
        <Button onclick={() => global.client.startGame()} class="ml-auto"
            >start</Button
        >
    {:else}
        <Subtitle>Waiting for the leader...</Subtitle>
    {/if}
</Card>
