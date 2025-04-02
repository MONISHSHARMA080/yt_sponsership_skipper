<script lang="ts">
	import { fade, fly } from 'svelte/transition';
	import { Spring } from 'svelte/motion';
	import { X, AlertCircle } from 'lucide-svelte';

	// Using $props() rune for component props
	let {
		message = 'An error occurred',
		type = 'error',
		duration = 5000, // Auto-dismiss after 5 seconds
		show = true
	} = $props<{
		message?: string;
		type?: 'error' | 'warning' | 'info';
		duration?: number;
		show?: boolean;
	}>();

	// Using $state for reactive state
	let timer = $state<ReturnType<typeof setTimeout> | null>(null);

	// Animation spring for hover effect
	const spring = new Spring({ scale: 1 });

	// Using $derived for computed values
	let bgColor = $derived(
		type === 'error' ? 'bg-red-500' : type === 'warning' ? 'bg-yellow-500' : 'bg-blue-500'
	);

	let borderColor = $derived(
		type === 'error'
			? 'border-red-600'
			: type === 'warning'
				? 'border-yellow-600'
				: 'border-blue-600'
	);

	// Using $effect for side effects
	$effect(() => {
		if (show && duration > 0) {
			timer = setTimeout(() => {
				show = false;
			}, duration);

			clearTimeout(timer);
		}

		return () => {
			if (timer) clearTimeout(timer);
		};
	});

	function handleMouseEnter() {
		spring.set({ scale: 1.02 });
		// Clear the auto-dismiss timer when hovering
		timer = setTimeout(() => {
			show = false;
		}, duration);

		clearTimeout(timer);
	}

	function handleMouseLeave() {
		spring.set({ scale: 1 });
		// Restart the timer when mouse leaves
		if (duration > 0) {
			timer = setTimeout(() => {
				show = false;
			}, duration);
		}
	}

	function close() {
		show = false;
		if (timer) clearTimeout(timer);
	}
</script>

{#if show}
	<!-- A11y: <div> with click handler must have an ARIA role -->
	<div
		class="fixed top-4 right-4 z-50 max-w-sm"
		in:fly={{ y: -20, duration: 300 }}
		out:fade={{ duration: 200 }}
		onmouseenter={handleMouseEnter}
		onmouseleave={handleMouseLeave}
		style="transform: scale({$spring.scale});"
	>
		<div
			class="relative border-4 border-black {borderColor} bg-white p-4 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]"
		>
			<!-- Type indicator -->
			<div
				class="absolute -top-4 -left-4 border-4 border-black {bgColor} px-3 py-1 font-bold text-black"
			>
				{type.charAt(0).toUpperCase() + type.slice(1)}
			</div>

			<!-- Close button -->
			<button
				class="absolute -top-4 -right-4 flex h-8 w-8 items-center justify-center border-4 border-black bg-white hover:bg-gray-100"
				onclick={close}
			>
				<X class="h-4 w-4" />
			</button>

			<!-- Content -->
			<div class="flex items-start gap-3 pt-2">
				<div class="mt-0.5">
					<AlertCircle
						class="h-5 w-5 {type === 'error'
							? 'text-red-500'
							: type === 'warning'
								? 'text-yellow-500'
								: 'text-blue-500'}"
					/>
				</div>
				<div>
					<p class="font-bold text-black">{message}</p>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	/* Animation for hover effect */
	@keyframes pulse {
		0% {
			transform: scale(1);
		}
		50% {
			transform: scale(1.02);
		}
		100% {
			transform: scale(1);
		}
	}
</style>
