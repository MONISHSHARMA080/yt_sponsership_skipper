<script lang="ts">
  import { onMount } from 'svelte';
  import { fade, fly, scale } from 'svelte/transition';
  
  // Svelte 5 runes for reactivity
  let count = $state(0);
  let isMenuOpen = $state(false);
  let activeSection = $state('hero');
  let mousePosition = $state({ x: 0, y: 0 });
  
  // Features list
  const features = $state([
    { 
      title: "Instant Skip", 
      description: "Automatically skip sponsored segments in videos", 
      icon: "⏭️" 
    },
    { 
      title: "Custom Rules", 
      description: "Set your own rules for what content to skip", 
      icon: "⚙️" 
    },
    { 
      title: "Statistics", 
      description: "Track how much time you've saved", 
      icon: "📊" 
    },
    { 
      title: "Cloud Sync", 
      description: "Sync your settings across devices", 
      icon: "☁️" 
    }
  ]);
  
  // Pricing tiers
  const pricingTiers = $state([
    {
      name: "Free",
      price: "$0",
      features: ["Skip up to 10 videos per day", "Basic statistics", "Manual skip mode"],
      recommended: false
    },
    {
      name: "Pro",
      price: "$4.99/mo",
      features: ["Unlimited skips", "Advanced statistics", "Priority updates", "Cloud sync"],
      recommended: true
    },
    {
      name: "Team",
      price: "$9.99/mo",
      features: ["Everything in Pro", "Up to 5 users", "Admin dashboard", "API access"],
      recommended: false
    }
  ]);
  
  // Animation functions
  function handleMouseMove(event: MouseEvent) {
    mousePosition.x = event.clientX;
    mousePosition.y = event.clientY;
  }
  
  function incrementCount() {
    count++;
  }
  
  function toggleMenu() {
    isMenuOpen = !isMenuOpen;
  }
  
  // Intersection observer to track active section
  $effect(() => {
    onMount(() => {
      const sections = document.querySelectorAll('section');
      
      const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
          if (entry.isIntersecting) {
            activeSection = entry.target.id;
          }
        });
      }, { threshold: 0.5 });
      
      sections.forEach(section => {
        observer.observe(section);
      });
      
      return () => {
        sections.forEach(section => {
          observer.unobserve(section);
        });
      };
    });
  });
</script>

<svelte:window on:mousemove={handleMouseMove} />

<div class="cursor" style="left: {mousePosition.x}px; top: {mousePosition.y}px"></div>

