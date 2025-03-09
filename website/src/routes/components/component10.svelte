<script>
  import { onMount } from 'svelte';
  
  const features = $state([
    { title: "Skip Sponsors", description: "Automatically skip sponsored segments in videos", icon: "⏭️" },
    { title: "Save Time", description: "Save hours of watching unwanted content", icon: "⏱️" },
    { title: "Community Powered", description: "Leverages crowd-sourced segment data", icon: "👥" },
    { title: "Custom Settings", description: "Configure what types of segments to skip", icon: "⚙️" }
  ]);
  
  let stats = $state({
    timesSaved: 0,
    userCount: 0,
    videosSkipped: 0
  });
  
  let pricingTiers = $state([
    { name: "Free", price: "$0", limit: "10 videos/day", features: ["Basic skipping", "Standard accuracy"] },
    { name: "Pro", price: "$4.99", limit: "Unlimited", features: ["Unlimited skipping", "Priority updates", "Advanced categories", "Custom skip rules"] }
  ]);
  
  let isHovering = $state(false);
  let activeSection = $state("hero");
  let mousePosition = $state({ x: 0, y: 0 });
  
  /**
	 * @param {{ clientX: number; clientY: number; }} event
	 */
  function updateMousePosition(event) {
    mousePosition.x = event.clientX;
    mousePosition.y = event.clientY;
  }
  
  function startAnimation() {
    stats.timesSaved = 0;
    stats.userCount = 0;
    stats.videosSkipped = 0;
    
    const interval = setInterval(() => {
      stats.timesSaved += Math.floor(Math.random() * 10);
      stats.userCount += Math.floor(Math.random() * 5);
      stats.videosSkipped += Math.floor(Math.random() * 20);
      
      if (stats.timesSaved > 1000) clearInterval(interval);
    }, 50);
  }
  
  onMount(() => {
    startAnimation();
    
    // Create intersection observer for sections
    const sections = document.querySelectorAll('section');
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          activeSection = entry.target.id;
        }
      });
    }, { threshold: 0.5 });
    
    sections.forEach(section => observer.observe(section));
    
    return () => observer.disconnect();
  });
  
  // Reactive effects with runes
  $effect(() => {
    // This will run whenever activeSection changes
    console.log(`Active section: ${activeSection}`);
  });
</script>

<svelte:window on:mousemove={updateMousePosition} />

