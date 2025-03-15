

<script>
  import { onMount } from 'svelte';

  // State using Svelte 5 runes
  let isPlaying = $state(false);
  let progress = $state(0);

  // Derived values using runes
  let progressPercent = $derived(`${progress}%`);
  let currentTime = $derived(Math.floor(progress * 0.6));
  let showSponsorOverlay = $derived(progress >= 30 && progress <= 60);

  // Animation states for geometric shapes
  let rotation = $state(0);
  let scale = $state(1);
  let skewX = $state(0);
  let translateX = $state(0);
  let borderRadius = $state('0%');
  
  // Header animation states
  let headerLogoRotation = $state(-90);
  let headerTitleOpacity = $state(0);
  let headerTitleY = $state(-20);
  let navOpacity = $state(0);
  
  // Hero section animation states
  let heroLeftOpacity = $state(0);
  let heroLeftX = $state(-50);
  let heroRightOpacity = $state(0);
  let heroRightScale = $state(0.8);
  
  // Features section animation states
  let featuresSectionOpacity = $state(0);
  let featuresSectionY = $state(30);
  let featuresVisible = $state([false, false, false]);

  // Pricing section animation states
  let pricingOpacity = $state(0);
  let pricingY = $state(30);
  let freeTierOpacity = $state(0);
  let freeTierX = $state(-30);
  let proTierOpacity = $state(0);
  let proTierX = $state(30);

  // CTA section animation states
  let ctaOpacity = $state(0);
  let ctaY = $state(30);

  // FAQ section animation states
  let faqHeaderOpacity = $state(0);
  let faqHeaderY = $state(30);
  let faqItemsOpacity = $state([0, 0, 0, 0]);
  let faqItemsY = $state([20, 20, 20, 20]);

  // Setup animation intervals and initial animations
  $effect(() => {
    /**
     * @type {number | undefined}
     */
    let progressInterval;
    let rotationInterval;
    let scaleInterval;
    let skewInterval;
    let borderRadiusInterval;

    // Initial animations (once)
    headerLogoRotation = 0;
    setTimeout(() => { headerTitleOpacity = 1; headerTitleY = 0; }, 200);
    setTimeout(() => { navOpacity = 1; }, 400);
    
    setTimeout(() => { 
      heroLeftOpacity = 1; 
      heroLeftX = 0; 
    }, 300);
    
    setTimeout(() => { 
      heroRightOpacity = 1; 
      heroRightScale = 1; 
    }, 600);

    // Animate features section on scroll
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          if (entry.target.id === 'features') {
            featuresSectionOpacity = 1;
            featuresSectionY = 0;
            
            // Stagger feature animations
            setTimeout(() => { featuresVisible[0] = true; }, 100);
            setTimeout(() => { featuresVisible[1] = true; }, 200);
            setTimeout(() => { featuresVisible[2] = true; }, 300);
          } else if (entry.target.id === 'pricing') {
            pricingOpacity = 1;
            pricingY = 0;
            
            setTimeout(() => { freeTierOpacity = 1; freeTierX = 0; }, 100);
            setTimeout(() => { proTierOpacity = 1; proTierX = 0; }, 200);
          } else if (entry.target.id === 'faq') {
            faqHeaderOpacity = 1;
            faqHeaderY = 0;
            
            // Stagger FAQ items animations
            faqs.forEach((_, index) => {
              setTimeout(() => { 
                faqItemsOpacity[index] = 1; 
                faqItemsY[index] = 0; 
              }, 100 + (index * 100));
            });
          }
        }
      });
    }, { threshold: 0.2 });
    
    const featuresSection = document.getElementById('features');
    const pricingSection = document.getElementById('pricing');
    const faqSection = document.getElementById('faq');
    
    if (featuresSection) {
      observer.observe(featuresSection);
    }
    
    if (pricingSection) {
      observer.observe(pricingSection);
    }
    
    if (faqSection) {
      observer.observe(faqSection);
    }

    // Observe CTA section
    const ctaObserver = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          ctaOpacity = 1;
          ctaY = 0;
        }
      });
    }, { threshold: 0.2 });
    
    const ctaSections = document.querySelectorAll('section');
    if (ctaSections.length > 3) { // Assuming CTA is the fourth section
      ctaObserver.observe(ctaSections[3]);
    }

    // Video progress animation
    if (isPlaying) {
      progressInterval = setInterval(() => {
        progress += 0.5;
        if (progress >= 100) {
          isPlaying = false;
          progress = 0;
        }
      }, 50);
    }

    // Background shape animations
    rotationInterval = setInterval(() => {
      rotation = (rotation + 1) % 360;
    }, 100);

    scaleInterval = setInterval(() => {
      scale = 1 + 0.2 * Math.sin(Date.now() / 1000);
    }, 50);

    skewInterval = setInterval(() => {
      skewX = 10 * Math.sin(Date.now() / 1500);
      translateX = 100 * Math.sin(Date.now() / 3000);
    }, 50);

    borderRadiusInterval = setInterval(() => {
      const value = 50 * Math.sin(Date.now() / 2000);
      borderRadius = `${Math.abs(value)}%`;
    }, 50);

    return () => {
      clearInterval(progressInterval);
      clearInterval(rotationInterval);
      clearInterval(scaleInterval);
      clearInterval(skewInterval);
      clearInterval(borderRadiusInterval);
      observer.disconnect();
      ctaObserver.disconnect();
    };
  });

  // Features data
  const features = [
    {
      icon: 'fast-forward',
      title: 'Automatic Detection',
      description: 'Our AI-powered algorithm identifies sponsorship segments with high accuracy.'
    },
    {
      icon: 'zap',
      title: 'Lightning Fast',
      description: 'Skip sponsors instantly without any lag or delay in your viewing experience.'
    },
    {
      icon: 'clock',
      title: 'Time Saved',
      description: 'Save hours of your time by automatically skipping through sponsored content.'
    }
  ];

  // Pricing tiers
  const freeTier = {
    name: 'Free',
    description: 'Basic sponsorship skipping',
    price: '$0',
    period: '/month',
    features: [
      'Skip up to 50 sponsorships per month',
      'Basic detection algorithm',
      'Standard support'
    ]
  };

  const proTier = {
    name: 'Pro',
    description: 'Unlimited sponsorship skipping',
    price: '$4.99',
    period: '/month',
    features: [
      'Unlimited sponsorship skipping',
      'Advanced AI detection algorithm',
      'Custom skip settings',
      'Priority support',
      'Early access to new features'
    ]
  };

  // FAQ data
  const faqs = [
    {
      question: 'How does SponsorSkip work?',
      answer: 'SponsorSkip uses machine learning algorithms to detect sponsorship segments in YouTube videos. When a sponsorship is detected, the extension automatically skips to the end of the segment.'
    },
    {
      question: 'Is there a limit to how many videos I can skip with the free version?',
      answer: 'Yes, the free version allows you to skip sponsorships in up to 50 videos per month. For unlimited skipping, you can upgrade to the Pro version.'
    },
    {
      question: 'Can I customize which types of segments to skip?',
      answer: 'Yes, with the Pro version you can customize which types of segments to skip, such as intros, outros, sponsorships, and more.'
    },
    {
      question: 'Does SponsorSkip work on all YouTube videos?',
      answer: 'SponsorSkip works on most YouTube videos, but the accuracy may vary depending on the content. Our algorithm is constantly improving to provide better detection.'
    }
  ];

  // Helper function for icons
  /**
   * @param {string} name
   */
  function getIcon(name) {
    switch (name) {
      case 'fast-forward':
        return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="13 19 22 12 13 5 13 19"></polygon><polygon points="2 19 11 12 2 5 2 19"></polygon></svg>`;
      case 'zap':
        return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon></svg>`;
      case 'clock':
        return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>`;
      case 'play':
        return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>`;
      case 'pause':
        return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="6" y="4" width="4" height="16"></rect><rect x="14" y="4" width="4" height="16"></rect></svg>`;
      case 'check':
        return `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>`;
      default:
        return '';
    }
  }
</script>

<div class="min-h-screen bg-black text-white overflow-hidden">
  <!-- Geometric background elements -->
  <div class="fixed inset-0 z-0 opacity-20">
    <div 
      class="absolute top-20 left-20 w-40 h-40 rounded-full bg-purple-500"
      style="transform: rotate({rotation}deg) scale({scale});"
    ></div>
    <div 
      class="absolute bottom-40 right-20 w-60 h-60 bg-yellow-400"
      style="transform: rotate({rotation}deg); border-radius: {borderRadius};"
    ></div>
    <div 
      class="absolute top-1/2 left-1/3 w-80 h-20 bg-blue-500"
      style="transform: skew({skewX}deg) translateX({translateX}px);"
    ></div>
  </div>
  
  <!-- Header -->
  <header class="relative z-10 flex justify-between items-center p-6 md:p-8">
    <div class="flex items-center gap-2">
      <div 
        class="transition-transform duration-500" 
        style="transform: rotate({headerLogoRotation}deg);"
      >
        <div class="h-8 w-8 text-pink-500">
          {@html getIcon('fast-forward')}
        </div>
      </div>
      <h1 
        class="text-2xl font-bold transition-all duration-500" 
        style="opacity: {headerTitleOpacity}; transform: translateY({headerTitleY}px);"
      >
        SponsorSkip
      </h1>
    </div>
    <nav 
      class="transition-opacity duration-500" 
      style="opacity: {navOpacity};"
    >
      <ul class="flex gap-6 md:gap-10 text-sm md:text-base">
        <li><a href="#features" class="hover:text-pink-400 transition-colors">Features</a></li>
        <li><a href="#pricing" class="hover:text-pink-400 transition-colors">Pricing</a></li>
        <li><a href="#faq" class="hover:text-pink-400 transition-colors">FAQ</a></li>
      </ul>
    </nav>
  </header>
  
  <!-- Hero Section -->
  <section class="relative z-10 mt-10 md:mt-20 px-6 md:px-10 max-w-7xl mx-auto">
    <div class="grid md:grid-cols-2 gap-10 items-center">
      <div 
        class="transition-all duration-700" 
        style="opacity: {heroLeftOpacity}; transform: translateX({heroLeftX}px);"
      >
        <h2 class="text-4xl md:text-6xl font-bold mb-6 leading-tight">
          Skip the <span class="text-pink-500">Sponsors.</span> 
          <br />
          Enjoy the <span class="text-blue-400">Content.</span>
        </h2>
        <p class="text-gray-300 text-lg md:text-xl mb-8">
          Our Chrome extension automatically detects and skips sponsorship segments in YouTube videos, saving you time and enhancing your viewing experience.
        </p>
        <div class="flex flex-col sm:flex-row gap-4">
          <button class="bg-gradient-to-r from-pink-500 to-purple-600 hover:from-pink-600 hover:to-purple-700 text-white font-bold py-3 px-8 rounded-md text-lg">
            Install Free
          </button>
          <button class="border-2 border-pink-500 text-pink-500 hover:bg-pink-500/10 font-bold py-3 px-8 rounded-md text-lg">
            Upgrade to Pro
          </button>
        </div>
      </div>
      
      <!-- Animated Demo -->
      <div 
        class="transition-all duration-700" 
        style="opacity: {heroRightOpacity}; transform: scale({heroRightScale});"
      >
        <div class="relative bg-gray-900 rounded-xl overflow-hidden border-4 border-gray-800">
          <div class="aspect-video relative">
            <!-- Video player mockup -->
            <div class="absolute inset-0 bg-gradient-to-br from-gray-800 to-gray-900 flex items-center justify-center">
              <div class="text-gray-500 text-center">
                <div class="h-16 w-16 mx-auto mb-4 opacity-50">
                  {@html getIcon('play')}
                </div>
                <p>YouTube Video</p>
              </div>
            </div>

            <!-- Sponsorship overlay -->
            {#if showSponsorOverlay}
              <div class="absolute inset-0 bg-gray-800/80 flex items-center justify-center transition-opacity duration-300">
                <div class="bg-pink-500/20 border-2 border-pink-500 rounded-lg p-6 max-w-xs text-center">
                  <div class="h-10 w-10 mx-auto mb-2 text-pink-500">
                    {@html getIcon('fast-forward')}
                  </div>
                  <p class="text-lg font-bold">Sponsorship Detected</p>
                  <p class="text-sm text-gray-300 mt-2">Skipping 30 seconds...</p>
                </div>
              </div>
            {/if}

            <!-- Video controls -->
            <div class="absolute bottom-0 left-0 right-0 p-4 bg-gradient-to-t from-black/80 to-transparent">
              <!-- Progress bar -->
              <div class="h-1 w-full bg-gray-700 rounded-full mb-4 overflow-hidden">
                <div class="h-full bg-red-500" style="width: {progressPercent};"></div>

                <!-- Sponsorship marker -->
                <div class="relative">
                  <div class="absolute top-0 left-[30%] h-3 w-[30%] -mt-2 bg-pink-500/50 rounded-full"></div>
                </div>
              </div>

              <div class="flex justify-between items-center">
                <button 
                  on:click={() => isPlaying = !isPlaying}
                  class="text-white hover:text-pink-400 transition-colors"
                >
                  {#if isPlaying}
                    {@html getIcon('pause')}
                  {:else}
                    {@html getIcon('play')}
                  {/if}
                </button>
                <div class="text-sm text-gray-400">
                  {currentTime}s / 60s
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
  
  <!-- Features Section -->
  <section id="features" class="relative z-10 mt-32 px-6 md:px-10 max-w-7xl mx-auto">
    <div 
      class="text-center mb-16 transition-all duration-700"
      style="opacity: {featuresSectionOpacity}; transform: translateY({featuresSectionY}px);"
    >
      <h2 class="text-3xl md:text-5xl font-bold mb-4">Powerful Features</h2>
      <p class="text-gray-300 text-lg max-w-2xl mx-auto">
        SponsorSkip uses advanced algorithms to detect and skip sponsorship segments automatically.
      </p>
    </div>
    <div class="grid md:grid-cols-3 gap-8">
      {#each features as feature, i}
        <div 
          class="bg-gray-900/50 backdrop-blur-sm border border-gray-800 rounded-xl p-6 hover:border-pink-500/50 transition-colors transition-all duration-500"
          style="opacity: {featuresVisible[i] ? 1 : 0}; transform: translateY({featuresVisible[i] ? 0 : 30}px);"
        >
          <div class="bg-gray-800 rounded-lg p-4 inline-block mb-4">
            <div class="h-10 w-10 {i === 0 ? 'text-pink-500' : i === 1 ? 'text-yellow-400' : 'text-blue-400'}">
              {@html getIcon(feature.icon)}
            </div>
          </div>
          <h3 class="text-xl font-bold mb-2">{feature.title}</h3>
          <p class="text-gray-400">{feature.description}</p>
        </div>
      {/each}
    </div>
  </section>

 <!-- Pricing Section -->
<section id="pricing" class="relative z-10 mt-32 px-6 md:px-10 max-w-7xl mx-auto">
  <div class="text-center mb-16" 
    style:opacity={pricingOpacity} 
    style:transform={`translateY(${pricingY}px)`}>
    <h2 class="text-3xl md:text-5xl font-bold mb-4">Simple Pricing</h2>
    <p class="text-gray-300 text-lg max-w-2xl mx-auto">
      Choose the plan that works best for you.
    </p>
  </div>

  <div class="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto">
    <!-- Free Tier -->
    <div class="bg-gray-900/50 backdrop-blur-sm border border-gray-800 rounded-xl p-8"
      style:opacity={freeTierOpacity} 
      style:transform={`translateX(${freeTierX}px)`}>
      <h3 class="text-2xl font-bold mb-2">{freeTier.name}</h3>
      <p class="text-gray-400 mb-6">{freeTier.description}</p>
      <div class="text-4xl font-bold mb-6">
        {freeTier.price} <span class="text-lg text-gray-400 font-normal">{freeTier.period}</span>
      </div>
      <ul class="space-y-4 mb-8">
        {#each freeTier.features as feature}
          <li class="flex items-start gap-2">
            <div class="mt-1 text-green-400">
              {@html getIcon('check')}
            </div>
            <span>{feature}</span>
          </li>
        {/each}
      </ul>
      <button class="w-full border-2 border-gray-700 hover:border-gray-600 hover:bg-gray-800/50 py-2 px-4 rounded-md">
        Install Free
      </button>
    </div>

    <!-- Pro Tier -->
    <div class="bg-gradient-to-br from-pink-500/20 to-purple-600/20 backdrop-blur-sm border border-pink-500/50 rounded-xl p-8 relative overflow-hidden"
      style:opacity={proTierOpacity} 
      style:transform={`translateX(${proTierX}px)`}>
      <!-- Memphis design element -->
      <div class="absolute -top-10 -right-10 w-40 h-40 bg-yellow-400/20 rounded-full z-0"></div>
      <div class="absolute -bottom-20 -left-20 w-60 h-60 bg-blue-500/10 z-0"></div>
      <div class="relative z-10">
        <div class="flex items-center gap-2 mb-2">
          <h3 class="text-2xl font-bold">{proTier.name}</h3>
          <span class="bg-pink-500 text-xs font-bold px-2 py-1 rounded-full">POPULAR</span>
        </div>
        <p class="text-gray-300 mb-6">{proTier.description}</p>
        <div class="text-4xl font-bold mb-6">
          {proTier.price} <span class="text-lg text-gray-300 font-normal">{proTier.period}</span>
        </div>
        <ul class="space-y-4 mb-8">
          {#each proTier.features as feature}
            <li class="flex items-start gap-2">
              <div class="mt-1 text-pink-400">
                {@html getIcon('check')}
              </div>
              <span>{feature}</span>
            </li>
          {/each}
        </ul>
        <button class="w-full bg-gradient-to-r from-pink-500 to-purple-600 hover:from-pink-600 hover:to-purple-700 py-2 px-4 rounded-md">
          Upgrade to Pro
        </button>
      </div>
    </div>
  </div>
</section>

<!-- CTA Section -->
<section class="relative z-10 mt-32 px-6 md:px-10 py-20">
  <div class="max-w-4xl mx-auto bg-gradient-to-br from-gray-900 to-gray-800 rounded-2xl p-10 border border-gray-700 relative overflow-hidden"
    style:opacity={ctaOpacity} 
    style:transform={`translateY(${ctaY}px)`}>
    <!-- Memphis design elements -->
    <div class="absolute top-0 right-0 w-40 h-40 bg-pink-500/10 rounded-full -translate-y-1/2 translate-x-1/2"></div>
    <div class="absolute bottom-0 left-0 w-60 h-60 bg-blue-500/10 rounded-full translate-y-1/2 -translate-x-1/2"></div>
    <div class="relative z-10 text-center">
      <h2 class="text-3xl md:text-4xl font-bold mb-4">Ready to skip the boring parts?</h2>
      <p class="text-gray-300 text-lg mb-8 max-w-2xl mx-auto">
        Join thousands of users who save hours every month by automatically skipping sponsorship segments.
      </p>
      <div class="flex flex-col sm:flex-row gap-4 justify-center">
        <button class="bg-gradient-to-r from-pink-500 to-purple-600 hover:from-pink-600 hover:to-purple-700 text-white font-bold py-3 px-8 rounded-md text-lg">
          Install Free
        </button>
        <button class="border-2 border-pink-500 text-pink-500 hover:bg-pink-500/10 font-bold py-3 px-8 rounded-md text-lg">
          Upgrade to Pro
        </button>
      </div>
    </div>
  </div>
</section>

<!-- FAQ Section -->
<section id="faq" class="relative z-10 mt-20 px-6 md:px-10 max-w-4xl mx-auto mb-32">
  <div class="text-center mb-16"
    style:opacity={faqHeaderOpacity} 
    style:transform={`translateY(${faqHeaderY}px)`}>
    <h2 class="text-3xl md:text-5xl font-bold mb-4">Frequently Asked Questions</h2>
  </div>
  
  <div class="space-y-6">
    {#each faqs as faq, i}
      <div class="bg-gray-900/50 backdrop-blur-sm border border-gray-800 rounded-xl p-6"
        style:opacity={faqItemsOpacity[i] || 0} 
        style:transform={`translateY(${faqItemsY[i] || 20}px)`}>
        <h3 class="text-xl font-bold mb-2">{faq.question}</h3>
        <p class="text-gray-400">{faq.answer}</p>
      </div>
    {/each}
  </div>
</section>



<!-- Footer -->
<footer class="relative z-10 border-t border-gray-800 mt-20 py-10 px-6 md:px-10">
  <div class="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-4 gap-10">
    <div>
      <div class="flex items-center gap-2 mb-4">
        <div class="h-6 w-6 text-primary">
          {@html getIcon('fast-forward')}
        </div>
        <h3 class="text-xl font-bold">SponsorSkip</h3>
      </div>
      <p class="text-muted-foreground text-sm">
        Skip sponsorships in YouTube videos automatically and save your valuable time.
      </p>
    </div>
    
    <div>
      <h4 class="font-bold mb-4">Product</h4>
      <ul class="space-y-2 text-sm text-muted-foreground">
        <li><a href="#features" class="hover:text-primary transition-colors duration-300">Features</a></li>
        <li><a href="#pricing" class="hover:text-primary transition-colors duration-300">Pricing</a></li>
        <li><a href="#faq" class="hover:text-primary transition-colors duration-300">FAQ</a></li>
      </ul>
    </div>
    
    <div>
      <h4 class="font-bold mb-4">Company</h4>
      <ul class="space-y-2 text-sm text-muted-foreground">
        <li><a href="#" class="hover:text-primary transition-colors duration-300">About</a></li>
        <li><a href="#" class="hover:text-primary transition-colors duration-300">Blog</a></li>
        <li><a href="#" class="hover:text-primary transition-colors duration-300">Careers</a></li>
      </ul>
    </div>
    
    <div>
      <h4 class="font-bold mb-4">Legal</h4>
      <ul class="space-y-2 text-sm text-muted-foreground">
        <li><a href="#" class="hover:text-primary transition-colors duration-300">Privacy Policy</a></li>
        <li><a href="#" class="hover:text-primary transition-colors duration-300">Terms of Service</a></li>
        <li><a href="#" class="hover:text-primary transition-colors duration-300">Cookie Policy</a></li>
      </ul>
    </div>
  </div>
  
  <div class="max-w-7xl mx-auto mt-10 pt-6 border-t border-gray-800 text-center text-sm text-muted-foreground">
    <p>© {new Date().getFullYear()} SponsorSkip. All rights reserved.</p>
  </div>
</footer>

</div>
<!-- 
<style>
  /* Animation classes */
    .animate-in {
      animation: fadeIn 0.7s ease forwards;
    }
    
    @keyframes fadeIn {
      from { opacity: 0; transform: translateY(20px); }
      to { opacity: 1; transform: translateY(0); }
    }
    
    .stagger-delay-1 { animation-delay: 0.1s; }
    .stagger-delay-2 { animation-delay: 0.2s; }
    .stagger-delay-3 { animation-delay: 0.3s; }




  /* Additional styles can be added here if needed */
  :global(html) {
    scroll-behavior: smooth;
  }
  
  /* Animation transitions */
  .transition-all {
    transition-property: all;
    transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  }
  
  .transition-transform {
    transition-property: transform;
    transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  }
  
  .transition-opacity {
    transition-property: opacity;
    transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  }
  
  .transition-colors {
    transition-property: color, background-color, border-color;
    transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  }
  
  .duration-300 {
    transition-duration: 300ms;
  }
  
  .duration-500 {
    transition-duration: 500ms;
  }
  
  .duration-700 {
    transition-duration: 700ms;
  }
</style> -->


<style>
  /* Animation classes */
  .animate-in {
    animation: fadeIn 0.7s ease forwards;
  }
    
  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
  }
    
  .stagger-delay-1 { animation-delay: 0.1s; }
  .stagger-delay-2 { animation-delay: 0.2s; }
  .stagger-delay-3 { animation-delay: 0.3s; }
  
  /* Footer specific styles */
  footer {
    position: relative;
    overflow: hidden;
  }
  
  footer::before {
    content: '';
    position: absolute;
    top: -50px;
    left: -10%;
    width: 120%;
    height: 50px;
    background: radial-gradient(ellipse at center, rgba(var(--primary), 0.15) 0%, rgba(0, 0, 0, 0) 70%);
    pointer-events: none;
  }
  
  footer a {
    display: inline-block;
    position: relative;
  }
  
  footer a::after {
    content: '';
    position: absolute;
    width: 0;
    height: 1px;
    bottom: -2px;
    left: 0;
    background-color: hsl(var(--primary));
    transition: width 0.3s ease;
  }
  
  footer a:hover::after {
    width: 100%;
  }
  
  /* Responsive adjustments */
  @media (max-width: 768px) {
    footer .grid {
      gap: 2rem;
    }
  }
  
  /* Animation transitions */
  .transition-all {
    transition-property: all;
    transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  }
  
  .transition-transform {
    transition-property: transform;
    transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  }
  
  .transition-opacity {
    transition-property: opacity;
    transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  }
  
  .transition-colors {
    transition-property: color, background-color, border-color;
    transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  }
  
  .duration-300 {
    transition-duration: 300ms;
  }
  
  .duration-500 {
    transition-duration: 500ms;
  }
  
  .duration-700 {
    transition-duration: 700ms;
  }
  
  /* Global scroll behavior */
  :global(html) {
    scroll-behavior: smooth;
  }
</style>
