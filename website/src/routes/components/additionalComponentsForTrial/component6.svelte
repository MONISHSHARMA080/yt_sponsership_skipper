<script>
  import { onMount } from 'svelte';
  import { fade, fly, scale } from 'svelte/transition';
  import { elasticOut } from 'svelte/easing';

  let isMenuOpen = false;
  let currentSection = 'hero';
  let isLoaded = false;
  let showPricingDetails = false;
  
  // Animation properties
  let heroImageLoaded = false;
  let featuresVisible = false;
  let pricingVisible = false;
  
  // Intersection observer for animations
  onMount(() => {
    setTimeout(() => {
      isLoaded = true;
    }, 500);
    
    const observerOptions = {
      threshold: 0.3
    };
    
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          if (entry.target.id === 'features') {
            featuresVisible = true;
          } else if (entry.target.id === 'pricing') {
            pricingVisible = true;
          }
          currentSection = entry.target.id;
        }
      });
    }, observerOptions);
    
    const sections = document.querySelectorAll('section');
    sections.forEach(section => {
      observer.observe(section);
    });
    
    return () => {
      sections.forEach(section => {
        observer.unobserve(section);
      });
    };
  });
  
  function toggleMenu() {
    isMenuOpen = !isMenuOpen;
  }
  
  /**
	 * @param {string} id
	 */
  function scrollToSection(id) {
    const element = document.getElementById(id);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
    }
    if (isMenuOpen) {
      isMenuOpen = false;
    }
  }
  
  function togglePricingDetails() {
    showPricingDetails = !showPricingDetails;
  }
  
  // Geometric shapes for background
  const shapes = [
    { type: 'circle', size: 80, x: 10, y: 15, color: 'var(--accent-color)', delay: 0 },
    { type: 'square', size: 60, x: 85, y: 25, color: 'var(--secondary-color)', delay: 0.2 },
    { type: 'triangle', size: 70, x: 20, y: 70, color: 'var(--primary-color)', delay: 0.4 },
    { type: 'circle', size: 40, x: 75, y: 80, color: 'var(--accent-color-2)', delay: 0.6 },
    { type: 'zigzag', size: 100, x: 50, y: 50, color: 'var(--secondary-color)', delay: 0.8 }
  ];
</script>

<svelte:head>
  <title>SponsorSkip - Skip YouTube Sponsorships Automatically</title>
  <meta name="description" content="SponsorSkip is a Chrome extension that automatically skips sponsorships in YouTube videos, saving you time and enhancing your viewing experience.">
</svelte:head>

