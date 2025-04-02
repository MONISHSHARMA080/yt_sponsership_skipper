<script>
	import { onDestroy } from 'svelte';
	import { fade } from 'svelte/transition';
	let { message = 'An error occurred!', duration = 3000 } = $props();
	let visible = $state(true);
	const timeout = setTimeout(() => (visible = false), duration);
	onDestroy(() => clearTimeout(timeout));
</script>

{#if visible}
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

