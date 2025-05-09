<script lang="ts">
	/**
	 * @description This component is used to show an error message to the user.
	 */
	import { onDestroy } from 'svelte';
	import { fade } from 'svelte/transition';


	/**
	 * checkTheConditionAfterTheWait - is used to check if the error message should be shown or not, intended use is to be a state var .
	 */
	interface propsType {message:String, duration?:number, checkTheConditionAfterTheWait:Boolean,waitToShowError?:number }
	let { message = 'An error occurred!', duration = 3000, checkTheConditionAfterTheWait, waitToShowError= 4300  }:propsType = $props();
	let visible = $state(false);
	const timeout = setTimeout(() => (visible = false), duration);
	if (waitToShowError > duration) {
		// If the waitToShowError is greater than the duration, set the duration to waitToShowError
		duration = waitToShowError;
		console.error('waitToShowError is greater than duration, setting duration to waitToShowError');
	}

	setTimeout(()=>{visible = true; console.log(`showing the error and the visible is ${visible} and the condition to see if we should is ${checkTheConditionAfterTheWait}`);
	}, waitToShowError)



	onDestroy(() => clearTimeout(timeout));

</script>

{#if visible && checkTheConditionAfterTheWait}
	<div class="error-message top" transition:fade>
		{message}
	</div>
{/if}

<style>
	.error-message {
		position: fixed;
		padding: 1rem 2rem;
		color: white;
		border-radius: 10px;
		box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
		background-color: rgb(220, 38, 38);
		text-align: center;
		z-index: 100;
		animation: bounceIn 0.5s ease-out;
		min-width: 250px;
		max-width: 80vw;
	}
	.top {
		top: 1rem;
		left: 50%;
		transform: translateX(-50%); /* This centers the element horizontally */
	}
	@keyframes bounceIn {
		0% {
			transform: scale(0.8) translateX(-50%);
			opacity: 0;
		}
		50% {
			transform: scale(1.1) translateX(-50%);
			opacity: 1;
		}
		100% {
			transform: scale(1) translateX(-50%);
		}
	}
</style>

