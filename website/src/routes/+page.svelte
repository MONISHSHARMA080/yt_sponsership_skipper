<script lang="ts">
	import { keyFromChromeExtensionState } from '$lib/sharedState/sharedKeyState.svelte';
	import { interactWithTheChromeExtensionAndStoreIt } from '$lib/utils/interactWithChromeExtension/getKeyFromTheChromeExtension';
	import { sendChromeExtensionNewKey } from '$lib/utils/interactWithChromeExtension/sendNewKeyAfterPayment';
	import PageDesign from './pageDesign.svelte';
	import Component2 from './components/component2.svelte';
	import Component3 from './components/component3.svelte';
	import Component4 from './components/component4.svelte';
	import Component from './components/component.svelte';
	import Component5 from './components/component5.svelte';
	import Component6 from './components/component6.svelte';
	import Component7 from './components/component7.svelte';
	import Component8 from './components/component8.svelte';
	import Component9 from './components/component9.svelte';
	import OriginalComponent1 from './components/originalComponent1.svelte';
	import { onMount } from 'svelte';
	import Component10 from './components/component10.svelte';
	import Component11 from './components/component11.svelte';
	import Component12 from './components/component12.svelte';
	
	let change = $state(false);
	let component = $state(0);
	
	onMount(() => {
		console.log("the event is running ->");
		setTimeout(() => {
			change = true;
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
	function nextComponent() {
		console.log("in the next component function", component);
		
		let nextComponent = component + 1;	
		if (nextComponent >= componentArray.length) {
			component = 0;
		}
		else {
			component = nextComponent;
		}
		
		console.log("the next component is ->", component);
	}

	function previousComponent() {
		console.log("in the previous component function", component);
		
		let nextNum = component - 1;
		if (nextNum < 0) {
			component = componentArray.length - 1;
		}
		else {
			component = nextNum;
		}
		
		console.log("the previous component is ->", component);
	}
	
	// Array of component references
	const componentArray = [
		Component,// implement video demo (ui/graphics) like other components here
		Component2, // implement video demo (ui/graphics) like other components here 
		Component3,
		Component4,
		Component5,
		Component6,
		Component7,
		Component8,
		Component9,
	];
	let autoIncrementCompoent  = () => {
	
		setTimeout(() => {
			nextComponent()
			autoIncrementCompoent()
		}, 4000)
	}
	// autoIncrementCompoent()
	
</script>

{#if change}
 

<!-- <Component2 /> -->
 <Component2 />
   
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
    <p style="height:min-content;">
        the key value it {JSON.stringify(keyFromChromeExtensionState)}
    </p>
{/if}