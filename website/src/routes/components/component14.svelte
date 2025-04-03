<script lang="ts">
	// main faq section
	import { onMount } from 'svelte';
	import { fade, scale, slide } from 'svelte/transition';
	import { Spring, Tween } from 'svelte/motion';
	import { ChevronRight, FastForward, Clock, Zap, Award, CreditCard, Loader2 } from 'lucide-svelte';
	import PremiumBenfitsSectionForPermiumUsers from './premiumBenfitsSectionForPermiumUsers/premiumBenfitsSectionForPermiumUsers.svelte';
	import { keyFromChromeExtensionState } from '$lib/sharedState/sharedKeyState.svelte';
	import HeroSection from '$lib/components/homepage/HeroSection.svelte';
	import FeatureSection from '$lib/components/homepage/FeatureSection.svelte';
	import CtaAndFooter from '$lib/components/homepage/CTAAndFooter.svelte';
	import TestimonialsAndFaqs from '$lib/components/homepage/TestimonialsAndFaqs.svelte';
	import type { RazorpayOptions } from '$lib/utils/razorpayIntegration/types/razorpayOption';
	import { razorpayOrderId } from '$lib/sharedState/razorPayKey.svelte';
	import {
		PUBLIC_CURRENCYTYPE,
		PUBLIC_ONETIMEPAYMENTPRICE,
		PUBLIC_RAZORPAY_KEY_ID,
		PUBLIC_RECURRINGPAYMENTPRICE
	} from '$env/static/public';
	import { validateCompletedPayment } from '$lib/utils/razorpayIntegration/ValidateCompletedPayment';
	import { didUserSelectOneTimePayment } from '$lib/sharedState/didUserSeletctOneTimePayment.svelte';
	import ShowUserErrorMessage from '$lib/components/showUserNotifications/ShowUserErrorMessage.svelte';
	import { askBackendForOrderId } from '$lib/utils/razorpayIntegration/AskBackendForOrderId.svelte';

	let yellowCircle = new Spring({ x: 0, y: 0 });
	const blueCircle = new Spring({ x: 0, y: 0 });

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

	onMount(() => {
		// animating the >> in the skipping the sponsor message
		springForFastForwardInVideo.target = {
			x: springForFastForwardInVideo.target.x + 12,
			y: springForFastForwardInVideo.target.y
		};

		// Animations for the circles
		const animateShapes = () => {
			const animateYellow = () => {
				yellowCircle.target = { x: 0, y: 0 };
				setTimeout(() => (yellowCircle.target = { x: 50, y: 30 }), 100);
				setTimeout(() => (yellowCircle.target = { x: 0, y: 0 }), 10100);
				setTimeout(animateYellow, 20000);
			};

			const animateBlue = () => {
				blueCircle.target = { x: 0, y: 0 };
				setTimeout(() => (blueCircle.target = { x: -70, y: 50 }), 100);
				setTimeout(() => (blueCircle.target = { x: 0, y: 0 }), 12600);
				setTimeout(animateBlue, 25000);
			};

			animateYellow();
			animateBlue();
		};

		animateShapes();

		return () => {};
	});

	async function paymentButtonClicked(textOnPaymentButton: string) {
		try {
			// Early return for the free tier
			if (textOnPaymentButton === 'Install Now') {
				console.log('Free tier selected, no payment needed');
				return; // as this is free tier and don't have to do anything on it
			}

			// Create a local variable that will be captured by the closure
			const isOneTimePayment = textOnPaymentButton === 'Try once';
			console.log(
				` is this a one time payment ->${isOneTimePayment} and the texxt on the button is ->${textOnPaymentButton}<- and ->${textOnPaymentButton === 'Try once'}`
			);

			// Also update the state for other parts of the app
			didUserSelectOneTimePayment.didUserSelectOneTimePayment = isOneTimePayment;
			didUserSelectOneTimePayment.valueChangedByMe = true;
			console.log(`Set one-time payment flag to: ${isOneTimePayment}`);

			// Get the appropriate order ID based on payment type
			let orderId = isOneTimePayment
				? razorpayOrderId.orderIdForOnetime
				: razorpayOrderId.orderIdForRecurring;

			// If order ID is not available, log error and return
			if (!orderId) {
				console.error("Order ID not found. Make sure it's fetched before payment.");
				return;
			}

			// Configure Razorpay options
			const options: RazorpayOptions = {
				key: PUBLIC_RAZORPAY_KEY_ID,
				amount: isOneTimePayment ? PUBLIC_ONETIMEPAYMENTPRICE : PUBLIC_RECURRINGPAYMENTPRICE,
				currency: PUBLIC_CURRENCYTYPE,
				name: 'Youtube Sponsorship Skipper',
				description: isOneTimePayment ? 'One-time payment' : 'Premium subscription',
				order_id: orderId,
				handler: async function (response) {
					console.log('Payment successful, validating payment...');
					console.log(`Using payment type (one-time): ${isOneTimePayment}`);

					if (
						keyFromChromeExtensionState.email === null ||
						keyFromChromeExtensionState.key === null
					) {
						return;
					}

					// Validate the payment - using the local variable instead of the state
					const validationResult = await validateCompletedPayment(
						response,
						keyFromChromeExtensionState.key,
						keyFromChromeExtensionState.email,
						isOneTimePayment // Using the local variable that's captured in the closure
					);

					if (validationResult !== null && validationResult.success) {
						console.log('Payment validation successful:', validationResult.message);
						// console.log(` the razor pay number of times before updating -> ${razorpayOrderId.numberOfTimesKeyUsed} `);
						// razorpayOrderId.numberOfTimesKeyUsed++
						// console.log(` the razor pay number of times after updating -> ${razorpayOrderId.numberOfTimesKeyUsed} `);
						askBackendForOrderId(keyFromChromeExtensionState);
						// Handle successful payment (e.g., update UI, redirect, etc.)
					} else {
						console.error('Payment validation failed:', validationResult);
						// Handle failed validation
						// updating the keys as this one is used and will not be required and if I updated the numberOfTImesKeyused++ I am stuck in a recursive loop
						askBackendForOrderId(keyFromChromeExtensionState);
						// console.log(` the razor pay number of times before updating -> ${razorpayOrderId.numberOfTimesKeyUsed} `);
						// razorpayOrderId.numberOfTimesKeyUsed++
						// console.log(` the razor pay number of times after updating -> ${razorpayOrderId.numberOfTimesKeyUsed} `);
					}
				},
				prefill: {
					email: keyFromChromeExtensionState.email || ''
				},
				theme: {
					color: '#2c15bf'
				}
			};

			// Initialize Razorpay
			const rzp = new (window as any).Razorpay(options);

			// Handle payment failures
			rzp.on('payment.failed', function (response: any) {
				console.error('Payment failed:', response.error);
				// alert(`Payment failed: ${response.error.description}`);

				// updating the razor pay state so that it can be updated for future use
				console.log(
					` the razor pay number of times before updating -> ${razorpayOrderId.numberOfTimesKeyUsed} `
				);
				razorpayOrderId.numberOfTimesKeyUsed++;
				console.log(
					` the razor pay number of times after updating -> ${razorpayOrderId.numberOfTimesKeyUsed} `
				);
			});

			// Open Razorpay payment modal
			rzp.open();
		} catch (error) {
			console.error('Error in payment button click handler:', error);
			alert('An error occurred while processing your payment. Please try again.');
		}
	}
