<script lang="ts">
	import { keyFromChromeExtensionState } from '$lib/sharedState/sharedKeyState.svelte';
	import { interactWithTheChromeExtensionAndStoreIt } from '$lib/utils/interactWithChromeExtension/getKeyFromTheChromeExtension';
	import { sendChromeExtensionNewKey } from '$lib/utils/interactWithChromeExtension/sendNewKeyAfterPayment';
	import { onMount } from 'svelte';
	import Component14 from './components/component14.svelte';
	import { checkIfKeyIsValidAndUpdateTheState } from '$lib/utils/seeIfTheKeyIsValidByBackend';
	import { razorpayOrderId } from '$lib/sharedState/razorPayKey.svelte';
	import { keyUpdatedState } from '$lib/sharedState/updatedKeyReceived.svelte';
	import getOrderIdRecursively from '$lib/utils/orderID/askForItRecursively';
	import { WriteSharedStateToStorageWhenItChanges } from '$lib/utils/SharedState/WriteSharedStateToLocalStorageOnChange.svelte';
	import { askBackendForOrderId } from '$lib/utils/razorpayIntegration/AskBackendForOrderId.svelte';
	import { shouldWeGetOrderIdRecursively } from '$lib/sharedState/getOrderIdRecursively.svelte';

	onMount(() => {
		$inspect(keyFromChromeExtensionState);

		let interactWithExtensionClass = new interactWithTheChromeExtensionAndStoreIt();
		let error = interactWithExtensionClass.start(
			// func that will run after we get the keys form the chrome extension and not form the local storage
			(a) => {
				console.log('the key we got in the func passed in the start() is', a);
				keyFromChromeExtensionState.key = a;
				interactWithExtensionClass.cleanup();
			}
		);
		if (error) {
			console.log(`the error in interacting with the extension class is ->${error}`);
		}
		console.log('about to ask backend for the order id');

		$inspect(razorpayOrderId);

		// make it is global state as I want it to be like a like a message or reactive functions that fetches when the value is changed
		// after the
		//
		// maybe I can do
		let shouldWeGetOrderIdMultipleTimes = $derived(shouldWeGetOrderIdRecursively);
		$effect(() => {
			if (
				shouldWeGetOrderIdMultipleTimes.shouldWeDoIt === true &&
				razorpayOrderId.areWeInAMiddleOfMultipleFetchCycle === false
			) {
				console.log(
					`shouldWeGetOrderIdMultipleTimes.shouldWeDoIt us true fething the order id recursively and are we still in fetching cycle ${razorpayOrderId.areWeInAMiddleOfMultipleFetchCycle}`
				);
				getOrderIdRecursively("page.svelte's effect");
				shouldWeGetOrderIdRecursively.shouldWeDoIt = false;
				console.log(`hope fully`);
			}
		});

		// WriteSharedStateToStorageWhenItChanges
		// write the keyFromChromeExtensionState in the storage when it changes
		let writeTolocalStorageInChnage = new WriteSharedStateToStorageWhenItChanges(
			keyFromChromeExtensionState,
			'KEY'
		);

		// ---- the fix ----
		// step 1) find when the key is changing using svelte runes etc and what updated them-- is it the key form the first time storage derivation(I donth think so)
		// or something else
		//
		//working theory:-> the updateKey clas is responsible see the code in it (also what is it used for, I think we can replcae it with WriteSharedStateToStorageWhenItChanges)
		// XXXXXXX
		// working theory 2-> see how we first retrive the chrome extension ->interactWithTheChromeExtensionAndStoreIt
		//
		// option 1) either change the functoTellWhenShouldYouUpdateIt fuc to detect wether the key is form the first time intitalization
		// option 2) make another shared state that is a bool and when we get the value form the local storage we will set it to true to indicate
		//           we can operate on the key
		//

		console.log(`about to watch for the key change an save it `);
		writeTolocalStorageInChnage.wathcAndSaveOnChange((sharedStateObj) => {
			if (sharedStateObj.key === '' || sharedStateObj.key === null) {
				console.log(
					`the key is either null or empty, returning false(not allowing the write),on key "KEY" and value ${sharedStateObj.key}  `
				);
				return false;
			} else {
				console.log(`about to return false and the key is ${sharedStateObj.isPaidUser}`);
				return true;
			}
		});
	});
	let keyUpdatedObj = $derived(keyUpdatedState);
	$effect(() => {
		if (keyUpdatedObj.newKeyReceived && keyFromChromeExtensionState.key) {
			// the key is updated and we are sending it to the chrome extension
			const sendNewKeyClass = new sendChromeExtensionNewKey(keyFromChromeExtensionState.key);
			sendNewKeyClass.sendKey().then(([error, didWeCorrectlySendIt]) => {
				console.log(
					`did we correctly send the new key to the chrome extension ->${didWeCorrectlySendIt}, and errors here is ->${error}`
				);
				if (error === null) {
					keyUpdatedObj.newKeyReceived = false;
				}
				sendNewKeyClass.clearEventListener();
			});
		}
	});

	$effect(() => {
		console.log(
			`the razor pay id is ->`,
			razorpayOrderId.orderIdForOnetime,
			razorpayOrderId.orderIdForRecurring
		);
	});
</script>

<Component14 />
