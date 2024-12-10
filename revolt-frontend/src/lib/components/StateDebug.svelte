<script lang="ts">
    import { global } from "../state.svelte";
    import { Card } from "../types";
</script>

<div class="font-mono">
    <p>
        game ID: {global.state.gameId}
    </p>
    <p>
        client ID: {global.state.self.id}
    </p>
    <p>
        status: {global.state.status}
    </p>
    <p>
        leading: {global.state.self.leading}
    </p>
    <p>
        my cards: {global.state.self.cards
            .filter((c) => c.card !== Card.Empty)
            .map((c) => c.card)
            .join(", ")}
    </p>
</div>
