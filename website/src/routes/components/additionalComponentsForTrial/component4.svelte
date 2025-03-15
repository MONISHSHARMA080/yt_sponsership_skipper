<script>
  import { onMount, onDestroy } from 'svelte';
  import { spring } from 'svelte/motion';
  import { ChevronDown, FastForward, Shield, Zap, Clock, Award, ArrowRight } from 'lucide-svelte';
//   import { fade, fly } from 'svelte/transition';
  import { fade, fly, scale } from 'svelte/transition';

  
  let scrollY = 0;
  
  // For animations
  const pulseScale = spring({ x: 0, y: 0, scale: 1 }, {
    stiffness: 0.1,
    damping: 0.4
  });
  
  const diamondRotate = spring({ rotate: 0, scale: 1 }, {
    stiffness: 0.05,
    damping: 0.3
  });
  
  const hexagonRotate = spring({ rotate: 0, y: 0 }, {
    stiffness: 0.06,
    damping: 0.3
  });
  
  const blobRadius = spring({ 
    borderRadius: "30% 70% 70% 30% / 30% 30% 70% 70%", 
    rotate: 0 
  }, {
    stiffness: 0.07,
    damping: 0.3
  });
  
  const progressBar = spring(30, {
    stiffness: 0.02,
    damping: 0.2
  });
  
  // Animate the scrolling arrow
  const scrollArrow = spring(0, {
    stiffness: 0.2,
    damping: 0.4
  });
  
  // Flag for button animation visibility
  let showSkipButton = false;
  
  // Update animations in intervals
  onMount(() => {
    const handleScroll = () => {
      scrollY = window.scrollY;
    };
    
    window.addEventListener('scroll', handleScroll);
    
    const pulseInterval = setInterval(() => {
      pulseScale.set({ x: 0, y: 0, scale: 1 });
      setTimeout(() => {
        pulseScale.set({ x: 10, y: 15, scale: 1.05 });
        setTimeout(() => {
          pulseScale.set({ x: 0, y: 0, scale: 1 });
        }, 4000);
      }, 4000);
    }, 8000);
    
    const diamondInterval = setInterval(() => {
      diamondRotate.set({ rotate: 0, scale: 1 });
      setTimeout(() => {
        diamondRotate.set({ rotate: 10, scale: 1.1 });
        setTimeout(() => {
          diamondRotate.set({ rotate: 0, scale: 1 });
        }, 6000);
      }, 6000);
    }, 12000);
    
    const hexagonInterval = setInterval(() => {
      hexagonRotate.set({ rotate: 0, y: 0 });
      setTimeout(() => {
        hexagonRotate.set({ rotate: -15, y: -20 });
        setTimeout(() => {
          hexagonRotate.set({ rotate: 0, y: 0 });
        }, 5000);
      }, 5000);
    }, 10000);
    
    const blobInterval = setInterval(() => {
      blobRadius.set({ 
        borderRadius: "30% 70% 70% 30% / 30% 30% 70% 70%", 
        rotate: 0 
      });
      setTimeout(() => {
        blobRadius.set({ 
          borderRadius: "50% 50% 50% 50% / 50% 50% 50% 50%", 
          rotate: 10 
        });
        setTimeout(() => {
          blobRadius.set({ 
            borderRadius: "30% 70% 70% 30% / 30% 30% 70% 70%", 
            rotate: 0 
          });
        }, 4000);
      }, 4000);
    }, 8000);
    
    const progressInterval = setInterval(() => {
      progressBar.set(30);
      setTimeout(() => {
        progressBar.set(70);
      }, 2500);
    }, 5000);
    
    const scrollArrowInterval = setInterval(() => {
      scrollArrow.set(0);
      setTimeout(() => {
        scrollArrow.set(10);
        setTimeout(() => {
          scrollArrow.set(0);
        }, 750);
      }, 750);
    }, 1500);
    

let  faqItems = [
    {
      question: "How does SponsorSkip detect sponsored segments?",
      answer: "SponsorSkip uses a combination of machine learning algorithms and community-contributed data to identify sponsored segments in videos with high accuracy. Our system analyzes audio, visual cues, and patterns to determine when a sponsorship begins and ends.",
    },
    {
      question: "Will SponsorSkip work on all YouTube videos?",
      answer: "SponsorSkip works on most YouTube videos, but detection accuracy may vary. Our system continuously improves as more users contribute data. Premium users get access to our most advanced detection algorithms for better accuracy.",
    },
    {
      question: "Can I customize which types of segments to skip?",
      answer: "Yes! Premium users can customize which types of segments to skip, including sponsorships, intros, outros, subscription reminders, and more. Free users have basic customization options.",
    },
    {
      question: "How do I cancel my Premium subscription?",
      answer: "You can cancel your Premium subscription at any time from your account settings. Your Premium features will remain active until the end of your billing period.",
    },
    {
      question: "Is my viewing data private?",
      answer: "Absolutely. We take privacy seriously and do not collect or store any personal viewing data. The extension works locally on your device, and we only collect anonymous usage statistics to improve our service.",
    },
  ];

let currentYear = $state(new Date().getFullYear());
  
  // Animation state
  let freeCardHover = false;
  let premiumCardHover = false;




    // Show the skip button after a delay
    setTimeout(() => {
      showSkipButton = true;
    }, 1200);
    
    return () => {
      window.removeEventListener('scroll', handleScroll);
      clearInterval(pulseInterval);
      clearInterval(diamondInterval);
      clearInterval(hexagonInterval);
      clearInterval(blobInterval);
      clearInterval(progressInterval);
      clearInterval(scrollArrowInterval);
    };
  });
