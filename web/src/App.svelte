<script>
	import Stats from './Stats.svelte';

	export let name;
</script>

<main>
	<mainhead>
		<h1>Hello {name}!</h1>
	</mainhead>	
	<Stats />
</main>

<style>
	main {
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	mainhead {
		text-align: center;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>