
<script lang="ts">
  import { onMount } from 'svelte';
  
  // Using the new $state rune from Svelte 5
  let activeTab = $state(0);
  let isMenuOpen = $state(false);
  let scrollY = $state(0);
  let pricing = $state([
    { tier: 'Free', price: '$0', skipsPerMonth: 50, features: ['Basic sponsorship detection', 'Skip button overlay', 'Simple statistics'] },
    { tier: 'Premium', price: '$3.99', skipsPerMonth: 'Unlimited', features: ['Advanced AI detection', 'Auto-skip functionality', 'Detailed analytics', 'Priority support', 'Early access to new features'] }
  ]);
  
  // Animation states
  let heroAnimationComplete = $state(false);
  let featuresInView = $state(false);
  let pricingInView = $state(false);
  
  // Intersection observer for animations
  function createIntersectionObserver(element: HTMLElement, callback: () => void) {
    const observer = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting) {
        callback();
        observer.disconnect();
      }
    }, { threshold: 0.3 });
    
    observer.observe(element);
    
    return {
      destroy() {
        observer.disconnect();
      }
    };
  }
  
  onMount(() => {
    // Start hero animation
    setTimeout(() => {
      heroAnimationComplete = true;
    }, 500);
    
    // Track scroll position for parallax effects
    window.addEventListener('scroll', () => {
      scrollY = window.scrollY;
    });
    
    return () => {
      window.removeEventListener('scroll', () => {});
    };
  });
  
  function toggleMenu() {
    isMenuOpen = !isMenuOpen;
  }
  
  function setFeaturesInView() {
    featuresInView = true;
  }
  
  function setPricingInView() {
    pricingInView = true;
  }
</script>

<svelte:head>
  <title>SkipSpot - Skip YouTube Sponsorships Automatically</title>
  <meta name="description" content="SkipSpot is a Chrome extension that automatically detects and skips sponsorship segments in YouTube videos, saving you time and enhancing your viewing experience.">
</svelte:head>

