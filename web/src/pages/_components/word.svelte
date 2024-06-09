<script>
	import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();
    
    export let word;
    let resultVisible = false;

    function showResult() {
        resultVisible = !resultVisible;
    }

    function onCorrect() {
		dispatch('correct', word);
    }

    function onFalse() {
		dispatch('false', word);
    }
</script>

<main>
    <texts>
        <original>{word.Name}</original>
        {#if resultVisible}
            <result>{word.Translation}</result>
        {/if}
    </texts>
    <buttons>
        <button on:click={onCorrect}>Richtig</button>
        <button on:click={showResult}>Lösung</button>
        <button on:click={onFalse}>Falsch</button>
    </buttons>
</main>

<style>
    main {
        margin: 0 auto;
        width: 300px;
        height: 300px;
        padding: 1em;
        display: flex;
        flex-direction: column;
        background-color: aquamarine;
    }

    buttons {
        display: flex;
        justify-content: space-between;
        flex: 1 1;
    }

    button {
        margin: 0;
    }

    texts {
        display: flex;
        align-items: center;
        flex-direction: column;
        justify-content: space-around;
        flex: 5 1;
    }
</style>