</script>

<svelte:head>
	<title>SkipIt - Skip the boring parts</title>
	<script src="https://checkout.razorpay.com/v1/checkout.js"></script>
</svelte:head>

<div class="min-h-screen overflow-hidden bg-white text-black">
	<!-- Abstract geometric shapes in background -->
	<div class="fixed inset-0 -z-10 overflow-hidden">
		<div
			class="absolute top-0 left-0 h-64 w-64 rounded-full bg-yellow-400 opacity-30"
			style="transform: translate({yellowCircle.current.x}px, {yellowCircle.current.y}px);"
		></div>
		<div
			class="absolute top-40 right-20 h-96 w-96 rounded-full bg-blue-500 opacity-20"
			style="transform: translate({blueCircle.current.x}px, {blueCircle.current.y}px);"
		></div>
		<div
			class="absolute bottom-20 left-40 h-80 w-80 bg-red-500 opacity-20"
			style="border-radius: 60% 40% 30% 70%/60% 30% 70% 40%; animation: rotate-blob 45s linear infinite;"
		></div>
		<div class="absolute inset-0 grid grid-cols-12 grid-rows-12 opacity-10">
			{#each Array(12) as _, rowIndex}
				{#each Array(12) as _, colIndex}
					<div
						class="border border-black {(rowIndex + colIndex) % 3 === 0
							? 'bg-purple-500'
							: (rowIndex + colIndex) % 3 === 1
								? 'bg-teal-400'
								: 'bg-transparent'}"
					></div>
				{/each}
			{/each}
		</div>
	</div>

	<!-- Header -->
	<header class="sticky top-0 z-50 border-b-4 border-black bg-white">
		<div class="container mx-auto flex items-center justify-between px-4 py-4">
			<div class="flex items-center gap-2">
				<div
					class="flex h-8 w-8 items-center justify-center rounded-full bg-red-500"
					style="animation: spin 2s linear infinite;"
				>
					<FastForward class="h-5 w-5 text-white" />
				</div>
				<span class="text-2xl font-black tracking-tight text-black">
					SKIP<span class="text-red-500">IT</span>
				</span>
			</div>
			<nav class="hidden gap-8 font-bold md:flex">
				<a href="#features" class="transition-colors hover:text-red-500"> Features </a>
				<a href="#pricing" class="transition-colors hover:text-red-500"> Pricing </a>
				<a href="#faq" class="transition-colors hover:text-red-500"> FAQ </a>
			</nav>
			<button
				class="transform border-2 border-black bg-black px-4 py-2 font-bold text-white transition-colors hover:scale-105 hover:bg-white hover:text-black active:scale-95"
			>
				Install Now
			</button>
		</div>
	</header>

	<!-- Hero Section -->
	<HeroSection />

	<!-- Features Section -->
	<FeatureSection />

	{#if keyFromChromeExtensionState.isPaidUser}
		<PremiumBenfitsSectionForPermiumUsers />
	{:else}
		<!-- Pricing Section -->
		<section
			id="pricing"
			class="relative border-b-4 border-black bg-gradient-to-b from-white to-gray-100 py-20 text-black"
		>
			<div class="container mx-auto px-4">
				<div class="mb-16 text-center" in:fade={{ duration: 500 }}>
					<h2 class="mb-4 text-5xl font-black">
						CHOOSE YOUR <span class="text-purple-600">PLAN</span>
					</h2>
					<p class="mx-auto max-w-2xl text-xl">
						Upgrade to Premium for unlimited skips and advanced features.
					</p>
				</div>
				<div class="mx-auto grid max-w-7xl gap-9 md:grid-cols-3">
					{#each [{ title: 'Free', price: '$0', period: 'forever', description: 'Basic sponsorship skipping for casual YouTube viewers', features: ['Skip up to 50 sponsorships per month', 'Basic sponsorship detection', 'Time saved tracker', 'Works on all YouTube videos'], buttonText: 'Install Now', buttonColor: 'bg-black text-white', messageOnTop: '' }, { title: 'One time', price: '$3.99', period: 'once', description: 'Try premium features without recurring payments', features: ['200 premium skips', 'Advanced detection algorithm', 'Basic custom skip rules', 'Skip intros & outros', 'No recurring charges'], buttonText: 'Try once', buttonColor: 'bg-yellow-500 text-black', messageOnTop: 'Try once' }, { title: 'Recurring', price: '$4.99', period: 'per month', description: 'Unlimited skipping and advanced features for power users', features: ['Unlimited sponsorship skipping', 'Advanced detection algorithm', 'Custom skip rules and preferences', 'Skip intros, outros & reminders', 'Detailed analytics dashboard', 'Priority support'], buttonText: 'Go Premium', buttonColor: 'bg-purple-600 text-white', messageOnTop: 'Popular' }] as plan, index}
						<div
							class="relative border-4 {plan.messageOnTop
								? plan.messageOnTop === 'Popular'
									? 'border-red-600'
									: plan.messageOnTop === 'Try once'
										? 'border-yellow-500'
										: 'border-black'
								: 'border-black'} flex flex-col bg-white p-8"
							in:fade={{ duration: 500, delay: index * 100 }}
						>
							{#if plan.messageOnTop}
								<div
									class="absolute -top-4 -right-4 border-4 border-black {plan.messageOnTop ===
									'Popular'
										? 'bg-red-500'
										: plan.messageOnTop === 'Try once'
											? 'bg-yellow-400'
											: 'bg-black'} px-4 py-1 font-bold text-black"
								>
									{plan.messageOnTop}
								</div>
							{/if}
							<h3 class="mb-2 text-3xl font-bold">{plan.title}</h3>
							<div class="mb-4 flex items-end">
								<span class="text-4xl font-black">{plan.price}</span>
								<span class="ml-1 text-gray-600">/{plan.period}</span>
							</div>
							<p class="mb-6 text-gray-600">{plan.description}</p>
							<ul class="mb-8 flex-grow space-y-3">
								{#each plan.features as feature}
									<li class="flex items-start font-medium">
										<div class="mr-2 rounded-sm bg-green-500 p-1 text-white">
											<ChevronRight class="h-4 w-4" />
										</div>
										<span>{feature}</span>
									</li>
								{/each}
							</ul>

							{#if razorpayOrderId.fetchingStatus === 'fetching'}
								<button
									class=" flex w-full transform items-center justify-center gap-2 rounded-md border-3 border-black bg-gray-200 px-8 py-3 font-bold text-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] active:translate-x-[4px] active:translate-y-[4px] active:shadow-none"
									on:click={() => {
										paymentButtonClicked(plan.buttonText);
										console.log('button clicked');
									}}
								>
									<Loader2 class="h-5 w-5 animate-spin" />
									<span>{plan.buttonText}</span>
								</button>
							{:else if razorpayOrderId.fetchingStatus === 'error'}
								<div
									class="flex w-full items-center justify-center gap-2 rounded-md border-3 border-black bg-red-600 px-8 py-3 font-bold text-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]"
								>
									<!-- {console.log(`hi form the component -------------------------------`)} -->
									<ShowUserErrorMessage
										message="there is a error that prevents you form making payment, redoing it "
									/>
									<span>{plan.buttonText}</span>
									<!-- <ShowUserErrorMessage message="there is a error in getting to the backend " /> -->
								</div>
							{:else}
								<button
									class="{plan.buttonColor} flex w-full transform items-center justify-center rounded-md border-3 border-black px-8 py-3 font-bold shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] active:translate-x-[4px] active:translate-y-[4px] active:shadow-none"
									on:click={() => {
										paymentButtonClicked(plan.buttonText);
										console.log('button clicked');
									}}
								>
									<span>{plan.buttonText}</span>
								</button>
							{/if}
						</div>
					{/each}
				</div>
				<div class="mt-16 text-center" in:fade={{ duration: 500 }}>
					<p class="mb-6 text-xl">Not convinced yet? Try Premium free for 7 days!</p>
					<button
						class="transform border-2 border-black bg-gradient-to-r from-purple-600 to-blue-500 px-8 py-3 font-bold text-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all hover:translate-x-1 hover:translate-y-1 hover:scale-105 hover:shadow-none active:scale-95"
					>
						Start Free
					</button>
				</div>
			</div>
			<!-- Memphis design elements -->
			<div
				class="absolute top-40 left-10 hidden h-4 w-20 border-4 border-black bg-red-500 lg:block"
			></div>
			<div
				class="absolute right-20 bottom-20 hidden h-12 w-12 rounded-full border-4 border-black bg-blue-400 lg:block"
			></div>
		</section>
	{/if}

	<!-- Testimonials -->
	<TestimonialsAndFaqs />

	<!-- CTA Section  and the Footer-->
	<CtaAndFooter />
</div>

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