<div class="app-container">
  <!-- Navigation -->
  <nav class:scrolled={scrollY > 50}>
    <div class="nav-container">
      <div class="logo">
        <div class="logo-shape"></div>
        <span>SkipSpot</span>
      </div>
      
      <div class="nav-links" class:open={isMenuOpen}>
        <a href="#features" on:click={() => isMenuOpen = false}>Features</a>
        <a href="#how-it-works" on:click={() => isMenuOpen = false}>How It Works</a>
        <a href="#pricing" on:click={() => isMenuOpen = false}>Pricing</a>
        <a href="#faq" on:click={() => isMenuOpen = false}>FAQ</a>
      </div>
      
      <div class="cta-button-container">
        <a href="#download" class="cta-button">Install Extension</a>
      </div>
      
      <button class="menu-toggle" on:click={toggleMenu} aria-label="Toggle menu">
        <div class="hamburger" class:open={isMenuOpen}>
          <span></span>
          <span></span>
          <span></span>
        </div>
      </button>
    </div>
  </nav>
  
  <!-- Hero Section -->
  <section class="hero">
    <div class="hero-bg">
      <div class="shape shape-1" style="transform: translateY({scrollY * 0.1}px)"></div>
      <div class="shape shape-2" style="transform: translateY({scrollY * -0.2}px)"></div>
      <div class="shape shape-3" style="transform: rotate({scrollY * 0.02}deg)"></div>
    </div>
    
    <div class="hero-content">
      <h1 class:animate={heroAnimationComplete}>
        <span class="gradient-text">Skip</span> YouTube Sponsorships
        <span class="subtitle">Automatically</span>
      </h1>
      
      <p class:animate={heroAnimationComplete}>
        Save time and enjoy uninterrupted viewing with our intelligent 
        sponsorship detection technology.
      </p>
      
      <div class="hero-cta" class:animate={heroAnimationComplete}>
        <a href="#download" class="cta-button primary">Install Free</a>
        <a href="#pricing" class="cta-button secondary">Go Premium</a>
      </div>
    </div>
    
    <div class="hero-demo" class:animate={heroAnimationComplete}>
      <div class="browser-mockup">
        <div class="browser-header">
          <div class="browser-actions">
            <span></span>
            <span></span>
            <span></span>
          </div>
          <div class="browser-address">youtube.com/watch</div>
        </div>
        <div class="browser-content">
          <div class="video-player">
            <div class="video-overlay">
              <div class="sponsor-alert">
                <div class="sponsor-alert-icon">⏭️</div>
                <div class="sponsor-alert-text">Sponsorship Detected & Skipped</div>
              </div>
            </div>
          </div>
          <div class="video-progress">
            <div class="progress-bar">
              <div class="progress-filled"></div>
              <div class="sponsor-segment"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
  
  <!-- Features Section -->
  <section id="features" class="features" use:createIntersectionObserver={setFeaturesInView}>
    <div class="section-header">
      <h2>Powerful <span class="gradient-text">Features</span></h2>
      <p>Everything you need for a seamless YouTube experience</p>
    </div>
    
    <div class="features-grid">
      <div class="feature-card" class:animate={featuresInView} style="animation-delay: 100ms">
        <div class="feature-icon">
          <div class="icon-bg ai-detection"></div>
        </div>
        <h3>AI-Powered Detection</h3>
        <p>Our advanced algorithm accurately identifies sponsorship segments in videos</p>
      </div>
      
      <div class="feature-card" class:animate={featuresInView} style="animation-delay: 200ms">
        <div class="feature-icon">
          <div class="icon-bg auto-skip"></div>
        </div>
        <h3>Automatic Skipping</h3>
        <p>Never manually skip again - we handle it automatically</p>
      </div>
      
      <div class="feature-card" class:animate={featuresInView} style="animation-delay: 300ms">
        <div class="feature-icon">
          <div class="icon-bg customization"></div>
        </div>
        <h3>Customizable Settings</h3>
        <p>Control exactly what gets skipped and how the extension behaves</p>
      </div>
      
      <div class="feature-card" class:animate={featuresInView} style="animation-delay: 400ms">
        <div class="feature-icon">
          <div class="icon-bg time-saved"></div>
        </div>
        <h3>Time-Saving Stats</h3>
        <p>See how much time you've saved by skipping sponsorships</p>
      </div>
    </div>
  </section>
  
  <!-- How It Works -->
  <section id="how-it-works" class="how-it-works">
    <div class="section-header">
      <h2>How <span class="gradient-text">It Works</span></h2>
      <p>Simple, effective, and non-intrusive</p>
    </div>
    
    <div class="steps">
      <div class="step">
        <div class="step-number">1</div>
        <div class="step-content">
          <h3>Install the Extension</h3>
          <p>Add SkipSpot to Chrome with just one click</p>
        </div>
      </div>
      
      <div class="step-connector"></div>
      
      <div class="step">
        <div class="step-number">2</div>
        <div class="step-content">
          <h3>Watch YouTube Normally</h3>
          <p>Continue using YouTube as you always do</p>
        </div>
      </div>
      
      <div class="step-connector"></div>
      
      <div class="step">
        <div class="step-number">3</div>
        <div class="step-content">
          <h3>Sponsorships Get Skipped</h3>
          <p>Our AI detects and skips sponsorship segments automatically</p>
        </div>
      </div>
    </div>
  </section>
  
  <!-- Pricing Section -->
  <section id="pricing" class="pricing" use:createIntersectionObserver={setPricingInView}>
    <div class="section-header">
      <h2>Simple <span class="gradient-text">Pricing</span></h2>
      <p>Choose the plan that works for you</p>
    </div>
    
    <div class="pricing-cards">
      {#each pricing as plan, i}
        <div class="pricing-card {i === 1 ? 'premium' : ''}" class:animate={pricingInView} style="animation-delay: {i * 200}ms">
          <div class="pricing-header">
            <h3>{plan.tier}</h3>
            <div class="price">{plan.price}</div>
            <div class="billing-cycle">{i === 0 ? 'Forever' : 'per month'}</div>
          </div>
          
          <div class="pricing-features">
            <div class="main-feature">
              <span class="feature-value">{plan.skipsPerMonth}</span>
              <span class="feature-name">skips per month</span>
            </div>
            
            <ul>
              {#each plan.features as feature}
                <li>{feature}</li>
              {/each}
            </ul>
          </div>
          
          <div class="pricing-cta">
            <a href="#download" class="cta-button {i === 1 ? 'primary' : 'secondary'}">
              {i === 0 ? 'Install Free' : 'Go Premium'}
            </a>
          </div>
        </div>
      {/each}
    </div>
  </section>
  
  <!-- FAQ Section -->
  <section id="faq" class="faq">
    <div class="section-header">
      <h2>Frequently Asked <span class="gradient-text">Questions</span></h2>
      <p>Everything you need to know about SkipSpot</p>
    </div>
    
    <div class="faq-container">
      <div class="faq-item">
        <button class="faq-question" on:click={() => activeTab = activeTab === 0 ? -1 : 0}>
          How accurate is the sponsorship detection?
          <span class="faq-icon">{activeTab === 0 ? '−' : '+'}</span>
        </button>
        <div class="faq-answer" class:open={activeTab === 0}>
          Our AI-powered detection is highly accurate and continuously improving. We use machine learning algorithms trained on thousands of videos to identify sponsorship patterns. The premium version offers even higher accuracy with our advanced detection model.
        </div>
      </div>
      
      <div class="faq-item">
        <button class="faq-question" on:click={() => activeTab = activeTab === 1 ? -1 : 1}>
          Will this work on all YouTube videos?
          <span class="faq-icon">{activeTab === 1 ? '−' : '+'}</span>
        </button>
        <div class="faq-answer" class:open={activeTab === 1}>
          SkipSpot works on most YouTube videos that contain standard sponsorship segments. Our detection is optimized for English content but works with many other languages as well. We're constantly expanding our language support.
        </div>
      </div>
      
      <div class="faq-item">
        <button class="faq-question" on:click={() => activeTab = activeTab === 2 ? -1 : 2}>
          What happens when I reach my monthly skip limit?
          <span class="faq-icon">{activeTab === 2 ? '−' : '+'}</span>
        </button>
        <div class="faq-answer" class:open={activeTab === 2}>
          Free users can skip up to 50 sponsorship segments per month. After reaching this limit, the extension will still detect sponsorships but will show a skip button instead of skipping automatically. Upgrade to Premium for unlimited skips.
        </div>
      </div>
      
      <div class="faq-item">
        <button class="faq-question" on:click={() => activeTab = activeTab === 3 ? -1 : 3}>
          Can I customize what gets skipped?
          <span class="faq-icon">{activeTab === 3 ? '−' : '+'}</span>
        </button>
        <div class="faq-answer" class:open={activeTab === 3}>
          Yes! SkipSpot offers customization options to control what types of segments get skipped (sponsorships, intros, outros, etc.) and whether to skip automatically or show a button. Premium users get access to more granular controls.
        </div>
      </div>
    </div>
  </section>
  
  <!-- Download Section -->
  <section id="download" class="download">
    <div class="download-content">
      <h2>Ready to <span class="gradient-text">Enhance</span> Your YouTube Experience?</h2>
      <p>Join thousands of users who save hours every month with SkipSpot</p>
      
      <div class="download-buttons">
        <a href="#" class="download-button">
          <div class="download-icon">
            <svg viewBox="0 0 24 24" width="24" height="24" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="8 17 12 21 16 17"></polyline>
              <line x1="12" y1="12" x2="12" y2="21"></line>
              <path d="M20.88 18.09A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.29"></path>
            </svg>
          </div>
          <div class="download-text">
            <span class="download-small">Available in the</span>
            <span class="download-large">Chrome Web Store</span>
          </div>
        </a>
      </div>
      
      <div class="users-count">
        <div class="users-icon">👥</div>
        <div class="users-text">
          <span class="users-number">10,000+</span> active users
        </div>
      </div>
    </div>
  </section>
  
  <!-- Footer -->
  <footer>
    <div class="footer-content">
      <div class="footer-logo">
        <div class="logo">
          <div class="logo-shape"></div>
          <span>SkipSpot</span>
        </div>
        <p>The smartest way to skip YouTube sponsorships</p>
      </div>
      
      <div class="footer-links">
        <div class="footer-column">
          <h4>Product</h4>
          <ul>
            <li><a href="#features">Features</a></li>
            <li><a href="#pricing">Pricing</a></li>
            <li><a href="#download">Download</a></li>
          </ul>
        </div>
        
        <div class="footer-column">
          <h4>Support</h4>
          <ul>
            <li><a href="#faq">FAQ</a></li>
            <li><a href="#">Contact</a></li>
            <li><a href="#">Privacy Policy</a></li>
          </ul>
        </div>
        
        <div class="footer-column">
          <h4>Connect</h4>
          <ul>
            <li><a href="#">Twitter</a></li>
            <li><a href="#">GitHub</a></li>
            <li><a href="#">Discord</a></li>
          </ul>
        </div>
      </div>
    </div>
    
    <div class="footer-bottom">
      <p>&copy; {new Date().getFullYear()} SkipSpot. All rights reserved.</p>
    </div>
  </footer>
</div>

<style>
  /* Global Styles */
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    background-color: #f8f9fc;
    color: #1a1a2e;
    overflow-x: hidden;
  }
  
  :global(*) {
    box-sizing: border-box;
  }
  
  .app-container {
    width: 100%;
    overflow-x: hidden;
  }
  
  section {
    padding: 6rem 2rem;
    position: relative;
  }
  
  .section-header {
    text-align: center;
    margin-bottom: 4rem;
  }
  
  .section-header h2 {
    font-size: 2.5rem;
    font-weight: 800;
    margin-bottom: 1rem;
  }
  
  .section-header p {
    font-size: 1.2rem;
    color: #4a4a68;
    max-width: 600px;
    margin: 0 auto;
  }
  
  .gradient-text {
    background: linear-gradient(90deg, #ff6b6b, #ff9e7d);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
    display: inline-block;
  }
  
  /* Navigation */
  nav {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    z-index: 1000;
    padding: 1.5rem 2rem;
    transition: all 0.3s ease;
  }
  
  nav.scrolled {
    background-color: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
    padding: 1rem 2rem;
  }
  
  .nav-container {
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  
  .logo {
    display: flex;
    align-items: center;
    font-weight: 800;
    font-size: 1.5rem;
    color: #1a1a2e;
  }
  
  .logo-shape {
    width: 24px;
    height: 24px;
    background: linear-gradient(135deg, #ff6b6b, #ff9e7d);
    border-radius: 6px;
    margin-right: 0.5rem;
    position: relative;
    overflow: hidden;
  }
  
  .logo-shape::before {
    content: '';
    position: absolute;
    width: 12px;
    height: 12px;
    background-color: white;
    top: 6px;
    left: 6px;
    clip-path: polygon(0 0, 100% 50%, 0 100%);
  }
  
  .nav-links {
    display: flex;
    gap: 2rem;
  }
  
  .nav-links a {
    color: #1a1a2e;
    text-decoration: none;
    font-weight: 500;
    transition: color 0.2s ease;
  }
  
  .nav-links a:hover {
    color: #ff6b6b;
  }
  
  .cta-button {
    display: inline-block;
    padding: 0.75rem 1.5rem;
    background: linear-gradient(90deg, #ff6b6b, #ff9e7d);
    color: white;
    font-weight: 600;
    border-radius: 50px;
    text-decoration: none;
    transition: all 0.3s ease;
    border: none;
    cursor: pointer;
  }
  
  .cta-button:hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 20px rgba(255, 107, 107, 0.2);
  }
  
  .cta-button.secondary {
    background: white;
    color: #1a1a2e;
    border: 2px solid #e0e0e0;
  }
  
  .cta-button.secondary:hover {
    border-color: #ff6b6b;
    color: #ff6b6b;
  }
  
  .menu-toggle {
    display: none;
    background: none;
    border: none;
    cursor: pointer;
  }
  
  .hamburger {
    width: 24px;
    height: 18px;
    position: relative;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
  }
  
  .hamburger span {
    display: block;
    width: 100%;
    height: 2px;
    background-color: #1a1a2e;
    transition: all 0.3s ease;
  }
  
  .hamburger.open span:nth-child(1) {
    transform: translateY(8px) rotate(45deg);
  }
  
  .hamburger.open span:nth-child(2) {
    opacity: 0;
  }
  
  .hamburger.open span:nth-child(3) {
    transform: translateY(-8px) rotate(-45deg);
  }
  
  /* Hero Section */
  .hero {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding-top: 80px;
    position: relative;
    overflow: hidden;
  }
  
  .hero-bg {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: -1;
  }
  
  .shape {
    position: absolute;
    border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%;
  }
  
  .shape-1 {
    width: 300px;
    height: 300px;
    background: linear-gradient(135deg, rgba(255, 107, 107, 0.1), rgba(255, 158, 125, 0.1));
    top: 10%;
    left: 5%;
  }
  
  .shape-2 {
    width: 400px;
    height: 400px;
    background: linear-gradient(135deg, rgba(255, 107, 107, 0.05), rgba(255, 158, 125, 0.05));
    bottom: -10%;
    right: 5%;
  }
  
  .shape-3 {
    width: 200px;
    height: 200px;
    background: linear-gradient(135deg, rgba(255, 107, 107, 0.08), rgba(255, 158, 125, 0.08));
    top: 40%;
    right: 15%;
  }
  
  .hero-content {
    text-align: center;
    max-width: 800px;
    margin-bottom: 3rem;
    padding: 0 2rem;
  }
  
  .hero-content h1 {
    font-size: 3.5rem;
    font-weight: 900;
    margin-bottom: 1.5rem;
    line-height: 1.1;
    opacity: 0;
    transform: translateY(20px);
    transition: all 0.6s ease;
  }
  
  .hero-content h1.animate {
    opacity: 1;
    transform: translateY(0);
  }
  
  .hero-content .subtitle {
    display: block;
    font-size: 2.5rem;
    margin-top: 0.5rem;
  }
  
  .hero-content p {
    font-size: 1.25rem;
    color: #4a4a68;
    margin-bottom: 2rem;
    opacity: 0;
    transform: translateY(20px);
    transition: all 0.6s ease;
    transition-delay: 0.2s;
  }
  
  .hero-content p.animate {
    opacity: 1;
    transform: translateY(0);
  }
  
  .hero-cta {
    display: flex;
    gap: 1rem;
    justify-content: center;
    opacity: 0;
    transform: translateY(20px);
    transition: all 0.6s ease;
    transition-delay: 0.4s;
  }
  
  .hero-cta.animate {
    opacity: 1;
    transform: translateY(0);
  }
  
  .hero-demo {
    width: 100%;
    max-width: 800px;
    opacity: 0;
    transform: translateY(40px);
    transition: all 0.8s ease;
    transition-delay: 0.6s;
  }
  
  .hero-demo.animate {
    opacity: 1;
    transform: translateY(0);
  }
  
  .browser-mockup {
    background-color: white;
    border-radius: 8px;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
    overflow: hidden;
  }
  
  .browser-header {
    background-color: #f1f1f1;
    padding: 0.75rem;
    display: flex;
    align-items: center;
  }
  
  .browser-actions {
    display: flex;
    gap: 6px;
    margin-right: 1rem;
  }
  
  .browser-actions span {
    width: 12px;
    height: 12px;
    border-radius: 50%;
  }
  
  .browser-actions span:nth-child(1) {
    background-color: #ff5f57;
  }
  
  .browser-actions span:nth-child(2) {
    background-color: #ffbd2e;
  }
  
  .browser-actions span:nth-child(3) {
    background-color: #28c940;
  }
  
  .browser-address {
    background-color: white;
    padding: 0.25rem 1rem;
    border-radius: 4px;
    font-size: 0.8rem;
    color: #666;
    flex-grow: 1;
  }
  
  .browser-content {
    padding: 0;
  }
  
  .video-player {
    background-color: #000;
    aspect-ratio: 16 / 9;
    position: relative;
  }
  
  .video-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .sponsor-alert {
    background-color: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 1rem 2rem;
    border-radius: 50px;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    animation: fadeInOut 3s infinite;
  }
  
  .sponsor-alert-icon {
    font-size: 1.5rem;
  }
  
  .sponsor-alert-text {
    font-weight: 600;
  }
  
  .video-progress {
    padding: 0.75rem;
    background-color: #f1f1f1;
  }
  
  .progress-bar {
    height: 4px;
    background-color: #e0e0e0;
    border-radius: 2px;
    position: relative;
    overflow: hidden;
  }
  
  .progress-filled {
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    width: 70%;
    background-color: #ff0000;
    animation: progressAnimation 10s linear infinite;
  }
  
  .sponsor-segment {
    position: absolute;
    top: 0;
    left: 40%;
    height: 100%;
    width: 15%;
    background-color: #ffbd2e;
    opacity: 0.8;
  }
  
  /* Features Section */
  .features {
    background-color: white;
  }
  
  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .feature-card {
    background-color: white;
    border-radius: 12px;
    padding: 2rem;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
    transition: all 0.3s ease;
    opacity: 0;
    transform: translateY(30px);
  }
  
  .feature-card.animate {
    opacity: 1;
    transform: translateY(0);
  }
  
  .feature-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
  }
  
  .feature-icon {
    margin-bottom: 1.5rem;
  }
  
  .icon-bg {
    width: 60px;
    height: 60px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
  }
  
  .icon-bg::before {
    content: '';
    position: absolute;
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, #ff6b6b, #ff9e7d);
    opacity: 0.15;
    border-radius: inherit;
  }
  
  .icon-bg::after {
    content: '';
    position: absolute;
    width: 24px;
    height: 24px;
    background-repeat: no-repeat;
    background-position: center;
    background-size: contain;
  }
  
  .ai-detection::after {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23ff6b6b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M12 2a10 10 0 1 0 10 10H12V2Z'/%3E%3Cpath d='M21.17 8H12V2.83c2 .44 3.8 1.5 5.17 3a10.02 10.02 0 0 1 4 5.17'/%3E%3C/svg%3E");
  }
  
  .auto-skip::after {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23ff6b6b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolygon points='5 4 15 12 5 20 5 4'/%3E%3Cline x1='19' y1='5' x2='19' y2='19'/%3E%3C/svg%3E");
  }
  
  .customization::after {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23ff6b6b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z'/%3E%3Ccircle cx='12' cy='12' r='3'/%3E%3C/svg%3E");
  }
  
  .time-saved::after {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23ff6b6b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='10'/%3E%3Cpolyline points='12 6 12 12 16 14'/%3E%3C/svg%3E");
  }
  
  .feature-card h3 {
    font-size: 1.25rem;
    font-weight: 700;
    margin-bottom: 0.75rem;
  }
  
  .feature-card p {
    color: #4a4a68;
    line-height: 1.6;
  }
  
  /* How It Works */
  .how-it-works {
    background-color: #f8f9fc;
  }
  
  .steps {
    max-width: 800px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .step {
    display: flex;
    align-items: flex-start;
    gap: 1.5rem;
    background-color: white;
    padding: 2rem;
    border-radius: 12px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  }
  
  .step-number {
    width: 40px;
    height: 40px;
    background: linear-gradient(135deg, #ff6b6b, #ff9e7d);
    color: white;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
    font-size: 1.25rem;
    flex-shrink: 0;
  }
  
  .step-content h3 {
    font-size: 1.25rem;
    font-weight: 700;
    margin-bottom: 0.5rem;
  }
  
  .step-content p {
    color: #4a4a68;
  }
  
  .step-connector {
    width: 2px;
    height: 40px;
    background-color: #e0e0e0;
    margin-left: 20px;
  }
  
  /* Pricing Section */
  .pricing {
    background-color: white;
  }
  
  .pricing-cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
    max-width: 1000px;
    margin: 0 auto;
  }
  
  .pricing-card {
    background-color: white;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
    transition: all 0.3s ease;
    opacity: 0;
    transform: translateY(30px);
    border: 2px solid #f1f1f1;
  }
  
  .pricing-card.animate {
    opacity: 1;
    transform: translateY(0);
  }
  
  .pricing-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
  }
  
  .pricing-card.premium {
    border-color: #ff6b6b;
    position: relative;
  }
  
  .pricing-card.premium::before {
    content: 'Most Popular';
    position: absolute;
    top: 1rem;
    right: 1rem;
    background: linear-gradient(90deg, #ff6b6b, #ff9e7d);
    color: white;
    padding: 0.25rem 0.75rem;
    border-radius: 50px;
    font-size: 0.75rem;
    font-weight: 600;
  }
  
  .pricing-header {
    padding: 2rem;
    text-align: center;
    border-bottom: 1px solid #f1f1f1;
  }
  
  .pricing-header h3 {
    font-size: 1.5rem;
    font-weight: 700;
    margin-bottom: 1rem;
  }
  
  .price {
    font-size: 3rem;
    font-weight: 800;
    line-height: 1;
    margin-bottom: 0.5rem;
  }
  
  .pricing-card.premium .price {
    color: #ff6b6b;
  }
  
  .billing-cycle {
    color: #4a4a68;
    font-size: 0.9rem;
  }
  
  .pricing-features {
    padding: 2rem;
  }
  
  .main-feature {
    text-align: center;
    margin-bottom: 1.5rem;
    padding-bottom: 1.5rem;
    border-bottom: 1px solid #f1f1f1;
  }
  
  .feature-value {
    font-size: 1.5rem;
    font-weight: 700;
    display: block;
    margin-bottom: 0.25rem;
  }
  
  .pricing-card.premium .feature-value {
    color: #ff6b6b;
  }
  
  .feature-name {
    color: #4a4a68;
    font-size: 0.9rem;
  }
  
  .pricing-features ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  
  .pricing-features li {
    padding: 0.5rem 0;
    display: flex;
    align-items: center;
    color: #4a4a68;
  }
  
  .pricing-features li::before {
    content: '✓';
    margin-right: 0.5rem;
    color: #ff6b6b;
    font-weight: 700;
  }
  
  .pricing-cta {
    padding: 0 2rem 2rem;
    text-align: center;
  }
  
  /* FAQ Section */
  .faq {
    background-color: #f8f9fc;
  }
  
  .faq-container {
    max-width: 800px;
    margin: 0 auto;
  }
  
  .faq-item {
    margin-bottom: 1rem;
    border-radius: 8px;
    overflow: hidden;
    background-color: white;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);
  }
  
  .faq-question {
    width: 100%;
    text-align: left;
    padding: 1.5rem;
    background: none;
    border: none;
    font-size: 1.1rem;
    font-weight: 600;
    cursor: pointer;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .faq-icon {
    font-size: 1.5rem;
    color: #ff6b6b;
    font-weight: 300;
  }
  
  .faq-answer {
    max-height: 0;
    overflow: hidden;
    transition: max-height 0.3s ease, padding 0.3s ease;
    padding: 0 1.5rem;
  }
  
  .faq-answer.open {
    max-height: 300px;
    padding: 0 1.5rem 1.5rem;
  }
  
  /* Download Section */
  .download {
    background: linear-gradient(135deg, #ff6b6b, #ff9e7d);
    color: white;
    text-align: center;
  }
  
  .download-content {
    max-width: 800px;
    margin: 0 auto;
  }
  
  .download h2 {
    font-size: 2.5rem;
    font-weight: 800;
    margin-bottom: 1.5rem;
  }
  
  .download .gradient-text {
    background: white;
    -webkit-background-clip: text;
    background-clip: text;
  }
  
  .download p {
    font-size: 1.2rem;
    margin-bottom: 2.5rem;
    opacity: 0.9;
  }
  
  .download-buttons {
    display: flex;
    justify-content: center;
    margin-bottom: 3rem;
  }
  
  .download-button {
    display: flex;
    align-items: center;
    background-color: white;
    color: #1a1a2e;
    padding: 0.75rem 1.5rem;
    border-radius: 8px;
    text-decoration: none;
    transition: all 0.3s ease;
  }
  
  .download-button:hover {
    transform: translateY(-3px);
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  }
  
  .download-icon {
    margin-right: 1rem;
  }
  
  .download-text {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    text-align: left;
  }
  
  .download-small {
    font-size: 0.8rem;
    opacity: 0.7;
  }
  
  .download-large {
    font-size: 1.1rem;
    font-weight: 600;
  }
  
  .users-count {
    display: inline-flex;
    align-items: center;
    background-color: rgba(255, 255, 255, 0.2);
    padding: 0.75rem 1.5rem;
    border-radius: 50px;
  }
  
  .users-icon {
    margin-right: 0.75rem;
    font-size: 1.25rem;
  }
  
  .users-number {
    font-weight: 700;
  }
  
  /* Footer */
  footer {
    background-color: #1a1a2e;
    color: white;
    padding: 4rem 2rem 2rem;
  }
  
  .footer-content {
    max-width: 1200px;
    margin: 0 auto;
    display: grid;
    grid-template-columns: 1fr;
    gap: 3rem;
  }
  
  .footer-logo p {
    margin-top: 1rem;
    color: rgba(255, 255, 255, 0.7);
  }
  
  .footer-links {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 2rem;
  }
  
  .footer-column h4 {
    font-size: 1.1rem;
    margin-bottom: 1.5rem;
    font-weight: 600;
  }
  
  .footer-column ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  
  .footer-column li {
    margin-bottom: 0.75rem;
  }
  
  .footer-column a {
    color: rgba(255, 255, 255, 0.7);
    text-decoration: none;
    transition: color 0.2s ease;
  }
  
  .footer-column a:hover {
    color: #ff6b6b;
  }
  
  .footer-bottom {
    max-width: 1200px;
    margin: 0 auto;
    padding-top: 2rem;
    margin-top: 2rem;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    text-align: center;
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.9rem;
  }
  
  /* Animations */
  @keyframes fadeInOut {
    0%, 100% { opacity: 0; }
    20%, 80% { opacity: 1; }
  }
  
  @keyframes progressAnimation {
    0% { width: 0%; }
    40% { width: 40%; }
    60% { width: 70%; }
    100% { width: 100%; }
  }
  
  /* Responsive Styles */
  @media (min-width: 768px) {
    .footer-content {
      grid-template-columns: 1fr 2fr;
    }
  }
  
  @media (max-width: 768px) {
    .nav-links {
      position: fixed;
      top: 80px;
      left: 0;
      width: 100%;
      background-color: white;
      flex-direction: column;
      padding: 2rem;
      gap: 1.5rem;
      box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
      transform: translateY(-100%);
      opacity: 0;
      transition: all 0.3s ease;
      z-index: 999;
    }
    
    .nav-links.open {
      transform: translateY(0);
      opacity: 1;
    }
    
    .menu-toggle {
      display: block;
    }
    
    .cta-button-container {
      display: none;
    }
    
    .hero-content h1 {
      font-size: 2.5rem;
    }
    
    .hero-content .subtitle {
      font-size: 2rem;
    }
    
    .section-header h2 {
      font-size: 2rem;
    }
    
    .pricing-cards {
      grid-template-columns: 1fr;
      max-width: 400px;
    }
  }
</style>