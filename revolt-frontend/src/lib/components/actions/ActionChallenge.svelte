<script lang="ts">
    import { global } from "../../state.svelte";
    import Button from "../atoms/Button.svelte";

    const challenge = () => {
        global.client.challenge();
    };
</script>

<Button onclick={challenge}>challenge</Button>
