<script lang="ts">
    import { global } from "../../state.svelte";
    import type { Block, Card } from "../../types";
    import Button from "../atoms/Button.svelte";

    let { card }: { card: Card } = $props();

    const block = () => {
        const block: Block = { card };
        global.client.attemptBlock(block);
    };
</script>

<Button onclick={block}>
    block with {card}
    {#if !global.state.self.cards.map((card) => card.card).includes(card)}
        (lie)
    {/if}
</Button>
