<script lang="ts">
    import Button from "./lib/components/atomic/Button.svelte";
    import Card from "./lib/components/atomic/Card.svelte";
    import TextInput from "./lib/components/atomic/TextInput.svelte";
    import Game from "./lib/components/Game.svelte";
    import Icon from "./lib/components/Icon.svelte";
    import Lobby from "./lib/components/Lobby.svelte";
    import { global } from "./lib/state.svelte";
    import { GameStatus } from "./lib/types";
    import { randomName } from "./lib/utils";

    let gameIdInput = $state("");
    let playerNameInput = $state(randomName());
    const connection = global.client.connect("ws://localhost:8080");
</script>

<main class="max-w-5xl mx-auto px-12 my-12">
    {#await connection}
        <p>establishing connection...</p>
    {:then _}
        {#if global.state.status === GameStatus.Default}
            <div class="flex flex-col gap-2">
                <Card>
                    <TextInput
                        bind:value={playerNameInput}
                        placeholder="name"
                    />
                    <Button onclick={() => (playerNameInput = randomName())}
                        ><Icon type="dice" />
                    </Button>
                </Card>

                <Card>
                    <Button
                        class="mr-auto"
                        onclick={() =>
                            global.client.createGame(playerNameInput)}
                        >create</Button
                    >
                </Card>

                <Card>
                    <TextInput bind:value={gameIdInput} placeholder="game id" />
                    <Button
                        onclick={() =>
                            global.client.joinGame(
                                gameIdInput,
                                playerNameInput,
                            )}
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
