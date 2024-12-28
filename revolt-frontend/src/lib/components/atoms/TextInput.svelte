<script lang="ts">
    let {
        value = $bindable(""),
        placeholder,
    }: {
        value: string;
        placeholder: string;
    } = $props();
</script>

<input
    class="px-2 rounded-sm bg-neutral-800 border-4 border-black border-double"
    type="text"
    {placeholder}
    bind:value
/>
