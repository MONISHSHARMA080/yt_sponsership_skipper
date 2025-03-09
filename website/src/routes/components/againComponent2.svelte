<script>
  import { onMount } from 'svelte';
  import { fade, scale } from 'svelte/transition';
  import { spring } from 'svelte/motion';
  import { ChevronRight, FastForward, Clock, Zap, Award, CreditCard } from 'lucide-svelte';

  let scrollY = 0;

  // Reactive animations
  const yellowCircle = spring({ x: 0, y: 0 });
  const blueCircle = spring({ x: 0, y: 0 });
  
  onMount(() => {
    const handleScroll = () => {
      scrollY = window.scrollY;
    };

    window.addEventListener("scroll", handleScroll);
    
    // Animations
    const animateShapes = () => {
      const animateYellow = () => {
        yellowCircle.set({ x: 0, y: 0 });
        setTimeout(() => yellowCircle.set({ x: 50, y: 30 }), 100);
        setTimeout(() => yellowCircle.set({ x: 0, y: 0 }), 10100);
        setTimeout(animateYellow, 20000);
      };

      const animateBlue = () => {
        blueCircle.set({ x: 0, y: 0 });
        setTimeout(() => blueCircle.set({ x: -70, y: 50 }), 100);
        setTimeout(() => blueCircle.set({ x: 0, y: 0 }), 12600);
        setTimeout(animateBlue, 25000);
      };

      animateYellow();
      animateBlue();
    };

    animateShapes();

    return () => {
      window.removeEventListener("scroll", handleScroll);
    };
  });
</script>

<svelte:head>
  <title>SkipIt - Skip the boring parts</title>
</svelte:head>

