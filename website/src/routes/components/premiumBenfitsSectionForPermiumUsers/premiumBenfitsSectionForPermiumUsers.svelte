<script lang="ts">
  import { onMount } from 'svelte';
  import { ChevronRight, FastForward, Clock, Zap, Award, CreditCard } from 'lucide-svelte';
  
  // Using $state for reactive variables in SvelteKit 5
  let scrollY = $state(0);
  let isSkipping = $state(false);
  let headerInView = $state(false);
  let mainCardInView = $state(false);
  let supportInView = $state(false);
  
  // References to DOM elements
  let headerRef: HTMLElement;
  let mainCardRef: HTMLElement;
  let supportRef: HTMLElement;
  
  onMount(() => {
    // Banner animation timer
    const skipTimer = setTimeout(() => {
      isSkipping = true;
      setTimeout(() => isSkipping = false, 3000); // Banner disappears after 3 seconds

    }, 2000); // Banner appears after 2 seconds
    
    // Scroll event listener
    const handleScroll = () => {
      scrollY = window.scrollY;
    };
    
    window.addEventListener("scroll", handleScroll);
    
    // Intersection Observers for animation triggers
    const observerHeader = new IntersectionObserver(
      ([entry]) => {
        headerInView = entry.isIntersecting;
      },
      { threshold: 0.1 },
    );
    
    const observerMainCard = new IntersectionObserver(
      ([entry]) => {
        mainCardInView = entry.isIntersecting;
      },
      { threshold: 0.1 },
    );
    
    const observerSupport = new IntersectionObserver(
      ([entry]) => {
        supportInView = entry.isIntersecting;
      },
      { threshold: 0.1 },
    );
    
    if (headerRef) observerHeader.observe(headerRef);
    if (mainCardRef) observerMainCard.observe(mainCardRef);
    if (supportRef) observerSupport.observe(supportRef);
    
    // Cleanup function
    return () => {
      clearTimeout(skipTimer);
      window.removeEventListener("scroll", handleScroll);
      if (headerRef) observerHeader.unobserve(headerRef);
      if (mainCardRef) observerMainCard.unobserve(mainCardRef);
      if (supportRef) observerSupport.unobserve(supportRef);
    };
  });
</script>

