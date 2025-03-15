<!-- src/routes/+page.svelte -->
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
        <span class="font-black text-2xl tracking-tight">
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
  <section id="features" class="py-20 border-b-4 border-black relative">
    <div class="container mx-auto px-4">
      <div class="text-center mb-16" in:fade={{ duration: 500 }}>
        <h2 class="text-5xl font-black mb-4">
          AWESOME <span class="text-blue-500">FEATURES</span>
        </h2>
        <p class="text-xl max-w-2xl mx-auto">
          Our extension is packed with powerful features to enhance your YouTube experience.
        </p>
      </div>

      <div class="grid md:grid-cols-3 gap-8">
        {#each [
          {
            icon: FastForward,
            title: "Auto-Skip Sponsorships",
            description:
              "Automatically detects and skips sponsored segments in videos so you don't have to manually skip them.",
            color: "bg-red-500",
          },
          {
            icon: Zap,
            title: "Lightning Fast",
            description:
              "Minimal impact on performance. Works silently in the background without slowing down your browsing.",
            color: "bg-yellow-400",
          },
          {
            icon: Clock,
            title: "Time Saved Tracker",
            description:
              "See exactly how much time you've saved by skipping sponsorships across all your watched videos.",
            color: "bg-blue-500",
          },
        ] as feature, index}
          <div 
            class="border-4 border-black bg-white p-6 relative overflow-hidden group"
            in:fade={{ duration: 500, delay: index * 100 }}
          >
            <div
              class="absolute top-0 right-0 w-20 h-20 {feature.color} border-l-4 border-b-4 border-black -mr-10 -mt-10 transition-all duration-300 group-hover:mr-0 group-hover:mt-0"
            ></div>
            <div class="relative z-10">
              <div class="{feature.color} text-black p-4 inline-block border-2 border-black mb-4">
                <svelte:component this={feature.icon} class="w-10 h-10" />
              </div>
              <h3 class="text-2xl font-bold mb-2">{feature.title}</h3>
              <p>{feature.description}</p>
            </div>
          </div>
        {/each}
      </div>

      <div 
        class="mt-16 p-8 border-4 border-black bg-purple-100 relative"
        in:fade={{ duration: 500 }}
      >
        <div class="absolute -top-5 -left-5 w-10 h-10 bg-yellow-300 border-4 border-black"></div>
        <div class="absolute -bottom-5 -right-5 w-10 h-10 bg-blue-400 border-4 border-black rounded-full"></div>

        <div class="grid md:grid-cols-2 gap-8 items-center">
          <div>
            <h3 class="text-3xl font-bold mb-4">Smart Detection Technology</h3>
            <p class="mb-4">
              Our advanced algorithm recognizes sponsorship segments with incredible accuracy, even when creators
              try to disguise them.
            </p>
            <ul class="space-y-2">
              {#each [
                "Recognizes common sponsorship phrases",
                "Detects visual sponsorship indicators",
                "Learns from user feedback",
                "Updates in real-time",
              ] as item}
                <li class="flex items-start">
                  <div class="bg-green-500 text-white p-1 mr-2 mt-1">
                    <ChevronRight class="w-4 h-4" />
                  </div>
                  {item}
                </li>
              {/each}
            </ul>
          </div>
          <div class="relative">
            <div class="border-4 border-black bg-white p-4 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
              <img
                src="/placeholder.svg?height=300&width=400"
                alt="Smart detection visualization"
                class="w-full border-2 border-black"
              />
            </div>
            <div class="absolute -bottom-4 -left-4 w-full h-full border-4 border-black bg-red-200 -z-10"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Memphis design elements -->
    <div class="absolute top-20 right-10 w-16 h-16 bg-yellow-300 border-4 border-black transform rotate-45 hidden lg:block"></div>
    <div class="absolute bottom-40 left-10 w-10 h-10 bg-blue-400 border-4 border-black rounded-full hidden lg:block"></div>
  </section>

  <!-- Pricing Section -->
  <section id="pricing" class="py-20 border-b-4 border-black relative bg-gradient-to-b from-white to-gray-100">
    <div class="container mx-auto px-4">
      <div 
        class="text-center mb-16"
        in:fade={{ duration: 500 }}
      >
        <h2 class="text-5xl font-black mb-4">
          CHOOSE YOUR <span class="text-purple-600">PLAN</span>
        </h2>
        <p class="text-xl max-w-2xl mx-auto">Upgrade to Premium for unlimited skips and advanced features.</p>
      </div>

      <div class="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto">
        {#each [
          {
            title: "Free",
            price: "$0",
            period: "forever",
            description: "Basic sponsorship skipping for casual YouTube viewers",
            features: [
              "Skip up to 50 sponsorships per month",
              "Basic sponsorship detection",
              "Time saved tracker",
              "Works on all YouTube videos",
            ],
            buttonText: "Install Now",
            buttonColor: "bg-black",
            popular: false,
          },
          {
            title: "Premium",
            price: "$4.99",
            period: "per month",
            description: "Unlimited skipping and advanced features for power users",
            features: [
              "Unlimited sponsorship skipping",
              "Advanced detection algorithm",
              "Custom skip rules and preferences",
              "Skip intros, outros & reminders",
              "Detailed analytics dashboard",
              "Priority support",
            ],
            buttonText: "Go Premium",
            buttonColor: "bg-purple-600",
            popular: true,
          },
        ] as plan, index}
          <div
            class="border-4 border-black bg-white p-8 relative {plan.popular ? 'transform md:-translate-y-4' : ''}"
            in:fade={{ duration: 500, delay: index * 100 }}
          >
            {#if plan.popular}
              <div class="absolute -top-4 -right-4 bg-yellow-400 text-black font-bold py-1 px-4 border-4 border-black">
                POPULAR
              </div>
            {/if}

            <h3 class="text-3xl font-bold mb-2">{plan.title}</h3>
            <div class="flex items-end mb-4">
              <span class="text-4xl font-black">{plan.price}</span>
              <span class="text-gray-600 ml-1">/{plan.period}</span>
            </div>
            <p class="mb-6 text-gray-600">{plan.description}</p>

            <ul class="space-y-3 mb-8">
              {#each plan.features as feature}
                <li class="flex items-start">
                  <div class="bg-green-500 text-white p-1 mr-2">
                    <ChevronRight class="w-4 h-4" />
                  </div>
                  {feature}
                </li>
              {/each}
            </ul>

            <button
              class="{plan.buttonColor} text-white font-bold py-3 px-8 w-full border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-1 hover:translate-y-1 transition-all transform hover:scale-105 active:scale-95"
            >
              {plan.buttonText}
            </button>
          </div>
        {/each}
      </div>

      <div 
        class="mt-16 text-center"
        in:fade={{ duration: 500 }}
      >
        <p class="text-xl mb-6">Not convinced yet? Try Premium free for 7 days!</p>
        <button
          class="bg-gradient-to-r from-purple-600 to-blue-500 text-white font-bold py-3 px-8 border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-1 hover:translate-y-1 transition-all transform hover:scale-105 active:scale-95"
        >
          Start Free 
        </button>
      </div>
    </div>

    <!-- Memphis design elements -->
    <div class="absolute top-40 left-10 w-20 h-4 bg-red-500 border-4 border-black hidden lg:block"></div>
    <div class="absolute bottom-20 right-20 w-12 h-12 bg-blue-400 border-4 border-black rounded-full hidden lg:block"></div>
  </section>

  <!-- Testimonials -->
  <section class="py-20 border-b-4 border-black relative">
    <div class="container mx-auto px-4">
      <div 
        class="text-center mb-16"
        in:fade={{ duration: 500 }}
      >
        <h2 class="text-5xl font-black mb-4">
          WHAT PEOPLE <span class="text-green-500">SAY</span>
        </h2>
        <p class="text-xl max-w-2xl mx-auto">Join thousands of happy users who save time every day.</p>
      </div>

      <div class="grid md:grid-cols-3 gap-8">
        {#each [
          {
            name: "Alex Johnson",
            role: "Tech Enthusiast",
            quote: "This extension has saved me hours of my life. No more sitting through boring sponsorships!",
            color: "bg-red-100",
          },
          {
            name: "Sarah Miller",
            role: "Daily YouTube User",
            quote: "The Premium version is worth every penny. I've saved over 3 hours this month alone.",
            color: "bg-blue-100",
          },
          {
            name: "Michael Chen",
            role: "Content Creator",
            quote: "As someone who watches a lot of tutorials, this extension is a game-changer for productivity.",
            color: "bg-yellow-100",
          },
        ] as testimonial, index}
          <div
            class="border-4 border-black {testimonial.color} p-6 relative"
            in:fade={{ duration: 500, delay: index * 100 }}
          >
            <div class="absolute -top-5 -left-5 w-10 h-10 bg-white border-4 border-black rounded-full flex items-center justify-center text-2xl font-bold">
              "
            </div>
            <p class="mb-6 text-lg italic">{testimonial.quote}</p>
            <div class="flex items-center">
              <div class="w-12 h-12 bg-gray-300 border-2 border-black rounded-full mr-4"></div>
              <div>
                <div class="font-bold">{testimonial.name}</div>
                <div class="text-sm text-gray-600">{testimonial.role}</div>
              </div>
            </div>
          </div>
        {/each}
      </div>

      <div 
        class="mt-16 p-8 border-4 border-black bg-green-100"
        in:fade={{ duration: 500 }}
      >
        <div class="flex flex-col md:flex-row items-center justify-between">
          <div class="mb-6 md:mb-0">
            <h3 class="text-3xl font-bold mb-2">Ready to start skipping?</h3>
            <p class="text-xl">Join over 100,000 users saving time every day.</p>
          </div>
          <div class="flex gap-4">
            <button
              class="bg-black text-white font-bold py-3 px-8 border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-1 hover:translate-y-1 transition-all transform hover:scale-105 active:scale-95"
            >
              Install Free
            </button>
            <button
              class="bg-purple-600 text-white font-bold py-3 px-8 border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-1 hover:translate-y-1 transition-all transform hover:scale-105 active:scale-95"
            >
              Go Premium
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- FAQ Section -->
  <section id="faq" class="py-20 border-b-4 border-black">
    <div class="container mx-auto px-4">
      <div 
        class="text-center mb-16"
        in:fade={{ duration: 500 }}
      >
        <h2 class="text-5xl font-black mb-4">
          FAQ<span class="text-orange-500">s</span>
        </h2>
        <p class="text-xl max-w-2xl mx-auto">Got questions? We've got answers.</p>
      </div>

      <div class="max-w-3xl mx-auto space-y-6">
        {#each [
          {
            question: "How does the extension detect sponsorships?",
            answer:
              "Our extension uses a combination of machine learning algorithms and community-reported data to identify sponsorship segments in videos. It recognizes patterns in speech, visual cues, and common sponsorship phrases.",
          },
          {
            question: "What's the difference between Free and Premium?",
            answer:
              "The Free version allows you to skip up to 50 sponsorships per month, while Premium offers unlimited skipping, advanced detection, custom skip rules, and additional features like intro/outro skipping and detailed analytics.",
          },
          {
            question: "Will this slow down my browser?",
            answer:
              "No, our extension is designed to be lightweight and efficient. It runs in the background with minimal impact on your browsing experience or computer performance.",
          },
          {
            question: "Can I customize what gets skipped?",
            answer:
              "Yes, Premium users can set custom rules for what types of segments to skip (sponsorships, intros, outros, etc.) and even create channel-specific settings.",
          },
          {
            question: "How do I cancel my Premium subscription?",
            answer:
              "You can cancel your Premium subscription at any time from your account settings. Your Premium features will remain active until the end of your billing period.",
          },
        ] as faq, index}
          <div
            class="border-4 border-black bg-white overflow-hidden"
            in:fade={{ duration: 500, delay: index * 100 }}
          >
            <div class="bg-gray-100 border-b-4 border-black p-4 font-bold text-lg flex justify-between items-center">
              {faq.question}
              <div class="w-6 h-6 bg-black text-white flex items-center justify-center">+</div>
            </div>
            <div class="p-4">
              <p>{faq.answer}</p>
            </div>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- CTA Section -->
  <section class="py-20 relative overflow-hidden">
    <div class="container mx-auto px-4 relative z-10">
      <div 
        class="max-w-4xl mx-auto text-center bg-white border-4 border-black p-10 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]"
        in:scale={{ duration: 500, start: 0.9 }}
      >
        <h2 class="text-5xl font-black mb-6">
          STOP WASTING <span class="text-red-500">TIME</span>
        </h2>
        <p class="text-xl mb-8">
          The average YouTube user wastes over 5 hours per month watching sponsorships. Get that time back with
          SkipIt!
        </p>

        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <button
            class="bg-black text-white font-bold py-3 px-8 border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-1 hover:translate-y-1 transition-all transform hover:scale-105 active:scale-95"
          >
            Install Free
          </button>
          <button
            class="bg-purple-600 text-white font-bold py-3 px-8 border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-1 hover:translate-y-1 transition-all flex items-center justify-center transform hover:scale-105 active:scale-95"
          >
            Go Premium <CreditCard class="ml-2 w-5 h-5" />
          </button>
        </div>

        <div class="mt-8 flex items-center justify-center gap-4">
          <div class="flex">
            {#each Array(5) as _}
              <Award class="w-6 h-6 text-yellow-500" />
            {/each}
          </div>
          <span class="font-bold">4.9/5 from 2,000+ reviews</span>
        </div>
      </div>
    </div>

    <!-- Memphis design background -->
    <div class="absolute inset-0 -z-10">
      <div class="absolute top-10 left-10 w-20 h-20 bg-yellow-300 border-4 border-black transform rotate-45"></div>
      <div class="absolute top-40 right-20 w-16 h-16 bg-blue-400 border-4 border-black rounded-full"></div>
      <div class="absolute bottom-20 left-1/3 w-24 h-8 bg-red-500 border-4 border-black"></div>
      <div class="absolute bottom-40 right-1/4 w-12 h-12 bg-purple-400 border-4 border-black transform rotate-12"></div>
    </div>
  </section>

  <!-- Footer -->
  <footer class="bg-black text-white py-12 border-t-4 border-white">
    <div class="container mx-auto px-4">
      <div class="grid md:grid-cols-4 gap-8">
        <div>
          <div class="flex items-center gap-2 mb-4">
            <div class="bg-red-500 w-8 h-8 rounded-full flex items-center justify-center">
              <FastForward class="text-white w-5 h-5" />
            </div>
            <span class="font-black text-2xl tracking-tight">
              SKIP<span class="text-red-500">IT</span>
            </span>
          </div>
          <p class="text-gray-400">Save time and enjoy uninterrupted YouTube content.</p>
        </div>

        <div>
          <h4 class="font-bold text-lg mb-4">Product</h4>
          <ul class="space-y-2">
            <li>
              <a href="#features" class="hover:text-red-500 transition-colors">
                Features
              </a>
            </li>
            <li>
              <a href="#pricing" class="hover:text-red-500 transition-colors">
                Pricing
              </a>
            </li>
            
          </ul>
        </div>
      </div>
    </div>
  </footer>
</div>
