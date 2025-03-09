<!-- src/routes/+page.svelte -->
<script>
  import { onMount } from 'svelte';
  import { FastForward, ChevronRight, Clock, Zap, Check } from 'lucide-svelte';
  import { fly, scale } from 'svelte/transition';
  import { tweened } from 'svelte/motion';
  import { cubicOut } from 'svelte/easing';

  let isScrolled = $state(false);

  // Sample testimonial data
  const testimonials = [
    {
      name: "Alex Johnson",
      role: "Tech Enthusiast",
      comment: "This extension has saved me hours of time. I watch YouTube daily and no longer have to manually skip through sponsorships!",
      avatar: "/placeholder.svg?height=80&width=80",
    },
    {
      name: "Sarah Miller",
      role: "Content Creator",
      comment: "As someone who watches tutorials all day, SponsorSkip has been a game changer. The premium version is absolutely worth it.",
      avatar: "/placeholder.svg?height=80&width=80",
    },
    {
      name: "Michael Chen",
      role: "Student",
      comment: "I upgraded to premium after hitting my free limit in just 3 days. Now I can binge educational content without interruptions.",
      avatar: "/placeholder.svg?height=80&width=80",
    },
  ];

  // Features data
  const features = [
    {
      icon: FastForward,
      title: "Auto-Skip Technology",
      description: "Our AI automatically detects and skips sponsorship segments without any manual input required.",
      color: "from-purple-600 to-purple-400",
    },
    {
      icon: Zap,
      title: "Lightning Fast",
      description: "Skips happen instantly with zero buffering or lag, keeping your viewing experience smooth.",
      color: "from-blue-500 to-blue-300",
    },
    {
      icon: Clock,
      title: "Time Saved Tracker",
      description: "See exactly how much time you've saved by skipping sponsorships across all your videos.",
      color: "from-indigo-600 to-indigo-400",
    },
  ];

  // How it works steps
  const steps = [
    {
      step: "01",
      title: "Install the Extension",
      description: "Add SponsorSkip to Chrome with just one click from the Chrome Web Store.",
      delay: 0,
    },
    {
      step: "02",
      title: "Watch YouTube Normally",
      description: "Continue watching your favorite videos as you always do.",
      delay: 0.2,
    },
    {
      step: "03",
      title: "Sponsorships Auto-Skip",
      description: "Our AI detects and automatically skips sponsored segments in real-time.",
      delay: 0.4,
    },
  ];

  // Current year for the footer
  const currentYear = new Date().getFullYear();

  // Underline animation progress
  const underlineProgress = tweened(0, {
    duration: 500,
    easing: cubicOut,
    delay: 800
  });

  onMount(() => {
    const handleScroll = () => {
      isScrolled = window.scrollY > 50;
    };

    window.addEventListener('scroll', handleScroll);
    $underlineProgress = 1;

    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  });
</script>