<main class="relative overflow-hidden">
  <!-- Abstract geometric background elements -->
  <div class="fixed inset-0 -z-10 bg-black pointer-events-none">
    <div class="absolute top-0 left-0 w-full h-full overflow-hidden">
      {#each Array(20) as _, i}
        <div 
          class="absolute rounded-full opacity-20" 
          style="
            background: hsl({(i * 20) % 360}deg, 100%, 70%); 
            width: {Math.random() * 300 + 50}px; 
            height: {Math.random() * 300 + 50}px; 
            top: {Math.random() * 100}vh; 
            left: {Math.random() * 100}vw; 
            filter: blur({Math.random() * 50 + 10}px);
            animation: float {Math.random() * 10 + 10}s infinite ease-in-out;
            animation-delay: -{Math.random() * 10}s;
          "
        ></div>
      {/each}
    </div>
    
    <!-- Memphis design elements -->
    <div class="absolute inset-0 opacity-10 pointer-events-none">
      {#each Array(15) as _, i}
        <div 
          class="absolute" 
          style="
            width: {Math.random() * 100 + 20}px; 
            height: {Math.random() * 100 + 20}px; 
            top: {Math.random() * 100}vh; 
            left: {Math.random() * 100}vw; 
            background: hsl({(i * 40) % 360}deg, 100%, 60%);
            clip-path: polygon(50% 0%, 100% 50%, 50% 100%, 0% 50%);
            transform: rotate({Math.random() * 360}deg);
            animation: spin {Math.random() * 20 + 10}s infinite linear;
          "
        ></div>
      {/each}
      
      {#each Array(10) as _, i}
        <div 
          class="absolute" 
          style="
            width: {Math.random() * 150 + 50}px; 
            height: {Math.random() * 10 + 5}px; 
            top: {Math.random() * 100}vh; 
            left: {Math.random() * 100}vw; 
            background: hsl({(i * 60) % 360}deg, 100%, 70%);
            transform: rotate({Math.random() * 360}deg);
          "
        ></div>
      {/each}
    </div>
  </div>

  <!-- Mouse follower -->
  <div 
    class="fixed w-40 h-40 rounded-full pointer-events-none mix-blend-screen opacity-50 z-10"
    style="
      background: radial-gradient(circle, rgba(255,255,255,0.8) 0%, rgba(255,255,255,0) 70%);
      transform: translate({mousePosition.x - 80}px, {mousePosition.y - 80}px);
    "
  ></div>

  <!-- Hero Section -->
  <section id="hero" class="min-h-screen flex flex-col items-center justify-center px-4 py-20 relative">
    <div class="absolute inset-0 grid grid-cols-8 grid-rows-8 opacity-10 pointer-events-none">
      {#each Array(64) as _, i}
        <div class="border border-white/20"></div>
      {/each}
    </div>
    
    <div class="max-w-6xl mx-auto text-center relative z-10">
      <div class="mb-6 inline-block relative">
        <div class="absolute -inset-1 bg-gradient-to-r from-pink-600 to-purple-600 rounded-lg blur opacity-75 group-hover:opacity-100 transition duration-1000 group-hover:duration-200 animate-pulse"></div>
        <span class="relative px-5 py-2.5 bg-black text-white rounded-md font-mono text-sm uppercase tracking-wider">Chrome Extension</span>
      </div>
      
      <h1 class="text-7xl md:text-9xl font-black mb-6 text-white leading-none tracking-tighter">
        <span class="block text-transparent bg-clip-text bg-gradient-to-r from-yellow-400 via-pink-500 to-purple-500 animate-gradient">
          SPONSOR
        </span>
        <span class="block text-transparent bg-clip-text bg-gradient-to-r from-blue-400 via-cyan-500 to-green-500 animate-gradient animation-delay-1000">
          SKIPPER
        </span>
      </h1>
      
      <p class="text-xl md:text-2xl mb-10 max-w-3xl mx-auto text-white/80 font-light">
        Skip the boring parts. Watch only what matters. 
        <span class="font-bold text-white">Save your precious time.</span>
      </p>
      
      <div class="stats flex flex-wrap justify-center gap-8 mb-12">
        {#each Object.entries({ "Hours Saved": stats.timesSaved, "Active Users": stats.userCount, "Sponsors Skipped": stats.videosSkipped }) as [label, value]}
          <div class="stat p-4 bg-white/10 backdrop-blur-md rounded-lg border border-white/20">
            <div class="text-5xl font-black text-white mb-2">{value.toLocaleString()}+</div>
            <div class="text-white/60 uppercase tracking-wider text-sm">{label}</div>
          </div>
        {/each}
      </div>
      
      <div class="flex flex-col sm:flex-row gap-6 justify-center">
        <button 
          class="relative group overflow-hidden rounded-lg text-lg font-bold py-4 px-10 bg-gradient-to-br from-purple-600 to-blue-500 text-white transform transition-all duration-300 hover:scale-105 hover:shadow-lg hover:shadow-purple-500/30"
          on:mouseenter={() => isHovering = true}
          on:mouseleave={() => isHovering = false}
        >
          <span class="relative z-10">Install Extension</span>
          <span class="absolute inset-0 bg-gradient-to-br from-pink-600 to-purple-600 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></span>
          <span class="absolute -inset-px bg-gradient-to-br from-pink-400 to-purple-400 opacity-30 group-hover:opacity-50 blur-sm transition-opacity duration-300"></span>
        </button>
        
        <button class="relative group overflow-hidden rounded-lg text-lg font-bold py-4 px-10 bg-white/10 backdrop-blur-sm border border-white/20 text-white transform transition-all duration-300 hover:scale-105 hover:bg-white/20">
          <span class="relative z-10">Learn More</span>
        </button>
      </div>
    </div>
    
    <div class="absolute bottom-10 left-1/2 transform -translate-x-1/2 animate-bounce">
      <div class="w-8 h-12 rounded-full border-2 border-white/50 flex items-start justify-center p-1">
        <div class="w-1 h-3 bg-white/50 rounded-full animate-pulse"></div>
      </div>
    </div>
  </section>

  <!-- Features Section -->
  <section id="features" class="min-h-screen py-20 px-4 relative">
    <div class="max-w-6xl mx-auto">
      <div class="mb-16 text-center">
        <h2 class="text-5xl md:text-7xl font-black mb-6 text-white">
          <span class="relative inline-block">
            <span class="absolute -inset-1 -skew-y-3 bg-pink-600 opacity-30"></span>
            <span class="relative">FEATURES</span>
          </span>
        </h2>
        <p class="text-xl text-white/70 max-w-2xl mx-auto">Our extension is packed with powerful features to enhance your YouTube experience.</p>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        {#each features as feature, i}
          <div 
            class="feature-card relative overflow-hidden p-8 rounded-2xl border border-white/10 backdrop-blur-sm bg-white/5 transform transition-all duration-500 hover:scale-105 hover:bg-white/10"
            style="animation: fadeIn 0.5s ease-out forwards; animation-delay: {i * 0.1}s; opacity: 0;"
          >
            <div class="absolute -right-4 -top-4 text-6xl opacity-10">{feature.icon}</div>
            <h3 class="text-2xl font-bold mb-4 text-white">{feature.title}</h3>
            <p class="text-white/70">{feature.description}</p>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- How It Works Section -->
  <section id="how-it-works" class="min-h-screen py-20 px-4 relative">
    <div class="max-w-6xl mx-auto">
      <div class="mb-16 text-center">
        <h2 class="text-5xl md:text-7xl font-black mb-6 text-white">
          <span class="relative inline-block">
            <span class="absolute -inset-1 skew-y-3 bg-blue-600 opacity-30"></span>
            <span class="relative">HOW IT WORKS</span>
          </span>
        </h2>
      </div>
      
      <div class="relative">
        <!-- Connection lines -->
        <div class="absolute left-1/2 top-0 bottom-0 w-1 bg-gradient-to-b from-pink-500 to-blue-500 hidden md:block"></div>
        
        {#each ['Detect', 'Analyze', 'Skip', 'Save'] as step, i}
          <div class="step-container relative mb-16 md:mb-32">
            <div class="flex flex-col md:flex-row items-center gap-8">
              <div class={`order-1 ${i % 2 === 0 ? 'md:order-1 md:text-right' : 'md:order-3 md:text-left'} flex-1`}>
                <h3 class="text-3xl font-bold mb-4 text-white">{step}</h3>
                <p class="text-white/70 max-w-md">
                  {#if step === 'Detect'}
                    Our extension automatically detects when you're watching a YouTube video and scans for sponsored segments.
                  {:else if step === 'Analyze'}
                    Using our advanced algorithms and community data, we analyze the video to identify sponsored content.
                  {:else if step === 'Skip'}
                    When a sponsored segment is detected, we automatically skip it, saving you valuable time.
                  {:else}
                    Track how much time you've saved and customize your experience in the extension settings.
                  {/if}
                </p>
              </div>
              
              <div class="order-2 z-10">
                <div class="w-16 h-16 rounded-full bg-gradient-to-br from-pink-500 to-purple-600 flex items-center justify-center text-white text-2xl font-bold">
                  {i + 1}
                </div>
              </div>
              
              <div class={`order-3 ${i % 2 === 0 ? 'md:order-3' : 'md:order-1'} flex-1`}>
                <div class="w-full aspect-video rounded-lg overflow-hidden bg-white/10 border border-white/20 flex items-center justify-center">
                  <div class="text-6xl opacity-50">{['📺', '🔍', '⏭️', '⏱️'][i]}</div>
                </div>
              </div>
            </div>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- Pricing Section -->
  <section id="pricing" class="min-h-screen py-20 px-4 relative">
    <div class="max-w-6xl mx-auto">
      <div class="mb-16 text-center">
        <h2 class="text-5xl md:text-7xl font-black mb-6 text-white">
          <span class="relative inline-block">
            <span class="absolute -inset-1 -skew-y-3 bg-green-600 opacity-30"></span>
            <span class="relative">PRICING</span>
          </span>
        </h2>
        <p class="text-xl text-white/70 max-w-2xl mx-auto">Choose the plan that works best for you.</p>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-4xl mx-auto">
        {#each pricingTiers as tier, i}
          <div 
            class={`pricing-card relative overflow-hidden p-8 rounded-2xl border backdrop-blur-sm transform transition-all duration-500 hover:scale-105 ${i === 0 ? 'border-white/10 bg-white/5' : 'border-purple-500/30 bg-gradient-to-br from-purple-900/40 to-blue-900/40'}`}
          >
            {#if i === 1}
              <div class="absolute top-0 right-0">
                <div class="bg-gradient-to-r from-pink-500 to-purple-500 text-white text-xs font-bold px-4 py-1 transform rotate-45 translate-x-6 translate-y-3">
                  POPULAR
                </div>
              </div>
            {/if}
            
            <h3 class="text-2xl font-bold mb-2 text-white">{tier.name}</h3>
            <div class="mb-6">
              <span class="text-4xl font-black text-white">{tier.price}</span>
              <span class="text-white/60">/month</span>
            </div>
            
            <div class="mb-6 pb-6 border-b border-white/10">
              <div class="text-white font-medium mb-2">
                {tier.limit}
              </div>
            </div>
            
            <ul class="mb-8 space-y-3">
              {#each tier.features as feature}
                <li class="flex items-start">
                  <span class="text-green-400 mr-2">✓</span>
                  <span class="text-white/80">{feature}</span>
                </li>
              {/each}
            </ul>
            
            <button 
              class={`w-full py-3 px-6 rounded-lg font-bold transition-all duration-300 ${i === 0 ? 'bg-white/10 text-white hover:bg-white/20' : 'bg-gradient-to-r from-pink-500 to-purple-600 text-white hover:shadow-lg hover:shadow-purple-500/30 hover:scale-105'}`}
            >
              {i === 0 ? 'Get Started' : 'Upgrade Now'}
            </button>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- Footer -->
  <footer class="py-12 px-4 border-t border-white/10 backdrop-blur-sm bg-black/30">
    <div class="max-w-6xl mx-auto">
      <div class="flex flex-col md:flex-row justify-between items-center mb-8">
        <div class="mb-6 md:mb-0">
          <div class="text-3xl font-black text-white mb-2">SPONSOR<span class="text-purple-500">SKIPPER</span></div>
          <p class="text-white/60">Skip the boring parts. Watch only what matters.</p>
        </div>
        
        <div class="flex gap-6">
          {#each ['Twitter', 'GitHub', 'Discord', 'Email'] as social}
            <a href="#" class="text-white/60 hover:text-white transition-colors">{social}</a>
          {/each}
        </div>
      </div>
      
      <div class="border-t border-white/10 pt-8 flex flex-col md:flex-row justify-between items-center">
        <div class="text-white/40 mb-4 md:mb-0">© 2024 SponsorSkipper. All rights reserved.</div>
        <div class="flex gap-6">
          {#each ['Privacy Policy', 'Terms of Service', 'Cookie Policy'] as link}
            <a href="#" class="text-white/40 hover:text-white/60 transition-colors text-sm">{link}</a>
          {/each}
        </div>
      </div>
    </div>
  </footer>
</main>

<style>
  @keyframes float {
    0%, 100% { transform: translateY(0) rotate(0deg); }
    50% { transform: translateY(-20px) rotate(5deg); }
  }
  
  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }
  
  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
  }
  
  @keyframes gradient {
    0% { background-position: 0% 50%; }
    50% { background-position: 100% 50%; }
    100% { background-position: 0% 50%; }
  }
  
  :global(body) {
    background-color: black;
    color: white;
    font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    overflow-x: hidden;
  }
  
  .animate-gradient {
    background-size: 200% 200%;
    animation: gradient 8s ease infinite;
  }
  
  .animation-delay-1000 {
    animation-delay: 1s;
  }
</style>