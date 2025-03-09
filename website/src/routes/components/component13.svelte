<script lang="ts">
  import { onMount } from 'svelte';
  import { fade, fly, scale } from 'svelte/transition';
  import { elasticOut, cubicOut } from 'svelte/easing';
  
  // State using runes
  let isMenuOpen = $state(false);
  let activeSection = $state('hero');
  let hasScrolled = $state(false);
  let mousePosition = $state({ x: 0, y: 0 });
  
  // Pricing toggle
  let isYearly = $state(false);
  
  // Derived values
  let priceFree = $derived(0);
  let pricePro = $derived(isYearly ? 29.99 : 3.99);
  let savingsText = $derived(isYearly ? 'Save 37%' : '');
  
  // Animation states
  const animationComplete = $state({
    hero: false,
    features: false,
    pricing: false,
    cta: false
  });
  
  // Track mouse position for interactive elements
  function handleMouseMove(event: MouseEvent) {
    mousePosition.x = event.clientX;
    mousePosition.y = event.clientY;
  }
  
  // Handle scroll to update active section
  function handleScroll() {
    hasScrolled = window.scrollY > 50;
    
    const sections = document.querySelectorAll('section[id]');
    const scrollPosition = window.scrollY + 200;
    
    sections.forEach(section => {
      const sectionTop = (section as HTMLElement).offsetTop;
      const sectionHeight = (section as HTMLElement).offsetHeight;
      const sectionId = section.getAttribute('id') || '';
      
      if (scrollPosition >= sectionTop && scrollPosition < sectionTop + sectionHeight) {
        activeSection = sectionId;
        
        // Mark section as viewed for animations
        if (!animationComplete[sectionId as keyof typeof animationComplete]) {
          animationComplete[sectionId as keyof typeof animationComplete] = true;
        }
      }
    });
  }
  
  // Initialize
  onMount(() => {
    window.addEventListener('scroll', handleScroll);
    window.addEventListener('mousemove', handleMouseMove);
    
    // Trigger initial animations
    setTimeout(() => {
      animationComplete.hero = true;
    }, 300);
    
    return () => {
      window.removeEventListener('scroll', handleScroll);
      window.removeEventListener('mousemove', handleMouseMove);
    };
  });
  
  // Scroll to section
  function scrollToSection(id: string) {
    const element = document.getElementById(id);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
    }
    if (isMenuOpen) {
      isMenuOpen = false;
    }
  }
</script>

<svelte:head>
  <title>SponsorSkipper - Skip YouTube Sponsorships Automatically</title>
  <meta name="description" content="Save time watching YouTube videos by automatically skipping sponsorship segments">
</svelte:head>

<svelte:window on:scroll={handleScroll} on:mousemove={handleMouseMove} />

<!-- Abstract geometric shapes background -->
<div class="geometric-shapes">
  <div class="shape shape-1"></div>
  <div class="shape shape-2"></div>
  <div class="shape shape-3"></div>
  <div class="shape shape-4"></div>
  <div class="shape shape-5"></div>
  <div class="shape shape-grid"></div>
</div>

<!-- Navigation -->
<header class:scrolled={hasScrolled}>
  <div class="logo">
    <div class="logo-icon">
      <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M5 5L19 19M12 12L19 5M12 12L5 19" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </div>
    <span>SPONSOR<span class="highlight">SKIPPER</span></span>
  </div>
  
  <nav class:open={isMenuOpen}>
    <ul>
      <li><button class:active={activeSection === 'hero'} on:click={() => scrollToSection('hero')}>Home</button></li>
      <li><button class:active={activeSection === 'features'} on:click={() => scrollToSection('features')}>Features</button></li>
      <li><button class:active={activeSection === 'pricing'} on:click={() => scrollToSection('pricing')}>Pricing</button></li>
      <li><button class:active={activeSection === 'cta'} on:click={() => scrollToSection('cta')}>Download</button></li>
    </ul>
  </nav>
  
  <button class="cta-button small">Get Extension</button>
  
  <button class="menu-toggle" on:click={() => isMenuOpen = !isMenuOpen}>
    <span></span>
    <span></span>
    <span></span>
  </button>
</header>

