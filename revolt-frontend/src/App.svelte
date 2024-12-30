<script lang="ts">
    import { ClientStatus } from "./lib/client";
    import Button from "./lib/components/atoms/Button.svelte";
    import Icon from "./lib/components/atoms/Icon.svelte";
    import Game from "./lib/components/pages/Game.svelte";
    import Lobby from "./lib/components/pages/Lobby.svelte";
    import TestGame from "./lib/components/pages/TestGame.svelte";
    import { global } from "./lib/state.svelte";
    import { GameStatus, type CreateGameResponse } from "./lib/types";
    import { randomName } from "./lib/utils";

    let gameIdInput = $state("");
    let playerNameInput = $state(randomName());

    $effect(() => {
        console.log("status", global.status);
    });

    const path = window.location.pathname.split("/").slice(1);
    let gameId = $state(path.length > 0 ? path[0] : undefined);
    const debug = $derived(gameId === "test");

    const createGame = async () => {
        const res = await fetch("http://localhost:8080/create", {
            method: "POST",
        });
        const response = (await res.json()) as CreateGameResponse;
        window.history.pushState({}, "", `/${response.id}`);
        gameId = response.id;
    };

    async function handleGameUrl() {
        if (gameId) {
            await global.client.connect(`ws://localhost:8080/${gameId}`);
        }
    }
</script>

<main class="px-4 my-4">
    {#if debug}
        <TestGame />
    {:else if gameId}
        {#await handleGameUrl() then _}
            {#if global.state.status === GameStatus.Lobby}
                <h1 class="text-5xl my-5">Lobby</h1>

                <Lobby />
            {:else if global.state.status === GameStatus.InProgress}
                <Game />
            {/if}
            {#if global.status === ClientStatus.Connected}{:else if global.client.status === ClientStatus.Connecting}
                <h1>loading...</h1>
            {:else if global.client.status === ClientStatus.Failed}
                <h1>Connection failed.</h1>
            {/if}
        {/await}
    {:else}
        <h1 class="text-5xl my-5">Revolt</h1>
        <div class="flex flex-col gap-2">
            <div class="panel flex-row">
                <input
                    type="text"
                    bind:value={playerNameInput}
                    placeholder="name"
                />
                <button
                    onclick={() => (playerNameInput = randomName())}
                    class="border-4 border-black border-double bg-gray-800 rounded-sm hover:brightness-75 duration-75"
                >
                    <Icon type="dice" />
                </button>
            </div>

            <div class="panel">
                <Button onclick={createGame}>Create Game</Button>
            </div>
        </div>
    {/if}
</main>
