<script lang="ts">
	import { interactWithTheChromeExtension } from '$lib/utils/getKeyFromTheChromeExtension';
	import { onMount } from 'svelte';
	
	
	
	onMount(() => {
		console.log("the event is running ->");
		
		let interactWithExtensionClass = new interactWithTheChromeExtension
		let error = interactWithExtensionClass.start()
		console.log("error in interacting with the chrome extension is -> ",error );
		

	});
	
	let a = $state('-----');
</script>

<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>
<h2>I am a god {a}</h2>