<script>
    import { global } from "../state.svelte";
</script>

<a href="/">
    <button onclick={() => global.client.leave()}>Leave Game</button>
</a>
