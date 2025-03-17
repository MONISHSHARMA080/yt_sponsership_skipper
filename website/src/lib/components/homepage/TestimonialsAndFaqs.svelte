<script>
	import { onMount } from "svelte";

    // new faq section 
	import { cubicOut, sineInOut } from "svelte/easing";
	import  { Tween } from "svelte/motion";
	import { fade, slide } from "svelte/transition";

const titleY = new Tween(20, { duration: 500, easing: cubicOut });

	const faqItemsY = [
		new Tween(20, { duration: 500, easing: cubicOut }),
		new Tween(20, { duration: 500, easing: cubicOut }),
		new Tween(20, { duration: 500, easing: cubicOut }),
		new Tween(20, { duration: 500, easing: cubicOut }),
		new Tween(20, { duration: 500, easing: cubicOut })
	];

	let titleVisible = $state(false);
	const faqItemsVisible = $state([false, false, false, false, false]);

	// Opacity tweens
	const titleOpacity = new Tween(0, { duration: 500, easing: cubicOut });
	const faqItemsOpacity = [
		new Tween(0, { duration: 500, easing: cubicOut }),
		new Tween(0, { duration: 500, easing: cubicOut }),
		new Tween(0, { duration: 500, easing: cubicOut }),
		new Tween(0, { duration: 500, easing: cubicOut }),
		new Tween(0, { duration: 500, easing: cubicOut })
	];
// for the skipping notification on the screen
	const springForFastForwardInVideo = new Tween(
		{ x: 0, y: 0 },
		{
			easing: (t) => t * t,
			duration: 400,
			interpolate: (a, b) => (t) => ({
				x: a.x + (b.x - a.x) * t,
				y: a.y + (b.y - a.y) * t
			})
		}
	);

let facsArray = $state([
		{
			shouldWeKeepItOpen: true,
			question: 'How does the extension detect sponsorships?',
			answer:
				'Our extension uses a combination of machine learning algorithms and community-reported data to identify sponsorship segments in videos. It recognizes patterns in speech, visual cues, and common sponsorship phrases.'
		},
		{
			question: "What's the difference between Free and Premium?",
			answer:
				'The Free version allows you to skip up to 50 sponsorships per month, while Premium offers unlimited skipping, advanced detection, custom skip rules, and additional features like intro/outro skipping and detailed analytics.',
			shouldWeKeepItOpen: true
		},
		{
			shouldWeKeepItOpen: true,
			question: 'Will this slow down my browser?',
			answer:
				'No, our extension is designed to be lightweight and efficient. It runs in the background with minimal impact on your browsing experience or computer performance.'
		},
		{
			shouldWeKeepItOpen: true,
			question: 'Can I customize what gets skipped?',
			answer:
				'Yes, Premium users can set custom rules for what types of segments to skip (sponsorships, intros, outros, etc.) and even create channel-specific settings.'
		},
		{
			shouldWeKeepItOpen: true,
			question: 'How do I cancel my Premium subscription?',
			answer:
				'You can cancel your Premium subscription at any time from your account settings. Your Premium features will remain active until the end of your billing period.'
		}
	]);
onMount(() => {
		const handleScroll = () => {
			scrollY = window.scrollY;

			// Check if elements are in viewport
			const faqSection = document.getElementById('faq');
			if (faqSection) {
				const rect = faqSection.getBoundingClientRect();
				const isInViewport = rect.top < window.innerHeight * 0.75 && rect.bottom > 0;

				if (isInViewport && !titleVisible) {
					titleVisible = true;
					titleOpacity.target = 1;
					titleY.target = 0;

					// Animate FAQ items with delay
					faqItemsVisible.forEach((_, index) => {
						setTimeout(() => {
							faqItemsVisible[index] = true;
							faqItemsOpacity[index].target = 1;
							faqItemsY[index].target = 0;
						}, 100 * index);
					});
				}
			}
		};
		// animating the >> in the skipping the sponsor message
		springForFastForwardInVideo.target = {
			x: springForFastForwardInVideo.target.x + 12,
			y: springForFastForwardInVideo.target.y
		};

		// setTimeout(() => {
		// 	springForFastForwardInVideo.set({x:0, y:0})
		// }, 400)
		window.addEventListener('scroll', handleScroll);
		// Initial check in case elements are already in viewport
		handleScroll();
    })