<section id="premium" class="py-20 border-b-4 border-black relative bg-gradient-to-b from-white to-gray-100">
  <div class="container mx-auto px-4">
    <div
      bind:this={headerRef}
      class="text-center mb-16 transition-all duration-500 transform {
        headerInView ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-5'
      }"
    >
      <span class="inline-block bg-yellow-400 text-black font-bold px-4 py-1 border-2 border-black mb-4">
        PREMIUM ACTIVATED
      </span>
      <h2 class="text-5xl font-black mb-4">
        YOUR <span class="text-purple-600">PREMIUM</span> BENEFITS
      </h2>
      <p class="text-xl max-w-2xl mx-auto">
        You're all set with unlimited skipping power. Simply browse YouTube and let SkipIt do the rest.
      </p>
    </div>

    <div class="max-w-4xl mx-auto">
      <div
        bind:this={mainCardRef}
        class="border-4 border-black bg-white p-8 relative transition-all duration-500 transform {
          mainCardInView ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-5'
        }"
      >
        <div class="absolute -top-4 -right-4 bg-purple-600 text-white font-bold py-1 px-4 border-4 border-black">
          PREMIUM
        </div>

        <div class="flex items-center gap-3 mb-8">
          <div class="bg-purple-600 p-3 border-2 border-black">
            <Award class="w-8 h-8 text-white" />
          </div>
          <h3 class="text-3xl font-bold">Thanks for being Premium</h3>
        </div>

        <p class="text-xl mb-8">
          Your YouTube experience has been elevated. Just watch YouTube as you normally would — SkipIt Premium
          works silently in the background, automatically detecting and skipping past sponsorships.
        </p>

        <div class="grid md:grid-cols-2 gap-8 mb-8">
          <div class="border-2 border-black p-6 bg-gray-50">
            <h4 class="text-xl font-bold mb-3 border-b-2 border-black pb-2">What You Have Access To</h4>
            <ul class="space-y-3">
              <li class="flex items-start">
                <div class="bg-green-500 text-white p-1 mr-2 mt-1">
                  <ChevronRight class="w-4 h-4" />
                </div>
                <span>
                  <span class="font-bold">Unlimited skipping</span> with no monthly limits
                </span>
              </li>
              <li class="flex items-start">
                <div class="bg-green-500 text-white p-1 mr-2 mt-1">
                  <ChevronRight class="w-4 h-4" />
                </div>
                <span>
                  <span class="font-bold">Advanced detection algorithm</span> with higher accuracy
                </span>
              </li>
              <li class="flex items-start">
                <div class="bg-green-500 text-white p-1 mr-2 mt-1">
                  <ChevronRight class="w-4 h-4" />
                </div>
                <span>
                  <span class="font-bold">Custom skip rules</span> and channel preferences
                </span>
              </li>
              <li class="flex items-start">
                <div class="bg-green-500 text-white p-1 mr-2 mt-1">
                  <ChevronRight class="w-4 h-4" />
                </div>
                <span>
                  <span class="font-bold">Skip intros, outros & reminders</span> automatically
                </span>
              </li>
            </ul>
          </div>

          <div class="border-2 border-black p-6 bg-purple-50">
            <h4 class="text-xl font-bold mb-3 border-b-2 border-black pb-2">How To Use It</h4>
            <ul class="space-y-4">
              <li class="flex items-start">
                <div class="bg-purple-600 text-white p-1 mr-3 rounded-full w-6 h-6 flex items-center justify-center mt-1">
                  <span class="font-bold">1</span>
                </div>
                <p>
                  Just <span class="font-bold">open YouTube</span> in your browser
                </p>
              </li>
              <li class="flex items-start">
                <div class="bg-purple-600 text-white p-1 mr-3 rounded-full w-6 h-6 flex items-center justify-center mt-1">
                  <span class="font-bold">2</span>
                </div>
                <p>
                  Watch your <span class="font-bold">favorite videos</span> as usual
                </p>
              </li>
              <li class="flex items-start">
                <div class="bg-purple-600 text-white p-1 mr-3 rounded-full w-6 h-6 flex items-center justify-center mt-1">
                  <span class="font-bold">3</span>
                </div>
                <p>
                  SkipIt <span class="font-bold">automatically detects and skips</span> sponsored segments
                </p>
              </li>
              <li class="flex items-start">
                <div class="bg-purple-600 text-white p-1 mr-3 rounded-full w-6 h-6 flex items-center justify-center mt-1">
                  <span class="font-bold">4</span>
                </div>
                <p>
                  Enjoy <span class="font-bold">uninterrupted viewing</span> without lifting a finger
                </p>
              </li>
            </ul>
          </div>
        </div>

        <div class="border-2 border-black p-6 bg-yellow-50 mb-8">
          <h4 class="text-xl font-bold mb-3">Your Time Savings</h4>
          <div class="flex flex-col md:flex-row justify-between items-center">
            <div class="flex items-center mb-4 md:mb-0">
              <Clock class="w-10 h-10 text-purple-600 mr-3" />
              <div>
                <p class="text-sm font-bold">This month you've saved</p>
                <p class="text-4xl font-black">2.7 hours</p>
              </div>
            </div>
            <div class="h-16 w-px bg-black hidden md:block"></div>
            <div class="flex items-center">
              <FastForward class="w-10 h-10 text-red-500 mr-3" />
              <div>
                <p class="text-sm font-bold">Total sponsorships skipped</p>
                <p class="text-4xl font-black">123</p>
              </div>
            </div>
          </div>
        </div>

        <div class="group bg-black text-white font-bold py-3 px-8 w-full border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-1 hover:translate-y-1 transition-all flex items-center justify-center">
          <a href="https://youtube.com" class="text-white no-underline flex items-center">
            Start Watching YouTube Now
            <ChevronRight class="ml-2 w-5 h-5" />
          </a>
        </div>
      </div>
    </div>

    <div
      bind:this={supportRef}
      class="mt-16 p-8 border-4 border-black bg-purple-100 transition-all duration-500 transform {
        supportInView ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-5'
      }"
    >
      <div class="flex flex-col md:flex-row items-center justify-between">
        <div class="mb-6 md:mb-0">
          <h3 class="text-3xl font-bold mb-2">Need help with your Premium account?</h3>
          <p class="text-xl">Our support team is here to assist you 24/7.</p>
        </div>
        <button class="bg-purple-600 text-white font-bold py-3 px-8 border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-1 hover:translate-y-1 transition-all group-hover:scale-105 active:scale-95">
          Contact Support
        </button>
      </div>
    </div>
  </div>

  <!-- Memphis design elements -->
  <div class="absolute top-40 left-10 w-20 h-4 bg-red-500 border-4 border-black hidden lg:block"></div>
  <div class="absolute bottom-20 right-20 w-12 h-12 bg-blue-400 border-4 border-black rounded-full hidden lg:block"></div>
</section>
