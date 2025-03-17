<script>
	import { FastForward, Zap, Clock, ChevronRight } from "lucide-svelte";
	import { fade } from "svelte/transition";

</script>
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