<header>
  <div class="logo">
    <div class="logo-shape"></div>
    <h1>SponsorSkip</h1>
  </div>
  
  <button class="menu-toggle" on:click={toggleMenu}>
    <div class="bar"></div>
    <div class="bar"></div>
    <div class="bar"></div>
  </button>
  
  {#if isMenuOpen}
    <nav transition:fly={{ y: -50, duration: 300 }}>
      <ul>
        <li><a href="#hero" class:active={activeSection === 'hero'}>Home</a></li>
        <li><a href="#features" class:active={activeSection === 'features'}>Features</a></li>
        <li><a href="#pricing" class:active={activeSection === 'pricing'}>Pricing</a></li>
        <li><a href="#download" class:active={activeSection === 'download'}>Download</a></li>
      </ul>
    </nav>
  {/if}
</header>

<main>
  <section id="hero" class:active={activeSection === 'hero'}>
    <div class="hero-content">
      <h1>SKIP THE <span class="highlight">BORING</span> STUFF</h1>
      <p>The ultimate Chrome extension to automatically skip sponsorships in YouTube videos</p>
      
      <div class="cta-container">
        <a href="#download" class="cta-button primary">Get Extension</a>
        <a href="#pricing" class="cta-button secondary">Upgrade Now</a>
      </div>
      
      <div class="counter-container">
        <p>Users have skipped <span class="counter">{count.toLocaleString()}</span> hours of sponsored content</p>
        <button class="hidden-button" on:click={incrementCount}>+</button>
      </div>
    </div>
    
    <div class="hero-graphics">
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
      <div class="youtube-mockup">
        <div class="video-player">
          <div class="sponsor-overlay">SPONSOR SEGMENT SKIPPED</div>
        </div>
      </div>
    </div>
  </section>
  
  <section id="features" class:active={activeSection === 'features'}>
    <h2>POWERFUL <span class="highlight">FEATURES</span></h2>
    
    <div class="features-grid">
      {#each features as feature, i}
        <div class="feature-card" 
             in:fly={{ y: 50, delay: i * 100, duration: 400 }}
             on:mouseenter={incrementCount}>
          <div class="feature-icon">{feature.icon}</div>
          <h3>{feature.title}</h3>
          <p>{feature.description}</p>
        </div>
      {/each}
    </div>
    
    <div class="feature-highlight">
      <h3>TIME SAVED = <span class="highlight">LIFE GAINED</span></h3>
      <p>The average YouTube user encounters 5 minutes of sponsored content daily. That's 30+ hours per year!</p>
    </div>
  </section>
  
  <section id="pricing" class:active={activeSection === 'pricing'}>
    <h2>UPGRADE YOUR <span class="highlight">EXPERIENCE</span></h2>
    
    <div class="pricing-grid">
      {#each pricingTiers as tier, i}
        <div class="pricing-card {tier.recommended ? 'recommended' : ''}"
             in:scale={{ start: 0.8, delay: i * 150, duration: 400 }}>
          {#if tier.recommended}
            <div class="recommended-badge">POPULAR</div>
          {/if}
          <h3>{tier.name}</h3>
          <div class="price">{tier.price}</div>
          <ul>
            {#each tier.features as feature}
              <li>{feature}</li>
            {/each}
          </ul>
          <a href="#" class="cta-button {tier.recommended ? 'primary' : 'secondary'}">
            {tier.name === 'Free' ? 'Get Started' : 'Upgrade Now'}
          </a>
        </div>
      {/each}
    </div>
  </section>
  
  <section id="download" class:active={activeSection === 'download'}>
    <h2>GET <span class="highlight">STARTED</span> NOW</h2>
    
    <div class="download-container">
      <div class="browser-icon">
        <div class="chrome-logo">
          <div class="chrome-red"></div>
          <div class="chrome-yellow"></div>
          <div class="chrome-green"></div>
          <div class="chrome-blue"></div>
        </div>
      </div>
      
      <a href="#" class="cta-button primary large">Add to Chrome</a>
      
      <p>Works with Chrome, Brave, Edge, and other Chromium-based browsers</p>
    </div>
  </section>
</main>

<footer>
  <div class="footer-grid">
    <div class="footer-column">
      <h3>SponsorSkip</h3>
      <p>Skip the boring stuff and get back to what matters.</p>
    </div>
    
    <div class="footer-column">
      <h4>Links</h4>
      <ul>
        <li><a href="#">Home</a></li>
        <li><a href="#">Features</a></li>
        <li><a href="#">Pricing</a></li>
        <li><a href="#">Download</a></li>
      </ul>
    </div>
    
    <div class="footer-column">
      <h4>Legal</h4>
      <ul>
        <li><a href="#">Privacy Policy</a></li>
        <li><a href="#">Terms of Service</a></li>
        <li><a href="#">Cookie Policy</a></li>
      </ul>
    </div>
    
    <div class="footer-column">
      <h4>Connect</h4>
      <div class="social-icons">
        <a href="#" class="social-icon">T</a>
        <a href="#" class="social-icon">F</a>
        <a href="#" class="social-icon">I</a>
        <a href="#" class="social-icon">G</a>
      </div>
    </div>
  </div>
  
  <div class="copyright">
    <p>&copy; {new Date().getFullYear()} SponsorSkip. All rights reserved.</p>
  </div>
</footer>

<style>
  /* Global styles */
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: 'Inter', 'Roboto Mono', monospace, sans-serif;
    background-color: #0f0f0f;
    color: #ffffff;
    overflow-x: hidden;
    cursor: none;
  }
  
  :global(*) {
    box-sizing: border-box;
  }
  
  /* Custom cursor */
  .cursor {
    position: fixed;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background-color: #ff3e00;
    mix-blend-mode: difference;
    pointer-events: none;
    z-index: 9999;
    transform: translate(-50%, -50%);
    transition: width 0.2s, height 0.2s;
  }
  
  /* Typography */
  h1, h2, h3, h4 {
    font-family: 'Roboto Mono', monospace;
    text-transform: uppercase;
    letter-spacing: -0.05em;
    font-weight: 900;
  }
  
  h1 {
    font-size: 4rem;
    line-height: 1;
    margin: 0;
  }
  
  h2 {
    font-size: 3rem;
    margin: 0 0 2rem 0;
    position: relative;
    display: inline-block;
  }
  
  h2::after {
    content: '';
    position: absolute;
    bottom: -10px;
    left: 0;
    width: 100%;
    height: 5px;
    background: linear-gradient(90deg, #ff3e00, #ff8700);
  }
  
  .highlight {
    color: #ff3e00;
    position: relative;
    display: inline-block;
  }
  
  /* Layout */
  header {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    padding: 1rem 2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    z-index: 100;
    background-color: rgba(15, 15, 15, 0.8);
    backdrop-filter: blur(10px);
    border-bottom: 2px solid #ff3e00;
  }
  
  .logo {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  .logo-shape {
    width: 40px;
    height: 40px;
    background-color: #ff3e00;
    clip-path: polygon(50% 0%, 100% 50%, 50% 100%, 0% 50%);
    animation: rotateLogo 10s linear infinite;
  }
  
  .menu-toggle {
    background: none;
    border: none;
    cursor: pointer;
    display: flex;
    flex-direction: column;
    gap: 6px;
    z-index: 101;
    padding: 10px;
  }
  
  .bar {
    width: 30px;
    height: 3px;
    background-color: white;
    transition: transform 0.3s, opacity 0.3s;
  }
  
  nav {
    position: absolute;
    top: 100%;
    right: 0;
    background-color: #1a1a1a;
    padding: 1rem;
    border: 2px solid #ff3e00;
    border-top: none;
  }
  
  nav ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  
  nav li {
    margin: 1rem 0;
  }
  
  nav a {
    color: white;
    text-decoration: none;
    font-size: 1.2rem;
    font-weight: bold;
    text-transform: uppercase;
    position: relative;
    padding: 0.5rem 0;
  }
  
  nav a::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    width: 0;
    height: 2px;
    background-color: #ff3e00;
    transition: width 0.3s;
  }
  
  nav a:hover::after,
  nav a.active::after {
    width: 100%;
  }
  
  main {
    padding-top: 80px;
  }
  
  section {
    min-height: 100vh;
    padding: 4rem 2rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    position: relative;
    overflow: hidden;
  }
  
  /* Hero section */
  #hero {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
    align-items: center;
  }
  
  .hero-content {
    z-index: 2;
  }
  
  .hero-graphics {
    position: relative;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  
  .shape {
    position: absolute;
    border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%;
    filter: blur(2px);
    opacity: 0.6;
    z-index: 1;
  }
  
  .shape-1 {
    width: 300px;
    height: 300px;
    background-color: #ff3e00;
    top: 10%;
    right: 20%;
    animation: floatAnimation 15s ease-in-out infinite alternate;
  }
  
  .shape-2 {
    width: 200px;
    height: 200px;
    background-color: #4285f4;
    bottom: 20%;
    left: 10%;
    animation: floatAnimation 12s ease-in-out infinite alternate-reverse;
  }
  
  .shape-3 {
    width: 150px;
    height: 150px;
    background-color: #fbbc05;
    top: 50%;
    right: 5%;
    animation: floatAnimation 10s ease-in-out infinite alternate;
  }
  
  .youtube-mockup {
    width: 80%;
    aspect-ratio: 16/9;
    background-color: #1a1a1a;
    border-radius: 8px;
    overflow: hidden;
    position: relative;
    z-index: 2;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
    border: 3px solid white;
  }
  
  .video-player {
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, #333, #111);
    position: relative;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  
  .sponsor-overlay {
    background-color: rgba(255, 62, 0, 0.8);
    color: white;
    padding: 1rem;
    font-weight: bold;
    border-radius: 4px;
    animation: pulseAnimation 2s infinite;
  }
  
  .cta-container {
    display: flex;
    gap: 1rem;
    margin: 2rem 0;
  }
  
  .cta-button {
    display: inline-block;
    padding: 1rem 2rem;
    border-radius: 4px;
    font-weight: bold;
    text-transform: uppercase;
    text-decoration: none;
    letter-spacing: 1px;
    transition: transform 0.3s, box-shadow 0.3s;
    position: relative;
    overflow: hidden;
  }
  
  .cta-button::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
    transition: left 0.7s;
  }
  
  .cta-button:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 20px rgba(0, 0, 0, 0.3);
  }
  
  .cta-button:hover::before {
    left: 100%;
  }
  
  .cta-button.primary {
    background-color: #ff3e00;
    color: white;
  }
  
  .cta-button.secondary {
    background-color: transparent;
    color: white;
    border: 2px solid #ff3e00;
  }
  
  .cta-button.large {
    padding: 1.5rem 3rem;
    font-size: 1.2rem;
  }
  
  .counter-container {
    margin-top: 2rem;
    font-size: 1.2rem;
  }
  
  .counter {
    color: #ff3e00;
    font-weight: bold;
    font-size: 1.5rem;
  }
  
  .hidden-button {
    opacity: 0;
    position: absolute;
  }
  
  /* Features section */
  #features {
    background-color: #1a1a1a;
    position: relative;
  }
  
  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 2rem;
    width: 100%;
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .feature-card {
    background-color: #2a2a2a;
    padding: 2rem;
    border-radius: 8px;
    transition: transform 0.3s, box-shadow 0.3s;
    position: relative;
    overflow: hidden;
    border: 2px solid transparent;
  }
  
  .feature-card::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 5px;
    background: linear-gradient(90deg, #ff3e00, #ff8700);
  }
  
  .feature-card:hover {
    transform: translateY(-10px);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
    border-color: #ff3e00;
  }
  
  .feature-icon {
    font-size: 3rem;
    margin-bottom: 1rem;
  }
  
  .feature-highlight {
    margin-top: 4rem;
    padding: 2rem;
    background-color: #2a2a2a;
    border-radius: 8px;
    max-width: 800px;
    text-align: center;
    border-left: 5px solid #ff3e00;
  }
  
  /* Pricing section */
  #pricing {
    background-color: #0f0f0f;
  }
  
  .pricing-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
    width: 100%;
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .pricing-card {
    background-color: #1a1a1a;
    padding: 2rem;
    border-radius: 8px;
    text-align: center;
    position: relative;
    transition: transform 0.3s, box-shadow 0.3s;
    border: 2px solid #333;
  }
  
  .pricing-card.recommended {
    border-color: #ff3e00;
    transform: scale(1.05);
    z-index: 2;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  }
  
  .pricing-card:hover {
    transform: translateY(-10px);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  }
  
  .pricing-card.recommended:hover {
    transform: translateY(-10px) scale(1.05);
  }
  
  .recommended-badge {
    position: absolute;
    top: 0;
    right: 2rem;
    background-color: #ff3e00;
    color: white;
    padding: 0.5rem 1rem;
    font-weight: bold;
    font-size: 0.8rem;
    transform: translateY(-50%);
    border-radius: 4px;
  }
  
  .price {
    font-size: 3rem;
    font-weight: bold;
    margin: 1rem 0;
  }
  
  .pricing-card ul {
    list-style: none;
    padding: 0;
    margin: 2rem 0;
    text-align: left;
  }
  
  .pricing-card li {
    margin: 0.5rem 0;
    padding-left: 1.5rem;
    position: relative;
  }
  
  .pricing-card li::before {
    content: '✓';
    position: absolute;
    left: 0;
    color: #ff3e00;
  }
  
  /* Download section */
  #download {
    background-color: #1a1a1a;
    text-align: center;
  }
  
  .download-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2rem;
    max-width: 600px;
  }
  
  .browser-icon {
    width: 120px;
    height: 120px;
    background-color: #fff;
    border-radius: 50%;
    display: flex;
    justify-content: center;
    align-items: center;
    margin-bottom: 2rem;
  }
  
  .chrome-logo {
    width: 80px;
    height: 80px;
    position: relative;
    border-radius: 50%;
    overflow: hidden;
  }
  
  .chrome-red, .chrome-yellow, .chrome-green, .chrome-blue {
    position: absolute;
    width: 80px;
    height: 80px;
  }
  
  .chrome-red {
    background-color: #EA4335;
    clip-path: polygon(50% 50%, 100% 0, 0 0);
    transform-origin: bottom center;
    animation: rotateChrome 10s linear infinite;
  }
  
  .chrome-yellow {
    background-color: #FBBC05;
    clip-path: polygon(50% 50%, 100% 0, 100% 100%);
    transform-origin: bottom left;
    animation: rotateChrome 10s linear infinite;
  }
  
  .chrome-green {
    background-color: #34A853;
    clip-path: polygon(50% 50%, 100% 100%, 0 100%);
    transform-origin: top center;
    animation: rotateChrome 10s linear infinite;
  }
  
  .chrome-blue {
    background-color: #4285F4;
    clip-path: polygon(50% 50%, 0 100%, 0 0);
    transform-origin: bottom right;
    animation: rotateChrome 10s linear infinite;
  }
  
  /* Footer */
  footer {
    background-color: #0a0a0a;
    padding: 4rem 2rem 2rem;
    border-top: 5px solid #ff3e00;
  }
  
  .footer-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .footer-column h3, .footer-column h4 {
    margin-top: 0;
    margin-bottom: 1rem;
  }
  
  .footer-column ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  
  .footer-column li {
    margin: 0.5rem 0;
  }
  
  .footer-column a {
    color: #ccc;
    text-decoration: none;
    transition: color 0.3s;
  }
  
  .footer-column a:hover {
    color: #ff3e00;
  }
  
  .social-icons {
    display: flex;
    gap: 1rem;
  }
  
  .social-icon {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 40px;
    height: 40px;
    background-color: #333;
    color: white;
    border-radius: 50%;
    text-decoration: none;
    transition: background-color 0.3s;
  }
  
  .social-icon:hover {
    background-color: #ff3e00;
  }
  
  .copyright {
    margin-top: 3rem;
    text-align: center;
    color: #666;
    font-size: 0.9rem;
  }
  
  /* Animations */
  @keyframes floatAnimation {
    0% {
      transform: translateY(0) rotate(0);
    }
    100% {
      transform: translateY(-20px) rotate(10deg);
    }
  }
  
  @keyframes pulseAnimation {
    0% {
      transform: scale(1);
    }
    50% {
      transform: scale(1.05);
    }
    100% {
      transform: scale(1);
    }
  }
  
  @keyframes rotateLogo {
    0% {
      transform: rotate(0);
    }
    100% {
      transform: rotate(360deg);
    }
  }
  
  @keyframes rotateChrome {
    0% {
      transform: rotate(0);
    }
    100% {
      transform: rotate(360deg);
    }
  }
  
  /* Responsive styles */
  @media (max-width: 768px) {
    h1 {
      font-size: 2.5rem;
    }
    
    h2 {
      font-size: 2rem;
    }
    
    #hero {
      grid-template-columns: 1fr;
    }
    
    .hero-graphics {
      order: -1;
    }
    
    .youtube-mockup {
      width: 100%;
    }
    
    .pricing-card.recommended {
      transform: scale(1);
    }
    
    .pricing-card.recommended:hover {
      transform: translateY(-10px);
    }
  }
</style>