<div class="app-container" class:loaded={isLoaded}>
  <!-- Background shapes -->
  <div class="background-shapes">
    {#each shapes as shape, i}
      <div 
        class="shape {shape.type}" 
        style="--x: {shape.x}%; --y: {shape.y}%; --size: {shape.size}px; --color: {shape.color}; --delay: {shape.delay}s;"
        in:fade={{ delay: 300 + (shape.delay * 1000), duration: 1000 }}
      ></div>
    {/each}
  </div>
  
  <!-- Loading animation -->
  {#if !isLoaded}
    <div class="loading-screen">
      <div class="loading-logo">
        <div class="logo-shape"></div>
        <div class="logo-text">SponsorSkip</div>
      </div>
    </div>
  {/if}
  
  <!-- Header -->
  <header class:scrolled={currentSection !== 'hero'}>
    <div class="logo" on:click={() => scrollToSection('hero')}>
      <div class="logo-icon"></div>
      <span>SponsorSkip</span>
    </div>
    
    <nav class:open={isMenuOpen}>
      <ul>
        <li class:active={currentSection === 'hero'}>
          <button on:click={() => scrollToSection('hero')}>Home</button>
        </li>
        <li class:active={currentSection === 'features'}>
          <button on:click={() => scrollToSection('features')}>Features</button>
        </li>
        <li class:active={currentSection === 'how-it-works'}>
          <button on:click={() => scrollToSection('how-it-works')}>How It Works</button>
        </li>
        <li class:active={currentSection === 'pricing'}>
          <button on:click={() => scrollToSection('pricing')}>Pricing</button>
        </li>
      </ul>
    </nav>
    
    <div class="cta-button-container">
      <button class="cta-button" on:click={() => scrollToSection('pricing')}>
        Upgrade Now
      </button>
    </div>
    
    <button class="menu-toggle" on:click={toggleMenu} aria-label="Toggle menu">
      <span></span>
      <span></span>
      <span></span>
    </button>
  </header>
  
  <main>
    <!-- Hero Section -->
    <section id="hero" class="hero-section">
      <div class="hero-content">
        <div class="hero-text" in:fly={{ y: 50, duration: 800, delay: 300 }}>
          <h1>Skip the Sponsors.<br>Watch What Matters.</h1>
          <p>SponsorSkip automatically detects and skips sponsorship segments in YouTube videos, saving you time and enhancing your viewing experience.</p>
          <div class="hero-buttons">
            <button class="primary-button" on:click={() => scrollToSection('pricing')}>
              Upgrade to Premium
            </button>
            <button class="secondary-button" on:click={() => window.open('https://chrome.google.com/webstore', '_blank')}>
              Install Free Version
            </button>
          </div>
        </div>
        
        <div class="hero-image-container" in:fly={{ x: 50, duration: 800, delay: 600 }}>
          <div class="hero-image">
            <div class="browser-mockup">
              <div class="browser-header">
                <div class="browser-controls">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
                <div class="browser-address-bar"></div>
              </div>
              <div class="browser-content">
                <div class="video-player">
                  <div class="video-content"></div>
                  <div class="sponsor-overlay">
                    <div class="sponsor-badge">SPONSOR</div>
                    <div class="skip-button">SKIP</div>
                  </div>
                  <div class="video-controls">
                    <div class="progress-bar">
                      <div class="progress-filled"></div>
                      <div class="sponsor-segment"></div>
                    </div>
                    <div class="control-buttons">
                      <div class="play-button"></div>
                      <div class="volume-button"></div>
                      <div class="time">3:24 / 10:15</div>
                      <div class="settings-button"></div>
                      <div class="fullscreen-button"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="scroll-indicator" on:click={() => scrollToSection('features')}>
        <div class="mouse">
          <div class="wheel"></div>
        </div>
        <div class="arrow"></div>
      </div>
    </section>
    
    <!-- Features Section -->
    <section id="features" class="features-section">
      <div class="section-header" in:fly={{ y: 30, duration: 800 }} class:visible={featuresVisible}>
        <h2>Powerful Features</h2>
        <p>SponsorSkip comes packed with features to enhance your YouTube experience</p>
      </div>
      
      <div class="features-grid" class:visible={featuresVisible}>
        {#each [
          { 
            title: 'Automatic Detection', 
            description: 'Advanced AI algorithms detect sponsorship segments with high accuracy', 
            icon: 'detection',
            delay: 0 
          },
          { 
            title: 'Seamless Skipping', 
            description: 'Skip sponsorships automatically without lifting a finger', 
            icon: 'skip',
            delay: 0.2 
          },
          { 
            title: 'Customizable Settings', 
            description: 'Configure which types of segments to skip and how to skip them', 
            icon: 'settings',
            delay: 0.4 
          },
          { 
            title: 'Time Saved Counter', 
            description: 'See how much time you\'ve saved by skipping sponsorships', 
            icon: 'time',
            delay: 0.6 
          }
        ] as feature, i}
          <div 
            class="feature-card" 
            in:fly={{ y: 30, duration: 800, delay: 300 + (feature.delay * 1000) }}
          >
            <div class="feature-icon {feature.icon}"></div>
            <h3>{feature.title}</h3>
            <p>{feature.description}</p>
          </div>
        {/each}
      </div>
    </section>
    
    <!-- How It Works Section -->
    <section id="how-it-works" class="how-it-works-section">
      <div class="section-header">
        <h2>How It Works</h2>
        <p>SponsorSkip uses advanced technology to make your viewing experience better</p>
      </div>
      
      <div class="steps-container">
        {#each [
          { 
            number: '01', 
            title: 'Install the Extension', 
            description: 'Add SponsorSkip to your Chrome browser with just one click' 
          },
          { 
            number: '02', 
            title: 'Watch YouTube', 
            description: 'Continue watching YouTube videos as you normally would' 
          },
          { 
            number: '03', 
            title: 'Automatic Skipping', 
            description: 'SponsorSkip detects and skips sponsorship segments automatically' 
          },
          { 
            number: '04', 
            title: 'Save Time', 
            description: 'Enjoy more content and less interruptions, saving valuable time' 
          }
        ] as step, i}
          <div class="step-card">
            <div class="step-number">{step.number}</div>
            <div class="step-content">
              <h3>{step.title}</h3>
              <p>{step.description}</p>
            </div>
          </div>
        {/each}
      </div>
      
      <div class="demo-container">
        <div class="demo-video">
          <div class="video-timeline">
            <div class="timeline-bar"></div>
            <div class="sponsor-segment-marker" style="--left: 30%; --width: 15%;"></div>
            <div class="sponsor-segment-marker" style="--left: 60%; --width: 10%;"></div>
            <div class="playhead"></div>
          </div>
          <div class="skip-animation">
            <div class="skip-notification">
              <span>Sponsorship Detected</span>
              <span class="skip-text">Skipping...</span>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- Pricing Section -->
    <section id="pricing" class="pricing-section">
      <div class="section-header" in:fly={{ y: 30, duration: 800 }} class:visible={pricingVisible}>
        <h2>Choose Your Plan</h2>
        <p>Upgrade to Premium for unlimited skipping and advanced features</p>
      </div>
      
      <div class="pricing-cards" class:visible={pricingVisible}>
        <div class="pricing-card free" in:fly={{ y: 30, x: -20, duration: 800, delay: 300 }}>
          <div class="pricing-header">
            <h3>Free</h3>
            <div class="price">$0</div>
            <p>Perfect for casual viewers</p>
          </div>
          <ul class="pricing-features">
            <li>Skip sponsorships on up to 50 videos per month</li>
            <li>Basic detection algorithm</li>
            <li>Standard skip animations</li>
            <li>Community-based detection</li>
          </ul>
          <button class="pricing-button secondary" on:click={() => window.open('https://chrome.google.com/webstore', '_blank')}>
            Install Now
          </button>
        </div>
        
        <div class="pricing-card premium" in:fly={{ y: 30, duration: 800, delay: 500 }}>
          <div class="popular-badge">Most Popular</div>
          <div class="pricing-header">
            <h3>Premium</h3>
            <div class="price">$4.99<span>/month</span></div>
            <p>For dedicated YouTube viewers</p>
          </div>
          <ul class="pricing-features">
            <li>Unlimited video skipping</li>
            <li>Advanced AI detection (99% accuracy)</li>
            <li>Custom skip animations</li>
            <li>Time saved statistics</li>
            <li>Priority support</li>
          </ul>
          <button class="pricing-button primary">
            Upgrade Now
          </button>
        </div>
        
        <div class="pricing-card lifetime" in:fly={{ y: 30, x: 20, duration: 800, delay: 700 }}>
          <div class="pricing-header">
            <h3>Lifetime</h3>
            <div class="price">$49.99<span>/once</span></div>
            <p>Best value for long-term users</p>
          </div>
          <ul class="pricing-features">
            <li>All Premium features</li>
            <li>One-time payment</li>
            <li>Early access to new features</li>
            <li>Exclusive themes and customizations</li>
            <li>Lifetime updates</li>
          </ul>
          <button class="pricing-button secondary">
            Get Lifetime
          </button>
        </div>
      </div>
      
      <div class="comparison-link" in:fade={{ delay: 900 }}>
        <button on:click={togglePricingDetails}>
          {showPricingDetails ? 'Hide detailed comparison' : 'Show detailed comparison'}
        </button>
      </div>
      
      {#if showPricingDetails}
        <div class="pricing-details" in:fly={{ y: 20, duration: 400 }}>
          <table>
            <thead>
              <tr>
                <th>Feature</th>
                <th>Free</th>
                <th>Premium</th>
                <th>Lifetime</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>Monthly video limit</td>
                <td>50 videos</td>
                <td>Unlimited</td>
                <td>Unlimited</td>
              </tr>
              <tr>
                <td>Detection accuracy</td>
                <td>85%</td>
                <td>99%</td>
                <td>99%</td>
              </tr>
              <tr>
                <td>Skip animations</td>
                <td>Standard</td>
                <td>Custom</td>
                <td>Custom</td>
              </tr>
              <tr>
                <td>Statistics dashboard</td>
                <td>Basic</td>
                <td>Advanced</td>
                <td>Advanced</td>
              </tr>
              <tr>
                <td>Support</td>
                <td>Community</td>
                <td>Priority</td>
                <td>Priority</td>
              </tr>
              <tr>
                <td>Early access</td>
                <td>No</td>
                <td>No</td>
                <td>Yes</td>
              </tr>
              <tr>
                <td>Custom themes</td>
                <td>No</td>
                <td>Limited</td>
                <td>Full access</td>
              </tr>
            </tbody>
          </table>
        </div>
      {/if}
    </section>
    
    <!-- Testimonials Section -->
    <section id="testimonials" class="testimonials-section">
      <div class="section-header">
        <h2>What Our Users Say</h2>
        <p>Join thousands of satisfied users who save time every day</p>
      </div>
      
      <div class="testimonials-container">
        {#each [
          {
            name: 'Alex Johnson',
            role: 'Daily YouTube Viewer',
            quote: 'SponsorSkip has saved me hours of time. I no longer have to manually skip through sponsorships!',
            avatar: '1'
          },
          {
            name: 'Sarah Miller',
            role: 'Content Creator',
            quote: 'As someone who watches a lot of tutorials, this extension is a game-changer. The premium version is worth every penny.',
            avatar: '2'
          },
          {
            name: 'Michael Chen',
            role: 'Tech Enthusiast',
            quote: 'The accuracy of sponsor detection is impressive. It rarely misses and the time saved counter is motivating.',
            avatar: '3'
          }
        ] as testimonial, i}
          <div class="testimonial-card">
            <div class="testimonial-avatar" style="--avatar-img: var(--avatar-{testimonial.avatar})"></div>
            <div class="testimonial-content">
              <p>"{testimonial.quote}"</p>
              <div class="testimonial-author">
                <span class="name">{testimonial.name}</span>
                <span class="role">{testimonial.role}</span>
              </div>
            </div>
          </div>
        {/each}
      </div>
    </section>
    
    <!-- CTA Section -->
    <section class="cta-section">
      <div class="cta-content">
        <h2>Ready to Skip the Sponsors?</h2>
        <p>Join thousands of users who save time and enjoy uninterrupted YouTube viewing</p>
        <div class="cta-buttons">
          <button class="primary-button" on:click={() => scrollToSection('pricing')}>
            Upgrade to Premium
          </button>
          <button class="secondary-button" on:click={() => window.open('https://chrome.google.com/webstore', '_blank')}>
            Try Free Version
          </button>
        </div>
      </div>
    </section>
  </main>
  
  <footer>
    <div class="footer-content">
      <div class="footer-logo">
        <div class="logo-icon"></div>
        <span>SponsorSkip</span>
      </div>
      
      <div class="footer-links">
        <div class="footer-column">
          <h4>Product</h4>
          <ul>
            <li><a href="#features">Features</a></li>
            <li><a href="#pricing">Pricing</a></li>
            <li><a href="#how-it-works">How It Works</a></li>
            <li><a href="#testimonials">Testimonials</a></li>
          </ul>
        </div>
        
        <div class="footer-column">
          <h4>Resources</h4>
          <ul>
            <li><a href="#">Documentation</a></li>
            <li><a href="#">FAQ</a></li>
            <li><a href="#">Blog</a></li>
            <li><a href="#">Support</a></li>
          </ul>
        </div>
        
        <div class="footer-column">
          <h4>Company</h4>
          <ul>
            <li><a href="#">About Us</a></li>
            <li><a href="#">Careers</a></li>
            <li><a href="#">Contact</a></li>
            <li><a href="#">Privacy Policy</a></li>
          </ul>
        </div>
      </div>
      
      <div class="footer-newsletter">
        <h4>Stay Updated</h4>
        <p>Subscribe to our newsletter for the latest updates and features</p>
        <div class="newsletter-form">
          <input type="email" placeholder="Your email address" />
          <button>Subscribe</button>
        </div>
      </div>
    </div>
    
    <div class="footer-bottom">
      <p>&copy; {new Date().getFullYear()} SponsorSkip. All rights reserved.</p>
      <div class="social-links">
        <a href="#" aria-label="Twitter"><div class="social-icon twitter">twitter</div></a>
        <a href="#" aria-label="Facebook"><div class="social-icon facebook">facebook</div></a>
        <a href="#" aria-label="Instagram"><div class="social-icon instagram">instagram</div></a>
        <a href="#" aria-label="GitHub"><div class="social-icon github">github</div></a>
      </div>
    </div>
  </footer>
