<script lang="ts">
  import { onMount } from 'svelte';
  import { fade, fly, scale } from 'svelte/transition';
  import { elasticOut, bounceOut } from 'svelte/easing';
  
  // SvelteKit 5 runes
  let count = $state(0);
  let isVisible = $state(false);
  let activeSection = $state('hero');
  let hoveredFeature = $state<string | null>(null);
  
  // Animation timing control
  let animationDelay = $state(0);
  
  // Features list
  const features = [
    { 
      id: 'detect', 
      title: 'Smart Detection', 
      description: 'AI-powered algorithm detects sponsorships with 99% accuracy',
      icon: '🔍',
      color: 'bg-gradient-to-r from-purple-500 to-pink-500'
    },
    { 
      id: 'skip', 
      title: 'Auto Skip', 
      description: 'Automatically jumps past sponsored segments so you never have to manually skip again',
      icon: '⏭️',
      color: 'bg-gradient-to-r from-yellow-400 to-orange-500'
    },
    { 
      id: 'stats', 
      title: 'Time Saved', 
      description: 'Track how much time you\'ve saved by skipping annoying sponsorships',
      icon: '⏱️',
      color: 'bg-gradient-to-r from-green-400 to-emerald-500'
    },
    { 
      id: 'customize', 
      title: 'Customizable', 
      description: 'Choose which types of segments to skip (sponsors, intros, outros, etc.)',
      icon: '⚙️',
      color: 'bg-gradient-to-r from-blue-500 to-cyan-500'
    }
  ];
  
  // Pricing tiers
  const tiers = [
    {
      name: 'Free',
      price: '$0',
      features: [
        'Skip up to 30 sponsorships per month',
        'Basic detection algorithm',
        'Standard support'
      ],
      cta: 'Download Now',
      highlight: false,
      color: 'bg-slate-700'
    },
    {
      name: 'Premium',
      price: '$4.99',
      period: '/month',
      features: [
        'Unlimited sponsorship skipping',
        'Advanced AI detection',
        'Priority support',
        'Early access to new features',
        'No waiting time'
      ],
      cta: 'Upgrade Now',
      highlight: true,
      color: 'bg-gradient-to-r from-violet-600 to-indigo-600'
    }
  ];
  
  // Animated shapes
  const shapes = Array(15).fill(0).map((_, i) => ({
    id: i,
    type: ['circle', 'square', 'triangle', 'hexagon'][Math.floor(Math.random() * 4)],
    size: 20 + Math.random() * 80,
    x: Math.random() * 100,
    y: Math.random() * 100,
    rotation: Math.random() * 360,
    color: [
      'bg-pink-500', 'bg-purple-500', 'bg-yellow-400', 
      'bg-blue-500', 'bg-green-400', 'bg-red-500'
    ][Math.floor(Math.random() * 6)],
    animationDuration: 20 + Math.random() * 40,
    animationDelay: Math.random() * 10
  }));
  
  // Intersection observer for scroll animations
  onMount(() => {
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          const section = entry.target.getAttribute('data-section');
          if (section) {
            activeSection = section;
            isVisible = true;
          }
        }
      });
    }, { threshold: 0.3 });
    
    document.querySelectorAll('section').forEach(section => {
      observer.observe(section);
    });
    
    // Start animation sequence
    isVisible = true;
    
    // Increment counter for time saved stats
    const interval = setInterval(() => {
      count += 1;
    }, 1000);
    
    return () => {
      clearInterval(interval);
      observer.disconnect();
    };
  });
  
  // Handle scroll to section
  function scrollToSection(section: string) {
    const element = document.querySelector(`[data-section="${section}"]`);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
    }
  }
  
  // Handle feature hover
  function setHoveredFeature(id: string | null) {
    hoveredFeature = id;
  }
</script>

