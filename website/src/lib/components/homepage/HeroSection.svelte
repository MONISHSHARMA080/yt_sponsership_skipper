<script>
	import { ChevronRight, FastForward, Clock } from "lucide-svelte";
	import { fade, scale } from "svelte/transition";
	import ProgressBar from "../../../routes/components/youtubeProgressBar/progressBar.svelte";



	let inSponsorSegment = $state(false);
	let sponsorStart = $state(16);
	let sponsorEnd = $state(52);
	let videoLength = $state(73);
	let isPlaying = $state(true);






</script>
	
    
    
    
    
    
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


<style>
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
