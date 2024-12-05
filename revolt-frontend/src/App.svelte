<script lang="ts">
    import { Client } from "./lib/Client";
    import Game from "./lib/components/Game.svelte";
    import Icon from "./lib/components/Icon.svelte";
    import Lobby from "./lib/components/Lobby.svelte";
    import StateDebug from "./lib/components/StateDebug.svelte";
    import { GameStatus, type State } from "./lib/types";
    import { randomName } from "./lib/utils";

    let gameIdInput = $state("");
    let playerNameInput = $state(randomName());
    let gameState: State | undefined = $state();

    // Callback here to trigger rerenders on state update.
    let client = $state(new Client((state) => (gameState = state)));

    const connection = client.connect("ws://localhost:8080");
</script>

<main class="max-w-xl mx-auto my-12">
    <div class="flex flex-col gap-4">
        <h1 class="font-bold text-3xl">revolt</h1>

        {#await connection}
            <p>establishing connection...</p>
        {:then _}
            <div class="card flex gap-2">
                <input
                    type="text"
                    bind:value={playerNameInput}
                    placeholder="name"
                />
                <button onclick={() => (playerNameInput = randomName())}
                    ><Icon type="dice" />
                </button>
            </div>

            <div class="card flex">
                <button
                    class="mr-auto"
                    onclick={() => client.createGame(playerNameInput)}
                    >create</button
                >
            </div>

            <div class="flex gap-2 card">
                <input
                    type="text"
                    bind:value={gameIdInput}
                    placeholder="game id"
                />
                <button
                    onclick={() =>
                        client.joinGame(gameIdInput, playerNameInput)}
                    >join
                </button>
            </div>

            <button onclick={() => client.startGame()}>start</button>

            {#if gameState}
                <StateDebug {gameState} />
                {#if gameState.status === GameStatus.Lobby}
                    <Lobby {gameState} />
                {:else if gameState.status === GameStatus.InProgress}
                    <Game {gameState} {client} />
                {/if}
            {/if}
        {/await}
    </div>
</main>
