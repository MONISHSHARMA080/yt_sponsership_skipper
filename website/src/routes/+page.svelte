<script lang="ts">
	import { keyFromChromeExtensionState } from '$lib/sharedState/sharedKeyState.svelte';
	import { interactWithTheChromeExtensionAndStoreIt } from '$lib/utils/interactWithChromeExtension/getKeyFromTheChromeExtension';
	import { sendChromeExtensionNewKey } from '$lib/utils/interactWithChromeExtension/sendNewKeyAfterPayment';
	import { onMount } from 'svelte';
	import Component14 from './components/component14.svelte';
	import { Tween } from 'svelte/motion';
	import { cubicIn, cubicInOut, cubicOut } from 'svelte/easing';
	import { Play, FastForward, Pause } from 'lucide-svelte';
	import { fade } from 'svelte/transition';
	import ProgressBar from './components/youtubeProgressBar/progressBar.svelte';
	import { checkIfKeyIsValidAndUpdateTheState } from '$lib/utils/seeIfTheKeyIsValidByBackend';
	import { askBackendForOrderId } from '$lib/utils/razorpayIntegration/AskBackendForOrderId';
	import { razorpayOrderId } from '$lib/sharedState/razorPayKey.svelte';
	
		// Commented extension code preserved as in original
		// let interactWithExtensionClass = new interactWithTheChromeExtensionAndStoreIt
		// let error = interactWithExtensionClass.start((key)=>{console.log("the key is received and it is ->",key," --- about to update the svelete store")
		// 	// keyFromChromeExtensionState.key = key
		// 	interactWithExtensionClass.cleanup()
		// let checkKeyAndnew = new checkIfKeyIsValidAndUpdateTheState()
		// //  checkKeyAndnew.seeIfKeyIsValid(key)
		// }
		// console.log("error in interacting with the chrome extension is -> ",error );
		// const sendNewKeyClass = new sendChromeExtensionNewKey("(((((((((((((((((((((((((((((((")
		// sendNewKeyClass.sendKey().then((response)=>{
		// 	console.log("the response after sending the key is  ->",response)
		// 	console.log('\n\n\n\n');
		// 	// sendNewKeyClass.clearEventListener()
		// })


		let interactWithExtensionClass = new interactWithTheChromeExtensionAndStoreIt
		let error =interactWithExtensionClass.start(
			// func that will run after we get the keys form the chrome extension and not form the local storage
			(a)=>{console.log("the key we got in the func passed in the start() is", a);
		})
		// let error = interactWithExtensionClass.start((key)=>{console.log("the key is received and it is ->",key," --- about to update the svelete store")

		
		let val = $derived(keyFromChromeExtensionState)
		$effect(()=>{
				 askBackendForOrderId(val).then((val)=>{
				console.log(`the svelte effect returned and the value is ->`, val);
			 })
		})

		$effect(()=>{
			console.log(`the razor pay id is ->`, razorpayOrderId.orderIdForOnetime, razorpayOrderId.orderIdForRecurring);
		})


  
</script>
   <Component14 /> 

<!--    
{:else}
    <h1>Welcome to SvelteKit</h1>
    {#if keyFromChromeExtensionState.key === null}
        <h2>bro the key is not there ->{keyFromChromeExtensionState.key}</h2>
    {:else}
        <h2>the key is -> {keyFromChromeExtensionState.key}</h2>
    {/if}
    {#if keyFromChromeExtensionState.isValidatedThroughBackend}
        <h1>Hello {keyFromChromeExtensionState.name} and your email is {keyFromChromeExtensionState.email}</h1>
    {/if}
    <p >
        the key value it {JSON.stringify(keyFromChromeExtensionState)}
    </p>
{/if} -->