<script lang="ts" generics="T">
    import type { Snippet } from "svelte";

    const {
        headers,
        data,
        row,
    }: {
        headers: string[];
        data: T[];
        row: Snippet<[T]>;
    } = $props();
</script>

<table>
    <thead>
        <tr>
            {#each headers as header}
                <th class="border-neutral-500 border">
                    {header}
                </th>
            {/each}
        </tr>
    </thead>
    <tbody>
        {#each data as r}
            {@render row(r)}
        {/each}
    </tbody>
</table>
