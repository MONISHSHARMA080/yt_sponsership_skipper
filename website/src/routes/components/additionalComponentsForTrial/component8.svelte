


<script lang="ts">
    import { onMount } from 'svelte';
  import { ChevronRight, Play, FastForward, Zap, Shield, Clock } from 'lucide-svelte';
  import { fade, fly } from 'svelte/transition';
  import { spring, tweened } from 'svelte/motion';
  import { cubicInOut } from 'svelte/easing';
  
  
  let scrollY = $state(0);
  let mounted = $state(false);
  
  const videoProgress = tweened(0, {
    duration: 2000,
    easing: cubicInOut
  });
  
  // Create inViewport action
  function inViewport(node: Element, options:any = {}) {
    const { callback = () => {}, threshold = 0.1 } = options;
    
    let observer: IntersectionObserver;
    
    function handleIntersect(entries: any[]) {
      const entry = entries[0];
      if (entry.isIntersecting) {
        callback(true);
        // Add animation classes or trigger animations here
        node.classList.add('in-viewport');
      } else {
        callback(false);
      }
    }
    
    function setupObserver() {
      if (observer) observer.disconnect();
      
      observer = new IntersectionObserver(handleIntersect, {
        rootMargin: '0px',
        threshold
      });
      
      observer.observe(node);
    }
    
    setupObserver();
    
    return {
      update(newOptions: any) {
        options = { ...options, ...newOptions };
        setupObserver();
      },
      destroy() {
        if (observer) observer.disconnect();
      }
    };
  }
  
  onMount(() => {
    mounted = true;
    
    const handleScroll = () => {
      scrollY = window.scrollY;
    };
    
    window.addEventListener("scroll", handleScroll);
    
    // Start the video progress animation after a delay
    setTimeout(() => {
      videoProgress.set(0.6);
    }, 1500);
    
    return () => window.removeEventListener("scroll", handleScroll);
  });
  
  const features = [
    {
      icon: Zap,
      title: "Automatic Detection",
      description: "Our AI automatically identifies sponsored segments in videos without any manual input required.",
    },
    {
      icon: FastForward,
      title: "Instant Skipping",
      description: "Skip straight past sponsored content with zero interaction needed - it just works.",
    },
    {
      icon: Shield,
      title: "Privacy Focused",
      description: "We don't track your browsing history or collect any personal data. Your privacy matters.",
    },
    {
      icon: Clock,
      title: "Time Saved Tracker",
      description: "See exactly how much time you've saved by skipping sponsored content across all videos.",
    },
    {
      icon: Play,
      title: "Custom Controls",
      description: "Fine-tune how the extension works with customizable settings for your perfect experience.",
    },
    {
      icon: 'custom',
      title: "Cross-Platform",
      description: "Works seamlessly across Chrome, Firefox, and Edge browsers for a consistent experience.",
    }
  ];
  
  const steps = [
    {
      number: "01",
      title: "Install the Extension",
      description: "Add SkipSpot to your browser with just one click. No account required.",
    },
    {
      number: "02",
      title: "Watch YouTube Normally",
      description: "Continue using YouTube as you always do. Our extension works silently in the background.",
    },
    {
      number: "03",
      title: "Enjoy Sponsor-Free Videos",
      description: "SkipSpot automatically detects and skips past sponsored segments in real-time.",
    }
  ];
  
  const testimonials = [
    {
      quote: "This extension has saved me hours of my life. No more sitting through boring sponsorships!",
      name: "Alex K.",
      title: "Daily YouTube User",
    },
    {
      quote: "The Pro version is absolutely worth it. Unlimited skips and the time-saved stats are eye-opening.",
      name: "Sarah M.",
      title: "Content Creator",
    },
    {
      quote: "I was skeptical at first, but now I can't watch YouTube without it. The detection is incredibly accurate.",
      name: "Michael T.",
      title: "Tech Enthusiast",
    }
  ];
</script>

