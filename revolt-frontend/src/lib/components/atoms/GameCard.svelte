<script lang="ts">
    import { Card, type CardState } from "../../types";
    import Icon from "./Icon.svelte";

    let { card }: { card?: CardState } = $props();

    function resolveImageSource(card: Card) {
        const sources: Record<Card, string> = {
            [Card.Duke]: "cards/duke.jpeg",
            [Card.Ambassador]: "cards/ambassador.jpeg",
            [Card.Contessa]: "cards/contessa.jpeg",
            [Card.Captain]: "cards/captain.jpeg",
            [Card.Assassin]: "cards/assassin.jpeg",
        };
        return `/${sources[card]}`;
    }
</script>

<div class="w-20 h-20 border-double border-black border-4 bg-gray-800">
    {#if card}
        <div class="relative">
            <img
                src={resolveImageSource(card.card)}
                alt={card.card}
                class={[!card.alive && "brightness-[0.4]"]}
            />
            {#if !card.alive}
                <span
                    class="absolute top-0 w-full h-full flex items-center justify-center"
                >
                    <Icon type="skull" />
                </span>
            {/if}
        </div>
    {:else}
        <div class="w-full h-full flex justify-center items-center text-2xl">
            ?
        </div>
    {/if}
</div>
