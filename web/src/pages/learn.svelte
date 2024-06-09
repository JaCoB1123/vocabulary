<script>
    import Word from './_components/word.svelte';

    let words = [];
    let loading = true;

    let index = 0;
    $: word = index>=words.length ? {} : words[index];
    
    fetch("/api/learn")
        .then(x => x.json())
        .then(x => words = x)
        .then(() => loading = false);

    function onCorrect(word) {
        index++;
    }

    function onFalse(word) {
        index++;
    }
</script>

<main>
    <pre>{JSON.stringify(word)}</pre>
    {#if loading}
        Loading...
    {:else}
        <Word word={word} on:correct={onCorrect} on:false={onFalse} />
    {/if}
</main>