<div class="min-h-screen bg-gradient-to-b from-purple-50 to-blue-50 overflow-hidden">
  <!-- Navbar -->
  <header class="fixed w-full z-50 transition-all duration-300 {isScrolled ? 'bg-white/90 backdrop-blur-md shadow-md py-3' : 'bg-transparent py-5'}">
    <div class="container mx-auto px-4 flex justify-between items-center">
      <div in:fly={{ x: -20, duration: 500, opacity: 0 }} class="flex items-center gap-2">
        <FastForward class="h-8 w-8 text-purple-600" />
        <span class="font-bold text-2xl bg-clip-text text-transparent bg-gradient-to-r from-purple-600 to-blue-500">
          SponsorSkip
        </span>
      </div>
      <nav in:fly={{ y: -10, duration: 500, delay: 200, opacity: 0 }} class="hidden md:flex items-center gap-8">
        <a href="#features" class="font-medium text-gray-700 hover:text-purple-600 transition">
          Features
        </a>
        <a href="#how-it-works" class="font-medium text-gray-700 hover:text-purple-600 transition">
          How It Works
        </a>
        <a href="#pricing" class="font-medium text-gray-700 hover:text-purple-600 transition">
          Pricing
        </a>
      </nav>
      <div in:fly={{ x: 20, duration: 500, delay: 300, opacity: 0 }}>
        <a
          href="#upgrade"
          class="bg-gradient-to-r from-purple-600 to-blue-500 text-white px-5 py-2 rounded-full font-medium hover:shadow-lg hover:scale-105 transition duration-300"
        >
          Upgrade Now
        </a>
      </div>
    </div>
  </header>

  <!-- Hero Section -->
  <section class="pt-32 pb-20 md:pt-40 md:pb-32 relative">
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute -top-40 -right-40 w-96 h-96 bg-purple-300 rounded-full opacity-20 blur-3xl"></div>
      <div class="absolute top-60 -left-20 w-72 h-72 bg-blue-300 rounded-full opacity-20 blur-3xl"></div>
    </div>
    <div class="container mx-auto px-4 relative z-10">
      <div class="flex flex-col md:flex-row items-center gap-12">
        <div in:fly={{ y: 20, duration: 500, opacity: 0 }} class="flex-1 text-center md:text-left">
          <div in:scale={{ duration: 500, delay: 200, start: 0.9 }} class="inline-block mb-4 px-4 py-1 bg-purple-100 text-purple-700 rounded-full font-medium text-sm">
            #1 Sponsorship Skipper for YouTube
          </div>
          <h1 class="text-4xl md:text-6xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-purple-600 via-blue-500 to-purple-600">
            Skip Sponsorships <br class=" md:block" />
            <span class="relative">
              Automatically
              <div class="absolute -bottom-2 left-0 w-full h-3 bg-blue-200 -z-10 rounded-full" style="width: {$underlineProgress * 100}%"></div>
            </span>
          </h1>
          <p class="text-xl text-gray-600 mb-8 max-w-xl mx-auto md:mx-0">
            Never waste time on sponsored segments again. Our smart AI detects and skips sponsorships in YouTube videos instantly.
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center md:justify-start">
            <a
              href="#upgrade"
              class="bg-gradient-to-r from-purple-600 to-blue-500 text-white px-8 py-3 rounded-full font-medium shadow-lg hover:shadow-xl transition duration-300 flex items-center justify-center gap-2 hover:scale-105 active:scale-98"
            >
              Upgrade to Premium <ChevronRight class="h-5 w-5" />
            </a>
            <a
              href="https://chrome.google.com/webstore"
              target="_blank"
              rel="noreferrer"
              class="bg-white text-gray-800 border border-gray-200 px-8 py-3 rounded-full font-medium shadow hover:shadow-md transition duration-300 flex items-center justify-center gap-2 hover:scale-105 active:scale-98"
            >
              Install Free Version
            </a>
          </div>
        </div>
        <div in:fly={{ x: 50, duration: 500, delay: 300, opacity: 0 }} class="flex-1">
          <div class="relative">
            <div class="absolute -inset-4 bg-gradient-to-r from-purple-400 to-blue-500 rounded-xl blur-xl opacity-30 animate-pulse"></div>
            <div class="relative bg-white rounded-xl shadow-2xl overflow-hidden border-8 border-white">
              <div class="bg-gray-800 h-8 flex items-center px-4 gap-2">
                <div class="w-3 h-3 rounded-full bg-red-500"></div>
                <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
                <div class="w-3 h-3 rounded-full bg-green-500"></div>
                <div class="ml-4 bg-gray-700 rounded h-5 w-64"></div>
              </div>
              <div class="relative">
                <!-- record the video of video getting skipped and then add it here  -->
                <img
                  src="/placeholder.svg?height=3000&width=5000"
                  alt="YouTube video with sponsorship being skipped"
                  class="w-full"
                />
                <div class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/70 to-transparent p-4">
                  <div class="flex items-center gap-2 text-white">
                    <FastForward class="h-5 w-5 text-purple-400" />
                    <span class="font-medium">Sponsorship detected! Skipping...</span>
                  </div>
                </div>
                <div in:scale={{ duration: 500, delay: 1000, start: 0 }} class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-purple-600/90 rounded-full p-6">
                  <FastForward class="h-10 w-10 text-white" />
                </div>
              </div>
            </div>
            <div in:fly={{ y: 20, duration: 500, delay: 1500, opacity: 0 }} class="absolute -right-4 -bottom-4 bg-white rounded-lg shadow-lg p-3 flex items-center gap-2">
              <Clock class="h-5 w-5 text-purple-600" />
              <span class="font-medium text-gray-800">Save 2.5 hours monthly</span>
            </div>
          </div>
        </div>
      </div>

      <div in:fly={{ y: 20, duration: 500, delay: 1800, opacity: 0 }} class="mt-20 text-center">
        <p class="text-gray-500 mb-6">Trusted by over 500,000+ users worldwide</p>
        <div class="flex flex-wrap justify-center gap-8 opacity-70">
          {#each ["Google", "Microsoft", "Adobe", "Spotify", "Netflix"] as brand, index}
            <div class="text-xl font-bold text-gray-400">
              {brand}
            </div>
          {/each}
        </div>
      </div>
    </div>
  </section>

  <!-- Features Section -->
  <section id="features" class="py-20 bg-white relative">
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute top-40 -right-20 w-72 h-72 bg-purple-100 rounded-full opacity-50 blur-3xl"></div>
      <div class="absolute -bottom-20 -left-20 w-80 h-80 bg-blue-100 rounded-full opacity-50 blur-3xl"></div>
    </div>
    <div class="container mx-auto px-4 relative z-10">
      <div class="text-center mb-16" in:fly={{ y: 20, duration: 500, opacity: 0 }}>
        <h2 class="text-3xl md:text-5xl font-bold mb-6">
          <span class="bg-clip-text text-transparent bg-gradient-to-r from-purple-600 to-blue-500">
            Smart Features
          </span> You'll Love
        </h2>
        <p class="text-xl text-gray-600 max-w-3xl mx-auto">
          Our extension is packed with intelligent features to make your YouTube experience seamless and sponsor-free.
        </p>
      </div>

      <div class="grid md:grid-cols-3 gap-8">
        {#each features as feature, index}
          <div 
            class="bg-white rounded-xl shadow-xl p-8 hover:shadow-2xl transition duration-300 border border-gray-100 relative overflow-hidden group"
            in:fly={{ y: 20, duration: 500, delay: index * 200, opacity: 0 }}
          >
            <div class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r {feature.color}"></div>
            <div class="mb-6 inline-flex rounded-lg p-3 bg-gradient-to-r {feature.color}">
              <svelte:component this={feature.icon} class="h-10 w-10 text-white" />
            </div>
            <h3 class="text-xl font-bold mb-3">{feature.title}</h3>
            <p class="text-gray-600">{feature.description}</p>
            <div class="absolute -right-20 -bottom-20 w-40 h-40 rounded-full bg-gradient-to-r {feature.color} opacity-10 group-hover:opacity-20 transition-opacity duration-300"></div>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- How It Works -->
  <section id="how-it-works" class="py-20 bg-gray-50 relative">
    <div class="container mx-auto px-4">
      <div class="text-center mb-16" in:fly={{ y: 20, duration: 500, opacity: 0 }}>
        <h2 class="text-3xl md:text-5xl font-bold mb-6">
          How <span class="bg-clip-text text-transparent bg-gradient-to-r from-purple-600 to-blue-500">SponsorSkip</span> Works
        </h2>
        <p class="text-xl text-gray-600 max-w-3xl mx-auto">
          Our technology makes skipping sponsorships effortless and automatic
        </p>
      </div>

      <div class="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
        {#each steps as item, index}
          <div class="relative" in:fly={{ y: 20, duration: 500, delay: item.delay * 1000, opacity: 0 }}>
            <div class="bg-white rounded-xl shadow-lg p-8 relative z-10">
              <div class="text-5xl font-bold text-gray-100 mb-4">{item.step}</div>
              <h3 class="text-xl font-bold mb-3">{item.title}</h3>
              <p class="text-gray-600">{item.description}</p>
            </div>
            {#if index < 2}
              <div 
                class="hidden md:block absolute top-1/2 -right-4 z-0 transform translate-x-full h-2 bg-gradient-to-r from-purple-400 to-blue-400 rounded-full" 
                style="width: {Math.min(100, 75 + index * 20)}%"
              ></div>
            {/if}
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- Pricing Section -->
  <section id="pricing" class="py-20 bg-white relative">
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute -top-40 -left-40 w-80 h-80 bg-purple-100 rounded-full opacity-30 blur-3xl"></div>
      <div class="absolute bottom-20 right-20 w-72 h-72 bg-blue-100 rounded-full opacity-30 blur-3xl"></div>
    </div>
    <div class="container mx-auto px-4 relative z-10">
      <div class="text-center mb-16" in:fly={{ y: 20, duration: 500, opacity: 0 }}>
        <h2 class="text-3xl md:text-5xl font-bold mb-6">
          Simple, Transparent <span class="bg-clip-text text-transparent bg-gradient-to-r from-purple-600 to-blue-500">Pricing</span>
        </h2>
        <p class="text-xl text-gray-600 max-w-3xl mx-auto">
          Choose the plan that works best for your YouTube watching habits
        </p>
      </div>

      <div class="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto">
        <div 
          class="bg-white border border-gray-200 rounded-xl shadow-lg p-8 hover:shadow-xl transition duration-300"
          in:fly={{ x: -20, duration: 500, opacity: 0 }}
        >
          <div class="text-center mb-6">
            <div class="text-lg font-medium text-gray-500 mb-2">Free Plan</div>
            <div class="text-4xl font-bold mb-1">$0</div>
            <div class="text-gray-500">Forever Free</div>
          </div>
          <div class="border-t border-gray-100 my-6"></div>
          <ul class="space-y-4 mb-8">
            {#each ["Auto-skip sponsorships", "Up to 50 videos per month", "Basic statistics", "Standard support"] as feature}
              <li class="flex items-center gap-3">
                <Check class="h-5 w-5 text-green-500 flex-shrink-0" />
                <span>{feature}</span>
              </li>
            {/each}
          </ul>
          <a
            href="https://chrome.google.com/webstore"
            target="_blank"
            class="block text-center bg-gray-100 hover:bg-gray-200 text-gray-800 font-medium py-3 px-6 rounded-lg transition duration-300"
            rel="noreferrer"
          >
            Install Free Version
          </a>
        </div>

        <div 
          id="upgrade"
          class="bg-gradient-to-br from-purple-50 to-blue-50 border border-purple-100 rounded-xl shadow-xl p-8 relative overflow-hidden"
          in:fly={{ x: 20, duration: 500, opacity: 0 }}
        >
          <div class="absolute top-0 right-0">
            <div class="bg-gradient-to-r from-purple-600 to-blue-500 text-white text-xs font-bold px-4 py-1 rounded-bl-lg">
              MOST POPULAR
            </div>
          </div>
          <div class="text-center mb-6">
            <div class="text-lg font-medium text-gray-700 mb-2">Premium Plan</div>
            <div class="text-4xl font-bold mb-1">$4.99</div>
            <div class="text-gray-600">per month</div>
          </div>
          <div class="border-t border-purple-100 my-6"></div>
          <ul class="space-y-4 mb-8">
            {#each [
              "Everything in Free plan",
              "Unlimited video skipping",
              "Advanced time-saving statistics",
              "Custom skip preferences",
              "Priority support",
              "Early access to new features",
            ] as feature}
              <li class="flex items-center gap-3">
                <div class="bg-gradient-to-r from-purple-600 to-blue-500 rounded-full p-1">
                  <Check class="h-3 w-3 text-white flex-shrink-0" />
                </div>
                <span class="font-medium">{feature}</span>
              </li>
            {/each}
          </ul>
          <a
            href="#"
            class="block text-center bg-gradient-to-r from-purple-600 to-blue-500 text-white font-medium py-3 px-6 rounded-lg shadow-lg hover:shadow-xl transition duration-300 hover:scale-103 active:scale-98"
          >
            Upgrade to Premium
          </a>
          <div class="text-center mt-4 text-sm text-gray-500">7-day money-back guarantee</div>
        </div>
      </div>
    </div>
  </section>

  <!-- Testimonials -->
  <section class="py-20 bg-gray-50">
    <div class="container mx-auto px-4">
      <div class="text-center mb-16" in:fly={{ y: 20, duration: 500, opacity: 0 }}>
        <h2 class="text-3xl md:text-5xl font-bold mb-6">
          What Our <span class="bg-clip-text text-transparent bg-gradient-to-r from-purple-600 to-blue-500">Users</span> Say
        </h2>
        <p class="text-xl text-gray-600 max-w-3xl mx-auto">
          Join thousands of happy users who save time every day
        </p>
      </div>

      <div class="grid md:grid-cols-3 gap-8">
        {#each testimonials as testimonial, index}
          <div 
            class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition duration-300"
            in:fly={{ y: 20, duration: 500, delay: index * 100, opacity: 0 }}
          >
            <div class="flex items-center gap-4 mb-4">
              <img
                src={testimonial.avatar}
                alt={testimonial.name}
                class="w-12 h-12 rounded-full object-cover"
              />
              <div>
                <div class="font-bold">{testimonial.name}</div>
                <div class="text-sm text-gray-500">{testimonial.role}</div>
              </div>
            </div>
            <p class="text-gray-600">"{testimonial.comment}"</p>
            <div class="mt-4 flex text-yellow-400">
              {#each Array(5) as _, i}
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-5 w-5"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                </svg>
              {/each}
            </div>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- CTA Section -->
  <section class="py-20 bg-gradient-to-r from-purple-600 to-blue-500 text-white">
    <div class="container mx-auto px-4">
      <div class="text-center max-w-3xl mx-auto" in:fly={{ y: 20, duration: 500, opacity: 0 }}>
        <h2 class="text-3xl md:text-5xl font-bold mb-6">Ready to Skip Sponsorships Forever?</h2>
        <p class="text-xl opacity-90 mb-8">
          Join thousands of users who save hours every month with SponsorSkip Premium.
        </p>
        <div class="inline-block hover:scale-105 active:scale-98">
          <a
            href="#upgrade"
            class="bg-white text-purple-600 px-8 py-4 rounded-full font-bold shadow-lg hover:shadow-xl transition duration-300 inline-flex items-center gap-2"
          >
            Upgrade to Premium <ChevronRight class="h-5 w-5" />
          </a>
        </div>
        <p class="mt-4 opacity-80">7-day money-back guarantee. No questions asked.</p>
      </div>
    </div>
  </section>

  <!-- Footer -->
  <footer class="bg-gray-900 text-gray-300 py-12">
    <div class="container mx-auto px-4">
      <div class="grid md:grid-cols-4 gap-8">
        <div>
          <div class="flex items-center gap-2 mb-4">
            <FastForward class="h-6 w-6 text-purple-400" />
            <span class="font-bold text-xl text-white">SponsorSkip</span>
          </div>
          <p class="text-gray-400 mb-4">
            The #1 Chrome extension for automatically skipping sponsorships in YouTube videos.
          </p>
          <div class="flex gap-4">
            {#each ["Twitter", "Facebook", "Instagram"] as social}
              <a href="#" class="text-gray-400 hover:text-white transition">
                {social}
              </a>
            {/each}
          </div>
        </div>

        <div>
          <h3 class="font-bold text-white mb-4">Product</h3>
          <ul class="space-y-2">
            {#each ["Features", "Pricing", "Premium", "Chrome Store"] as item}
              <li>
                <a href="#" class="text-gray-400 hover:text-white transition">
                  {item}
                </a>
              </li>
            {/each}
          </ul>
        </div>

        <div>
          <h3 class="font-bold text-white mb-4">Support</h3>
          <ul class="space-y-2">
            {#each ["Help Center", "Contact Us", "FAQ", "Troubleshooting", "Status"] as item}
              <li>
                <a href="#" class="text-gray-400 hover:text-white transition">
                  {item}
                </a>
              </li>
            {/each}
          </ul>
        </div>

        <div>
          <h3 class="font-bold text-white mb-4">Legal</h3>
          <ul class="space-y-2">
            {#each ["Terms of Service", "Privacy Policy", "Cookie Policy", "GDPR"] as item}
              <li>
                <a href="#" class="text-gray-400 hover:text-white transition">
                  {item}
                </a>
              </li>
            {/each}
          </ul>
        </div>
      </div>

      <div class="border-t border-gray-800 mt-12 pt-8 text-center text-gray-400 text-sm">
        <p>© {currentYear} SponsorSkip. All rights reserved.</p>
      </div>
    </div>
  </footer>
</div>