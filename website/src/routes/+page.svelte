<script lang="ts">
	import { keyFromChromeExtensionState } from '$lib/sharedState/sharedKeyState.svelte';
	import { interactWithTheChromeExtensionAndStoreIt } from '$lib/utils/interactWithChromeExtension/getKeyFromTheChromeExtension';
	import { sendChromeExtensionNewKey } from '$lib/utils/interactWithChromeExtension/sendNewKeyAfterPayment';
	import { onMount } from 'svelte';
	import Component14 from './components/component14.svelte';
	
	let change = $state(false);
	
	onMount(() => {
		setTimeout(() => {
			change = true;
		console.log("the event is running ->");
		}, 100);
		
		// Commented extension code preserved as in original
		// let interactWithExtensionClass = new interactWithTheChromeExtensionAndStoreIt
		// let error = interactWithExtensionClass.start((key)=>{console.log("the key is received and it is ->",key," --- about to update the svelete store")
		// 	// keyFromChromeExtensionState.key = key
		// 	interactWithExtensionClass.cleanup()
		// // let checkKeyAndnew = new checkIfKeyIsValidAndUpdateTheState()
		// //  checkKeyAndnew.seeIfKeyIsValid(key)
		// }
		// console.log("error in interacting with the chrome extension is -> ",error );
		// const sendNewKeyClass = new sendChromeExtensionNewKey("(((((((((((((((((((((((((((((((")
		// sendNewKeyClass.sendKey().then((response)=>{
		// 	console.log("the response after sending the key is  ->",response)
		// 	console.log('\n\n\n\n');
		// 	// sendNewKeyClass.clearEventListener()
		// })
	});
	
	// Fixed function to cycle through components
	// function nextComponent() {
	// 	console.log("in the next component function", component);
		
	// 	let nextComponent = component + 1;	
	// 	if (nextComponent >= componentArray.length) {
	// 		component = 0;
	// 	}
	// 	else {
	// 		component = nextComponent;
	// 	}
		
	// 	console.log("the next component is ->", component);
	// }

	// function previousComponent() {
	// 	console.log("in the previous component function", component);
		
	// 	let nextNum = component - 1;
	// 	if (nextNum < 0) {
	// 		component = componentArray.length - 1;
	// 	}
	// 	else {
	// 		component = nextNum;
	// 	}
		
	// 	console.log("the previous component is ->", component);
	// }
	
	// Array of component references
	// const componentArray = [
	// 	Component,// implement video demo (ui/graphics) like other components here
	// 	Component2, // implement video demo (ui/graphics) like other components here 
	// 	Component3,
	// 	Component4,
	// 	Component5,
	// 	Component6,
	// 	Component7,
	// 	Component8,
	// 	Component9,
	// ];
	// let autoIncrementCompoent  = () => {
	
	// 	setTimeout(() => {
	// 		nextComponent()
	// 		autoIncrementCompoent()
	// 	}, 4000)
	// }
	// autoIncrementCompoent()
	
</script>

{#if change}
 

<!-- <Component2 /> -->
 <Component14/>
  <!-- <AgainComponent2/> -->
   
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
{/if}