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
    class="py-1 px-2 rounded-sm bg-neutral-800"
    type="text"
    {placeholder}
    bind:value
/>
