<script lang="ts">
    import type { Snippet } from "svelte";

    const {
        children,
        ...rest
    }: {
        onclick?: () => void;
        disabled?: boolean;
        children: Snippet;
        class?: string;
    } = $props();
</script>

<button
    {...rest}
    class={`
        rounded-sm border-transparent border bg-neutral-800
        py-1 px-2
        hover:border-blue-400 duration-200
        focus:outline-4 
        disabled:brightness-75
        disabled:border-none
        ${rest.class}
        `}
>
    {@render children()}
</button>
