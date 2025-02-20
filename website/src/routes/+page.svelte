<script lang="ts">
	import { getTheKeyFromTheEvent } from '$lib/utils/getKeyFromTheChromeExtension';
	import { onMount } from 'svelte';
	
	let cleanup: (() => void) | undefined;
	
	function messageHandler(event: MessageEvent) {
		if (event.origin !== window.location.origin) return;
		
		if (event.data.type === "GET_KEY") {
			console.log("the key is  ->", event.data.key);
			if (!event.data.key) return;
			
			// Send message to remove listeners and remove our own listener
			console.log("closing all the event listeners as the key is received");
			window.postMessage({
				type: "removeAllEventListener"
			}, window.location.origin);
			
			// Remove this event listener since we got the key
			if (typeof window !== 'undefined') {
				window.removeEventListener('message', messageHandler);
			}
		}
	}
	
	onMount(() => {
		console.log("the event is running ->");
		
		try {
			// Only add listeners in browser environment
			if (typeof window !== 'undefined') {
				// Add the event listener before sending the message
				window.addEventListener('message', messageHandler);
				
				// Store cleanup function
				cleanup = () => window.removeEventListener('message', messageHandler);
				
				// Request the key
				window.postMessage({ type: 'GET_KEY' }, window.location.origin);
			}
			
		} catch (error) {
			console.log("got an error in the try catch ->", error);
		}
		
		getTheKeyFromTheEvent();

		// Return cleanup function that will run when component is destroyed
		return () => {
			if (cleanup) cleanup();
		};
	});
	
	let a = $state('-----');
</script>

<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>
<h2>I am a god {a}</h2>