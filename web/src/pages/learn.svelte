<script>
    import Word from './_components/word.svelte';

    let words = [];
    let loading = true;

    let index = 0;
    $: word = words[index];
    
    fetch("/api/learn")
        .then(x => x.json())
        .then(x => words = x)
        .then(() => loading = false);
</script>

<main>
    {#if loading}
        Loading...
    {:else}
        <Word word={word} />
        <ul>
        {#each words as word}
            <li>{JSON.stringify(word)}</li>
        {/each}
        </ul>
        <pre></pre>
    {/if}
</main>

<style>
</style>