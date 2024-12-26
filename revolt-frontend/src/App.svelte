<script lang="ts">
    import { ClientStatus } from "./lib/client";
    import Button from "./lib/components/atoms/Button.svelte";
    import Card from "./lib/components/atoms/Card.svelte";
    import Icon from "./lib/components/atoms/Icon.svelte";
    import TextInput from "./lib/components/atoms/TextInput.svelte";
    import Title from "./lib/components/atoms/Title.svelte";
    import Game from "./lib/components/pages/Game.svelte";
    import Lobby from "./lib/components/pages/Lobby.svelte";
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

    const post = async () => {
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

<main class="max-w-5xl mx-auto px-12 my-12">
    {#if gameId}
        {#await handleGameUrl() then _}
            {#if global.state.status === GameStatus.Lobby}
                <Lobby />
            {:else if global.state.status === GameStatus.InProgress}
                <Game />
            {/if}
            {#if global.status === ClientStatus.Connected}{:else if global.client.status === ClientStatus.Connecting}
                <h1>loading...</h1>
            {:else if global.client.status === ClientStatus.Failed}
                <Title>Connection failed.</Title>
            {/if}
        {/await}
    {:else}
        <div class="flex flex-col gap-2">
            <Card row>
                <TextInput bind:value={playerNameInput} placeholder="name" />
                <Button onclick={() => (playerNameInput = randomName())}>
                    <Icon type="dice" />
                </Button>
            </Card>

            <Card>
                <Button onclick={post}>create game</Button>
            </Card>
        </div>
    {/if}
</main>
