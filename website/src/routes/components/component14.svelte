<script>
	import { onMount } from 'svelte';
	import { fade, scale, slide } from 'svelte/transition';
	import { Spring, Tween } from 'svelte/motion';
	import { ChevronRight, FastForward, Clock, Zap, Award, CreditCard } from 'lucide-svelte';
	import { cubicOut, sineIn, sineInOut } from 'svelte/easing';
	import ProgressBar from './youtubeProgressBar/progressBar.svelte';
	import PremiumBenfitsSectionForPermiumUsers from './premiumBenfitsSectionForPermiumUsers/premiumBenfitsSectionForPermiumUsers.svelte';
	import { keyFromChromeExtensionState } from '$lib/sharedState/sharedKeyState.svelte';

	let scrollY = $state(0);
	let yellowCircle = new Spring({ x: 0, y: 0 });
	const blueCircle = new Spring({ x: 0, y: 0 });

	let inSponsorSegment = $state(false);
	let shouldWeStartRemovingSponsorSkippingMessage = $state(false);

	// Initialize element visibility states
	let titleVisible = $state(false);
	const faqItemsVisible = $state([false, false, false, false, false]);

	// Create tweens for Y position animations
	const titleY = new Tween(20, { duration: 500, easing: cubicOut });

	const faqItemsY = [
		new Tween(20, { duration: 500, easing: cubicOut }),
		new Tween(20, { duration: 500, easing: cubicOut }),
		new Tween(20, { duration: 500, easing: cubicOut }),
		new Tween(20, { duration: 500, easing: cubicOut }),
		new Tween(20, { duration: 500, easing: cubicOut })
	];

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

		return () => {
			window.removeEventListener('scroll', handleScroll);
		};
	});
	// after on sight
	/**
	 * Creates a YouTube-style progress bar as an HTML element.
	 * @param {string} baseColorForTailwindProgressBar - The Tailwind color for the filled progress bar.
	 * @param {string} baseColorForBarWhenNotUse - The Tailwind color for the unfilled part of the progress bar.
	 * @param {number} currentStateOfTheProgressBar - The progress percentage (0-100).
	 * @returns {HTMLDivElement} - The progress bar element.
	 */
	function factoryYoutubeProgressBar(
		baseColorForTailwindProgressBar,
		baseColorForBarWhenNotUse,
		currentStateOfTheProgressBar
	) {
		// Create the main progress bar container
		const progressBar = document.createElement('div');
		progressBar.className = `relative w-full h-4 ${baseColorForBarWhenNotUse} rounded overflow-hidden`;

		// Create the filled part of the progress bar
		const progressFill = document.createElement('div');
		progressFill.className = `absolute top-0 left-0 h-full ${baseColorForTailwindProgressBar}`;
		progressFill.style.width = `${currentStateOfTheProgressBar}%`;

		// Append the fill to the main container
		progressBar.appendChild(progressFill);

		return progressBar;
	}

	let isPlaying = $state(true);
	let progress = $state(0);

	// Define the sponsor segment (in real app, this would come from an API or database)
	let sponsorStart = $state(16);
	let sponsorEnd = $state(52);
	let videoLength = $state(73);
</script>

