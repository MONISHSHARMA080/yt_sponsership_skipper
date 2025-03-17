<script>
	import { CreditCard, Award, FastForward } from "lucide-svelte";
	import { scale } from "svelte/transition";

</script>
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