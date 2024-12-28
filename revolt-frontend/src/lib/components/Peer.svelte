<script lang="ts">
    import type { Peer } from "../types";
    import { formatCurrency } from "../utils";
    import GameCard from "./atoms/GameCard.svelte";
    import Icon from "./atoms/Icon.svelte";
    import Panel from "./atoms/Panel.svelte";

    const { peer }: { peer: Peer } = $props();
</script>

<Panel>
    <div class="flex flex-row gap-2 text-lg items-center">
        <div class="font-bold">
            {peer.name}
        </div>
        ~
        <div>
            {formatCurrency(peer.credits)}
        </div>
        {#if peer.leading}
            <div
                class="ml-auto bg-yellow-300 border-4 border-black border-double text-black rounded-sm"
            >
                <Icon type="crown" size={20} />
            </div>
        {/if}
    </div>
    <div class="flex flex-1 gap-1 justify-between">
        {#each { length: 2 } as _, i}
            <GameCard card={peer.cards[i]} />
        {/each}
    </div>
</Panel>
