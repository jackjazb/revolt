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
    class={[
        `bg-gray-800 border-4 border-black  border-double
        rounded-sm 
        px-2
        text-base
        hover:border-blue-400 duration-200
        focus:outline-4 
        disabled:brightness-75
        disabled:hover:border-black`,
        rest.class,
    ]}
>
    {@render children()}
</button>
