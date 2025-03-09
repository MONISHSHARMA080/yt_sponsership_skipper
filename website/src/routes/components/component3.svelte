<script>
// @ts-nocheck

  import { onMount } from 'svelte';
  
  // Animation states
  let isIntersecting = {};
  let videoTime = 0;
  let isPremiumHovered = false;
  
  // Features for comparison
  const features = [
    { name: "Skip Sponsorships", free: true, premium: true },
    { name: "Skip Intros", free: true, premium: true },
    { name: "Skip End Cards", free: true, premium: true },
    { name: "Unlimited Videos", free: false, premium: true },
    { name: "Custom Skip Rules", free: false, premium: true },
    { name: "Priority Updates", free: false, premium: true }
  ];
  
  // Testimonials
  const testimonials = [
    { name: "Alex K.", text: "This extension saves me hours every week. Premium is worth every penny!" },
    { name: "Jamie T.", text: "I was skeptical at first, but now I can't watch YouTube without it." },
    { name: "Sam R.", text: "The premium tier unlocked so much more value. Highly recommend upgrading." }
  ];
  
  // Set up intersection observers for animations
  onMount(() => {
    const sections = document.querySelectorAll('.animate-on-scroll');
    
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        // @ts-ignore
        isIntersecting[entry.target.id] = entry.isIntersecting;
        isIntersecting = {...isIntersecting};
      });
    }, { threshold: 0.2 });
    
    sections.forEach(section => {
      observer.observe(section);
    });
    
    // Simulate video progress for the demo
    const videoInterval = setInterval(() => {
      videoTime += 1;
      if (videoTime > 100) videoTime = 0;
      
      // Simulate sponsor segment
      if (videoTime > 30 && videoTime < 45) {
        // @ts-ignore
        document.querySelector('.sponsor-overlay').style.opacity = '1';
      } else {
        // @ts-ignore
        document.querySelector('.sponsor-overlay').style.opacity = '0';
      }
    }, 100);
    
    return () => {
      clearInterval(videoInterval);
      sections.forEach(section => observer.unobserve(section));
    };
  });
  
  // Handle premium hover
  // @ts-ignore
  function handlePremiumHover(hovering) {
    isPremiumHovered = hovering;
  }
</script>