</script>

<div class="min-h-screen overflow-hidden bg-black text-white">
  <!-- Abstract geometric shapes background -->
  <div class="fixed inset-0 -z-10 opacity-30">
    <div class="absolute top-20 left-10 w-64 h-64 rounded-full bg-purple-500"
         style="transform: translate({$pulseScale.x}px, {$pulseScale.y}px) scale({$pulseScale.scale});">
    </div>
    
    <div class="absolute top-40 right-20 w-80 h-80 bg-yellow-400"
         style="clip-path: polygon(50% 0%, 100% 50%, 50% 100%, 0% 50%); transform: rotate({$diamondRotate.rotate}deg) scale({$diamondRotate.scale});">
    </div>
    
    <div class="absolute bottom-20 left-1/4 w-72 h-72 bg-cyan-500"
         style="clip-path: polygon(25% 0%, 75% 0%, 100% 50%, 75% 100%, 25% 100%, 0% 50%); transform: rotate({$hexagonRotate.rotate}deg) translateY({$hexagonRotate.y}px);">
    </div>
    
    <div class="absolute top-1/2 right-1/4 w-56 h-56 bg-pink-500"
         style="border-radius: {$blobRadius.borderRadius}; transform: rotate({$blobRadius.rotate}deg);">
    </div>
  </div>

  <!-- Header -->
  <header class="fixed top-0 left-0 right-0 z-50 backdrop-blur-md bg-black/50 border-b border-white/10">
    <div class="container mx-auto px-4 py-4 flex justify-between items-center">
      <div class="flex items-center gap-2 font-bold text-xl" in:fade={{ duration: 500, delay: 0 }}>
        <FastForward class="text-yellow-400" />
        <span class="bg-clip-text text-transparent bg-gradient-to-r from-yellow-400 via-pink-500 to-cyan-500">
          SponsorSkip
        </span>
      </div>
      
      <div class="hidden md:flex gap-8 items-center" in:fade={{ duration: 500, delay: 200 }}>
        <a href="#features" class="hover:text-yellow-400 transition-colors">Features</a>
        <a href="#how-it-works" class="hover:text-pink-500 transition-colors">How It Works</a>
        <a href="#pricing" class="hover:text-cyan-500 transition-colors">Pricing</a>
        <a href="#faq" class="hover:text-purple-500 transition-colors">FAQ</a>
      </div>
      
      <div in:fade={{ duration: 500, delay: 400 }}>
        <button class="bg-gradient-to-r from-yellow-400 to-pink-500 hover:from-yellow-500 hover:to-pink-600 text-black font-bold px-4 py-2 rounded-md">
          Install Extension
        </button>
      </div>
    </div>
  </header>

  <!-- Hero Section -->
  <section class="relative min-h-screen flex items-center justify-center pt-20">
    <div class="container mx-auto px-4 py-20 flex flex-col items-center text-center">
      <h1 class="text-5xl md:text-7xl font-black mb-6 leading-tight" in:fade={{ duration: 800, delay: 0 }}>
        <span class="block">Skip The</span>
        <span class="bg-clip-text text-transparent bg-gradient-to-r from-yellow-400 via-pink-500 to-cyan-500">
          Sponsorships
        </span>
        <span class="block">Save Your Time</span>
      </h1>
      
      <p class="text-xl md:text-2xl mb-10 max-w-2xl text-gray-300" in:fade={{ duration: 800, delay: 200 }}>
        The Chrome extension that automatically detects and skips sponsored segments in YouTube videos, so you can enjoy uninterrupted content.
      </p>
      
      <div class="flex flex-col sm:flex-row gap-4 mb-16" in:fade={{ duration: 800, delay: 400 }}>
        <button class="bg-gradient-to-r from-yellow-400 to-pink-500 hover:from-yellow-500 hover:to-pink-600 text-black font-bold text-lg px-8 py-6 rounded-md">
          Install Free Version
        </button>
        <button class="border-2 border-cyan-500 text-cyan-500 hover:bg-cyan-500/20 font-bold text-lg px-8 py-6 rounded-md">
          Upgrade to Premium
        </button>
      </div>
      
      <div class="relative w-full max-w-4xl rounded-xl overflow-hidden shadow-2xl border-4 border-white/10" in:fade={{ duration: 1000, delay: 600 }}>
        <div class="relative pt-[56.25%] bg-gray-900">
          <div class="absolute inset-0 flex items-center justify-center">
            <div class="w-full h-full bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center">
              <div class="relative w-3/4 h-3/4 border-2 border-white/20 rounded-lg overflow-hidden">
                <div class="absolute top-0 left-0 right-0 h-12 bg-gray-800 flex items-center px-4">
                  <div class="flex gap-2">
                    <div class="w-3 h-3 rounded-full bg-red-500"></div>
                    <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
                    <div class="w-3 h-3 rounded-full bg-green-500"></div>
                  </div>
                  <div class="ml-4 bg-gray-700 rounded-full px-4 py-1 text-xs text-gray-300 flex-1 max-w-md">
                    youtube.com/watch?v=example
                  </div>
                </div>
                <div class="absolute top-12 left-0 right-0 bottom-0 bg-gray-700 flex items-center justify-center">
                  <div class="relative w-full h-full">
                    <div class="absolute inset-0 bg-gradient-to-br from-gray-800 to-gray-900 flex items-center justify-center">
                      <img src="/api/placeholder/600/400" alt="YouTube video player" class="w-full h-full object-cover opacity-50" />
                      
                      {#if showSkipButton}
                        <button 
                          class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-yellow-400 text-black font-bold px-6 py-3 rounded-full flex items-center gap-2 shadow-lg"
                          in:scale={{ duration: 500 }}
                        >
                          <FastForward size={24} />
                          <span>Sponsorship Skipped!</span>
                        </button>
                      {/if}
                    </div>
                    <div class="absolute bottom-0 left-0 right-0 h-10 bg-gray-800 flex items-center px-4">
                      <div class="w-full bg-gray-600 h-1 rounded-full overflow-hidden">
                        <div class="h-full bg-red-500" style="width: {$progressBar}%;"></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="mt-12" in:fade={{ duration: 1000, delay: 1200 }}>
        <a href="#features" class="flex flex-col items-center text-gray-400 hover:text-white transition-colors">
          <span>Scroll to explore</span>
          <div style="transform: translateY({$scrollArrow}px);">
            <ChevronDown size={24} />
          </div>
        </a>
      </div>
    </div>
  </section>

  <!-- Features Section -->
  <section id="features" class="py-20 relative">
    <div class="container mx-auto px-4">
      <div class="text-center mb-16" in:fade={{ duration: 800 }}>
        <h2 class="text-4xl md:text-5xl font-bold mb-4">
          <span class="bg-clip-text text-transparent bg-gradient-to-r from-yellow-400 to-pink-500">
            Powerful Features
          </span>
        </h2>
        <p class="text-xl text-gray-300 max-w-2xl mx-auto">
          SponsorSkip comes packed with everything you need to enjoy YouTube without interruptions
        </p>
      </div>
      
      <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        {#each [
          {
            icon: FastForward,
            color: "text-yellow-400",
            title: "Automatic Detection",
            description: "Our AI-powered algorithm automatically detects sponsored segments in videos with high accuracy."
          },
          {
            icon: Zap,
            color: "text-pink-500",
            title: "Instant Skipping",
            description: "Skip sponsored content instantly without having to manually forward through videos."
          },
          {
            icon: Shield,
            color: "text-cyan-500",
            title: "Privacy Focused",
            description: "We don't collect any personal data. Your viewing habits remain completely private."
          },
          {
            icon: Clock,
            color: "text-purple-500",
            title: "Time Saved Counter",
            description: "See how much time you've saved by skipping sponsored segments across all videos."
          },
          {
            icon: Award,
            color: "text-green-400",
            title: "Premium Features",
            description: "Upgrade to Premium for unlimited skips, custom settings, and priority updates."
          },
          {
            icon: ArrowRight,
            color: "text-orange-400",
            title: "Cross-Platform",
            description: "Works seamlessly across all your devices with Chrome browser installed."
          }
        ] as feature, i (feature.title)}
          <!-- svelte-ignore a11y_no_static_element_interactions -->
          <div 
            class="bg-gradient-to-br from-gray-900 to-black border border-white/10 rounded-xl p-6 hover:shadow-xl hover:shadow-pink-500/10 transition-all duration-300"
            in:fade={{ duration: 500, delay: i * 100 }}
            on:mouseenter={() => {}}
            on:mouseleave={() => {}}
          >
            <div class="mb-4">
              <svelte:component this={feature.icon} size={32} class={feature.color} />
            </div>
            <h3 class="text-xl font-bold mb-2">{feature.title}</h3>
            <p class="text-gray-400">{feature.description}</p>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- How It Works Section -->
  <section id="how-it-works" class="py-20 relative bg-gradient-to-b from-black to-gray-900">
    <div class="container mx-auto px-4">
      <div class="text-center mb-16" in:fade={{ duration: 800 }}>
        <h2 class="text-4xl md:text-5xl font-bold mb-4">
          <span class="bg-clip-text text-transparent bg-gradient-to-r from-cyan-500 to-purple-500">
            How It Works
          </span>
        </h2>
        <p class="text-xl text-gray-300 max-w-2xl mx-auto">
          SponsorSkip is incredibly easy to use. Just install and forget about it!
        </p>
      </div>
      
      <div class="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
        {#each [
          {
            step: "01",
            title: "Install the Extension",
            description: "Add SponsorSkip to your Chrome browser with just one click from the Chrome Web Store."
          },
          {
            step: "02",
            title: "Watch YouTube Videos",
            description: "Continue watching YouTube as you normally would. No configuration needed."
          },
          {
            step: "03",
            title: "Enjoy Uninterrupted Content",
            description: "SponsorSkip automatically detects and skips sponsored segments in your videos."
          }
        ] as step, i (step.step)}
          <div class="relative" in:fade={{ duration: 500, delay: i * 200 }}>
            <div class="bg-gradient-to-br from-gray-800 to-gray-900 border border-white/10 rounded-xl p-6 h-full">
              <div class="text-6xl font-black text-white/10 mb-4">{step.step}</div>
              <h3 class="text-xl font-bold mb-2">{step.title}</h3>
              <p class="text-gray-400">{step.description}</p>
            </div>
            
            {#if i < 2}
              <div 
                class="hidden md:block absolute top-1/2 right-0 transform translate-x-1/2 -translate-y-1/2 text-pink-500"
                style="transform: translateX(calc(50% + {Math.sin(Date.now() / 1500) * 10}px)) translateY(-50%);"
              >
                <ArrowRight size={32} />
              </div>
            {/if}
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- Pricing Section -->
  <section id="pricing" class="py-20 relative">
    <div class="container mx-auto px-4">
      <div class="text-center mb-16" in:fade={{ duration: 800 }}>
        <h2 class="text-4xl md:text-5xl font-bold mb-4">
          <span class="bg-clip-text text-transparent bg-gradient-to-r from-yellow-400 via-pink-500 to-cyan-500">
            Choose Your Plan
          </span>
        </h2>
        <p class="text-xl text-gray-300 max-w-2xl mx-auto">
          Start with our free plan or upgrade to Premium for unlimited skipping power
        </p>
      </div>
      
      <div class="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto">
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div 
          class="bg-gradient-to-br from-gray-900 to-black border border-white/10 rounded-xl overflow-hidden"
          in:fade={{ duration: 500,  }}
          on:mouseenter={() => {}}
          on:mouseleave={() => {}}
        >
          <div class="p-8">
            <h3 class="text-2xl font-bold mb-2">Free</h3>
            <div class="text-4xl font-bold mb-4">
              $0<span class="text-lg text-gray-400">/month</span>
            </div>
            <p class="text-gray-400 mb-6">Perfect for casual YouTube viewers</p>
            <ul class="space-y-3 mb-8">
              {#each ["Up to 10 skips per day", "Basic detection algorithm", "Standard support", "Single device"] as feature}
                <li class="flex items-center gap-2">
                  <div class="text-green-400">✓</div>
                  <span>{feature}</span>
                </li>
              {/each}
            </ul>
            <button class="w-full bg-white hover:bg-gray-200 text-black font-bold py-2 px-4 rounded-md">
              Install Free Version
            </button>
          </div>
        </div>
        
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div 
          class="bg-gradient-to-br from-yellow-500/20 via-pink-500/20 to-cyan-500/20 border-2 border-pink-500 rounded-xl overflow-hidden relative"
          in:fade={{ duration: 500, }}
          on:mouseenter={() => {}}
          on:mouseleave={() => {}}
        >
          <div class="absolute top-0 right-0 bg-gradient-to-r from-yellow-400 to-pink-500 text-black font-bold px-4 py-1 text-sm">
            MOST POPULAR
          </div>
          <div class="p-8">
            <h3 class="text-2xl font-bold mb-2">Premium</h3>
            <div class="text-4xl font-bold mb-4">
              $4.99<span class="text-lg text-gray-400">/month</span>
            </div>
            <p class="text-gray-400 mb-6">For serious YouTube enthusiasts</p>
            <ul class="space-y-3 mb-8">
              {#each [
                "Unlimited skips",
                "Advanced AI detection",
                "Priority support",
                "Use on multiple devices",
                "Custom skip settings",
                "No advertisements"
              ] as feature}
                <li class="flex items-center gap-2">
                  <div class="text-pink-500">✓</div>
                  <span>{feature}</span>
                </li>
              {/each}
            </ul>
            <button class="w-full bg-gradient-to-r from-yellow-400 to-pink-500 hover:from-yellow-500 hover:to-pink-600 text-black font-bold py-2 px-4 rounded-md">
              Upgrade to Premium
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>
</div>

<!-- FAQ Section -->
  <section id="faq" class="py-20 relative bg-gradient-to-b from-gray-900 to-black">
    <div class="container mx-auto px-4">
      <div class="text-center mb-16" in:fade={{ duration: 800}}>
        <h2 class="text-4xl md:text-5xl font-bold mb-4">
          <span class="bg-clip-text text-transparent bg-gradient-to-r from-purple-500 to-cyan-500">
            Frequently Asked Questions
          </span>
        </h2>
        <p class="text-xl text-gray-300 max-w-2xl mx-auto">Got questions? We've got answers</p>
      </div>
      
      <div class="max-w-3xl mx-auto space-y-6">
    {#each [
    {
      question: "How does SponsorSkip detect sponsored segments?",
      answer: "SponsorSkip uses a combination of machine learning algorithms and community-contributed data to identify sponsored segments in videos with high accuracy. Our system analyzes audio, visual cues, and patterns to determine when a sponsorship begins and ends.",
    },
    {
      question: "Will SponsorSkip work on all YouTube videos?",
      answer: "SponsorSkip works on most YouTube videos, but detection accuracy may vary. Our system continuously improves as more users contribute data. Premium users get access to our most advanced detection algorithms for better accuracy.",
    },
    {
      question: "Can I customize which types of segments to skip?",
      answer: "Yes! Premium users can customize which types of segments to skip, including sponsorships, intros, outros, subscription reminders, and more. Free users have basic customization options.",
    },
    {
      question: "How do I cancel my Premium subscription?",
      answer: "You can cancel your Premium subscription at any time from your account settings. Your Premium features will remain active until the end of your billing period.",
    },
    {
      question: "Is my viewing data private?",
      answer: "Absolutely. We take privacy seriously and do not collect or store any personal viewing data. The extension works locally on your device, and we only collect anonymous usage statistics to improve our service.",
    },
  ] as faq, index}

          <div 
            class="bg-gradient-to-br from-gray-800 to-gray-900 border border-white/10 rounded-xl overflow-hidden"
            in:fly={{ duration: 500, delay: index * 100, y: 20 }}
          >
            <details class="group">
              <summary class="flex justify-between items-center p-6 cursor-pointer list-none">
                <h3 class="text-xl text-white font-medium">{faq.question}</h3>
                <div class="transition-transform duration-300 group-open:rotate-180">
                  <ChevronDown size={20} />
                </div>
              </summary>
              <div class="px-6 pb-6 pt-0">
                <p class="text-gray-400">{faq.answer}</p>
              </div>
            </details>
          </div>
        {/each}
      </div>
    </div>
  </section>
  
  <!-- CTA Section -->
  <section class="py-20 relative">
    <div class="container mx-auto px-4">
      <div 
        class="max-w-4xl mx-auto bg-gradient-to-br from-gray-900 to-black border border-white/10 rounded-2xl p-8 md:p-12 text-center relative overflow-hidden"
        in:fade={{ duration: 800 }}
      >
        <div 
          class="absolute -top-24 -right-24 w-64 h-64 bg-yellow-400/20 rounded-full blur-3xl" 
          style="animation: pulse1 8s infinite;"
        ></div>
        <div 
          class="absolute -bottom-32 -left-32 w-80 h-80 bg-pink-500/20 rounded-full blur-3xl"
          style="animation: pulse2 10s 1s infinite;"
        ></div>
        
        <h2 class="text-3xl md:text-4xl text-white font-bold mb-4 relative z-10">Ready to Skip the Boring Parts?</h2>
        <p class="text-xl text-gray-300 mb-8 relative z-10 max-w-2xl mx-auto">
          Join thousands of users who save hours every month by skipping sponsored content on YouTube.
        </p>
        
        <div class="flex flex-col sm:flex-row gap-4 justify-center relative z-10">
          <button class="bg-gradient-to-r from-yellow-400 to-pink-500 hover:from-yellow-500 hover:to-pink-600 text-black font-bold text-lg px-8 py-6 rounded-md">
            Install Free Version
          </button>
          <button class="border-2 border-cyan-500 text-cyan-500 hover:bg-cyan-500/20 font-bold text-lg px-8 py-6 rounded-md">
            Upgrade to Premium
          </button>
        </div>
      </div>
    </div>
  </section>
  
  <!-- Footer -->
  <footer class="py-12 border-t border-white/10 bg-black">
    <div class="container mx-auto px-4">
      <div class="grid md:grid-cols-4 gap-8">
        <div>
          <div class="flex items-center gap-2 font-bold text-xl mb-4">
            <FastForward class="text-yellow-400" />
            <span class="bg-clip-text text-transparent bg-gradient-to-r from-yellow-400 via-pink-500 to-cyan-500">
              SponsorSkip
            </span>
          </div>
          <p class="text-gray-400 mb-4">
            The ultimate browser extension for skipping sponsored content on YouTube videos.
          </p>
          <div class="flex gap-4">
            <!-- svelte-ignore a11y_consider_explicit_label -->
            <a href="#" class="text-gray-400 hover:text-white transition-colors">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="lucide lucide-twitter"
              >
                <path d="M22 4s-.7 2.1-2 3.4c1.6 10-9.4 17.3-18 11.6 2.2.1 4.4-.6 6-2C3 15.5.5 9.6 3 5c2.2 2.6 5.6 4.1 9 4-.9-4.2 4-6.6 7-3.8 1.1 0 3-1.2 3-1.2z" />
              </svg>
            </a>
            <!-- svelte-ignore a11y_consider_explicit_label -->
            <a href="#" class="text-gray-400 hover:text-white transition-colors">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="lucide lucide-facebook"
              >
                <path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3z" />
              </svg>
            </a>
            <!-- svelte-ignore a11y_consider_explicit_label -->
            <a href="#" class="text-gray-400 hover:text-white transition-colors">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="lucide lucide-instagram"
              >
                <rect width="20" height="20" x="2" y="2" rx="5" ry="5" />
                <path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z" />
                <line x1="17.5" x2="17.51" y1="6.5" y2="6.5" />
              </svg>
            </a>
          </div>
        </div>
        
        <div>
          <h3 class="font-bold text-lg mb-4">Product</h3>
          <ul class="space-y-2">
            <li>
              <a href="#features" class="text-gray-400 hover:text-white transition-colors">
                Features
              </a>
            </li>
            <li>
              <a href="#pricing" class="text-gray-400 hover:text-white transition-colors">
                Pricing
              </a>
            </li>
            <li>
              <a href="#" class="text-gray-400 hover:text-white transition-colors">
                Testimonials
              </a>
            </li>
            <li>
              <a href="#faq" class="text-gray-400 hover:text-white transition-colors">
                FAQ
              </a>
            </li>
          </ul>
        </div>
        
        <div>
          <h3 class="font-bold text-lg mb-4">Resources</h3>
          <ul class="space-y-2">
            <li>
              <a href="#" class="text-gray-400 hover:text-white transition-colors">
                Blog
              </a>
            </li>
            <li>
              <a href="#" class="text-gray-400 hover:text-white transition-colors">
                Documentation
              </a>
            </li>
            <li>
              <a href="#" class="text-gray-400 hover:text-white transition-colors">
                Community
              </a>
            </li>
            <li>
              <a href="#" class="text-gray-400 hover:text-white transition-colors">
                Support
              </a>
            </li>
          </ul>
        </div>
        
        <div>
          <h3 class="font-bold text-lg mb-4">Company</h3>
          <ul class="space-y-2">
            <li>
              <a href="#" class="text-gray-400 hover:text-white transition-colors">
                About Us
              </a>
            </li>
            <li>
              <a href="#" class="text-gray-400 hover:text-white transition-colors">
                Careers
              </a>
            </li>
            <li>
              <a href="#" class="text-gray-400 hover:text-white transition-colors">
                Privacy Policy
              </a>
            </li>
            <li>
              <a href="#" class="text-gray-400 hover:text-white transition-colors">
                Terms of Service
              </a>
            </li>
          </ul>
        </div>
      </div>
      
      <div class="mt-12 pt-8 border-t border-white/10 text-center text-gray-500">
        <p>© {new Date().getFullYear()} SponsorSkip. All rights reserved.</p>
      </div>
    </div>
  </footer>