<main>
  <!-- Hero Section -->
  <section id="hero" class="hero-section">
    <div class="container">
      <div class="hero-content">
        {#if animationComplete.hero}
          <div class="hero-text" in:fly={{ y: 50, duration: 800, delay: 200, easing: cubicOut }}>
            <div class="badge" in:scale={{ start: 0.8, duration: 500, delay: 300, easing: elasticOut }}>
              CHROME EXTENSION
            </div>
            <h1>
              <span class="line" in:fly={{ y: 30, duration: 600, delay: 400, easing: cubicOut }}>SKIP THE</span>
              <span class="line highlight" in:fly={{ y: 30, duration: 600, delay: 600, easing: cubicOut }}>BORING PARTS</span>
              <span class="line" in:fly={{ y: 30, duration: 600, delay: 800, easing: cubicOut }}>OF YOUTUBE</span>
            </h1>
            <p in:fly={{ y: 20, duration: 600, delay: 1000, easing: cubicOut }}>
              Automatically skip sponsorships, intros, and outros in YouTube videos.
              Save time and enjoy uninterrupted content.
            </p>
            
            <div class="hero-buttons" in:fly={{ y: 20, duration: 600, delay: 1200, easing: cubicOut }}>
              <button class="cta-button primary">Install Free</button>
              <button class="cta-button secondary" on:click={() => scrollToSection('pricing')}>
                Go Premium
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M9 18L15 12L9 6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </button>
            </div>
          </div>
        {/if}
        
        <div class="hero-image">
          {#if animationComplete.hero}
            <div class="browser-mockup" in:fly={{ x: 50, duration: 800, delay: 500, easing: cubicOut }}>
              <div class="browser-header">
                <div class="browser-dots">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
                <div class="browser-address">youtube.com</div>
              </div>
              <div class="browser-content">
                <div class="video-player">
                  <div class="video-placeholder"></div>
                  <div class="sponsor-overlay" in:fly={{ y: 20, duration: 400, delay: 1500, easing: cubicOut }}>
                    <div class="sponsor-badge">
                      <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M5 3L19 12L5 21V3Z" fill="currentColor"/>
                      </svg>
                      SPONSORSHIP DETECTED
                    </div>
                    <div class="sponsor-progress">
                      <div class="sponsor-bar"></div>
                    </div>
                    <button class="sponsor-skip">
                      SKIPPING...
                      <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M13 5L21 12L13 19M3 5L11 12L3 19" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                      </svg>
                    </button>
                  </div>
                </div>
                <div class="video-info">
                  <div class="video-title">How to Build a Website in 2023 (Complete Tutorial)</div>
                  <div class="video-stats">
                    <div class="time-saved">
                      <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2"/>
                        <path d="M12 7V12L15 15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                      </svg>
                      47s saved
                    </div>
                  </div>
                </div>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>
    
    <div class="memphis-pattern bottom-pattern"></div>
  </section>
  
  <!-- Features Section -->
  <section id="features" class="features-section">
    <div class="container">
      <div class="section-header">
        <h2>
          <span class="highlight">AWESOME</span> FEATURES
        </h2>
        <p>Our extension is packed with powerful features to enhance your YouTube experience.</p>
      </div>
      
      <div class="features-grid">
        {#if animationComplete.features}
          <div class="feature-card" in:fly={{ y: 30, duration: 600, delay: 200, easing: cubicOut }}>
            <div class="feature-icon red">
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M13 5L21 12L13 19M3 5L11 12L3 19" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <h3>Auto-Skip Sponsorships</h3>
            <p>Automatically detects and skips sponsored segments in videos so you don't have to manually skip them.</p>
          </div>
          
          <div class="feature-card" in:fly={{ y: 30, duration: 600, delay: 300, easing: cubicOut }}>
            <div class="feature-icon yellow">
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M13 2L3 14H12L11 22L21 10H12L13 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <h3>Lightning Fast</h3>
            <p>Minimal impact on performance. Works silently in the background without slowing down your browsing.</p>
          </div>
          
          <div class="feature-card" in:fly={{ y: 30, duration: 600, delay: 400, easing: cubicOut }}>
            <div class="feature-icon blue">
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2"/>
                <path d="M12 7V12L15 15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </div>
            <h3>Time Saved Tracker</h3>
            <p>See exactly how much time you've saved by skipping sponsorships across all your watched videos.</p>
          </div>
        {/if}
      </div>
      
      {#if animationComplete.features}
        <div class="feature-highlight" in:fly={{ y: 30, duration: 600, delay: 500, easing: cubicOut }}>
          <div class="highlight-content">
            <h3>Smart Detection Technology</h3>
            <p>Our advanced algorithm recognizes sponsorship segments with incredible accuracy, even when creators try to disguise them.</p>
            <ul class="feature-list">
              <li>
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M5 12L10 17L20 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                Recognizes common sponsorship phrases
              </li>
              <li>
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M5 12L10 17L20 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                Detects visual sponsorship indicators
              </li>
              <li>
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M5 12L10 17L20 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                Learns from user feedback
              </li>
              <li>
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M5 12L10 17L20 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                Updates in real-time
              </li>
            </ul>
          </div>
          <div class="highlight-image">
            <div class="image-container">
              <div class="image-placeholder"></div>
            </div>
          </div>
        </div>
      {/if}
    </div>
    
    <div class="memphis-pattern right-pattern"></div>
  </section>
  
  <!-- Pricing Section -->
  <section id="pricing" class="pricing-section">
    <div class="container">
      <div class="section-header">
        <h2>
          CHOOSE YOUR <span class="highlight">PLAN</span>
        </h2>
        <p>Upgrade to Premium for unlimited skips and advanced features.</p>
      </div>
      
      {#if animationComplete.pricing}
        <div class="pricing-toggle" in:fade={{ duration: 400, delay: 200 }}>
          <span class:active={!isYearly}>Monthly</span>
          <label class="toggle">
            <input type="checkbox" bind:checked={isYearly}>
            <span class="slider"></span>
          </label>
          <span class:active={isYearly}>Yearly</span>
          {#if savingsText}
            <span class="savings">{savingsText}</span>
          {/if}
        </div>
        
        <div class="pricing-cards">
          <div class="pricing-card" in:fly={{ y: 30, duration: 600, delay: 300, easing: cubicOut }}>
            <div class="card-header">
              <h3>Free</h3>
              <div class="price">
                <span class="amount">${priceFree}</span>
                <span class="period">forever</span>
              </div>
              <p>Basic sponsorship skipping for casual YouTube viewers</p>
            </div>
            <div class="card-features">
              <ul>
                <li>
                  <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M5 12L10 17L20 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  Skip up to 50 sponsorships per month
                </li>
                <li>
                  <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M5 12L10 17L20 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  Basic sponsorship detection
                </li>
                <li>
                  <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M5 12L10 17L20 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  Time saved tracker
                </li>
                <li>
                  <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M5 12L10 17L20 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  Works on all YouTube videos
                </li>
              </ul>
            </div>
            <div class="card-footer">
              <button class="cta-button outline">Install Now</button>
            </div>
          </div>
          
          <div class="pricing-card premium" in:fly={{ y: 30, duration: 600, delay: 400, easing: cubicOut }}>
            <div class="popular-badge">POPULAR</div>
            <div class="card-header">
              <h3>Premium</h3>
              <div class="price">
                <span class="amount">${pricePro}</span>
                <span class="period">/{isYearly ? 'year' : 'month'}</span>
              </div>
              <p>Unlimited skipping and advanced features for power users</p>
            </div>
            <div class="card-features">
              <ul>
                <li>
                  <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M5 12L10 17L20 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  <strong>Unlimited</strong> sponsorship skipping
                </li>
                <li>
                  <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M5 12L10 17L20 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  Advanced detection algorithm
                </li>
                <li>
                  <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M5 12L10 17L20 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  Custom skip rules and preferences
                </li>
                <li>
                  <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M5 12L10 17L20 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  Skip intros, outros & reminders
                </li>
                <li>
                  <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M5 12L10 17L20 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  Detailed analytics dashboard
                </li>
                <li>
                  <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M5 12L10 17L20 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  Priority support
                </li>
              </ul>
            </div>
            <div class="card-footer">
              <button class="cta-button primary">Go Premium</button>
            </div>
          </div>
        </div>
        
        <div class="pricing-note" in:fade={{ duration: 400, delay: 600 }}>
          <p>Not convinced yet? Try Premium free for 7 days!</p>
          <button class="cta-button gradient">Start Free Trial</button>
        </div>
      {/if}
    </div>
    
    <div class="memphis-pattern left-pattern"></div>
  </section>
  
  <!-- CTA Section -->
  <section id="cta" class="cta-section">
    <div class="container">
      {#if animationComplete.cta}
        <div class="cta-box" in:fly={{ y: 30, duration: 600, delay: 200, easing: cubicOut }}>
          <h2>
            STOP WASTING <span class="highlight">TIME</span>
          </h2>
          <p>The average YouTube user wastes over 5 hours per month watching sponsorships. Get that time back with SponsorSkipper!</p>
          
          <div class="cta-buttons">
            <button class="cta-button primary large">Install Free</button>
            <button class="cta-button secondary large">Go Premium</button>
          </div>
          
          <div class="ratings">
            <div class="stars">
              {#each Array(5) as _, i}
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z" fill="currentColor" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              {/each}
            </div>
            <span>4.9/5 from 2,000+ reviews</span>
          </div>
        </div>
      {/if}
    </div>
  </section>
</main>

<footer>
  <div class="container">
    <div class="footer-content">
      <div class="footer-logo">
        <div class="logo-icon">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M5 5L19 19M12 12L19 5M12 12L5 19" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <span>SPONSOR<span class="highlight">SKIPPER</span></span>
      </div>
      
      <div class="footer-links">
        <div class="link-group">
          <h4>Product</h4>
          <ul>
            <li><a href="#features">Features</a></li>
            <li><a href="#pricing">Pricing</a></li>
            <li><a href="#cta">Download</a></li>
            <li><a href="#">Chrome Web Store</a></li>
          </ul>
        </div>
        
        <div class="link-group">
          <h4>Company</h4>
          <ul>
            <li><a href="#">About Us</a></li>
            <li><a href="#">Blog</a></li>
            <li><a href="#">Careers</a></li>
            <li><a href="#">Contact</a></li>
          </ul>
        </div>
        
        <div class="link-group">
          <h4>Legal</h4>
          <ul>
            <li><a href="#">Terms of Service</a></li>
            <li><a href="#">Privacy Policy</a></li>
            <li><a href="#">Cookie Policy</a></li>
          </ul>
        </div>
      </div>
    </div>
    
    <div class="footer-bottom">
      <p>&copy; {new Date().getFullYear()} SponsorSkipper. All rights reserved.</p>
      <div class="social-links">
        <a href="#" aria-label="Twitter">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M23 3.01006C22.0424 3.68553 20.9821 4.20217 19.86 4.54006C19.2577 3.84757 18.4573 3.35675 17.567 3.13398C16.6767 2.91122 15.7395 2.96725 14.8821 3.29451C14.0247 3.62177 13.2884 4.20446 12.773 4.96377C12.2575 5.72309 11.9877 6.62239 12 7.54006V8.54006C10.2426 8.58562 8.50127 8.19587 6.93101 7.4055C5.36074 6.61513 4.01032 5.44869 3 4.01006C3 4.01006 -1 13.0101 8 17.0101C5.94053 18.408 3.48716 19.109 1 19.0101C10 24.0101 21 19.0101 21 7.51006C20.9991 7.23151 20.9723 6.95365 20.92 6.68006C21.9406 5.67355 22.6608 4.40277 23 3.01006Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </a>
        <a href="#" aria-label="GitHub">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M16 22.0268V19.1568C16.0375 18.68 15.9731 18.2006 15.811 17.7506C15.6489 17.3006 15.3929 16.8902 15.06 16.5468C18.2 16.1968 21.5 15.0068 21.5 9.54679C21.4997 8.15062 20.9627 6.80799 20 5.79679C20.4558 4.5753 20.4236 3.22514 19.91 2.02679C19.91 2.02679 18.73 1.67679 16 3.50679C13.708 2.88561 11.292 2.88561 9 3.50679C6.27 1.67679 5.09 2.02679 5.09 2.02679C4.57638 3.22514 4.54424 4.5753 5 5.79679C4.03013 6.81549 3.49252 8.17026 3.5 9.57679C3.5 14.9968 6.8 16.1868 9.94 16.5768C9.611 16.9168 9.35726 17.3222 9.19531 17.7667C9.03335 18.2112 8.96681 18.6849 9 19.1568V22.0268" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M9 20.0267C6 20.9999 3.5 20.0267 2 17.0267" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </a>
        <a href="#" aria-label="Discord">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M9 11C9 11.5523 8.55228 12 8 12C7.44772 12 7 11.5523 7 11C7 10.4477 7.44772 10 8 10C8.55228 10 9 10.4477 9 11Z" fill="currentColor"/>
            <path d="M16 12C16.5523 12 17 11.5523 17 11C17 10.4477 16.5523 10 16 10C15.4477 10 15 10.4477 15 11C15 11.5523 15.4477 12 16 12Z" fill="currentColor"/>
            <path d="M9.09 17.3C9.51 17.7 10.12 17.95 11 17.95H13C13.88 17.95 14.49 17.7 14.91 17.3" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M18.2 4.8C19.9673 5.17284 21.3027 5.90686 22 8.5C23 12.5 22 20 18.5 20.5C15.5 21 12.5 21 9.5 20.5C6 20 5 12.5 6 8.5C6.69732 5.90686 8.03274 5.17284 9.8 4.8C10.2 4.8 10.7 5.1 11 5.5L11.4 6.1C12.8 6 14.2 6 15.6 6.1L16 5.5C16.3 5.1 16.8 4.8 17.2 4.8H18.2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </a>
      </div>
    </div>
  </div>
</footer>

<style>
  /* Global Styles */
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    background-color: #0f0f0f;
    color: #ffffff;
    overflow-x: hidden;
  }
  
  :global(*) {
    box-sizing: border-box;
  }
  
  .container {
    width: 100%;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 1.5rem;
    position: relative;
    z-index: 2;
  }
  
  .highlight {
    color: #ff3e3e;
  }
  
  /* Abstract Geometric Shapes */
  .geometric-shapes {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 0;
    overflow: hidden;
    pointer-events: none;
  }
  
  .shape {
    position: absolute;
    opacity: 0.15;
  }
  
  .shape-1 {
    top: 10%;
    left: 5%;
    width: 300px;
    height: 300px;
    border-radius: 50%;
    background: linear-gradient(45deg, #ff3e3e, #ff9d00);
    animation: float 20s ease-in-out infinite;
  }
  
  .shape-2 {
    top: 40%;
    right: 10%;
    width: 400px;
    height: 400px;
    border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%;
    background: linear-gradient(45deg, #4c00ff, #00b8ff);
    animation: float 25s ease-in-out infinite reverse;
  }
  
  .shape-3 {
    bottom: 10%;
    left: 20%;
    width: 250px;
    height: 250px;
    background: linear-gradient(45deg, #00ff88, #00b8ff);
    transform: rotate(45deg);
    animation: float 18s ease-in-out infinite;
  }
  
  .shape-4 {
    top: 60%;
    left: 60%;
    width: 150px;
    height: 150px;
    border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%;
    background: linear-gradient(45deg, #ff3e3e, #ff00aa);
    animation: float 15s ease-in-out infinite;
  }
  
  .shape-5 {
    top: 20%;
    right: 30%;
    width: 100px;
    height: 100px;
    background: #ffcc00;
    border-radius: 50%;
    animation: float 12s ease-in-out infinite;
  }
  
  .shape-grid {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-image: 
      linear-gradient(to right, rgba(255, 255, 255, 0.05) 1px, transparent 1px),
      linear-gradient(to bottom, rgba(255, 255, 255, 0.05) 1px, transparent 1px);
    background-size: 50px 50px;
  }
  
  @keyframes float {
    0% {
      transform: translate(0, 0) rotate(0deg);
    }
    50% {
      transform: translate(30px, 20px) rotate(5deg);
    }
    100% {
      transform: translate(0, 0) rotate(0deg);
    }
  }
  
  /* Memphis Design Patterns */
  .memphis-pattern {
    position: absolute;
    z-index: 1;
    pointer-events: none;
  }
  
  .bottom-pattern {
    bottom: 0;
    left: 0;
    width: 100%;
    height: 80px;
    background-image: 
      radial-gradient(circle at 10% 50%, #ff3e3e 10px, transparent 10px),
      radial-gradient(circle at 30% 50%, #4c00ff 15px, transparent 15px),
      radial-gradient(circle at 50% 50%, #ffcc00 8px, transparent 8px),
      radial-gradient(circle at 70% 50%, #00ff88 12px, transparent 12px),
      radial-gradient(circle at 90% 50%, #ff00aa 10px, transparent 10px);
  }
  
  .right-pattern {
    top: 20%;
    right: 0;
    width: 60px;
    height: 60%;
    background-image: 
      linear-gradient(45deg, #ff3e3e 25%, transparent 25%),
      linear-gradient(-45deg, #ff3e3e 25%, transparent 25%),
      linear-gradient(45deg, transparent 75%, #ff3e3e 75%),
      linear-gradient(-45deg, transparent 75%, #ff3e3e 75%);
    background-size: 20px 20px;
    background-position: 0 0, 0 10px, 10px -10px, -10px 0px;
  }
  
  .left-pattern {
    top: 10%;
    left: 0;
    width: 40px;
    height: 80%;
    background-image: 
      linear-gradient(to right, #ffcc00, #ffcc00 10px, transparent 10px, transparent 20px),
      linear-gradient(to bottom, #ffcc00, #ffcc00 10px, transparent 10px, transparent 20px);
    background-size: 20px 20px;
  }
  
  /* Header */
  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    z-index: 100;
    transition: all 0.3s ease;
  }
  
  header.scrolled {
    background-color: rgba(15, 15, 15, 0.9);
    backdrop-filter: blur(10px);
    box-shadow: 0 4px 30px rgba(0, 0, 0, 0.1);
    padding: 1rem 1.5rem;
  }
  
  .logo {
    display: flex;
    align-items: center;
    font-weight: 800;
    font-size: 1.5rem;
    letter-spacing: -0.03em;
  }
  
  .logo-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 2.5rem;
    height: 2.5rem;
    background-color: #ff3e3e;
    border-radius: 8px;
    margin-right: 0.75rem;
    color: white;
  }
  
  nav {
    display: flex;
  }
  
  nav ul {
    display: flex;
    list-style: none;
    margin: 0;
    padding: 0;
    gap: 2rem;
  }
  
  nav button {
    background: none;
    border: none;
    color: #ffffff;
    font-size: 1rem;
    cursor: pointer;
    padding: 0.5rem;
    position: relative;
    transition: color 0.2s ease;
  }
  
  nav button:hover, nav button.active {
    color: #ff3e3e;
  }
  
  nav button.active::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0.5rem;
    right: 0.5rem;
    height: 2px;
    background-color: #ff3e3e;
  }
  
  .menu-toggle {
    display: none;
    flex-direction: column;
    justify-content: space-between;
    width: 2rem;
    height: 1.5rem;
    background: none;
    border: none;
    cursor: pointer;
    padding: 0;
  }
  
  .menu-toggle span {
    display: block;
    width: 100%;
    height: 3px;
    background-color: #ffffff;
    border-radius: 3px;
    transition: all 0.3s ease;
  }
  
  /* CTA Buttons */
  .cta-button {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0.75rem 1.5rem;
    border: 3px solid #000000;
    border-radius: 8px;
    font-weight: 700;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.2s ease;
    text-transform: uppercase;
    box-shadow: 4px 4px 0 #000000;
    position: relative;
    overflow: hidden;
  }
  
  .cta-button svg {
    width: 1.25rem;
    height: 1.25rem;
    margin-left: 0.5rem;
  }
  
  .cta-button.primary {
    background-color: #ff3e3e;
    color: #ffffff;
  }
  
  .cta-button.secondary {
    background-color: #4c00ff;
    color: #ffffff;
  }
  
  .cta-button.outline {
    background-color: transparent;
    color: #ffffff;
    border-color: #ffffff;
  }
  
  .cta-button.gradient {
    background: linear-gradient(45deg, #ff3e3e, #ff9d00);
    color: #ffffff;
  }
  
  .cta-button.large {
    padding: 1rem 2rem;
    font-size: 1.125rem;
  }
  
  .cta-button.small {
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
    box-shadow: 3px 3px 0 #000000;
  }
  
  .cta-button:hover {
    transform: translate(2px, 2px);
    box-shadow: 2px 2px 0 #000000;
  }
  
  .cta-button:active {
    transform: translate(4px, 4px);
    box-shadow: 0 0 0 #000000;
  }
  
  /* Hero Section */
  .hero-section {
    min-height: 100vh;
    display: flex;
    align-items: center;
    padding: 8rem 0 4rem;
    position: relative;
    overflow: hidden;
  }
  
  .hero-content {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
    align-items: center;
  }
  
  .hero-text {
    max-width: 600px;
  }
  
  .badge {
    display: inline-block;
    background-color: #ffcc00;
    color: #000000;
    font-weight: 700;
    font-size: 0.875rem;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    margin-bottom: 1.5rem;
    border: 2px solid #000000;
  }
  
  h1 {
    font-size: 4rem;
    line-height: 1.1;
    margin: 0 0 1.5rem;
    font-weight: 900;
    letter-spacing: -0.03em;
    text-transform: uppercase;
  }
  
  h1 .line {
    display: block;
  }
  
  .hero-text p {
    font-size: 1.25rem;
    line-height: 1.6;
    margin-bottom: 2rem;
    color: rgba(255, 255, 255, 0.8);
  }
  
  .hero-buttons {
    display: flex;
    gap: 1rem;
  }
  
  .hero-image {
    position: relative;
  }
  
  .browser-mockup {
    background-color: #1a1a1a;
    border-radius: 8px;
    border: 3px solid #000000;
    box-shadow: 8px 8px 0 #000000;
    overflow: hidden;
  }
  
  .browser-header {
    display: flex;
    align-items: center;
    background-color: #2a2a2a;
    padding: 0.75rem;
    border-bottom: 3px solid #000000;
  }
  
  .browser-dots {
    display: flex;
    gap: 0.5rem;
    margin-right: 1rem;
  }
  
  .browser-dots span {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background-color: #ff3e3e;
  }
  
  .browser-dots span:nth-child(2) {
    background-color: #ffcc00;
  }
  
  .browser-dots span:nth-child(3) {
    background-color: #00ff88;
  }
  
  .browser-address {
    background-color: #3a3a3a;
    padding: 0.25rem 1rem;
    border-radius: 4px;
    font-size: 0.875rem;
    color: rgba(255, 255, 255, 0.7);
    flex-grow: 1;
  }
  
  .browser-content {
    padding: 1rem;
  }
  
  .video-player {
    position: relative;
    border-radius: 4px;
    overflow: hidden;
    margin-bottom: 1rem;
  }
  
  .video-placeholder {
    width: 100%;
    padding-top: 56.25%; /* 16:9 aspect ratio */
    background-color: #2a2a2a;
    background-image: url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="rgba(255,255,255,0.2)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polygon points="10 8 16 12 10 16 10 8"/></svg>');
    background-repeat: no-repeat;
    background-position: center;
    background-size: 48px;
  }
  
  .sponsor-overlay {
    position: absolute;
    bottom: 1rem;
    left: 1rem;
    right: 1rem;
    background-color: rgba(0, 0, 0, 0.8);
    border: 2px solid #ff3e3e;
    border-radius: 4px;
    padding: 0.75rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }
  
  .sponsor-badge {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-weight: 700;
    font-size: 0.875rem;
    color: #ff3e3e;
  }
  
  .sponsor-badge svg {
    width: 1rem;
    height: 1rem;
  }
  
  .sponsor-progress {
    width: 100%;
    height: 4px;
    background-color: rgba(255, 255, 255, 0.2);
    border-radius: 2px;
    overflow: hidden;
  }
  
  .sponsor-bar {
    height: 100%;
    width: 70%;
    background-color: #ff3e3e;
    animation: progress 2s linear infinite;
  }
  
  @keyframes progress {
    0% {
      width: 0%;
    }
    100% {
      width: 100%;
    }
  }
  
  .sponsor-skip {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    background-color: #ff3e3e;
    color: #ffffff;
    border: none;
    border-radius: 4px;
    padding: 0.5rem 1rem;
    font-weight: 700;
    font-size: 0.75rem;
    cursor: pointer;
  }
  
  .sponsor-skip svg {
    width: 1rem;
    height: 1rem;
  }
  
  .video-info {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .video-title {
    font-weight: 600;
    font-size: 1rem;
  }
  
  .video-stats {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  .time-saved {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    font-size: 0.875rem;
    color: #00ff88;
  }
  
  .time-saved svg {
    width: 1rem;
    height: 1rem;
  }
  
  /* Features Section */
  .features-section {
    padding: 6rem 0;
    position: relative;
    background-color: #151515;
  }
  
  .section-header {
    text-align: center;
    max-width: 800px;
    margin: 0 auto 4rem;
  }
  
  h2 {
    font-size: 3rem;
    font-weight: 900;
    margin: 0 0 1rem;
    letter-spacing: -0.03em;
    text-transform: uppercase;
  }
  
  .section-header p {
    font-size: 1.25rem;
    color: rgba(255, 255, 255, 0.8);
  }
  
  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
    margin-bottom: 4rem;
  }
  
  .feature-card {
    background-color: #1a1a1a;
    border: 3px solid #000000;
    border-radius: 8px;
    padding: 2rem;
    box-shadow: 8px 8px 0 #000000;
    transition: transform 0.3s ease;
  }
  
  .feature-card:hover {
    transform: translateY(-5px);
  }
  
  .feature-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 4rem;
    height: 4rem;
    border-radius: 8px;
    margin-bottom: 1.5rem;
    color: #000000;
    border: 3px solid #000000;
  }
  
  .feature-icon svg {
    width: 2rem;
    height: 2rem;
  }
  
  .feature-icon.red {
    background-color: #ff3e3e;
  }
  
  .feature-icon.yellow {
    background-color: #ffcc00;
  }
  
  .feature-icon.blue {
    background-color: #4c00ff;
  }
  
  h3 {
    font-size: 1.5rem;
    font-weight: 700;
    margin: 0 0 1rem;
  }
  
  .feature-card p {
    color: rgba(255, 255, 255, 0.8);
    line-height: 1.6;
  }
  
  .feature-highlight {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 3rem;
    background-color: #1a1a1a;
    border: 3px solid #000000;
    border-radius: 8px;
    padding: 3rem;
    box-shadow: 8px 8px 0 #000000;
    position: relative;
    overflow: hidden;
  }
  
  .feature-highlight::before {
    content: '';
    position: absolute;
    top: 0;
    right: 0;
    width: 100px;
    height: 100px;
    background-color: #ff3e3e;
    transform: translate(50%, -50%) rotate(45deg);
  }
  
  .highlight-content {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .feature-list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .feature-list li {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    color: rgba(255, 255, 255, 0.8);
  }
  
  .feature-list li svg {
    width: 1.25rem;
    height: 1.25rem;
    color: #00ff88;
    flex-shrink: 0;
  }
  
  .highlight-image {
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .image-container {
    position: relative;
    width: 100%;
    max-width: 400px;
  }
  
  .image-placeholder {
    width: 100%;
    padding-top: 75%;
    background-color: #2a2a2a;
    border: 3px solid #000000;
    border-radius: 8px;
  }
  
  .image-container::after {
    content: '';
    position: absolute;
    top: 1rem;
    left: 1rem;
    width: 100%;
    height: 100%;
    background-color: #4c00ff;
    border: 3px solid #000000;
    border-radius: 8px;
    z-index: -1;
  }
  
  /* Pricing Section */
  .pricing-section {
    padding: 6rem 0;
    position: relative;
  }
  
  .pricing-toggle {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    margin-bottom: 3rem;
  }
  
  .pricing-toggle span {
    font-size: 1rem;
    color: rgba(255, 255, 255, 0.7);
    transition: color 0.2s ease;
  }
  
  .pricing-toggle span.active {
    color: #ffffff;
    font-weight: 600;
  }
  
  .toggle {
    position: relative;
    display: inline-block;
    width: 60px;
    height: 30px;
  }
  
  .toggle input {
    opacity: 0;
    width: 0;
    height: 0;
  }
  
  .slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: #2a2a2a;
    transition: .4s;
    border-radius: 34px;
    border: 2px solid #000000;
  }
  
  .slider:before {
    position: absolute;
    content: "";
    height: 22px;
    width: 22px;
    left: 2px;
    bottom: 2px;
    background-color: #ffffff;
    transition: .4s;
    border-radius: 50%;
    border: 2px solid #000000;
  }
  
  input:checked + .slider {
    background-color: #4c00ff;
  }
  
  input:checked + .slider:before {
    transform: translateX(30px);
  }
  
  .savings {
    background-color: #ffcc00;
    color: #000000;
    font-weight: 700;
    font-size: 0.75rem;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    border: 2px solid #000000;
  }
  
  .pricing-cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
    margin-bottom: 3rem;
  }
  
  .pricing-card {
    background-color: #1a1a1a;
    border: 3px solid #000000;
    border-radius: 8px;
    padding: 2rem;
    box-shadow: 8px 8px 0 #000000;
    display: flex;
    flex-direction: column;
    position: relative;
  }
  
  .pricing-card.premium {
    background-color: #2a2a2a;
    transform: translateY(-1rem);
  }
  
  .popular-badge {
    position: absolute;
    top: 0;
    right: 2rem;
    background-color: #ffcc00;
    color: #000000;
    font-weight: 700;
    font-size: 0.875rem;
    padding: 0.5rem 1rem;
    border: 3px solid #000000;
    border-top: none;
    border-radius: 0 0 8px 8px;
  }
  
  .card-header {
    margin-bottom: 2rem;
  }
  
  .price {
    display: flex;
    align-items: baseline;
    margin: 1rem 0;
  }
  
  .amount {
    font-size: 3rem;
    font-weight: 800;
    line-height: 1;
  }
  
  .period {
    font-size: 1rem;
    color: rgba(255, 255, 255, 0.7);
    margin-left: 0.25rem;
  }
  
  .card-header p {
    color: rgba(255, 255, 255, 0.8);
  }
  
  .card-features {
    flex-grow: 1;
    margin-bottom: 2rem;
  }
  
  .card-features ul {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .card-features li {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    color: rgba(255, 255, 255, 0.8);
  }
  
  .card-features li svg {
    width: 1.25rem;
    height: 1.25rem;
    color: #00ff88;
    flex-shrink: 0;
  }
  
  .card-footer {
    margin-top: auto;
  }
  
  .card-footer .cta-button {
    width: 100%;
  }
  
  .pricing-note {
    text-align: center;
    margin-top: 3rem;
  }
  
  .pricing-note p {
    font-size: 1.25rem;
    margin-bottom: 1.5rem;
  }
  
  /* CTA Section */
  .cta-section {
    padding: 6rem 0;
    background-color: #151515;
    position: relative;
    overflow: hidden;
  }
  
  .cta-box {
    background-color: #1a1a1a;
    border: 3px solid #000000;
    border-radius: 8px;
    padding: 3rem;
    box-shadow: 8px 8px 0 #000000;
    text-align: center;
    max-width: 800px;
    margin: 0 auto;
    position: relative;
  }
  
  .cta-box::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100px;
    height: 100px;
    background-color: #ffcc00;
    transform: translate(-50%, -50%) rotate(45deg);
  }
  
  .cta-box::after {
    content: '';
    position: absolute;
    bottom: 0;
    right: 0;
    width: 100px;
    height: 100px;
    background-color: #4c00ff;
    transform: translate(50%, 50%) rotate(45deg);
  }
  
  .cta-box h2 {
    font-size: 3.5rem;
    margin-bottom: 1.5rem;
  }
  
  .cta-box p {
    font-size: 1.25rem;
    color: rgba(255, 255, 255, 0.8);
    margin-bottom: 2rem;
    max-width: 600px;
    margin-left: auto;
    margin-right: auto;
  }
  
  .cta-buttons {
    display: flex;
    justify-content: center;
    gap: 1.5rem;
    margin-bottom: 2rem;
  }
  
  .ratings {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.75rem;
  }
  
  .stars {
    display: flex;
    gap: 0.25rem;
  }
  
  .stars svg {
    width: 1.25rem;
    height: 1.25rem;
    color: #ffcc00;
  }
  
  /* Footer */
  footer {
    background-color: #0f0f0f;
    border-top: 3px solid #2a2a2a;
    padding: 4rem 0 2rem;
  }
  
  .footer-content {
    display: grid;
    grid-template-columns: 1fr 2fr;
    gap: 4rem;
    margin-bottom: 3rem;
  }
  
  .footer-logo {
    display: flex;
    align-items: center;
    font-weight: 800;
    font-size: 1.5rem;
    letter-spacing: -0.03em;
    margin-bottom: 1rem;
  }
  
  .footer-links {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 2rem;
  }
  
  .link-group h4 {
    font-size: 1.125rem;
    font-weight: 700;
    margin: 0 0 1.5rem;
  }
  
  .link-group ul {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }
  
  .link-group a {
    color: rgba(255, 255, 255, 0.7);
    text-decoration: none;
    transition: color 0.2s ease;
  }
  
  .link-group a:hover {
    color: #ffffff;
  }
  
  .footer-bottom {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 2rem;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  .footer-bottom p {
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.875rem;
  }
  
  .social-links {
    display: flex;
    gap: 1rem;
  }
  
  .social-links a {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 2.5rem;
    height: 2.5rem;
    background-color: #1a1a1a;
    border-radius: 50%;
    color: rgba(255, 255, 255, 0.7);
    transition: all 0.2s ease;
  }
  
  .social-links a:hover {
    background-color: #ff3e3e;
    color: #ffffff;
    transform: translateY(-3px);
  }
  
  .social-links svg {
    width: 1.25rem;
    height: 1.25rem;
  }
  
  /* Responsive Styles */
  @media (max-width: 1024px) {
    h1 {
      font-size: 3.5rem;
    }
    
    h2 {
      font-size: 2.5rem;
    }
    
    .hero-content {
      grid-template-columns: 1fr;
    }
    
    .hero-text {
      max-width: 100%;
      text-align: center;
      margin-bottom: 2rem;
    }
    
    .hero-buttons {
      justify-content: center;
    }
    
    .feature-highlight {
      grid-template-columns: 1fr;
      padding: 2rem;
    }
    
    .footer-content {
      grid-template-columns: 1fr;
      gap: 2rem;
    }
  }
  
  @media (max-width: 768px) {
    header {
      padding: 1rem;
    }
    
    nav {
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100vh;
      background-color: #0f0f0f;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      z-index: 99;
      transform: translateX(-100%);
      transition: transform 0.3s ease;
    }
    
    nav.open {
      transform: translateX(0);
    }
    
    nav ul {
      flex-direction: column;
      align-items: center;
    }
    
    nav button {
      font-size: 1.5rem;
      padding: 1rem;
    }
    
    .menu-toggle {
      display: flex;
      z-index: 100;
    }
    
    .menu-toggle span:nth-child(1) {
      transform-origin: top left;
    }
    
    .menu-toggle span:nth-child(3) {
      transform-origin: bottom left;
    }
    
    nav.open + .menu-toggle span:nth-child(1) {
      transform: rotate(45deg) translate(5px, -5px);
    }
    
    nav.open + .menu-toggle span:nth-child(2) {
      opacity: 0;
    }
    
    nav.open + .menu-toggle span:nth-child(3) {
      transform: rotate(-45deg) translate(5px, 5px);
    }
    
    .cta-button.small {
      display: none;
    }
    
    h1 {
      font-size: 2.5rem;
    }
    
    h2 {
      font-size: 2rem;
    }
    
    .pricing-cards {
      grid-template-columns: 1fr;
    }
    
    .pricing-card.premium {
      transform: none;
    }
    
    .cta-buttons {
      flex-direction: column;
      gap: 1rem;
    }
    
    .footer-links {
      grid-template-columns: 1fr;
    }
    
    .footer-bottom {
      flex-direction: column;
      gap: 1.5rem;
    }
  }
</style>