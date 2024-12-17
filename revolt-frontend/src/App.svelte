<script lang="ts">
    import Button from "./lib/components/atoms/Button.svelte";
    import Card from "./lib/components/atoms/Card.svelte";
    import Icon from "./lib/components/atoms/Icon.svelte";
    import TextInput from "./lib/components/atoms/TextInput.svelte";
    import Game from "./lib/components/pages/Game.svelte";
    import Lobby from "./lib/components/pages/Lobby.svelte";
    import { global } from "./lib/state.svelte";
    import { GameStatus, type CreateGameResponse } from "./lib/types";
    import { randomName } from "./lib/utils";

    let gameIdInput = $state("");
    let playerNameInput = $state(randomName());

    const path = window.location.pathname.split("/").slice(1);
    let gameId = $state(path.length > 0 ? path[0] : undefined);

    $effect(() => {
        console.log(gameId);
    });

    const post = async () => {
        const res = await fetch("http://localhost:8080/create", {
            method: "POST",
        });
        const response = (await res.json()) as CreateGameResponse;
        console.log(response);
        window.history.pushState({}, "", `/${response.id}`);
        gameId = response.id;
    };

    async function handleGameUrl() {
        if (gameId) {
            await global.client.connect("ws://localhost:8080");
            global.client.joinGame(gameId, randomName());
        }
    }
</script>

<main class="max-w-5xl mx-auto px-12 my-12">
    {#await handleGameUrl() then _}
        {#if global.state.status === GameStatus.Default}
            <div class="flex flex-col gap-2">
                <Card row>
                    <TextInput
                        bind:value={playerNameInput}
                        placeholder="name"
                    />
                    <Button onclick={() => (playerNameInput = randomName())}>
                        <Icon type="dice" />
                    </Button>
                </Card>

                <Card>
                    <Button onclick={post}>test</Button>
                </Card>

                <Card row>
                    <TextInput bind:value={gameIdInput} placeholder="game id" />
                    <Button
                        onclick={async () => {
                            await global.client.connect("ws://localhost:8080");
                            global.client.joinGame(
                                gameIdInput,
                                playerNameInput,
                            );
                        }}
                        >join
                    </Button>
                </Card>
            </div>
        {:else if global.state.status === GameStatus.Lobby}
            <Lobby />
        {:else if global.state.status === GameStatus.InProgress}
            <Game />
        {/if}
    {/await}
</main>
