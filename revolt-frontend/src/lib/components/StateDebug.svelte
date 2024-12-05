<script lang="ts">
    import { Card, type State } from "../types";

    let { gameState }: { gameState: State } = $props();
</script>

<div class="font-mono">
    <p>
        game ID: {gameState.gameId}
    </p>
    <p>
        client ID: {gameState.clientId}
    </p>
    <p>
        status: {gameState.status}
    </p>
    <p>
        leading : {gameState.isLeader}
    </p>
    <p>
        my cards: {gameState.self.cards
            .filter((c) => c.card !== Card.Empty)
            .map((c) => c.card)
            .join(", ")}
    </p>
</div>
