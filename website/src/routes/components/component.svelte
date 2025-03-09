<script>
  import { onMount } from 'svelte';
  import { fade, fly, scale } from 'svelte/transition';
  import { spring } from 'svelte/motion';
  import { quintOut } from 'svelte/easing';

  // Using Svelte 5 runes for state management
  let isMenuOpen = $state(false);
  let scrollY = $state(0);
  
  // FAQ state

  /**@type {number|null} openFaqIndex*/
  let openFaqIndex = $state(null);
  
  // Handle scroll effect
  function handleScroll() {
    scrollY = window.scrollY;
  }
  
  onMount(() => {
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  });
  
  // Features data
  let features = [
    {
      icon: 'fast-forward',
      title: 'Automatic Skip',
      description: 'Our AI automatically detects and skips sponsorship segments in YouTube videos, saving you time and frustration.'
    },
    {
      icon: 'zap',
      title: 'Lightning Fast',
      description: 'SponsorSkip works instantly, with no noticeable delay or buffering when skipping segments.'
    },
    {
      icon: 'shield',
      title: 'Privacy Focused',
      description: "We don't track your browsing history or collect any personal data. Your privacy is our priority."
    },
    {
      icon: 'message-square',
      title: 'Community Driven',
      description: 'Our detection algorithm improves with user feedback, making it more accurate over time.'
    },
    {
      icon: 'settings',
      title: 'Customizable',
      description: 'Premium users can customize skip preferences and choose which types of segments to skip.'
    },
    {
      icon: 'bolt',
      title: 'Works Everywhere',
      description: 'SponsorSkip works on all YouTube platforms including desktop, mobile web, and embedded videos.'
    }
  ];
  
  // FAQ data
  const faqs = [
    {
      question: 'How does SponsorSkip detect sponsorship segments?',
      answer: 'SponsorSkip uses a combination of machine learning algorithms and community-contributed data to identify sponsorship segments in videos. Our AI analyzes audio and visual cues that typically indicate sponsored content.'
    },
    {
      question: 'Will SponsorSkip work on all YouTube videos?',
      answer: 'SponsorSkip works on most YouTube videos, but its effectiveness may vary depending on the video content. Our detection algorithm is constantly improving to cover more videos and different types of sponsored content.'
    },
    {
      question: 'Is there a limit to how many videos I can use SponsorSkip on?',
      answer: 'Free users can skip sponsorships on up to 10 videos per day. Premium users enjoy unlimited skips across all their devices.'
    },
    {
      question: 'Does SponsorSkip collect my personal data?',
      answer: 'No, SponsorSkip does not collect any personal data or browsing history. We only collect anonymous usage statistics to improve our service.'
    },
    {
      question: 'Can I customize what types of segments get skipped?',
      answer: 'Yes, Premium users can customize which types of segments to skip, such as sponsorships, intros, outros, and more. Free users have access to basic sponsorship skipping only.'
    }
  ];
  
  // Toggle FAQ item
  /**
	 * @param {number } index
	 */
  function toggleFaq(index) {
    openFaqIndex = openFaqIndex === index ? null : index;
  }
  
  // Current year for footer
  let currentYear = $derived(new Date().getFullYear());


let gridCells = Array(100).fill(0).map(() => ({
    shouldAnimate: Math.random() > 0.8,
    duration: Math.random() * 5 + 2,
    delay: Math.random() * 2,
    isHighlighted: Math.random() > 0.92
  }));
  
  // Animation properties for floating blobs
  const blob1Position = spring({ x: 0, y: 0 }, {
    stiffness: 0.05,
    damping: 0.3
  });
  
  const blob2Position = spring({ x: 0, y: 0 }, {
    stiffness: 0.05,
    damping: 0.3
  });
  
  // Animate the blobs
  const animateBlobs = () => {
    const loop = () => {
      blob1Position.set({ x: Math.sin(Date.now() / 4000) * 50, y: Math.cos(Date.now() / 5000) * 30 });
      blob2Position.set({ x: Math.cos(Date.now() / 4500) * -50, y: Math.sin(Date.now() / 4000) * -30 });
      requestAnimationFrame(loop);
    };
    
    requestAnimationFrame(loop);
  };



  $effect(()=>{
 if (typeof window !== 'undefined') {
      animateBlobs();
    }
  })




