<script lang="ts">
	import { keyFromChromeExtensionState } from '$lib/sharedState/sharedKeyState.svelte';
	import { interactWithTheChromeExtensionAndStoreIt } from '$lib/utils/interactWithChromeExtension/getKeyFromTheChromeExtension';
	import { sendChromeExtensionNewKey } from '$lib/utils/interactWithChromeExtension/sendNewKeyAfterPayment';
	import { onMount } from 'svelte';
	
	
	
	onMount(() => {
		console.log("the event is running ->");
		
		// let interactWithExtensionClass = new interactWithTheChromeExtensionAndStoreIt
		// let error = interactWithExtensionClass.start((key)=>{console.log("the key is received and it is ->",key," --- about to update the svelete store")
		// 	// keyFromChromeExtensionState.key = key
		// 	interactWithExtensionClass.cleanup()
		// // let checkKeyAndnew = new checkIfKeyIsValidAndUpdateTheState()
		// //  checkKeyAndnew.seeIfKeyIsValid(key)
		// })
		// console.log("error in interacting with the chrome extension is -> ",error );
		const sendNewKeyClass = new sendChromeExtensionNewKey("(((((((((((((((((((((((((((((((")
		sendNewKeyClass.sendKey().then((response)=>{
			console.log("the response after sending the key is  ->",response)
			console.log('\n\n\n\n');
			// sendNewKeyClass.clearEventListener()
			
		})
	});
	
</script>

<h1>Welcome to SvelteKit</h1>
{#if keyFromChromeExtensionState.key === null}
	<h2>bro the key is not there ->{keyFromChromeExtensionState.key}</h2>
{:else}
<h2> the key is -> {keyFromChromeExtensionState.key} </h2>
{/if}

{#if keyFromChromeExtensionState.isValidatedThroughBackend}
	<h1> Hello {keyFromChromeExtensionState.name} and your email is {keyFromChromeExtensionState.email}</h1>
{/if}

<p style="height:min-content;">
	the key value it {JSON.stringify(keyFromChromeExtensionState)}
</p>