<svelte:head>
  <style>
    @keyframes float1 {
      0%, 100% { transform: translate(0, 0); }
      50% { transform: translate(20px, -30px); }
    }
    @keyframes float2 {
      0%, 100% { transform: translate(0, 0); }
      50% { transform: translate(-20px, 20px); }
    }
    @keyframes float3 {
      0%, 100% { transform: translate(0, 0); }
      50% { transform: translate(15px, 25px); }
    }
    @keyframes skipNotification {
      0% { opacity: 0; transform: translate(-50%, -50%) scale(0.8); }
      10%, 90% { opacity: 1; transform: translate(-50%, -50%) scale(1); }
      100% { opacity: 0; transform: translate(-50%, -50%) scale(0.8); }
    }
    
    .in-viewport {
      animation: fadeIn 0.8s forwards;
    }
    
    @keyframes fadeIn {
      from { opacity: 0; transform: translateY(20px); }
      to { opacity: 1; transform: translateY(0); }
    }
  </style>
</svelte:head>


<div class="min-h-screen bg-black text-white overflow-hidden">
  <!-- Floating geometric shapes -->
  <div class="fixed inset-0 pointer-events-none z-0">
    <div class="absolute top-1/4 left-1/4 w-64 h-64 rounded-full bg-purple-500/20 blur-3xl"
      style="animation: float1 8s ease-in-out infinite;">
    </div>
    <div class="absolute bottom-1/3 right-1/3 w-96 h-96 rounded-full bg-cyan-500/20 blur-3xl"
      style="animation: float2 10s ease-in-out infinite;">
    </div>
    <div class="absolute top-2/3 right-1/4 w-72 h-72 rounded-full bg-yellow-500/20 blur-3xl"
      style="animation: float3 9s ease-in-out infinite;">
    </div>
  </div>

  <!-- Header -->
  <header class="fixed top-0 left-0 right-0 z-50 bg-black/80 backdrop-blur-md border-b border-white/10">
    <div class="container mx-auto px-4 py-4 flex justify-between items-center">
      {#if mounted}
        <div in:fly={{ x: -20, duration: 500 }} class="flex items-center gap-2">
          <FastForward class="h-8 w-8 text-yellow-400" />
          <span class="text-2xl font-bold tracking-tight">SkipSpot</span>
        </div>
        <nav in:fly={{ y: -10, duration: 500, delay: 200 }} class="hidden md:flex items-center gap-8">
          <a href="#features" class="hover:text-yellow-400 transition-colors">Features</a>
          <a href="#how-it-works" class="hover:text-yellow-400 transition-colors">How it works</a>
          <a href="#pricing" class="hover:text-yellow-400 transition-colors">Pricing</a>
        </nav>
        <div in:fly={{ x: 20, duration: 500, delay: 400 }}>
          <button class="bg-yellow-400 text-black hover:bg-yellow-300 font-bold px-4 py-2 rounded-md">
            Get Extension
          </button>
        </div>
      {/if}
    </div>
  </header>

  <!-- Hero Section -->
  <section class="relative pt-32 pb-20 md:pt-40 md:pb-32">
    <div class="container mx-auto px-4">
      <div class="max-w-4xl mx-auto text-center">
        {#if mounted}
          <div in:fly={{ y: 20, duration: 700 }} class="mb-6">
            <span class="inline-block px-3 py-1 bg-yellow-400 text-black font-medium rounded-full text-sm mb-4">
              SKIP THE BORING PARTS
            </span>
          </div>
          <h1 in:fly={{ y: 20, duration: 700, delay: 200 }} class="text-5xl md:text-7xl font-extrabold mb-6 leading-tight">
            Watch YouTube <span class="text-yellow-400">Without</span> The Sponsorships
          </h1>
          <p in:fly={{ y: 20, duration: 700, delay: 400 }} class="text-xl text-gray-300 mb-10 max-w-2xl mx-auto">
            Our smart AI automatically detects and skips sponsored segments in YouTube videos, saving you time and improving your viewing experience.
          </p>
          <div in:fly={{ y: 20, duration: 700, delay: 600 }} class="flex flex-col sm:flex-row gap-4 justify-center">
            <button class="bg-yellow-400 text-black hover:bg-yellow-300 font-bold text-lg px-8 py-6 rounded-md">
              Install Free Version
            </button>
            <button class="border border-yellow-400 text-yellow-400 hover:bg-yellow-400/10 font-bold text-lg px-8 py-6 rounded-md flex items-center justify-center">
              Upgrade to Pro
              <ChevronRight class="ml-2 h-5 w-5" />
            </button>
          </div>
        {/if}
      </div>
    </div>

    <!-- Animated video player mockup -->
    {#if mounted}
      <div in:fly={{ y: 40, duration: 800, delay: 800 }} class="max-w-4xl mx-auto mt-20 relative">
        <div class="relative rounded-xl overflow-hidden border-2 border-white/20 shadow-2xl">
          <div class="bg-zinc-900 h-12 flex items-center px-4 gap-2">
            <div class="flex gap-2">
              <div class="w-3 h-3 rounded-full bg-red-500"></div>
              <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
              <div class="w-3 h-3 rounded-full bg-green-500"></div>
            </div>
            <div class="ml-4 bg-zinc-800 rounded px-2 py-1 text-xs flex-1 text-center text-gray-400">
              youtube.com/watch?v=example
            </div>
          </div>
          <div class="aspect-video bg-zinc-800 relative">
            <div class="absolute inset-0 flex items-center justify-center">
              <img src="/placeholder.svg?height=400&width=720" alt="YouTube video player" class="w-full h-full object-cover" />
              {#if mounted}
                <div in:fade={{ delay: 1200, duration: 500 }} class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent p-4">
                  <div class="flex items-center gap-4">
                    <Play class="h-8 w-8" />
                    <div class="h-2 bg-gray-700 rounded-full flex-1 overflow-hidden">
                      <div class="h-full bg-red-600" style={`width: ${$videoProgress * 100}%`}></div>
                    </div>
                    <span class="text-sm">10:24</span>
                  </div>
                </div>
              {/if}
              <div
                class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-yellow-400 text-black font-bold px-6 py-3 rounded-lg flex items-center gap-2 shadow-lg"
                style="animation: skipNotification 2s 2s forwards;">
                <FastForward class="h-5 w-5" />
                Sponsorship Skipped!
              </div>
            </div>
          </div>
        </div>
      </div>
    {/if}
  </section>

  <!-- Features Section -->
  <section id="features" class="py-20 bg-zinc-900">
    <div class="container mx-auto px-4">
      <div use:inViewport class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-4">Powerful Features</h2>
        <p class="text-xl text-gray-400 max-w-2xl mx-auto">
          Our extension comes packed with features to enhance your YouTube experience
        </p>
      </div>
      <div class="grid md:grid-cols-3 gap-8">
        {#each features as feature, index}
          <div
            use:inViewport={{ callback: (_inView: any) => {} }}
            class="bg-zinc-800 p-6 rounded-xl border border-white/10 hover:border-yellow-400/50 transition-colors group">
            <div class="mb-4">
              {#if feature.icon === 'custom'}
                <svg class="h-10 w-10 text-yellow-400" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 4V20M12 4L8 8M12 4L16 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              {:else}
                <div class="h-10 w-10 text-yellow-400">
                {#if feature.icon === Zap}
                    <Zap />
                {:else if feature.icon === FastForward}
                    <FastForward />
                {:else if feature.icon === Shield}
                    <Shield />
                {:else if feature.icon === Clock}
                    <Clock />
                {:else if feature.icon === Play}
                    <Play />
                {/if}
                </div>
              {/if}
            </div>
            <h3 class="text-xl font-bold mb-2 group-hover:text-yellow-400 transition-colors">
              {feature.title}
            </h3>
            <p class="text-gray-400">{feature.description}</p>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- How It Works -->
  <section id="how-it-works" class="py-20 bg-black">
    <div class="container mx-auto px-4">
      <div use:inViewport class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-4">How It Works</h2>
        <p class="text-xl text-gray-400 max-w-2xl mx-auto">
          SkipSpot uses advanced AI to detect and skip sponsored content automatically
        </p>
      </div>
      <div class="grid md:grid-cols-3 gap-8 max-w-4xl mx-auto">
        {#each steps as step, index}
          <div use:inViewport={{ callback: () => {} }} class="relative">
            <div class="text-8xl font-bold text-yellow-400/20 absolute -top-10 left-0">{step.number}</div>
            <div class="pt-8 pl-4 relative">
              <h3 class="text-xl font-bold mb-2">{step.title}</h3>
              <p class="text-gray-400">{step.description}</p>
            </div>
            {#if index < 2}
              <div use:inViewport={{ callback: (_inView:any) => {} }} class="hidden md:block absolute top-1/2 right-0 transform translate-x-1/2">
                <ChevronRight class="h-8 w-8 text-yellow-400" />
              </div>
            {/if}
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- Pricing Section -->
  <section id="pricing" class="py-20 bg-zinc-900 relative overflow-hidden">
    <div class="absolute inset-0 opacity-30 pointer-events-none">
      <div class="absolute top-0 left-0 w-full h-full bg-[radial-gradient(circle_at_30%_20%,rgba(255,204,0,0.2),transparent_40%)]"></div>
      <div class="absolute bottom-0 right-0 w-full h-full bg-[radial-gradient(circle_at_70%_80%,rgba(255,204,0,0.2),transparent_40%)]"></div>
    </div>
    <div class="container mx-auto px-4 relative z-10">
      <div use:inViewport class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-4">Choose Your Plan</h2>
        <p class="text-xl text-gray-400 max-w-2xl mx-auto">
          Start with our free plan or upgrade for unlimited skipping
        </p>
      </div>
      <div class="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto">
        <div use:inViewport={{ callback: (_inView:any) => {}, oneshot: true }} class="bg-zinc-800 rounded-xl p-8 border border-white/10">
          <div class="mb-6">
            <h3 class="text-2xl font-bold mb-2">Free</h3>
            <p class="text-gray-400">Perfect for casual YouTube viewers</p>
          </div>
          <div class="mb-6">
            <span class="text-4xl font-bold">$0</span>
            <span class="text-gray-400 ml-2">forever</span>
          </div>
          <ul class="space-y-3 mb-8">
            <li class="flex items-start gap-2">
              <svg class="h-5 w-5 text-yellow-400 mt-0.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
              <span>Skip up to 30 sponsored segments per day</span>
            </li>
            <li class="flex items-start gap-2">
              <svg class="h-5 w-5 text-yellow-400 mt-0.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
              <span>Basic detection algorithm</span>
            </li>
            <li class="flex items-start gap-2">
              <svg class="h-5 w-5 text-yellow-400 mt-0.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
              <span>Standard support</span>
            </li>
            <li class="flex items-start gap-2 text-gray-500">
              <svg class="h-5 w-5 text-gray-500 mt-0.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
              <span>No time saved statistics</span>
            </li>
            <li class="flex items-start gap-2 text-gray-500">
              <svg class="h-5 w-5 text-gray-500 mt-0.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
              <span>No customization options</span>
            </li>
          </ul>
            <button class="w-full bg-zinc-700 hover:bg-zinc-600 text-white px-4 py-2 rounded-md transition-colors duration-300">Install Free Version</button>
        </div>
        <div use:inViewport={{ callback: (_inView:any) => {}, oneshot: true, delay: 200 }} class="bg-gradient-to-br from-yellow-400/10 to-yellow-600/10 rounded-xl p-8 border border-yellow-400/30 relative">
          <div class="absolute -top-4 right-4">
            <span class="bg-yellow-400 text-black font-bold px-3 py-1 rounded-full text-sm">MOST POPULAR</span>
          </div>
          <div class="mb-6">
            <h3 class="text-2xl font-bold mb-2">Pro</h3>
            <p class="text-gray-400">For serious YouTube enthusiasts</p>
          </div>
          <div class="mb-6">
            <span class="text-4xl font-bold">$4.99</span>
            <span class="text-gray-400 ml-2">per month</span>
          </div>
          <ul class="space-y-3 mb-8">
            <li class="flex items-start gap-2">
              <svg class="h-5 w-5 text-yellow-400 mt-0.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
              <span><strong>Unlimited</strong> sponsored segment skipping</span>
            </li>
            <li class="flex items-start gap-2">
              <svg class="h-5 w-5 text-yellow-400 mt-0.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
              <span>Advanced AI detection (99% accuracy)</span>
            </li>
            <li class="flex items-start gap-2">
              <svg class="h-5 w-5 text-yellow-400 mt-0.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
              <span>Priority support</span>
            </li>
            <li class="flex items-start gap-2">
              <svg class="h-5 w-5 text-yellow-400 mt-0.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
              <span>Detailed time-saved statistics</span>
            </li>
            <li class="flex items-start gap-2">
              <svg class="h-5 w-5 text-yellow-400 mt-0.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
              <span>Full customization options</span>
            </li>
            <li class="flex items-start gap-2">
              <svg class="h-5 w-5 text-yellow-400 mt-0.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
              <span>Early access to new features</span>
            </li>
          </ul>
          <button class="w-full bg-yellow-400 hover:bg-yellow-300 text-black font-bold px-4 py-2 rounded-md transition-colors duration-300">Upgrade to Pro</button>
        </div>
      </div>
    </div>
  </section>

  <!-- Testimonials -->
  <section class="py-20 bg-black">
    <div class="container mx-auto px-4">
      <div use:inViewport class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-4">What Users Say</h2>
        <p class="text-xl text-gray-400 max-w-2xl mx-auto">
          Join thousands of satisfied users who save time every day
        </p>
      </div>
      <div class="grid md:grid-cols-3 gap-8">
        {#each testimonials as testimonial, index}
          <div use:inViewport={{ callback: (_inView:any) => {}, delay: index * 100 }} class="bg-zinc-800 p-6 rounded-xl border border-white/10">
            <div class="mb-4 text-yellow-400">
              {#each Array(5) as _, i}
                <span class="text-xl">★</span>
              {/each}
            </div>
            <p class="text-gray-300 mb-6 italic">"{testimonial.quote}"</p>
            <div>
              <p class="font-bold">{testimonial.name}</p>
              <p class="text-gray-400 text-sm">{testimonial.title}</p>
            </div>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- CTA Section -->
  <section class="py-20 bg-yellow-400 text-black">
    <div class="container mx-auto px-4">
      <div class="max-w-4xl mx-auto text-center">
        <h2 use:inViewport class="text-4xl md:text-5xl font-bold mb-6">
          Ready to Take Back Your Time?
        </h2>
        <p use:inViewport={{ delay: 200 }} class="text-xl mb-8 max-w-2xl mx-auto">
          Join thousands of users who save hours every month by skipping sponsored content automatically.
        </p>
        <div use:inViewport={{ delay: 400 }} class="flex flex-col sm:flex-row gap-4 justify-center">
          <button class="bg-black text-white hover:bg-zinc-800 font-bold text-lg px-8 py-6 rounded-md">
            Install Free Version
          </button>
          <button class="border border-black text-black hover:bg-black/10 font-bold text-lg px-8 py-6 rounded-md flex items-center justify-center">
            Upgrade to Pro
            <ChevronRight class="ml-2 h-5 w-5" />
          </button>
        </div>
      </div>
    </div>
  </section>

  <!-- Footer -->
<footer class="bg-zinc-900 text-white py-12">
  <div class="container mx-auto px-4">
    <div class="grid md:grid-cols-4 gap-8">
      <div>
        <div class="flex items-center gap-2 mb-4">
          <FastForward class="h-6 w-6 text-yellow-400" />
          <span class="text-xl font-bold">SkipSpot</span>
        </div>
        <p class="text-gray-400 mb-4">The smartest way to skip sponsored content on YouTube.</p>
        <div class="flex gap-4">
          <!-- svelte-ignore a11y_consider_explicit_label -->
          <a href="#" class="text-gray-400 hover:text-yellow-400 transition-colors">
            <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M24 4.557c-.883.392-1.832.656-2.828.775 1.017-.609 1.798-1.574 2.165-2.724-.951.564-2.005.974-3.127 1.195-.897-.957-2.178-1.555-3.594-1.555-3.179 0-5.515 2.966-4.797 6.045-4.091-.205-7.719-2.165-10.148-5.144-1.29 2.213-.669 5.108 1.523 6.574-.806-.026-1.566-.247-2.229-.616-.054 2.281 1.581 4.415 3.949 4.89-.693.188-1.452.232-2.224.084.626 1.956 2.444 3.379 4.6 3.419-2.07 1.623-4.678 2.348-7.29 2.04 2.179 1.397 4.768 2.212 7.548 2.212 9.142 0 14.307-7.721 13.995-14.646.962-.695 1.797-1.562 2.457-2.549z"></path>
            </svg>
          </a>
          <!-- svelte-ignore a11y_consider_explicit_label -->
          <a href="#" class="text-gray-400 hover:text-yellow-400 transition-colors">
            <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z"></path>
            </svg>
          </a>
          <!-- svelte-ignore a11y_consider_explicit_label -->
          <a href="#" class="text-gray-400 hover:text-yellow-400 transition-colors">
            <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M9 8h-3v4h3v12h5v-12h3.642l.358-4h-4v-1.667c0-.955.192-1.333 1.115-1.333h2.885v-5h-3.808c-3.596 0-5.192 1.583-5.192 4.615v3.385z"></path>
            </svg>
          </a>
        </div>
      </div>
      <div>
        <h3 class="font-bold text-lg mb-4">Product</h3>
        <ul class="space-y-2">
          <li>
            <a href="#" class="text-gray-400 hover:text-yellow-400 transition-colors">
              Features
            </a>
          </li>
          <li>
            <a href="#" class="text-gray-400 hover:text-yellow-400 transition-colors">
              Pricing
            </a>
          </li>
          <li>
            <a href="#" class="text-gray-400 hover:text-yellow-400 transition-colors">
              FAQ
            </a>
          </li>
          <li>
            <a href="#" class="text-gray-400 hover:text-yellow-400 transition-colors">
              Download
            </a>
          </li>
        </ul>
      </div>
      <div>
        <h3 class="font-bold text-lg mb-4">Company</h3>
        <ul class="space-y-2">
          <li>
            <a href="#" class="text-gray-400 hover:text-yellow-400 transition-colors">
              About
            </a>
          </li>
          <li>
            <a href="#" class="text-gray-400 hover:text-yellow-400 transition-colors">
              Blog
            </a>
          </li>
          <li>
            <a href="#" class="text-gray-400 hover:text-yellow-400 transition-colors">
              Careers
            </a>
          </li>
          <li>
            <a href="#" class="text-gray-400 hover:text-yellow-400 transition-colors">
              Contact
            </a>
          </li>
        </ul>
      </div>
      <div>
        <h3 class="font-bold text-lg mb-4">Legal</h3>
        <ul class="space-y-2">
          <li>
            <a href="#" class="text-gray-400 hover:text-yellow-400 transition-colors">
              Privacy Policy
            </a>
          </li>
          <li>
            <a href="#" class="text-gray-400 hover:text-yellow-400 transition-colors">
              Terms of Service
            </a>
          </li>
          <li>
            <a href="#" class="text-gray-400 hover:text-yellow-400 transition-colors">
              Cookie Policy
            </a>
          </li>
        </ul>
      </div>
    </div>
    <div class="border-t border-white/10 mt-12 pt-8 text-center text-gray-400 text-sm">
      <p>© {new Date().getFullYear()} SkipSpot. All rights reserved.</p>
    </div>
  </div>
</footer>
</div>