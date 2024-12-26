<script lang="ts">
    import { global } from "../../state.svelte";
    import Button from "../atoms/Button.svelte";

    const { text }: { text?: string } = $props();
    const commit = () => {
        global.client.commitTurn();
    };
</script>

<Button onclick={commit}>{text ?? "commit"}</Button>
