<script lang="ts">
    import { global } from "../../state.svelte";
    import Button from "../atoms/Button.svelte";

    const commit = () => {
        global.client.commitTurn();
    };
</script>

<Button onclick={commit}>commit</Button>