<div class="min-h-screen bg-white text-black overflow-hidden">
  <!-- Abstract geometric shapes in background -->
  <div class="fixed inset-0 -z-10 overflow-hidden">
    <div
      class="absolute top-0 left-0 w-64 h-64 bg-yellow-400 rounded-full opacity-30"
      style="transform: translate({$yellowCircle.x}px, {$yellowCircle.y}px);"
    ></div>
    <div
      class="absolute top-40 right-20 w-96 h-96 bg-blue-500 rounded-full opacity-20"
      style="transform: translate({$blueCircle.x}px, {$blueCircle.y}px);"
    ></div>
    <div
      class="absolute bottom-20 left-40 w-80 h-80 bg-red-500 opacity-20"
      style="border-radius: 60% 40% 30% 70%/60% 30% 70% 40%; animation: rotate-blob 45s linear infinite;"
    ></div>
    <div class="absolute inset-0 grid grid-cols-12 grid-rows-12 opacity-10">
      {#each Array(12) as _, rowIndex}
        {#each Array(12) as _, colIndex}
          <div
            class="border border-black {(rowIndex + colIndex) % 3 === 0 ? 'bg-purple-500' : (rowIndex + colIndex) % 3 === 1 ? 'bg-teal-400' : 'bg-transparent'}"
          ></div>
        {/each}
      {/each}
    </div>
  </div>

  <!-- Header -->
  <header class="sticky top-0 z-50 border-b-4 border-black bg-white">
    <div class="container mx-auto px-4 py-4 flex justify-between items-center">
      <div class="flex items-center gap-2">
        <div
          class="bg-red-500 w-8 h-8 rounded-full flex items-center justify-center"
          style="animation: spin 2s linear infinite;"
        >
          <FastForward class="text-white w-5 h-5" />
        </div>
        <span class="font-black text-2xl text-black tracking-tight">
          SKIP<span class="text-red-500">IT</span>
        </span>
      </div>
      <nav class="hidden md:flex gap-8 font-bold">
        <a href="#features" class="hover:text-red-500 transition-colors">
          Features
        </a>
        <a href="#pricing" class="hover:text-red-500 transition-colors">
          Pricing
        </a>
        <a href="#faq" class="hover:text-red-500 transition-colors">
          FAQ
        </a>
      </nav>
      <button
        class="bg-black text-white font-bold py-2 px-4 border-2 border-black hover:bg-white hover:text-black transition-colors transform hover:scale-105 active:scale-95"
      >
        Install Now
      </button>
    </div>
  </header>

  <!-- Hero Section -->
  <section class="relative pt-20 pb-32 overflow-hidden border-b-4 border-black">
    <div class="container mx-auto px-4">
      <div class="grid md:grid-cols-2 gap-12 items-center">
        <div in:fade={{ duration: 500, delay: 100 }}>
          <div class="mb-6">
            <span class="inline-block bg-yellow-300 text-black font-bold px-4 py-1 border-2 border-black mb-4">
              CHROME EXTENSION
            </span>
            <h1 class="text-6xl md:text-7xl font-black text-black leading-none mb-4">
                  SKIP THE <span class="text-red-500">BORING</span> PARTS
                </h1>
            <p class="text-xl mb-8">
              Automatically skip sponsorships, intros, and outros in YouTube videos. Save time and enjoy
              uninterrupted content.
            </p>
          </div>

          <div class="flex flex-col sm:flex-row gap-4">
            <button
              class="bg-red-500 text-white font-bold py-3 px-8 border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-1 hover:translate-y-1 transition-all transform hover:scale-105 active:scale-95"
            >
              Install Free
            </button>
            <a
              href="#pricing"
              class="bg-blue-500 text-white font-bold py-3 px-8 border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-1 hover:translate-y-1 transition-all flex items-center justify-center transform hover:scale-105 active:scale-95"
            >
              Go Premium <ChevronRight class="ml-1 w-5 h-5" />
            </a>
          </div>
        </div>

        <div class="relative" in:scale={{ duration: 500, delay: 200, start: 0.8 }}>
          <div class="relative z-10">
            <div class="bg-white p-2 border-4 border-black rounded-lg shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
              <div class="bg-gray-100 rounded-md p-2 mb-2">
                <div class="flex gap-2 mb-2">
                  <div class="w-3 h-3 bg-red-500 rounded-full"></div>
                  <div class="w-3 h-3 bg-yellow-500 rounded-full"></div>
                  <div class="w-3 h-3 bg-green-500 rounded-full"></div>
                </div>
                <div class="aspect-video bg-gray-800 rounded-md relative overflow-hidden">
                  <img
                    src="/placeholder.svg?height=400&width=600"
                    alt="YouTube video with sponsorship section highlighted"
                    class="w-full h-full object-cover"
                  />
                  <div
                    class="absolute bottom-4 left-4 right-4 bg-red-500 text-white font-bold py-2 px-4 rounded flex items-center justify-center"
                    style="animation: pulse 2s ease-in-out infinite;"
                  >
                    <FastForward class="mr-2" /> Sponsorship Detected - Skipping...
                  </div>
                </div>
              </div>
              <div class="flex justify-between items-center">
                <div class="font-bold">SkipIt Extension</div>
                <div class="text-green-600 font-bold flex items-center">
                  <Clock class="w-4 h-4 mr-1" /> 47s saved
                </div>
              </div>
            </div>
          </div>

          <!-- Decorative elements -->
          <div class="absolute -top-10 -right-10 w-20 h-20 bg-yellow-300 border-4 border-black z-0"></div>
          <div class="absolute -bottom-10 -left-10 w-16 h-16 bg-blue-400 border-4 border-black rounded-full z-0"></div>
          <div class="absolute top-1/2 -right-5 transform -translate-y-1/2 w-10 h-40 bg-purple-400 border-4 border-black z-0"></div>
        </div>
      </div>
    </div>

    <!-- Memphis design elements -->
    <div class="absolute bottom-0 left-0 w-full h-8 bg-black"></div>
    <div class="absolute bottom-8 left-0 w-full h-4 bg-yellow-300"></div>
    <div class="absolute -bottom-2 left-1/4 w-8 h-8 bg-blue-500 rounded-full border-4 border-black"></div>
    <div class="absolute -bottom-2 left-2/3 w-12 h-12 bg-red-500 border-4 border-black transform rotate-45"></div>
  </section>


  <!-- Features Section -->
   <div class="bg-white">
  <section id="features" class="py-20 border-b-4  bg-white  border-black relative">

    <!-- Memphis design elements -->
    <div class="absolute top-20 right-10 w-16 h-16 bg-yellow-300 border-4 border-black transform rotate-45 hidden lg:block"></div>
    <div class="absolute bottom-40 left-10 w-10 h-10 bg-blue-400 border-4 border-black rounded-full hidden lg:block"></div>
  </section>

   </div>

 
</div>



<style>
    /* cause the tailwind is not helping there  */
    /* #features {
        background-color: #ffffff;
    } */
</style>
