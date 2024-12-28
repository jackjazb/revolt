<script>
    import { global } from "../state.svelte";
    import Button from "./atoms/Button.svelte";
</script>

<a href="/">
    <Button onclick={() => global.client.leave()}>Leave Game</Button>
</a>