<main>
  <!-- Hero Section -->
  <section class="hero">
    <div class="hero-content">
      <div class="hero-text">
        <h1>Skip the Boring Parts</h1>
        <p>Watch YouTube videos without the sponsorships, intros, and outros.</p>
        <div class="cta-buttons">
          <a href="#" class="cta-button primary">Install Free</a>
          <a href="#premium" class="cta-button secondary">Go Premium</a>
        </div>
      </div>
      
      <div class="hero-visual">
        <div class="video-demo">
          <div class="video-player">
            <div class="video-content"></div>
            <div class="sponsor-overlay">
              <div class="skip-button">
                <span>SPONSOR</span>
                <button>SKIP →</button>
              </div>
            </div>
            <div class="progress-bar">
              <div class="progress" style="width: {videoTime}%"></div>
              <div class="sponsor-segment" style="left: 30%; width: 15%;"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div class="hero-shapes">
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
    </div>
  </section>
  
  <!-- How It Works -->
  <section id="how-it-works" class="how-it-works animate-on-scroll">
    <h2 class={isIntersecting['how-it-works'] ? 'visible' : ''}>How It Works</h2>
    
    <div class="steps">
      <div class="step">
        <div class="step-number">1</div>
        <h3>Install the Extension</h3>
        <p>Add our extension to Chrome with just one click.</p>
      </div>
      
      <div class="step">
        <div class="step-number">2</div>
        <h3>Watch YouTube</h3>
        <p>Continue watching YouTube as you normally would.</p>
      </div>
      
      <div class="step">
        <div class="step-number">3</div>
        <h3>Skip Automatically</h3>
        <p>Our AI detects and skips sponsorships automatically.</p>
      </div>
    </div>
  </section>
  
  <!-- Features Comparison -->
  <section id="premium" class="premium animate-on-scroll">
    <h2 class={isIntersecting['premium'] ? 'visible' : ''}>Choose Your Plan</h2>
    
    <div class="plans">
      <div class="plan free">
        <h3>Free</h3>
        <p class="price">$0</p>
        <p class="description">Basic sponsorship skipping with limits</p>
        <a href="#" class="plan-cta">Install Now</a>
        <div class="plan-features">
          {#each features as feature}
            <div class="feature">
              <span>{feature.name}</span>
              <span class="feature-available">{feature.free ? '✓' : '✗'}</span>
            </div>
          {/each}
        </div>
      </div>
      
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div 
        class="plan premium" 
        class:hovered={isPremiumHovered}
        on:mouseenter={() => handlePremiumHover(true)}
        on:mouseleave={() => handlePremiumHover(false)}
      >
        <div class="premium-badge">RECOMMENDED</div>
        <h3>Premium</h3>
        <p class="price">$4.99<span>/month</span></p>
        <p class="description">Unlimited skipping and advanced features</p>
        <a href="#" class="plan-cta premium-cta">Upgrade Now</a>
        <div class="plan-features">
          {#each features as feature}
            <div class="feature">
              <span>{feature.name}</span>
              <span class="feature-available premium">{feature.premium ? '✓' : '✗'}</span>
            </div>
          {/each}
        </div>
      </div>
    </div>
  </section>
  
  <!-- Testimonials -->
  <section id="testimonials" class="testimonials animate-on-scroll">
    <h2 class={isIntersecting['testimonials'] ? 'visible' : ''}>What Users Say</h2>
    
    <div class="testimonials-grid">
      {#each testimonials as testimonial, i}
        <div class="testimonial" style="animation-delay: {i * 0.2}s">
          <p>"{testimonial.text}"</p>
          <div class="testimonial-author">— {testimonial.name}</div>
        </div>
      {/each}
    </div>
  </section>
  
  <!-- FAQ -->
  <section id="faq" class="faq animate-on-scroll">
    <h2 class={isIntersecting['faq'] ? 'visible' : ''}>Frequently Asked Questions</h2>
    
    <div class="faq-items">
      <div class="faq-item">
        <h3>How does the extension detect sponsorships?</h3>
        <p>Our advanced AI analyzes video content and audio to identify sponsored segments with high accuracy.</p>
      </div>
      
      <div class="faq-item">
        <h3>What's the difference between free and premium?</h3>
        <p>The free version limits the number of videos you can skip sponsorships on each month. Premium gives you unlimited skips and additional features.</p>
      </div>
      
      <div class="faq-item">
        <h3>Can I customize what gets skipped?</h3>
        <p>Premium users can create custom rules for what content to skip and what to keep.</p>
      </div>
    </div>
  </section>
  
  <!-- Final CTA -->
  <section class="final-cta">
    <div class="cta-content">
      <h2>Ready to Enhance Your YouTube Experience?</h2>
      <div class="cta-buttons">
        <a href="#" class="cta-button primary">Install Free</a>
        <a href="#premium" class="cta-button secondary">Go Premium</a>
      </div>
    </div>
  </section>
  
  <!-- Footer -->
  <footer>
    <div class="footer-content">
      <div class="footer-logo">
        <span>SponsorSkipper</span>
      </div>
      
      <div class="footer-links">
        <a href="#">Privacy Policy</a>
        <a href="#">Terms of Service</a>
        <a href="#">Contact</a>
      </div>
      
      <div class="footer-copyright">
        © {new Date().getFullYear()} SponsorSkipper. All rights reserved.
      </div>
    </div>
  </footer>
</main>

<style>
  /* Global Styles */
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
    background-color: #f8f9fa;
    color: #212529;
    overflow-x: hidden;
  }
  
  main {
    width: 100%;
    overflow-x: hidden;
  }
  
  h1, h2, h3 {
    margin: 0;
    font-weight: 800;
  }
  
  h1 {
    font-size: 4rem;
    line-height: 1.1;
    margin-bottom: 1rem;
  }
  
  h2 {
    font-size: 2.5rem;
    margin-bottom: 2rem;
    opacity: 0;
    transform: translateY(30px);
    transition: opacity 0.6s ease, transform 0.6s ease;
  }
  
  h2.visible {
    opacity: 1;
    transform: translateY(0);
  }
  
  section {
    padding: 5rem 2rem;
    position: relative;
  }
  
  /* Hero Section */
  .hero {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    overflow: hidden;
    background-color: #0f172a;
    color: white;
  }
  
  .hero-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-between;
    max-width: 1200px;
    width: 100%;
    z-index: 2;
    padding: 2rem;
    gap: 4rem;
  }
  
  @media (min-width: 768px) {
    .hero-content {
      flex-direction: row;
    }
  }
  
  .hero-text {
    flex: 1;
    max-width: 600px;
  }
  
  .hero-text p {
    font-size: 1.25rem;
    margin-bottom: 2rem;
    opacity: 0.9;
  }
  
  .hero-visual {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  
  .video-demo {
    width: 100%;
    max-width: 500px;
    aspect-ratio: 16 / 9;
    background-color: #1e293b;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
    position: relative;
  }
  
  .video-player {
    width: 100%;
    height: 100%;
    position: relative;
    overflow: hidden;
  }
  
  .video-content {
    width: 100%;
    height: 100%;
    background-image: url('/placeholder.svg?height=720&width=1280');
    background-size: cover;
    background-position: center;
  }
  
  .sponsor-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.3s ease;
  }
  
  .skip-button {
    background-color: rgba(255, 255, 255, 0.2);
    padding: 0.75rem 1.5rem;
    border-radius: 8px;
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  .skip-button span {
    font-weight: 600;
    color: #ff4d4d;
  }
  
  .skip-button button {
    background-color: #ff4d4d;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s ease;
  }
  
  .skip-button button:hover {
    background-color: #ff3333;
  }
  
  .progress-bar {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 5px;
    background-color: rgba(255, 255, 255, 0.2);
  }
  
  .progress {
    height: 100%;
    background-color: #ff4d4d;
    transition: width 0.1s linear;
  }
  
  .sponsor-segment {
    position: absolute;
    height: 100%;
    background-color: rgba(255, 77, 77, 0.5);
    top: 0;
  }
  
  .hero-shapes {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    overflow: hidden;
    z-index: 1;
  }
  
  .shape {
    position: absolute;
    border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%;
    opacity: 0.1;
  }
  
  .shape-1 {
    background-color: #ff4d4d;
    width: 500px;
    height: 500px;
    top: -100px;
    right: -100px;
    animation: float 15s ease-in-out infinite;
  }
  
  .shape-2 {
    background-color: #4d79ff;
    width: 600px;
    height: 600px;
    bottom: -200px;
    left: -200px;
    animation: float 20s ease-in-out infinite reverse;
  }
  
  .shape-3 {
    background-color: #ffcc4d;
    width: 300px;
    height: 300px;
    top: 50%;
    right: 10%;
    animation: float 18s ease-in-out infinite 2s;
  }
  
  @keyframes float {
    0% { transform: translate(0, 0) rotate(0deg); }
    50% { transform: translate(50px, 50px) rotate(180deg); }
    100% { transform: translate(0, 0) rotate(360deg); }
  }
  
  /* CTA Buttons */
  .cta-buttons {
    display: flex;
    gap: 1rem;
    flex-wrap: wrap;
  }
  
  .cta-button {
    display: inline-block;
    padding: 0.75rem 1.5rem;
    border-radius: 8px;
    font-weight: 600;
    text-decoration: none;
    transition: all 0.2s ease;
    text-align: center;
  }
  
  .cta-button.primary {
    background-color: #ff4d4d;
    color: white;
    box-shadow: 0 4px 6px -1px rgba(255, 77, 77, 0.2), 0 2px 4px -1px rgba(255, 77, 77, 0.1);
  }
  
  .cta-button.primary:hover {
    background-color: #ff3333;
    transform: translateY(-2px);
    box-shadow: 0 10px 15px -3px rgba(255, 77, 77, 0.2), 0 4px 6px -2px rgba(255, 77, 77, 0.1);
  }
  
  .cta-button.secondary {
    background-color: transparent;
    color: white;
    border: 2px solid white;
  }
  
  .cta-button.secondary:hover {
    background-color: rgba(255, 255, 255, 0.1);
    transform: translateY(-2px);
  }
  
  /* How It Works */
  .how-it-works {
    background-color: white;
    text-align: center;
  }
  
  .steps {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 2rem;
    margin-top: 3rem;
  }
  
  .step {
    flex: 1;
    min-width: 250px;
    max-width: 350px;
    padding: 2rem;
    background-color: #f8f9fa;
    border-radius: 12px;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
  }
  
  .step:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  }
  
  .step-number {
    width: 50px;
    height: 50px;
    background-color: #ff4d4d;
    color: white;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.5rem;
    font-weight: 700;
    margin: 0 auto 1.5rem;
  }
  
  /* Premium Section */
  .premium {
    background-color: #f0f4f8;
    text-align: center;
  }
  
  .plans {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 2rem;
    margin-top: 3rem;
  }
  
  .plan {
    flex: 1;
    min-width: 280px;
    max-width: 400px;
    padding: 3rem 2rem;
    background-color: white;
    border-radius: 12px;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    position: relative;
  }
  
  .plan.premium {
    border: 2px solid transparent;
    transition: transform 0.3s ease, box-shadow 0.3s ease, border-color 0.3s ease;
  }
  
  .plan.premium.hovered {
    transform: translateY(-10px) scale(1.05);
    border-color: #ff4d4d;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  }
  
  .premium-badge {
    position: absolute;
    top: -12px;
    left: 50%;
    transform: translateX(-50%);
    background-color: #ff4d4d;
    color: white;
    padding: 0.25rem 1rem;
    border-radius: 20px;
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.05em;
  }
  
  .plan h3 {
    font-size: 1.5rem;
    margin-bottom: 1rem;
  }
  
  .price {
    font-size: 2.5rem;
    font-weight: 800;
    margin-bottom: 1rem;
  }
  
  .price span {
    font-size: 1rem;
    font-weight: 400;
    opacity: 0.7;
  }
  
  .description {
    color: #64748b;
    margin-bottom: 2rem;
  }
  
  .plan-cta {
    display: inline-block;
    padding: 0.75rem 2rem;
    border-radius: 8px;
    font-weight: 600;
    text-decoration: none;
    margin-bottom: 2rem;
    transition: all 0.2s ease;
  }
  
  .plan-cta {
    background-color: #e2e8f0;
    color: #334155;
  }
  
  .plan-cta:hover {
    background-color: #cbd5e1;
  }
  
  .premium-cta {
    background-color: #ff4d4d;
    color: white;
  }
  
  .premium-cta:hover {
    background-color: #ff3333;
    transform: translateY(-2px);
    box-shadow: 0 4px 6px -1px rgba(255, 77, 77, 0.2), 0 2px 4px -1px rgba(255, 77, 77, 0.1);
  }
  
  .plan-features {
    text-align: left;
  }
  
  .feature {
    display: flex;
    justify-content: space-between;
    padding: 0.75rem 0;
    border-bottom: 1px solid #e2e8f0;
  }
  
  .feature:last-child {
    border-bottom: none;
  }
  
  .feature-available {
    font-weight: 600;
  }
  
  .feature-available.premium {
    color: #ff4d4d;
  }
  
  /* Testimonials */
  .testimonials {
    background-color: white;
    text-align: center;
  }
  
  .testimonials-grid {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 2rem;
    margin-top: 3rem;
  }
  
  .testimonial {
    flex: 1;
    min-width: 280px;
    max-width: 350px;
    padding: 2rem;
    background-color: #f8f9fa;
    border-radius: 12px;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
    text-align: left;
    opacity: 0;
    animation: fadeIn 0.5s ease forwards;
  }
  
  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
  }
  
  .testimonial p {
    font-size: 1.1rem;
    line-height: 1.6;
    margin-bottom: 1.5rem;
    font-style: italic;
  }
  
  .testimonial-author {
    font-weight: 600;
    color: #64748b;
  }
  
  /* FAQ */
  .faq {
    background-color: #f0f4f8;
    text-align: center;
  }
  
  .faq-items {
    max-width: 800px;
    margin: 0 auto;
    text-align: left;
  }
  
  .faq-item {
    margin-bottom: 2rem;
    padding: 2rem;
    background-color: white;
    border-radius: 12px;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  }
  
  .faq-item h3 {
    margin-bottom: 1rem;
    font-size: 1.25rem;
  }
  
  .faq-item p {
    color: #64748b;
    line-height: 1.6;
  }
  
  /* Final CTA */
  .final-cta {
    background-color: #0f172a;
    color: white;
    text-align: center;
    padding: 5rem 2rem;
  }
  
  .cta-content {
    max-width: 800px;
    margin: 0 auto;
  }
  
  .cta-content h2 {
    font-size: 2.5rem;
    margin-bottom: 2rem;
    opacity: 1;
    transform: none;
  }
  
  .final-cta .cta-buttons {
    justify-content: center;
    margin-top: 2rem;
  }
  
  /* Footer */
  footer {
    background-color: #1e293b;
    color: white;
    padding: 3rem 2rem;
  }
  
  .footer-content {
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2rem;
  }
  
  .footer-logo {
    font-size: 1.5rem;
    font-weight: 700;
  }
  
  .footer-links {
    display: flex;
    gap: 2rem;
    flex-wrap: wrap;
    justify-content: center;
  }
  
  .footer-links a {
    color: #cbd5e1;
    text-decoration: none;
    transition: color 0.2s ease;
  }
  
  .footer-links a:hover {
    color: white;
  }
  
  .footer-copyright {
    color: #64748b;
    font-size: 0.875rem;
  }
  
  /* Responsive Adjustments */
  @media (max-width: 768px) {
    h1 {
      font-size: 2.5rem;
    }
    
    h2 {
      font-size: 2rem;
    }
    
    .hero-content {
      padding: 1rem;
    }
    
    .cta-buttons {
      flex-direction: column;
      width: 100%;
    }
    
    .cta-button {
      width: 100%;
    }
  }
</style>