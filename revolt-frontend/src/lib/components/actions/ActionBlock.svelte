<script lang="ts">
    import { global } from "../../state.svelte";
    import type { Block, Card } from "../../types";
    import { formatCard } from "../../utils";
    import LieButton from "../atoms/LieButton.svelte";

    let { card }: { card: Card } = $props();

    const block = () => {
        const block: Block = { card };
        global.client.attemptBlock(block);
    };
</script>

<LieButton
    onclick={block}
    lie={!global.state.self.cards.map((card) => card.card).includes(card)}
>
    Block with {formatCard(card)}
</LieButton>