<!-- </div> -->

<style>
  @keyframes pulse1 {
    0% { transform: scale(1); opacity: 0.2; }
    50% { transform: scale(1.2); opacity: 0.3; }
    100% { transform: scale(1); opacity: 0.2; }
  }
  
  @keyframes pulse2 {
    0% { transform: scale(1); opacity: 0.2; }
    50% { transform: scale(1.3); opacity: 0.3; }
    100% { transform: scale(1); opacity: 0.2; }
  }
  
  :global(.group-open) :global(.group-open\:rotate-180) {
    transform: rotate(180deg);
  }
 /* Fix for CTA section background */
  :global(section.py-20.relative) {
    background-color: black;
  }
  
  :global(.max-w-4xl.mx-auto.bg-gradient-to-br) {
    background: linear-gradient(to bottom right, #111827, #000000) !important;
    border: 1px solid rgba(255, 255, 255, 0.1) !important;
  } 
  .container {
    width: 100%;
    margin-left: auto;
    margin-right: auto;
    padding-left: 1rem;
    padding-right: 1rem;
  }
  
  @media (min-width: 640px) {
    .container { max-width: 640px; }
  }
  
  @media (min-width: 768px) {
    .container { max-width: 768px; }
  }
  
  @media (min-width: 1024px) {
    .container { max-width: 1024px; }
  }
  
  @media (min-width: 1280px) {
    .container { max-width: 1280px; }
  }
 
</style>