</script>

<div class="min-h-screen bg-black text-white overflow-hidden">
  <!-- Header -->
  <header class="fixed top-0 left-0 right-0 z-50 px-4 py-4 md:px-8 md:py-6">
    <div class="container mx-auto flex justify-between items-center">
      <div class="flex items-center gap-2">
        <div class="animate-spin-slow">
          <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-purple-500">
            <polygon points="13 19 22 12 13 5 13 19"></polygon>
            <polygon points="2 19 11 12 2 5 2 19"></polygon>
          </svg>
        </div>
        <span class="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-purple-500 to-pink-500">
          SponsorSkip
        </span>
      </div>
      
      <!-- Desktop Navigation -->
      <nav class="hidden md:flex items-center gap-8">
        <a href="#features" class="hover:text-purple-400 transition-colors">
          Features
        </a>
        <a href="#pricing" class="hover:text-purple-400 transition-colors">
          Pricing
        </a>
        <a href="#faq" class="hover:text-purple-400 transition-colors">
          FAQ
        </a>
        <button class="bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white rounded-full px-6 py-2">
          Get Extension
        </button>
      </nav>
      
      <!-- Mobile Menu Button -->
      <button 
        class="md:hidden text-white"
        onclick={() => isMenuOpen = !isMenuOpen}
      >
        {#if isMenuOpen}
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="3" y1="12" x2="21" y2="12"></line>
            <line x1="3" y1="6" x2="21" y2="6"></line>
            <line x1="3" y1="18" x2="21" y2="18"></line>
          </svg>
        {/if}
      </button>
    </div>
    
    <!-- Mobile Navigation -->
    {#if isMenuOpen}
      <div 
        transition:fly={{ y: -20, duration: 300 }}
        class="absolute top-16 left-0 right-0 bg-black border-t border-gray-800 py-4 md:hidden"
      >
        <div class="container mx-auto flex flex-col gap-4 px-4">
          <a 
            href="#features" 
            class="py-2 hover:text-purple-400 transition-colors"
            onclick={() => isMenuOpen = false}
          >
            Features
          </a>
          <a 
            href="#pricing" 
            class="py-2 hover:text-purple-400 transition-colors"
            onclick={() => isMenuOpen = false}
          >
            Pricing
          </a>
          <a 
            href="#faq" 
            class="py-2 hover:text-purple-400 transition-colors"
            onclick={() => isMenuOpen = false}
          >
            FAQ
          </a>
          <button class="bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white rounded-full py-2">
            Get Extension
          </button>
        </div>
      </div>
    {/if}
  </header>

  <!-- Hero Section -->
  <section class="relative pt-32 pb-20 md:pt-40 md:pb-32 overflow-hidden">
    <!-- Animated Background Elements -->
    <div class="absolute inset-0 -z-10 overflow-hidden">
      <div 
        class="absolute top-20 left-10 w-64 h-64 rounded-full bg-purple-700/20 blur-3xl"
        style="animation: float1 8s infinite ease-in-out;"
      ></div>
      <div 
        class="absolute bottom-20 right-10 w-80 h-80 rounded-full bg-pink-700/20 blur-3xl"
        style="animation: float2 10s infinite ease-in-out;"
      ></div>
      <div 
        class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 rounded-full bg-blue-700/20 blur-3xl"
        style="animation: pulse 12s infinite ease-in-out;"
      ></div>
    </div>

    <div class="container mx-auto px-4 md:px-8">
      <div class="max-w-4xl mx-auto text-center">
        <div in:fly={{ y: 20, duration: 600 }}>
          <h1 class="text-4xl md:text-6xl lg:text-7xl font-bold mb-6 leading-tight">
            <span class="bg-clip-text text-transparent bg-gradient-to-r from-purple-400 via-pink-500 to-red-500">
              Skip the Sponsors,
            </span>
            <br />
            <span class="text-white">
              Enjoy the Content
            </span>
          </h1>
          <p class="text-xl md:text-2xl text-gray-300 mb-8 max-w-2xl mx-auto">
            The smartest Chrome extension that automatically detects and skips sponsorship segments in YouTube videos.
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <button class="bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white text-lg rounded-full px-8 py-6">
              Get SponsorSkip Free
            </button>
            <button class="border-purple-500 text-purple-400 hover:bg-purple-950/30 text-lg rounded-full px-8 py-6">
              See How It Works
            </button>
          </div>
        </div>

        <!-- Browser Mockup -->
        <div
          in:fly={{ y: 40, duration: 800, delay: 300 }}
          class="mt-16 relative"
        >
          <div class="relative mx-auto w-full max-w-3xl rounded-xl shadow-2xl overflow-hidden border border-gray-800">
            <div class="h-8 bg-gray-900 flex items-center px-4">
              <div class="flex space-x-2">
                <div class="w-3 h-3 rounded-full bg-red-500"></div>
                <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
                <div class="w-3 h-3 rounded-full bg-green-500"></div>
              </div>
            </div>
            <div class="bg-gray-950 aspect-video relative">
              <div class="absolute inset-0 flex items-center justify-center">
                <div class="w-full max-w-xl bg-gray-900 rounded-lg overflow-hidden">
                  <div class="h-12 bg-gray-800 flex items-center px-4">
                    <div class="w-3/4 h-4 bg-gray-700 rounded"></div>
                  </div>
                  <div class="relative aspect-video bg-gray-800">
                    <div class="absolute inset-0 flex items-center justify-center">
                      <div class="text-center">
                        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-purple-500 mx-auto mb-4">
                          <polygon points="13 19 22 12 13 5 13 19"></polygon>
                          <polygon points="2 19 11 12 2 5 2 19"></polygon>
                        </svg>
                        <div class="h-2 bg-gradient-to-r from-purple-600 to-pink-600 rounded-full max-w-xs mx-auto progress-bar"></div>
                        <p class="text-white mt-4 font-medium">Skipping Sponsorship...</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Animated Geometric Shapes -->
          <div 
            class="absolute -top-8 -left-8 w-16 h-16 bg-purple-500 rounded-lg"
            style="animation: rotate-scale 8s infinite ease-in-out;"
          ></div>
          <div 
            class="absolute -bottom-8 -right-8 w-20 h-20 bg-pink-500 rounded-full"
            style="animation: rotate-scale-reverse 10s infinite ease-in-out;"
          ></div>
          <div 
            class="absolute top-1/2 -right-12 w-24 h-6 bg-blue-500"
            style="animation: skew-x 6s infinite ease-in-out;"
          ></div>
          <div 
            class="absolute bottom-1/4 -left-12 w-6 h-24 bg-yellow-500"
            style="animation: skew-y 7s infinite ease-in-out;"
          ></div>
        </div>
      </div>
    </div>
  </section>

  <!-- Stats Section -->
  <section class="py-16 bg-gradient-to-b from-black to-purple-950/20">
    <div class="container mx-auto px-4 md:px-8">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        {#each [
          { icon: 'clock', value: '2.5M+', label: 'Hours Saved', color: 'purple' },
          { icon: 'zap', value: '10M+', label: 'Sponsorships Skipped', color: 'pink' },
          { icon: 'crown', value: '500K+', label: 'Happy Users', color: 'blue' }
        ] as stat, i}
          <div 
            in:fly={{ y: 20, duration: 500, delay: i * 200 }}
            class="bg-gray-900/50 backdrop-blur-lg rounded-xl p-8 text-center border border-{stat.color}-500/20"
          >
            <div class="float-animation mx-auto mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-{stat.color}-400 mx-auto">
                {#if stat.icon === 'clock'}
                  <circle cx="12" cy="12" r="10"></circle>
                  <polyline points="12 6 12 12 16 14"></polyline>
                {:else if stat.icon === 'zap'}
                  <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>
                {:else if stat.icon === 'crown'}
                  <path d="M2 4l3 12h14l3-12-6 7-4-7-4 7-6-7zm3 16h14"></path>
                {/if}
              </svg>
            </div>
            <h3 class="text-4xl font-bold text-white mb-2">{stat.value}</h3>
            <p class="text-gray-400">{stat.label}</p>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- Features Section -->
<section id="features" class="py-20 md:py-32 relative overflow-hidden">
  <!-- Animated Background Elements -->
  <div class="absolute inset-0 -z-10">
    <div class="absolute top-0 left-0 w-full h-full opacity-10">
      <div class="grid grid-cols-10 grid-rows-10 h-full w-full">
        {#each gridCells as cell, i}
          <div 
            class="border-[0.5px] border-purple-500/10"
            style="animation: {cell.shouldAnimate ? `grid-cell-pulse ${cell.duration}s infinite ease-in-out ${cell.delay}s` : ''};
                   background-color: {cell.isHighlighted ? 'rgba(168, 85, 247, 0.2)' : 'transparent'};"
          ></div>
        {/each}
      </div>
    </div>
  </div>

  <div class="container mx-auto px-4 md:px-8">
    <div 
      in:fly={{ y: 20, duration: 500 }}
      class="text-center mb-16"
    >
      <h2 class="text-3xl md:text-5xl font-bold mb-6">
        <span class="bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-pink-500">
          Powerful Features
        </span>
      </h2>
      <p class="text-xl text-gray-300 max-w-2xl mx-auto">
        SponsorSkip comes packed with intelligent features to enhance your YouTube experience.
      </p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
      {#each features as feature, index}
        <div
          in:fly={{ y: 20, duration: 500, delay: index * 100 }}
          class="bg-gray-900/50 backdrop-blur-lg rounded-xl p-8 border border-gray-800 hover:border-purple-500/50 transition-colors group"
        >
          <div 
            class="w-14 h-14 rounded-lg bg-purple-900/50 flex items-center justify-center mb-6 group-hover:bg-purple-800/50 transition-colors feature-icon"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-purple-400">
              {#if feature.icon === 'fast-forward'}
                <polygon points="13 19 22 12 13 5 13 19"></polygon>
                <polygon points="2 19 11 12 2 5 2 19"></polygon>
              {:else if feature.icon === 'zap'}
                <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>
              {:else if feature.icon === 'shield'}
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
              {:else if feature.icon === 'message-square'}
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
              {:else if feature.icon === 'settings'}
                <circle cx="12" cy="12" r="3"></circle>
                <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
              {:else if feature.icon === 'bolt'}
                <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"></path>
              {/if}
            </svg>
          </div>
          <h3 class="text-xl font-bold mb-3 text-white">{feature.title}</h3>
          <p class="text-gray-400">{feature.description}</p>
        </div>
      {/each}
    </div>
  </div>
</section>

<!-- Pricing Section -->
<section id="pricing" class="py-20 md:py-32 bg-gradient-to-b from-purple-950/20 to-black relative overflow-hidden">
  <!-- Animated Geometric Shapes -->
  <div 
    class="absolute top-20 left-10 w-40 h-40 rounded-full bg-purple-700/10 blur-3xl"
    style="transform: translate({$blob1Position.x}px, {$blob1Position.y}px);"
  ></div>
  <div 
    class="absolute bottom-20 right-10 w-60 h-60 rounded-full bg-pink-700/10 blur-3xl"
    style="transform: translate({$blob2Position.x}px, {$blob2Position.y}px);"
  ></div>

  <div class="container mx-auto px-4 md:px-8">
    <div 
      in:fly={{ y: 20, duration: 500 }}
      class="text-center mb-16"
    >
      <h2 class="text-3xl md:text-5xl font-bold mb-6">
        <span class="bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-pink-500">
          Choose Your Plan
        </span>
      </h2>
      <p class="text-xl text-gray-300 max-w-2xl mx-auto">
        Upgrade to Premium for unlimited skips and exclusive features.
      </p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-4xl mx-auto">
      <!-- Free Plan -->
      <div
        in:fly={{ x: -20, duration: 500 }}
        class="bg-gray-900/50 backdrop-blur-lg rounded-xl p-8 border border-gray-800"
      >
        <div class="text-center mb-6">
          <h3 class="text-xl font-bold text-gray-300 mb-2">Free</h3>
          <div class="flex items-center justify-center">
            <span class="text-4xl font-bold text-white">$0</span>
            <span class="text-gray-400 ml-2">/month</span>
          </div>
        </div>
        <ul class="space-y-4 mb-8">
          <li class="flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-purple-500 mr-2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
            <span class="text-gray-300">Up to 10 skips per day</span>
          </li>
          <li class="flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-purple-500 mr-2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
            <span class="text-gray-300">Basic detection algorithm</span>
          </li>
          <li class="flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-purple-500 mr-2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
            <span class="text-gray-300">Standard support</span>
          </li>
          <li class="flex items-center opacity-50">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-gray-500 mr-2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
            <span class="text-gray-500">No customization options</span>
          </li>
          <li class="flex items-center opacity-50">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-gray-500 mr-2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
            <span class="text-gray-500">No advanced features</span>
          </li>
        </ul>
        <button class="w-full bg-gray-800 hover:bg-gray-700 text-white rounded-full py-6">
          Get Started
        </button>
      </div>

      <!-- Premium Plan -->
      <div
        in:fly={{ x: 20, duration: 500, delay: 200 }}
        class="bg-gradient-to-br from-purple-900/50 to-pink-900/50 backdrop-blur-lg rounded-xl p-8 border border-purple-500/30 relative overflow-hidden"
      >
        <!-- Animated Background -->
        <div 
          class="absolute inset-0 -z-10 opacity-20 radial-pulse"
          style="background-image: radial-gradient(circle at center, rgba(168, 85, 247, 0.4) 0%, transparent 70%);
                 background-size: 100% 100%;"
        ></div>

        <div class="absolute -top-4 -right-4">
          <div class="bg-gradient-to-r from-purple-600 to-pink-600 text-white text-xs font-bold px-4 py-1 rounded-full transform rotate-12">
            POPULAR
          </div>
        </div>

        <div class="text-center mb-6">
          <h3 class="text-xl font-bold text-white mb-2">Premium</h3>
          <div class="flex items-center justify-center">
            <span class="text-4xl font-bold text-white">$4.99</span>
            <span class="text-gray-300 ml-2">/month</span>
          </div>
        </div>
        <ul class="space-y-4 mb-8">
          <li class="flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-pink-500 mr-2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
            <span class="text-white">Unlimited skips</span>
          </li>
          <li class="flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-pink-500 mr-2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
            <span class="text-white">Advanced AI detection</span>
          </li>
          <li class="flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-pink-500 mr-2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
            <span class="text-white">Priority support</span>
          </li>
          <li class="flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-pink-500 mr-2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
            <span class="text-white">Custom skip preferences</span>
          </li>
          <li class="flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-pink-500 mr-2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
            <span class="text-white">Early access to new features</span>
          </li>
        </ul>
        <div class="scale-on-hover">
          <button class="w-full bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white rounded-full py-6">
            Upgrade Now
          </button>
        </div>
      </div>
    </div>
  </div>
</section>


  <!-- FAQ Section -->
  <section id="faq" class="py-20 md:py-32">
    <div class="container mx-auto px-4 md:px-8">
      <div 
        in:fly={{ y: 20, duration: 500 }}
        class="text-center mb-16"
      >
        <h2 class="text-3xl md:text-5xl font-bold mb-6">
          <span class="bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-pink-500">
            Frequently Asked Questions
          </span>
        </h2>
        <p class="text-xl text-gray-300 max-w-2xl mx-auto">
          Got questions? We've got answers.
        </p>
      </div>

      <div class="max-w-3xl mx-auto space-y-6">
        {#each faqs as faq, index}
          <div 
            in:fly={{ y: 20, duration: 500, delay: index * 100 }}
            class="border border-gray-800 rounded-lg overflow-hidden"
          >
            <button
              onclick={() => toggleFaq(index)}
              class="flex justify-between items-center w-full p-6 text-left bg-gray-900/50 hover:bg-gray-900/80 transition-colors"
            >
              <span class="text-lg font-medium text-white">{faq.question}</span>
              <div class="transform transition-transform duration-300" class:rotate-180={openFaqIndex === index}>
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-purple-500 transform rotate-90">
                  <polyline points="9 18 15 12 9 6"></polyline>
                </svg>
              </div>
            </button>
            {#if openFaqIndex === index}
              <div transition:fly={{ y: -20, duration: 300 }} class="p-6 bg-gray-900/30 text-gray-300">
                {faq.answer}
              </div>
            {/if}
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- CTA Section -->
  <section class="py-20 md:py-32 bg-gradient-to-b from-black to-purple-950/20 relative overflow-hidden">
    <!-- Animated Geometric Shapes -->
    <div class="absolute inset-0 -z-10">
      <div 
        class="absolute top-0 right-0 w-64 h-64 bg-purple-700/10 blur-3xl"
        style="animation: float-alt 8s infinite ease-in-out;"
      ></div>
      <div 
        class="absolute bottom-0 left-0 w-80 h-80 bg-pink-700/10 blur-3xl"
        style="animation: float-alt-reverse 10s infinite ease-in-out;"
      ></div>
    </div>

    <div class="container mx-auto px-4 md:px-8">
      <div class="max-w-4xl mx-auto bg-gradient-to-br from-purple-900/30 to-pink-900/30 backdrop-blur-lg rounded-2xl p-8 md:p-12 border border-purple-500/20">
        <div 
          in:fly={{ y: 20, duration: 500 }}
          class="text-center"
        >
          <h2 class="text-3xl md:text-4xl font-bold mb-6 text-white">
            Ready to Skip the Boring Parts?
          </h2>
          <p class="text-xl text-gray-300 mb-8 max-w-2xl mx-auto">
            Join thousands of users who save time every day with SponsorSkip. Get started for free or upgrade to Premium for unlimited skips.
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <div class="scale-on-hover">
              <button class="bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white text-lg rounded-full px-8 py-6 w-full sm:w-auto">
                Upgrade to Premium
              </button>
            </div>
            <button class="border border-purple-500 text-purple-400 hover:bg-purple-950/30 text-lg rounded-full px-8 py-6">
              Install Free Version
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Footer -->
  <footer class="py-12 border-t border-gray-800">
    <div class="container mx-auto px-4 md:px-8">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
        <div class="md:col-span-1">
          <div class="flex items-center gap-2 mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-purple-500">
              <polygon points="13 19 22 12 13 5 13 19"></polygon>
              <polygon points="2 19 11 12 2 5 2 19"></polygon>
            </svg>
            <span class="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-purple-500 to-pink-500">
              SponsorSkip
            </span>
          </div>
          <p class="text-gray-400 mb-4">
            Skip the sponsors, enjoy the content. The smartest way to enhance your YouTube experience.
          </p>
          <div class="flex space-x-4">
            <!-- svelte-ignore a11y_consider_explicit_label -->
            <a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">
              <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path fill-rule="evenodd" d="M22 12c0-5.523-4.477-10-10-10S2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.988C18.343 21.128 22 16.991 22 12z" clip-rule="evenodd" />
              </svg>
            </a>
            <!-- svelte-ignore a11y_consider_explicit_label -->
            <a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">
              <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path d="M8.29 20.251c7.547 0 11.675-6.253 11.675-11.675 0-.178 0-.355-.012-.53A8.348 8.348 0 0022 5.92a8.19 8.19 0 01-2.357.646 4.118 4.118 0 001.804-2.27 8.224 8.224 0 01-2.605.996 4.107 4.107 0 00-6.993 3.743 11.65 11.65 0 01-8.457-4.287 4.106 4.106 0 001.27 5.477A4.072 4.072 0 012.8 9.713v.052a4.105 4.105 0 003.292 4.022 4.095 4.095 0 01-1.853.07 4.108 4.108 0 003.834 2.85A8.233 8.233 0 012 18.407a11.616 11.616 0 006.29 1.84" />
              </svg>
            </a>
            <!-- svelte-ignore a11y_consider_explicit_label -->
            <a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">
              <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd" />
              </svg>
            </a>
          </div>
        </div>
        <div>
          <h3 class="text-white font-medium mb-4">Product</h3>
          <ul class="space-y-2">
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Features</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Pricing</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">FAQ</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Changelog</a></li>
          </ul>
        </div>
        <div>
          <h3 class="text-white font-medium mb-4">Resources</h3>
          <ul class="space-y-2">
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Documentation</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Support</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Privacy Policy</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Terms of Service</a></li>
          </ul>
        </div>
        <div>
          <h3 class="text-white font-medium mb-4">Company</h3>
          <ul class="space-y-2">
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">About</a></li>
            <!-- svelte-ignore a11y_invalid_attribute -->
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Blog</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Careers</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Contact</a></li>
          </ul>
        </div>
      </div>
      <div class="mt-12 pt-8 border-t border-gray-800 text-center">
        <p class="text-gray-400">
          &copy; {currentYear} SponsorSkip. All rights reserved.
        </p>
      </div>
    </div>
  </footer>
</div>

<style>
  /* Animations */
  @keyframes float1 {
    0%, 100% { transform: translate(0, 0); }
    50% { transform: translate(50px, 30px); }
  }
  
  @keyframes float2 {
    0%, 100% { transform: translate(0, 0); }
    50% { transform: translate(-50px, -30px); }
  }
  
  @keyframes float-alt {
    0%, 100% { transform: translate(0, 0); }
    50% { transform: translate(-30px, 30px); }
  }
  
  @keyframes float-alt-reverse {
    0%, 100% { transform: translate(0, 0); }
    50% { transform: translate(30px, -30px); }
  }
  
  @keyframes pulse {
    0%, 100% { transform: translate(-50%, -50%) scale(1); }
    50% { transform: translate(-50%, -50%) scale(1.2); }
  }
  
  @keyframes rotate-scale {
    0%, 100% { transform: rotate(0deg) scale(1); }
    50% { transform: rotate(180deg) scale(1.2); }
  }
  
  @keyframes rotate-scale-reverse {
    0%, 100% { transform: rotate(0deg) scale(1); }
    50% { transform: rotate(-180deg) scale(1.3); }
  }
  
  @keyframes skew-x {
    0%, 100% { transform: translateX(0) skewX(0deg); }
    50% { transform: translateX(-20px) skewX(10deg); }
  }
  
  @keyframes skew-y {
    0%, 100% { transform: translateY(0) skewY(0deg); }
    50% { transform: translateY(20px) skewY(10deg); }
  }
  
  @keyframes grid-cell-pulse {
    0%, 100% { opacity: 0.1; background-color: transparent; }
    50% { opacity: 0.5; background-color: rgba(168, 85, 247, 0.2); }
  }
  
  .animate-spin-slow {
    animation: spin 2s linear infinite;
  }
  
  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }
  
  .progress-bar {
    animation: progress 2s linear infinite;
  }
  
  @keyframes progress {
    0% { width: 0%; }
    100% { width: 100%; }
  }
  
  .radial-pulse {
    animation: radial-pulse 10s linear infinite alternate;
  }
  
  @keyframes radial-pulse {
    0% { background-position: 0% 0%; }
    100% { background-position: 100% 100%; }
  }
  
  .float-animation {
    animation: float 2s ease-in-out infinite;
  }
  
  @keyframes float {
    0%, 100% { transform: translateY(0); }
    50% { transform: translateY(-10px); }
  }
  
  .feature-icon:hover {
    transform: rotate(5deg) scale(1.1);
    transition: transform 0.3s ease;
  }
  
  .scale-on-hover:hover {
    transform: scale(1.05);
    transition: transform 0.3s ease;
  }


  @keyframes grid-cell-pulse {
    0%, 100% { opacity: 0.1; background-color: transparent; }
    50% { opacity: 0.5; background-color: rgba(168, 85, 247, 0.2); }
  }
  
  .radial-pulse {
    animation: radial-pulse 10s linear infinite alternate;
  }
  
  @keyframes radial-pulse {
    0% { background-position: 0% 0%; }
    100% { background-position: 100% 100%; }
  }
  
  .feature-icon:hover {
    transform: rotate(5deg) scale(1.1);
    transition: transform 0.3s ease;
  }
  
  .scale-on-hover:hover {
    transform: scale(1.05);
    transition: transform 0.3s ease;
  }










</style>