<script lang="ts">
    import { global } from "../../state.svelte";
    import { formatCurrency } from "../../utils";
    import Card from "../atoms/Card.svelte";
    import Subtitle from "../atoms/Subtitle.svelte";
    import Table from "../atoms/Table.svelte";
    import LeaderView from "../GameState.svelte";
    import LeaveGame from "../LeaveGame.svelte";
    import PeerCards from "../PeerCards.svelte";
    import StateSummary from "../StateSummary.svelte";
</script>

<div class="grid grid-cols-4 gap-2 w-full h-full">
    <Card row class="col-span-5"
        ><StateSummary /><span class="ml-auto"><LeaveGame /></span></Card
    >
    <Card class="col-span-1 row-span-2">
        <Table data={global.state.peers}>
            {#snippet row(peer)}
                <tr>
                    <td class="p-1 border-neutral-500 border">
                        <div class="font-bold">
                            {peer.name}
                        </div>
                        <div>
                            {formatCurrency(peer.credits)}
                        </div>
                    </td>
                    <td class="p-1 border-neutral-500 border">
                        <PeerCards cards={peer.cards} />
                    </td>
                </tr>
            {/snippet}
        </Table>
    </Card>
    <Card class="col-start-2 col-span-4">
        <div class="flex flex-row gap-2 items-center">
            <div>
                <Subtitle>{global.state.self.name}</Subtitle>
            </div>
            ~
            <div class="mr-auto">
                <Subtitle>{formatCurrency(global.state.self.credits)}</Subtitle>
            </div>
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
        <LeaderView />
    </Card>
</div>
