<script lang="ts">
    import { global } from "../state.svelte";
    import Card from "./atomic/Card.svelte";
    import Table from "./atomic/Table.svelte";
    import Title from "./atomic/Title.svelte";
    import LeaderView from "./LeaderView.svelte";
    import PeerCards from "./PeerCards.svelte";
    import PeerView from "./PeerView.svelte";
</script>

<div class="flex w-full gap-4">
    <Card>
        <Title>players</Title>
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
    <Card>
        {#if global.state.self.leading}
            <LeaderView />
        {:else}
            <PeerView />
        {/if}
    </Card>
    <Card>
        <Title>cards</Title>
    </Card>
</div>
