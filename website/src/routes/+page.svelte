<script lang="ts">
	import { getTheKeyFromTheEvent } from '$lib/utils/getKeyFromTheChromeExtension';
	import {onMount} from 'svelte';

	
	onMount(()=>{
		console.log("the event is running  ->");
	let b = 2;
	try {
		window.postMessage({type: 'GET_KEY'}, window.location.origin);
		 console.log("hi ", b);
		window.addEventListener('message', (event) => {
		console.log("\n\n--+++++---message received in svelte and it is ->",event,"----+++----\n\n\n");
		if (event.origin !== window.location.origin) return;
      
      if (event.data.type === 'KEY_RESPONSE') {
        var key = event.data.key;
      }
    });


		 
	} catch (error) {
		console.log("got a error in the try catch ->",error);
		}
		getTheKeyFromTheEvent()
	})

	let a = $state('-----');
    
</script>

<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>
<h2>I am a god {a}</h2>