</script>
	<section id="testimonials" class="relative border-b-4 border-black bg-white py-20">
		<div class=" container mx-auto px-4">
			<div class="mb-16 text-center" in:fade={{ duration: 500 }}>
				<h2 class="mb-4 text-5xl font-black">
					WHAT PEOPLE <span class="text-green-500">SAY</span>
				</h2>
				<p class="mx-auto max-w-2xl text-xl">
					Join thousands of happy users who save time every day.
				</p>
			</div>

			<div class="grid gap-8 md:grid-cols-3">
				{#each [{ name: 'Alex Johnson', role: 'Tech Enthusiast', quote: 'This extension has saved me hours of my life. No more sitting through boring sponsorships!', color: 'bg-red-100' }, { name: 'Sarah Miller', role: 'Daily YouTube User', quote: "The Premium version is worth every penny. I've saved over 3 hours this month alone.", color: 'bg-blue-100' }, { name: 'Michael Chen', role: 'Content Creator', quote: 'As someone who watches a lot of tutorials, this extension is a game-changer for productivity.', color: 'bg-yellow-100' }] as testimonial, index}
					<div
						class="border-4 border-black {testimonial.color} relative p-6"
						in:fade={{ duration: 500, delay: index * 100 }}
					>
						<div
							class="absolute -top-5 -left-5 flex h-10 w-10 items-center justify-center rounded-full border-4 border-black bg-white text-2xl font-bold"
						>
							"
						</div>
						<p class="mb-6 text-lg italic">{testimonial.quote}</p>
						<div class="flex items-center">
							<div class="mr-4 h-12 w-12 rounded-full border-2 border-black bg-gray-300"></div>
							<div>
								<div class="font-bold">{testimonial.name}</div>
								<div class="text-sm text-gray-600">{testimonial.role}</div>
							</div>
						</div>
					</div>
				{/each}
			</div>

			<div class="mt-16 border-4 border-black bg-green-100 p-8" in:fade={{ duration: 500 }}>
				<div class="flex flex-col items-center justify-between md:flex-row">
					<div class="mb-6 md:mb-0">
						<h3 class="mb-2 text-3xl font-bold">Ready to start skipping?</h3>
						<p class="text-xl">Join over 100,000 users saving time every day.</p>
					</div>
					<div class="flex gap-4">
						<button
							class="transform border-2 border-black bg-black px-8 py-3 font-bold text-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all hover:translate-x-1 hover:translate-y-1 hover:scale-105 hover:shadow-none active:scale-95"
						>
							Install Free
						</button>
						<button
							class="transform border-2 border-black bg-purple-600 px-8 py-3 font-bold text-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all hover:translate-x-1 hover:translate-y-1 hover:scale-105 hover:shadow-none active:scale-95"
						>
							Go Premium
						</button>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- FAQ Section -->

	<section id="faq" class="border-b-4 border-black py-20 text-black">
		<div class="container mx-auto px-4">
			<div
				class="mb-16 text-center"
				style:opacity={titleOpacity.current}
				style:transform={`translateY(${titleY.current}px)`}
			>
				<h2 class="mb-4 text-5xl font-black">
					FAQ<span class="text-orange-500">s</span>
				</h2>
				<p class="mx-auto max-w-2xl text-xl">Got questions? We've got answers.</p>
			</div>

			<div class="mx-auto max-w-3xl space-y-6">
				{#each facsArray as faq, index}
					<!-- svelte-ignore a11y_click_events_have_key_events -->
					<!-- svelte-ignore a11y_no_static_element_interactions -->
					<div
						class="overflow-hidden border-4 border-black bg-white"
						style:opacity={faqItemsOpacity[index].current}
						style:transform={`translateY(${faqItemsY[index].current}px)`}
						onclick={() => {
							// close the question in the array at that index
							// console.log(`shouldWeKeepItOpen: ${faq.shouldWeKeepItOpen}`);
							facsArray[index].shouldWeKeepItOpen = !facsArray[index].shouldWeKeepItOpen;
							// console.log(`shouldWeKeepItOpen: ${faq.shouldWeKeepItOpen}`);
						}}
					>
						<div
							class="flex items-center justify-between border-b-4 border-black bg-gray-100 p-4 text-lg font-bold"
						>
							{faq.question}
							<div class="flex h-6 w-6 items-center justify-center bg-black text-white">
								{faq.shouldWeKeepItOpen ? '−' : '+'}
							</div>
						</div>

						{#if faq.shouldWeKeepItOpen}
							<div transition:slide={{ duration: 198, easing: sineInOut }} class="p-4">
								{faq.answer}
							</div>
						{/if}
					</div>
				{/each}
			</div>
		</div>
	</section>


<style>
	.body {
		font-family: 'Inter', sans-serif;
		background-color: white;
	}

	.section {
		background-color: white;
	}

	@keyframes spin {
		from {
			transform: rotate(0deg);
		}
		to {
			transform: rotate(360deg);
		}
	}

	@keyframes rotate-blob {
		from {
			transform: rotate(0deg);
		}
		to {
			transform: rotate(360deg);
		}
	}

	@keyframes rotate-blob {
		from {
			transform: rotate(0deg);
		}
		to {
			transform: rotate(360deg);
		}
	}

	@keyframes spin {
		from {
			transform: rotate(0deg);
		}
		to {
			transform: rotate(360deg);
		}
	}

	@keyframes pulse {
		0% {
			opacity: 0.8;
			transform: scale(1);
		}
		50% {
			opacity: 1;
			transform: scale(1.05);
		}
		100% {
			opacity: 0.8;
			transform: scale(1);
		}
	}

	@keyframes moveBackForth {
		0%,
		100% {
			transform: translateX(0px);
		}
		50% {
			transform: translateX(9px);
		}
	}

	.back-forth {
		animation: moveBackForth 1s ease-in-out infinite;
	}

	@keyframes slideUp {
		from {
			transform: translateY(50px);
			opacity: 0;
		}
		to {
			transform: translateY(0);
			opacity: 1;
		}
	}

	@keyframes slideDown {
		from {
			transform: translateY(0);
			opacity: 1;
		}
		to {
			transform: translateY(50px);
			opacity: 0;
		}
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	@keyframes slideInRight {
		from {
			transform: translateX(-100px);
			opacity: 0;
		}
		to {
			transform: translateX(0);
			opacity: 1;
		}
	}

	@keyframes moveBackForth {
		0%,
		100% {
			transform: translateX(0px);
		}
		50% {
			transform: translateX(10px);
		}
	}

	/* Animation classes */
	.animate-slide-up {
		animation: slideUp 0.59s cubic-bezier(0.16, 1, 0.3, 1) forwards;
	}

	.animate-slide-down {
		animation: slideDown 0.79s cubic-bezier(0.16, 1, 0.3, 1) forwards;
	}

	.animate-fade-in {
		animation: fadeIn 0.3s ease-out 0.2s both;
	}

	.animate-delayed-fade {
		animation: fadeIn 0.38s ease-out 0.43s both;
	}

	.animate-slide-in-right {
		animation: slideInRight 0.6s cubic-bezier(0.16, 1, 0.3, 1) 0.2s forwards;
	}

	.back-forth {
		animation: moveBackForth 1s ease-in-out infinite;
	}
</style>
