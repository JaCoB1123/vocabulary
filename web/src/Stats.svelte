<script>
	let stats;
	fetch("/vocabulary/stats")
				.then(x =>x.json())
                .then(x => stats = x);
                
    function formatTime(time) {
        time = time/1000/1000/1000/60;
        if( time > 153722867){
            return "never";
        }

        if(time < 60 ){
            return  Math.round(time) + " minutes"
        }

        time = time / 60;
        if(time < 24 ){
            return  Math.round(time) + " hours"
        }

        time = time/24
        if(time < 7 ){
            return  Math.round(time) + " days"
        }

        if(time < 30 ){
            return  Math.round(time/7) + " weeks"
        }

        time = time/30
    	return  Math.round(time) + " months"
    }
</script>

<main>
    {#if stats}
        <h1>Stats</h1>
        <stat>Total Words<value>{stats.TotalWords}</value></stat>
        <stat>Total Answers<value>{stats.TotalAnswers}</value></stat>
        <stat>Always Correct<value>{stats.AlwaysCorrect}</value></stat>
        <stat>Words Answered<value>{stats.WordsAnswered}</value></stat>
        <stat>Total Due<value>{stats.TotalDue}</value></stat>

        <h2>Words by Levels</h2>
        {#each stats.LevelStats as count, level}
            <stat>{level}<value>{count}</value></stat>
        {/each}

        <h2>Words by Tags</h2>
        {#each Object.entries(stats.TagStats) as [tag, count]}
            <stat>{tag}<value>{count}</value></stat>
        {/each}

        <h2>Words by Recency</h2>
        {#each Object.entries(stats.TimeStats) as [time, count]}
            <stat>{formatTime(time)}<value>{count}</value></stat>
        {/each}
	{/if}
</main>

<style>
    stat {
        display: block;
        width: 300px;
    }

    value {
        float: right;
    }
</style>