<script lang="ts">
	import '../app.css';
	let { children } = $props();
</script>

<svelte:head>
	<script src="https://checkout.razorpay.com/v1/checkout.js"></script>
</svelte:head>

{@render children()}
