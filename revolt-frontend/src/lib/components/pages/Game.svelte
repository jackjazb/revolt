<script lang="ts">
    import { global } from "../../state.svelte";
    import Card from "../atoms/Card.svelte";
    import Subtitle from "../atoms/Subtitle.svelte";
    import Table from "../atoms/Table.svelte";
    import Title from "../atoms/Title.svelte";
    import LeaderView from "../LeaderView.svelte";
    import PeerCards from "../PeerCards.svelte";
    import PeerView from "../PeerView.svelte";
</script>

<div class="grid grid-cols-4 gap-2 w-full h-full">
    <Card class="col-span-1 row-span-2">
        <Title>game {global.state.gameId}</Title>
        <Table headers={["name", "credits", "cards"]} data={global.state.peers}>
            {#snippet row(peer)}
                <tr>
                    <td class="p-1 border-neutral-500 border">
                        {peer.name}
                    </td>
                    <td class="p-1 border-neutral-500 border">
                        {peer.credits}
                    </td>
                    <td class="p-1 border-neutral-500 border">
                        <PeerCards cards={peer.cards} />
                    </td>
                </tr>
            {/snippet}
        </Table>
    </Card>
    <Card class="col-start-2 col-span-4">
        <div class="bg-neutral-800 p-2 mr-auto">
            <Subtitle>{global.state.self.credits}₡</Subtitle>
        </div>
        <div class="flex gap-2">
            {#each global.state.self.cards as card}
                <div class="bg-neutral-800 p-2">
                    <Subtitle>{card.card}</Subtitle>
                    {#if card.alive}
                        alive
                    {:else}
                        dead
                    {/if}
                </div>
            {/each}
        </div>
    </Card>
    <Card class="col-span-4">
        {#if global.state.self.leading}
            <LeaderView />
        {:else}
            <PeerView />
        {/if}
    </Card>
</div>
