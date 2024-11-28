<script lang="ts">
    import { Client } from "../Client";

    let gameIdInput = $state("");

    let client = new Client();
    void client.connect("ws://localhost:8080");
</script>

<div class="flex flex-col gap-4">
    <button onclick={() => client.createGame()}>create</button>
    <div>
        <input type="text" bind:value={gameIdInput} placeholder="game id" />
        <button onclick={() => client.joinGame(gameIdInput)}>join </button>
    </div>
    <button onclick={() => client.startGame()}>start</button>
</div>

<p>my cards: {client.state?.self.cards.map((c) => c.card).join(", ")}</p>