<div class="relative overflow-hidden bg-slate-900 text-white">
  <!-- Animated background shapes -->
  <div class="fixed inset-0 overflow-hidden opacity-20 pointer-events-none">
    {#each shapes as shape}
      <div 
        class="{shape.color} absolute rounded-full opacity-70"
        style="
          width: {shape.size}px; 
          height: {shape.size}px; 
          left: {shape.x}%; 
          top: {shape.y}%; 
          transform: rotate({shape.rotation}deg);
          animation: float {shape.animationDuration}s infinite alternate ease-in-out;
          animation-delay: {shape.animationDelay}s;
        "
      ></div>
    {/each}
  </div>
  
  <!-- Navigation -->
  <nav class="fixed top-0 left-0 right-0 z-50 bg-slate-900/80 backdrop-blur-md border-b border-slate-700">
    <div class="container mx-auto px-4 py-4 flex justify-between items-center">
      <div class="flex items-center gap-2">
        <div class="h-10 w-10 rounded-lg bg-gradient-to-br from-purple-600 to-pink-600 flex items-center justify-center">
          <span class="text-xl">⏭️</span>
        </div>
        <span class="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-pink-500">
          SponsorSkip
        </span>
      </div>
      
      <div class="hidden md:flex gap-8">
        {#each ['features', 'pricing', 'testimonials', 'faq'] as section}
          <button 
            class="relative font-medium hover:text-purple-400 transition-colors"
            class:text-purple-400={activeSection === section}
            on:click={() => scrollToSection(section)}
          >
            <span class="capitalize">{section}</span>
            {#if activeSection === section}
              <span 
                class="absolute -bottom-1 left-0 right-0 h-0.5 bg-purple-400 rounded-full"
                in:scale={{ duration: 300, easing: elasticOut }}
              ></span>
            {/if}
          </button>
        {/each}
      </div>
      
      <button 
        class="px-5 py-2 rounded-full bg-gradient-to-r from-purple-600 to-pink-600 font-medium hover:shadow-lg hover:shadow-purple-500/20 transition-all hover:-translate-y-0.5"
      >
        Download
      </button>
    </div>
  </nav>
  
  <!-- Hero Section -->
  <section data-section="hero" class="min-h-screen pt-24 relative">
    <div class="container mx-auto px-4 py-20">
      <div class="max-w-4xl mx-auto text-center">
        {#if isVisible}
          <div in:fly={{ y: 50, duration: 800, delay: 200 }}>
            <h1 class="text-5xl md:text-7xl font-bold mb-6 leading-tight">
              <span class="bg-clip-text text-transparent bg-gradient-to-r from-purple-400 via-pink-500 to-red-500">
                Skip the Sponsors,
              </span>
              <br />
              <span class="bg-clip-text text-transparent bg-gradient-to-r from-yellow-400 via-green-400 to-blue-500">
                Watch What Matters
              </span>
            </h1>
          </div>
          
          <div in:fly={{ y: 50, duration: 800, delay: 400 }}>
            <p class="text-xl md:text-2xl text-slate-300 mb-10 max-w-2xl mx-auto">
              The smartest Chrome extension that automatically detects and skips sponsorship segments in YouTube videos, saving you hours of time.
            </p>
          </div>
          
          <div class="flex flex-col sm:flex-row gap-4 justify-center mb-16" in:fly={{ y: 50, duration: 800, delay: 600 }}>
            <button class="px-8 py-4 rounded-full bg-gradient-to-r from-purple-600 to-pink-600 text-lg font-medium hover:shadow-lg hover:shadow-purple-500/20 transition-all hover:-translate-y-1">
              Download Extension
            </button>
            <button class="px-8 py-4 rounded-full border border-purple-500/30 text-lg font-medium hover:bg-purple-500/10 transition-all">
              How It Works
            </button>
          </div>
          
          <div class="relative" in:fly={{ y: 50, duration: 800, delay: 800 }}>
            <div class="absolute -inset-0.5 bg-gradient-to-r from-purple-600 to-pink-600 rounded-2xl blur opacity-30"></div>
            <div class="relative bg-slate-800 p-1 rounded-2xl">
              <div class="aspect-video rounded-xl overflow-hidden border border-slate-700 shadow-2xl">
                <div class="w-full h-full bg-slate-900 relative flex items-center justify-center">
                  <div class="absolute inset-0 flex items-center justify-center">
                    <div class="w-16 h-16 rounded-full bg-purple-600 flex items-center justify-center">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                    </div>
                  </div>
                  <img 
                    src="https://via.placeholder.com/1280x720" 
                    alt="SponsorSkip in action" 
                    class="w-full h-full object-cover opacity-60"
                  />
                </div>
              </div>
            </div>
          </div>
          
          <div class="mt-16 flex justify-center" in:fly={{ y: 50, duration: 800, delay: 1000 }}>
            <div class="flex gap-4 items-center px-6 py-3 bg-slate-800/50 rounded-full border border-slate-700">
              <div class="flex -space-x-2">
                {#each Array(4) as _, i}
                  <div class="w-8 h-8 rounded-full bg-gradient-to-r from-purple-500 to-pink-500 flex items-center justify-center text-xs font-bold border-2 border-slate-800">
                    {['JD', 'MK', 'TS', 'AR'][i]}
                  </div>
                {/each}
              </div>
              <div class="text-slate-300">
                <span class="font-bold text-white">10,000+</span> users already saving time
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>
    
    <!-- Animated wave divider -->
    <div class="absolute bottom-0 left-0 right-0 h-16 overflow-hidden">
      <svg viewBox="0 0 1200 120" preserveAspectRatio="none" class="absolute bottom-0 w-full h-full">
        <path 
          d="M321.39,56.44c58-10.79,114.16-30.13,172-41.86,82.39-16.72,168.19-17.73,250.45-.39C823.78,31,906.67,72,985.66,92.83c70.05,18.48,146.53,26.09,214.34,3V120H0V0C0,0,0,0,0,0z" 
          class="fill-slate-800"
        ></path>
      </svg>
    </div>
  </section>
  
  <!-- Features Section -->
  <section data-section="features" class="py-20 bg-slate-800 relative">
    <div class="container mx-auto px-4">
      <div class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-pink-500">
          Powerful Features
        </h2>
        <p class="text-xl text-slate-300 max-w-2xl mx-auto">
          Our extension is packed with smart features to enhance your YouTube experience
        </p>
      </div>
      
      <div class="grid md:grid-cols-2 gap-8 max-w-5xl mx-auto">
        {#each features as feature, i}
          <div 
            class="relative group"
            in:fly={{ y: 50, duration: 800, delay: 200 + i * 100 }}
            on:mouseenter={() => setHoveredFeature(feature.id)}
            on:mouseleave={() => setHoveredFeature(null)}
          >
            <div class="absolute -inset-0.5 {feature.color} rounded-2xl blur opacity-30 group-hover:opacity-100 transition duration-300"></div>
            <div class="relative p-6 bg-slate-900 rounded-xl border border-slate-700 h-full transition-all duration-300 group-hover:-translate-y-1">
              <div class="w-14 h-14 rounded-lg {feature.color} flex items-center justify-center mb-4 text-2xl">
                {feature.icon}
              </div>
              <h3 class="text-xl font-bold mb-2">{feature.title}</h3>
              <p class="text-slate-300">{feature.description}</p>
              
              {#if hoveredFeature === feature.id}
                <div 
                  class="absolute -bottom-1 -right-1 w-12 h-12 rounded-br-xl {feature.color}"
                  in:scale={{ duration: 200, easing: elasticOut }}
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 absolute bottom-2 right-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>
                </div>
              {/if}
            </div>
          </div>
        {/each}
      </div>
      
      <div class="mt-20 max-w-4xl mx-auto bg-slate-900 rounded-2xl p-8 border border-slate-700">
        <div class="flex flex-col md:flex-row gap-8 items-center">
          <div class="flex-1">
            <h3 class="text-2xl font-bold mb-4">Time Saved Counter</h3>
            <p class="text-slate-300 mb-6">
              Our users have collectively saved thousands of hours by skipping sponsorships. Join them today!
            </p>
            <div class="flex gap-4">
              <div class="flex-1 bg-slate-800 rounded-lg p-4 border border-slate-700">
                <div class="text-sm text-slate-400 mb-1">Your Time Saved</div>
                <div class="text-2xl font-bold">{Math.floor(count / 60)}m {count % 60}s</div>
              </div>
              <div class="flex-1 bg-slate-800 rounded-lg p-4 border border-slate-700">
                <div class="text-sm text-slate-400 mb-1">Community Saved</div>
                <div class="text-2xl font-bold">127,482 hours</div>
              </div>
            </div>
          </div>
          <div class="w-40 h-40 relative">
            <svg viewBox="0 0 100 100" class="w-full h-full">
              <circle cx="50" cy="50" r="45" fill="none" stroke="#1e293b" stroke-width="8" />
              <circle 
                cx="50" 
                cy="50" 
                r="45" 
                fill="none" 
                stroke="url(#gradient)" 
                stroke-width="8" 
                stroke-dasharray="283" 
                stroke-dashoffset={283 - (283 * Math.min(count / 100, 1))}
                transform="rotate(-90 50 50)"
              />
              <defs>
                <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="0%">
                  <stop offset="0%" stop-color="#8b5cf6" />
                  <stop offset="100%" stop-color="#ec4899" />
                </linearGradient>
              </defs>
            </svg>
            <div class="absolute inset-0 flex items-center justify-center">
              <div class="text-2xl font-bold">{Math.min(Math.floor(count), 100)}%</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
  
  <!-- Pricing Section -->
  <section data-section="pricing" class="py-20 bg-slate-900 relative">
    <div class="container mx-auto px-4">
      <div class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-pink-500">
          Choose Your Plan
        </h2>
        <p class="text-xl text-slate-300 max-w-2xl mx-auto">
          Start for free and upgrade when you need more skips
        </p>
      </div>
      
      <div class="flex flex-col md:flex-row gap-8 max-w-5xl mx-auto">
        {#each tiers as tier, i}
          <div 
            class="flex-1 relative group"
            in:fly={{ y: 50, duration: 800, delay: 200 + i * 200 }}
          >
            {#if tier.highlight}
              <div class="absolute -inset-0.5 bg-gradient-to-r from-purple-600 to-pink-600 rounded-2xl blur opacity-30 group-hover:opacity-100 transition duration-300"></div>
            {/if}
            <div class="relative h-full flex flex-col p-8 bg-slate-800 rounded-xl border border-slate-700 transition-all duration-300 group-hover:-translate-y-1">
              {#if tier.highlight}
                <div class="absolute -top-4 left-0 right-0 flex justify-center">
                  <div class="px-4 py-1 bg-gradient-to-r from-purple-600 to-pink-600 rounded-full text-sm font-medium">
                    Most Popular
                  </div>
                </div>
              {/if}
              
              <h3 class="text-2xl font-bold mb-2">{tier.name}</h3>
              <div class="mb-6">
                <span class="text-4xl font-bold">{tier.price}</span>
                {#if tier.period}
                  <span class="text-slate-400">{tier.period}</span>
                {/if}
              </div>
              
              <ul class="mb-8 flex-1">
                {#each tier.features as feature}
                  <li class="flex items-start gap-2 mb-3">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-400 mt-0.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                    </svg>
                    <span class="text-slate-300">{feature}</span>
                  </li>
                {/each}
              </ul>
              
              <button 
                class={`w-full py-3 rounded-lg font-medium transition-all ${tier.highlight ? 'bg-gradient-to-r from-purple-600 to-pink-600 hover:shadow-lg hover:shadow-purple-500/20' : 'bg-slate-700 hover:bg-slate-600'}`}
              >
                {tier.cta}
              </button>
            </div>
          </div>
        {/each}
      </div>
      
      <div class="mt-16 text-center">
        <p class="text-slate-400 mb-4">All plans include a 14-day money-back guarantee</p>
        <div class="flex flex-wrap justify-center gap-4">
          <div class="flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <span>No credit card required</span>
          </div>
          <div class="flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <span>Cancel anytime</span>
          </div>
          <div class="flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <span>Secure payment</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Animated wave divider -->
    <div class="absolute bottom-0 left-0 right-0 h-16 overflow-hidden">
      <svg viewBox="0 0 1200 120" preserveAspectRatio="none" class="absolute bottom-0 w-full h-full transform rotate-180">
        <path 
          d="M321.39,56.44c58-10.79,114.16-30.13,172-41.86,82.39-16.72,168.19-17.73,250.45-.39C823.78,31,906.67,72,985.66,92.83c70.05,18.48,146.53,26.09,214.34,3V120H0V0C0,0,0,0,0,0z" 
          class="fill-slate-800"
        ></path>
      </svg>
    </div>
  </section>
  
  <!-- Testimonials Section -->
  <section data-section="testimonials" class="py-20 bg-slate-800 relative">
    <div class="container mx-auto px-4">
      <div class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-pink-500">
          What Users Say
        </h2>
        <p class="text-xl text-slate-300 max-w-2xl mx-auto">
          Join thousands of happy users who save time every day
        </p>
      </div>
      
      <div class="grid md:grid-cols-3 gap-8 max-w-6xl mx-auto">
        {#each Array(3) as _, i}
          <div 
            class="bg-slate-900 p-6 rounded-xl border border-slate-700 relative"
            in:fly={{ y: 50, duration: 800, delay: 200 + i * 100 }}
          >
            <div class="absolute -top-5 -left-5 w-10 h-10 rounded-full bg-gradient-to-r from-purple-600 to-pink-600 flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />
              </svg>
            </div>
            
            <div class="flex gap-2 mb-4">
              {#each Array(5) as _}
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-yellow-400" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                </svg>
              {/each}
            </div>
            
            <p class="text-slate-300 mb-6">
              {[
                "This extension is a game changer! I've saved so much time by automatically skipping those long sponsor segments. Worth every penny for the premium version.",
                "I was skeptical at first, but now I can't imagine watching YouTube without SponsorSkip. The AI detection is incredibly accurate.",
                "The time saved counter is addictive to watch. I've saved over 3 hours in just two weeks. The premium tier is definitely worth it!"
              ][i]}
            </p>
            
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-full bg-gradient-to-r from-purple-500 to-pink-500 flex items-center justify-center text-xs font-bold">
                {['JD', 'MK', 'TS'][i]}
              </div>
              <div>
                <div class="font-medium">
                  {['John Doe', 'Mary Kim', 'Tom Smith'][i]}
                </div>
                <div class="text-sm text-slate-400">YouTube Creator</div>
              </div>
            </div>
          </div>
        {/each}
      </div>
    </div>
  </section>
  
  <!-- FAQ Section -->
  <section data-section="faq" class="py-20 bg-slate-900 relative">
    <div class="container mx-auto px-4">
      <div class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-pink-500">
          Frequently Asked Questions
        </h2>
        <p class="text-xl text-slate-300 max-w-2xl mx-auto">
          Everything you need to know about SponsorSkip
        </p>
      </div>
      
      <div class="max-w-3xl mx-auto">
        {#each [
          {
            q: 'How does SponsorSkip detect sponsorships?',
            a: 'Our extension uses a combination of machine learning algorithms and community-contributed data to accurately identify sponsored segments in videos. The AI analyzes audio and visual cues to determine when a sponsorship begins and ends.'
          },
          {
            q: 'Will this work on all YouTube videos?',
            a: 'SponsorSkip works on the vast majority of YouTube videos. Our detection algorithm is constantly improving, and with our community-driven approach, coverage continues to expand every day.'
          },
          {
            q: 'What\'s the difference between free and premium?',
            a: 'The free version allows you to skip up to 30 sponsorships per month. Premium gives you unlimited skips, advanced AI detection for better accuracy, priority support, and early access to new features.'
          },
          {
            q: 'Can I customize what types of segments to skip?',
            a: 'Yes! You can choose to skip sponsors, intros, outros, subscription reminders, and more. Premium users get access to more granular controls.'
          },
          {
            q: 'Is my data being collected?',
            a: 'We only collect anonymous usage statistics to improve our service. We never track your browsing history or personal information.'
          }
        ] as faq, i}
          <div 
            class="mb-6 border border-slate-700 rounded-xl overflow-hidden"
            in:fly={{ y: 20, duration: 400, delay: 100 + i * 100 }}
          >
            <button 
              class="w-full p-5 text-left bg-slate-800 hover:bg-slate-700 transition-colors flex justify-between items-center"
              on:click={() => {
                if (activeSection === `faq-${i}`) {
                  activeSection = 'faq';
                } else {
                  activeSection = `faq-${i}`;
                }
              }}
            >
              <span class="font-medium">{faq.q}</span>
              <svg 
                xmlns="http://www.w3.org/2000/svg" 
                class="h-5 w-5 transition-transform" 
                class:rotate-180={activeSection === `faq-${i}`}
                fill="none" 
                viewBox="0 0 24 24" 
                stroke="currentColor"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            
            {#if activeSection === `faq-${i}`}
              <div 
                class="p-5 bg-slate-900 border-t border-slate-700"
                in:fly={{ y: -10, duration: 200 }}
              >
                <p class="text-slate-300">{faq.a}</p>
              </div>
            {/if}
          </div>
        {/each}
      </div>
      
      <div class="mt-16 text-center">
        <p class="text-slate-300 mb-6">Still have questions?</p>
        <button class="px-6 py-3 rounded-lg bg-gradient-to-r from-purple-600 to-pink-600 font-medium hover:shadow-lg hover:shadow-purple-500/20 transition-all hover:-translate-y-0.5">
          Contact Support
        </button>
      </div>
    </div>
  </section>
  
  <!-- CTA Section -->
  <section class="py-20 bg-gradient-to-r from-purple-900 to-pink-900 relative overflow-hidden">
    <div class="absolute inset-0 opacity-30">
      {#each Array(20) as _, i}
        <div 
          class="absolute rounded-full bg-white"
          style="
            width: {5 + Math.random() * 10}px; 
            height: {5 + Math.random() * 10}px; 
            left: {Math.random() * 100}%; 
            top: {Math.random() * 100}%; 
            animation: float {10 + Math.random() * 20}s infinite alternate ease-in-out;
            animation-delay: {Math.random() * 5}s;
          "
        ></div>
      {/each}
    </div>
    
    <div class="container mx-auto px-4 relative">
      <div class="max-w-4xl mx-auto text-center">
        <h2 class="text-4xl md:text-5xl font-bold mb-6">
          Ready to Skip the Boring Parts?
        </h2>
        <p class="text-xl mb-10 text-purple-100 max-w-2xl mx-auto">
          Join thousands of users who save hours of time every month with SponsorSkip
        </p>
        
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <button class="px-8 py-4 rounded-full bg-white text-purple-900 text-lg font-medium hover:shadow-lg hover:shadow-purple-500/20 transition-all hover:-translate-y-1">
            Download Now
          </button>
          <button class="px-8 py-4 rounded-full border border-white/30 text-lg font-medium hover:bg-white/10 transition-all">
            View Pricing
          </button>
        </div>
      </div>
    </div>
  </section>
  
  <!-- Footer -->
  <footer class="bg-slate-900 border-t border-slate-800 py-12">
    <div class="container mx-auto px-4">
      <div class="grid md:grid-cols-4 gap-8">
        <div>
          <div class="flex items-center gap-2 mb-4">
            <div class="h-10 w-10 rounded-lg bg-gradient-to-br from-purple-600 to-pink-600 flex items-center justify-center">
              <span class="text-xl">⏭️</span>
            </div>
            <span class="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-pink-500">
              SponsorSkip
            </span>
          </div>
          <p class="text-slate-400 mb-4">
            The smartest way to skip sponsorships on YouTube and save hours of your time.
          </p>
          <div class="flex gap-4">
            <a href="#" class="text-slate-400 hover:text-white transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M24 4.557c-.883.392-1.832.656-2.828.775 1.017-.609 1.798-1.574 2.165-2.724-.951.564-2.005.974-3.127 1.195-.897-.957-2.178-1.555-3.594-1.555-3.179 0-5.515 2.966-4.797 6.045-4.091-.205-7.719-2.165-10.148-5.144-1.29 2.213-.669 5.108 1.523 6.574-.806-.026-1.566-.247-2.229-.616-.054 2.281 1.581 4.415 3.949 4.89-.693.188-1.452.232-2.224.084.626 1.956 2.444 3.379 4.6 3.419-2.07 1.623-4.678 2.348-7.29 2.04 2.179 1.397 4.768 2.212 7.548 2.212 9.142 0 14.307-7.721 13.995-14.646.962-.695 1.797-1.562 2.457-2.549z"/>
              </svg>
            </a>
            <a href="#" class="text-slate-400 hover:text-white transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
              </svg>
            </a>
            <a href="#" class="text-slate-400 hover:text-white transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M9 8h-3v4h3v12h5v-12h3.642l.358-4h-4v-1.667c0-.955.192-1.333 1.115-1.333h2.885v-5h-3.808c-3.596 0-5.192 1.583-5.192 4.615v3.385z"/>
              </svg>
            </a>
          </div>
        </div>
        
        <div>
          <h3 class="font-bold mb-4">Product</h3>
          <ul class="space-y-2">
            <li><a href="#" class="text-slate-400 hover:text-white transition-colors">Features</a></li>
            <li><a href="#" class="text-slate-400 hover:text-white transition-colors">Pricing</a></li>
            <li><a href="#" class="text-slate-400 hover:text-white transition-colors">Download</a></li>
            <li><a href="#" class="text-slate-400 hover:text-white transition-colors">Updates</a></li>
          </ul>
        </div>
        
        <div>
          <h3 class="font-bold mb-4">Support</h3>
          <ul class="space-y-2">
            <li><a href="#" class="text-slate-400 hover:text-white transition-colors">Help Center</a></li>
            <li><a href="#" class="text-slate-400 hover:text-white transition-colors">Contact Us</a></li>
            <li><a href="#" class="text-slate-400 hover:text-white transition-colors">FAQ</a></li>
            <li><a href="#" class="text-slate-400 hover:text-white transition-colors">Community</a></li>
          </ul>
        </div>
        
        <div>
          <h3 class="font-bold mb-4">Legal</h3>
          <ul class="space-y-2">
            <li><a href="#" class="text-slate-400 hover:text-white transition-colors">Privacy Policy</a></li>
            <li><a href="#" class="text-slate-400 hover:text-white transition-colors">Terms of Service</a></li>
            <li><a href="#" class="text-slate-400 hover:text-white transition-colors">Cookie Policy</a></li>
            <li><a href="#" class="text-slate-400 hover:text-white transition-colors">GDPR</a></li>
          </ul>
        </div>
      </div>
      
      <div class="mt-12 pt-8 border-t border-slate-800 flex flex-col md:flex-row justify-between items-center">
        <p class="text-slate-400 text-sm mb-4 md:mb-0">
          &copy; {new Date().getFullYear()} SponsorSkip. All rights reserved.
        </p>
        <div class="flex gap-4">
          <a href="#" class="text-sm text-slate-400 hover:text-white transition-colors">Privacy</a>
          <a href="#" class="text-sm text-slate-400 hover:text-white transition-colors">Terms</a>
          <a href="#" class="text-sm text-slate-400 hover:text-white transition-colors">Cookies</a>
        </div>
      </div>
    </div>
  </footer>
</div>

<style>
  @keyframes float {
    0% {
      transform: translateY(0) rotate(0deg);
    }
    50% {
      transform: translateY(-20px) rotate(10deg);
    }
    100% {
      transform: translateY(0) rotate(0deg);
    }
  }
  
  :global(html) {
    scroll-behavior: smooth;
  }
</style>