</div>

<style>
  :root {
    --primary-color: #ff3e00;
    --primary-color-dark: #d93600;
    --secondary-color: #40b3ff;
    --accent-color: #ffca28;
    --accent-color-2: #9c27b0;
    --text-color: #2c3e50;
    --text-light: #546e7a;
    --background-color: #ffffff;
    --background-alt: #f9f9f9;
    --card-background: #ffffff;
    --border-color: #e0e0e0;
    
    /* Avatar placeholders */
    --avatar-1: #ffca28;
    --avatar-2: #40b3ff;
    --avatar-3: #9c27b0;
    
    /* Fonts */
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  }
  
  * {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
  }
  
  html, body {
    width: 100%;
    height: 100%;
    scroll-behavior: smooth;
    overflow-x: hidden;
  }
  
  .app-container {
    position: relative;
    width: 100%;
    min-height: 100vh;
    background-color: var(--background-color);
    color: var(--text-color);
    opacity: 0;
    transition: opacity 0.5s ease;
  }
  
  .app-container.loaded {
    opacity: 1;
  }
  
  /* Loading Screen */
  .loading-screen {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: var(--background-color);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }
  
  .loading-logo {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }
  
  .logo-shape {
    width: 60px;
    height: 60px;
    background-color: var(--primary-color);
    border-radius: 12px;
    position: relative;
    animation: pulse 1.5s infinite;
  }
  
  .logo-shape::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 30px;
    height: 30px;
    background-color: white;
    clip-path: polygon(0 0, 100% 50%, 0 100%);
  }
  
  .logo-text {
    font-size: 1.5rem;
    font-weight: 700;
    letter-spacing: 1px;
  }
  
  @keyframes pulse {
    0% {
      transform: scale(1);
    }
    50% {
      transform: scale(1.1);
    }
    100% {
      transform: scale(1);
    }
  }
  
  /* Background Shapes */
  .background-shapes {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
    z-index: -1;
    overflow: hidden;
  }
  
  .shape {
    position: absolute;
    opacity: 0.05;
    transform: translate(-50%, -50%);
    top: var(--y);
    left: var(--x);
    animation: float 15s infinite ease-in-out;
    animation-delay: var(--delay);
  }
  
  .shape.circle {
    width: var(--size);
    height: var(--size);
    border-radius: 50%;
    background-color: var(--color);
  }
  
  .shape.square {
    width: var(--size);
    height: var(--size);
    background-color: var(--color);
    transform: translate(-50%, -50%) rotate(45deg);
  }
  
  .shape.triangle {
    width: var(--size);
    height: var(--size);
    background-color: transparent;
    clip-path: polygon(50% 0%, 0% 100%, 100% 100%);
    background-color: var(--color);
  }
  
  .shape.zigzag {
    width: var(--size);
    height: calc(var(--size) / 5);
    background-color: transparent;
    background: linear-gradient(135deg, var(--color) 25%, transparent 25%) -10px 0,
                linear-gradient(225deg, var(--color) 25%, transparent 25%) -10px 0,
                linear-gradient(315deg, var(--color) 25%, transparent 25%),
                linear-gradient(45deg, var(--color) 25%, transparent 25%);
    background-size: 20px 20px;
  }
  
  @keyframes float {
    0% {
      transform: translate(-50%, -50%) rotate(0deg);
    }
    50% {
      transform: translate(-50%, -50%) translate(20px, 20px) rotate(5deg);
    }
    100% {
      transform: translate(-50%, -50%) rotate(0deg);
    }
  }
  
  /* Header */
  header {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 80px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 2rem;
    z-index: 100;
    transition: all 0.3s ease;
    background-color: transparent;
  }
  
  header.scrolled {
    background-color: rgba(255, 255, 255, 0.95);
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    height: 70px;
  }
  
  .logo {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
  }
  
  .logo-icon {
    width: 32px;
    height: 32px;
    background-color: var(--primary-color);
    border-radius: 8px;
    position: relative;
  }
  
  .logo-icon::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 16px;
    height: 16px;
    background-color: white;
    clip-path: polygon(0 0, 100% 50%, 0 100%);
  }
  
  .logo span {
    font-weight: 700;
    font-size: 1.25rem;
    letter-spacing: 0.5px;
  }
  
  nav {
    position: absolute;
    top: 80px;
    left: 0;
    width: 100%;
    background-color: white;
    transform: translateY(-100%);
    opacity: 0;
    visibility: hidden;
    transition: all 0.3s ease;
    box-shadow: 0 10px 15px rgba(0, 0, 0, 0.1);
  }
  
  nav.open {
    transform: translateY(0);
    opacity: 1;
    visibility: visible;
  }
  
  nav ul {
    list-style: none;
    padding: 1rem;
  }
  
  nav li {
    margin: 0.5rem 0;
  }
  
  nav button {
    background: none;
    border: none;
    font-size: 1rem;
    color: var(--text-color);
    cursor: pointer;
    padding: 0.5rem;
    width: 100%;
    text-align: left;
    transition: color 0.2s ease;
  }
  
  nav button:hover {
    color: var(--primary-color);
  }
  
  li.active button {
    color: var(--primary-color);
    font-weight: 600;
  }
  
  .menu-toggle {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    width: 30px;
    height: 20px;
    background: none;
    border: none;
    cursor: pointer;
    padding: 0;
  }
  
  .menu-toggle span {
    width: 100%;
    height: 2px;
    background-color: var(--text-color);
    transition: all 0.3s ease;
  }
  
  .cta-button-container {
    display: none;
  }
  
  .cta-button {
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: 6px;
    padding: 0.5rem 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s ease;
  }
  
  .cta-button:hover {
    background-color: var(--primary-color-dark);
  }
  
  /* Main Content */
  main {
    padding-top: 80px;
  }
  
  section {
    padding: 4rem 1rem;
  }
  
  .section-header {
    text-align: center;
    margin-bottom: 3rem;
    opacity: 0;
    transform: translateY(20px);
    transition: all 0.6s ease;
  }
  
  .section-header.visible {
    opacity: 1;
    transform: translateY(0);
  }
  
  .section-header h2 {
    font-size: 2.5rem;
    font-weight: 800;
    margin-bottom: 1rem;
    background: linear-gradient(90deg, var(--primary-color), var(--accent-color));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    display: inline-block;
  }
  
  .section-header p {
    font-size: 1.1rem;
    color: var(--text-light);
    max-width: 600px;
    margin: 0 auto;
  }
  
  /* Hero Section */
  .hero-section {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    justify-content: center;
    position: relative;
    overflow: hidden;
    padding: 6rem 1rem 4rem;
  }
  
  .hero-content {
    display: flex;
    flex-direction: column;
    gap: 3rem;
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .hero-text {
    max-width: 600px;
  }
  
  .hero-text h1 {
    font-size: 3rem;
    font-weight: 900;
    line-height: 1.2;
    margin-bottom: 1.5rem;
  }
  
  .hero-text p {
    font-size: 1.2rem;
    color: var(--text-light);
    margin-bottom: 2rem;
    line-height: 1.6;
  }
  
  .hero-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
  }
  
  .primary-button {
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: 8px;
    padding: 0.75rem 1.5rem;
    font-weight: 600;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.2s ease;
  }
  
  .primary-button:hover {
    background-color: var(--primary-color-dark);
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(255, 62, 0, 0.2);
  }
  
  .secondary-button {
    background-color: transparent;
    color: var(--text-color);
    border: 2px solid var(--border-color);
    border-radius: 8px;
    padding: 0.75rem 1.5rem;
    font-weight: 600;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.2s ease;
  }
  
  .secondary-button:hover {
    border-color: var(--primary-color);
    color: var(--primary-color);
    transform: translateY(-2px);
  }
  
  .hero-image-container {
    max-width: 600px;
    margin: 0 auto;
  }
  
  .hero-image {
    width: 100%;
    position: relative;
  }
  
  .browser-mockup {
    width: 100%;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
    background-color: white;
    border: 1px solid var(--border-color);
  }
  
  .browser-header {
    background-color: #f5f5f5;
    padding: 0.5rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  
  .browser-controls {
    display: flex;
    gap: 6px;
  }
  
  .browser-controls span {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background-color: #ddd;
  }
  
  .browser-controls span:nth-child(1) {
    background-color: #ff5f57;
  }
  
  .browser-controls span:nth-child(2) {
    background-color: #ffbd2e;
  }
  
  .browser-controls span:nth-child(3) {
    background-color: #28c940;
  }
  
  .browser-address-bar {
    flex: 1;
    height: 24px;
    background-color: white;
    border-radius: 4px;
    margin-left: 0.5rem;
  }
  
  .browser-content {
    padding: 1rem;
    background-color: #f9f9f9;
  }
  
  .video-player {
    width: 100%;
    aspect-ratio: 16/9;
    background-color: #000;
    position: relative;
    border-radius: 4px;
    overflow: hidden;
  }
  
  .video-content {
    width: 100%;
    height: 100%;
    background: linear-gradient(rgba(0,0,0,0.3), rgba(0,0,0,0.3)), url('/placeholder.svg?height=720&width=1280');
    background-size: cover;
    background-position: center;
  }
  
  .sponsor-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    padding: 1rem;
  }
  
  .sponsor-badge {
    align-self: flex-start;
    background-color: rgba(255, 255, 255, 0.2);
    color: white;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.75rem;
    font-weight: 600;
    backdrop-filter: blur(4px);
  }
  
  .skip-button {
    align-self: flex-end;
    background-color: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    font-size: 0.875rem;
    font-weight: 600;
    cursor: pointer;
    animation: pulse 2s infinite;
  }
  
  .video-controls {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    padding: 1rem;
    background: linear-gradient(transparent, rgba(0,0,0,0.7));
  }
  
  .progress-bar {
    width: 100%;
    height: 4px;
    background-color: rgba(255, 255, 255, 0.3);
    border-radius: 2px;
    position: relative;
    margin-bottom: 0.5rem;
  }
  
  .progress-filled {
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    width: 35%;
    background-color: var(--primary-color);
    border-radius: 2px;
  }
  
  .sponsor-segment {
    position: absolute;
    top: 0;
    left: 40%;
    height: 100%;
    width: 15%;
    background-color: var(--accent-color);
    border-radius: 2px;
    opacity: 0.8;
  }
  
  .control-buttons {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  .play-button, .volume-button, .settings-button, .fullscreen-button {
    width: 12px;
    height: 12px;
    background-color: white;
    border-radius: 50%;
  }
  
  .time {
    color: white;
    font-size: 0.75rem;
    margin-left: auto;
    margin-right: 1rem;
  }
  
  .scroll-indicator {
    position: absolute;
    bottom: 2rem;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
  }
  
  .mouse {
    width: 30px;
    height: 50px;
    border: 2px solid var(--text-light);
    border-radius: 15px;
    position: relative;
  }
  
  .wheel {
    width: 4px;
    height: 8px;
    background-color: var(--text-light);
    border-radius: 2px;
    position: absolute;
    top: 10px;
    left: 50%;
    transform: translateX(-50%);
    animation: scroll 2s infinite;
  }
  
  .arrow {
    width: 10px;
    height: 10px;
    border-right: 2px solid var(--text-light);
    border-bottom: 2px solid var(--text-light);
    transform: rotate(45deg);
  }
  
  @keyframes scroll {
    0% {
      opacity: 1;
      transform: translateX(-50%) translateY(0);
    }
    100% {
      opacity: 0;
      transform: translateX(-50%) translateY(15px);
    }
  }
  
  /* Features Section */
  .features-section {
    background-color: var(--background-alt);
  }
  
  .features-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
    opacity: 0;
    transform: translateY(20px);
    transition: all 0.6s ease;
  }
  
  .features-grid.visible {
    opacity: 1;
    transform: translateY(0);
  }
  
  .feature-card {
    background-color: var(--card-background);
    border-radius: 12px;
    padding: 2rem;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
    transition: all 0.3s ease;
    border: 1px solid var(--border-color);
  }
  
  .feature-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
  }
  
  .feature-icon {
    width: 60px;
    height: 60px;
    border-radius: 12px;
    margin-bottom: 1.5rem;
    position: relative;
  }
  
  .feature-icon.detection {
    background-color: var(--primary-color);
  }
  
  .feature-icon.skip {
    background-color: var(--secondary-color);
  }
  
  .feature-icon.settings {
    background-color: var(--accent-color);
  }
  
  .feature-icon.time {
    background-color: var(--accent-color-2);
  }
  
  .feature-card h3 {
    font-size: 1.5rem;
    font-weight: 700;
    margin-bottom: 1rem;
  }
  
  .feature-card p {
    color: var(--text-light);
    line-height: 1.6;
  }
  
  /* How It Works Section */
  .how-it-works-section {
    background-color: var(--background-color);
  }
  
  .steps-container {
    display: grid;
    grid-template-columns: 1fr;
    gap: 2rem;
    max-width: 1000px;
    margin: 0 auto 4rem;
  }
  
  .step-card {
    display: flex;
    gap: 1.5rem;
    align-items: flex-start;
  }
  
  .step-number {
    font-size: 2.5rem;
    font-weight: 900;
    color: var(--primary-color);
    opacity: 0.2;
    line-height: 1;
  }
  
  .step-content {
    flex: 1;
  }
  
  .step-content h3 {
    font-size: 1.5rem;
    font-weight: 700;
    margin-bottom: 0.5rem;
  }
  
  .step-content p {
    color: var(--text-light);
    line-height: 1.6;
  }
  
  .demo-container {
    max-width: 800px;
    margin: 0 auto;
    background-color: var(--card-background);
    border-radius: 12px;
    padding: 2rem;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
    border: 1px solid var(--border-color);
  }
  
  .demo-video {
    width: 100%;
    position: relative;
  }
  
  .video-timeline {
    width: 100%;
    height: 40px;
    background-color: #f5f5f5;
    border-radius: 8px;
    position: relative;
    margin-bottom: 2rem;
  }
  
  .timeline-bar {
    position: absolute;
    top: 50%;
    left: 0;
    transform: translateY(-50%);
    width: 100%;
    height: 6px;
    background-color: #ddd;
    border-radius: 3px;
  }
  
  .sponsor-segment-marker {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    height: 6px;
    left: var(--left);
    width: var(--width);
    background-color: var(--accent-color);
    border-radius: 3px;
  }
  
  .playhead {
    position: absolute;
    top: 50%;
    left: 30%;
    transform: translate(-50%, -50%);
    width: 16px;
    height: 16px;
    background-color: var(--primary-color);
    border-radius: 50%;
    animation: playhead 8s infinite linear;
  }
  
  @keyframes playhead {
    0% {
      left: 0%;
    }
    25% {
      left: 30%;
    }
    30% {
      left: 45%;
    }
    50% {
      left: 60%;
    }
    55% {
      left: 70%;
    }
    100% {
      left: 100%;
    }
  }
  
  .skip-animation {
    width: 100%;
    height: 100px;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  
  .skip-notification {
    background-color: rgba(0, 0, 0, 0.8);
    color: white;
    padding: 0.75rem 1.5rem;
    border-radius: 8px;
    display: flex;
    align-items: center;
    gap: 1rem;
    animation: fadeInOut 8s infinite;
  }
  
  .skip-text {
    color: var(--primary-color);
    font-weight: 600;
  }
  
  @keyframes fadeInOut {
    0%, 100% {
      opacity: 0;
    }
    25%, 35% {
      opacity: 1;
    }
    55%, 65% {
      opacity: 1;
    }
    40%, 60% {
      opacity: 0;
    }
  }
  
  /* Pricing Section */
  .pricing-section {
    background-color: var(--background-alt);
  }
  
  .pricing-cards {
    display: grid;
    grid-template-columns: 1fr;
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
    opacity: 0;
    transform: translateY(20px);
    transition: all 0.6s ease;
  }
  
  .pricing-cards.visible {
    opacity: 1;
    transform: translateY(0);
  }
  
  .pricing-card {
    background-color: var(--card-background);
    border-radius: 12px;
    padding: 2rem;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
    transition: all 0.3s ease;
    border: 1px solid var(--border-color);
    position: relative;
  }
  
  .pricing-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
  }
  
  .pricing-card.premium {
    border: 2px solid var(--primary-color);
    box-shadow: 0 10px 30px rgba(255, 62, 0, 0.1);
  }
  
  .popular-badge {
    position: absolute;
    top: -12px;
    left: 50%;
    transform: translateX(-50%);
    background-color: var(--primary-color);
    color: white;
    padding: 0.25rem 1rem;
    border-radius: 20px;
    font-size: 0.875rem;
    font-weight: 600;
  }
  
  .pricing-header {
    text-align: center;
    margin-bottom: 2rem;
  }
  
  .pricing-header h3 {
    font-size: 1.75rem;
    font-weight: 700;
    margin-bottom: 1rem;
  }
  
  .price {
    font-size: 3rem;
    font-weight: 800;
    color: var(--primary-color);
    margin-bottom: 0.5rem;
  }
  
  .price span {
    font-size: 1rem;
    font-weight: 400;
    color: var(--text-light);
  }
  
  .pricing-header p {
    color: var(--text-light);
  }
  
  .pricing-features {
    list-style: none;
    margin-bottom: 2rem;
  }
  
  .pricing-features li {
    padding: 0.75rem 0;
    border-bottom: 1px solid var(--border-color);
    color: var(--text-color);
    position: relative;
    padding-left: 1.5rem;
  }
  
  .pricing-features li::before {
    content: '✓';
    position: absolute;
    left: 0;
    color: var(--primary-color);
    font-weight: 600;
  }
  
  .pricing-button {
    width: 100%;
    padding: 0.75rem;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
  }
  
  .pricing-button.primary {
    background-color: var(--primary-color);
    color: white;
    border: none;
  }
  
  .pricing-button.primary:hover {
    background-color: var(--primary-color-dark);
  }
  
  .pricing-button.secondary {
    background-color: transparent;
    color: var(--text-color);
    border: 2px solid var(--border-color);
  }
  
  .pricing-button.secondary:hover {
    border-color: var(--primary-color);
    color: var(--primary-color);
  }
  
  .comparison-link {
    text-align: center;
    margin-top: 2rem;
  }
  
  .comparison-link button {
    background: none;
    border: none;
    color: var(--primary-color);
    font-weight: 600;
    cursor: pointer;
    text-decoration: underline;
    font-size: 0.875rem;
  }
  
  .pricing-details {
    max-width: 1000px;
    margin: 2rem auto 0;
    overflow-x: auto;
  }
  
  table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 1rem;
  }
  
  th, td {
    padding: 1rem;
    text-align: left;
    border-bottom: 1px solid var(--border-color);
  }
  
  th {
    font-weight: 600;
    background-color: #f5f5f5;
  }
  
  tr:last-child td {
    border-bottom: none;
  }
  
  /* Testimonials Section */
  .testimonials-section {
    background-color: var(--background-color);
  }
  
  .testimonials-container {
    display: grid;
    grid-template-columns: 1fr;
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .testimonial-card {
    background-color: var(--card-background);
    border-radius: 12px;
    padding: 2rem;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
    border: 1px solid var(--border-color);
    display: flex;
    gap: 1.5rem;
    align-items: flex-start;
  }
  
  .testimonial-avatar {
    width: 60px;
    height: 60px;
    border-radius: 50%;
    background-color: var(--avatar-img);
    flex-shrink: 0;
  }
  
  .testimonial-content {
    flex: 1;
  }
  
  .testimonial-content p {
    font-size: 1.1rem;
    line-height: 1.6;
    margin-bottom: 1rem;
    font-style: italic;
  }
  
  .testimonial-author {
    display: flex;
    flex-direction: column;
  }
  
  .name {
    font-weight: 600;
  }
  
  .role {
    font-size: 0.875rem;
    color: var(--text-light);
  }
  
  /* CTA Section */
  .cta-section {
    background: linear-gradient(135deg, var(--primary-color), var(--accent-color-2));
    color: white;
    text-align: center;
  }
  
  .cta-content {
    max-width: 800px;
    margin: 0 auto;
  }
  
  .cta-content h2 {
    font-size: 2.5rem;
    font-weight: 800;
    margin-bottom: 1rem;
  }
  
  .cta-content p {
    font-size: 1.1rem;
    margin-bottom: 2rem;
    opacity: 0.9;
  }
  
  .cta-buttons {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 1rem;
  }
  
  .cta-section .primary-button {
    background-color: white;
    color: var(--primary-color);
  }
  
  .cta-section .primary-button:hover {
    background-color: rgba(255, 255, 255, 0.9);
    box-shadow: 0 4px 12px rgba(255, 255, 255, 0.3);
  }
  
  .cta-section .secondary-button {
    border-color: white;
    color: white;
  }
  
  .cta-section .secondary-button:hover {
    background-color: rgba(255, 255, 255, 0.1);
    border-color: white;
    color: white;
  }
  
  /* Footer */
  footer {
    background-color: var(--background-alt);
    padding: 4rem 1rem 2rem;
  }
  
  .footer-content {
    max-width: 1200px;
    margin: 0 auto;
    display: grid;
    grid-template-columns: 1fr;
    gap: 3rem;
  }
  
  .footer-logo {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  
  .footer-links {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 2rem;
  }
  
  .footer-column h4 {
    font-size: 1.1rem;
    font-weight: 600;
    margin-bottom: 1.5rem;
  }
  
  .footer-column ul {
    list-style: none;
  }
  
  .footer-column li {
    margin-bottom: 0.75rem;
  }
  
  .footer-column a {
    color: var(--text-light);
    text-decoration: none;
    transition: color 0.2s ease;
  }
  
  .footer-column a:hover {
    color: var(--primary-color);
  }
  
  .footer-newsletter {
    max-width: 400px;
  }
  
  .footer-newsletter h4 {
    font-size: 1.1rem;
    font-weight: 600;
    margin-bottom: 1rem;
  }
  
  .footer-newsletter p {
    color: var(--text-light);
    margin-bottom: 1.5rem;
  }
  
  .newsletter-form {
    display: flex;
    gap: 0.5rem;
  }
  
  .newsletter-form input {
    flex: 1;
    padding: 0.75rem;
    border: 1px solid var(--border-color);
    border-radius: 6px;
    font-size: 0.875rem;
  }
  
  .newsletter-form button {
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: 6px;
    padding: 0.75rem 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s ease;
  }
  
  .newsletter-form button:hover {
    background-color: var(--primary-color-dark);
  }
  
  .footer-bottom {
    max-width: 1200px;
    margin: 3rem auto 0;
    padding-top: 2rem;
    border-top: 1px solid var(--border-color);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }
  
  .footer-bottom p {
    font-size: 0.875rem;
    color: var(--text-light);
  }
  
  .social-links {
    display: flex;
    gap: 1rem;
  }
  
  .social-icon {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background-color: #f5f5f5;
    transition: all 0.2s ease;
  }
  
  .social-icon:hover {
    background-color: var(--primary-color);
    transform: translateY(-3px);
  }
  
  /* Responsive Styles */
  @media (min-width: 640px) {
    .hero-text h1 {
      font-size: 3.5rem;
    }
    
    .testimonial-card {
      padding: 2.5rem;
    }
  }
  
  @media (min-width: 768px) {
    section {
      padding: 5rem 2rem;
    }
    
    .features-grid {
      grid-template-columns: repeat(2, 1fr);
    }
    
    .pricing-cards {
      grid-template-columns: repeat(2, 1fr);
    }
    
    .pricing-card.lifetime {
      grid-column: span 2;
      max-width: 500px;
      margin: 0 auto;
    }
    
    .testimonials-container {
      grid-template-columns: repeat(2, 1fr);
    }
    
    .footer-content {
      grid-template-columns: 1fr 2fr;
    }
    
    .footer-bottom {
      flex-direction: row;
      justify-content: space-between;
    }
  }
  
  @media (min-width: 1024px) {
    header {
      padding: 0 4rem;
    }
    
    .menu-toggle {
      display: none;
    }
    
    nav {
      position: static;
      transform: none;
      opacity: 1;
      visibility: visible;
      background-color: transparent;
      box-shadow: none;
      width: auto;
    }
    
    nav ul {
      display: flex;
      gap: 2rem;
      padding: 0;
    }
    
    nav li {
      margin: 0;
    }
    
    nav button {
      padding: 0;
      width: auto;
    }
    
    .cta-button-container {
      display: block;
    }
    
    .hero-content {
      flex-direction: row;
      align-items: center;
    }
    
    .hero-text {
      flex: 1;
    }
    
    .hero-image-container {
      flex: 1;
    }
    
    .features-grid {
      grid-template-columns: repeat(4, 1fr);
    }
    
    .steps-container {
      grid-template-columns: repeat(2, 1fr);
    }
    
    .pricing-cards {
      grid-template-columns: repeat(3, 1fr);
    }
    
    .pricing-card.lifetime {
      grid-column: auto;
      max-width: none;
    }
    
    .testimonials-container {
      grid-template-columns: repeat(3, 1fr);
    }
    
    .footer-content {
      grid-template-columns: 1fr 2fr 1fr;
    }
  }
  
  @media (min-width: 1280px) {
    .hero-text h1 {
      font-size: 4rem;
    }
    
    .section-header h2 {
      font-size: 3rem;
    }
  }
</style>