<svelte:head>
	<title>SkipIt - Skip the boring parts</title>
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
	<section class="relative overflow-hidden border-b-4 border-black pt-20 pb-32">
		<div class="container mx-auto px-4">
			<div class="grid items-center gap-12 md:grid-cols-2">
				<div in:fade={{ duration: 500, delay: 100 }}>
					<div class="mb-6">
						<span
							class="mb-4 inline-block border-2 border-black bg-yellow-300 px-4 py-1 font-bold text-black"
						>
							CHROME EXTENSION
						</span>
						<h1 class="mb-4 text-6xl leading-none font-black text-black md:text-7xl">
							SKIP THE <span class="text-red-500">BORING</span> PARTS
						</h1>
						<p class="mb-8 text-xl">
							Automatically skip sponsorships, intros, and outros in YouTube videos. Save time and
							enjoy uninterrupted content.
						</p>
					</div>

					<div class="flex flex-col gap-4 sm:flex-row">
						<button
							class="transform border-2 border-black bg-red-500 px-8 py-3 font-bold text-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all hover:translate-x-1 hover:translate-y-1 hover:scale-105 hover:shadow-none active:scale-95"
						>
							Install Free
						</button>
						<a
							href="#pricing"
							class="flex transform items-center justify-center border-2 border-black bg-blue-500 px-8 py-3 font-bold text-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all hover:translate-x-1 hover:translate-y-1 hover:scale-105 hover:shadow-none active:scale-95"
						>
							Go Premium <ChevronRight class="ml-1 h-5 w-5" />
						</a>
					</div>
				</div>

				<div class="relative" in:scale={{ duration: 500, delay: 200, start: 0.8 }}>
					<div class="relative z-10">
						<div
							class="rounded-lg border-4 border-black bg-white p-2 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]"
						>
							<div class="mb-2 rounded-md bg-gray-100 p-2">
								<div class="mb-2 flex gap-2">
									<div class="h-3 w-3 rounded-full bg-red-500"></div>
									<div class="h-3 w-3 rounded-full bg-yellow-500"></div>
									<div class="h-3 w-3 rounded-full bg-green-500"></div>
								</div>
								<div class="relative aspect-video overflow-hidden rounded-md bg-gray-800">
									<img
										src="/placeholder.svg?height=400&width=600"
										alt="YouTube video with sponsorship section highlighted"
										class="h-full w-full object-cover"
									/>
									<!-- <div
										class="absolute right-4 bottom-4 left-4 flex items-center justify-center rounded bg-red-500 px-4 py-2 font-bold text-white"
										style="animation: pulse 1s ease-in-out infinite;"
									>
										<FastForward class="mr-2" /> Sponsorship Detected - Skipping...
									</div> -->

									{#if inSponsorSegment}
										<div class="animate-slide-up absolute right-4 bottom-8 left-4 z-50">
											<div
												class="flex items-center justify-between overflow-hidden rounded-sm border-4 border-black bg-white p-2 shadow-lg"
											>
												<div class="animate-fade-in flex items-center">
													<div class="mr-2 rounded-sm bg-red-500 p-1">
														<FastForward class="h-4 w-4 text-white" />
													</div>
													<span class="animate-delayed-fade text-sm font-bold text-black">
														Upcomming sponsership detected
													</span>
												</div>
												<div class="animate-slide-in-right flex items-center">
													<span
														class="mr-2 border-2 border-black bg-yellow-400 px-2 py-0.5 text-xs font-bold"
													>
														SKIPPING
													</span>
													<div class="back-forth">
														<FastForward class="h-4 w-4 text-black" />
													</div>
												</div>
											</div>
										</div>
									{:else}
										<!-- if the  inSponsorSegment is false then try to slide it down -->
										<!-- When inSponsorSegment is false, apply slide-down animation -->
										<div class="animate-slide-down absolute right-4 bottom-6 left-4 z-50">
											<div
												class="flex items-center justify-between overflow-hidden rounded-sm border-4 border-black bg-white p-2 shadow-lg"
											>
												<div class="flex items-center">
													<div class="mr-2 rounded-sm bg-red-500 p-1">
														<FastForward class="h-4 w-4 text-white" />
													</div>
													<span class="text-sm font-bold text-black"> SPONSORSHIP DETECTED </span>
												</div>
												<div class="flex items-center">
													<span
														class="mr-2 border-2 border-black bg-yellow-400 px-2 py-0.5 text-xs font-bold"
													>
														SKIPPING
													</span>
													<div class="back-forth">
														<FastForward class="h-4 w-4 text-black" />
													</div>
												</div>
											</div>
										</div>
									{/if}

									<ProgressBar
										funcToRunWhenInTheSponSorSection={(areWeInSponsorSegment) => {}}
										funToRunFewSecBeforeSponsorSegment={{
											time: 6,
											func: () => {
												// console.log("running the function before the sponsor segment")
												inSponsorSegment = true;
											}
										}}
										{sponsorStart}
										{sponsorEnd}
										{videoLength}
										playVideo={isPlaying}
										funcToRunAfterVideoCompletion={() => {
											//console.log("video ended and stopping it ")
										}}
										funcToRunAfterTheSponsorSegment={() => {
											// console.log(" sponsor segment ended and stopping it ")
											//   isPlaying = true
											// here make a set timeout so that it goes in a slow fashion
											setTimeout(() => {
												inSponsorSegment = false;
											}, 1090);
										}}
										sponsorShipDetectedFastForward={true}
									/>

									<!-- the youtube video progress bar -->
									<!-- <div class="absolute bottom-0 left-0 h-4 rounded-3xl my w-full bg-blue-800" > -->
									<!-- the video progress bar  -->
									<!-- <div class="absolute bottom-0  h-4 rounded-3xl my  bg-red-600" style="left:20% ;width:70%" ></div> -->
									<!-- <div class="h-full bg-green-500" style="width: 50%"></div> -->
									<div></div>
								</div>
							</div>
							<div class="flex items-center justify-between">
								<div class="font-bold">SkipIt Extension</div>
								<div class="flex items-center font-bold text-green-600">
									<Clock class="mr-1 h-4 w-4" /> 47s saved
								</div>
							</div>
						</div>
					</div>

					<!-- Decorative elements -->
					<div
						class="absolute -top-10 -right-10 z-0 h-20 w-20 border-4 border-black bg-yellow-300"
					></div>
					<div
						class="absolute -bottom-10 -left-10 z-0 h-16 w-16 rounded-full border-4 border-black bg-blue-400"
					></div>
					<div
						class="absolute top-1/2 -right-5 z-0 h-40 w-10 -translate-y-1/2 transform border-4 border-black bg-purple-400"
					></div>
				</div>
			</div>
		</div>

		<!-- Memphis design elements -->
		<div class="absolute bottom-0 left-0 h-8 w-full bg-black"></div>
		<div class="absolute bottom-8 left-0 h-4 w-full bg-yellow-300"></div>
		<div
			class="absolute -bottom-2 left-1/4 h-8 w-8 rounded-full border-4 border-black bg-blue-500"
		></div>
		<div
			class="absolute -bottom-2 left-2/3 h-12 w-12 rotate-45 transform border-4 border-black bg-red-500"
		></div>
	</section>

	<!-- Features Section -->
	<div class="bg-white">
		<section id="features" class="relative border-b-4 border-black py-20">
			<div class="container mx-auto bg-white px-4">
				<div class="mb-16 text-center" in:fade={{ duration: 500 }}>
					<h2 class="mb-4 text-5xl font-black text-black">
						AWESOME <span class="text-blue-500">FEATURES</span>
					</h2>
					<p class="mx-auto max-w-2xl text-xl">
						Our extension is packed with powerful features to enhance your YouTube experience.
					</p>
				</div>

				<div class="grid gap-8 md:grid-cols-3">
					{#each [{ icon: FastForward, title: 'Auto-Skip Sponsorships', description: "Automatically detects and skips sponsored segments in videos so you don't have to manually skip them.", color: 'bg-red-500' }, { icon: Zap, title: 'Lightning Fast', description: 'Minimal impact on performance. Works silently in the background without slowing down your browsing.', color: 'bg-yellow-400' }, { icon: Clock, title: 'Time Saved Tracker', description: "See exactly how much time you've saved by skipping sponsorships across all your watched videos.", color: 'bg-blue-500' }] as feature, index}
						<div
							class="group relative overflow-hidden border-4 border-black bg-white p-6"
							in:fade={{ duration: 500, delay: index * 100 }}
						>
							<div
								class="absolute top-0 right-0 h-20 w-20 {feature.color} -mt-10 -mr-10 border-b-4 border-l-4 border-black transition-all duration-300 group-hover:mt-0 group-hover:mr-0"
							></div>
							<div class="relative z-20">
								<div class="{feature.color} mb-4 inline-block border-2 border-black p-4 text-black">
									<feature.icon class="h-10 w-10" />
								</div>
								<h3 class="mb-2 text-2xl font-bold text-black">{feature.title}</h3>
								<p class="text-black">{feature.description}</p>
							</div>
						</div>
					{/each}
				</div>

				<div
					class="relative mt-16 border-4 border-black bg-purple-100 p-8"
					in:fade={{ duration: 500 }}
				>
					<div class="absolute -top-5 -left-5 h-10 w-10 border-4 border-black bg-yellow-300"></div>
					<div
						class="absolute -right-5 -bottom-5 h-10 w-10 rounded-full border-4 border-black bg-blue-400"
					></div>

					<div class="grid items-center gap-8 md:grid-cols-2">
						<div>
							<h3 class="mb-4 text-3xl font-bold">Smart Detection Technology</h3>
							<p class="mb-4">
								Our advanced algorithm recognizes sponsorship segments with incredible accuracy,
								even when creators try to disguise them.
							</p>
							<ul class="space-y-2">
								{#each ['Recognizes common sponsorship phrases', 'Detects visual sponsorship indicators', 'Learns from user feedback', 'Updates in real-time'] as item}
									<li class="flex items-start">
										<div class="mt-1 mr-2 bg-green-500 p-1 text-white">
											<ChevronRight class="h-4 w-4" />
										</div>
										{item}
									</li>
								{/each}
							</ul>
						</div>
						<div class="relative">
							<div class="relative border-4 border-black bg-white p-6">
								<h4 class="mb-3 border-b-2 border-black pb-2 text-2xl font-bold">How It Works</h4>
								<div class="space-y-4">
									<div class="flex items-start">
										<div class="mt-1 mr-3 bg-red-500 p-1 text-white">
											<span class="font-bold">1</span>
										</div>
										<div>
											<h5 class="text-lg font-bold">Audio Analysis</h5>
											<p>
												Our AI listens for key phrases like "this video is sponsored by" or "use
												code" that indicate sponsorship content.
											</p>
										</div>
									</div>

									<div class="flex items-start">
										<div class="mt-1 mr-3 bg-yellow-400 p-1 text-black">
											<span class="font-bold">2</span>
										</div>
										<div>
											<h5 class="text-lg font-bold">Visual Detection</h5>
											<p>
												We identify sponsor logos, product placements, and visual cues that creators
												use during sponsored segments.
											</p>
										</div>
									</div>

									<div class="flex items-start">
										<div class="mt-1 mr-3 bg-blue-500 p-1 text-white">
											<span class="font-bold">3</span>
										</div>
										<div>
											<h5 class="text-lg font-bold">Community Data</h5>
											<p>
												Our database is constantly updated with user-reported sponsorship segments
												across millions of videos.
											</p>
										</div>
									</div>

									<div class="flex items-start">
										<div class="mt-1 mr-3 bg-green-500 p-1 text-white">
											<span class="font-bold">4</span>
										</div>
										<div>
											<h5 class="text-lg font-bold">Instant Skipping</h5>
											<p>
												When a sponsorship is detected, we automatically skip ahead to the relevant
												content - no manual interaction needed!
											</p>
										</div>
									</div>
								</div>
							</div>
							<div
								class="absolute -bottom-4 -left-4 -z-10 h-full w-full border-4 border-black bg-red-200"
							></div>
						</div>
					</div>
				</div>
			</div>

			<!-- Memphis design elements -->
			<div
				class="absolute top-20 right-10 hidden h-16 w-16 rotate-45 transform border-4 border-black bg-yellow-300 lg:block"
			></div>
			<div
				class="absolute bottom-40 left-10 hidden h-10 w-10 rounded-full border-4 border-black bg-blue-400 lg:block"
			></div>
		</section>
	</div>

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
    {#each [
        { 
            title: 'Free', 
            price: '$0', 
            period: 'forever', 
            description: 'Basic sponsorship skipping for casual YouTube viewers', 
            features: [
                'Skip up to 50 sponsorships per month', 
                'Basic sponsorship detection', 
                'Time saved tracker', 
                'Works on all YouTube videos'
            ], 
            buttonText: 'Install Now', 
            buttonColor: 'bg-black text-white', 
            messageOnTop: ""
        },
        { 
            title: 'One time', 
            price: '$3.99', 
            period: 'once', 
            description: 'Try premium features without recurring payments', 
            features: [
                '200 premium skips', 
                'Advanced detection algorithm', 
                'Basic custom skip rules', 
                'Skip intros & outros',
                'No recurring charges'
            ], 
            buttonText: 'Try once', 
            buttonColor: 'bg-yellow-500 text-black', 
            messageOnTop: "Try once"
        },
        { 
            title: 'Recurring', 
            price: '$4.99', 
            period: 'per month', 
            description: 'Unlimited skipping and advanced features for power users', 
            features: [
                'Unlimited sponsorship skipping', 
                'Advanced detection algorithm', 
                'Custom skip rules and preferences', 
                'Skip intros, outros & reminders', 
                'Detailed analytics dashboard', 
                'Priority support'
            ], 
            buttonText: 'Go Premium', 
            buttonColor: 'bg-purple-600 text-white', 
            messageOnTop: "Popular"
        }
    ] as plan, index}

        <div
            class="relative border-4 {plan.messageOnTop ? (plan.messageOnTop === 'Popular' ? 'border-red-500' : (plan.messageOnTop === 'Try once' ? 'border-yellow-500' : 'border-black')) : 'border-black'} flex flex-col bg-white p-8"
            in:fade={{ duration: 500, delay: index * 100 }}
        >
            {#if plan.messageOnTop}
                <div
                    class="absolute -top-4 -right-4 border-4 border-black {plan.messageOnTop === 'Popular' ? 'bg-red-500' : (plan.messageOnTop === 'Try once' ? 'bg-yellow-400' : 'bg-black')} px-4 py-1 font-bold text-black"
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
									<li class="flex  font-medium items-start">
										<div class="mr-2 bg-green-500 p-1 text-white">
											<ChevronRight class="h-4 w-4" />
										</div>
										{feature}
									</li>
								{/each}
							</ul>
							<button
								class="{plan.buttonColor} w-full transform border-3 border-black px-8 py-3 font-bold  shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all hover:translate-x-1 hover:translate-y-1 hover:scale-105 hover:shadow-none active:scale-95"
							>
								{plan.buttonText}
							</button>
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

	<!-- CTA Section -->

	<section id="cta" class="relative overflow-hidden bg-white py-20">
		<div class=" relative z-10 container mx-auto px-4">
			<div
				class="mx-auto max-w-4xl border-4 border-black bg-white p-10 text-center shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]"
				in:scale={{ duration: 500, start: 0.9 }}
			>
				<h2 class="mb-6 text-5xl font-black">
					STOP WASTING <span class="text-red-500">TIME</span>
				</h2>
				<p class="mb-8 text-xl">
					The average YouTube user wastes over 5 hours per month watching sponsorships. Get that
					time back with SkipIt!
				</p>

				<div class="flex flex-col justify-center gap-4 sm:flex-row">
					<button
						class="transform border-2 border-black bg-black px-8 py-3 font-bold text-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all hover:translate-x-1 hover:translate-y-1 hover:scale-105 hover:shadow-none active:scale-95"
					>
						Install Free
					</button>
					<button
						class="flex transform items-center justify-center border-2 border-black bg-purple-600 px-8 py-3 font-bold text-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all hover:translate-x-1 hover:translate-y-1 hover:scale-105 hover:shadow-none active:scale-95"
					>
						Go Premium <CreditCard class="ml-2 h-5 w-5" />
					</button>
				</div>

				<div class="mt-8 flex items-center justify-center gap-4">
					<div class="flex">
						{#each Array(5) as _}
							<Award class="h-6 w-6 text-yellow-500" />
						{/each}
					</div>
					<span class="font-bold">4.9/5 from 2,000+ reviews</span>
				</div>
			</div>
		</div>

		<!-- Memphis design background -->
		<div class="absolute inset-0 -z-10">
			<div
				class="absolute top-10 left-10 h-20 w-20 rotate-45 transform border-4 border-black bg-yellow-300"
			></div>
			<div
				class="absolute top-40 right-20 h-16 w-16 rounded-full border-4 border-black bg-blue-400"
			></div>
			<div class="absolute bottom-20 left-1/3 h-8 w-24 border-4 border-black bg-red-500"></div>
			<div
				class="absolute right-1/4 bottom-40 h-12 w-12 rotate-12 transform border-4 border-black bg-purple-400"
			></div>
		</div>
	</section>

	<!-- Footer -->
	<footer class="border-t-4 border-white bg-black py-12 text-white">
		<div class="container mx-auto px-4">
			<div class="flex flex-wrap items-center">
				<div class="mr-auto flex items-center gap-2">
					<div class="flex h-8 w-9 items-center justify-center rounded-full bg-red-500">
						<FastForward class="h-5 w-5 text-white" />
					</div>
					<span class="text-2xl font-black tracking-tight">
						SKIP<span class="text-red-500">IT</span>
					</span>
				</div>
				<div class="mt-4 ml-8 md:mt-0 md:ml-8">
					<p class="text-gray-400">
						Save time and enjoy uninterrupted YouTube content. Don't waste time getting stuck on the
						youtube sponsership
					</p>
				</div>
			</div>
		</div>
	</footer>
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
