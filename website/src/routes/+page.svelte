<script lang="ts">
	import { keyFromChromeExtensionState } from '$lib/sharedState/sharedKeyState.svelte';
	import { interactWithTheChromeExtensionAndStoreIt } from '$lib/utils/interactWithChromeExtension/getKeyFromTheChromeExtension';
	import { sendChromeExtensionNewKey } from '$lib/utils/interactWithChromeExtension/sendNewKeyAfterPayment';
	import { onMount } from 'svelte';
	import Component14 from './components/component14.svelte';
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


		
		onMount(()=>{
			let interactWithExtensionClass = new interactWithTheChromeExtensionAndStoreIt
		let error =interactWithExtensionClass.start(
			// func that will run after we get the keys form the chrome extension and not form the local storage
			(a)=>{console.log("the key we got in the func passed in the start() is", a);
			keyFromChromeExtensionState.key = a
			interactWithExtensionClass.cleanup()
		})
		if (error) {
			console.log(`the error in interacting with the extension class is ->${error}`);
		}
		console.log("about to ask backend for the order id");
		
			 askBackendForOrderId(keyFromChromeExtensionState).then((val)=>{
				console.log(`the order id form the backend's value is ->${val}`);
				
			 })


		}
	)
		// let error = interactWithExtensionClass.start((key)=>{console.log("the key is received and it is ->",key," --- about to update the svelete store")

		
		let val = $derived(keyFromChromeExtensionState)
		$effect(()=>{
			console.log(`asking backend for the order id and the key is -> ${keyFromChromeExtensionState.key}`);
				 askBackendForOrderId(val).then((val)=>{
				console.log(`the svelte effect returned and the value is ->`, val);
			 })
		})

		// if we don't do this and the ordered id is used in a failure , it will not update, also update it only when the failure is registered
		// let newRazorPayOrderID = $derived(razorpayOrderId)
		// $effect(()=>{
		// 	console.log(`the razor pay order id after change is -> ${JSON.stringify(newRazorPayOrderID)}`);
		// 	if (newRazorPayOrderID.numberOfTimesKeyUsed === 1 || newRazorPayOrderID.numberOfTimesKeyUsed >1) {
		// 		console.log(`the razor pay ordered id has been used one time`);
		// 		askBackendForOrderId(keyFromChromeExtensionState).then((val)=>{
		// 			console.log(`the razorpay order ID is used one time and we are updated it and the returned value is  -> ${val} `);
		// 		})
		// 	}else{
		// 		return
		// 	}
		// })

		$effect(()=>{
			console.log(`the razor pay id is ->`, razorpayOrderId.orderIdForOnetime, razorpayOrderId.orderIdForRecurring);
		})


  
</script>
   <Component14 /> 
