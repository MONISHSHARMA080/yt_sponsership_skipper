<script lang="ts">
	import { keyFromChromeExtensionState } from '$lib/sharedState/sharedKeyState.svelte';
	import { interactWithTheChromeExtensionAndStoreItInTheState } from '$lib/utils/getKeyFromTheChromeExtension';
	import { onMount } from 'svelte';
	
	
	
	onMount(() => {
		console.log("the event is running ->");
		
		let interactWithExtensionClass = new interactWithTheChromeExtensionAndStoreItInTheState
		let error = interactWithExtensionClass.start((key)=>{console.log("the key is received and it is ->",key," --- about to update the svelete store")
			keyFromChromeExtensionState.key = key
		})
		console.log("error in interacting with the chrome extension is -> ",error );
		

	});
	
</script>

<h1>Welcome to SvelteKit</h1>
{#if keyFromChromeExtensionState.key === null}
	<h2>bro the key is not there ->{keyFromChromeExtensionState.key}</h2>
{:else}
<h2> the key is -> {keyFromChromeExtensionState.key} </h2>
{/if}
