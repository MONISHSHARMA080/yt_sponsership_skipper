<script>
  import { onMount } from 'svelte';
  import { fade, fly, scale } from 'svelte/transition';
  import { spring } from 'svelte/motion';
  
  // Runes for state management
  let isMenuOpen = false;
  let activeSection = 'hero';
  let isPremiumModalOpen = false;
  let mousePosition = { x: 0, y: 0 };
  
  // Spring animation for the hero image
  const heroImagePosition = spring({ x: 50, y: 50 }, {
    stiffness: 0.1,
    damping: 0.4
  });
  
  // Features list
  const features = [
    { 
      title: 'Auto-Skip', 
      description: 'Automatically skips sponsored segments in videos',
      icon: '⏭️'
    },
    { 
      title: 'Time Saved', 
      description: 'See how much time you\'ve saved by skipping sponsors',
      icon: '⏱️'
    },
    { 
      title: 'Customizable', 
      description: 'Choose which types of segments to skip',
      icon: '⚙️'
    },
    { 
      title: 'Smart Detection', 
      description: 'Uses AI to identify sponsored content',
      icon: '🧠'
    }
  ];
  
  // Premium features
  const premiumFeatures = [
    'Unlimited videos processed',
    'Early access to new features',
    'Priority support',
    'No waiting time between skips',
    'Custom skip animations'
  ];
  
  // Handle mouse movement for interactive elements
  /**
	 * @param {{ clientX: number; clientY: number; }} event
	 */
  function handleMouseMove(event) {
    mousePosition = { 
      x: event.clientX, 
      y: event.clientY 
    };
    
    // Update spring animation based on mouse position
    heroImagePosition.set({ 
      x: 50 + (event.clientX / window.innerWidth - 0.5) * 10,
      y: 50 + (event.clientY / window.innerHeight - 0.5) * 10
    });
  }
  
  // Intersection observer to detect active section
  onMount(() => {
    const sections = document.querySelectorAll('section[id]');
    
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
</script>

<svelte:window on:mousemove={handleMouseMove} />

<div class="app" class:dark={activeSection === 'premium'}>
  <!-- Navigation -->
  <nav>
    <div class="logo">
      <span class="logo-text">SponsorSkip</span>
      <div class="logo-dot"></div>
    </div>
    
    <div class="nav-links" class:open={isMenuOpen}>
      <a href="#hero" class:active={activeSection === 'hero'}>Home</a>
      <a href="#features" class:active={activeSection === 'features'}>Features</a>
      <a href="#how-it-works" class:active={activeSection === 'how-it-works'}>How It Works</a>
      <a href="#premium" class:active={activeSection === 'premium'}>Premium</a>
    </div>
    
    <!-- svelte-ignore a11y_consider_explicit_label -->
    <button class="menu-toggle" on:click={() => isMenuOpen = !isMenuOpen}>
      <div class="hamburger" class:open={isMenuOpen}></div>
    </button>
  </nav>

  <!-- Hero Section -->
  <section id="hero" class="hero">
    <div class="hero-content">
      <h1 in:fly={{ y: 50, duration: 800 }}>
        <span class="gradient-text">Skip</span> the Sponsors,
        <br />
        <span class="outline-text">Enjoy</span> the Content
      </h1>
      
      <p in:fly={{ y: 50, duration: 800, delay: 200 }}>
        SponsorSkip automatically detects and skips sponsored segments in YouTube videos,
        saving you time and enhancing your viewing experience.
      </p>
      
      <div class="cta-buttons" in:fly={{ y: 50, duration: 800, delay: 400 }}>
        <button class="primary-button">
          <span class="button-text">Add to Chrome</span>
          <span class="button-icon">→</span>
        </button>
        
        <button class="secondary-button" on:click={() => isPremiumModalOpen = true}>
          Upgrade to Premium
        </button>
      </div>
      
      <div class="stats" in:fly={{ y: 50, duration: 800, delay: 600 }}>
        <div class="stat">
          <span class="stat-number">1M+</span>
          <span class="stat-label">Users</span>
        </div>
        
        <div class="stat">
          <span class="stat-number">4.8</span>
          <span class="stat-label">Rating</span>
        </div>
        
        <div class="stat">
          <span class="stat-number">10M+</span>
          <span class="stat-label">Hours Saved</span>
        </div>
      </div>
    </div>
    
    <div class="hero-image-container">
      <div 
        class="hero-image" 
        style="transform: translate({$heroImagePosition.x - 50}px, {$heroImagePosition.y - 50}px)"
      >
        <div class="browser-mockup">
          <div class="browser-controls">
            <div class="browser-dot"></div>
            <div class="browser-dot"></div>
            <div class="browser-dot"></div>
          </div>
          <div class="browser-content">
            <div class="video-player">
              <div class="video-timeline">
                <div class="sponsor-segment"></div>
                <div class="sponsor-segment"></div>
                <div class="sponsor-segment"></div>
              </div>
              <div class="skip-button">SKIP</div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="abstract-shapes">
        <div class="shape shape-1"></div>
        <div class="shape shape-2"></div>
        <div class="shape shape-3"></div>
      </div>
    </div>
  </section>

  <!-- Features Section -->
  <section id="features" class="features">
    <h2>Powerful <span class="gradient-text">Features</span></h2>
    
    <div class="features-grid">
      {#each features as feature, i}
        <div 
          class="feature-card" 
          in:fly={{ y: 50, duration: 800, delay: i * 100 }}
        >
          <div class="feature-icon">{feature.icon}</div>
          <h3>{feature.title}</h3>
          <p>{feature.description}</p>
        </div>
      {/each}
    </div>
  </section>

  <!-- How It Works Section -->
  <section id="how-it-works" class="how-it-works">
    <h2>How It <span class="outline-text">Works</span></h2>
    
    <div class="steps">
      <div class="step" in:fly={{ x: -50, duration: 800 }}>
        <div class="step-number">1</div>
        <h3>Install the Extension</h3>
        <p>Add SponsorSkip to your Chrome browser with just one click</p>
      </div>
      
      <div class="step-connector"></div>
      
      <div class="step" in:fly={{ y: 50, duration: 800, delay: 200 }}>
        <div class="step-number">2</div>
        <h3>Watch YouTube</h3>
        <p>Continue watching your favorite videos as usual</p>
      </div>
      
      <div class="step-connector"></div>
      
      <div class="step" in:fly={{ x: 50, duration: 800, delay: 400 }}>
        <div class="step-number">3</div>
        <h3>Skip Automatically</h3>
        <p>SponsorSkip detects and skips sponsored segments automatically</p>
      </div>
    </div>
    
    <div class="demo-container" in:fade={{ duration: 800, delay: 600 }}>
      <div class="demo-video">
        <div class="video-progress">
          <div class="progress-bar"></div>
          <div class="sponsor-marker" style="left: 30%"></div>
          <div class="sponsor-marker" style="left: 60%"></div>
          <div class="sponsor-marker" style="left: 85%"></div>
        </div>
        <div class="video-controls">
          <div class="control-button">▶</div>
          <div class="time-display">2:45 / 10:30</div>
          <div class="saved-time">Time Saved: 1:15</div>
        </div>
      </div>
    </div>
  </section>

  <!-- Premium Section -->
  <section id="premium" class="premium">
    <div class="premium-content">
      <h2 in:fly={{ y: 50, duration: 800 }}>
        Upgrade to <span class="gradient-text">Premium</span>
      </h2>
      
      <p in:fly={{ y: 50, duration: 800, delay: 200 }}>
        Take your viewing experience to the next level with SponsorSkip Premium
      </p>
      
      <ul class="premium-features" in:fly={{ y: 50, duration: 800, delay: 400 }}>
        {#each premiumFeatures as feature, i}
          <li>
            <span class="check-icon">✓</span>
            {feature}
          </li>
        {/each}
      </ul>
      
      <button 
        class="premium-button" 
        on:click={() => isPremiumModalOpen = true}
        in:fly={{ y: 50, duration: 800, delay: 600 }}
      >
        Get Premium
      </button>
    </div>
    
    <div class="premium-decoration">
      <div class="premium-shape premium-shape-1"></div>
      <div class="premium-shape premium-shape-2"></div>
      <div class="premium-shape premium-shape-3"></div>
    </div>
  </section>

  <!-- Footer -->
  <footer>
    <div class="footer-content">
      <div class="footer-logo">
        <span class="logo-text">SponsorSkip</span>
        <div class="logo-dot"></div>
      </div>
      
      <div class="footer-links">
        <div class="footer-column">
          <h4>Product</h4>
          <a href="#features">Features</a>
          <a href="#premium">Premium</a>
          <a href="#how-it-works">How It Works</a>
        </div>
        
        <div class="footer-column">
          <h4>Company</h4>
          <a href="#">About Us</a>
          <a href="#">Blog</a>
          <a href="#">Careers</a>
        </div>
        
        <div class="footer-column">
          <h4>Support</h4>
          <a href="#">Help Center</a>
          <a href="#">Contact Us</a>
          <a href="#">Privacy Policy</a>
        </div>
      </div>
    </div>
    
    <div class="footer-bottom">
      <p>© {new Date().getFullYear()} SponsorSkip. All rights reserved.</p>
    </div>
  </footer>

  <!-- Premium Modal -->
  {#if isPremiumModalOpen}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="modal-overlay" on:click={() => isPremiumModalOpen = false} in:fade={{ duration: 300 }}>
      <div class="premium-modal" on:click|stopPropagation in:scale={{ start: 0.8, duration: 400 }}>
        <button class="close-modal" on:click={() => isPremiumModalOpen = false}>×</button>
        
        <h2>SponsorSkip Premium</h2>
        
        <div class="pricing-options">
          <div class="pricing-card">
            <h3>Monthly</h3>
            <div class="price">$4.99<span>/month</span></div>
            <ul>
              {#each premiumFeatures as feature}
                <li>✓ {feature}</li>
              {/each}
            </ul>
            <button class="primary-button">Subscribe</button>
          </div>
          
          <div class="pricing-card featured">
            <div class="best-value">Best Value</div>
            <h3>Yearly</h3>
            <div class="price">$39.99<span>/year</span></div>
            <div class="savings">Save 33%</div>
            <ul>
              {#each premiumFeatures as feature}
                <li>✓ {feature}</li>
              {/each}
            </ul>
            <button class="primary-button">Subscribe</button>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  /* Global Styles */
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    overflow-x: hidden;
    background-color: #ffffff;
    color: #1a1a1a;
  }
  
  :global(*) {
    box-sizing: border-box;
  }
  
  .app {
    width: 100%;
    min-height: 100vh;
    transition: background-color 0.5s ease;
  }
  
  .app.dark {
    background-color: #121212;
    color: #ffffff;
  }
  
  section {
    padding: 6rem 2rem;
    position: relative;
    overflow: hidden;
  }
  
  h1, h2, h3, h4 {
    margin: 0;
    line-height: 1.2;
  }
  
  h1 {
    font-size: clamp(2.5rem, 5vw, 4.5rem);
    font-weight: 800;
    margin-bottom: 1.5rem;
  }
  
  h2 {
    font-size: clamp(2rem, 4vw, 3.5rem);
    font-weight: 700;
    margin-bottom: 3rem;
    text-align: center;
  }
  
  p {
    line-height: 1.6;
    margin: 0 0 1.5rem 0;
    font-size: clamp(1rem, 1.5vw, 1.125rem);
  }
  
  a {
    text-decoration: none;
    color: inherit;
    transition: color 0.3s ease;
  }
  
  button {
    cursor: pointer;
    border: none;
    font-family: inherit;
    font-weight: 600;
    transition: all 0.3s ease;
  }
  
  .gradient-text {
    background: linear-gradient(90deg, #FF5E62, #FF9966);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
    display: inline-block;
  }
  
  .outline-text {
    -webkit-text-stroke: 1px currentColor;
    color: transparent;
    display: inline-block;
  }
  
  /* Navigation */
  nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem 2rem;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 100;
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(10px);
    transition: background 0.3s ease;
  }
  
  .dark nav {
    background: rgba(18, 18, 18, 0.8);
  }
  
  .logo {
    display: flex;
    align-items: center;
    font-weight: 800;
    font-size: 1.5rem;
  }
  
  .logo-dot {
    width: 8px;
    height: 8px;
    background: linear-gradient(90deg, #FF5E62, #FF9966);
    border-radius: 50%;
    margin-left: 4px;
  }
  
  .nav-links {
    display: flex;
    gap: 2rem;
  }
  
  .nav-links a {
    position: relative;
    font-weight: 500;
  }
  
  .nav-links a::after {
    content: '';
    position: absolute;
    bottom: -4px;
    left: 0;
    width: 0;
    height: 2px;
    background: linear-gradient(90deg, #FF5E62, #FF9966);
    transition: width 0.3s ease;
  }
  
  .nav-links a:hover::after,
  .nav-links a.active::after {
    width: 100%;
  }
  
  .menu-toggle {
    display: none;
    background: transparent;
    width: 40px;
    height: 40px;
    position: relative;
    z-index: 101;
  }
  
  .hamburger {
    position: relative;
    width: 24px;
    height: 2px;
    background-color: currentColor;
    margin: 0 auto;
    transition: all 0.3s ease;
  }
  
  .hamburger::before,
  .hamburger::after {
    content: '';
    position: absolute;
    width: 24px;
    height: 2px;
    background-color: currentColor;
    transition: all 0.3s ease;
  }
  
  .hamburger::before {
    top: -8px;
  }
  
  .hamburger::after {
    bottom: -8px;
  }
  
  .hamburger.open {
    background-color: transparent;
  }
  
  .hamburger.open::before {
    transform: rotate(45deg);
    top: 0;
  }
  
  .hamburger.open::after {
    transform: rotate(-45deg);
    bottom: 0;
  }
  
  /* Hero Section */
  .hero {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
    min-height: 100vh;
    padding-top: 8rem;
  }
  
  .hero-content {
    display: flex;
    flex-direction: column;
    justify-content: center;
  }
  
  .cta-buttons {
    display: flex;
    gap: 1rem;
    margin-bottom: 2rem;
  }
  
  .primary-button {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    background: linear-gradient(90deg, #FF5E62, #FF9966);
    color: white;
    padding: 0.75rem 1.5rem;
    border-radius: 50px;
    font-size: 1rem;
    font-weight: 600;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
  }
  
  .primary-button:hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 20px rgba(255, 94, 98, 0.2);
  }
  
  .secondary-button {
    background: transparent;
    color: currentColor;
    padding: 0.75rem 1.5rem;
    border-radius: 50px;
    font-size: 1rem;
    border: 2px solid #e0e0e0;
    transition: border-color 0.3s ease;
  }
  
  .secondary-button:hover {
    border-color: #FF5E62;
  }
  
  .stats {
    display: flex;
    gap: 2rem;
  }
  
  .stat {
    display: flex;
    flex-direction: column;
  }
  
  .stat-number {
    font-size: 2rem;
    font-weight: 700;
    background: linear-gradient(90deg, #FF5E62, #FF9966);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
  }
  
  .stat-label {
    font-size: 0.875rem;
    color: #666;
  }
  
  .dark .stat-label {
    color: #aaa;
  }
  
  .hero-image-container {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .hero-image {
    position: relative;
    z-index: 2;
    transition: transform 0.1s ease-out;
  }
  
  .browser-mockup {
    background: white;
    border-radius: 8px;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    width: 100%;
    max-width: 500px;
  }
  
  .browser-controls {
    display: flex;
    gap: 6px;
    padding: 12px;
    background: #f1f1f1;
  }
  
  .browser-dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background: #ff5f5f;
  }
  
  .browser-dot:nth-child(2) {
    background: #ffbe2e;
  }
  
  .browser-dot:nth-child(3) {
    background: #2aca44;
  }
  
  .browser-content {
    padding: 20px;
  }
  
  .video-player {
    background: #000;
    border-radius: 4px;
    padding: 20px;
    position: relative;
    aspect-ratio: 16/9;
  }
  
  .video-timeline {
    position: absolute;
    bottom: 20px;
    left: 20px;
    right: 20px;
    height: 4px;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 2px;
  }
  
  .sponsor-segment {
    position: absolute;
    height: 100%;
    background: #FF5E62;
    border-radius: 2px;
  }
  
  .sponsor-segment:nth-child(1) {
    left: 20%;
    width: 10%;
  }
  
  .sponsor-segment:nth-child(2) {
    left: 45%;
    width: 15%;
  }
  
  .sponsor-segment:nth-child(3) {
    left: 70%;
    width: 8%;
  }
  
  .skip-button {
    position: absolute;
    right: 40px;
    bottom: 40px;
    background: rgba(255, 255, 255, 0.9);
    color: #000;
    padding: 6px 12px;
    border-radius: 4px;
    font-size: 0.875rem;
    font-weight: 600;
    animation: pulse 2s infinite;
  }
  
  @keyframes pulse {
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
  
  .abstract-shapes {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 1;
  }
  
  .shape {
    position: absolute;
    border-radius: 50%;
    opacity: 0.6;
    filter: blur(40px);
  }
  
  .shape-1 {
    width: 300px;
    height: 300px;
    background: #FF5E62;
    top: -100px;
    right: -100px;
    animation: float 8s ease-in-out infinite;
  }
  
  .shape-2 {
    width: 200px;
    height: 200px;
    background: #FF9966;
    bottom: -50px;
    left: 10%;
    animation: float 10s ease-in-out infinite reverse;
  }
  
  .shape-3 {
    width: 150px;
    height: 150px;
    background: #9966FF;
    top: 30%;
    right: 20%;
    animation: float 12s ease-in-out infinite;
  }
  
  @keyframes float {
    0% {
      transform: translate(0, 0);
    }
    50% {
      transform: translate(20px, -20px);
    }
    100% {
      transform: translate(0, 0);
    }
  }
  
  /* Features Section */
  .features {
    background: #f9f9f9;
    padding-top: 8rem;
    padding-bottom: 8rem;
  }
  
  .dark .features {
    background: #1a1a1a;
  }
  
  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .feature-card {
    background: white;
    border-radius: 16px;
    padding: 2rem;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
  }
  
  .dark .feature-card {
    background: #222;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  }
  
  .feature-card:hover {
    transform: translateY(-10px);
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  }
  
  .feature-icon {
    font-size: 2.5rem;
    margin-bottom: 1rem;
  }
  
  .feature-card h3 {
    font-size: 1.25rem;
    margin-bottom: 0.5rem;
  }
  
  .feature-card p {
    color: #666;
    font-size: 0.875rem;
    margin: 0;
  }
  
  .dark .feature-card p {
    color: #aaa;
  }
  
  /* How It Works Section */
  .how-it-works {
    padding-top: 8rem;
    padding-bottom: 8rem;
  }
  
  .steps {
    display: flex;
    justify-content: space-between;
    align-items: center;
    max-width: 1000px;
    margin: 0 auto 4rem;
  }
  
  .step {
    text-align: center;
    max-width: 250px;
  }
  
  .step-number {
    width: 60px;
    height: 60px;
    border-radius: 50%;
    background: linear-gradient(90deg, #FF5E62, #FF9966);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.5rem;
    font-weight: 700;
    margin: 0 auto 1rem;
  }
  
  .step h3 {
    margin-bottom: 0.5rem;
  }
  
  .step p {
    font-size: 0.875rem;
    color: #666;
    margin: 0;
  }
  
  .dark .step p {
    color: #aaa;
  }
  
  .step-connector {
    flex-grow: 1;
    height: 2px;
    background: linear-gradient(90deg, #FF5E62, #FF9966);
    max-width: 100px;
    position: relative;
  }
  
  .step-connector::before,
  .step-connector::after {
    content: '';
    position: absolute;
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: #FF5E62;
    top: 50%;
    transform: translateY(-50%);
  }
  
  .step-connector::before {
    left: 0;
  }
  
  .step-connector::after {
    right: 0;
    background: #FF9966;
  }
  
  .demo-container {
    max-width: 800px;
    margin: 0 auto;
  }
  
  .demo-video {
    background: #000;
    border-radius: 8px;
    padding: 20px;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
  }
  
  .video-progress {
    height: 300px;
    background: #222;
    border-radius: 4px;
    position: relative;
    margin-bottom: 1rem;
  }
  
  .progress-bar {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 70%;
    height: 4px;
    background: linear-gradient(90deg, #FF5E62, #FF9966);
    border-radius: 2px;
  }
  
  .sponsor-marker {
    position: absolute;
    bottom: 0;
    width: 8px;
    height: 8px;
    background: #FF5E62;
    border-radius: 50%;
    transform: translateY(50%);
  }
  
  .video-controls {
    display: flex;
    align-items: center;
    justify-content: space-between;
    color: white;
    padding: 0.5rem 0;
  }
  
  .control-button {
    width: 30px;
    height: 30px;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .saved-time {
    background: linear-gradient(90deg, #FF5E62, #FF9966);
    padding: 0.25rem 0.75rem;
    border-radius: 4px;
    font-size: 0.875rem;
    font-weight: 500;
  }
  
  /* Premium Section */
  .premium {
    background: #121212;
    color: white;
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
    padding-top: 8rem;
    padding-bottom: 8rem;
  }
  
  .premium-content {
    max-width: 500px;
    justify-self: end;
  }
  
  .premium-content h2,
  .premium-content p {
    text-align: left;
  }
  
  .premium-features {
    list-style: none;
    padding: 0;
    margin: 2rem 0;
  }
  
  .premium-features li {
    display: flex;
    align-items: center;
    margin-bottom: 1rem;
    font-size: 1.125rem;
  }
  
  .check-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    background: linear-gradient(90deg, #FF5E62, #FF9966);
    border-radius: 50%;
    margin-right: 1rem;
    font-size: 0.875rem;
  }
  
  .premium-button {
    background: white;
    color: #121212;
    padding: 1rem 2rem;
    border-radius: 50px;
    font-size: 1.125rem;
    font-weight: 600;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
  }
  
  .premium-button:hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 20px rgba(255, 255, 255, 0.1);
  }
  
  .premium-decoration {
    position: relative;
    height: 100%;
  }
  
  .premium-shape {
    position: absolute;
    border-radius: 20px;
  }
  
  .premium-shape-1 {
    width: 300px;
    height: 300px;
    background: #FF5E62;
    top: 10%;
    left: 10%;
    transform: rotate(15deg);
    opacity: 0.7;
  }
  
  .premium-shape-2 {
    width: 200px;
    height: 200px;
    background: #FF9966;
    bottom: 20%;
    right: 20%;
    transform: rotate(-10deg);
    opacity: 0.5;
  }
  
  .premium-shape-3 {
    width: 150px;
    height: 150px;
    background: #9966FF;
    top: 40%;
    right: 40%;
    transform: rotate(45deg);
    opacity: 0.3;
  }
  
  /* Footer */
  footer {
    padding: 4rem 2rem 2rem;
    background: #f9f9f9;
  }
  
  .dark footer {
    background: #1a1a1a;
  }
  
  .footer-content {
    display: flex;
    justify-content: space-between;
    max-width: 1200px;
    margin: 0 auto 3rem;
  }
  
  .footer-links {
    display: flex;
    gap: 4rem;
  }
  
  .footer-column {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .footer-column h4 {
    font-size: 1.125rem;
    margin-bottom: 0.5rem;
  }
  
  .footer-column a {
    font-size: 0.875rem;
    color: #666;
    transition: color 0.3s ease;
  }
  
  .dark .footer-column a {
    color: #aaa;
  }
  
  .footer-column a:hover {
    color: #FF5E62;
  }
  
  .footer-bottom {
    border-top: 1px solid #e0e0e0;
    padding-top: 2rem;
    text-align: center;
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .dark .footer-bottom {
    border-color: #333;
  }
  
  .footer-bottom p {
    font-size: 0.875rem;
    color: #666;
    margin: 0;
  }
  
  .dark .footer-bottom p {
    color: #aaa;
  }
  
  /* Premium Modal */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 1rem;
  }
  
  .premium-modal {
    background: white;
    border-radius: 16px;
    padding: 2rem;
    max-width: 800px;
    width: 100%;
    position: relative;
    color: #1a1a1a;
  }
  
  .close-modal {
    position: absolute;
    top: 1rem;
    right: 1rem;
    background: transparent;
    font-size: 1.5rem;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #666;
    transition: background 0.3s ease;
  }
  
  .close-modal:hover {
    background: #f1f1f1;
  }
  
  .premium-modal h2 {
    text-align: center;
    margin-bottom: 2rem;
  }
  
  .pricing-options {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
  }
  
  .pricing-card {
    background: #f9f9f9;
    border-radius: 16px;
    padding: 2rem;
    text-align: center;
    position: relative;
  }
  
  .pricing-card.featured {
    background: #fff;
    box-shadow: 0 10px 30px rgba(255, 94, 98, 0.2);
    border: 2px solid #FF5E62;
  }
  
  .best-value {
    position: absolute;
    top: -12px;
    left: 50%;
    transform: translateX(-50%);
    background: linear-gradient(90deg, #FF5E62, #FF9966);
    color: white;
    padding: 0.25rem 1rem;
    border-radius: 50px;
    font-size: 0.75rem;
    font-weight: 600;
  }
  
  .pricing-card h3 {
    font-size: 1.5rem;
    margin-bottom: 1rem;
  }
  
  .price {
    font-size: 2.5rem;
    font-weight: 700;
    margin-bottom: 0.5rem;
  }
  
  .price span {
    font-size: 1rem;
    font-weight: 400;
    color: #666;
  }
  
  .savings {
    display: inline-block;
    background: #e6f7e6;
    color: #2aca44;
    padding: 0.25rem 0.75rem;
    border-radius: 50px;
    font-size: 0.875rem;
    font-weight: 600;
    margin-bottom: 1rem;
  }
  
  .pricing-card ul {
    list-style: none;
    padding: 0;
    margin: 1.5rem 0;
    text-align: left;
  }
  
  .pricing-card li {
    margin-bottom: 0.75rem;
    font-size: 0.875rem;
  }
  
  /* Responsive Styles */
  @media (max-width: 992px) {
    .hero {
      grid-template-columns: 1fr;
      text-align: center;
    }
    
    .hero-content {
      order: 1;
    }
    
    .hero-image-container {
      order: 0;
      margin-bottom: 2rem;
    }
    
    .cta-buttons {
      justify-content: center;
    }
    
    .stats {
      justify-content: center;
    }
    
    .premium {
      grid-template-columns: 1fr;
    }
    
    .premium-content {
      justify-self: center;
      text-align: center;
    }
    
    .premium-content h2,
    .premium-content p {
      text-align: center;
    }
    
    .premium-features li {
      justify-content: center;
    }
    
    .premium-decoration {
      display: none;
    }
    
    .steps {
      flex-direction: column;
      gap: 2rem;
    }
    
    .step-connector {
      width: 2px;
      height: 50px;
      max-width: none;
    }
    
    .footer-content {
      flex-direction: column;
      gap: 2rem;
    }
    
    .footer-links {
      flex-direction: column;
      gap: 2rem;
    }
    
    .pricing-options {
      grid-template-columns: 1fr;
    }
  }
  
  @media (max-width: 768px) {
    .nav-links {
      position: fixed;
      top: 0;
      right: -100%;
      bottom: 0;
      width: 70%;
      background: white;
      flex-direction: column;
      padding: 6rem 2rem 2rem;
      transition: right 0.3s ease;
      box-shadow: -10px 0 30px rgba(0, 0, 0, 0.1);
    }
    
    .dark .nav-links {
      background: #121212;
    }
    
    .nav-links.open {
      right: 0;
    }
    
    .menu-toggle {
      display: block;
    }
  }
</style>